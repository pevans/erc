package mos65c02

import (
	"fmt"
)

func (s *mosSuite) TestAdc() {
	var (
		d10    uint8 = 10
		d250   uint8 = 250
		d127   uint8 = 127
		offset uint8 = 10
		one    uint8 = 1
	)

	s.Run("result is A + EffVal when carry isn't set", func() {
		s.cpu.A = d10
		s.cpu.EffVal = offset
		s.cpu.P = 0
		Adc(s.cpu)

		s.Equal(d10+offset, s.cpu.A)
	})

	s.Run("result is A + EffVal + 1 when carry is set", func() {
		s.cpu.A = d10
		s.cpu.EffVal = offset
		s.cpu.P = CARRY
		Adc(s.cpu)

		s.Equal(d10+offset+one, s.cpu.A)
	})

	s.Run("carry is set when result is larger than 255", func() {
		s.cpu.A = d250
		s.cpu.EffVal = offset
		s.cpu.P = 0

		Adc(s.cpu)
		s.Equal(d250+offset, s.cpu.A)
		s.Equal(CARRY, s.cpu.P&CARRY)
	})

	s.Run("overflow is set when going from positive to negative", func() {
		// Test going from positive to negative sets OVERFLOW
		s.cpu.A = d127
		s.cpu.EffVal = offset
		s.cpu.P = 0
		Adc(s.cpu)

		s.Equal(d127+offset, s.cpu.A)
		s.Equal(OVERFLOW, s.cpu.P&OVERFLOW)
	})

	s.Run("overflow is set when going from negative to positive", func() {
		s.cpu.A = d250
		s.cpu.EffVal = offset
		s.cpu.P = 0
		Adc(s.cpu)

		s.Equal(d250+offset, s.cpu.A)
		s.Equal(OVERFLOW, s.cpu.P&OVERFLOW)
	})
}

func (s *mosSuite) testCompare(val *uint8, fn func(*CPU)) {
	var (
		d10 uint8 = 10
		d20 uint8 = 20
	)

	s.Run("zero is set when given equal values", func() {
		*val = d10
		s.cpu.EffVal = d10
		fn(s.cpu)

		s.Equal(ZERO, s.cpu.P&ZERO)
	})

	s.Run("negative is set when effval > accum", func() {
		*val = d10
		s.cpu.EffVal = d20
		fn(s.cpu)

		s.Equal(NEGATIVE, s.cpu.P&NEGATIVE)
	})

	s.Run("carry is set when accum >= effval", func() {
		*val = d20
		s.cpu.EffVal = d10
		fn(s.cpu)

		s.Equal(CARRY, s.cpu.P&CARRY)

		*val = d10
		s.cpu.P = 0
		fn(s.cpu)
		s.Equal(CARRY, s.cpu.P&CARRY)
	})
}

func (s *mosSuite) TestCmp() {
	s.testCompare(&s.cpu.A, Cmp)
}

func (s *mosSuite) TestCpx() {
	s.testCompare(&s.cpu.X, Cpx)
}

func (s *mosSuite) TestCpy() {
	s.testCompare(&s.cpu.Y, Cpy)
}

func (s *mosSuite) testDecrement(
	funcName string,
	setVal func(*CPU, uint8),
	getVal func(*CPU) uint8,
	fn func(*CPU),
) {
	var (
		d1 uint8 = 1
		d0 uint8 = 0
	)

	runVal := fmt.Sprintf("%s decrements by one", funcName)
	runNegative := fmt.Sprintf("%s sets negative flag when flipping sign", funcName)
	runZero := fmt.Sprintf("%s sets zero flag when result is zero", funcName)

	s.Run(runVal, func() {
		setVal(s.cpu, d1)
		fn(s.cpu)

		s.Equal(d0, getVal(s.cpu))
	})

	s.Run(runNegative, func() {
		setVal(s.cpu, d0)
		fn(s.cpu)

		s.Equal(NEGATIVE, s.cpu.P&NEGATIVE)
	})

	s.Run(runZero, func() {
		setVal(s.cpu, d1)
		fn(s.cpu)

		s.Equal(ZERO, s.cpu.P&ZERO)
	})
}

func (s *mosSuite) TestDex() {
	s.testDecrement("DEX",
		func(c *CPU, b uint8) { c.X = b },
		func(c *CPU) uint8 { return c.X },
		Dex,
	)
}

func (s *mosSuite) TestDey() {
	s.testDecrement("DEY",
		func(c *CPU, b uint8) { c.Y = b },
		func(c *CPU) uint8 { return c.Y },
		Dey,
	)
}

