package mos65c02

import "github.com/pevans/erc/pkg/mach"

// jumpIf will jump to EffAddr if the value in bits is greater than
// zero. The value of bits is typically computed with a bitwise
// operation, so we can say that if bits is non-zero then the operation
// was "successful".
func (c *CPU) jumpIf(bits mach.Byte) {
	if bits > 0 {
		c.PC = c.EffAddr
	} else {
		c.PC += 2
	}
}

// Bcc implements the BCC (branch on carry clear) instruction.
func Bcc(c *CPU) {
	c.jumpIf(^c.P & CARRY)
}

// Bcs implements the BCS (branch on carry set) instruction.
func Bcs(c *CPU) {
	c.jumpIf(c.P & CARRY)
}

// Beq implements the BEQ (branch on zero set, or branch on "equal to
// zero") instruction.
func Beq(c *CPU) {
	c.jumpIf(c.P & ZERO)
}

// Bmi implements the BMI (branch on negative set, or branch when
// "minus") instruction.
func Bmi(c *CPU) {
	c.jumpIf(c.P & NEGATIVE)
}

// Bne implements the BNE (branch on zero clear, or branch on "not
// equal to zero") instruction.
func Bne(c *CPU) {
	c.jumpIf(^c.P & ZERO)
}

// Bpl implements the BPL (branch on negative clear, or branch on
// "plus") instruction.
func Bpl(c *CPU) {
	c.jumpIf(^c.P & NEGATIVE)
}

// Bra implements the BRA (branch always) instruction.
func Bra(c *CPU) {
	c.PC = c.EffAddr
}

// Bvc implements the BVC (branch on overflow clear) instruction.
func Bvc(c *CPU) {
	c.jumpIf(^c.P & OVERFLOW)
}

// Bvs implements the BVS (branch on overflow set) instruction.
func Bvs(c *CPU) {
	c.jumpIf(c.P & OVERFLOW)
}
