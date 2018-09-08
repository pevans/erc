package mos65c02

import (
	"github.com/pevans/erc/pkg/mach"
	"github.com/stretchr/testify/assert"
)

func (s *mosSuite) TestAdc() {
	cases := []struct {
		a, oper, p, aWant, pWant mach.Byte
	}{
		{0, 1, 0, 1, 0},
		{80, 80, CARRY, 161, NEGATIVE},
		{160, 160, 0, (160 + 160 - 256), CARRY},
	}

	for _, cas := range cases {
		s.cpu.A = cas.a
		s.cpu.EffVal = cas.oper
		s.cpu.P = cas.p

		Adc(s.cpu)
		assert.Equal(s.T(), cas.aWant, s.cpu.A)
		assert.Equal(s.T(), cas.pWant, s.cpu.P)
	}
}

func (s *mosSuite) TestCmp() {
	cases := []struct {
		a, x, y mach.Byte
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
		a, aWant mach.Byte
		x, xWant mach.Byte
		y, yWant mach.Byte
		pWant    mach.Byte
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
}

func (s *mosSuite) TestInxy() {
	cases := []struct {
		fn       func(c *CPU)
		a, aWant mach.Byte
		x, xWant mach.Byte
		y, yWant mach.Byte
		pWant    mach.Byte
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
}

func (s *mosSuite) TestSbc() {
	cases := []struct {
		a, oper, aWant, p, pWant mach.Byte
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
