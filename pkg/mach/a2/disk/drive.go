package disk

import (
	"fmt"
	"io/ioutil"
	"log"
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
	Latch        data.Byte
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
func (d *Drive) Position() data.Int {
	if d.Data == nil {
		return 0
	}

	pos := data.Int(((d.TrackPos / 2) * sixtwo.PhysTrackLen) + d.SectorPos)

	return pos
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
func Phase(addr data.DByte) int {
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

// StepPhase will step the drive head forward or backward based upon the
// given address (from which we decipher the motor phase).
func (d *Drive) StepPhase(addr data.DByte) {
	newPhase := Phase(addr)
	curPhase := d.Phase
	offset := 0

	if newPhase < 1 || newPhase > 4 {
		return
	}

	switch {
	case newPhase == 1 && curPhase == 4:
		offset = 1
	case newPhase == 4 && curPhase == 1:
		offset = -1

	case newPhase > curPhase:
		offset = 1
	case newPhase < curPhase:
		offset = -1
	}

	d.Step(offset)
	d.Phase = newPhase

	log.Printf("step phase: new=%d, cur=%d, off=%d, phase=%d", newPhase, curPhase, offset, d.Phase)
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
func (d *Drive) Load(file string) error {
	var err error

	// See if we can figure out what type of image this is
	d.ImageType, err = ImageType(file)
	if err != nil {
		return errors.Wrapf(err, "failed to understand image type")
	}

	// Read the bytes from the file into a buffer
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.Wrapf(err, "failed to read file %s", file)
	}

	// Copy directly into the image segment
	d.Image = data.NewSegment(len(bytes))
	_, err = d.Image.CopySlice(0, data.ByteSlice(bytes))
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

func (d *Drive) Read() data.Byte {
	// Set the latch value to the byte at our current position, then
	// shift our position by one place
	d.Latch = d.Data.Get(d.Position())
	log.Printf("Reading position %x byte %x", d.Position(), d.Latch)

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
