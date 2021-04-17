package mos65c02

import (
	"github.com/pevans/erc/pkg/data"
	"github.com/stretchr/testify/assert"
)

func (s *mosSuite) TestAdc() {
	var (
		d10    data.Byte = 10
		d250   data.Byte = 250
		d127   data.Byte = 127
		offset data.Byte = 10
		one    data.Byte = 1
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

func (s *mosSuite) TestCmp() {
	cases := []struct {
		a, x, y data.Byte
		fn      func(c *CPU)
	}{
		{5, 0, 0, Cmp},
		{0, 5, 0, Cpx},
		{0, 0, 5, Cpy},
	}

	for _, cas := range cases {
		s.cpu.A, s.cpu.X, s.cpu.Y = cas.a, cas.x, cas.y
		s.cpu.EffVal = 3
		s.cpu.P = 0

		cas.fn(s.cpu)
		assert.Equal(s.T(), CARRY, s.cpu.P)
	}
}

func (s *mosSuite) TestDexy() {
	cases := []struct {
		fn       func(c *CPU)
		a, aWant data.Byte
		x, xWant data.Byte
		y, yWant data.Byte
		pWant    data.Byte
	}{
		{Dec, 3, 2, 0, 0, 0, 0, 0},
		{Dec, 1, 0, 0, 0, 0, 0, ZERO},
		{Dec, 0, 0xFF, 0, 0, 0, 0, NEGATIVE},
		{Dex, 0, 0, 3, 2, 0, 0, 0},
		{Dex, 0, 0, 1, 0, 0, 0, ZERO},
		{Dex, 0, 0, 0, 255, 0, 0, NEGATIVE},
		{Dey, 0, 0, 0, 0, 3, 2, 0},
		{Dey, 0, 0, 0, 0, 1, 0, ZERO},
		{Dey, 0, 0, 0, 0, 0, 255, NEGATIVE},
	}

	s.cpu.AddrMode = amAcc

	for _, cas := range cases {
		s.cpu.A = cas.a
		s.cpu.X = cas.x
		s.cpu.Y = cas.y

		Acc(s.cpu)
		cas.fn(s.cpu)

		assert.Equal(s.T(), cas.aWant, s.cpu.A)
		assert.Equal(s.T(), cas.xWant, s.cpu.X)
		assert.Equal(s.T(), cas.yWant, s.cpu.Y)
	}

	// We still need to see if we can decrement a spot in memory
	s.cpu.Set16(s.cpu.PC+1, data.DByte(0x1122))
	s.cpu.Set(data.DByte(0x1122), data.Byte(0x34))
	Abs(s.cpu)
	Dec(s.cpu)
	s.Equal(data.Byte(0x33), s.cpu.Get(data.DByte(0x1122)))
}

func (s *mosSuite) TestInxy() {
	cases := []struct {
		fn       func(c *CPU)
		a, aWant data.Byte
		x, xWant data.Byte
		y, yWant data.Byte
		pWant    data.Byte
	}{
		{Inc, 11, 12, 2, 2, 0, 0, 0},
		{Inc, 0xFF, 0, 2, 2, 0, 0, ZERO},
		{Inc, 0x7F, 0x80, 2, 2, 0, 0, NEGATIVE},
		{Inx, 0, 0, 2, 3, 0, 0, 0},
		{Inx, 0, 0, 255, 0, 0, 0, ZERO},
		{Inx, 0, 0, 127, 128, 0, 0, NEGATIVE},
		{Iny, 0, 0, 0, 0, 2, 3, 0},
		{Iny, 0, 0, 0, 0, 255, 0, ZERO},
		{Iny, 0, 0, 0, 0, 127, 128, NEGATIVE},
	}

	s.cpu.AddrMode = amAcc

	for _, cas := range cases {
		s.cpu.X = cas.x
		s.cpu.Y = cas.y
		s.cpu.A = cas.a

		// This is a bit janky, but INC needs to know that the A
		// register is loaded into EffVal for INC with amAcc to work.
		Acc(s.cpu)
		cas.fn(s.cpu)

		assert.Equal(s.T(), cas.aWant, s.cpu.A)
		assert.Equal(s.T(), cas.xWant, s.cpu.X)
		assert.Equal(s.T(), cas.yWant, s.cpu.Y)
	}

	// One final test -- we want to test that we can set a value in some
	// other spot in memory
	s.cpu.Set16(s.cpu.PC+1, data.DByte(0x1122))
	s.cpu.Set(data.DByte(0x1122), data.Byte(0x34))
	Abs(s.cpu)
	Inc(s.cpu)
	s.Equal(data.Byte(0x35), s.cpu.Get(data.DByte(0x1122)))
}

func (s *mosSuite) TestSbc() {
	cases := []struct {
		a, oper, aWant, p, pWant data.Byte
	}{
		{3, 2, 1, CARRY, CARRY},
		{3, 2, 0, 0, CARRY | ZERO},
		{3, 3, 0, CARRY, CARRY | ZERO},
		{3, 4, 255, CARRY, NEGATIVE},
	}

	for _, cas := range cases {
		s.cpu.A = cas.a
		s.cpu.P = cas.p
		s.cpu.EffVal = cas.oper

		Sbc(s.cpu)
		assert.Equal(s.T(), cas.aWant, s.cpu.A)
		assert.Equal(s.T(), cas.pWant, s.cpu.P)
	}
}
