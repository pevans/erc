package mos_test

import (
	"testing"

	"github.com/pevans/erc/memory"
	"github.com/pevans/erc/mos"
	"github.com/stretchr/testify/suite"
)

type mosSuite struct {
	suite.Suite

	cpu *mos.CPU
}

func (s *mosSuite) SetupTest() {
	seg := memory.NewSegment(0x10000)
	s.cpu = new(mos.CPU)
	s.cpu.State = new(memory.StateMap)
	s.cpu.RMem = seg
	s.cpu.WMem = seg
}

func TestMosSuite(t *testing.T) {
	suite.Run(t, new(mosSuite))
}

func (s *mosSuite) TestAcc() {
	cases := []struct {
		want uint8
	}{
		{123},
		{0xFF},
		{0x00},
	}

	for _, c := range cases {
		s.cpu.A = c.want
		mos.Acc(s.cpu)

		s.Equal(uint16(0), s.cpu.EffAddr)
		s.Equal(c.want, s.cpu.EffVal)
	}
}

func (s *mosSuite) TestAbs() {
	cases := []struct {
		oper uint16
		want uint8
	}{
		{0x1234, 0xFB},
		{0x6012, 0x33},
		{0xFE01, 0x11},
	}

	for _, c := range cases {
		s.cpu.Set16(s.cpu.PC+1, c.oper)
		s.cpu.Set(c.oper, c.want)

		mos.Abs(s.cpu)

		s.Equal(c.oper, s.cpu.EffAddr)
		s.Equal(c.want, s.cpu.EffVal)
	}
}

func (s *mosSuite) TestAbx() {
	cases := []struct {
		oper uint16
		x    uint8
		want uint8
	}{
		{0x1234, 0x11, 0xFB},
		{0x6012, 0x21, 0x33},
		{0xFE01, 0x31, 0x11},
	}

	for _, c := range cases {
		s.cpu.Set16(s.cpu.PC+1, c.oper)
		s.cpu.Set(c.oper+uint16(c.x), c.want)

		s.cpu.X = c.x
		mos.Abx(s.cpu)

		s.Equal(c.oper+uint16(c.x), s.cpu.EffAddr)
		s.Equal(c.want, s.cpu.EffVal)
	}
}

func (s *mosSuite) TestAby() {
	cases := []struct {
		oper uint16
		y    uint8
		want uint8
	}{
		{0x1234, 0x65, 0xFB},
		{0x6012, 0x55, 0x33},
		{0xFE01, 0x45, 0x11},
	}

	for _, c := range cases {
		s.cpu.Set16(s.cpu.PC+1, c.oper)
		s.cpu.Set(c.oper+uint16(c.y), c.want)

		s.cpu.Y = c.y
		mos.Aby(s.cpu)

		s.Equal(c.oper+uint16(c.y), s.cpu.EffAddr)
		s.Equal(c.want, s.cpu.EffVal)
	}
}

func (s *mosSuite) TestImm() {
	cases := []struct {
		want uint8
	}{
		{0x12},
		{0x34},
		{0x56},
	}

	for _, c := range cases {
		s.cpu.Set(s.cpu.PC+1, c.want)

		mos.Imm(s.cpu)

		s.Equal(uint16(0), s.cpu.EffAddr)
		s.Equal(c.want, s.cpu.EffVal)
	}
}

func (s *mosSuite) TestInd() {
	cases := []struct {
		oper uint16
		addr uint16
		want uint8
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

		mos.Ind(s.cpu)

		s.Equal(c.addr, s.cpu.EffAddr)
		s.Equal(c.want, s.cpu.EffVal)
	}
}

func (s *mosSuite) TestIdx() {
	cases := []struct {
		name   string
		oper   uint8
		x      uint8
		atOper uint16 // The 16-bit pointer value
		want   uint8  // The value at the target address
	}{
		{"normal case", 0x05, 0x03, 0x3333, 0x34},
		{"normal case 2", 0xD0, 0xFF, 0x4444, 0x33},
		{"zero page boundary", 0xFE, 0x01, 0x1234, 0xAB}, // oper+x = 0xFF
	}

	for _, c := range cases {
		s.Run(c.name, func() {
			// Set the operand `$NN`.
			s.cpu.Set(s.cpu.PC+1, c.oper)

			// Set the 16-bit pointer at (oper + x) location
			// We need to handle the zero page boundary case where oper+x = 0xFF
			zpAddr := uint16(c.oper+c.x) & 0xFF
			if zpAddr == 0xFF {
				// Low byte at 0xFF, high byte wraps to 0x00
				s.cpu.Set(0xFF, uint8(c.atOper&0xFF))
				s.cpu.Set(0x00, uint8(c.atOper>>8))
			} else {
				// Normal case: set the 16-bit pointer
				s.cpu.Set16(zpAddr, c.atOper)
			}

			// Set the value we want to see at the target address
			s.cpu.Set(c.atOper, c.want)

			s.cpu.X = c.x
			mos.Idx(s.cpu)

			s.Equal(c.atOper, s.cpu.EffAddr)
			s.Equal(c.want, s.cpu.EffVal)
		})
	}
}

