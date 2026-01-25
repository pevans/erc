package a2bank

import (
	"testing"

	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/memory"
	"github.com/stretchr/testify/suite"
)

type bankSuite struct {
	suite.Suite

	state *memory.StateMap
	main  *memory.Segment
	aux   *memory.Segment
	rom   *memory.Segment
}

func (s *bankSuite) SetupTest() {
	// Create memory segments (need extra space for bank2 addressing: addr +
	// 0x3000)
	s.main = memory.NewSegment(0x14000)
	s.aux = memory.NewSegment(0x14000)
	s.rom = memory.NewSegment(0x5000)

	// Create state map
	s.state = memory.NewStateMap()

	// Set up segment mappings
	s.state.SetSegment(a2state.MemMainSegment, s.main)
	s.state.SetSegment(a2state.MemAuxSegment, s.aux)
	s.state.SetSegment(a2state.MemReadSegment, s.main)
	s.state.SetSegment(a2state.MemWriteSegment, s.main)
	s.state.SetSegment(a2state.BankROMSegment, s.rom)
	s.state.SetSegment(a2state.BankSysBlockSegment, s.main)
}

func TestBankSuite(t *testing.T) {
	suite.Run(t, new(bankSuite))
}

func (s *bankSuite) TestUseDefaults() {
	UseDefaults(s.state, s.main, s.rom)
	s.False(s.state.Bool(a2state.BankReadRAM))
	s.True(s.state.Bool(a2state.BankWriteRAM))
	s.False(s.state.Bool(a2state.BankDFBlockBank2))
	s.False(s.state.Bool(a2state.BankSysBlockAux))
}

func (s *bankSuite) TestSwitchRead() {
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
		_ = SwitchRead(int(addr), s.state)
		return s.state.Bool(a2state.BankReadRAM)
	}

	wr := func(addr int) bool {
		// Because the read attempts are adjusted in the computer Process
		// method, we simulate that here.
		_ = SwitchRead(int(addr), s.state)

		switch addr {
		case
			0xC081, 0xC083, 0xC085, 0xC087,
			0xC089, 0xC08B, 0xC08D, 0xC08F:
			s.state.SetInt(
				a2state.BankReadAttempts,
				s.state.Int(a2state.BankReadAttempts)+1,
			)
			s.state.SetBool(a2state.InstructionReadOp, true)
		default:
			s.state.SetBool(a2state.InstructionReadOp, false)
			s.state.SetInt(a2state.BankReadAttempts, 0)
		}

		return s.state.Bool(a2state.BankWriteRAM)
	}

	df := func(addr int) bool {
		_ = SwitchRead(int(addr), s.state)
		return s.state.Bool(a2state.BankDFBlockBank2)
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

		s.state.SetBool(a2state.BankDFBlockBank2, true)
		s.Equal(hi7, SwitchRead(int(0xC011), s.state))
		s.state.SetBool(a2state.BankDFBlockBank2, false)
		s.Equal(lo7, SwitchRead(int(0xC011), s.state))

		s.state.SetBool(a2state.BankReadRAM, true)
		s.Equal(hi7, SwitchRead(int(0xC012), s.state))
		s.state.SetBool(a2state.BankReadRAM, false)
		s.Equal(lo7, SwitchRead(int(0xC012), s.state))

		s.state.SetBool(a2state.BankSysBlockAux, true)
		s.Equal(hi7, SwitchRead(int(0xC016), s.state))
		s.state.SetBool(a2state.BankSysBlockAux, false)
		s.Equal(lo7, SwitchRead(int(0xC016), s.state))
	})
}

func (s *bankSuite) TestSwitchWrite() {
	var (
		d123 uint8 = 123
		d45  uint8 = 45
		addr       = 0x11
	)

	s.Run("switching main to aux", func() {
		s.main.DirectSet(addr, d123)
		s.aux.DirectSet(addr, d45)
		s.state.SetBool(a2state.BankSysBlockAux, false)
		SwitchWrite(int(0xC009), d45, s.state)
		s.True(s.state.Bool(a2state.BankSysBlockAux))
		s.Equal(d45, s.aux.DirectGet(addr))
	})

	s.Run("switching aux to main", func() {
		s.aux.DirectSet(addr, d45)
		SwitchWrite(int(0xC008), d123, s.state)
		s.False(s.state.Bool(a2state.BankSysBlockAux))
		s.Equal(d123, s.main.DirectGet(addr))
	})
}

