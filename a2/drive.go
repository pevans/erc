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

	newLatchData bool

	// motorOn is true if the motor is on. When the drive motor is on, the disk
	// contained in the drive will spin.
	motorOn bool

	// cyclesSinceLastSpin is number of cycles the CPU has executed at the time
	// we last needed to spin the disk platter (the wafer contained within the
	// plastic shell of a floppy disk).
	cyclesSinceLastSpin uint64
}

// NewDrive returns a new disk drive ready for DOS 3.3 images.
func NewDrive() *Drive {
	drive := new(Drive)

	drive.Mode = ReadMode
	drive.ImageType = a2enc.DOS33

	return drive
}

// StartMotor turns the drive motor on and starts spinning the disk platter.
func (d *Drive) StartMotor(cycles uint64) {
	d.motorOn = true
	d.cyclesSinceLastSpin = cycles
}

// StopMotor turns off the drive motor, and ceases spinning the disk platter.
func (d *Drive) StopMotor() {
	d.motorOn = false
}

// MotorOn is true if the drive motor is on.
func (d *Drive) MotorOn() bool {
	return d.motorOn
}

// SpinDisk will, if a drive motor is on, spin the disk -- adjusting the
// sector position based on the number of cycles executed.
func (d *Drive) SpinDisk(cycles uint64) {
	// If the drive is not on, then we don't want to adjust our position
	if !d.MotorOn() {
		return
	}

	// We can assume that cycles is an essentially monotonic number that can
	// only go up, and thus will always be equal to or greater than the cycles
	// since last spin
	diff := cycles - d.cyclesSinceLastSpin

	// Since we don't know how many cycles it's been since we last shifted our
	// position, we may need to shift by many positions. Note that the final
	// cyclesPerLastSpin value may _not_ be equal to the given cycles.
	bytes := diff / cyclesPerByte
	d.Shift(int(bytes))
	d.cyclesSinceLastSpin += (bytes * cyclesPerByte)

	// Set the latch and let everyone know that there's new data
	if bytes > 0 && d.Mode == ReadMode {
		d.Latch = d.Data.DirectGet(d.Position())
		d.newLatchData = true
	}
}

// Position returns the segment position that the drive is currently at,
// based upon track and sector position.
func (d *Drive) Position() int {
	return ((d.TrackPos / 2) * a2enc.PhysTrackLen) + d.SectorPos
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

	if d.SectorPos >= a2enc.PhysTrackLen || d.SectorPos < 0 {
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

func (d *Drive) Read() uint8 {
	if d.Data == nil {
		return 0xFF
	}

	if !d.newLatchData {
		return d.Latch & 0x7F
	}

	d.newLatchData = false
	return d.Latch
}

func (d *Drive) Write() {
	if d.Data == nil {
		return
	}

	// We can only write our latch value if the high-bit is set
	if d.Mode == WriteMode && d.MotorOn() && d.Latch&0x80 > 0 {
		d.Data.DirectSet(d.Position(), d.Latch)
	}
}

func (d *Drive) RandomByte() uint8 {
	return uint8(rand.IntN(0xFF))
}
