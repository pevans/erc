package mos65c02

import (
	"github.com/pevans/erc/pkg/data"
	"github.com/stretchr/testify/assert"
)

func (s *mosSuite) TestSaveResult() {
	cases := []struct {
		mode int
		addr data.DByte
		want data.Byte
	}{
		{amAcc, 0x1234, 0x12},
		{amAbs, 0x1234, 0x34},
	}

	for _, c := range cases {
		s.cpu.AddrMode = c.mode
		s.cpu.EffAddr = c.addr

		s.cpu.A = 0
		s.cpu.Set(s.cpu.EffAddr, 0)
		s.cpu.saveResult(c.want)

		switch c.mode {
		case amAcc:
			assert.Equal(s.T(), c.want, s.cpu.A)
		case amAbs:
			assert.Equal(s.T(), c.want, s.cpu.Get(s.cpu.EffAddr))
		}
	}
}

type bitCase struct {
	a     data.Byte
	in    data.Byte
	want  data.Byte
	initp data.Byte
	p     data.Byte
}

// And implements the AND instruction, which performs a bitwise-and on A
// and the effective value and saves the result there.
func (s *mosSuite) TestAnd() {
	cases := []bitCase{
		{0x4, 0x7, 0x4, 0, 0},
		{0x80, 0x81, 0x80, 0, NEGATIVE},
		{0x80, 0x7F, 0x0, 0, ZERO},
	}

	for _, c := range cases {
		s.cpu.A = c.a
		s.cpu.EffVal = c.in
		s.cpu.P = 0

		And(s.cpu)
		assert.Equal(s.T(), c.want, s.cpu.A)
		assert.Equal(s.T(), c.p, s.cpu.P)
	}
}

func (s *mosSuite) TestAsl() {
	cases := []bitCase{
		{0x0, 0x4, 0x8, 0, 0},
		{0x0, 0x40, 0x80, 0, NEGATIVE},
		{0x0, 0x0, 0x0, 0, ZERO},
		{0x0, 0x80, 0x0, 0, ZERO | CARRY},
	}

	for _, c := range cases {
		s.cpu.EffVal = c.in
		s.cpu.AddrMode = amAcc
		s.cpu.P = 0

		Asl(s.cpu)
		assert.Equal(s.T(), c.want, s.cpu.A)
		assert.Equal(s.T(), c.p, s.cpu.P)
	}
}

func (s *mosSuite) TestBit() {
	cases := []bitCase{
		{0x1, 0x3, 0, 0, 0},
		{0x1, 0x2, 0, 0, ZERO},
		{0x0, 0x84, 0, 0, ZERO | NEGATIVE},
		{0x4, 0x84, 0, 0, NEGATIVE},
		{0x4, 0x44, 0, 0, OVERFLOW},
		{0x4, 0xC4, 0, 0, NEGATIVE | OVERFLOW},
	}

	for _, c := range cases {
		s.cpu.A = c.a
		s.cpu.EffVal = c.in
		s.cpu.P = c.initp

		Bit(s.cpu)
		assert.Equal(s.T(), c.p, s.cpu.P)
	}
}

func (s *mosSuite) TestBim() {
	cases := []bitCase{
		{0x1, 0x3, 0, 0, 0},
		{0x1, 0x2, 0, 0, ZERO},
	}

	for _, c := range cases {
		s.cpu.A = c.a
		s.cpu.EffVal = c.in
		s.cpu.P = c.initp

		Bim(s.cpu)
		assert.Equal(s.T(), c.p, s.cpu.P)
	}
}

func (s *mosSuite) TestEor() {
	cases := []bitCase{
		{0x4, 0x7, 0x3, 0, 0},
		{0x81, 0x1, 0x80, 0, NEGATIVE},
		{0x00, 0x00, 0x00, 0, ZERO},
	}

	for _, c := range cases {
		s.cpu.A = c.a
		s.cpu.EffVal = c.in
		s.cpu.P = 0

		Eor(s.cpu)
		assert.Equal(s.T(), c.want, s.cpu.A)
		assert.Equal(s.T(), c.p, s.cpu.P)
	}
}

