package a2

import (
	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/memory"
)

func (s *a2Suite) TestUseDefaults() {
	bankUseDefaults(s.comp)
	s.False(s.comp.State.Bool(a2state.BankReadRAM))
	s.True(s.comp.State.Bool(a2state.BankWriteRAM))
	s.False(s.comp.State.Bool(a2state.BankDFBlockBank2))
	s.False(s.comp.State.Bool(a2state.BankSysBlockAux))
}

func (s *a2Suite) TestSwitchRead() {
	rmodes := map[int]bool{
		0xC080: true,
		0xC081: false,
		0xC082: false,
		0xC083: true,
		0xC088: true,
		0xC089: false,
		0xC08A: false,
		0xC08B: true,
	}

	wmodes := map[int]bool{
		0xC080: false,
		0xC081: true,
		0xC082: false,
		0xC083: true,
		0xC088: false,
		0xC089: true,
		0xC08A: false,
		0xC08B: true,
	}

	dfmodes := map[int]bool{
		0xC080: true,
		0xC081: true,
		0xC082: true,
		0xC083: true,
		0xC088: false,
		0xC089: false,
		0xC08A: false,
		0xC08B: false,
	}

	rd := func(addr int) bool {
		_ = bankSwitchRead(int(addr), s.comp.State)
		return s.comp.State.Bool(a2state.BankReadRAM)
	}

	wr := func(addr int) bool {
		// Because the read attempts are adjusted in the computer
		// Process method, we simulate that here.
		_ = bankSwitchRead(int(addr), s.comp.State)

		switch addr {
		case
			0xC081, 0xC083, 0xC085, 0xC087,
			0xC089, 0xC08B, 0xC08D, 0xC08F:
			s.comp.State.SetInt(
				a2state.BankReadAttempts,
				s.comp.State.Int(a2state.BankReadAttempts)+1,
			)
			s.comp.State.SetBool(a2state.InstructionReadOp, true)
		default:
			s.comp.State.SetBool(a2state.InstructionReadOp, false)
			s.comp.State.SetInt(a2state.BankReadAttempts, 0)
		}

		return s.comp.State.Bool(a2state.BankWriteRAM)
	}

	df := func(addr int) bool {
		_ = bankSwitchRead(int(addr), s.comp.State)
		return s.comp.State.Bool(a2state.BankDFBlockBank2)
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
			s.False(wr(addr))
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

		s.comp.State.SetBool(a2state.BankDFBlockBank2, true)
		s.Equal(hi7, bankSwitchRead(int(0xC011), s.comp.State))
		s.comp.State.SetBool(a2state.BankDFBlockBank2, false)
		s.Equal(lo7, bankSwitchRead(int(0xC011), s.comp.State))

		s.comp.State.SetBool(a2state.BankReadRAM, true)
		s.Equal(hi7, bankSwitchRead(int(0xC012), s.comp.State))
		s.comp.State.SetBool(a2state.BankReadRAM, false)
		s.Equal(lo7, bankSwitchRead(int(0xC012), s.comp.State))

		s.comp.State.SetBool(a2state.BankSysBlockAux, true)
		s.Equal(hi7, bankSwitchRead(int(0xC016), s.comp.State))
		s.comp.State.SetBool(a2state.BankSysBlockAux, false)
		s.Equal(lo7, bankSwitchRead(int(0xC016), s.comp.State))
	})
}

