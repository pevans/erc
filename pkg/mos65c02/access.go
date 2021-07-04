package mos65c02

import (
	"github.com/pevans/erc/pkg/data"
)

// Get will return the byte at a given address.
func (c *CPU) Get(addr data.DByte) data.Byte {
	return c.RMem.Get(addr.Int())
}

// Set will set the byte at a given address to the given value.
func (c *CPU) Set(addr data.DByte, val data.Byte) {
	c.WMem.Set(addr.Int(), val)
}

// Get16 returns a 16-bit value at a given address, which is read in
// little-endian order.
func (c *CPU) Get16(addr data.DByte) data.DByte {
	lsb := c.RMem.Get(addr.Int())
	msb := c.RMem.Get(addr.Int() + 1)

	return (data.DByte(msb) << 8) | data.DByte(lsb)
}

// Set16 sets the two bytes beginning at the given address to the given
// value. The bytes are set in little-endian order.
func (c *CPU) Set16(addr data.DByte, val data.DByte) {
	lsb := data.Byte(val & 0xFF)
	msb := data.Byte(val >> 8)

	c.WMem.Set(addr.Int(), lsb)
	c.WMem.Set(addr.Int()+1, msb)
}
