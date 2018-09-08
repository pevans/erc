package mach

// Processor is the interface that abstracts how the computer can
// process operations (opcodes).
type Processor interface {
	Process() error
}
