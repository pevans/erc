package a2enc

import (
	"fmt"

	"github.com/pevans/erc/memory"
)

type decodeMap map[uint8]uint8

type decoder struct {
	logicalSegment  *memory.Segment
	physicalSegment *memory.Segment
	decodeMap       decodeMap
	imageType       int
	logicalOffset   int
	physicalOffset  int
}

// addressField holds the decoded metadata from an address field
type addressField struct {
	Volume   uint8
	Track    uint8
	Sector   uint8
	Checksum uint8
}

var addressFieldPrologue = []uint8{
	0xD5, 0xAA, 0x96,
}

var addressFieldEpilogue = []uint8{
	0xDE, 0xAA, 0xEB,
}

var dataFieldPrologue = []uint8{
	0xD5, 0xAA, 0xAD,
}

var dataFieldEpilogue = []uint8{
	0xDE, 0xAA, 0xEB,
}

// Decode62 returns a new segment that is the six-and-two decoded form
// (translating a physical to a logical data structure), based on a kind
// of input image and segment.
func Decode62(imageType int, src *memory.Segment) (*memory.Segment, error) {
	dec := &decoder{
		physicalSegment: src,
		logicalSegment:  memory.NewSegment(DosSize),
		imageType:       imageType,
		decodeMap:       newDecodeMap(),
	}

	for track := 0; track < NumTracks; track++ {
		if err := dec.writeTrack(track); err != nil {
			return nil, err
		}
	}

	return dec.logicalSegment, nil
}

func (d *decoder) writeTrack(track int) error {
	logTrackOffset := LogTrackLen * track
	physTrackOffset := PhysTrackLen * track

	for sect := 0; sect < NumSectors; sect++ {
		var (
			logSect  = LogicalSector(d.imageType, sect)
			physSect = sect
		)

		// The logical offset is based on logTrackOffset, with the
		// sector length times the logical sector we should be copying
		d.logicalOffset = logTrackOffset + (LogSectorLen * logSect)

		// However, the physical offset is based on the physical sector,
		// which may need to be encoded in a different order
		d.physicalOffset = physTrackOffset + (PhysSectorLen * physSect)

		if err := d.writeSector(track, sect); err != nil {
			return err
		}
	}

	return nil
}

func newDecodeMap() decodeMap {
	m := make(decodeMap)

	for i, b := range encGCR62 {
		m[b] = uint8(i)
	}

	return m
}

func (d *decoder) logByte(b uint8) uint8 {
	lb, ok := d.decodeMap[b]
	if !ok {
		panic(fmt.Errorf("strange byte in decoding: %x", b))
	}

	return lb
}

func (d *decoder) writeByte(b uint8) {
	d.logicalSegment.Set(d.logicalOffset, b)
	d.logicalOffset++
}

func (d *decoder) readByte() uint8 {
	byt := d.physicalSegment.Get(d.physicalOffset)
	d.physicalOffset++

	return byt
}

func (d *decoder) scanForBytes(want []uint8) bool {
	maxlen := d.physicalSegment.Size() - len(want)

	// look for the prologue bytes
	for d.physicalOffset < maxlen {
		found := true

		for i, byt := range want {
			if d.physicalSegment.Get(d.physicalOffset+i) != byt {
				found = false
				break
			}
		}

		if found {
			// Move past the pattern bytes
			d.physicalOffset += len(want)
			return true
		}

		d.physicalOffset++
	}

	return false
}

