package a2

import (
	"github.com/pevans/erc/pkg/data"
)

func (s *a2Suite) TestBankMode() {
	cases := []struct {
		bankMode int
		want     int
	}{
		{4, 4},
		{0, 0},
	}

	for _, c := range cases {
		s.comp.BankMode = c.bankMode
		s.Equal(c.want, bankMode(s.comp))
	}
}

func (s *a2Suite) TestBankSetMode() {
	// Test that just any copy works
	s.comp.BankMode = BankWrite
	bankSetMode(s.comp, BankRAM)
	s.Equal(BankRAM, s.comp.BankMode)

	val := data.Byte(0x12)
	idx := data.DByte(0x1)

	// Test that a bank set for aux memory works
	s.comp.Main.Mem[idx] = val
	bankSetMode(s.comp, BankAuxiliary)
	s.Equal(val, s.comp.Aux.Mem[idx])

	// Test that a bank set for main memory from aux works
	s.comp.Aux.Mem[idx] = val + 1
	bankSetMode(s.comp, BankDefault)
	s.Equal(val+1, s.comp.Main.Mem[idx])
}

func (s *a2Suite) TestBankRead() {
	cases := []struct {
		bankMode int
		addr     data.DByte
		romSet   data.Byte
		ram1Set  data.Byte
		ram2Set  data.Byte
		want     data.Byte
	}{
		{0, 0xD012, 0x11, 0x22, 0x33, 0x11},
		{BankRAM, 0xD012, 0x11, 0x22, 0x33, 0x22},
		{BankRAM | BankRAM2, 0xD012, 0x11, 0x22, 0x33, 0x33},
	}

	for _, c := range cases {
		s.comp.BankMode = c.bankMode
		s.comp.ROM.Set(data.Plus(c.addr, -SysRomOffset), c.romSet)
		s.comp.ReadSegment().Set(c.addr, c.ram1Set)
		s.comp.ReadSegment().Set(data.Plus(c.addr, 0x3000), c.ram2Set)

		s.Equal(c.want, bankRead(s.comp, c.addr))
	}
}

func (s *a2Suite) TestBankWrite() {
	idx := data.DByte(0x1)
	val := data.Byte(0x12)

	// Test that writes without write mode will fail
	s.comp.BankMode = BankDefault
	bankWrite(s.comp, idx, val)
	s.NotEqual(val, s.comp.Main.Mem[idx])

	// Test that nominal writes will succeed
	s.comp.BankMode = BankWrite
	bankWrite(s.comp, idx, val)
	s.Equal(val, s.comp.Main.Mem[idx])

	// Test that writes in bank-switchable memory for RAM2 go to the
	// write place.
	s.comp.BankMode |= BankRAM2
	bankWrite(s.comp, idx, val)
	s.Equal(val, s.comp.ReadSegment().Get(idx))
}
