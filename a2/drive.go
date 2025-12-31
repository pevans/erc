package a2

import (
	"fmt"
	"io"
	"math/rand/v2"
	"strings"

	"github.com/pevans/erc/a2/a2enc"
	"github.com/pevans/erc/memory"

	"github.com/pkg/errors"
)

// The drive mode helps us determine whether to read or write from
// the disk, but is actually unrelated to write protect mode!
const (
	// readMode is read mode for the drive.
	readMode = iota

	// writeMode indicates that we are in write mode for the drive.
	writeMode
)

// A Drive represents the state of a virtual Disk II drive.
type Drive struct {
	// phase is the stepper motor phase that the drive head is currently at.
	// We track the phase to determine when to adjust our track position.
	phase int

	// latch is the byte that we last read from the disk, or is the byte that
	// we may write to the disk. Anything coming out of the disk, or going
	// into it, must go through the latch (which makes it a bit like an
	// airlock for a spaceship).
	latch uint8

	// trackPos is the current track that the drive head is stationed at. Any
	// sectors we read or write will be found in that track.
	trackPos int

	// sectorPos is the position of the drive head within a given track.
	// Whenever you read a byte or write a byte, you are doing so at the
	// drive's sectorPos.
	sectorPos int

	// data is the physically encoded form of the bytes that we read from a
	// disk image.
	data *memory.Segment

	// image is the memory segment containing the bytes of the image file
	// loaded in the drive. These bytes may be the logical form of the data,
	// or they may be the physical form if the image was a nibble file.
	image *memory.Segment

	// imageType is the type of the image file loaded in the drive (DOS33,
	// ProDOS).
	imageType int

	// imageName is the name of the image file loaded in the drive.
	imageName string

	// mode is the read/write mode of the drive. A drive can either be in read
	// mode or in write mode; never both together, and never neither mode.
	mode int

	// writeProtect is true when the disk in the drive is considered
	// write-protected. A write-protected disk may not be written to by the
	// drive.
	writeProtect bool

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

	drive.SetReadMode()
	drive.imageType = a2enc.DOS33

	return drive
}

