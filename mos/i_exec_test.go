package mos_test

import "github.com/pevans/erc/mos"

const (
	execAddr = uint16(0x4444)
	execPC   = uint16(0x1234)
	execP    = uint8(0x12)
)

func (s *mosSuite) TestBrk() {
	s.Run("sets interrupt flag", func() {
		s.op(mos.Brk, with{pc: execPC, p: execP})
		s.Equal(mos.INTERRUPT, s.cpu.P&mos.INTERRUPT)
	})

	s.Run("sets PC register to new PC + 2", func() {
		s.Equal(execPC+2, s.cpu.PC)
	})

	s.Run("new p register value is on the stack", func() {
		s.Equal(execP, s.cpu.PopStack())
	})

	s.Run("new PC is on the stack", func() {
		lsb := uint16(s.cpu.PopStack())
		msb := uint16(s.cpu.PopStack())

		s.Equal(execPC, (msb<<8)|lsb)
	})
}

func (s *mosSuite) TestJmp() {
	s.Run("moves PC register to the new location", func() {
		s.op(mos.Jmp, with{pc: execPC, addr: execAddr})
		s.Equal(execAddr, s.cpu.PC)
	})
}

func (s *mosSuite) TestJsr() {
	s.Run("moves PC register to the new location", func() {
		s.op(mos.Jsr, with{pc: execPC, addr: execAddr})
		s.Equal(execAddr, s.cpu.PC)
	})

	s.Run("adds original PC + 2 to the stack", func() {
		lsb := uint16(s.cpu.PopStack())
		msb := uint16(s.cpu.PopStack())

		s.Equal(execPC+2, (msb<<8)|lsb)
	})
}

func (s *mosSuite) TestRti() {
	s.Run("sets P and PC registers to values from the stack", func() {
		msb := uint8(execPC >> 8)
		lsb := uint8(execPC & 0xFF)

		s.cpu.PushStack(msb)
		s.cpu.PushStack(lsb)
		s.cpu.PushStack(execP)

		// Have to pass s: s.cpu.S, or else s.op will overwrite the s register
		// value
		s.op(mos.Rti, with{s: s.cpu.S})

		s.Equal(execP, s.cpu.P)
		s.Equal(execPC, s.cpu.PC)
	})
}

func (s *mosSuite) TestRts() {
	s.Run("sets the PC register value to the one from the stack", func() {
		pc := execPC + 2
		msb := uint8(pc >> 8)
		lsb := uint8(pc & 0xFF)

		s.cpu.PushStack(msb)
		s.cpu.PushStack(lsb)

		s.op(mos.Rts, with{s: s.cpu.S})
		s.Equal(pc+1, s.cpu.PC)
	})
}
