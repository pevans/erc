package mos_test

import "github.com/pevans/erc/mos"

func (s *mosSuite) TestClc() {
	s.op(mos.Clc, with{p: mos.CARRY})
	s.NotEqual(mos.CARRY, s.cpu.P&mos.CARRY)
}

func (s *mosSuite) TestCld() {
	s.op(mos.Cld, with{p: mos.DECIMAL})
	s.NotEqual(mos.DECIMAL, s.cpu.P&mos.DECIMAL)
}

func (s *mosSuite) TestCli() {
	s.op(mos.Cli, with{p: mos.INTERRUPT})
	s.NotEqual(mos.INTERRUPT, s.cpu.P&mos.INTERRUPT)
}

func (s *mosSuite) TestClv() {
	s.op(mos.Clv, with{p: mos.OVERFLOW})
	s.NotEqual(mos.OVERFLOW, s.cpu.P&mos.OVERFLOW)
}

func (s *mosSuite) TestSec() {
	s.op(mos.Sec, with{})
	s.Equal(mos.CARRY, s.cpu.P&mos.CARRY)
}

func (s *mosSuite) TestSed() {
	s.op(mos.Sed, with{})
	s.Equal(mos.DECIMAL, s.cpu.P&mos.DECIMAL)
}

func (s *mosSuite) TestSei() {
	s.op(mos.Sei, with{})
	s.Equal(mos.INTERRUPT, s.cpu.P&mos.INTERRUPT)
}
