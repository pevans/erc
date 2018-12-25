package a2

import "github.com/pevans/erc/pkg/mach"

// NewEmulator returns a new Emulator which is configured to behave as
// an Apple II.
func NewEmulator() *mach.Emulator {
	comp := NewComputer()

	emu := &mach.Emulator{
		Booter:    comp,
		Loader:    comp,
		Ender:     comp,
		Processor: comp,
	}

	return emu
}
