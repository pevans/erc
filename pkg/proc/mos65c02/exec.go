package mos65c02

import "github.com/pevans/erc/pkg/mach"

// Brk implements the BRK instruction, which is a hardware interrupt.
// This isn't something that normally happens in software, but you might
// see it in the system monitor (which was a debugger in Apple IIs).
func Brk(c *CPU) {
	// This pushes the the current PC register value in little-endian
	// order.
	c.PushStack(mach.Byte(c.PC >> 8))
	c.PushStack(mach.Byte(c.PC & 0xFF))

	// Also hang onto the status
	c.PushStack(c.P)

	// Always set INTERRUPT, always remove DECIMAL
	c.P |= INTERRUPT
	c.P &^= DECIMAL

	c.PC += 2
}

// Jmp implements the JMP instruction, which sets the program counter to
// the effective address.
func Jmp(c *CPU) {
	c.PC = c.EffAddr
}

// Jsr implements the JSR (jump to subroutine) instruction, which saves
// the current program counter to the stack and then sets PC to EffAddr.
func Jsr(c *CPU) {
	nextPos := c.PC + 2

	// We have to save the position that we should jump back to after we
	// return from subroutine (RTS) in the stack.
	c.PushStack(mach.Byte(nextPos >> 8))
	c.PushStack(mach.Byte(nextPos & 0xFF))

	c.PC = c.EffAddr
}

// Nop implements the NOP (no operation) instruction, which--well, it
// does nothing. On purpose.
func Nop(c *CPU) {
}

// Np2 implements the NP2 instruction, which like NOP does nothing.
func Np2(c *CPU) {
}

// Np3 implements the NP3 instruction, which like NOP does nothing.
func Np3(c *CPU) {
}

// Rti implements the RTI (return from interrupt) instruction, which
// recovers the program state from a previous BRK operation.
func Rti(c *CPU) {
	c.P = c.PopStack()

	// Since we saved the bytes in BRK in order of msb, then lsb, we
	// need to pop them in the reverse order; lsb, then msb.
	lsb := mach.DByte(c.PopStack())
	msb := mach.DByte(c.PopStack())
	c.PC = (msb << 8) | lsb
}

// Rts implements the RTS (return from subroutine) instruction, which
// sets the program counter to the position saved from a previous JSR.
func Rts(c *CPU) {
	lsb := mach.DByte(c.PopStack())
	msb := mach.DByte(c.PopStack())

	c.PC = ((msb << 8) | lsb) + 1
}
