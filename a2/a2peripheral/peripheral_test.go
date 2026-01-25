package a2peripheral

import (
	"testing"

	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/memory"
	"github.com/stretchr/testify/suite"
)

type peripheralSuite struct {
	suite.Suite

	state *memory.StateMap
	rom   *memory.Segment
}

func (s *peripheralSuite) SetupTest() {
	s.rom = memory.NewSegment(0x5000)
	s.state = memory.NewStateMap()
	s.state.SetSegment(a2state.PCROMSegment, s.rom)
}

func TestPeriphSuite(t *testing.T) {
	suite.Run(t, new(peripheralSuite))
}

func (s *peripheralSuite) TestUseDefaults() {
	UseDefaults(s.state, s.rom)
	s.False(s.state.Bool(a2state.PCExpansion))
	s.False(s.state.Bool(a2state.PCSlotC3))
	s.True(s.state.Bool(a2state.PCSlotCX))
}

func (s *peripheralSuite) TestSwitchWrite() {
	s.Run("slot c3 rom writes work", func() {
		s.state.SetBool(a2state.PCSlotC3, false)
		SwitchWrite(int(0xC00B), 0x0, s.state)
		s.True(s.state.Bool(a2state.PCSlotC3))

		SwitchWrite(int(0xC00A), 0x0, s.state)
		s.False(s.state.Bool(a2state.PCSlotC3))
	})

	s.Run("slot cx rom writes work", func() {
		s.state.SetBool(a2state.PCSlotCX, false)
		s.state.SetBool(a2state.PCSlotC3, false)
		SwitchWrite(int(0xC006), 0x0, s.state)
		s.True(s.state.Bool(a2state.PCSlotCX))
		s.False(s.state.Bool(a2state.PCSlotC3)) // CX switch should NOT affect C3

		SwitchWrite(int(0xC007), 0x0, s.state)
		s.False(s.state.Bool(a2state.PCSlotCX))
		s.False(s.state.Bool(a2state.PCSlotC3)) // CX switch should NOT affect C3
	})
}

func (s *peripheralSuite) TestSwitchRead() {
	var (
		hi uint8 = 0x80
		lo uint8 = 0x00
	)

	s.Run("read of slotc3 returns hi", func() {
		s.state.SetBool(a2state.PCSlotC3, true)
		s.Equal(hi, SwitchRead(int(0xC017), s.state))
	})

	s.Run("read of slot cx returns lo", func() {
		s.state.SetBool(a2state.PCSlotCX, true)
		s.Equal(lo, SwitchRead(int(0xC015), s.state))
	})
}

func (s *peripheralSuite) TestRead() {
	var (
		c301   = 0xC301
		uc301  = int(c301)
		prc301 = promAddr(c301)
		irc301 = iromAddr(c301)
		c401   = 0xC401
		uc401  = int(c401)
		prc401 = promAddr(c401)
		irc401 = iromAddr(c401)
	)

	s.Run("reads from c3 rom space", func() {
		s.state.SetBool(a2state.PCSlotC3, true)
		s.state.SetBool(a2state.PCSlotCX, false)
		s.Equal(s.rom.Get(prc301), Read(uc301, s.state))

		s.state.SetBool(a2state.PCSlotC3, false)
		s.Equal(s.rom.Get(irc301), Read(uc301, s.state))

		s.state.SetBool(a2state.PCSlotCX, true)
		s.Equal(s.rom.Get(irc301), Read(uc301, s.state))
	})

	s.Run("reads from cx rom space", func() {
		s.state.SetBool(a2state.PCSlotCX, true)
		s.Equal(s.rom.DirectGet(prc401), Read(uc401, s.state))

		s.state.SetBool(a2state.PCSlotCX, false)
		s.Equal(s.rom.DirectGet(irc401), Read(uc401, s.state))
	})
}
