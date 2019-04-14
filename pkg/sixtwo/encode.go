package disk

import (
	"github.com/pevans/erc/pkg/mach"
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
var encGCR62 = []mach.Byte{
	0X96, 0X97, 0X9A, 0X9B, 0X9D, 0X9E, 0X9F, 0XA6, 0XA7, 0XAB, 0XAC, 0XAD, 0XAE, 0XAF, 0XB2, 0XB3, // 00
	0XB4, 0XB5, 0XB6, 0XB7, 0XB9, 0XBA, 0XBB, 0XBC, 0XBD, 0XBE, 0XBF, 0XCB, 0XCD, 0XCE, 0XCF, 0XD3, // 10
	0XD6, 0XD7, 0XD9, 0XDA, 0XDB, 0XDC, 0XDD, 0XDE, 0XDF, 0XE5, 0XE6, 0XE7, 0XE9, 0XEA, 0XEB, 0XEC, // 20
	0XED, 0XEE, 0XEF, 0XF2, 0XF3, 0XF4, 0XF5, 0XF6, 0XF7, 0XF9, 0XFA, 0XFB, 0XFC, 0XFD, 0XFE, 0XFF, // 30
}

// Define the physical sector order in which we write encoded data
var encPhysOrder = []int{
	0x0, 0xD, 0xB, 0x9, 0x7, 0x5, 0x3, 0x1,
	0xE, 0xC, 0xA, 0x8, 0x6, 0x4, 0x2, 0xF,
}

// This is the sector table for DOS 3.3.
var dosSectorTable = []int{
	0x0, 0x7, 0xe, 0x6, 0xd, 0x5, 0xc, 0x4,
	0xb, 0x3, 0xa, 0x2, 0x9, 0x1, 0x8, 0xf,
}

// This is the sector table for ProDOS.
var proSectorTable = []int{
	0x0, 0x8, 0x1, 0x9, 0x2, 0xa, 0x3, 0xb,
	0x4, 0xc, 0x5, 0xd, 0x6, 0xe, 0x7, 0xf,
}

// An Encoder is a struct which defines the pieces we need to encode
// logical data into a physical format.
type Encoder struct {
	src       *mach.Segment
	dst       *mach.Segment
	imageType int
}

// NewEncoder returns a new encoder struct based upon a given image type
// and source segment.
func NewEncoder(imgType int, src *mach.Segment) *Encoder {
	return &Encoder{
		src:       src,
		imageType: imgType,
	}
}

// EncodeDOS returns a segment that is dos-encoded based on the given
// encoder struct.
func (e *Encoder) EncodeDOS() (*mach.Segment, error) {
	e.dst = mach.NewSegment(NibSize)
	doff := 0

	for track := 0; track < NumTracks; track++ {
		doff += e.EncodeTrack(track, doff)
	}

	return e.dst, nil
}

// EncodeNIB returns a segment that is nibble-encoded based on the given
// encoder struct.
func (e *Encoder) EncodeNIB() (*mach.Segment, error) {
	dst := mach.NewSegment(e.src.Size())
	_, err := dst.CopySlice(0, e.src.Mem)

	if err != nil {
		return nil, err
	}

	return dst, nil
}

// LogicalSector returns the logical sector number, given the current
// image type and a physical sector number (sect).
func LogicalSector(imageType, sect int) int {
	if sect < 0 || sect > 15 {
		return 0
	}

	switch imageType {
	case DOS33:
		return dosSectorTable[sect]

	case ProDOS:
		return proSectorTable[sect]
	}

	// Note: logical nibble sectors are the same as the "physical"
	// sectors.
	return sect
}

// Write will write a set of bytes into the destination segment at a
// given offset. The number of bytes written is returned.
func (e *Encoder) Write(doff int, bytes []mach.Byte) int {
	off, _ := e.dst.CopySlice(doff, bytes)
	return off
}

// EncodeTrack will write a physically encoded track into the
// destination segment based on a logically encoded source.
func (e *Encoder) EncodeTrack(track, doff int) int {
	// toff is the offset where we can find the logical track that we
	// are looking to write out
	toff := LogTrackLen * track

	// Write the track header
	doff += e.Write(doff, []mach.Byte{
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	})

	orig := doff
	for i := 0; i < NumSectors; i++ {
		logSect := LogicalSector(e.imageType, i)
		physSect := encPhysOrder[i]

		soff := toff + (LogSectorLen * logSect)
		doff := orig + (PhysSectorLen * physSect)

		_ = e.EncodeSector(track, i, doff, soff)
	}

	return PhysTrackLen
}

// Encode4n4 writes the given byte in 4-and-4 encoded form, which is
// used in sector headers.
func (e *Encoder) Encode4n4(doff int, val mach.Byte) int {
	return e.Write(doff, []mach.Byte{
		((val >> 1) & 0x55) | 0xAA,
		(val & 0x55) | 0xAA,
	})
}

// EncodeSector writes a physically encoded sector into the destination
// segment based on the logically encoded source segment.
func (e *Encoder) EncodeSector(track, sect, doff, soff int) int {
	// Write the sector header prologue
	doff += e.Write(doff, []mach.Byte{
		0xD5, 0xAA, 0x96,
	})

	// Write the metadata
	doff += e.Encode4n4(doff, VolumeMarker)
	doff += e.Encode4n4(doff, mach.Byte(track))
	doff += e.Encode4n4(doff, mach.Byte(sect))
	doff += e.Encode4n4(doff, mach.Byte(VolumeMarker^track^sect))

	// Write the sector header epilogue
	doff += e.Write(doff, []mach.Byte{
		0xDE, 0xAA, 0xEB,
		0xFF, 0xFF, 0xFF,
		0xFF, 0xFF,
	})

	// This is the initial preparation of data to be encoded. It's
	// written in an intermediate form, which is used to build the xor
	// table and ultimately to pass through the GCR table.
	init := make([]mach.Byte, 0x156)
	xor := make([]mach.Byte, 0x157)

	// This is a bit hard to explain, but the first 86 bytes (0x56) are
	// built from parts of all of the bytes in the sector.
	for i := 0; i < 0x56; i++ {
		var (
			v, vac, v56, v00 mach.Byte
			offac, off56     mach.DByte
		)

		offac = mach.DByte(i + 0xAC)
		off56 = mach.DByte(i + 0x56)

		vac = e.src.Get(mach.DByte(soff) + offac)
		v56 = e.src.Get(mach.DByte(soff) + off56)
		v00 = e.src.Get(mach.DByte(soff + i))

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
		init[i+0x56] = e.src.Get(mach.DByte(soff + i))
	}

	// Here we XOR each byte with each other byte.
	var last mach.Byte
	for i := 0; i < 0x156; i++ {
		cur := init[i]
		xor[i] = cur ^ last
		last = cur
	}

	// One more...
	xor[0x156] = last

	// Write out the marker that begins our sector data block
	doff += e.Write(doff, []mach.Byte{
		0xD5, 0xAA, 0xAD,
	})

	// Now we copy everything we XOR'd into the destination segment,
	// except that the bytes must be passed through the GCR table.
	for i := 0; i < 0x157; i++ {
		e.dst.Set(mach.DByte(doff+i), encGCR62[xor[i]>>2])
	}

	// Finally, we write the end marker for sector data, plus 48 bytes
	// (0x30) of padding. Note the offset is doff + 0x157, given that
	// we'd just written those 0x157 bytes above.
	_ = e.Write(doff+0x157, []mach.Byte{
		0xDE, 0xAA, 0xEB,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	})

	// We return the physical sector length, since that is invariably
	// what we've written.
	return PhysSectorLen
}
