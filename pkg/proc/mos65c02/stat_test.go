package mos65c02

import (
	"github.com/pevans/erc/pkg/mach"
	"github.com/stretchr/testify/assert"
)

func (s *mosSuite) TestClc() {
	s.cpu.P = CARRY
	Clc(s.cpu)

	assert.Equal(s.T(), mach.Byte(0), s.cpu.P)
}

func (s *mosSuite) TestCld() {
	s.cpu.P = DECIMAL
	Cld(s.cpu)

	assert.Equal(s.T(), mach.Byte(0), s.cpu.P)
}

func (s *mosSuite) TestCli() {
	s.cpu.P = INTERRUPT
	Cli(s.cpu)

	assert.Equal(s.T(), mach.Byte(0), s.cpu.P)
}

func (s *mosSuite) TestClv() {
	s.cpu.P = OVERFLOW
	Clv(s.cpu)

	assert.Equal(s.T(), mach.Byte(0), s.cpu.P)
}

func (s *mosSuite) TestSec() {
	s.cpu.P = 0
	Sec(s.cpu)

	assert.Equal(s.T(), CARRY, s.cpu.P)
}

func (s *mosSuite) TestSed() {
	s.cpu.P = 0
	Sed(s.cpu)

	assert.Equal(s.T(), DECIMAL, s.cpu.P)
}

func (s *mosSuite) TestSei() {
	s.cpu.P = 0
	Sei(s.cpu)

	assert.Equal(s.T(), INTERRUPT, s.cpu.P)
}
