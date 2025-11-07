package a2disk

import (
	"fmt"

	"github.com/pevans/erc/a2/a2enc"
	"github.com/pevans/erc/memory"
)

// The VTOC is a "volume table of contents" that can be found on a floppy
// disk. It's a packed binary struct located in track 17.
type VTOC struct {
	FirstCatalogSectorTrackNumber  uint8
	FirstCatalogSectorSectorNumber uint8
	ReleaseNumberOfDOS             uint8
	DisketteVolume                 uint8
	MaxTrackSectorPairs            uint8
	LastTrackAllocated             uint8
	DirectionOfAllocation          int8
	TracksPerDiskette              uint8
	SectorsPerTrack                uint8
	BytesPerSector                 uint16
	FreeSectors                    map[int]string

	// NOTE: The VTOC can technically hold bitmaps of additional tracks beyond
	// 35, if your disk has such a thing. This was not typical in Apple
	// software.
}

func (v *VTOC) Parse(seg *memory.Segment) error {
	offset := a2enc.LogTrackLen * 17

	v.FirstCatalogSectorTrackNumber = seg.Get(offset + 0x1)
	v.FirstCatalogSectorSectorNumber = seg.Get(offset + 0x2)
	v.ReleaseNumberOfDOS = seg.Get(offset + 0x3)
	v.DisketteVolume = seg.Get(offset + 0x6)
	v.MaxTrackSectorPairs = seg.Get(offset + 0x27)
	v.LastTrackAllocated = seg.Get(offset + 0x30)

	// This is intended to be a positive or negative number, so we want to
	// keep the sign intact.
	v.DirectionOfAllocation = int8(seg.Get(offset + 0x31))

	v.TracksPerDiskette = seg.Get(offset + 0x34)
	v.SectorsPerTrack = seg.Get(offset + 0x35)

	// This number is meant to be stored as 16-bit
	v.BytesPerSector = (uint16(seg.Get(offset+0x37)) << 8) |
		uint16(seg.Get(offset+0x36))

	v.FreeSectors = make(map[int]string)
	for i := 0x38; i < 0xC4; i += 0x4 {
		bitmap1 := seg.Get(offset + i)
		bitmap2 := seg.Get(offset + i + 1)

		v.FreeSectors[i-0x38] = freeSectors(bitmap1, bitmap2)
	}

	return nil
}

func (v *VTOC) Valid() bool {
	// This is a really peculiar set of criteria for a "valid" VTOC. It was
	// chosen because:
	// - $FE is essentially always the "volume" of a disk when reading tracks
	// - 122 is the number for 256 byte sectors
	//
	// The idea is, if we see other values here, this is likely to be a disk
	// that happens to use track 17 for other kinds of data. It may not be
	// corrupted data, but it's not a VTOC.
	return v.DisketteVolume == 0xFE && v.MaxTrackSectorPairs == 122
}

// If every sector were free, we'd show the template below.
const freeSectorTemplate = "FEDCBA98 76543210"

func freeSectors(bitmap1, bitmap2 uint8) string {
	asBinary := fmt.Sprintf("%08b %08b", bitmap1, bitmap2)
	asSectors := ""

	for i := 0; i < len(asBinary); i++ {
		if asBinary[i] == '1' || asBinary[i] == ' ' {
			asSectors += string(freeSectorTemplate[i])
		} else {
			asSectors += "."
		}
	}

	return asSectors
}
