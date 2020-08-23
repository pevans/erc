package mos65c02

import (
	"testing"

	"github.com/pevans/erc/pkg/data"
	"github.com/stretchr/testify/assert"
)

func TestInstructionString(t *testing.T) {
	var i Instruction = Brk
	assert.Equal(t, "BRK", i.String())
}

func TestExecute(t *testing.T) {
	c := new(CPU)
	c.RMem = data.NewSegment(0x10000)
	c.WMem = data.NewSegment(0x10000)

	// In just a blank default template, this should error out.
	assert.Error(t, c.Execute())
}
