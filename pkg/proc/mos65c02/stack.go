package mos65c02

import "github.com/pevans/erc/pkg/mach"

// stackAddr returns the address indicated by the current position of
// the stack register.
func (c *CPU) stackAddr() mach.DByte {
	return 0x100 + mach.DByte(c.S)
}

// PushStack adds the given byt to the stack, and decrements the stack
// counter. Note that in MOS 6502 chips, the stack counter begins life
// at 0xFF and we add to the stack from the _end_ of the stack page.
func (c *CPU) PushStack(byt mach.Byte) {
	c.Set(c.stackAddr(), byt)
	c.S--
}

// PopStack increments the stack counter and returns byte at the current
// end of the stack.
func (c *CPU) PopStack() mach.Byte {
	c.S++
	return c.Get(c.stackAddr())
}
