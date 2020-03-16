package a2

import (
	"github.com/pevans/erc/pkg/data"
	"github.com/stretchr/testify/assert"
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
		assert.Equal(s.T(), c.want, bankMode(s.comp))
	}
}

func (s *a2Suite) TestBankSetMode() {
	cases := []struct {
		bankMode int
		newMode  int
		want     int
	}{
		{4, 3, 3},
		{0, 2, 2},
		{1, 0, 0},
	}

	for _, c := range cases {
		s.comp.BankMode = c.bankMode
		bankSetMode(s.comp, c.newMode)
		assert.Equal(s.T(), c.want, s.comp.BankMode)
	}
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

		assert.Equal(s.T(), c.want, bankRead(s.comp, c.addr))
	}
}
