package a2

import "github.com/pevans/erc/pkg/data"

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
		ps.SwitchWrite(s.comp, data.DByte(0xC00B), 0x0)
		s.Equal(true, ps.slotC3)

		ps.SwitchWrite(s.comp, data.DByte(0xC00A), 0x0)
		s.Equal(false, ps.slotC3)
	})

	s.Run("slot cx rom writes work", func() {
		ps.slotCX = false
		ps.slotC3 = false
		ps.SwitchWrite(s.comp, data.DByte(0xC006), 0x0)
		s.Equal(true, ps.slotCX)
		s.Equal(true, ps.slotC3)

		ps.SwitchWrite(s.comp, data.DByte(0xC007), 0x0)
		s.Equal(false, ps.slotCX)
		s.Equal(false, ps.slotC3)
	})
}

func (s *a2Suite) TestPCSwitcherSwitchRead() {
	var (
		ps pcSwitcher
		hi data.Byte = 0x80
		lo data.Byte = 0x00
	)

	s.Run("read of slotc3 returns hi", func() {
		ps.slotC3 = true
		s.Equal(hi, ps.SwitchRead(s.comp, data.DByte(0xC017)))
	})

	s.Run("read of slot cx returns lo", func() {
		ps.slotCX = true
		s.Equal(lo, ps.SwitchRead(s.comp, data.DByte(0xC016)))
	})
}

func (s *a2Suite) TestPCRead() {
	var (
		c801   = data.DByte(0xC801)
		prc801 = data.DByte(pcPROMAddr(c801.Addr()))
		irc801 = data.DByte(pcIROMAddr(c801.Addr()))
		c301   = data.DByte(0xC301)
		prc301 = data.DByte(pcPROMAddr(c301.Addr()))
		irc301 = data.DByte(pcIROMAddr(c301.Addr()))
		c401   = data.DByte(0xC401)
		prc401 = data.DByte(pcPROMAddr(c401.Addr()))
		irc401 = data.DByte(pcIROMAddr(c401.Addr()))
	)

	s.Run("reads from expansion space", func() {
		s.comp.pc.expansion = true
		s.Equal(s.comp.ROM.Get(prc801), PCRead(s.comp, c801))

		s.comp.pc.expansion = false
		s.Equal(s.comp.ROM.Get(irc801), PCRead(s.comp, c801))
	})

	s.Run("reads from c3 rom space", func() {
		s.comp.pc.slotC3 = true
		s.comp.pc.slotCX = false
		s.Equal(s.comp.ROM.Get(prc301), PCRead(s.comp, c301))

		s.comp.pc.slotC3 = false
		s.Equal(s.comp.ROM.Get(irc301), PCRead(s.comp, c301))

		// slotCX is a trick; it enables us to read C3 ROM even if C3 ROM is
		// off.
		s.comp.pc.slotCX = true
		s.Equal(s.comp.ROM.Get(prc301), PCRead(s.comp, c301))
	})

	s.Run("reads from cx rom space", func() {
		s.comp.pc.slotCX = true
		s.Equal(s.comp.ROM.Get(prc401), PCRead(s.comp, c401))

		s.comp.pc.slotCX = false
		s.Equal(s.comp.ROM.Get(irc401), PCRead(s.comp, c401))
	})
}
