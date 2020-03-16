package a2

import (
	"io"

	"github.com/pkg/errors"
)

func (c *Computer) Load(r io.Reader, fileName string) error {
	if err := c.Drive1.Load(r, fileName); err != nil {
		return errors.Wrapf(err, "could not read file: %s", fileName)
	}

	return nil
}
