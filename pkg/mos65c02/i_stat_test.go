package mos65c02

func (s *mosSuite) TestClc() {
	s.op(Clc, with{p: CARRY})
	s.NotEqual(CARRY, s.cpu.P&CARRY)
}

func (s *mosSuite) TestCld() {
	s.op(Cld, with{p: DECIMAL})
	s.NotEqual(DECIMAL, s.cpu.P&DECIMAL)
}

func (s *mosSuite) TestCli() {
	s.op(Cli, with{p: INTERRUPT})
	s.NotEqual(INTERRUPT, s.cpu.P&INTERRUPT)
}

func (s *mosSuite) TestClv() {
	s.op(Clv, with{p: OVERFLOW})
	s.NotEqual(OVERFLOW, s.cpu.P&OVERFLOW)
}

func (s *mosSuite) TestSec() {
	s.op(Sec, with{})
	s.Equal(CARRY, s.cpu.P&CARRY)
}

func (s *mosSuite) TestSed() {
	s.op(Sed, with{})
	s.Equal(DECIMAL, s.cpu.P&DECIMAL)
}

func (s *mosSuite) TestSei() {
	s.op(Sei, with{})
	s.Equal(INTERRUPT, s.cpu.P&INTERRUPT)
}
