package sixtwo

import (
	"fmt"

	"github.com/pevans/erc/pkg/data"
)

type decodeMap map[data.Byte]data.Byte

type decoder struct {
	ls        *data.Segment
	ps        *data.Segment
	decMap    decodeMap
	imageType int
	loff      int
	poff      int
}

// Decode returns a new segment that is the six-and-two decoded form
// (translating a physical to a logical data structure), based on a kind
// of input image and segment.
func Decode(imageType int, src *data.Segment) (*data.Segment, error) {
	dec := &decoder{
		ps:        src,
		ls:        data.NewSegment(DosSize),
		imageType: imageType,
		decMap:    newDecodeMap(),
	}

	for track := 0; track < NumTracks; track++ {
		dec.writeTrack(track)
	}

	return dec.ls, nil
}

func (d *decoder) writeTrack(track int) {
	logTrackOffset := LogTrackLen * track
	physTrackOffset := PhysTrackLen * track

	for sect := 0; sect < NumSectors; sect++ {
		var (
			logSect  = logicalSector(d.imageType, sect)
			physSect = encPhysOrder[sect]
		)

		// The logical offset is based on logTrackOffset, with the
		// sector length times the logical sector we should be copying
		d.loff = logTrackOffset + (LogSectorLen * logSect)

		// However, the physical offset is based on the physical sector,
		// which may need to be encoded in a different order
		d.poff = physTrackOffset + (PhysSectorLen * physSect)

		d.writeSector(track, sect)
	}
}

func newDecodeMap() decodeMap {
	m := make(decodeMap)

	for i, b := range encGCR62 {
		m[b] = data.Byte(i)
	}

	return m
}

func (d *decoder) logByte(b data.Byte) data.Byte {
	lb, ok := d.decMap[b]
	if !ok {
		panic(fmt.Errorf("strange byte in decoding: %x", b))
	}

	return lb
}

func (d *decoder) writeByte(b data.Byte) {
	d.ls.Set(data.Int(d.loff), b)
	d.loff++
}

func (d *decoder) writeSector(track, sect int) {
	var (
		six = make([]data.Byte, SixBlock)
		two = make([]data.Byte, TwoBlock)
	)

	// There's going to be some opening metadata bytes that we will want
	// to skip.
	d.poff += PhysSectorHeader

	checksum := d.logByte(d.ps.Get(data.Int(uint(d.poff))))
	two[0] = checksum

	for i := uint(1); i < TwoBlock; i++ {
		lb := d.logByte(d.ps.Get(data.Int(uint(d.poff) + i)))

		checksum ^= lb
		two[i] = checksum
	}

	d.poff += int(TwoBlock)

	for i := uint(0); i < SixBlock; i++ {
		lb := d.logByte(d.ps.Get(data.Int(uint(d.poff) + i)))

		checksum ^= lb
		six[i] = checksum
	}

	d.poff += int(SixBlock)

	checksum ^= d.logByte(d.ps.Get(data.Int(uint(d.poff))))
	if checksum != 0 {
		panic(fmt.Errorf("track %d, sector %d: checksum does not match", track, sect))
	}

	for i := uint(0); i < SixBlock; i++ {
		var (
			div = i / TwoBlock
			rem = i % TwoBlock
			byt data.Byte
		)

		switch div {
		case 0:
			byt = ((two[rem] & 2) >> 1) | ((two[rem] & 1) << 1)
		case 1:
			byt = ((two[rem] & 8) >> 3) | ((two[rem] & 4) >> 1)
		case 2:
			byt = ((two[rem] & 0x20) >> 5) | ((two[rem] & 0x10) >> 3)
		}

		d.writeByte((six[i] << 2) | byt)
	}
}
