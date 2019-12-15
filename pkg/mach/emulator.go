package mach

import "os"

// An Emulator is a struct which contains the means to "run" a computer,
// for all intents and purposes.
type Emulator struct {
	Booter    Booter
	Drawer    Drawer
	Ender     Ender
	Loader    Loader
	Processor Processor

	// InstructionLog is the file where records of our assembly
	// instructions are recorded.
	InstructionLog *os.File
}
