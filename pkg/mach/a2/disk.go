package a2

import (
	"os"

	"github.com/pevans/erc/pkg/mach"
)

const (
	// DDDOS33 is the image type for DOS 3.3, which is the
	// generally-used image type for Apple II DOS images. There are
	// other DOS versions, which are formatted differently, but we don't
	// handle them here.
	DDDOS33 = iota

	// DDProDOS indicates that the image type is ProDOS.
	DDProDOS

	// DDNibble is the disk type for nibble (*.NIB) disk images.
	// Although this is an image type, it's not something that actual
	// disks would have been formatted in during the Apple II era.
	DDNibble
)

// The drive mode helps us determine whether to read or write from
// the disk, but is actually unrelated to write protect mode!
const (
	// DDRead is read mode for the drive.
	DDRead = iota

	// DDWrite indicates that we are in write mode for the drive.
	DDWrite
)

// A DiskDrive represents the state of a virtual Disk II drive.
type DiskDrive struct {
	Phase        int
	Latch        mach.Byte
	TrackPos     int
	SectorPos    int
	Data         *mach.Segment
	Image        *mach.Segment
	ImageType    int
	Stream       *os.File
	Online       bool
	Mode         int
	WriteProtect bool
	Locked       bool
}

var dosSectorTable = []int{
	0x0, 0x7, 0xe, 0x6, 0xd, 0x5, 0xc, 0x4,
	0xb, 0x3, 0xa, 0x2, 0x9, 0x1, 0x8, 0xf,
}

var proSectorTable = []int{
	0x0, 0x8, 0x1, 0x9, 0x2, 0xa, 0x3, 0xb,
	0x4, 0xc, 0x5, 0xd, 0x6, 0xe, 0x7, 0xf,
}

// NewDiskDrive returns a new disk drive ready for DOS 3.3 images.
func NewDiskDrive() *DiskDrive {
	drive := new(DiskDrive)

	drive.Mode = DDRead
	drive.ImageType = DDDOS33

	return drive
}

// LogicalSector returns the logical sector number, given the current
// image type and a physical sector number (sect).
func (d *DiskDrive) LogicalSector(sect int) int {
	if sect < 0 || sect > 15 {
		return 0
	}

	switch d.ImageType {
	case DDDOS33:
		return dosSectorTable[sect]

	case DDProDOS:
		return proSectorTable[sect]
	}

	// Note: logical nibble sectors are the same as the "physical"
	// sectors.
	return sect
}
