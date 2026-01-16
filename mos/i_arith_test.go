package mos_test

import (
	"fmt"
	"testing"

	"github.com/pevans/erc/mos"
	"github.com/stretchr/testify/assert"
)

func TestNewDecimal(t *testing.T) {
	assert.Equal(t, mos.Decimal{Result: 3}, mos.NewDecimal(0b00000011))
	assert.Equal(t, mos.Decimal{Result: 13}, mos.NewDecimal(0b00010011))
	assert.Equal(t, mos.Decimal{Result: 99}, mos.NewDecimal(0b10011001))

	d := mos.NewDecimal(0b00001010)
	assert.Error(t, d.Error)

	d = mos.NewDecimal(0b10100000)
	assert.Error(t, d.Error)
}

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
		mos.Adc(s.cpu)

		s.Equal(d10+offset, s.cpu.A)
	})

	s.Run("result is A + EffVal + 1 when carry is set", func() {
		s.cpu.A = d10
		s.cpu.EffVal = offset
		s.cpu.P = mos.CARRY
		mos.Adc(s.cpu)

		s.Equal(d10+offset+one, s.cpu.A)
	})

	s.Run("carry is set when result is larger than 255", func() {
		s.cpu.A = d250
		s.cpu.EffVal = offset
		s.cpu.P = 0

		mos.Adc(s.cpu)
		s.Equal(d250+offset, s.cpu.A)
		s.Equal(mos.CARRY, s.cpu.P&mos.CARRY)
	})

	s.Run("overflow is set when going from positive to negative", func() {
		// Test going from positive to negative sets OVERFLOW
		s.cpu.A = d127
		s.cpu.EffVal = offset
		s.cpu.P = 0
		mos.Adc(s.cpu)

		s.Equal(d127+offset, s.cpu.A)
		s.Equal(mos.OVERFLOW, s.cpu.P&mos.OVERFLOW)
	})

	s.Run("overflow is set when going from negative to positive", func() {
		s.cpu.A = d250
		s.cpu.EffVal = offset
		s.cpu.P = 0
		mos.Adc(s.cpu)

		s.Equal(d250+offset, s.cpu.A)
		s.Equal(mos.OVERFLOW, s.cpu.P&mos.OVERFLOW)
	})

	s.Run("decimal works", func() {
		s.cpu.A = 0x6
		s.cpu.EffVal = 0x6
		s.cpu.P = mos.DECIMAL

		mos.Adc(s.cpu)
		s.Equal(uint8(0x12), s.cpu.A)
	})
}

func (s *mosSuite) testCompare(val *uint8, fn func(*mos.CPU)) {
	var (
		d10 uint8 = 10
		d20 uint8 = 20
	)

	s.Run("zero is set when given equal values", func() {
		*val = d10
		s.cpu.EffVal = d10
		fn(s.cpu)

		s.Equal(mos.ZERO, s.cpu.P&mos.ZERO)
	})

	s.Run("negative is set when effval > accum", func() {
		*val = d10
		s.cpu.EffVal = d20
		fn(s.cpu)

		s.Equal(mos.NEGATIVE, s.cpu.P&mos.NEGATIVE)
	})

	s.Run("carry is set when accum >= effval", func() {
		*val = d20
		s.cpu.EffVal = d10
		fn(s.cpu)

		s.Equal(mos.CARRY, s.cpu.P&mos.CARRY)

		*val = d10
		s.cpu.P = 0
		fn(s.cpu)
		s.Equal(mos.CARRY, s.cpu.P&mos.CARRY)
	})
}

func (s *mosSuite) TestCmp() {
	s.testCompare(&s.cpu.A, mos.Cmp)
}

func (s *mosSuite) TestCpx() {
	s.testCompare(&s.cpu.X, mos.Cpx)
}

func (s *mosSuite) TestCpy() {
	s.testCompare(&s.cpu.Y, mos.Cpy)
}

func (s *mosSuite) testDecrement(
	funcName string,
	setVal func(*mos.CPU, uint8),
	getVal func(*mos.CPU) uint8,
	fn func(*mos.CPU),
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

		s.Equal(mos.NEGATIVE, s.cpu.P&mos.NEGATIVE)
	})

	s.Run(runZero, func() {
		setVal(s.cpu, d1)
		fn(s.cpu)

		s.Equal(mos.ZERO, s.cpu.P&mos.ZERO)
	})
}

func (s *mosSuite) TestDex() {
	s.testDecrement("DEX",
		func(c *mos.CPU, b uint8) { c.X = b },
		func(c *mos.CPU) uint8 { return c.X },
		mos.Dex,
	)
}

func (s *mosSuite) TestDey() {
	s.testDecrement("DEY",
		func(c *mos.CPU, b uint8) { c.Y = b },
		func(c *mos.CPU) uint8 { return c.Y },
		mos.Dey,
	)
}

