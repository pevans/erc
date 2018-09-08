package a2

import "github.com/pevans/erc/pkg/mach"

func NewEmulator() *mach.Emulator {
	comp := NewComputer()
	emu := &mach.Emulator{
		Booter: comp,
	}

	return emu
}
