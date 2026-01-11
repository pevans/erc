package mos_test

import (
	"testing"

	"github.com/pevans/erc/memory"
	"github.com/pevans/erc/mos"
	"github.com/stretchr/testify/assert"
)

func TestCycleCounter(t *testing.T) {
	c := new(mos.CPU)
	seg := memory.NewSegment(0x10000)
	c.RMem = seg
	c.WMem = seg
	c.State = new(memory.StateMap)

	// NOP (0xEA) takes 2 cycles
	seg.Set(0, 0xEA)
	c.PC = 0

	assert.Equal(t, uint64(0), c.CycleCounter())
	assert.NoError(t, c.Execute())
	assert.Equal(t, uint64(2), c.CycleCounter())

	// Execute another NOP
	seg.Set(1, 0xEA)
	assert.NoError(t, c.Execute())
	assert.Equal(t, uint64(4), c.CycleCounter())
}
