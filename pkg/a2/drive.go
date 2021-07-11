package a2

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pevans/erc/pkg/data"
	"github.com/pevans/erc/pkg/sixtwo"

	"github.com/pkg/errors"
)

const (
	// Nibble is the disk type for nibble (*.NIB) disk images.
	// Although this is an image type, it's not something that actual
	// disks would have been formatted in during the Apple II era.
	Nibble = iota
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
	Latch        uint8
	TrackPos     int
	SectorPos    int
	Data         *data.Segment
	Image        *data.Segment
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
	drive.ImageType = sixtwo.DOS33

	return drive
}

// Position returns the segment position that the drive is currently at,
// based upon track and sector position.
func (d *Drive) Position() int {
	return ((d.TrackPos / 2) * sixtwo.PhysTrackLen) + d.SectorPos
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

	if d.SectorPos >= sixtwo.PhysTrackLen || d.SectorPos < 0 {
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
	case d.TrackPos > sixtwo.MaxSteps:
		d.TrackPos = sixtwo.MaxSteps
	case d.TrackPos < 0:
		d.TrackPos = 0
	}

	// The sector position also resets when the drive motor steps
	d.SectorPos = 0
}

// Phase returns the motor phase based upon the given address.
func Phase(addr uint16) int {
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

// SwitchPhase will figure out what phase we should be moving to based on a
// given address.
func (d *Drive) SwitchPhase(addr uint16) {
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

	// Because of the above check, we can assert that the formula we use for the
	// phase transition ((curPhase * 5) + phase) will match something, so we
	// step immediately.
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
		return sixtwo.DOS33, nil
	case strings.HasSuffix(lower, ".nib"):
		return Nibble, nil
	case strings.HasSuffix(lower, ".po"):
		return sixtwo.ProDOS, nil
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
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return errors.Wrapf(err, "failed to read file %s", file)
	}

	// Copy directly into the image segment
	d.Image = data.NewSegment(len(bytes))
	_, err = d.Image.CopySlice(0, []uint8(bytes))
	if err != nil {
		d.Image = nil
		return errors.Wrapf(err, "failed to copy bytes into image segment")
	}

	// Decode into the data segment
	d.Data, err = sixtwo.Encode(d.ImageType, d.Image)
	if err != nil {
		d.Image = nil
		return errors.Wrapf(err, "failed to decode image")
	}

	return nil
}

func (d *Drive) Read() uint8 {
	// Set the latch value to the byte at our current position, then
	// shift our position by one place
	d.Latch = d.Data.Get(d.Position())

	d.Shift(1)

	return d.Latch
}

func (d *Drive) Write() {
	// We can only write our latch value if the high-bit is set
	if d.Latch&0x80 > 0 {
		d.Data.Set(d.Position(), d.Latch)
		d.Shift(1)
	}
}
