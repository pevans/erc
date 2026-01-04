package a2

import "github.com/pevans/erc/a2/a2state"

func (s *a2Suite) TestPCSwitcherUseDefaults() {
	pcUseDefaults(s.comp)
	s.Equal(false, s.comp.State.Bool(a2state.PCExpansion))
	s.Equal(false, s.comp.State.Bool(a2state.PCSlotC3))
	s.Equal(true, s.comp.State.Bool(a2state.PCSlotCX))
}

func (s *a2Suite) TestPCSwitcherSwitchWrite() {
	s.Run("slot c3 rom writes work", func() {
		s.comp.State.SetBool(a2state.PCSlotC3, false)
		pcSwitchWrite(int(0xC00B), 0x0, s.comp.State)
		s.Equal(true, s.comp.State.Bool(a2state.PCSlotC3))

		pcSwitchWrite(int(0xC00A), 0x0, s.comp.State)
		s.Equal(false, s.comp.State.Bool(a2state.PCSlotC3))
	})

	s.Run("slot cx rom writes work", func() {
		s.comp.State.SetBool(a2state.PCSlotCX, false)
		s.comp.State.SetBool(a2state.PCSlotC3, false)
		pcSwitchWrite(int(0xC006), 0x0, s.comp.State)
		s.Equal(true, s.comp.State.Bool(a2state.PCSlotCX))
		s.Equal(false, s.comp.State.Bool(a2state.PCSlotC3)) // CX switch should NOT affect C3

		pcSwitchWrite(int(0xC007), 0x0, s.comp.State)
		s.Equal(false, s.comp.State.Bool(a2state.PCSlotCX))
		s.Equal(false, s.comp.State.Bool(a2state.PCSlotC3)) // CX switch should NOT affect C3
	})
}

func (s *a2Suite) TestPCSwitcherSwitchRead() {
	var (
		hi uint8 = 0x80
		lo uint8 = 0x00
	)

	s.Run("read of slotc3 returns hi", func() {
		s.comp.State.SetBool(a2state.PCSlotC3, true)
		s.Equal(hi, pcSwitchRead(int(0xC017), s.comp.State))
	})

	s.Run("read of slot cx returns lo", func() {
		s.comp.State.SetBool(a2state.PCSlotCX, true)
		s.Equal(lo, pcSwitchRead(int(0xC016), s.comp.State))
	})
}

func (s *a2Suite) TestPCRead() {
	var (
		c301   = 0xC301
		uc301  = int(c301)
		prc301 = pcPROMAddr(c301)
		irc301 = pcIROMAddr(c301)
		c401   = 0xC401
		uc401  = int(c401)
		prc401 = pcPROMAddr(c401)
		irc401 = pcIROMAddr(c401)
	)

	s.Run("reads from c3 rom space", func() {
		s.comp.State.SetBool(a2state.PCSlotC3, true)
		s.comp.State.SetBool(a2state.PCSlotCX, false)
		s.Equal(s.comp.ROM.Get(prc301), PCRead(uc301, s.comp.State))

		s.comp.State.SetBool(a2state.PCSlotC3, false)
		s.Equal(s.comp.ROM.Get(irc301), PCRead(uc301, s.comp.State))

		s.comp.State.SetBool(a2state.PCSlotCX, true)
		s.Equal(s.comp.ROM.Get(irc301), PCRead(uc301, s.comp.State))
	})

	s.Run("reads from cx rom space", func() {
		s.comp.State.SetBool(a2state.PCSlotCX, true)
		s.Equal(s.comp.ROM.DirectGet(prc401), PCRead(uc401, s.comp.State))

		s.comp.State.SetBool(a2state.PCSlotCX, false)
		s.Equal(s.comp.ROM.DirectGet(irc401), PCRead(uc401, s.comp.State))
	})
}
