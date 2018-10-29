package a2

import (
	"os"

	"github.com/pevans/erc/pkg/mach"
)

const (
	// VolumeMarker is a hard-coded number that encoding software will
	// expect to see in the track padding.
	DDVolumeMarker = 0xFE

	// NumTracks is the number of tracks that can be contained on a
	// disk.
	DDNumTracks = 35

	// NumSectors is the number of sectors that each track contains.
	DDNumSectors = 16

	// LogSectorLen is the length of a logical sector, which is 256
	// bytes.
	LogSectorLen = 0x100

	// LogTrackLen is the length of a logical track, consisting of 16
	// logical sectors, which are each 256 bytes long. It thus holds 4
	// kilobytes of data.
	LogTrackLen = LogSectorLen * DDNumSectors

	// PhysSectorLen is the length of a physical sector
	PhysSectorLen = 0x18C

	// PhysTrackLen is the length of a physical track, consisting of 16
	// physical sectors.
	PhysTrackLen = PhysSectorLen * DDNumSectors

	// A sector header consists of some byte markers--all byte markers
	// in 6-and-2 encoding are 3 bytes long--and also some metadata,
	// such as the track number, the sector number, the volume, and an
	// XOR'd combination of all three. It also contains some padding.
	PhysSectorHeader = 0x14

	// The track header is 48 bytes of--well, nothing really, just
	// padding.
	PhysTrackHeader = 0x30
)

const (
	// DDMaxSteps is the maximum number of steps we can move the drive
	// head before running out of tracks on the disk. (Note that steps
	// are half of the length of a track; 35 tracks, 70 steps.)
	DDMaxSteps = 70

	// DDMaxSectorPos is the highest sector position that we can allow
	// within a given track. (0xFFF = 4k - 1.)
	DDMaxSectorPos = 0xFFF

	// DD140K is the number of bytes in 140 kilobytes.
	DD140K = 143360

	// DD140KNib is the capacity of the segment we will create for
	// nibblized data, whether from 140k logical data or just any-old
	// NIB file.
	DD140KNib = 250000
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

// NewDiskDrive returns a new disk drive ready for DOS 3.3 images.
func NewDiskDrive() *DiskDrive {
	drive := new(DiskDrive)

	drive.Mode = DDRead
	drive.ImageType = DDDOS33

	return drive
}

// Position returns the segment position that the drive is currently at,
// based upon track and sector position.
func (d *DiskDrive) Position() int {
	if d.Data == nil {
		return 0
	}

	return ((d.TrackPos / 2) * PhysTrackLen) + d.SectorPos
}

// Shift moves the sector position forward, or backward, depending on
// the sign of the given offset. If this would involve moving beyond the
// beginning or end of a track, then the sector position is instead set
// to zero.
func (d *DiskDrive) Shift(offset int) {
	if d.Locked {
		return
	}

	d.SectorPos += offset

	if d.SectorPos >= PhysTrackLen || d.SectorPos < 0 {
		d.SectorPos = 0
	}
}

// Step moves the track position forward or backward, depending on the
// sign of the offset. This simulates the stepper motor that moves the
// drive head further into the center of the disk platter (offset > 0)
// or further out (offset < 0).
func (d *DiskDrive) Step(offset int) {
	d.TrackPos += offset

	switch {
	case d.TrackPos > DDMaxSteps:
		d.TrackPos = DDMaxSteps
	case d.TrackPos < 0:
		d.TrackPos = 0
	}

	// The sector position also resets when the drive motor steps
	d.SectorPos = 0
}

func DiskPhase(addr mach.DByte) int {
	phase := -1

	switch addr & 0xF {
	case 0x1:
		phase = 1
	case 0x3:
		phase = 2
	case 0x5:
		phase = 3
	case 0x7:
		phase = 4
	}

	return phase
}

//  0  1  2  3  4     phase transition
var phaseTransitions = []int{
	0, 0, 0, 0, 0, // no phases
	0, 0, 1, 0, -1, // phase 1
	0, -1, 0, 1, 0, // phase 2
	0, 0, -1, 0, 1, // phase 3
	0, 1, 0, -1, 0, // phase 4
}

func (d *DiskDrive) StepPhase(addr mach.DByte) {
	phase := DiskPhase(addr)

	if phase < 0 || phase > 4 {
		return
	}

	offset := phaseTransitions[(d.Phase*5)+phase]
	d.Step(offset)

	d.Phase = phase
}