// ImageName returns the name of the image file loaded in the drive
func (d *Drive) ImageName() string {
	return d.imageName
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

// SetWriteProtect will change the writeProtect status of a drive to the given
// status.
func (d *Drive) SetWriteProtect(status bool) {
	d.writeProtect = status
}

// ToggleWriteProtect flips the status of write protection for the disk in a
// drive. If it was true, it becomes false, and vice-versa.
func (d *Drive) ToggleWriteProtect() {
	d.writeProtect = !d.writeProtect
}

// WriteProtected returns true if the disk in the drive is write-protected
// (can't be written to).
func (d *Drive) WriteProtected() bool {
	return d.writeProtect
}

// ReadMode returns true if the drive is in read mode (is able to read data
// from the disk)
func (d *Drive) ReadMode() bool {
	return d.mode == readMode
}

// SetReadMode sets the drive to read mode.
func (d *Drive) SetReadMode() {
	d.mode = readMode
}

// WriteMode returns true if the drive is in write mode (is able to write data
// to the disk). This does not take write protection into account; it's
// possible for a drive to be in write mode but still be unable to write to a
// disk that is write-protected.
func (d *Drive) WriteMode() bool {
	return d.mode == writeMode
}

// SetWriteMode sets the drive to write mode.
func (d *Drive) SetWriteMode() {
	d.mode = writeMode
}

// Position returns the segment position that the drive is currently at,
// based upon track and sector position.
func (d *Drive) Position() int {
	return ((d.trackPos / 2) * a2enc.PhysTrackLen) + d.sectorPos
}

// Shift updates the sector position of the drive forward or backward by the
// given offset in bytes. Since tracks are circular and the disk is
// spinning, offsets that carry us beyond the bounds of the track instead
// bring us to the other end of the track.
func (d *Drive) Shift(offset int) {
	d.sectorPos += offset

	// In practice, these for loops are mutually exclusive; only one of them
	// would ever be entered.
	for d.sectorPos >= a2enc.PhysTrackLen {
		d.sectorPos -= a2enc.PhysTrackLen
	}

	for d.sectorPos < 0 {
		d.sectorPos += a2enc.PhysTrackLen
	}

	d.diskShifted = true
}

// Step moves the track position forward or backward, depending on the
// sign of the offset. This simulates the stepper motor that moves the
// drive head further into the center of the disk platter (offset > 0)
// or further out (offset < 0).
func (d *Drive) Step(offset int) {
	d.trackPos += offset

	switch {
	case d.trackPos >= a2enc.MaxSteps:
		d.trackPos = a2enc.MaxSteps - 1
	case d.trackPos < 0:
		d.trackPos = 0
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

	offset := phaseTable[(d.phase*5)+phase]

	// Because of the above check, we can assert that the formula we use for
	// the phase transition ((curPhase * 5) + phase) will match something, so
	// we step immediately.
	d.Step(offset)

	// We also have to update our current phase
	d.phase = phase
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
	d.imageType, err = ImageType(file)
	if err != nil {
		return errors.Wrapf(err, "failed to understand image type")
	}

	// Read the bytes from the file into a buffer
	bytes, err := io.ReadAll(r)
	if err != nil {
		return errors.Wrapf(err, "failed to read file %s", file)
	}

	// Copy directly into the image segment
	d.image = memory.NewSegment(len(bytes))
	_, err = d.image.CopySlice(0, []uint8(bytes))
	if err != nil {
		d.image = nil
		return errors.Wrapf(err, "failed to copy bytes into image segment")
	}

	// Decode into the data segment
	d.data, err = a2enc.Encode(d.imageType, d.image)
	if err != nil {
		d.image = nil
		return errors.Wrapf(err, "failed to decode image")
	}

	// Reset the sector position, but leave track alone; the drive head
	// has not shifted since replacing the disk.
	d.sectorPos = 0

	// If the disk had write-protected status, we should assume the next disk
	// loaded does not have it
	d.writeProtect = false

	d.imageName = file

	return nil
}

// RemoveDisk will essentially treat the drive as empty. This method DOES NOT
// SAVE ANY DATA -- please call the Save method to do that. Additionally, this
// method is not strictly necessary if you are swapping one disk for another.
// Instead, you can simply call Load to do that. RemoveDisk is only useful if
// you have a use-case to treat the drive as functionally empty.
func (d *Drive) RemoveDisk() {
	d.imageName = ""
	d.image = nil
	d.data = nil
}

// Write the contents of the drive's disk back to the filesystem
func (d *Drive) Save() error {
	// There's no file, so there's nothing to save.
	if d.imageName == "" || d.data == nil {
		return nil
	}

	logSegment, err := a2enc.Decode(d.imageType, d.data)
	if err != nil {
		return fmt.Errorf("could not decode image: %w", err)
	}

	return logSegment.WriteFile(d.imageName)
}

// ReadLatch returns the byte that is currently loaded in the drive latch.
// If this data has not been read before, it is returned unmodified. If it has
// been read before, then it will be returned with the high bit set to zero.
func (d *Drive) ReadLatch() uint8 {
	if d.data == nil {
		return 0xFF
	}

	if d.latchWasRead {
		return d.latch & 0x7F
	}

	d.latchWasRead = true

	return d.latch
}

// WriteLatch writes the byte in the drive's latch to the disk loaded in the
// drive. The drive must be in WriteMode for this operation to succeed, and
// the motor must be on. WriteLatch will not write any data if the latch byte
// does not have its high bit set to 1. The byte in the latch will be written
// to the current position of the drive head on the disk (with respect to
// track and sector).
func (d *Drive) WriteLatch() {
	if d.data == nil {
		return
	}

	if d.WriteMode() && d.MotorOn() && d.latch&0x80 > 0 {
		d.data.DirectSet(d.Position(), d.latch)
	}
}

// LoadLatch reads the byte at the current drive position with respect to
// track and position. No change to the latch will occur if we cannot detect
// that the disk position has shifted since the last LoadLatch was called.
func (d *Drive) LoadLatch() {
	if d.data == nil {
		return
	}

	if d.diskShifted {
		d.latch = d.data.DirectGet(d.Position())
		d.diskShifted = false
		d.latchWasRead = false
	}
}

// RandomByte returns a random byte as might be returned by the drive.
func (d *Drive) RandomByte() uint8 {
	return uint8(rand.IntN(0xFF))
}
