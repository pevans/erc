package a2

import (
	"github.com/pevans/erc/pkg/data"
	"github.com/stretchr/testify/assert"
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

		assert.Equal(s.T(), c.want, s.comp.Get(addr))
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

		assert.Equal(s.T(), c.main, s.comp.Main.Get(addr))
		assert.Equal(s.T(), c.aux, s.comp.Aux.Get(addr))
	}
}
