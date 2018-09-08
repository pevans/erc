package mos65c02

import (
	"testing"

	"github.com/pevans/erc/pkg/mach"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type mosSuite struct {
	suite.Suite

	cpu *CPU
}

func (s *mosSuite) SetupTest() {
	s.cpu = new(CPU)
	s.cpu.RSeg = mach.NewSegment(AddrSpace)
	s.cpu.WSeg = s.cpu.RSeg
}

func TestMosSuite(t *testing.T) {
	suite.Run(t, new(mosSuite))
}

func (s *mosSuite) TestAcc() {
	cases := []struct {
		want mach.Byte
	}{
		{123},
		{0xFF},
		{0x00},
	}

	for _, c := range cases {
		s.cpu.A = c.want
		Acc(s.cpu)

		assert.Equal(s.T(), mach.DByte(0), s.cpu.EffAddr)
		assert.Equal(s.T(), c.want, s.cpu.EffVal)
	}
}

func (s *mosSuite) TestAbs() {
	cases := []struct {
		oper mach.DByte
		want mach.Byte
	}{
		{0x1234, 0xFB},
		{0x6012, 0x33},
		{0xFE01, 0x11},
	}

	for _, c := range cases {
		s.cpu.Set16(s.cpu.PC+1, c.oper)
		s.cpu.Set(c.oper, c.want)

		Abs(s.cpu)

		assert.Equal(s.T(), c.oper, s.cpu.EffAddr)
		assert.Equal(s.T(), c.want, s.cpu.EffVal)
	}
}

func (s *mosSuite) TestAbx() {
	cases := []struct {
		oper mach.DByte
		x    mach.Byte
		want mach.Byte
	}{
		{0x1234, 0x11, 0xFB},
		{0x6012, 0x21, 0x33},
		{0xFE01, 0x31, 0x11},
	}

	for _, c := range cases {
		s.cpu.Set16(s.cpu.PC+1, c.oper)
		s.cpu.Set(c.oper+mach.DByte(c.x), c.want)

		s.cpu.X = c.x
		Abx(s.cpu)

		assert.Equal(s.T(), c.oper+mach.DByte(c.x), s.cpu.EffAddr)
		assert.Equal(s.T(), c.want, s.cpu.EffVal)
	}
}

func (s *mosSuite) TestAby() {
	cases := []struct {
		oper mach.DByte
		y    mach.Byte
		want mach.Byte
	}{
		{0x1234, 0x65, 0xFB},
		{0x6012, 0x55, 0x33},
		{0xFE01, 0x45, 0x11},
	}

	for _, c := range cases {
		s.cpu.Set16(s.cpu.PC+1, c.oper)
		s.cpu.Set(c.oper+mach.DByte(c.y), c.want)

		s.cpu.Y = c.y
		Aby(s.cpu)

		assert.Equal(s.T(), c.oper+mach.DByte(c.y), s.cpu.EffAddr)
		assert.Equal(s.T(), c.want, s.cpu.EffVal)
	}
}

func (s *mosSuite) TestImm() {
	cases := []struct {
		want mach.Byte
	}{
		{0x12},
		{0x34},
		{0x56},
	}

	for _, c := range cases {
		s.cpu.Set(s.cpu.PC+1, c.want)

		Imm(s.cpu)

		assert.Equal(s.T(), mach.DByte(0), s.cpu.EffAddr)
		assert.Equal(s.T(), c.want, s.cpu.EffVal)
	}
}

func (s *mosSuite) TestInd() {
	cases := []struct {
		oper mach.DByte
		addr mach.DByte
		want mach.Byte
	}{
		{0x1111, 0x2222, 0xFE},
		{0x3333, 0x4444, 0xEA},
		{0x5555, 0x6666, 0x12},
	}

	for _, c := range cases {
		// Set the operand
		s.cpu.Set16(s.cpu.PC+1, c.oper)

		// Set the pointer value at the operand address
		s.cpu.Set16(c.oper, c.addr)

		// And, finally, the value.
		s.cpu.Set(c.addr, c.want)

		Ind(s.cpu)

		assert.Equal(s.T(), c.addr, s.cpu.EffAddr)
		assert.Equal(s.T(), c.want, s.cpu.EffVal)
	}
}

func (s *mosSuite) TestIdx() {
	cases := []struct {
		oper   mach.Byte
		atOper mach.DByte
		x      mach.Byte
		want   mach.Byte
	}{
		{0x05, 0x3333, 0x03, 0x34},
		{0xD0, 0x4444, 0xFF, 0x33},
	}

	for _, c := range cases {
		// Set the operand `$NN`.
		s.cpu.Set(s.cpu.PC+1, c.oper)

		// And at the operand (+ X).
		s.cpu.Set16(mach.DByte(c.oper+c.x), c.atOper)

		// Finally, the value we want to see.
		s.cpu.Set(c.atOper, c.want)

		s.cpu.X = c.x
		Idx(s.cpu)

		assert.Equal(s.T(), s.cpu.EffAddr, c.atOper)
		assert.Equal(s.T(), s.cpu.EffVal, c.want)
	}
}

func (s *mosSuite) TestIdy() {
	cases := []struct {
		oper   mach.Byte
		atOper mach.DByte
		y      mach.Byte
		want   mach.Byte
	}{
		{0x05, 0x3102, 0x03, 0x34},
		{0xD0, 0x3156, 0xFF, 0x33},
	}

	for _, c := range cases {
		// Set the `$NN` part of the operand
		s.cpu.Set(s.cpu.PC+1, c.oper)

		// Now set the base address we want at `$NN`
		s.cpu.Set16(mach.DByte(c.oper), c.atOper)

		// The value we want to see will be set at the base address + Y
		addr := c.atOper + mach.DByte(c.y)
		s.cpu.Set(addr, c.want)

		// And now we resolve the address.
		s.cpu.Y = c.y
		Idy(s.cpu)

		assert.Equal(s.T(), c.want, s.cpu.EffVal)
		assert.Equal(s.T(), addr, s.cpu.EffAddr)
	}
}

func (s *mosSuite) TestImp() {
	Imp(s.cpu)
	assert.Equal(s.T(), mach.Byte(0), s.cpu.EffVal)
	assert.Equal(s.T(), mach.DByte(0), s.cpu.EffAddr)

	By2(s.cpu)
	assert.Equal(s.T(), mach.Byte(0), s.cpu.EffVal)
	assert.Equal(s.T(), mach.DByte(0), s.cpu.EffAddr)

	By3(s.cpu)
	assert.Equal(s.T(), mach.Byte(0), s.cpu.EffVal)
	assert.Equal(s.T(), mach.DByte(0), s.cpu.EffAddr)
}

func (s *mosSuite) TestRel() {
	cases := []struct {
		pc   mach.DByte
		next mach.Byte
		want mach.DByte
	}{
		{0x00, 0x30, 0x32},
		{0xFF, 0x02, 0x103},
		{0x36, 0xFF, 0x37},
	}

	for _, c := range cases {
		s.cpu.PC = c.pc
		s.cpu.Set(c.pc+1, c.next)

		Rel(s.cpu)
		assert.Equal(s.T(), c.want, s.cpu.EffAddr)
	}
}

func (s *mosSuite) TestZpg() {
	cases := []struct {
		addr mach.Byte
		want mach.Byte
	}{
		{0x30, 82},
		{0x00, 28},
		{0xFF, 34},
	}

	for _, c := range cases {
		// Set `$NN`
		s.cpu.Set(s.cpu.PC+1, c.addr)

		// Set the value for `$NN`
		s.cpu.Set(mach.DByte(c.addr), c.want)

		Zpg(s.cpu)
		assert.Equal(s.T(), mach.DByte(c.addr), s.cpu.EffAddr)
		assert.Equal(s.T(), c.want, s.cpu.EffVal)
	}
}

func (s *mosSuite) TestZpx() {
	cases := []struct {
		oper mach.Byte
		x    mach.Byte
		want mach.Byte
	}{
		{0x30, 0xF, 82},
		{0x83, 0x1, 28},
		{0xFE, 0x5, 34},
	}

	for _, c := range cases {
		addr := mach.DByte(c.oper + c.x)

		// Set `$NN`
		s.cpu.Set(s.cpu.PC+1, c.oper)

		// Set the value at `$NN,X`
		s.cpu.Set(addr, c.want)

		s.cpu.X = c.x
		Zpx(s.cpu)

		assert.Equal(s.T(), addr, s.cpu.EffAddr)
		assert.Equal(s.T(), c.want, s.cpu.EffVal)
	}
}

func (s *mosSuite) TestZpy() {
	cases := []struct {
		oper mach.Byte
		y    mach.Byte
		want mach.Byte
	}{
		{0x30, 0xF, 82},
		{0x84, 0x1, 28},
		{0xFE, 0x5, 34},
	}

	for _, c := range cases {
		addr := mach.DByte(c.oper + c.y)

		// Set the `$NN` part.
		s.cpu.Set(s.cpu.PC+1, c.oper)

		// But set `$NN,Y` to what we want.
		s.cpu.Set(addr, c.want)

		s.cpu.Y = c.y
		Zpy(s.cpu)

		assert.Equal(s.T(), addr, s.cpu.EffAddr)
		assert.Equal(s.T(), c.want, s.cpu.EffVal)
	}
}
