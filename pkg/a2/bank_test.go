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
		_ = bank.SwitchRead(s.comp, data.Int(addr))
		return bank.read
	}

	wr := func(addr int) int {
		_ = bank.SwitchRead(s.comp, data.Int(addr))
		return bank.write
	}

	df := func(addr int) int {
		_ = bank.SwitchRead(s.comp, data.Int(addr))
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
		s.Equal(hi7, bank.SwitchRead(s.comp, data.Int(0xC011)))
		bank.dfBlock = bank1
		s.Equal(lo7, bank.SwitchRead(s.comp, data.Int(0xC011)))

		bank.read = bankRAM
		s.Equal(hi7, bank.SwitchRead(s.comp, data.Int(0xC012)))
		bank.read = bankROM
		s.Equal(lo7, bank.SwitchRead(s.comp, data.Int(0xC012)))

		bank.sysBlock = bankAux
		s.Equal(hi7, bank.SwitchRead(s.comp, data.Int(0xC016)))
		bank.sysBlock = bankMain
		s.Equal(lo7, bank.SwitchRead(s.comp, data.Int(0xC016)))
	})
}

func (s *a2Suite) TestSwitchWrite() {
	var (
		bank bankSwitcher
		d123 data.Byte = 123
		d45  data.Byte = 45
		addr data.Int  = 0x11
	)

	s.Run("switching main to aux", func() {
		s.comp.Main.Mem[addr] = d123
		bank.sysBlock = bankMain
		bank.SwitchWrite(s.comp, data.Int(0xC009), d45)
		s.Equal(bankAux, bank.sysBlock)
		s.Equal(d123, s.comp.Aux.Mem[addr])
	})

	s.Run("switching aux to main", func() {
		s.comp.Aux.Mem[addr] = d45
		bank.SwitchWrite(s.comp, data.Int(0xC008), d123)
		s.Equal(bankMain, bank.sysBlock)
		s.Equal(d45, s.comp.Main.Mem[addr])
	})

	s.Run("not changing the mode should not copy pages", func() {
		s.comp.Aux.Mem[addr] = d123
		bank.SwitchWrite(s.comp, data.Int(0xC008), d123)
		s.Equal(bankMain, bank.sysBlock)
		s.Equal(d45, s.comp.Main.Mem[addr])
	})
}

func (s *a2Suite) TestBankDFRead() {
	var (
		xd000  data.Int  = 0xD000
		xe000  data.Int  = 0xE000
		x1000  data.Int  = 0x1000
		x2000  data.Int  = 0x2000
		x10000 data.Int  = 0x10000
		d123   data.Byte = 123
		d111   data.Byte = 111
	)

	s.comp.WriteSegment().Set(xd000, d123)
	s.comp.WriteSegment().Set(xe000, d123)
	s.comp.WriteSegment().Set(x10000, d111)

	s.Run("read from rom", func() {
		s.comp.bank.read = bankROM
		s.Equal(s.comp.Get(xd000), s.comp.ROM.Get(x1000))
		s.Equal(s.comp.Get(xe000), s.comp.ROM.Get(x2000))
	})

	s.Run("read from bank2 ram", func() {
		s.comp.bank.read = bankRAM
		s.comp.bank.dfBlock = bank2
		// The first read should use bank 2, but the second read should not,
		// since it's in the E0 page.
		s.Equal(s.comp.Get(xd000), s.comp.ReadSegment().Get(x10000))
		s.Equal(s.comp.Get(xe000), s.comp.ReadSegment().Get(xe000))
	})

	s.Run("read from normal (bank1) ram", func() {
		s.comp.bank.dfBlock = bank1
		s.Equal(s.comp.Get(xd000), s.comp.ReadSegment().Get(xd000))
	})
}

func (s *a2Suite) TestBankDFWrite() {
	var (
		dfaddr data.Int  = 0xD011
		efaddr data.Int  = 0xE011
		val1   data.Byte = 123
		val2   data.Byte = 111
	)

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
		s.Equal(val2, s.comp.ReadSegment().Get(data.Int(0x10000)))

		s.comp.Set(efaddr, val1)
		s.Equal(val1, s.comp.ReadSegment().Get(efaddr))
	})
}
