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
	s.Equal(0, s.state.Int(a2state.PCExpSlot))
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
	s.Run("read of slotc3 returns hi with keyboard latch", func() {
		s.state.SetBool(a2state.PCSlotC3, true)
		s.state.SetUint8(a2state.KBLastKey, 0x41)
		s.Equal(uint8(0x80|0x41), SwitchRead(int(0xC017), s.state))
	})

	s.Run("read of slotc3 returns lo with keyboard latch", func() {
		s.state.SetBool(a2state.PCSlotC3, false)
		s.state.SetUint8(a2state.KBLastKey, 0x41)
		s.Equal(uint8(0x41), SwitchRead(int(0xC017), s.state))
	})

	s.Run("read of slot cx returns lo with keyboard latch", func() {
		s.state.SetBool(a2state.PCSlotCX, true)
		s.state.SetUint8(a2state.KBLastKey, 0x41)
		s.Equal(uint8(0x41), SwitchRead(int(0xC015), s.state))
	})

	s.Run("read of slot cx returns hi with keyboard latch", func() {
		s.state.SetBool(a2state.PCSlotCX, false)
		s.state.SetUint8(a2state.KBLastKey, 0x41)
		s.Equal(uint8(0x80|0x41), SwitchRead(int(0xC015), s.state))
	})

	s.Run("read of control switches returns keyboard data latch", func() {
		s.state.SetUint8(a2state.KBLastKey, 0x41)
		s.state.SetUint8(a2state.KBStrobe, 0x80)

		for _, addr := range []int{0xC006, 0xC007, 0xC00A, 0xC00B} {
			s.Equal(uint8(0x80|0x41), SwitchRead(addr, s.state))
		}

		s.state.SetUint8(a2state.KBStrobe, 0x00)

		for _, addr := range []int{0xC006, 0xC007, 0xC00A, 0xC00B} {
			s.Equal(uint8(0x41), SwitchRead(addr, s.state))
		}
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

	s.Run("sets IOSelect and ExpSlot on slot ROM read", func() {
		s.state.SetBool(a2state.PCSlotCX, true)
		s.state.SetBool(a2state.PCIOSelect, false)
		s.state.SetInt(a2state.PCExpSlot, 0)

		Read(uc401, s.state)

		s.True(s.state.Bool(a2state.PCIOSelect))
		s.Equal(4, s.state.Int(a2state.PCExpSlot))
	})

	s.Run("sets IOStrobe on expansion ROM range read", func() {
		s.state.SetBool(a2state.PCSlotCX, true)
		s.state.SetBool(a2state.PCIOSelect, false)
		s.state.SetBool(a2state.PCIOStrobe, false)

		Read(0xC900, s.state)

		s.True(s.state.Bool(a2state.PCIOStrobe))
	})

	s.Run("activates expansion ROM when IOSelect and IOStrobe", func() {
		s.state.SetBool(a2state.PCSlotCX, true)
		s.state.SetBool(a2state.PCIOSelect, true)
		s.state.SetBool(a2state.PCExpansion, false)

		Read(0xC900, s.state)

		s.True(s.state.Bool(a2state.PCIOStrobe))
		s.True(s.state.Bool(a2state.PCExpansion))
	})

	s.Run("does not activate expansion ROM without IOSelect", func() {
		s.state.SetBool(a2state.PCSlotCX, true)
		s.state.SetBool(a2state.PCIOSelect, false)
		s.state.SetBool(a2state.PCExpansion, false)

		Read(0xC900, s.state)

		s.True(s.state.Bool(a2state.PCIOStrobe))
		s.False(s.state.Bool(a2state.PCExpansion))
	})

	s.Run("CFFF read returns expansion ROM when active", func() {
		s.state.SetBool(a2state.PCSlotCX, true)
		s.state.SetBool(a2state.PCIOSelect, true)
		s.state.SetBool(a2state.PCIOStrobe, true)
		s.state.SetBool(a2state.PCExpansion, true)
		s.state.SetInt(a2state.PCExpSlot, 6)

		expected := s.rom.DirectGet(iromAddr(0xCFFF))
		val := Read(0xCFFF, s.state)

		s.Equal(expected, val)
		s.False(s.state.Bool(a2state.PCIOSelect))
		s.False(s.state.Bool(a2state.PCIOStrobe))
		s.False(s.state.Bool(a2state.PCExpansion))
		s.Equal(0, s.state.Int(a2state.PCExpSlot))
	})

	s.Run("CFFF read returns peripheral ROM when SlotCX without IOSelect", func() {
		s.state.SetBool(a2state.PCSlotCX, true)
		s.state.SetBool(a2state.PCIOSelect, false)
		s.state.SetBool(a2state.PCExpansion, false)

		s.rom.DirectSet(iromAddr(0xCFFF), 0xAA)
		s.rom.DirectSet(promAddr(0xCFFF), 0xBB)

		val := Read(0xCFFF, s.state)
		s.Equal(uint8(0xBB), val)
	})

	s.Run("CFFF read returns internal ROM when SlotCX is false", func() {
		s.state.SetBool(a2state.PCSlotCX, false)
		s.state.SetBool(a2state.PCExpansion, false)

		s.rom.DirectSet(iromAddr(0xCFFF), 0xAA)
		s.rom.DirectSet(promAddr(0xCFFF), 0xBB)

		val := Read(0xCFFF, s.state)
		s.Equal(uint8(0xAA), val)
	})
}

func (s *peripheralSuite) TestWrite() {
	s.Run("CFFF write disables expansion", func() {
		s.state.SetBool(a2state.PCIOSelect, true)
		s.state.SetBool(a2state.PCIOStrobe, true)
		s.state.SetBool(a2state.PCExpansion, true)
		s.state.SetInt(a2state.PCExpSlot, 6)

		Write(0xCFFF, 0x00, s.state)

		s.False(s.state.Bool(a2state.PCIOSelect))
		s.False(s.state.Bool(a2state.PCIOStrobe))
		s.False(s.state.Bool(a2state.PCExpansion))
		s.Equal(0, s.state.Int(a2state.PCExpSlot))
	})

	s.Run("non-CFFF write does not disable expansion", func() {
		s.state.SetBool(a2state.PCExpansion, true)
		s.state.SetInt(a2state.PCExpSlot, 6)

		Write(0xC500, 0x00, s.state)

		s.True(s.state.Bool(a2state.PCExpansion))
		s.Equal(6, s.state.Int(a2state.PCExpSlot))
	})
}