// TestDec is a bit tricky, because DEC does two very different things based on
// its address mode.
func (s *mosSuite) TestDec() {
	s.cpu.AddrMode = amAcc
	s.testDecrement("DEC (accumulator)",
		func(c *CPU, b uint8) { c.A = b; c.EffVal = b },
		func(c *CPU) uint8 { return c.A },
		Dec,
	)

	var (
		addr uint16 = 1
	)

	s.cpu.AddrMode = amAbs
	s.testDecrement("DEC (memory)",
		func(c *CPU, b uint8) { c.Set(addr, b); c.EffAddr = addr; c.EffVal = b },
		func(c *CPU) uint8 { return c.Get(addr) },
		Dec,
	)
}

func (s *mosSuite) testIncrement(
	funcName string,
	setVal func(*CPU, uint8),
	getVal func(*CPU) uint8,
	fn func(*CPU),
) {
	var (
		d1   uint8 = 1
		d0   uint8 = 0
		d127 uint8 = 127
		d255 uint8 = 255
	)

	runVal := fmt.Sprintf("%s increments by one", funcName)
	runNegative := fmt.Sprintf("%s sets negative flag when flipping sign", funcName)
	runZero := fmt.Sprintf("%s sets zero flag when result is zero", funcName)

	s.Run(runVal, func() {
		setVal(s.cpu, d0)
		fn(s.cpu)

		s.Equal(d1, getVal(s.cpu))
	})

	s.Run(runNegative, func() {
		setVal(s.cpu, d127)
		fn(s.cpu)

		s.Equal(NEGATIVE, s.cpu.P&NEGATIVE)
	})

	s.Run(runZero, func() {
		setVal(s.cpu, d255)
		fn(s.cpu)

		s.Equal(ZERO, s.cpu.P&ZERO)
	})
}

func (s *mosSuite) TestInx() {
	s.testIncrement("INX",
		func(c *CPU, b uint8) { c.X = b },
		func(c *CPU) uint8 { return c.X },
		Inx,
	)
}

func (s *mosSuite) TestIny() {
	s.testIncrement("INY",
		func(c *CPU, b uint8) { c.Y = b },
		func(c *CPU) uint8 { return c.Y },
		Iny,
	)
}

// Same deal as TestDec above -- the INC instruction does some very different
// things based on address mode.
func (s *mosSuite) TestInc() {
	s.cpu.AddrMode = amAcc
	s.testIncrement("INC",
		func(c *CPU, b uint8) { c.A = b; c.EffVal = b },
		func(c *CPU) uint8 { return c.A },
		Inc,
	)

	var (
		addr uint16 = 1
	)

	s.cpu.AddrMode = amAbs
	s.testIncrement("INC",
		func(c *CPU, b uint8) { c.Set(addr, b); c.EffAddr = addr; c.EffVal = b },
		func(c *CPU) uint8 { return c.Get(addr) },
		Inc,
	)
}

func (s *mosSuite) TestSbc() {
	var (
		d10    uint8 = 10
		d127   uint8 = 127
		offset uint8 = 20
		one    uint8 = 1
		d0     uint8 = 0
	)

	s.Run("subtracting with carry set results in A = A - EV", func() {
		s.cpu.A = d127
		s.cpu.EffVal = offset
		s.cpu.P = CARRY
		Sbc(s.cpu)

		s.Equal(d127-offset, s.cpu.A)
	})

	s.Run("subtraction with nonzero result sets carry", func() {
		s.Equal(CARRY, s.cpu.P&CARRY)
	})

	s.Run("subtracting without carry sets results in A = A - EV - 1", func() {
		s.cpu.A = d127
		s.cpu.EffVal = offset
		s.cpu.P = 0
		Sbc(s.cpu)

		s.Equal(d127-offset-one, s.cpu.A)
	})

	s.Run("subtracting a larger from a smaller number sets negative", func() {
		s.cpu.A = d10
		s.cpu.EffVal = offset
		s.cpu.P = CARRY
		Sbc(s.cpu)

		s.Equal(d10-offset, s.cpu.A)
		s.Equal(NEGATIVE, s.cpu.P&NEGATIVE)
	})

	s.Run("subtracting that flips the sign positive to negative sets overflow", func() {
		s.Equal(OVERFLOW, s.cpu.P&OVERFLOW)
	})

	s.Run("subtracting that flips the sign negative to positve sets overflow", func() {
		s.cpu.A = d127 + offset
		s.cpu.EffVal = offset
		s.cpu.P = CARRY
		Sbc(s.cpu)

		s.Equal(d127, s.cpu.A)
		s.Equal(OVERFLOW, s.cpu.P&OVERFLOW)
	})

	s.Run("subtracting that results in zero sets zero", func() {
		s.cpu.A = d127
		s.cpu.EffVal = d127
		s.cpu.P = CARRY
		Sbc(s.cpu)

		s.Equal(d0, s.cpu.A)
		s.Equal(ZERO, s.cpu.P&ZERO)
	})
}
