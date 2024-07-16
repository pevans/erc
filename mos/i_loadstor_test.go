package mos_test

import (
	"github.com/pevans/erc/mos"
	"github.com/stretchr/testify/assert"
)

var loadCases = []struct {
	want uint8
	p    uint8
}{
	{0x00, mos.ZERO},
	{0x01, 0},
	{0x81, mos.NEGATIVE},
}

func (s *mosSuite) TestLda() {
	for _, c := range loadCases {
		s.cpu.EffVal = c.want
		mos.Lda(s.cpu)

		assert.Equal(s.T(), c.p, s.cpu.P)
		assert.Equal(s.T(), c.want, s.cpu.A)
	}
}

func (s *mosSuite) TestLdx() {
	for _, c := range loadCases {
		s.cpu.EffVal = c.want
		mos.Ldx(s.cpu)

		assert.Equal(s.T(), c.p, s.cpu.P)
		assert.Equal(s.T(), c.want, s.cpu.X)
	}
}

func (s *mosSuite) TestLdy() {
	for _, c := range loadCases {
		s.cpu.EffVal = c.want
		mos.Ldy(s.cpu)

		assert.Equal(s.T(), c.p, s.cpu.P)
		assert.Equal(s.T(), c.want, s.cpu.Y)
	}
}

func (s *mosSuite) TestPha() {
	s.cpu.A = 123
	mos.Pha(s.cpu)

	assert.Equal(s.T(), uint8(123), s.cpu.PopStack())
}

func (s *mosSuite) TestPhp() {
	s.cpu.P = 123
	mos.Php(s.cpu)

	assert.Equal(s.T(), uint8(123), s.cpu.PopStack())
}

func (s *mosSuite) TestPhx() {
	s.cpu.X = 123
	mos.Phx(s.cpu)

	assert.Equal(s.T(), uint8(123), s.cpu.PopStack())
}

func (s *mosSuite) TestPhy() {
	s.cpu.Y = 123
	mos.Phy(s.cpu)

	assert.Equal(s.T(), uint8(123), s.cpu.PopStack())
}

func (s *mosSuite) TestPla() {
	s.cpu.PushStack(123)
	mos.Pla(s.cpu)

	assert.Equal(s.T(), uint8(123), s.cpu.A)
}

func (s *mosSuite) TestPlp() {
	s.cpu.PushStack(123)
	mos.Plp(s.cpu)

	assert.Equal(s.T(), uint8(123), s.cpu.P)
}

func (s *mosSuite) TestPlx() {
	s.cpu.PushStack(123)
	mos.Plx(s.cpu)

	assert.Equal(s.T(), uint8(123), s.cpu.X)
}

func (s *mosSuite) TestPly() {
	s.cpu.PushStack(123)
	mos.Ply(s.cpu)

	assert.Equal(s.T(), uint8(123), s.cpu.Y)
}

func (s *mosSuite) TestSta() {
	s.cpu.EffAddr = 123
	s.cpu.A = 234
	mos.Sta(s.cpu)

	assert.Equal(s.T(), uint8(234), s.cpu.Get(s.cpu.EffAddr))
}

func (s *mosSuite) TestStx() {
	s.cpu.EffAddr = 123
	s.cpu.X = 234
	mos.Stx(s.cpu)

	assert.Equal(s.T(), uint8(234), s.cpu.Get(s.cpu.EffAddr))
}

func (s *mosSuite) TestSty() {
	s.cpu.EffAddr = 123
	s.cpu.Y = 234
	mos.Sty(s.cpu)

	assert.Equal(s.T(), uint8(234), s.cpu.Get(s.cpu.EffAddr))
}

func (s *mosSuite) TestStz() {
	s.cpu.EffAddr = 123
	s.cpu.Set(s.cpu.EffAddr, 234)
	mos.Stz(s.cpu)

	assert.Equal(s.T(), uint8(0), s.cpu.Get(s.cpu.EffAddr))
}

func (s *mosSuite) TestTax() {
	for _, c := range loadCases {
		s.cpu.A = c.want
		mos.Tax(s.cpu)

		assert.Equal(s.T(), c.p, s.cpu.P)
		assert.Equal(s.T(), c.want, s.cpu.X)
	}
}

func (s *mosSuite) TestTay() {
	for _, c := range loadCases {
		s.cpu.A = c.want
		mos.Tay(s.cpu)

		assert.Equal(s.T(), c.p, s.cpu.P)
		assert.Equal(s.T(), c.want, s.cpu.Y)
	}
}

func (s *mosSuite) TestTsx() {
	for _, c := range loadCases {
		s.cpu.S = c.want
		mos.Tsx(s.cpu)

		assert.Equal(s.T(), c.p, s.cpu.P)
		assert.Equal(s.T(), c.want, s.cpu.X)
	}
}

func (s *mosSuite) TestTxa() {
	for _, c := range loadCases {
		s.cpu.X = c.want
		mos.Txa(s.cpu)

		assert.Equal(s.T(), c.p, s.cpu.P)
		assert.Equal(s.T(), c.want, s.cpu.A)
	}
}

func (s *mosSuite) TestTxs() {
	for _, c := range loadCases {
		s.cpu.X = c.want
		mos.Txs(s.cpu)

		assert.Equal(s.T(), c.want, s.cpu.S)
	}
}

func (s *mosSuite) TestTya() {
	for _, c := range loadCases {
		s.cpu.Y = c.want
		mos.Tya(s.cpu)

		assert.Equal(s.T(), c.p, s.cpu.P)
		assert.Equal(s.T(), c.want, s.cpu.A)
	}
}
