package a2enc

import (
	"github.com/pevans/erc/memory"
)

const (
	// SixBlock is the length of a data block that is the "six" of
	// six-and-two encoding.
	SixBlock = 0x100

	// TwoBlock is, vice-versa, the length of the "two" data block.
	TwoBlock = 0x56
)

// This is the table that holds the bytes that represent 6-and-2 encoded
// memory. Note the table goes from $00..$3F; that is the amount of values
// that six bits can hold. Each of those six-bit combinations maps to a
// different byte value that would be literally written to and read from
// the disk media. Apple II's RWTS subroutine would then translate them
// back into data that is useful to the software being run.
//
// Also, since I forget: gcr is short for "group coded recording".
//
//	00    01    02    03    04    05    06    07    08    09    0a    0b    0c    0d    0e    0f
var encGCR62 = []uint8{
	0x96, 0x97, 0x9A, 0x9B, 0x9D, 0x9E, 0x9F, 0xA6, 0xA7, 0xAB, 0xAC, 0xAD, 0xAE, 0xAF, 0xB2, 0xB3, // 00
	0xB4, 0xB5, 0xB6, 0xB7, 0xB9, 0xBA, 0xBB, 0xBC, 0xBD, 0xBE, 0xBF, 0xCB, 0xCD, 0xCE, 0xCF, 0xD3, // 10
	0xD6, 0xD7, 0xD9, 0xDA, 0xDB, 0xDC, 0xDD, 0xDE, 0xDF, 0xE5, 0xE6, 0xE7, 0xE9, 0xEA, 0xEB, 0xEC, // 20
	0xED, 0xEE, 0xEF, 0xF2, 0xF3, 0xF4, 0xF5, 0xF6, 0xF7, 0xF9, 0xFA, 0xFB, 0xFC, 0xFD, 0xFE, 0xFF, // 30
}

// An encoder is a struct which defines the pieces we need to encode
// logical data into a physical format.
type encoder struct {
	logicalSegment  *memory.Segment
	physicalSegment *memory.Segment
	imageType       int
	logicalOffset   int
	physicalOffset  int
}

func newEncoder(logSize, physSize int) *encoder {
	return &encoder{
		logicalSegment:  memory.NewSegment(logSize),
		physicalSegment: memory.NewSegment(physSize),
	}
}

// Encode62 returns a segment that is the six-and-two encoded form of the
// input segment, essentially translating from a logical to a physical
// structure.
func Encode62(imageType int, src *memory.Segment) (*memory.Segment, error) {
	enc := &encoder{
		physicalSegment: memory.NewSegment(NibSize),
		logicalSegment:  src,
		imageType:       imageType,
	}

	for track := range NumTracks {
		enc.physicalOffset = PhysTrackLen * track
		enc.writeTrack(track)
	}

	return enc.physicalSegment, nil
}

// Write will write a set of bytes into the destination segment at the
// current offset.
func (e *encoder) write(bytes []uint8) {
	for _, b := range bytes {
		e.writeByte(b)
	}
}

// writeByte simply writes a single byte into the physical segment
// without having to deal with passing around a slice
func (e *encoder) writeByte(byt uint8) {
	e.physicalSegment.Set(e.physicalOffset, byt)
	e.physicalOffset++
}

// encodeTrack will write a physically encoded track into the
// destination segment based on a logically encoded source.
func (e *encoder) writeTrack(track int) {
	logTrackOffset := LogTrackLen * track
	physTrackOffset := (PhysTrackLen * track) + len(gap1)

	// We need to write the gap1 bytes before we do anything.
	e.write(gap1)

	for sect := range NumSectors {
		logSect := LogicalSector(e.imageType, sect)

		// The logical offset is based on logTrackOffset, with the
		// sector length times the logical sector we should be copying
		e.logicalOffset = logTrackOffset + (LogSectorLen * logSect)

		// The physical offset for which we need to write will need to account
		// for the gap1 bytes that we wrote before the loop
		e.physicalOffset = physTrackOffset + (PhysSectorLen * sect)

		e.writeSector(track, sect)
	}
}

func (e *encoder) writeAddressField(track, sect int) {
	// Write the address field, starting with the prologue bytes
	e.write([]uint8{
		0xD5, 0xAA, 0x96,
	})

	// The address field consists of metadata that tells the software where to
	// organize this sector (e.g. which sector, which track)
	e.write4n4(VolumeMarker)
	e.write4n4(uint8(track))
	e.write4n4(uint8(sect))
	e.write4n4(uint8(VolumeMarker ^ track ^ sect))

	e.write([]uint8{
		// These are the epilogue of the address field, which tells the
		// software that the field is completed
		0xDE, 0xAA, 0xEB,
	})
}

func (e *encoder) writeDataField(track, sect int) {
	six := make([]uint8, SixBlock)
	two := make([]uint8, TwoBlock)

	e.write([]uint8{
		// These 3 bytes mark the beginning of the data field
		0xD5, 0xAA, 0xAD,
	})

	// Loop on the logical sector data block and build up the six-block
	// and two-block buffers
	for i := range 0x100 {
		byt := e.logicalSegment.Get(e.logicalOffset + i)

		// These are the final two bits, but their order is reversed
		rev := ((byt & 2) >> 1) | ((byt & 1) << 1)

		// The "first" six bits, which are the most significant bits
		six[i] = byt >> 2

		// And then we encode the "last" two bits by OR-ing them with
		// other two-bit segments we've seen.
		two[i%TwoBlock] |= rev << ((i / TwoBlock) * 2)
	}

	// As we write out the physical sector data block, we must XOR each
	// byte with each other byte. But the first byte is written
	// unmodified.
	e.writeByte(encGCR62[two[0]])
	for i := uint(1); i < TwoBlock; i++ {
		e.writeByte(encGCR62[two[i]^two[i-1]])
	}

	// A similar strategy is employed while writing the six-block
	// buffer, with the note that we must XOR the final two-block byte
	// with the initial six-block one.
	e.writeByte(encGCR62[six[0]^two[TwoBlock-1]])
	for i := uint(1); i < SixBlock; i++ {
		e.writeByte(encGCR62[six[i]^six[i-1]])
	}

	// We still need to write out the last byte of the six-block buffer,
	// but here, there's no need to XOR since there's no other byte.
	e.writeByte(encGCR62[six[SixBlock-1]])

	// Finally, we write the end marker
	e.write([]uint8{
		0xDE, 0xAA, 0xEB,
	})
}

// encode4n4 writes the given byte in 4-and-4 encoded form, which is
// used in sector headers.
func (e *encoder) write4n4(val uint8) {
	e.write([]uint8{
		((val >> 1) & 0x55) | 0xAA,
		(val & 0x55) | 0xAA,
	})
}

// encodeSector writes a physically encoded sector into the destination
// segment based on the logically encoded source segment.
func (e *encoder) writeSector(track, sect int) {
	e.writeAddressField(track, sect)
	e.write(gap2)
	e.writeDataField(track, sect)
	e.write(gap3)
}
