// a2enc provides code to encode logically formatted disk images into a
// physical format, or vice versa, decode physically formatted images to
// a logical format.
//
// This may not make much sense to the layperson. Apple II floppy disks
// were considered prone to error, and disk drives were unable to
// distinguish intentional zero bits from errors. An encoding scheme was
// devised to ensure there would never be more than one zero bit in a
// row; this scheme is the aforementioned physical format, so named
// because it represents the format of data written to a disk.
//
// The Apple II disk operating system would expect to decode data in
// such a form to something which represents real data, like program
// code. Many of the existing disk images produced from Apple II
// software were written in a logical form. This means successful
// emulation often requires that you physically encode a logical disk
// image, so that the Apple II system can later _decode_ it.
//
// Some images are encoded in what is curiously named a "nibble" format.
// These are just physically encoded disk images that do not require any
// further encoding. The purpose for nibble-formatted images stems from
// tricks that software may use to read or store data in areas otherwise
// reserved for padding by the encoding scheme.
//
// In general, you can think of the logical data as a set of 35 tracks, each
// containing 16 sectors, all laid out one after the other in ascending order.
//
// By contrast, physical data is the same set of 35 tracks in ascending order,
// but the sectors are interleaved rather than sorted in ascending order. In
// between each track and sector are "gaps" comprised of self-sync bytes
// (typically written as 0xFF). These gaps will occur within the sector also.
// Sectors have an address field with metadata about the track and sector,
// plus a data field with the actual data.
//
// It looks somewhat like this:
//
// track N:
// [ gap1 ... ][ sector 0 ][ sector 1 ][ sector N... ]
//
// and each sector looks like:
// [ address field ][ gap2 ][ data field ][ gap3 ]
package a2enc

import (
	"fmt"

	"github.com/pevans/erc/memory"
)

// Here we have some static (as in, unchanging) variables that are used in
// physically-encoded data.

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

// self-sync bytes that are at the beginning of every track
var gap1 = []uint8{
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
}

// self-sync bytes that separate the address and data fields of a sector
var gap2 = []uint8{
	0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF,
}

// self-sync bytes that are written after every data field
var gap3 = []uint8{
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF,
}

func Encode(imageType int, seg *memory.Segment) (*memory.Segment, error) {
	switch imageType {
	case DOS33, ProDOS:
		return Encode62(imageType, seg)

	case Nibble:
		return seg, nil
	}

	return nil, fmt.Errorf("unknown image type: %v", imageType)
}