func (s *mosSuite) TestIdy() {
	cases := []struct {
		name     string
		oper     uint8  // The zero page address
		baseAddr uint16 // The 16-bit pointer value at oper
		y        uint8  // Y register value
		want     uint8  // The value at (baseAddr + y)
	}{
		{"normal case", 0x05, 0x3102, 0x03, 0x34},
		{"normal case 2", 0xD0, 0x3156, 0xFF, 0x33},
		{"zero page boundary", 0xFF, 0x2000, 0x10, 0xCD}, // oper = 0xFF
	}

	for _, c := range cases {
		s.Run(c.name, func() {
			// Set the `$NN` part of the operand
			s.cpu.Set(s.cpu.PC+1, c.oper)

			// Set the 16-bit base address pointer at the operand location
			// Handle the zero page boundary case where oper = 0xFF
			if c.oper == 0xFF {
				// Low byte at 0xFF, high byte wraps to 0x00
				s.cpu.Set(0xFF, uint8(c.baseAddr&0xFF))
				s.cpu.Set(0x00, uint8(c.baseAddr>>8))
			} else {
				// Normal case: set the 16-bit pointer
				s.cpu.Set16(uint16(c.oper), c.baseAddr)
			}

			// The final target address is base + Y
			targetAddr := c.baseAddr + uint16(c.y)
			s.cpu.Set(targetAddr, c.want)

			// And now we resolve the address.
			s.cpu.Y = c.y
			mos.Idy(s.cpu)

			s.Equal(c.want, s.cpu.EffVal)
			s.Equal(targetAddr, s.cpu.EffAddr)
		})
	}
}

func (s *mosSuite) TestImp() {
	mos.Imp(s.cpu)
	s.Equal(uint8(0), s.cpu.EffVal)
	s.Equal(uint16(0), s.cpu.EffAddr)

	mos.By2(s.cpu)
	s.Equal(uint8(0), s.cpu.EffVal)
	s.Equal(uint16(0), s.cpu.EffAddr)

	mos.By3(s.cpu)
	s.Equal(uint8(0), s.cpu.EffVal)
	s.Equal(uint16(0), s.cpu.EffAddr)
}

func (s *mosSuite) TestRel() {
	cases := []struct {
		pc   uint16
		next uint8
		want uint16
	}{
		{0x00, 0x30, 0x32},
		{0xFF, 0x02, 0x103},
		{0x36, 0xFF, 0x37},
	}

	for _, c := range cases {
		s.cpu.PC = c.pc
		s.cpu.Set(c.pc+1, c.next)

		mos.Rel(s.cpu)
		s.Equal(c.want, s.cpu.EffAddr)
	}
}

func (s *mosSuite) TestZpg() {
	cases := []struct {
		addr uint8
		want uint8
	}{
		{0x30, 82},
		{0x00, 28},
		{0xFF, 34},
	}

	for _, c := range cases {
		// Set `$NN`
		s.cpu.Set(s.cpu.PC+1, c.addr)

		// Set the value for `$NN`
		s.cpu.Set(uint16(c.addr), c.want)

		mos.Zpg(s.cpu)
		s.Equal(uint16(c.addr), s.cpu.EffAddr)
		s.Equal(c.want, s.cpu.EffVal)
	}
}

func (s *mosSuite) TestZpx() {
	cases := []struct {
		oper uint8
		x    uint8
		want uint8
	}{
		{0x30, 0xF, 82},
		{0x83, 0x1, 28},
		{0xFE, 0x5, 34},
	}

	for _, c := range cases {
		addr := uint16(c.oper + c.x)

		// Set `$NN`
		s.cpu.Set(s.cpu.PC+1, c.oper)

		// Set the value at `$NN,X`
		s.cpu.Set(addr, c.want)

		s.cpu.X = c.x
		mos.Zpx(s.cpu)

		s.Equal(addr, s.cpu.EffAddr)
		s.Equal(c.want, s.cpu.EffVal)
	}
}

func (s *mosSuite) TestZpy() {
	cases := []struct {
		oper uint8
		y    uint8
		want uint8
	}{
		{0x30, 0xF, 82},
		{0x84, 0x1, 28},
		{0xFE, 0x5, 34},
	}

	for _, c := range cases {
		addr := uint16(c.oper + c.y)

		// Set the `$NN` part.
		s.cpu.Set(s.cpu.PC+1, c.oper)

		// But set `$NN,Y` to what we want.
		s.cpu.Set(addr, c.want)

		s.cpu.Y = c.y
		mos.Zpy(s.cpu)

		s.Equal(addr, s.cpu.EffAddr)
		s.Equal(c.want, s.cpu.EffVal)
	}
}
