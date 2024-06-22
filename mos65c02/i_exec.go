package mos65c02

import (
	"fmt"

	"github.com/pevans/erc/a2/a2sym"
	"github.com/pevans/erc/internal/metrics"
	"github.com/pevans/erc/statemap"
)

// Brk implements the BRK instruction, which is a hardware interrupt.
// This isn't something that normally happens in software, but you might
// see it in the system monitor (which was a debugger in Apple IIs).
func Brk(c *CPU) {
	metrics.Increment("instruction_brk", 1)

	// This pushes the the current PC register value in little-endian
	// order.
	c.PushStack(uint8(c.PC >> 8))
	c.PushStack(uint8(c.PC & 0xFF))

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
	c.PushStack(uint8(nextPos >> 8))
	c.PushStack(uint8(nextPos & 0xFF))

	// This is a bad hack to allow us only to look for builtin routines
	// when we're reading ROM. TODO: I should get rid of this hack.
	if c.State.Int(statemap.BankRead) == 1 {
		if routine := a2sym.Subroutine(int(c.EffAddr)); routine != "" {
			metrics.Increment(fmt.Sprintf("jsr_builtin_%s", routine), 1)
		}
	}

	c.PC = c.EffAddr
}

// Nop implements the NOP (no operation) instruction, which--well, it
// does nothing. On purpose.
func Nop(c *CPU) {
	// The only intentional NOP instruction is executed from opcode $EA.
	// Anything else may indicate that something is wrong -- like a bug
	// in the emulator that led us to execute data as if it were program
	// code.
	if c.Opcode != 0xEA {
		metrics.Increment("bad_opcodes", 1)
	}
}

// Np2 implements the NP2 instruction, which like NOP does nothing.
func Np2(c *CPU) {
	metrics.Increment("instruction_np2", 1)
}

// Np3 implements the NP3 instruction, which like NOP does nothing.
func Np3(c *CPU) {
	metrics.Increment("instruction_np3", 1)
}

// Rti implements the RTI (return from interrupt) instruction, which
// recovers the program state from a previous BRK operation.
func Rti(c *CPU) {
	c.P = c.PopStack()

	// Since we saved the bytes in BRK in order of msb, then lsb, we
	// need to pop them in the reverse order; lsb, then msb.
	lsb := uint16(c.PopStack())
	msb := uint16(c.PopStack())
	c.PC = (msb << 8) | lsb
}

// Rts implements the RTS (return from subroutine) instruction, which
// sets the program counter to the position saved from a previous JSR.
func Rts(c *CPU) {
	lsb := uint16(c.PopStack())
	msb := uint16(c.PopStack())

	c.PC = ((msb << 8) | lsb) + 1
}
