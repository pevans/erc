package a2

import "github.com/pevans/erc/statemap"

func (s *a2Suite) TestPCSwitcherUseDefaults() {
	pcUseDefaults(s.comp)
	s.Equal(false, s.comp.state.Bool(statemap.PCExpansion))
	s.Equal(false, s.comp.state.Bool(statemap.PCSlotC3))
	s.Equal(true, s.comp.state.Bool(statemap.PCSlotCX))
}

func (s *a2Suite) TestPCSwitcherSwitchWrite() {
	s.Run("slot c3 rom writes work", func() {
		s.comp.state.SetBool(statemap.PCSlotC3, false)
		pcSwitchWrite(int(0xC00B), 0x0, s.comp.state)
		s.Equal(true, s.comp.state.Bool(statemap.PCSlotC3))

		pcSwitchWrite(int(0xC00A), 0x0, s.comp.state)
		s.Equal(false, s.comp.state.Bool(statemap.PCSlotC3))
	})

	s.Run("slot cx rom writes work", func() {
		s.comp.state.SetBool(statemap.PCSlotCX, false)
		s.comp.state.SetBool(statemap.PCSlotC3, false)
		pcSwitchWrite(int(0xC006), 0x0, s.comp.state)
		s.Equal(true, s.comp.state.Bool(statemap.PCSlotCX))
		s.Equal(true, s.comp.state.Bool(statemap.PCSlotC3))

		pcSwitchWrite(int(0xC007), 0x0, s.comp.state)
		s.Equal(false, s.comp.state.Bool(statemap.PCSlotCX))
		s.Equal(false, s.comp.state.Bool(statemap.PCSlotC3))
	})
}

func (s *a2Suite) TestPCSwitcherSwitchRead() {
	var (
		hi uint8 = 0x80
		lo uint8 = 0x00
	)

	s.Run("read of slotc3 returns hi", func() {
		s.comp.state.SetBool(statemap.PCSlotC3, true)
		s.Equal(hi, pcSwitchRead(int(0xC017), s.comp.state))
	})

	s.Run("read of slot cx returns lo", func() {
		s.comp.state.SetBool(statemap.PCSlotCX, true)
		s.Equal(lo, pcSwitchRead(int(0xC016), s.comp.state))
	})
}

func (s *a2Suite) TestPCRead() {
	var (
		c801   = 0xC801
		uc801  = int(c801)
		prc801 = pcPROMAddr(c801)
		irc801 = pcIROMAddr(c801)
		c301   = 0xC301
		uc301  = int(c301)
		prc301 = pcPROMAddr(c301)
		irc301 = pcIROMAddr(c301)
		c401   = 0xC401
		uc401  = int(c401)
		prc401 = pcPROMAddr(c401)
		irc401 = pcIROMAddr(c401)
	)

	s.Run("reads from expansion space", func() {
		s.comp.state.SetBool(statemap.PCExpansion, true)
		s.Equal(s.comp.ROM.Get(prc801), PCRead(uc801, s.comp.state))

		s.comp.state.SetBool(statemap.PCExpansion, false)
		s.Equal(s.comp.ROM.Get(irc801), PCRead(uc801, s.comp.state))
	})

	s.Run("reads from c3 rom space", func() {
		s.comp.state.SetBool(statemap.PCSlotC3, true)
		s.comp.state.SetBool(statemap.PCSlotCX, false)
		s.Equal(s.comp.ROM.Get(prc301), PCRead(uc301, s.comp.state))

		s.comp.state.SetBool(statemap.PCSlotC3, false)
		s.Equal(s.comp.ROM.Get(irc301), PCRead(uc301, s.comp.state))

		// slotCX is a trick; it enables us to read C3 ROM even if C3 ROM is
		// off.
		s.comp.state.SetBool(statemap.PCSlotCX, true)
		s.Equal(s.comp.ROM.Get(prc301), PCRead(uc301, s.comp.state))
	})

	s.Run("reads from cx rom space", func() {
		s.comp.state.SetBool(statemap.PCSlotCX, true)
		s.Equal(s.comp.ROM.DirectGet(prc401), PCRead(uc401, s.comp.state))

		s.comp.state.SetBool(statemap.PCSlotCX, false)
		s.Equal(s.comp.ROM.DirectGet(irc401), PCRead(uc401, s.comp.state))
	})
}
