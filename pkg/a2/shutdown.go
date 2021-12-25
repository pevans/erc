package a2

import (
	"github.com/pevans/erc/pkg/asmrec"
	"github.com/pevans/erc/pkg/disasm"
)

// Shutdown will execute whatever is necessary to basically cease operation of
// the computer.
func (c *Computer) Shutdown() error {
	disasm.Shutdown()
	asmrec.Shutdown()

	return nil
}
