package a2memory_test

import (
	"testing"

	"github.com/pevans/erc/a2/a2memory"
	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/memory"
	"github.com/stretchr/testify/assert"
)

func TestUseDefaults(t *testing.T) {
	stm := memory.NewStateMap()
	main := memory.NewSegment(0x10000)
	aux := memory.NewSegment(0x10000)

	a2memory.UseDefaults(stm, main, aux)

	assert.False(t, stm.Bool(a2state.MemReadAux))
	assert.False(t, stm.Bool(a2state.MemWriteAux))
}

func TestSwitchRead(t *testing.T) {
	var (
		c013       = 0xC013
		c014       = 0xC014
		hi   uint8 = 0x80
		lo   uint8 = 0x00
		stm        = memory.NewStateMap()
	)

	t.Run("read profile", func(t *testing.T) {
		stm.SetBool(a2state.MemReadAux, true)
		assert.Equal(t, hi, a2memory.SwitchRead(c013, stm))

		stm.SetBool(a2state.MemReadAux, false)
		assert.Equal(t, lo, a2memory.SwitchRead(c013, stm))
	})

	t.Run("write profile", func(t *testing.T) {
		stm.SetBool(a2state.MemWriteAux, true)
		assert.Equal(t, hi, a2memory.SwitchRead(c014, stm))

		stm.SetBool(a2state.MemWriteAux, false)
		assert.Equal(t, lo, a2memory.SwitchRead(c014, stm))
	})
}

func TestSwitchWrite(t *testing.T) {
	var (
		c002 = 0xC002
		c003 = 0xC003
		c004 = 0xC004
		c005 = 0xC005
		stm  = memory.NewStateMap()
		main = memory.NewSegment(0x10000)
		aux  = memory.NewSegment(0x10000)
	)

	// Set up segments
	stm.SetSegment(a2state.MemMainSegment, main)
	stm.SetSegment(a2state.MemAuxSegment, aux)

	t.Run("set aux works", func(t *testing.T) {
		stm.SetBool(a2state.MemReadAux, false)
		a2memory.SwitchWrite(c003, 0, stm)
		assert.True(t, stm.Bool(a2state.MemReadAux))

		stm.SetBool(a2state.MemWriteAux, false)
		a2memory.SwitchWrite(c005, 0, stm)
		assert.True(t, stm.Bool(a2state.MemWriteAux))
	})

	t.Run("set main works", func(t *testing.T) {
		stm.SetBool(a2state.MemReadAux, true)
		a2memory.SwitchWrite(c002, 0, stm)
		assert.False(t, stm.Bool(a2state.MemReadAux))

		stm.SetBool(a2state.MemWriteAux, true)
		a2memory.SwitchWrite(c004, 0, stm)
		assert.False(t, stm.Bool(a2state.MemWriteAux))
	})
}
