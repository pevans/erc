package a2

import "github.com/pevans/erc/pkg/obj"

// Boot steps through the boot procedures for the Apple II computer.
func (c *Computer) Boot() error {
	rom, err := obj.Slice(0, RomMemorySize)
	if err != nil {
		return err
	}

	return c.ROM.CopySlice(0, RomMemorySize, rom)
}
