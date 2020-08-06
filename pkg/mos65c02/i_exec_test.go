package mos65c02

import (
	"github.com/pevans/erc/pkg/data"
	"github.com/stretchr/testify/assert"
)

const (
	execAddr = data.DByte(0x4444)
	execPC   = data.DByte(0x1234)
	execP    = data.Byte(0x12)
)

func (s *mosSuite) TestBrk() {
	s.cpu.PC, s.cpu.P = execPC, execP
	Brk(s.cpu)

	assert.Equal(s.T(), INTERRUPT, s.cpu.P&INTERRUPT)
	assert.Equal(s.T(), execPC+2, s.cpu.PC)

	assert.Equal(s.T(), execP, s.cpu.PopStack())

	lsb := data.DByte(s.cpu.PopStack())
	msb := data.DByte(s.cpu.PopStack())

	assert.Equal(s.T(), execPC, (msb<<8)|lsb)
}

func (s *mosSuite) TestJmp() {
	s.cpu.EffAddr = execAddr
	s.cpu.PC = execPC
	Jmp(s.cpu)

	assert.Equal(s.T(), execAddr, s.cpu.PC)
}

func (s *mosSuite) TestJsr() {
	s.cpu.PC = execPC
	s.cpu.EffAddr = execAddr
	Jsr(s.cpu)

	assert.Equal(s.T(), execAddr, s.cpu.PC)

	lsb := data.DByte(s.cpu.PopStack())
	msb := data.DByte(s.cpu.PopStack())

	assert.Equal(s.T(), execPC+2, (msb<<8)|lsb)
}

func (s *mosSuite) TestNops() {
	// These functions do nothing, and this test also does nothing.
	Nop(s.cpu)
	Np2(s.cpu)
	Np3(s.cpu)
}

func (s *mosSuite) TestRti() {
	msb := data.Byte(execPC >> 8)
	lsb := data.Byte(execPC & 0xFF)

	s.cpu.PushStack(msb)
	s.cpu.PushStack(lsb)
	s.cpu.PushStack(execP)

	Rti(s.cpu)

	assert.Equal(s.T(), execP, s.cpu.P)
	assert.Equal(s.T(), execPC, s.cpu.PC)
}

func (s *mosSuite) TestRts() {
	pc := execPC + 2
	msb := data.Byte(pc >> 8)
	lsb := data.Byte(pc & 0xFF)

	s.cpu.PushStack(msb)
	s.cpu.PushStack(lsb)

	Rts(s.cpu)

	assert.Equal(s.T(), pc+1, s.cpu.PC)
}
