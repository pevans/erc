package a2

import (
	"github.com/pevans/erc/pkg/mach"
	"github.com/stretchr/testify/assert"
)

func (s *a2Suite) TestDefineSoftSwitches() {
	var ok bool

	for addr := 0x0; addr < 0x200; addr++ {
		_, ok = s.comp.RMap[addr]
		assert.Equal(s.T(), true, ok)

		_, ok = s.comp.WMap[addr]
		assert.Equal(s.T(), true, ok)
	}
}

func (s *a2Suite) TestZeroPageRead() {
	addr := mach.DByte(0x123)
	cases := []struct {
		mode int
		main mach.Byte
		aux  mach.Byte
		want mach.Byte
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
	addr := mach.DByte(0x123)
	cases := []struct {
		mode int
		main mach.Byte
		aux  mach.Byte
		want mach.Byte
	}{
		{BankAuxiliary, 0x0, 0x2, 0x2},
		{0, 0x3, 0x0, 0x3},
	}

	for _, c := range cases {
		s.comp.BankMode = c.mode
		s.comp.Set(addr, c.want)

		assert.Equal(s.T(), c.main, s.comp.Main.Get(addr))
		assert.Equal(s.T(), c.aux, s.comp.Aux.Get(addr))
	}
}
