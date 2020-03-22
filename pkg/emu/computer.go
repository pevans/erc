package emu

import (
	"io"

	"github.com/pevans/erc/pkg/boot"
)

// A Computer is an interface by which architectures can implement the
// ways that we can execute code.
type Computer interface {
	Boot() error
	DrawLoop()
	Load(io.Reader, string) error
	Process() error
	ProcessLoop()
	Shutdown() error
	SetLogger(*boot.Logger)
	SetRecorderWriter(io.Writer)
}
