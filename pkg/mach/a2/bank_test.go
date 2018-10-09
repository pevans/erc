package a2

import (
	"github.com/pevans/erc/pkg/mach"
	"github.com/stretchr/testify/assert"
)

func (s *a2Suite) TestBankRead() {
	cases := []struct {
		bankMode int
		addr     mach.DByte
		romSet   mach.Byte
		ram1Set  mach.Byte
		ram2Set  mach.Byte
		want     mach.Byte
	}{
		{0, 0xD012, 0x11, 0x22, 0x33, 0x11},
		{BankRAM, 0xD012, 0x11, 0x22, 0x33, 0x22},
		{BankRAM | BankRAM2, 0xD012, 0x11, 0x22, 0x33, 0x33},
	}

	for _, c := range cases {
		s.comp.BankMode = c.bankMode
		s.comp.ROM.Set(mach.Plus(c.addr, -SysRomOffset), c.romSet)
		s.comp.ReadSegment().Set(c.addr, c.ram1Set)
		s.comp.ReadSegment().Set(mach.Plus(c.addr, 0x3000), c.ram2Set)

		assert.Equal(s.T(), c.want, bankRead(s.comp, c.addr))
	}
}
