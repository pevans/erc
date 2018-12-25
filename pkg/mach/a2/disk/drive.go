package disk

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pevans/erc/pkg/mach"
	"github.com/pkg/errors"
)

const (
	// VolumeMarker is a hard-coded number that encoding software will
	// expect to see in the track padding.
	VolumeMarker = 0xFE

	// NumTracks is the number of tracks that can be contained on a
	// disk.
	NumTracks = 35

	// NumSectors is the number of sectors that each track contains.
	NumSectors = 16

	// LogSectorLen is the length of a logical sector, which is 256
	// bytes.
	LogSectorLen = 0x100

	// LogTrackLen is the length of a logical track, consisting of 16
	// logical sectors, which are each 256 bytes long. It thus holds 4
	// kilobytes of data.
	LogTrackLen = LogSectorLen * NumSectors

	// PhysSectorLen is the length of a physical sector
	PhysSectorLen = 0x1A0

	// PhysSectorHeader is the length of a sector header
	PhysSectorHeader = 0x13

	// PhysTrackLen is the length of a physical track, consisting of 16
	// physical sectors.
	PhysTrackLen = (PhysSectorLen * NumSectors) + PhysTrackHeader

	// PhysTrackHeader is the length of a track header.
	PhysTrackHeader = 0x30
)

const (
	// MaxSteps is the maximum number of steps we can move the drive
	// head before running out of tracks on the disk. (Note that steps
	// are half of the length of a track; 35 tracks, 70 steps.)
	MaxSteps = 70

	// MaxSectorPos is the highest sector position that we can allow
	// within a given track. (0xFFF = 4k - 1.)
	MaxSectorPos = 0xFFF

	// DosSize is the number of bytes in 140 kilobytes.
	DosSize = 143360

	// NibSize is the capacity of the segment we will create for
	// nibblized data, whether from 140k logical data or just any-old
	// NIB file.
	NibSize = 234640
)

const (
	// DOS33 is the image type for DOS 3.3, which is the
	// generally-used image type for Apple II DOS images. There are
	// other DOS versions, which are formatted differently, but we don't
	// handle them here.
	DOS33 = iota

	// ProDOS indicates that the image type is ProDOS.
	ProDOS

	// Nibble is the disk type for nibble (*.NIB) disk images.
	// Although this is an image type, it's not something that actual
	// disks would have been formatted in during the Apple II era.
	Nibble
)

// The drive mode helps us determine whether to read or write from
// the disk, but is actually unrelated to write protect mode!
const (
	// ReadMode is read mode for the drive.
	ReadMode = iota

	// WriteMode indicates that we are in write mode for the drive.
	WriteMode
)

// A Drive represents the state of a virtual Disk II drive.
type Drive struct {
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

// NewDrive returns a new disk drive ready for DOS 3.3 images.
func NewDrive() *Drive {
	drive := new(Drive)

	drive.Mode = ReadMode
	drive.ImageType = DOS33

	return drive
}

// Position returns the segment position that the drive is currently at,
// based upon track and sector position.
func (d *Drive) Position() int {
	if d.Data == nil {
		return 0
	}

	return ((d.TrackPos / 2) * PhysTrackLen) + d.SectorPos
}

// Shift moves the sector position forward, or backward, depending on
// the sign of the given offset. If this would involve moving beyond the
// beginning or end of a track, then the sector position is instead set
// to zero.
func (d *Drive) Shift(offset int) {
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
func (d *Drive) Step(offset int) {
	d.TrackPos += offset

	switch {
	case d.TrackPos > MaxSteps:
		d.TrackPos = MaxSteps
	case d.TrackPos < 0:
		d.TrackPos = 0
	}

	// The sector position also resets when the drive motor steps
	d.SectorPos = 0
}

// Phase returns the motor phase based upon the given address.
func Phase(addr mach.DByte) int {
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

// StepPhase will step the drive head forward or backward based upon the
// given address (from which we decipher the motor phase).
func (d *Drive) StepPhase(addr mach.DByte) {
	phase := Phase(addr)

	if phase < 0 || phase > 4 {
		return
	}

	offset := phaseTransitions[(d.Phase*5)+phase]
	d.Step(offset)

	d.Phase = phase
}

// ImageType returns the type of image that is suggested by the suffix
// of the given filename.
func ImageType(file string) (int, error) {
	lower := strings.ToLower(file)

	switch {
	case strings.HasSuffix(lower, ".do"):
		return DOS33, nil
	case strings.HasSuffix(lower, ".nib"):
		return Nibble, nil
	case strings.HasSuffix(lower, ".po"):
		return ProDOS, nil
	}

	return -1, fmt.Errorf("Unrecognized suffix for file %s", file)
}

// Load will read a file from the filesystem and set its contents as the
// image in the drive. It also decodes the contents according to the
// (detected) image type.
func (d *Drive) Load(file string) error {
	var err error

	// See if we can figure out what type of image this is
	d.ImageType, err = ImageType(file)
	if err != nil {
		return errors.Wrapf(err, "Failed to understand image type")
	}

	// Read the bytes from the file into a buffer
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.Wrapf(err, "Failed to read file %s", file)
	}

	// Copy directly into the image segment
	d.Image = mach.NewSegment(len(bytes))
	_, err = d.Image.CopySlice(0, mach.ByteSlice(bytes))
	if err != nil {
		d.Image = nil
		return errors.Wrapf(err, "Failed to copy bytes into image segment")
	}

	// Decode into the data segment
	dec := NewDecoder(d.ImageType, d.Image)
	d.Data, err = dec.Decode()
	if err != nil {
		d.Image = nil
		return errors.Wrapf(err, "Failed to decode image")
	}

	return nil
}