func (s *a2Suite) TestSwitchWrite() {
	var (
		d123 uint8 = 123
		d45  uint8 = 45
		addr int   = 0x11
	)

	s.Run("switching main to aux", func() {
		s.comp.Main.Mem[addr] = d123
		s.comp.Aux.Mem[addr] = d45
		s.comp.State.SetBool(a2state.BankSysBlockAux, false)
		bankSwitchWrite(int(0xC009), d45, s.comp.State)
		s.True(s.comp.State.Bool(a2state.BankSysBlockAux))
		s.Equal(d45, s.comp.Aux.Mem[addr])
	})

	s.Run("switching aux to main", func() {
		s.comp.Aux.Mem[addr] = d45
		bankSwitchWrite(int(0xC008), d123, s.comp.State)
		s.False(s.comp.State.Bool(a2state.BankSysBlockAux))
		s.Equal(d123, s.comp.Main.Mem[addr])
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

	testForBankAux := func(useAux bool) {
		s.comp.State.SetBool(a2state.BankSysBlockAux, useAux)

		BankSegment(s.comp.State).Set(xd000, val1)
		BankSegment(s.comp.State).Set(xe000, val1)
		BankSegment(s.comp.State).Set(x10000, val2)

		s.Run("read from rom", func() {
			s.comp.State.SetBool(a2state.BankReadRAM, false)
			s.comp.State.SetBool(a2state.BankDFBlockBank2, false)
			s.Equal(s.comp.Get(xd000), s.comp.ROM.DirectGet(x1000))
			s.Equal(s.comp.Get(xe000), s.comp.ROM.DirectGet(x2000))

			s.comp.State.SetBool(a2state.BankDFBlockBank2, true)
			s.Equal(s.comp.Get(xd000), s.comp.ROM.DirectGet(x1000))
			s.Equal(s.comp.Get(xe000), s.comp.ROM.DirectGet(x2000))
		})

		s.Run("read from bank2 ram", func() {
			s.comp.State.SetBool(a2state.BankReadRAM, true)
			s.comp.State.SetBool(a2state.BankDFBlockBank2, true)
			// The first read should use bank 2, but the second read should not,
			// since it's in the E0 page.
			s.Equal(s.comp.Get(xd000), BankSegment(s.comp.State).Get(x10000))
			s.Equal(s.comp.Get(xe000), BankSegment(s.comp.State).Get(xe000))
		})

		s.Run("read from normal (bank1) ram", func() {
			s.comp.State.SetBool(a2state.BankDFBlockBank2, false)
			s.Equal(s.comp.Get(xd000), BankSegment(s.comp.State).Get(xd000))
		})
	}

	testForBankAux(false)
	testForBankAux(true)
}

func (s *a2Suite) TestBankDFWrite() {
	var (
		dfaddr       = 0xD011
		efaddr       = 0xE011
		val1   uint8 = 87
		val2   uint8 = 89
	)

	testForBankAux := func(useAux bool) {
		s.comp.State.SetBool(a2state.BankSysBlockAux, useAux)
		s.Run("writes respect the value of the write mode", func() {
			s.comp.State.SetBool(a2state.BankReadRAM, true)
			s.comp.State.SetBool(a2state.BankWriteRAM, true)
			s.comp.State.SetBool(a2state.BankDFBlockBank2, false)
			s.comp.Set(dfaddr, val1)
			s.Equal(val1, s.comp.Get(dfaddr))

			s.comp.State.SetBool(a2state.BankWriteRAM, false)
			s.comp.Set(efaddr, val2)
			s.NotEqual(val2, s.comp.Get(efaddr))
		})

		s.Run("writes use bank2 in the D0-DF page range", func() {
			s.comp.State.SetBool(a2state.BankWriteRAM, true)
			s.comp.State.SetBool(a2state.BankDFBlockBank2, true)
			s.comp.Set(dfaddr, val2)
			s.Equal(val2, ReadSegment(s.comp.State).Get(0x10011))

			s.comp.Set(efaddr, val1)
			s.Equal(val1, ReadSegment(s.comp.State).Get(efaddr))
		})
	}

	testForBankAux(false)
	testForBankAux(true)
}

func (s *a2Suite) TestBankZPRead() {
	addr := 0x123
	cases := []struct {
		useAux bool
		seg    *memory.Segment
		main   uint8
		aux    uint8
		want   uint8
	}{
		{true, s.comp.Aux, 0x1, 0x2, 0x2},
		{false, s.comp.Main, 0x3, 0x2, 0x3},
	}

	for _, c := range cases {
		s.comp.Main.DirectSet(addr, c.main)
		s.comp.Aux.DirectSet(addr, c.aux)
		s.comp.State.SetBool(a2state.BankSysBlockAux, c.useAux)
		s.comp.State.SetSegment(a2state.BankSysBlockSegment, c.seg)

		s.Equal(c.want, s.comp.Get(addr))
	}
}

func (s *a2Suite) TestBankZPWrite() {
	addr := 0x123
	cases := []struct {
		useAux bool
		seg    *memory.Segment
		main   uint8
		aux    uint8
		want   uint8
	}{
		{true, s.comp.Aux, 0x0, 0x2, 0x2},
		{false, s.comp.Main, 0x3, 0x0, 0x3},
	}

	for _, c := range cases {
		s.comp.Main.Set(addr, 0x0)
		s.comp.Aux.Set(addr, 0x0)

		s.comp.State.SetBool(a2state.BankSysBlockAux, c.useAux)
		s.comp.State.SetSegment(a2state.BankSysBlockSegment, c.seg)
		s.comp.Set(addr, c.want)

		s.Equal(c.main, s.comp.Main.DirectGet(addr))
		s.Equal(c.aux, s.comp.Aux.DirectGet(addr))
	}
}
