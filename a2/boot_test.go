package a2

import (
	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/mos"
)

func (s *a2Suite) TestBoot() {
	c := NewComputer(123)

	s.NoError(c.Boot())

	// We know as part of the boot procedure that we copy in rom, but we
	// don't necessarily want to test the entirety of that; let's just
	// test ROM doesn't look empty.
	s.NotEqual(uint8(0), c.ROM.DirectGet(0x100))

	s.Equal(uint8(AppleSoft&0xFF), c.Main.Get(BootVector))
	s.Equal(uint8(AppleSoft>>8), c.Main.Get(BootVector+1))
}

func (s *a2Suite) TestReset() {
	c := NewComputer(123)

	c.Reset()

	s.Equal(
		mos.INTERRUPT|mos.BREAK|mos.UNUSED,
		c.CPU.P,
	)
	s.Equal(c.CPU.Get16(ResetPC), c.CPU.PC)
	s.Equal(uint8(0xFF), c.CPU.S)
	s.True(c.State.Bool(a2state.DisplayText))
}
