package a2

import (
	"os"

	"github.com/pevans/erc/pkg/mach"
)

// NewEmulator returns a new Emulator which is configured to behave as
// an Apple II.
func NewEmulator(instLogFile *os.File) *mach.Emulator {
	comp := NewComputer()
	comp.CPU.RecFile = instLogFile

	emu := &mach.Emulator{
		Booter:    comp,
		Drawer:    comp,
		Loader:    comp,
		Ender:     comp,
		Processor: comp,
	}

	return emu
}