func (s *bankSuite) TestDFRead() {
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
		s.state.SetBool(a2state.BankSysBlockAux, useAux)

		Segment(s.state).Set(xd000, val1)
		Segment(s.state).Set(xe000, val1)
		Segment(s.state).Set(x10000, val2)

		s.Run("read from rom", func() {
			s.state.SetBool(a2state.BankReadRAM, false)
			s.state.SetBool(a2state.BankDFBlockBank2, false)
			s.Equal(DFRead(xd000, s.state), s.rom.DirectGet(x1000))
			s.Equal(DFRead(xe000, s.state), s.rom.DirectGet(x2000))

			s.state.SetBool(a2state.BankDFBlockBank2, true)
			s.Equal(DFRead(xd000, s.state), s.rom.DirectGet(x1000))
			s.Equal(DFRead(xe000, s.state), s.rom.DirectGet(x2000))
		})

		s.Run("read from bank2 ram", func() {
			s.state.SetBool(a2state.BankReadRAM, true)
			s.state.SetBool(a2state.BankDFBlockBank2, true)
			// The first read should use bank 2, but the second read should
			// not, since it's in the E0 page.
			s.Equal(DFRead(xd000, s.state), Segment(s.state).Get(x10000))
			s.Equal(DFRead(xe000, s.state), Segment(s.state).Get(xe000))
		})

		s.Run("read from normal (bank1) ram", func() {
			s.state.SetBool(a2state.BankDFBlockBank2, false)
			s.Equal(DFRead(xd000, s.state), Segment(s.state).Get(xd000))
		})
	}

	testForBankAux(false)
	testForBankAux(true)
}

func (s *bankSuite) TestDFWrite() {
	var (
		dfaddr       = 0xD011
		efaddr       = 0xE011
		val1   uint8 = 87
		val2   uint8 = 89
	)

	testForBankAux := func(useAux bool) {
		s.state.SetBool(a2state.BankSysBlockAux, useAux)
		s.Run("writes respect the value of the write mode", func() {
			s.state.SetBool(a2state.BankReadRAM, true)
			s.state.SetBool(a2state.BankWriteRAM, true)
			s.state.SetBool(a2state.BankDFBlockBank2, false)
			DFWrite(dfaddr, val1, s.state)
			s.Equal(val1, DFRead(dfaddr, s.state))

			s.state.SetBool(a2state.BankWriteRAM, false)
			DFWrite(efaddr, val2, s.state)
			s.NotEqual(val2, DFRead(efaddr, s.state))
		})

		s.Run("writes use bank2 in the D0-DF page range", func() {
			s.state.SetBool(a2state.BankWriteRAM, true)
			s.state.SetBool(a2state.BankDFBlockBank2, true)
			DFWrite(dfaddr, val2, s.state)
			s.Equal(val2, s.state.Segment(a2state.MemReadSegment).Get(0x10011))

			DFWrite(efaddr, val1, s.state)
			s.Equal(val1, s.state.Segment(a2state.MemReadSegment).Get(efaddr))
		})
	}

	testForBankAux(false)
	testForBankAux(true)
}

func (s *bankSuite) TestZPRead() {
	addr := 0x123
	cases := []struct {
		useAux bool
		seg    *memory.Segment
		main   uint8
		aux    uint8
		want   uint8
	}{
		{true, s.aux, 0x1, 0x2, 0x2},
		{false, s.main, 0x3, 0x2, 0x3},
	}

	for _, c := range cases {
		s.main.DirectSet(addr, c.main)
		s.aux.DirectSet(addr, c.aux)
		s.state.SetBool(a2state.BankSysBlockAux, c.useAux)
		s.state.SetSegment(a2state.BankSysBlockSegment, c.seg)

		s.Equal(c.want, ZPRead(addr, s.state))
	}
}

func (s *bankSuite) TestZPWrite() {
	addr := 0x123
	cases := []struct {
		useAux bool
		seg    *memory.Segment
		main   uint8
		aux    uint8
		want   uint8
	}{
		{true, s.aux, 0x0, 0x2, 0x2},
		{false, s.main, 0x3, 0x0, 0x3},
	}

	for _, c := range cases {
		s.main.Set(addr, 0x0)
		s.aux.Set(addr, 0x0)

		s.state.SetBool(a2state.BankSysBlockAux, c.useAux)
		s.state.SetSegment(a2state.BankSysBlockSegment, c.seg)
		ZPWrite(addr, c.want, s.state)

		s.Equal(c.main, s.main.DirectGet(addr))
		s.Equal(c.aux, s.aux.DirectGet(addr))
	}
}
