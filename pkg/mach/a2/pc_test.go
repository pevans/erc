package a2

import (
	"github.com/pevans/erc/pkg/mach"
	"github.com/stretchr/testify/assert"
)

// This is kind of a nasty function, in that it handles a lot of
// different cases. Hence the rather large table of test cases.
func (s *a2Suite) TestPCROMAddr() {
	cases := []struct {
		addr   mach.DByte
		pcMode int
		want   mach.DByte
	}{
		{0xC000, 0, 0},
		{0xC000, PCExpROM, 0},
		{0xC800, PCExpROM, 0x4800},
		{0xCFFF, PCExpROM, 0x4FFF},
		{0xC000, PCSlotCxROM, 0},
		{0xC100, PCSlotCxROM, 0x4100},
		{0xC7FF, PCSlotCxROM, 0x47FF},
		{0xC000, PCSlotC3ROM, 0},
		{0xC300, PCSlotC3ROM, 0x4300},
		{0xC3FF, PCSlotC3ROM, 0x43FF},
		{0xC800, PCExpROM | PCSlotCxROM, 0x4800},
		{0xC800, PCExpROM | PCSlotC3ROM, 0x4800},
		{0xC100, PCExpROM | PCSlotCxROM, 0x4100},
		{0xC100, PCExpROM | PCSlotC3ROM, 0x100},
		{0xC300, PCExpROM | PCSlotCxROM, 0x4300},
		{0xC300, PCExpROM | PCSlotC3ROM, 0x4300},
	}

	for _, c := range cases {
		assert.Equal(s.T(), c.want, pcROMAddr(c.addr, c.pcMode))
	}
}

func (s *a2Suite) TestPCRead() {
	cases := []struct {
		addr mach.DByte
		want mach.Byte
	}{
		{0xC111, 123},
		{0xC222, 223},
	}

	for _, c := range cases {
		s.comp.ROM.Set(c.addr-0xC000, c.want)
		assert.Equal(s.T(), c.want, pcRead(s.comp, c.addr))
	}
}

func (s *a2Suite) TestPCWrite() {
	cases := []struct {
		addr mach.DByte
		want mach.Byte
	}{
		{0xC111, 123},
		{0xC222, 223},
	}

	for _, c := range cases {
		pcWrite(s.comp, c.addr, c.want)

		// We test here that the value is NOT equal, because pcWrite()
		// should prevent any writes to ROM.
		assert.NotEqual(s.T(), c.want, s.comp.ROM.Get(c.addr-0xC000))
	}
}

func (s *a2Suite) TestNewPCSwitchCheck() {
	assert.NotEqual(s.T(), nil, newPCSwitchCheck())
}

func (s *a2Suite) TestPCMode() {
	s.comp.PCMode = 123
	assert.Equal(s.T(), 123, pcMode(s.comp))

	s.comp.PCMode = 124
	assert.Equal(s.T(), 124, pcMode(s.comp))
}

func (s *a2Suite) TestPCSetMode() {
	pcSetMode(s.comp, 123)
	assert.Equal(s.T(), 123, s.comp.PCMode)

	pcSetMode(s.comp, 124)
	assert.Equal(s.T(), 124, s.comp.PCMode)
}
