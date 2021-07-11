package a2

func (s *a2Suite) TestPCSwitcherUseDefaults() {
	var ps pcSwitcher

	ps.UseDefaults()
	s.Equal(false, ps.expansion)
	s.Equal(false, ps.slotC3)
	s.Equal(true, ps.slotCX)
}

func (s *a2Suite) TestPCSwitcherSwitchWrite() {
	var ps pcSwitcher

	s.Run("slot c3 rom writes work", func() {
		ps.slotC3 = false
		ps.SwitchWrite(s.comp, uint16(0xC00B), 0x0)
		s.Equal(true, ps.slotC3)

		ps.SwitchWrite(s.comp, uint16(0xC00A), 0x0)
		s.Equal(false, ps.slotC3)
	})

	s.Run("slot cx rom writes work", func() {
		ps.slotCX = false
		ps.slotC3 = false
		ps.SwitchWrite(s.comp, uint16(0xC006), 0x0)
		s.Equal(true, ps.slotCX)
		s.Equal(true, ps.slotC3)

		ps.SwitchWrite(s.comp, uint16(0xC007), 0x0)
		s.Equal(false, ps.slotCX)
		s.Equal(false, ps.slotC3)
	})
}

func (s *a2Suite) TestPCSwitcherSwitchRead() {
	var (
		ps pcSwitcher
		hi uint8 = 0x80
		lo uint8 = 0x00
	)

	s.Run("read of slotc3 returns hi", func() {
		ps.slotC3 = true
		s.Equal(hi, ps.SwitchRead(s.comp, uint16(0xC017)))
	})

	s.Run("read of slot cx returns lo", func() {
		ps.slotCX = true
		s.Equal(lo, ps.SwitchRead(s.comp, uint16(0xC016)))
	})
}

func (s *a2Suite) TestPCRead() {
	var (
		c801   = 0xC801
		uc801  = uint16(c801)
		prc801 = pcPROMAddr(c801)
		irc801 = pcIROMAddr(c801)
		c301   = 0xC301
		uc301  = uint16(c301)
		prc301 = pcPROMAddr(c301)
		irc301 = pcIROMAddr(c301)
		c401   = 0xC401
		uc401  = uint16(c401)
		prc401 = pcPROMAddr(c401)
		irc401 = pcIROMAddr(c401)
	)

	s.Run("reads from expansion space", func() {
		s.comp.pc.expansion = true
		s.Equal(s.comp.ROM.Get(prc801), PCRead(s.comp, uc801))

		s.comp.pc.expansion = false
		s.Equal(s.comp.ROM.Get(irc801), PCRead(s.comp, uc801))
	})

	s.Run("reads from c3 rom space", func() {
		s.comp.pc.slotC3 = true
		s.comp.pc.slotCX = false
		s.Equal(s.comp.ROM.Get(prc301), PCRead(s.comp, uc301))

		s.comp.pc.slotC3 = false
		s.Equal(s.comp.ROM.Get(irc301), PCRead(s.comp, uc301))

		// slotCX is a trick; it enables us to read C3 ROM even if C3 ROM is
		// off.
		s.comp.pc.slotCX = true
		s.Equal(s.comp.ROM.Get(prc301), PCRead(s.comp, uc301))
	})

	s.Run("reads from cx rom space", func() {
		s.comp.pc.slotCX = true
		s.Equal(s.comp.ROM.Get(prc401), PCRead(s.comp, uc401))

		s.comp.pc.slotCX = false
		s.Equal(s.comp.ROM.Get(irc401), PCRead(s.comp, uc401))
	})
}
