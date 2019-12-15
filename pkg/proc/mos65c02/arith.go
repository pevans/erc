package mos65c02

import "github.com/pevans/erc/pkg/data"

// Adc implements the ADC (add with carry) instruction. ADC is used to
// add integers to the accumulator; if the carry flag is set, then the
// result is further incremented by one.
func Adc(c *CPU) {
	// It's useful to make an accounting of how the result looks in a
	// 16-bit context (to determine if the C bit should be set)
	res16 := data.DByte(c.A)
	res16 += data.DByte(c.EffVal)
	res16 += data.DByte(c.P & CARRY)

	// But we mostly care about the 8-bit result, even if the unsigned
	// 8-bit value overflows
	res8 := data.Byte(res16)

	c.ApplyNZ(res8)
	c.ApplyStatus(res16 > 0xFF, CARRY)

	// That said, bear in mind that "overflow" in the MOS6502 means that
	// the value went from positive to negative, or from negative to
	// positive.
	c.ApplyStatus(
		(c.A^c.EffVal)&0x80 > 0 && (c.A^res8)&0x80 > 0,
		OVERFLOW,
	)

	c.A = res8
}

// Cmp implements the CMP (compare A) instruction, and compares with the A
// register. See the Compare method for more details.
func Cmp(c *CPU) {
	Compare(c, c.A)
}

// Cpx implements the CPX (compare X) instruction, and compares with the
// X register. See the Compare method for more details.
func Cpx(c *CPU) {
	Compare(c, c.X)
}

// Cpy implements the CPY (compare Y) instruction, and compares with the
// Y register. See the Compare method for more details.
func Cpy(c *CPU) {
	Compare(c, c.Y)
}

// Dec implements the DEC (decrement) instruction. DEC can decrement
// from the A register (if in the amAcc address mode), or can decrement
// from any address (depending on the other address modes used).
func Dec(c *CPU) {
	c.EffVal--
	c.ApplyNZ(c.EffVal)

	if c.AddrMode == amAcc {
		c.A = c.EffVal
		return
	}

	c.Set(c.EffAddr, c.EffVal)
}

// Dex implements the DEX (decrement X) instruction. DEX decrements only
// from the X register.
func Dex(c *CPU) {
	c.X--
	c.ApplyNZ(c.X)
}

// Dey implements the DEY (decrement Y) instruction. DEY decrements only
// from the Y register.
func Dey(c *CPU) {
	c.Y--
	c.ApplyNZ(c.Y)
}

// Inc implements the INC (increment) instruction. Like DEC, INC can
// increment from the A register or from any address in memory,
// depending on the addr mode.
func Inc(c *CPU) {
	c.EffVal++
	c.ApplyNZ(c.EffVal)

	if c.AddrMode == amAcc {
		c.A = c.EffVal
		return
	}

	c.WMem.Set(c.EffAddr, c.EffVal)
}

// Inx implements the INX (increment X) instruction. INX can only
// increment the X register.
func Inx(c *CPU) {
	c.X++
	c.ApplyNZ(c.X)
}

// Iny implements the INY (increment Y) instruction. INY can only
// increment the Y register.
func Iny(c *CPU) {
	c.Y++
	c.ApplyNZ(c.Y)
}

// Sbc implements the SBC (subtract with carry) instruction. SBC
// subtracts from the A register. If the carry flag is NOT set, then an
// additional one is subtracted from the result.
func Sbc(c *CPU) {
	res := int(c.A)
	res -= int(c.EffVal)

	if c.P&CARRY == 0 {
		res--
	}

	res8 := data.Byte(res)

	c.ApplyZ(res8)
	c.ApplyStatus(res < 0, NEGATIVE)
	c.ApplyStatus(res >= 0, CARRY)
	c.ApplyStatus(
		(c.A^c.EffVal)&0x80 > 0 && (c.A^res8)&0x80 > 0,
		OVERFLOW,
	)

	c.A = res8
}
