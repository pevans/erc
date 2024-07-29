package a2

import (
	"fmt"
	"io"

	"github.com/pevans/erc/a2/a2state"
	"github.com/pkg/errors"
)

// Load will load a disk image from the given reader. The given filename is not
// strictly important -- we're not reading from a filesystem, we already have a
// reader -- but it will be used to help us determine what _kind_ of image it is
// (nibble, dos).
func (c *Computer) Load(r io.Reader, fileName string) error {
	if c.diskLog != nil {
		_ = c.diskLog.WriteToFile()
	}

	if err := c.Drive1.Load(r, fileName); err != nil {
		return errors.Wrapf(err, "could not read file: %s", fileName)
	}

	c.diskLog = nil

	if c.State.Bool(a2state.DebugImage) {
		c.diskLog = NewDiskLog(fileName)
		c.Drive1.Data.WriteFile(fmt.Sprintf("%v.physical", fileName))
	}

	return nil
}
