package a2

import (
	"github.com/pkg/errors"
)

// Load a file as a disk into the main disk drive
func (c *Computer) Load(file string) error {
	if err := c.Drive1.Load(file); err != nil {
		return errors.Wrapf(err, "could not read file: %s", file)
	}

	return nil
}
