package a2

import (
	"io"

	"github.com/pkg/errors"
)

// Load will load a disk image from the given reader. The given filename is not
// strictly important -- we're not reading from a filesystem, we already have a
// reader -- but it will be used to help us determine what _kind_ of image it is
// (nibble, dos).
func (c *Computer) Load(r io.Reader, fileName string) error {
	if err := c.Drive1.Load(r, fileName); err != nil {
		return errors.Wrapf(err, "could not read file: %s", fileName)
	}

	return nil
}
