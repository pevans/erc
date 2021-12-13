package a2

import "github.com/pevans/erc/pkg/data"

func (s *a2Suite) TestUseDefaults() {
	s.comp.bank.UseDefaults(s.comp)
	s.Equal(bankROM, s.comp.state.Int(bankRead))
	s.Equal(bankRAM, s.comp.state.Int(bankWrite))
	s.Equal(bank2, s.comp.state.Int(bankDFBlock))
	s.Equal(bankMain, s.comp.state.Int(bankSysBlock))
}

func (s *a2Suite) TestSwitchRead() {
	bank := bankSwitcher{}

	rmodes := map[int]int{
		0xC080: bankRAM,
		0xC081: bankROM,
		0xC082: bankROM,
		0xC083: bankRAM,
		0xC088: bankRAM,
		0xC089: bankROM,
		0xC08A: bankROM,
		0xC08B: bankRAM,
	}

	wmodes := map[int]int{
		0xC080: bankNone,
		0xC081: bankRAM,
		0xC082: bankNone,
		0xC083: bankRAM,
		0xC088: bankNone,
		0xC089: bankRAM,
		0xC08A: bankNone,
		0xC08B: bankRAM,
	}

	dfmodes := map[int]int{
		0xC080: bank2,
		0xC081: bank2,
		0xC082: bank2,
		0xC083: bank2,
		0xC088: bank1,
		0xC089: bank1,
		0xC08A: bank1,
		0xC08B: bank1,
	}

	rd := func(addr int) int {
		_ = bank.SwitchRead(s.comp, int(addr))
		return s.comp.state.Int(bankRead)
	}

	wr := func(addr int) int {
		_ = bank.SwitchRead(s.comp, int(addr))
		return s.comp.state.Int(bankWrite)
	}

	df := func(addr int) int {
		_ = bank.SwitchRead(s.comp, int(addr))
		return s.comp.state.Int(bankDFBlock)
	}

	s.Run("read modes are set properly", func() {
		for addr, mode := range rmodes {
			s.Equal(mode, rd(addr))
		}
	})

	s.Run("write modes are set properly", func() {
		for addr, mode := range wmodes {
			// Setting the write mode always requires _two_ writes, so we want
			// to confirm that the first attempt always keeps write off
			s.Equal(bankNone, wr(addr))
			s.Equal(mode, wr(addr))
		}
	})

	s.Run("df block modes are set properly", func() {
		for addr, mode := range dfmodes {
			s.Equal(mode, df(addr))
		}
	})

	s.Run("bit 7 is high", func() {
		hi7 := uint8(0x80)
		lo7 := uint8(0x00)

		s.comp.state.SetInt(bankDFBlock, bank2)
		s.Equal(hi7, bank.SwitchRead(s.comp, int(0xC011)))
		s.comp.state.SetInt(bankDFBlock, bank1)
		s.Equal(lo7, bank.SwitchRead(s.comp, int(0xC011)))

		s.comp.state.SetInt(bankRead, bankRAM)
		s.Equal(hi7, bank.SwitchRead(s.comp, int(0xC012)))
		s.comp.state.SetInt(bankRead, bankROM)
		s.Equal(lo7, bank.SwitchRead(s.comp, int(0xC012)))

		s.comp.state.SetInt(bankSysBlock, bankAux)
		s.Equal(hi7, bank.SwitchRead(s.comp, int(0xC016)))
		s.comp.state.SetInt(bankSysBlock, bankMain)
		s.Equal(lo7, bank.SwitchRead(s.comp, int(0xC016)))
	})
}

func (s *a2Suite) TestSwitchWrite() {
	var (
		bank bankSwitcher
		d123 uint8 = 123
		d45  uint8 = 45
		addr int   = 0x11
	)

	s.Run("switching main to aux", func() {
		s.comp.Main.Mem[addr] = d123
		s.comp.state.SetInt(bankSysBlock, bankMain)
		bank.SwitchWrite(s.comp, int(0xC009), d45)
		s.Equal(bankAux, s.comp.state.Int(bankSysBlock))
		s.Equal(d123, s.comp.Aux.Mem[addr])
	})

	s.Run("switching aux to main", func() {
		s.comp.Aux.Mem[addr] = d45
		bank.SwitchWrite(s.comp, int(0xC008), d123)
		s.Equal(bankMain, s.comp.state.Int(bankSysBlock))
		s.Equal(d45, s.comp.Main.Mem[addr])
	})

	s.Run("not changing the mode should not copy pages", func() {
		s.comp.Aux.Mem[addr] = d123
		bank.SwitchWrite(s.comp, int(0xC008), d123)
		s.Equal(bankMain, s.comp.state.Int(bankSysBlock))
		s.Equal(d45, s.comp.Main.Mem[addr])
	})
}

