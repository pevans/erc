package a2

import (
	"github.com/pevans/erc/pkg/data"
	"github.com/pevans/erc/pkg/mos65c02"
)

func (s *a2Suite) TestBoot() {
	c := NewComputer()

	s.NoError(c.Boot(""))

	// We know as part of the boot procedure that we copy in rom, but we
	// don't necessarily want to test the entirety of that; let's just
	// test ROM doesn't look empty.
	s.NotEqual(data.Byte(0), c.ROM.Get(0x100))

	s.Equal(data.Byte(AppleSoft&0xFF), c.Main.Get(BootVector))
	s.Equal(data.Byte(AppleSoft>>8), c.Main.Get(BootVector+1))
}

func (s *a2Suite) TestReset() {
	c := NewComputer()

	c.Reset()

	s.Equal(
		mos65c02.INTERRUPT|mos65c02.BREAK|mos65c02.UNUSED,
		c.CPU.P,
	)
	s.Equal(c.CPU.Get16(ResetPC), c.CPU.PC)
	s.Equal(data.Byte(0xFF), c.CPU.S)
	s.True(c.disp.text)
}
