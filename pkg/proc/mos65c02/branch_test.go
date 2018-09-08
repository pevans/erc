package mos65c02

import (
	"github.com/pevans/erc/pkg/mach"
	"github.com/stretchr/testify/assert"
)

func (s *mosSuite) TestJumpIf() {
	cases := []struct {
		in   mach.Byte
		pc   mach.DByte
		addr mach.DByte
		want mach.DByte
	}{
		{0, 0x1111, 0x2222, 0x1113},
		{1, 0x1111, 0x2222, 0x2222},
	}

	for _, c := range cases {
		s.cpu.PC = c.pc
		s.cpu.EffAddr = c.addr
		s.cpu.jumpIf(c.in)

		assert.Equal(s.T(), c.want, s.cpu.PC)
	}
}

var branchCases = []struct {
	on   bool
	pc   mach.DByte
	addr mach.DByte
	want mach.DByte
}{
	{false, 0x1234, 0x1234 + 0x28, 0x1234 + 2},
	{true, 0x1234, 0x1234 + 0x28, 0x1234 + 0x28},
	{true, 0x1234, 0x1234 + 0x88 - 0x100, 0x1234 + 0x88 - 0x100},
}

func (s *mosSuite) TestBcc() {
	for _, c := range branchCases {
		if c.on {
			s.cpu.P = 0
		} else {
			s.cpu.P = CARRY
		}

		s.cpu.PC = c.pc
		s.cpu.EffAddr = c.addr
		Bcc(s.cpu)

		assert.Equal(s.T(), c.want, s.cpu.PC)
	}
}

func (s *mosSuite) TestBcs() {
	for _, c := range branchCases {
		if c.on {
			s.cpu.P = CARRY
		} else {
			s.cpu.P = 0
		}

		s.cpu.PC = c.pc
		s.cpu.EffAddr = c.addr
		Bcs(s.cpu)

		assert.Equal(s.T(), c.want, s.cpu.PC)
	}
}

func (s *mosSuite) TestBeq() {
	for _, c := range branchCases {
		if c.on {
			s.cpu.P = ZERO
		} else {
			s.cpu.P = 0
		}

		s.cpu.PC = c.pc
		s.cpu.EffAddr = c.addr
		Beq(s.cpu)

		assert.Equal(s.T(), c.want, s.cpu.PC)
	}
}

func (s *mosSuite) TestBmi() {
	for _, c := range branchCases {
		if c.on {
			s.cpu.P = NEGATIVE
		} else {
			s.cpu.P = 0
		}

		s.cpu.PC = c.pc
		s.cpu.EffAddr = c.addr
		Bmi(s.cpu)

		assert.Equal(s.T(), c.want, s.cpu.PC)
	}
}

func (s *mosSuite) TestBne() {
	for _, c := range branchCases {
		if c.on {
			s.cpu.P = 0
		} else {
			s.cpu.P = ZERO
		}

		s.cpu.PC = c.pc
		s.cpu.EffAddr = c.addr
		Bne(s.cpu)

		assert.Equal(s.T(), c.want, s.cpu.PC)
	}
}

func (s *mosSuite) TestBpl() {
	for _, c := range branchCases {
		if c.on {
			s.cpu.P = 0
		} else {
			s.cpu.P = NEGATIVE
		}

		s.cpu.PC = c.pc
		s.cpu.EffAddr = c.addr
		Bpl(s.cpu)

		assert.Equal(s.T(), c.want, s.cpu.PC)
	}
}

func (s *mosSuite) TestBra() {
	for _, c := range branchCases {
		// The BRA instruction always works according to the "true"
		// branchCases.
		if !c.on {
			continue
		}

		s.cpu.PC = c.pc
		s.cpu.EffAddr = c.addr
		Bra(s.cpu)

		assert.Equal(s.T(), c.want, s.cpu.PC)
	}
}

func (s *mosSuite) TestBvc() {
	for _, c := range branchCases {
		if c.on {
			s.cpu.P = 0
		} else {
			s.cpu.P = OVERFLOW
		}

		s.cpu.PC = c.pc
		s.cpu.EffAddr = c.addr
		Bvc(s.cpu)

		assert.Equal(s.T(), c.want, s.cpu.PC)
	}
}

func (s *mosSuite) TestBvs() {
	for _, c := range branchCases {
		if c.on {
			s.cpu.P = OVERFLOW
		} else {
			s.cpu.P = 0
		}

		s.cpu.PC = c.pc
		s.cpu.EffAddr = c.addr
		Bvs(s.cpu)

		assert.Equal(s.T(), c.want, s.cpu.PC)
	}
}
