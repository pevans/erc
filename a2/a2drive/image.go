package a2drive

import (
	"fmt"
	"io"
	"strings"

	"github.com/pevans/erc/a2/a2enc"
	"github.com/pevans/erc/memory"
)

// ImageType returns the type of image that is suggested by the suffix of the
// given filename.
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

// Load will read a file from the filesystem and set its contents as the image
// in the drive. It also decodes the contents according to the (detected)
// image type.
func (d *Drive) Load(r io.Reader, file string) error {
	var err error

	// See if we can figure out what type of image this is
	d.imageType, err = ImageType(file)
	if err != nil {
		return fmt.Errorf("failed to understand image type: %w", err)
	}

	// Read the bytes from the file into a buffer
	bytes, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", file, err)
	}

	// Validate the file size based on image type
	expectedSize, err := a2enc.Size(d.imageType)
	if err != nil {
		return fmt.Errorf("failed to determine expected size for image type: %w", err)
	}

	if len(bytes) != expectedSize {
		return fmt.Errorf(
			"invalid %s file size: got %d bytes, expected %d bytes",
			file, len(bytes), expectedSize,
		)
	}

	// Copy directly into the image segment
	d.image = memory.NewSegment(len(bytes))
	_, err = d.image.CopySlice(0, []uint8(bytes))
	if err != nil {
		d.image = nil
		return fmt.Errorf("failed to copy bytes into image segment: %w", err)
	}

	// Decode into the data segment
	d.data, err = a2enc.Encode(d.imageType, d.image)
	if err != nil {
		d.image = nil
		return fmt.Errorf("failed to decode image: %w", err)
	}

	// Reset the sector position, but leave track alone; the drive head has
	// not shifted since replacing the disk.
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
