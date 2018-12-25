package mach

// An Emulator is a struct which contains the means to "run" a computer,
// for all intents and purposes.
type Emulator struct {
	Booter    Booter
	Ender     Ender
	Processor Processor
	Loader    Loader
}
