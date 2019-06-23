package sixtwo

import (
	"github.com/pevans/erc/pkg/data"
)

// This is the table that holds the bytes that represent 6-and-2 encoded
// data. Note the table goes from $00..$3F; that is the amount of values
// that six bits can hold. Each of those six-bit combinations maps to a
// different byte value that would be literally written to and read from
// the disk media. Apple II's RWTS subroutine would then translate them
// back into data that is useful to the software being run.
//
// Also, since I forget: gcr is short for "group coded recording".
//  00    01    02    03    04    05    06    07    08    09    0a    0b    0c    0d    0e    0f
var encGCR62 = []data.Byte{
	0x96, 0x97, 0x9A, 0x9B, 0x9D, 0x9E, 0x9F, 0xA6, 0xA7, 0xAB, 0xAC, 0xAD, 0xAE, 0xAF, 0xB2, 0xB3, // 00
	0xB4, 0xB5, 0xB6, 0xB7, 0xB9, 0xBA, 0xBB, 0xBC, 0xBD, 0xBE, 0xBF, 0xCB, 0xCD, 0xCE, 0xCF, 0xD3, // 10
	0xD6, 0xD7, 0xD9, 0xDA, 0xDB, 0xDC, 0xDD, 0xDE, 0xDF, 0xE5, 0xE6, 0xE7, 0xE9, 0xEA, 0xEB, 0xEC, // 20
	0xED, 0xEE, 0xEF, 0xF2, 0xF3, 0xF4, 0xF5, 0xF6, 0xF7, 0xF9, 0xFA, 0xFB, 0xFC, 0xFD, 0xFE, 0xFF, // 30
}

// Define the physical sector order in which we write encoded data
var encPhysOrder = []int{
	0x0, 0xD, 0xB, 0x9, 0x7, 0x5, 0x3, 0x1,
	0xE, 0xC, 0xA, 0x8, 0x6, 0x4, 0x2, 0xF,
}

// An encoder is a struct which defines the pieces we need to encode
// logical data into a physical format.
type encoder struct {
	ls        *data.Segment
	ps        *data.Segment
	imageType int
	loff      int
	poff      int
}

func newEncoder(logSize, physSize int) *encoder {
	return &encoder{
		ls: data.NewSegment(logSize),
		ps: data.NewSegment(physSize),
	}
}

// EncodeDOS returns a segment that is dos-encoded based on the given
// encoder struct.
func Encode(imageType int, src *data.Segment) (*data.Segment, error) {
	enc := &encoder{
		ps:        data.NewSegment(NibSize),
		ls:        src,
		imageType: imageType,
	}

	for track := 0; track < NumTracks; track++ {
		enc.writeTrack(track)
	}

	return enc.ps, nil
}

// Write will write a set of bytes into the destination segment at the
// current offset.
func (e *encoder) write(bytes []data.Byte) {
	for _, b := range bytes {
		e.writeByte(b)
	}
}

// writeByte simply writes a single byte into the physical segment
// without having to deal with passing around a slice
func (e *encoder) writeByte(byt data.Byte) {
	e.ps.Set(data.Int(e.poff), byt)
	e.poff++
}

// encodeTrack will write a physically encoded track into the
// destination segment based on a logically encoded source.
func (e *encoder) writeTrack(track int) {
	// This is the offset where we can find the logical track that we
	// are looking to write out
	logTrackOffset := LogTrackLen * track

	// Whereas, this is where we should begin writing.
	e.poff = PhysTrackLen * track

	// Write the track header
	e.write([]data.Byte{
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	})

	orig := e.poff
	for sect := 0; sect < NumSectors; sect++ {
		var (
			logSect  = logicalSector(e.imageType, sect)
			physSect = encPhysOrder[sect]
		)

		// The logical offset is based on logTrackOffset, with the
		// sector length times the logical sector we should be copying
		e.loff = logTrackOffset + (LogSectorLen * logSect)

		// However, the physical offset is based on the physical sector,
		// which may need to be encoded in a different order
		e.poff = orig + (PhysSectorLen * physSect)

		e.writeSector(track, sect)
	}
}

// encode4n4 writes the given byte in 4-and-4 encoded form, which is
// used in sector headers.
func (e *encoder) write4n4(val data.Byte) {
	e.write([]data.Byte{
		((val >> 1) & 0x55) | 0xAA,
		(val & 0x55) | 0xAA,
	})
}

// encodeSector writes a physically encoded sector into the destination
// segment based on the logically encoded source segment.
func (e *encoder) writeSector(track, sect int) {
	// Write the sector header prologue
	e.write([]data.Byte{
		0xD5, 0xAA, 0x96,
	})

	// Write the metadata
	e.write4n4(VolumeMarker)
	e.write4n4(data.Byte(track))
	e.write4n4(data.Byte(sect))
	e.write4n4(data.Byte(VolumeMarker ^ track ^ sect))

	// Write the sector header epilogue
	e.write([]data.Byte{
		0xDE, 0xAA, 0xEB,
		0xFF, 0xFF, 0xFF,
		0xFF, 0xFF,
	})

	// This is the initial preparation of data to be encoded. It's
	// written in an intermediate form, which is used to build the xor
	// table and ultimately to pass through the GCR table.
	var (
		init = make([]data.Byte, 0x156)
		xor  = make([]data.Byte, 0x157)
	)

	// This is a bit hard to explain, but the first 86 bytes (0x56) are
	// built from parts of all of the bytes in the sector.
	for i := 0; i < 0x56; i++ {
		var (
			offac = data.DByte(i + 0xAC)
			off56 = data.DByte(i + 0x56)
			vac   = e.ls.Get((data.DByte(e.loff) + offac) % 256)
			v56   = e.ls.Get(data.DByte(e.loff) + off56)
			v00   = e.ls.Get(data.DByte(e.loff + i))
			v     data.Byte
		)

		v = (v << 2) | ((vac & 0x1) << 1) | ((vac & 0x2) >> 1)
		v = (v << 2) | ((v56 & 0x1) << 1) | ((v56 & 0x2) >> 1)
		v = (v << 2) | ((v00 & 0x1) << 1) | ((v00 & 0x2) >> 1)

		init[i] = v << 2
	}

	// The last two bytes we wrote can't contain more than 6 bits of
	// 1s, so we AND them with 0x3F.
	init[0x54] &= 0x3F
	init[0x55] &= 0x3F

	// The rest of init is filled in with the src bytes unmodified. But
	// note we are writing from 0x56 onward, since we already wrote the
	// bytes before 0x56 in the block above.
	for i := 0; i < 0x100; i++ {
		init[i+0x56] = e.ls.Get(data.DByte(e.loff + i))
	}

	// Here we XOR each byte with each other byte.
	var last data.Byte
	for i := 0; i < 0x156; i++ {
		xor[i] = init[i] ^ last
		last = init[i]
	}

	// One more...
	xor[0x156] = last

	// Write out the marker that begins our sector data block
	e.write([]data.Byte{
		0xD5, 0xAA, 0xAD,
	})

	// Now we copy everything we XOR'd into the destination segment,
	// except that the bytes must be passed through the GCR table.
	for i := 0; i < 0x157; i++ {
		e.writeByte(encGCR62[xor[i]>>2])
	}

	// Finally, we write the end marker for sector data, plus 48 bytes
	// (0x30) of padding. Note the offset is doff + 0x157, given that
	// we'd just written those 0x157 bytes above.
	e.write([]data.Byte{
		0xDE, 0xAA, 0xEB,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	})
}