// Look for the address field of a sector and return that. There are likely to
// be some self-sync padding bytes (either from gap3 or gap1), and if so those
// will be skipped past. We use the address field to confirm that the
// sector is what we think it should be.
func (d *decoder) decodeAddressField() (*addressField, error) {
	if !d.scanForBytes(addressFieldPrologue) {
		return nil, fmt.Errorf("address field prologue not found")
	}

	// Parse the 4-and-4 encoded metadata
	volume := d.decode4n4()
	track := d.decode4n4()
	sector := d.decode4n4()
	checksum := d.decode4n4()

	if !d.scanForBytes(addressFieldEpilogue) {
		return nil, fmt.Errorf("address field prologue not found")
	}

	metadata := &addressField{
		Volume:   volume,
		Track:    track,
		Sector:   sector,
		Checksum: checksum,
	}

	// There's also a checksum byte, which is the XOR of each other byte in
	// the field.
	expected := volume ^ track ^ sector
	if checksum != expected {
		return nil, fmt.Errorf(
			"address field checksum mismatch: got %02X, expected %02X",
			checksum, expected,
		)
	}

	return metadata, nil
}

// Data in the address field is 4-and-4 encoded, which essentially splits one
// byte into two. We need some way to put that byte back together.
func (d *decoder) decode4n4() uint8 {
	first := d.readByte()
	second := d.readByte()

	return ((first & 0x55) << 1) | (second & 0x55)
}

// Find and return the data field for a sector. Note that there may be
// self-sync padding bytes in front of the data field, which are called gap2
// bytes. This makes an assumption that data is 6-and-2 encoded.
func (d *decoder) decodeDataField() ([]uint8, error) {
	if !d.scanForBytes(dataFieldPrologue) {
		return nil, fmt.Errorf("data field prologue not found")
	}

	var (
		six = make([]uint8, SixBlock)
		two = make([]uint8, TwoBlock)
	)

	// The checksum begins its life as the first byte of the block. We'll use
	// this checksum to confirm that the data looks compared to a checksum
	// byte at the end by XOR-ing it with each other byte.
	checksum := d.logByte(d.readByte())
	two[0] = checksum

	for i := 1; i < TwoBlock; i++ {
		lb := d.logByte(d.readByte())
		checksum ^= lb
		two[i] = checksum
	}

	for i := 0; i < SixBlock; i++ {
		lb := d.logByte(d.readByte())
		checksum ^= lb
		six[i] = checksum
	}

	checksum ^= d.logByte(d.readByte())

	if checksum != 0 {
		return nil, fmt.Errorf("data field checksum mismatch")
	}

	if !d.scanForBytes(dataFieldEpilogue) {
		return nil, fmt.Errorf("data field epilogue not found")
	}

	// Now we need to put the 256-byte sector together from the six and two
	// blocks
	data := make([]uint8, SixBlock)
	for i := 0; i < SixBlock; i++ {
		var (
			div = i / TwoBlock
			rem = i % TwoBlock
			byt uint8
		)

		switch div {
		case 0:
			byt = ((two[rem] & 2) >> 1) | ((two[rem] & 1) << 1)
		case 1:
			byt = ((two[rem] & 8) >> 3) | ((two[rem] & 4) >> 1)
		case 2:
			byt = ((two[rem] & 0x20) >> 5) | ((two[rem] & 0x10) >> 3)
		}

		data[i] = (six[i] << 2) | byt
	}

	return data, nil
}

// Given some physical sector, write the logical form of the sector. A
// physical sector will be comprised of an address field and a data field with
// some padding between (that we don't worry about in here).
func (d *decoder) writeSector(track, sect int) error {
	addrField, err := d.decodeAddressField()
	if err != nil {
		return fmt.Errorf("track %d, sector %d: %w", track, sect, err)
	}

	if addrField.Track != uint8(track) {
		return fmt.Errorf(
			"track mismatch in address field: expected %d, got %d",
			track, addrField.Track,
		)
	}

	if addrField.Sector != uint8(sect) {
		return fmt.Errorf(
			"sector mismatch in address field: expected %d, got %d",
			sect, addrField.Sector,
		)
	}

	data, err := d.decodeDataField()
	if err != nil {
		return fmt.Errorf("track %d, sector %d: %w", track, sect, err)
	}

	for _, byt := range data {
		d.writeByte(byt)
	}

	return nil
}
