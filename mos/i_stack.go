package mos

// stackAddr returns the address indicated by the current position of the
// stack register.
func (c *CPU) stackAddr() uint16 {
	return 0x100 + uint16(c.S)
}

// PushStack adds the given byt to the stack, and decrements the stack
// counter. Note that in MOS 6502 chips, the stack counter begins life at 0xFF
// and we add to the stack from the _end_ of the stack page.
func (c *CPU) PushStack(byt uint8) {
	c.Set(c.stackAddr(), byt)
	c.S--
}

// PopStack increments the stack counter and returns byte at the current end
// of the stack.
func (c *CPU) PopStack() uint8 {
	c.S++
	return c.Get(c.stackAddr())
}
