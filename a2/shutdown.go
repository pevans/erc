package a2

import (
	"github.com/pevans/erc/asmrec"
	"github.com/pevans/erc/disasm"
	"github.com/pevans/erc/input"
)

// Shutdown will execute whatever is necessary to basically cease operation of
// the computer.
func (c *Computer) Shutdown() error {
	disasm.Shutdown()
	asmrec.Shutdown()
	input.Shutdown()

	return nil
}