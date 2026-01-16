package mos

// Get will return the byte at a given address.
func (c *CPU) Get(addr uint16) uint8 {
	return c.RMem.Get(int(addr))
}

// Set will set the byte at a given address to the given value.
func (c *CPU) Set(addr uint16, val uint8) {
	c.WMem.Set(int(addr), val)
}

// Get16 returns a 16-bit value at a given address, which is read in
// little-endian order.
func (c *CPU) Get16(addr uint16) uint16 {
	return c.RMem.Get16(int(addr))
}

// Set16 sets the two bytes beginning at the given address to the given value.
// The bytes are set in little-endian order.
func (c *CPU) Set16(addr uint16, val uint16) {
	c.WMem.Set16(int(addr), val)
}