// TestDec is a bit tricky, because DEC does two very different things based
// on its address mode.
func (s *mosSuite) TestDec() {
	s.cpu.AddrMode = mos.AmACC
	s.testDecrement("DEC (accumulator)",
		func(c *mos.CPU, b uint8) { c.A = b; c.EffVal = b },
		func(c *mos.CPU) uint8 { return c.A },
		mos.Dec,
	)

	var addr uint16 = 1

	s.cpu.AddrMode = mos.AmABS
	s.testDecrement("DEC (memory)",
		func(c *mos.CPU, b uint8) { c.Set(addr, b); c.EffAddr = addr; c.EffVal = b },
		func(c *mos.CPU) uint8 { return c.Get(addr) },
		mos.Dec,
	)
}

func (s *mosSuite) testIncrement(
	funcName string,
	setVal func(*mos.CPU, uint8),
	getVal func(*mos.CPU) uint8,
	fn func(*mos.CPU),
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

		s.Equal(mos.NEGATIVE, s.cpu.P&mos.NEGATIVE)
	})

	s.Run(runZero, func() {
		setVal(s.cpu, d255)
		fn(s.cpu)

		s.Equal(mos.ZERO, s.cpu.P&mos.ZERO)
	})
}

func (s *mosSuite) TestInx() {
	s.testIncrement("INX",
		func(c *mos.CPU, b uint8) { c.X = b },
		func(c *mos.CPU) uint8 { return c.X },
		mos.Inx,
	)
}

func (s *mosSuite) TestIny() {
	s.testIncrement("INY",
		func(c *mos.CPU, b uint8) { c.Y = b },
		func(c *mos.CPU) uint8 { return c.Y },
		mos.Iny,
	)
}

// Same deal as TestDec above -- the INC instruction does some very different
// things based on address mode.
func (s *mosSuite) TestInc() {
	s.cpu.AddrMode = mos.AmACC
	s.testIncrement("INC",
		func(c *mos.CPU, b uint8) { c.A = b; c.EffVal = b },
		func(c *mos.CPU) uint8 { return c.A },
		mos.Inc,
	)

	var addr uint16 = 1

	s.cpu.AddrMode = mos.AmABS
	s.testIncrement("INC",
		func(c *mos.CPU, b uint8) { c.Set(addr, b); c.EffAddr = addr; c.EffVal = b },
		func(c *mos.CPU) uint8 { return c.Get(addr) },
		mos.Inc,
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
		s.cpu.P = mos.CARRY
		mos.Sbc(s.cpu)

		s.Equal(d127-offset, s.cpu.A)
	})

	s.Run("subtraction with nonzero result sets carry", func() {
		s.Equal(mos.CARRY, s.cpu.P&mos.CARRY)
	})

	s.Run("subtracting without carry sets results in A = A - EV - 1", func() {
		s.cpu.A = d127
		s.cpu.EffVal = offset
		s.cpu.P = 0
		mos.Sbc(s.cpu)

		s.Equal(d127-offset-one, s.cpu.A)
	})

	s.Run("subtracting a larger from a smaller number sets negative", func() {
		s.cpu.A = d10
		s.cpu.EffVal = offset
		s.cpu.P = mos.CARRY
		mos.Sbc(s.cpu)

		s.Equal(d10-offset, s.cpu.A)
		s.Equal(mos.NEGATIVE, s.cpu.P&mos.NEGATIVE)
	})

	s.Run("subtracting that flips the sign positive to negative sets overflow", func() {
		s.Equal(mos.OVERFLOW, s.cpu.P&mos.OVERFLOW)
	})

	s.Run("subtracting that flips the sign negative to positve sets overflow", func() {
		s.cpu.A = d127 + offset
		s.cpu.EffVal = offset
		s.cpu.P = mos.CARRY
		mos.Sbc(s.cpu)

		s.Equal(d127, s.cpu.A)
		s.Equal(mos.OVERFLOW, s.cpu.P&mos.OVERFLOW)
	})

	s.Run("subtracting that results in zero sets zero", func() {
		s.cpu.A = d127
		s.cpu.EffVal = d127
		s.cpu.P = mos.CARRY
		mos.Sbc(s.cpu)

		s.Equal(d0, s.cpu.A)
		s.Equal(mos.ZERO, s.cpu.P&mos.ZERO)
	})

	s.Run("decimal works", func() {
		s.cpu.A = 0x16
		s.cpu.EffVal = 0x7
		s.cpu.P = mos.DECIMAL | mos.CARRY

		mos.Sbc(s.cpu)
		s.Equal(uint8(0x9), s.cpu.A)
	})
}
