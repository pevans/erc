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
package a2enc

import (
	"fmt"

	"github.com/pevans/erc/memory"
)

func Encode(imageType int, seg *memory.Segment) (*memory.Segment, error) {
	switch imageType {
	case DOS33, ProDOS:
		return Encode62(imageType, seg)

	case Nibble:
		return seg, nil
	}

	return nil, fmt.Errorf("unknown image type: %v", imageType)
}