func (s *a2Suite) TestBankDFRead() {
	var (
		xd000        = 0xD000
		xe000        = 0xE000
		x1000        = 0x1000
		x2000        = 0x2000
		x10000       = 0x10000
		val1   uint8 = 124
		val2   uint8 = 112
	)

	testFor := func(sblock int) {
		s.comp.state.SetInt(bankSysBlock, sblock)

		BankSegment(s.comp.state).Set(xd000, val1)
		BankSegment(s.comp.state).Set(xe000, val1)
		BankSegment(s.comp.state).Set(x10000, val2)

		s.Run("read from rom", func() {
			s.comp.state.SetInt(bankRead, bankROM)
			s.comp.state.SetInt(bankDFBlock, bank1)
			s.Equal(s.comp.Get(xd000), s.comp.ROM.Get(x1000))
			s.Equal(s.comp.Get(xe000), s.comp.ROM.Get(x2000))

			s.comp.state.SetInt(bankDFBlock, bank2)
			s.NotEqual(s.comp.Get(xd000), s.comp.ROM.Get(x1000))
			s.Equal(s.comp.Get(xe000), s.comp.ROM.Get(x2000))
		})

		s.Run("read from bank2 ram", func() {
			s.comp.state.SetInt(bankRead, bankRAM)
			s.comp.state.SetInt(bankDFBlock, bank2)
			// The first read should use bank 2, but the second read should not,
			// since it's in the E0 page.
			s.Equal(s.comp.Get(xd000), BankSegment(s.comp.state).Get(x10000))
			s.Equal(s.comp.Get(xe000), BankSegment(s.comp.state).Get(xe000))
		})

		s.Run("read from normal (bank1) ram", func() {
			s.comp.state.SetInt(bankDFBlock, bank1)
			s.Equal(s.comp.Get(xd000), BankSegment(s.comp.state).Get(xd000))
		})
	}

	testFor(bankMain)
	testFor(bankAux)
}

func (s *a2Suite) TestBankDFWrite() {
	var (
		dfaddr       = 0xD011
		efaddr       = 0xE011
		val1   uint8 = 87
		val2   uint8 = 89
	)

	testFor := func(sblock int) {
		s.comp.state.SetInt(bankSysBlock, sblock)
		s.Run("writes respect the value of the write mode", func() {
			s.comp.state.SetInt(bankRead, bankRAM)
			s.comp.state.SetInt(bankWrite, bankRAM)
			s.comp.state.SetInt(bankDFBlock, bank1)
			s.comp.Set(dfaddr, val1)
			s.Equal(val1, s.comp.Get(dfaddr))

			s.comp.state.SetInt(bankWrite, bankNone)
			s.comp.Set(efaddr, val2)
			s.NotEqual(val2, s.comp.Get(efaddr))
		})

		s.Run("writes use bank2 in the D0-DF page range", func() {
			s.comp.state.SetInt(bankWrite, bankRAM)
			s.comp.state.SetInt(bankDFBlock, bank2)
			s.comp.Set(dfaddr, val2)
			s.Equal(val2, ReadSegment(s.comp.state).Get(0x10011))

			s.comp.Set(efaddr, val1)
			s.Equal(val1, ReadSegment(s.comp.state).Get(efaddr))
		})
	}

	testFor(bankMain)
	testFor(bankAux)
}

func (s *a2Suite) TestBankZPRead() {
	addr := 0x123
	cases := []struct {
		mode int
		seg  *data.Segment
		main uint8
		aux  uint8
		want uint8
	}{
		{bankAux, s.comp.Aux, 0x1, 0x2, 0x2},
		{bankMain, s.comp.Main, 0x3, 0x2, 0x3},
	}

	for _, c := range cases {
		s.comp.Main.DirectSet(addr, c.main)
		s.comp.Aux.DirectSet(addr, c.aux)
		s.comp.state.SetInt(bankSysBlock, c.mode)
		s.comp.state.SetSegment(bankSysBlockSegment, c.seg)

		s.Equal(c.want, s.comp.Get(addr))
	}
}

func (s *a2Suite) TestBankZPWrite() {
	addr := 0x123
	cases := []struct {
		mode int
		seg  *data.Segment
		main uint8
		aux  uint8
		want uint8
	}{
		{bankAux, s.comp.Aux, 0x0, 0x2, 0x2},
		{bankMain, s.comp.Main, 0x3, 0x0, 0x3},
	}

	for _, c := range cases {
		s.comp.Main.Set(addr, 0x0)
		s.comp.Aux.Set(addr, 0x0)

		s.comp.state.SetInt(bankSysBlock, c.mode)
		s.comp.state.SetSegment(bankSysBlockSegment, c.seg)
		s.comp.Set(addr, c.want)

		s.Equal(c.main, s.comp.Main.DirectGet(addr))
		s.Equal(c.aux, s.comp.Aux.DirectGet(addr))
	}
}
