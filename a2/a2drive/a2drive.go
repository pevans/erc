package a2drive

import (
	"math/rand/v2"

	"github.com/pevans/erc/a2/a2enc"
	"github.com/pevans/erc/memory"
)

// The drive mode helps us determine whether to read or write from the disk,
// but is actually unrelated to write protect mode!
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
	// sectors we read or write will be found in that track. Note that this
	// number is technically stored as half tracks.
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

	// motorOn is true if the motor is on. When the drive motor is on, the
	// disk contained in the drive will spin.
	motorOn bool
}

// phaseTable is a really small set of state transitions we can make when
// accounting for the current phase (which is the row in the table) and the
// new phase (which is a column in that row). The 0 column and 0 row are not
// used since there is no 0 phase. We could refactor this a bit by removing
// those and subtracting 1 from the phase when mapping into this table.
var phaseTable = []int{
	0, 0, 0, 0, 0,
	0, 0, 1, 0, -1,
	0, -1, 0, 1, 0,
	0, 0, -1, 0, 1,
	0, 1, 0, -1, 0,
}

// NewDrive returns a new disk drive ready for DOS 3.3 images.
func NewDrive() *Drive {
	drive := new(Drive)

	drive.SetReadMode()
	drive.imageType = a2enc.DOS33

	return drive
}

// MotorOn is true if the drive motor is on.
func (d *Drive) MotorOn() bool {
	return d.motorOn
}

// RandomByte returns a random byte as might be returned by the drive.
func (d *Drive) RandomByte() uint8 {
	return uint8(rand.IntN(0xFF))
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

// WriteDataToFile writes the data segment's data to the provided filename. If
// that operation is not successful, a non-nil error is returned.
func (d *Drive) WriteDataToFile(filename string) error {
	if d.data == nil {
		return nil
	}

	return d.data.WriteFile(filename)
}
