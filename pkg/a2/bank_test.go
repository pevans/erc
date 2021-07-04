package a2

import "github.com/pevans/erc/pkg/data"

func (s *a2Suite) TestUseDefaults() {
	s.comp.bank.UseDefaults()
	s.Equal(bankROM, s.comp.bank.read)
	s.Equal(bankRAM, s.comp.bank.write)
	s.Equal(bank2, s.comp.bank.dfBlock)
	s.Equal(bankMain, s.comp.bank.sysBlock)
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
		_ = bank.SwitchRead(s.comp, data.DByte(addr))
		return bank.read
	}

	wr := func(addr int) int {
		_ = bank.SwitchRead(s.comp, data.DByte(addr))
		return bank.write
	}

	df := func(addr int) int {
		_ = bank.SwitchRead(s.comp, data.DByte(addr))
		return bank.dfBlock
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
		hi7 := data.Byte(0x80)
		lo7 := data.Byte(0x00)

		bank.dfBlock = bank2
		s.Equal(hi7, bank.SwitchRead(s.comp, data.DByte(0xC011)))
		bank.dfBlock = bank1
		s.Equal(lo7, bank.SwitchRead(s.comp, data.DByte(0xC011)))

		bank.read = bankRAM
		s.Equal(hi7, bank.SwitchRead(s.comp, data.DByte(0xC012)))
		bank.read = bankROM
		s.Equal(lo7, bank.SwitchRead(s.comp, data.DByte(0xC012)))

		bank.sysBlock = bankAux
		s.Equal(hi7, bank.SwitchRead(s.comp, data.DByte(0xC016)))
		bank.sysBlock = bankMain
		s.Equal(lo7, bank.SwitchRead(s.comp, data.DByte(0xC016)))
	})
}

func (s *a2Suite) TestSwitchWrite() {
	var (
		bank bankSwitcher
		d123 data.Byte  = 123
		d45  data.Byte  = 45
		addr data.DByte = 0x11
	)

	s.Run("switching main to aux", func() {
		s.comp.Main.Mem[addr] = d123
		bank.sysBlock = bankMain
		bank.SwitchWrite(s.comp, data.DByte(0xC009), d45)
		s.Equal(bankAux, bank.sysBlock)
		s.Equal(d123, s.comp.Aux.Mem[addr])
	})

	s.Run("switching aux to main", func() {
		s.comp.Aux.Mem[addr] = d45
		bank.SwitchWrite(s.comp, data.DByte(0xC008), d123)
		s.Equal(bankMain, bank.sysBlock)
		s.Equal(d45, s.comp.Main.Mem[addr])
	})

	s.Run("not changing the mode should not copy pages", func() {
		s.comp.Aux.Mem[addr] = d123
		bank.SwitchWrite(s.comp, data.DByte(0xC008), d123)
		s.Equal(bankMain, bank.sysBlock)
		s.Equal(d45, s.comp.Main.Mem[addr])
	})
}

func (s *a2Suite) TestBankDFRead() {
	var (
		xd000            = 0xD000
		xe000            = 0xE000
		x1000            = 0x1000
		x2000            = 0x2000
		x10000           = 0x10000
		val1   data.Byte = 124
		val2   data.Byte = 112
	)

	testFor := func(sblock int) {
		s.comp.bank.sysBlock = sblock

		s.comp.BankSegment().Set(xd000, val1)
		s.comp.BankSegment().Set(xe000, val1)
		s.comp.BankSegment().Set(x10000, val2)

		s.Run("read from rom", func() {
			s.comp.bank.read = bankROM
			s.comp.bank.dfBlock = bank1
			s.Equal(s.comp.Get(xd000), s.comp.ROM.Get(x1000))
			s.Equal(s.comp.Get(xe000), s.comp.ROM.Get(x2000))

			s.comp.bank.dfBlock = bank2
			s.NotEqual(s.comp.Get(xd000), s.comp.ROM.Get(x1000))
			s.Equal(s.comp.Get(xe000), s.comp.ROM.Get(x2000))
		})

		s.Run("read from bank2 ram", func() {
			s.comp.bank.read = bankRAM
			s.comp.bank.dfBlock = bank2
			// The first read should use bank 2, but the second read should not,
			// since it's in the E0 page.
			s.Equal(s.comp.Get(xd000), s.comp.BankSegment().Get(x10000))
			s.Equal(s.comp.Get(xe000), s.comp.BankSegment().Get(xe000))
		})

		s.Run("read from normal (bank1) ram", func() {
			s.comp.bank.dfBlock = bank1
			s.Equal(s.comp.Get(xd000), s.comp.BankSegment().Get(xd000))
		})
	}

	testFor(bankMain)
	testFor(bankAux)
}

func (s *a2Suite) TestBankDFWrite() {
	var (
		dfaddr           = 0xD011
		efaddr           = 0xE011
		val1   data.Byte = 87
		val2   data.Byte = 89
	)

	testFor := func(sblock int) {
		s.comp.bank.sysBlock = sblock
		s.Run("writes respect the value of the write mode", func() {
			s.comp.bank.read = bankRAM
			s.comp.bank.write = bankRAM
			s.comp.bank.dfBlock = bank1
			s.comp.Set(dfaddr, val1)
			s.Equal(val1, s.comp.Get(dfaddr))

			s.comp.bank.write = bankNone
			s.comp.Set(efaddr, val2)
			s.NotEqual(val2, s.comp.Get(efaddr))
		})

		s.Run("writes use bank2 in the D0-DF page range", func() {
			s.comp.bank.write = bankRAM
			s.comp.bank.dfBlock = bank2
			s.comp.Set(dfaddr, val2)
			s.Equal(val2, s.comp.ReadSegment().Get(0x10011))

			s.comp.Set(efaddr, val1)
			s.Equal(val1, s.comp.ReadSegment().Get(efaddr))
		})
	}

	testFor(bankMain)
	testFor(bankAux)
}

func (s *a2Suite) TestBankZPRead() {
	addr := 0x123
	cases := []struct {
		mode int
		main data.Byte
		aux  data.Byte
		want data.Byte
	}{
		{bankAux, 0x1, 0x2, 0x2},
		{bankMain, 0x3, 0x2, 0x3},
	}

	for _, c := range cases {
		s.comp.Main.Set(addr, c.main)
		s.comp.Aux.Set(addr, c.aux)
		s.comp.bank.sysBlock = c.mode

		s.Equal(c.want, s.comp.Get(addr))
	}
}

func (s *a2Suite) TestBankZPWrite() {
	addr := 0x123
	cases := []struct {
		mode int
		main data.Byte
		aux  data.Byte
		want data.Byte
	}{
		{bankAux, 0x0, 0x2, 0x2},
		{bankMain, 0x3, 0x0, 0x3},
	}

	for _, c := range cases {
		s.comp.Main.Set(addr, 0x0)
		s.comp.Aux.Set(addr, 0x0)

		s.comp.bank.sysBlock = c.mode
		s.comp.Set(addr, c.want)

		s.Equal(c.main, s.comp.Main.Get(addr))
		s.Equal(c.aux, s.comp.Aux.Get(addr))
	}
}
