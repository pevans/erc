package a2

import (
	"fmt"
	"io"
	"math/rand/v2"
	"os"
	"strings"

	"github.com/pevans/erc/a2/a2enc"
	"github.com/pevans/erc/memory"

	"github.com/pkg/errors"
)

// The drive mode helps us determine whether to read or write from
// the disk, but is actually unrelated to write protect mode!
const (
	// ReadMode is read mode for the drive.
	ReadMode = iota

	// WriteMode indicates that we are in write mode for the drive.
	WriteMode
)

// cyclesPerByte is the number of cycles whereby a specific byte may be read
// or written before drive spin would carry us to the next byte.
const cyclesPerByte uint64 = 32

// A Drive represents the state of a virtual Disk II drive.
type Drive struct {
	Phase        int
	Latch        uint8
	TrackPos     int
	SectorPos    int
	Data         *memory.Segment
	Image        *memory.Segment
	ImageType    int
	ImageName    string
	Stream       *os.File
	Mode         int
	WriteProtect bool
	Locked       bool

	// diskShifted is true if the disk has shifted after the last time data
	// was loaded into the latch.
	diskShifted bool

	// latchWasRead is true if the data in the latch has already been read.
	latchWasRead bool

	// motorOn is true if the motor is on. When the drive motor is on, the disk
	// contained in the drive will spin.
	motorOn bool
}

// NewDrive returns a new disk drive ready for DOS 3.3 images.
func NewDrive() *Drive {
	drive := new(Drive)

	drive.Mode = ReadMode
	drive.ImageType = a2enc.DOS33

	return drive
}

// StartMotor turns the drive motor on. In theory, this would cause the disk
// in the drive to spin.
func (d *Drive) StartMotor() {
	d.motorOn = true
}

// StopMotor turns off the drive motor, and theoretically stops the disk in
// the drive from spinning.
func (d *Drive) StopMotor() {
	d.motorOn = false
}

// MotorOn is true if the drive motor is on.
func (d *Drive) MotorOn() bool {
	return d.motorOn
}

// Position returns the segment position that the drive is currently at,
// based upon track and sector position.
func (d *Drive) Position() int {
	return ((d.TrackPos / 2) * a2enc.PhysTrackLen) + d.SectorPos
}

// Shift updates the sector position of the drive forward or backward by the
// given offset in bytes. Since tracks are circular and the disk is
// spinning, offsets that carry us beyond the bounds of the track instead
// bring us to the other end of the track.
func (d *Drive) Shift(offset int) {
	if d.Locked {
		return
	}

	d.SectorPos += offset

	// In practice, these for loops are mutually exclusive; only one of them
	// would ever be entered.
	for d.SectorPos >= a2enc.PhysTrackLen {
		d.SectorPos -= a2enc.PhysTrackLen
	}

	for d.SectorPos < 0 {
		d.SectorPos += a2enc.PhysTrackLen
	}

	d.diskShifted = true
}

// Step moves the track position forward or backward, depending on the
// sign of the offset. This simulates the stepper motor that moves the
// drive head further into the center of the disk platter (offset > 0)
// or further out (offset < 0).
func (d *Drive) Step(offset int) {
	d.TrackPos += offset

	switch {
	case d.TrackPos >= a2enc.MaxSteps:
		d.TrackPos = a2enc.MaxSteps - 1
	case d.TrackPos < 0:
		d.TrackPos = 0
	}
}

// SwitchPhase will figure out what phase we should be moving to based on a
// given address.
func (d *Drive) SwitchPhase(addr int) {
	phase := -1

	switch addr & 0xf {
	case 0x1:
		phase = 1
	case 0x3:
		phase = 2
	case 0x5:
		phase = 3
	case 0x7:
		phase = 4
	}

	d.PhaseTransition(phase)
}

// phaseTable is a really small set of state transitions we can make when
// accounting for the current phase (which is the row in the table) and the new
// phase (which is a column in that row). The 0 column and 0 row are not used
// since there is no 0 phase. We could refactor this a bit by removing those and
// subtracting 1 from the phase when mapping into this table.
var phaseTable = []int{
	0, 0, 0, 0, 0,
	0, 0, 1, 0, -1,
	0, -1, 0, 1, 0,
	0, 0, -1, 0, 1,
	0, 1, 0, -1, 0,
}

