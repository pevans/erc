package mos_test

import (
	"github.com/pevans/erc/mos"
	"github.com/stretchr/testify/assert"
)

var branchCases = []struct {
	on   bool
	pc   uint16
	addr uint16
	want uint16
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
			s.cpu.P = mos.CARRY
		}

		s.cpu.PC = c.pc
		s.cpu.EffAddr = c.addr
		mos.Bcc(s.cpu)

		assert.Equal(s.T(), c.want, s.cpu.PC)
	}
}

func (s *mosSuite) TestBcs() {
	for _, c := range branchCases {
		if c.on {
			s.cpu.P = mos.CARRY
		} else {
			s.cpu.P = 0
		}

		s.cpu.PC = c.pc
		s.cpu.EffAddr = c.addr
		mos.Bcs(s.cpu)

		assert.Equal(s.T(), c.want, s.cpu.PC)
	}
}

func (s *mosSuite) TestBeq() {
	for _, c := range branchCases {
		if c.on {
			s.cpu.P = mos.ZERO
		} else {
			s.cpu.P = 0
		}

		s.cpu.PC = c.pc
		s.cpu.EffAddr = c.addr
		mos.Beq(s.cpu)

		assert.Equal(s.T(), c.want, s.cpu.PC)
	}
}

func (s *mosSuite) TestBmi() {
	for _, c := range branchCases {
		if c.on {
			s.cpu.P = mos.NEGATIVE
		} else {
			s.cpu.P = 0
		}

		s.cpu.PC = c.pc
		s.cpu.EffAddr = c.addr
		mos.Bmi(s.cpu)

		assert.Equal(s.T(), c.want, s.cpu.PC)
	}
}

func (s *mosSuite) TestBne() {
	for _, c := range branchCases {
		if c.on {
			s.cpu.P = 0
		} else {
			s.cpu.P = mos.ZERO
		}

		s.cpu.PC = c.pc
		s.cpu.EffAddr = c.addr
		mos.Bne(s.cpu)

		assert.Equal(s.T(), c.want, s.cpu.PC)
	}
}

func (s *mosSuite) TestBpl() {
	for _, c := range branchCases {
		if c.on {
			s.cpu.P = 0
		} else {
			s.cpu.P = mos.NEGATIVE
		}

		s.cpu.PC = c.pc
		s.cpu.EffAddr = c.addr
		mos.Bpl(s.cpu)

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
		mos.Bra(s.cpu)

		assert.Equal(s.T(), c.want, s.cpu.PC)
	}
}

func (s *mosSuite) TestBvc() {
	for _, c := range branchCases {
		if c.on {
			s.cpu.P = 0
		} else {
			s.cpu.P = mos.OVERFLOW
		}

		s.cpu.PC = c.pc
		s.cpu.EffAddr = c.addr
		mos.Bvc(s.cpu)

		assert.Equal(s.T(), c.want, s.cpu.PC)
	}
}

func (s *mosSuite) TestBvs() {
	for _, c := range branchCases {
		if c.on {
			s.cpu.P = mos.OVERFLOW
		} else {
			s.cpu.P = 0
		}

		s.cpu.PC = c.pc
		s.cpu.EffAddr = c.addr
		mos.Bvs(s.cpu)

		assert.Equal(s.T(), c.want, s.cpu.PC)
	}
}
