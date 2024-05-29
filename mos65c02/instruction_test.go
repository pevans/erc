package mos65c02

import (
	"testing"

	"github.com/pevans/erc/memory"
	"github.com/stretchr/testify/assert"
)

func TestInstructionString(t *testing.T) {
	var i Instruction = Brk
	assert.Equal(t, "BRK", i.String())
}

func TestExecute(t *testing.T) {
	c := new(CPU)
	seg := memory.NewSegment(0x10000)
	c.RMem = seg
	c.WMem = seg
	c.State = new(memory.StateMap)

	// In just a blank default template, this should error out.
	assert.NoError(t, c.Execute())
}
