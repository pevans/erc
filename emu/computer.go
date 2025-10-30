package emu

import (
	"io"

	"github.com/pevans/erc/memory"
)

// A Computer is an interface by which architectures can implement the
// ways that we can execute code.
type Computer interface {
	Boot() error
	Load(io.Reader, string) error

	// Process returns the number of cycles executed and an error status
	Process() (int, error)
	Shutdown() error

	StateMap() *memory.StateMap
}