func (s *mosSuite) TestLsr() {
	// There aren't so many cases here; you might ask "what about
	// negative"? But it's impossible to set the N status with LSR,
	// because there is no form of LSR which can result in a negative
	// number: the eighth bit is always set to zero.
	cases := []bitCase{
		{0x0, 0x4, 0x2, 0, 0},
		{0x0, 0x1, 0x0, 0, ZERO | CARRY},
	}

	for _, c := range cases {
		s.cpu.EffVal = c.in
		s.cpu.AddrMode = amAcc
		s.cpu.P = 0

		Lsr(s.cpu)
		assert.Equal(s.T(), c.want, s.cpu.A)
		assert.Equal(s.T(), c.p, s.cpu.P)
	}
}

func (s *mosSuite) TestOra() {
	cases := []bitCase{
		{0x4, 0x7, 0x7, 0, 0},
		{0x80, 0x81, 0x81, 0, NEGATIVE},
		{0x00, 0x00, 0x00, 0, ZERO},
	}

	for _, c := range cases {
		s.cpu.A = c.a
		s.cpu.EffVal = c.in
		s.cpu.P = 0

		Ora(s.cpu)
		assert.Equal(s.T(), c.want, s.cpu.A)
		assert.Equal(s.T(), c.p, s.cpu.P)
	}
}

func (s *mosSuite) TestRol() {
	cases := []bitCase{
		{0x0, 0x4, 0x8, 0, 0},
		{0x0, 0x4, 0x9, CARRY, 0},
		{0x0, 0x84, 0x8, 0, CARRY},
		{0x0, 0x80, 0x0, 0, ZERO | CARRY},
		{0x0, 0x0, 0x0, 0, ZERO},
		{0x0, 0x40, 0x80, 0, NEGATIVE},
	}

	for _, c := range cases {
		s.cpu.A = c.a
		s.cpu.P = c.initp
		s.cpu.EffVal = c.in
		s.cpu.AddrMode = amAcc

		Rol(s.cpu)
		assert.Equal(s.T(), c.want, s.cpu.A)
		assert.Equal(s.T(), c.p, s.cpu.P)
	}
}

func (s *mosSuite) TestRor() {
	cases := []bitCase{
		{0x0, 0x4, 0x2, 0, 0},
		{0x0, 0x81, 0x40, 0, CARRY},
		{0x0, 0x1, 0x0, 0, ZERO | CARRY},
		{0x0, 0x0, 0x0, 0, ZERO},
		// Note that there is no scenario in which a pre-existing CARRY
		// status does not cause us to then set the NEGATIVE flag after
		// an ROR
		{0x0, 0x2, 0x81, CARRY, NEGATIVE},
	}

	for _, c := range cases {
		s.cpu.A = c.a
		s.cpu.P = c.initp
		s.cpu.EffVal = c.in
		s.cpu.AddrMode = amAcc

		Ror(s.cpu)
		assert.Equal(s.T(), c.want, s.cpu.A)
		assert.Equal(s.T(), c.p, s.cpu.P)
	}
}

func (s *mosSuite) TestTrb() {
	addr := data.DByte(0x1234)
	cases := []bitCase{
		{0x1, 0x3, 0, 0, 0},
		{0x1, 0x2, 0, 0, ZERO},
	}

	for _, c := range cases {
		s.cpu.A = c.a
		s.cpu.EffVal = c.in
		s.cpu.EffAddr = addr
		s.cpu.P = c.initp

		Trb(s.cpu)
		assert.Equal(s.T(), c.p, s.cpu.P)
		assert.Equal(s.T(), s.cpu.Get(addr), (s.cpu.A^0xFF)&s.cpu.EffVal)
	}
}

func (s *mosSuite) TestTsb() {
	addr := data.DByte(0x1234)
	cases := []bitCase{
		{0x1, 0x3, 0, 0, 0},
		{0x1, 0x2, 0, 0, ZERO},
	}

	for _, c := range cases {
		s.cpu.A = c.a
		s.cpu.EffVal = c.in
		s.cpu.EffAddr = addr
		s.cpu.P = c.initp

		Tsb(s.cpu)
		assert.Equal(s.T(), c.p, s.cpu.P)
		assert.Equal(s.T(), s.cpu.Get(addr), s.cpu.A|s.cpu.EffVal)
	}
}
