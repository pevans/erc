package a2

import (
	"github.com/pevans/erc/pkg/data"
)

func (s *a2Suite) TestZeroPageRead() {
	addr := data.DByte(0x123)
	cases := []struct {
		mode int
		main data.Byte
		aux  data.Byte
		want data.Byte
	}{
		{BankAuxiliary, 0x1, 0x2, 0x2},
		{0, 0x3, 0x2, 0x3},
	}

	for _, c := range cases {
		s.comp.Main.Set(addr, c.main)
		s.comp.Aux.Set(addr, c.aux)
		s.comp.BankMode = c.mode

		s.Equal(c.want, s.comp.Get(addr))
	}
}

func (s *a2Suite) TestZeroPageWrite() {
	addr := data.DByte(0x123)
	cases := []struct {
		mode int
		main data.Byte
		aux  data.Byte
		want data.Byte
	}{
		{BankAuxiliary, 0x0, 0x2, 0x2},
		{0, 0x3, 0x0, 0x3},
	}

	for _, c := range cases {
		s.comp.Main.Set(addr, 0x0)
		s.comp.Aux.Set(addr, 0x0)

		s.comp.BankMode = c.mode
		s.comp.Set(addr, c.want)

		s.Equal(c.main, s.comp.Main.Get(addr))
		s.Equal(c.aux, s.comp.Aux.Get(addr))
	}
}