// PhaseTransition will make use of a given drive phase to step the drive
// forward or backward by some number of tracks. This is always going to be -1
// (step backward); 1 (step forward); or 0 (no change).
func (d *Drive) PhaseTransition(phase int) {
	if phase < 0 || phase > 4 {
		return
	}

	offset := phaseTable[(d.Phase*5)+phase]

	// Because of the above check, we can assert that the formula we use for
	// the phase transition ((curPhase * 5) + phase) will match something, so
	// we step immediately.
	d.Step(offset)

	// We also have to update our current phase
	d.Phase = phase
}

// ImageType returns the type of image that is suggested by the suffix
// of the given filename.
func ImageType(file string) (int, error) {
	lower := strings.ToLower(file)

	switch {
	case strings.HasSuffix(lower, ".do"), strings.HasSuffix(lower, ".dsk"):
		return a2enc.DOS33, nil
	case strings.HasSuffix(lower, ".nib"):
		return a2enc.Nibble, nil
	case strings.HasSuffix(lower, ".po"):
		return a2enc.ProDOS, nil
	}

	return -1, fmt.Errorf("unrecognized suffix for file %s", file)
}

// Load will read a file from the filesystem and set its contents as the
// image in the drive. It also decodes the contents according to the
// (detected) image type.
func (d *Drive) Load(r io.Reader, file string) error {
	var err error

	// See if we can figure out what type of image this is
	d.ImageType, err = ImageType(file)
	if err != nil {
		return errors.Wrapf(err, "failed to understand image type")
	}

	// Read the bytes from the file into a buffer
	bytes, err := io.ReadAll(r)
	if err != nil {
		return errors.Wrapf(err, "failed to read file %s", file)
	}

	// Copy directly into the image segment
	d.Image = memory.NewSegment(len(bytes))
	_, err = d.Image.CopySlice(0, []uint8(bytes))
	if err != nil {
		d.Image = nil
		return errors.Wrapf(err, "failed to copy bytes into image segment")
	}

	// Decode into the data segment
	d.Data, err = a2enc.Encode(d.ImageType, d.Image)
	if err != nil {
		d.Image = nil
		return errors.Wrapf(err, "failed to decode image")
	}

	// Reset the sector position, but leave track alone; the drive head
	// has not shifted since replacing the disk.
	d.SectorPos = 0

	d.ImageName = file

	return nil
}

// Write the contents of the drive's disk back to the filesystem
func (d *Drive) Save() error {
	// There's no file, so there's nothing to save.
	if d.ImageName == "" || d.Data == nil {
		return nil
	}

	logSegment, err := a2enc.Decode(d.ImageType, d.Data)
	if err != nil {
		return fmt.Errorf("could not decode image: %w", err)
	}

	return logSegment.WriteFile(d.ImageName)
}

// ReadLatch returns the byte that is currently loaded in the drive latch.
// If this data has not been read before, it is returned unmodified. If it has
// been read before, then it will be returned with the high bit set to zero.
func (d *Drive) ReadLatch() uint8 {
	if d.Data == nil {
		return 0xFF
	}

	if d.latchWasRead {
		return d.Latch & 0x7F
	}

	d.latchWasRead = true

	return d.Latch
}

// WriteLatch writes the byte in the drive's latch to the disk loaded in the
// drive. The drive must be in WriteMode for this operation to succeed, and
// the motor must be on. WriteLatch will not write any data if the latch byte
// does not have its high bit set to 1. The byte in the latch will be written
// to the current position of the drive head on the disk (with respect to
// track and sector).
func (d *Drive) WriteLatch() {
	if d.Data == nil {
		return
	}

	if d.Mode == WriteMode && d.MotorOn() && d.Latch&0x80 > 0 {
		d.Data.DirectSet(d.Position(), d.Latch)
	}
}

// LoadLatch reads the byte at the current drive position with respect to
// track and position. No change to the latch will occur if we cannot detect
// that the disk position has shifted since the last LoadLatch was called.
func (d *Drive) LoadLatch() {
	if d.Data == nil {
		return
	}

	if d.diskShifted {
		d.Latch = d.Data.DirectGet(d.Position())
		d.diskShifted = false
		d.latchWasRead = false
	}
}

// RandomByte returns a random byte as might be returned by the drive.
func (d *Drive) RandomByte() uint8 {
	return uint8(rand.IntN(0xFF))
}
