package mos65c02

// Get will return the byte at a given address.
func (c *CPU) Get(addr uint16) uint8 {
	return c.Memory.Get(int(addr))
}

// Set will set the byte at a given address to the given value.
func (c *CPU) Set(addr uint16, val uint8) {
	c.Memory.Set(int(addr), val)
}

// Get16 returns a 16-bit value at a given address, which is read in
// little-endian order.
func (c *CPU) Get16(addr uint16) uint16 {
	lsb := c.Memory.Get(int(addr))
	msb := c.Memory.Get(int(addr) + 1)

	return (uint16(msb) << 8) | uint16(lsb)
}

// Set16 sets the two bytes beginning at the given address to the given
// value. The bytes are set in little-endian order.
func (c *CPU) Set16(addr uint16, val uint16) {
	lsb := uint8(val & 0xFF)
	msb := uint8(val >> 8)

	c.Memory.Set(int(addr), lsb)
	c.Memory.Set(int(addr)+1, msb)
}
