package mos65c02

import (
	"github.com/pevans/erc/pkg/mach"
	"github.com/stretchr/testify/assert"
)

var loadCases = []struct {
	want mach.Byte
	p    mach.Byte
}{
	{0x00, ZERO},
	{0x01, 0},
	{0x81, NEGATIVE},
}

func (s *mosSuite) TestLda() {
	for _, c := range loadCases {
		s.cpu.EffVal = c.want
		Lda(s.cpu)

		assert.Equal(s.T(), c.p, s.cpu.P)
		assert.Equal(s.T(), c.want, s.cpu.A)
	}
}

func (s *mosSuite) TestLdx() {
	for _, c := range loadCases {
		s.cpu.EffVal = c.want
		Ldx(s.cpu)

		assert.Equal(s.T(), c.p, s.cpu.P)
		assert.Equal(s.T(), c.want, s.cpu.X)
	}
}

func (s *mosSuite) TestLdy() {
	for _, c := range loadCases {
		s.cpu.EffVal = c.want
		Ldy(s.cpu)

		assert.Equal(s.T(), c.p, s.cpu.P)
		assert.Equal(s.T(), c.want, s.cpu.Y)
	}
}

func (s *mosSuite) TestPha() {
	s.cpu.A = 123
	Pha(s.cpu)

	assert.Equal(s.T(), mach.Byte(123), s.cpu.PopStack())
}

func (s *mosSuite) TestPhp() {
	s.cpu.P = 123
	Php(s.cpu)

	assert.Equal(s.T(), mach.Byte(123), s.cpu.PopStack())
}

func (s *mosSuite) TestPhx() {
	s.cpu.X = 123
	Phx(s.cpu)

	assert.Equal(s.T(), mach.Byte(123), s.cpu.PopStack())
}

func (s *mosSuite) TestPhy() {
	s.cpu.Y = 123
	Phy(s.cpu)

	assert.Equal(s.T(), mach.Byte(123), s.cpu.PopStack())
}

func (s *mosSuite) TestPla() {
	s.cpu.PushStack(123)
	Pla(s.cpu)

	assert.Equal(s.T(), mach.Byte(123), s.cpu.A)
}

func (s *mosSuite) TestPlp() {
	s.cpu.PushStack(123)
	Plp(s.cpu)

	assert.Equal(s.T(), mach.Byte(123), s.cpu.P)
}

func (s *mosSuite) TestPlx() {
	s.cpu.PushStack(123)
	Plx(s.cpu)

	assert.Equal(s.T(), mach.Byte(123), s.cpu.X)
}

func (s *mosSuite) TestPly() {
	s.cpu.PushStack(123)
	Ply(s.cpu)

	assert.Equal(s.T(), mach.Byte(123), s.cpu.Y)
}

func (s *mosSuite) TestSta() {
	s.cpu.EffAddr = 123
	s.cpu.A = 234
	Sta(s.cpu)

	assert.Equal(s.T(), mach.Byte(234), s.cpu.Get(s.cpu.EffAddr))
}

func (s *mosSuite) TestStx() {
	s.cpu.EffAddr = 123
	s.cpu.X = 234
	Stx(s.cpu)

	assert.Equal(s.T(), mach.Byte(234), s.cpu.Get(s.cpu.EffAddr))
}

func (s *mosSuite) TestSty() {
	s.cpu.EffAddr = 123
	s.cpu.Y = 234
	Sty(s.cpu)

	assert.Equal(s.T(), mach.Byte(234), s.cpu.Get(s.cpu.EffAddr))
}

func (s *mosSuite) TestStz() {
	s.cpu.EffAddr = 123
	s.cpu.Set(s.cpu.EffAddr, 234)
	Stz(s.cpu)

	assert.Equal(s.T(), mach.Byte(0), s.cpu.Get(s.cpu.EffAddr))
}

func (s *mosSuite) TestTax() {
	for _, c := range loadCases {
		s.cpu.A = c.want
		Tax(s.cpu)

		assert.Equal(s.T(), c.p, s.cpu.P)
		assert.Equal(s.T(), c.want, s.cpu.X)
	}
}

func (s *mosSuite) TestTay() {
	for _, c := range loadCases {
		s.cpu.A = c.want
		Tay(s.cpu)

		assert.Equal(s.T(), c.p, s.cpu.P)
		assert.Equal(s.T(), c.want, s.cpu.Y)
	}
}

func (s *mosSuite) TestTsx() {
	for _, c := range loadCases {
		s.cpu.S = c.want
		Tsx(s.cpu)

		assert.Equal(s.T(), c.p, s.cpu.P)
		assert.Equal(s.T(), c.want, s.cpu.X)
	}
}

func (s *mosSuite) TestTxa() {
	for _, c := range loadCases {
		s.cpu.X = c.want
		Txa(s.cpu)

		assert.Equal(s.T(), c.p, s.cpu.P)
		assert.Equal(s.T(), c.want, s.cpu.A)
	}
}

func (s *mosSuite) TestTxs() {
	for _, c := range loadCases {
		s.cpu.X = c.want
		Txs(s.cpu)

		assert.Equal(s.T(), c.p, s.cpu.P)
		assert.Equal(s.T(), c.want, s.cpu.S)
	}
}

func (s *mosSuite) TestTya() {
	for _, c := range loadCases {
		s.cpu.Y = c.want
		Tya(s.cpu)

		assert.Equal(s.T(), c.p, s.cpu.P)
		assert.Equal(s.T(), c.want, s.cpu.A)
	}
}
