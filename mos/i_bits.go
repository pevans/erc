package mos

// saveResult makes a decision on the instruction level where to save
// the result of an operation. If we're in accumulator mode, then we
// save the result in the A register; if not, then we save it in memory
// at the effective address.
func (c *CPU) saveResult(res uint8) {
	if c.AddrMode == AmACC {
		c.A = res
	} else {
		c.Set(c.EffAddr, res)
	}
}

// And implements the AND instruction, which performs a bitwise-and on A
// and the effective value and saves the result there.
func And(c *CPU) {
	c.A &= c.EffVal
	c.ApplyNZ(c.A)
}

// Asl implements the ASL ("arithmetic" shift left) instruction, which
// shifts the EffVal left by one bit and saves the result. If the eighth
// bit is 1, then the carry flag is set. Don't ask why it's
// "arithmetic", because it doesn't make any sense to me.
func Asl(c *CPU) {
	res := int(c.EffVal) << 1

	c.ApplyNZ(uint8(res))
	c.ApplyStatus(res > 0xFF, CARRY)

	c.saveResult(uint8(res))
}

// Bit implements the BIT instruction, which performs several "bit"
// tests. It performs bitwise-and on A and EffVal; if that's zero, then
// the zero flag is set. (That result is not saved.) If the eighth bit is
// 1, then the negative flag is set. If the seventh bit is 1, then the
// overflow flag is set.
func Bit(c *CPU) {
	c.ApplyStatus((c.A&c.EffVal) == 0, ZERO)
	c.ApplyStatus((c.EffVal&0x80) > 0, NEGATIVE)
	c.ApplyStatus((c.EffVal&0x40) > 0, OVERFLOW)
}

// Bim implements the BIM instruction, which is like the BIT instruction
// but only executes the zero flag logic. (I don't really know what the
// "M" stands for.) (Update: This might abbreviate for (B)IT in
// (IM)mediate mode, as that is the only address mode where this
// instruction is configured. The Apple IIe technical reference doesn't
// distinguish the instructions; it merely notes that BIT operates
// differently, in the way that BIM does here, as opcode $89 (which is
// the code used here).
func Bim(c *CPU) {
	c.ApplyStatus((c.A&c.EffVal) == 0, ZERO)
}

// Eor implements the EOR instruction, which performs an exclusive-or
// operation between A and EffVal and saves the result in A.
func Eor(c *CPU) {
	c.A ^= c.EffVal
	c.ApplyNZ(c.A)
}

// Lsr implements the LSR ("logical" shift right) instruction, which
// shifts the EffVal right by one and saves the result. If the first bit
// is 1, then the carry flag is set. The negative flag is always UNSET,
// because the eighth bit is always zero as a result of this operation.
// Don't ask why this is "logical" and ASL is "arithmetic"; it doesn't
// make sense to me.
func Lsr(c *CPU) {
	res := c.EffVal >> 1

	// The result of a shift right is to _always_ unset the negative
	// flag.
	c.P &= ^NEGATIVE

	c.ApplyZ(res)
	c.ApplyStatus((c.EffVal&0x1) > 0, CARRY)

	c.saveResult(res)
}

// Ora implements the ORA instruction, which performs a bitwise-or
// operation on the A register and EffVal, and saves the result in A.
func Ora(c *CPU) {
	c.A |= c.EffVal
	c.ApplyNZ(c.A)
}

// Rol implements the ROL (rotate left) instruction, which rotates the
// EffVal by 1. The eighth bit is saved as the carry flag; if 1, then
// carry is on; if zero, it's off. The old carry flag value is saved in
// the first bit as a result of the operation.
//
// You can think of the rotate operation as a kind of nine-bit rotation,
// where the ninth bit is the carry flag.
func Rol(c *CPU) {
	res := int(c.EffVal) << 1

	if c.P&CARRY > 0 {
		res |= 0x1
	}

	c.ApplyStatus(res > 0xFF, CARRY)
	c.ApplyNZ(uint8(res))

	c.saveResult(uint8(res))
}

// Ror implements the ROR (rotate right) instruction, which rotates
// EffVal by 1. Again, like ROL, this is like a nine-bit rotation: the
// first bit is saved into the carry flag; the old carry flag value is
// saved as the eighth bit as a result of the operation.
func Ror(c *CPU) {
	res := c.EffVal >> 1

	if c.P&CARRY > 0 {
		res |= 0x80
	}

	c.ApplyStatus(c.EffVal&0x1 > 0, CARRY)
	c.ApplyNZ(res)

	c.saveResult(res)
}

// Trb implements the TRB instruction, which sets the zero flag if A &
// EffVal is not zero, and also saves result of (A exclusive-or 0xFF) &
// EffVal.
func Trb(c *CPU) {
	c.ApplyStatus(c.A&c.EffVal == 0, ZERO)
	c.Set(c.EffAddr, (c.A^0xff)&c.EffVal)
}

// Tsb implements the TSB instruction, which sets the zero flag if A &
// EffVal is not zero, and also saves the result of A | EffVal.
func Tsb(c *CPU) {
	c.ApplyStatus(c.A&c.EffVal == 0, ZERO)
	c.Set(c.EffAddr, c.A|c.EffVal)
}
