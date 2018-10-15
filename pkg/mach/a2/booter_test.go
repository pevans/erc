package a2

import (
	"testing"

	"github.com/pevans/erc/pkg/mach"
	"github.com/pevans/erc/pkg/obj"
	"github.com/pevans/erc/pkg/proc/mos65c02"
	"github.com/stretchr/testify/assert"
)

func TestBoot(t *testing.T) {
	c := NewComputer()

	rom, err := obj.Slice(0, RomMemorySize)
	assert.Equal(t, nil, err)

	// Boot up
	assert.Equal(t, nil, c.Boot())

	for i := 0; i < 4; i++ {
		assert.Equal(t, rom[i], c.ROM.Get(mach.DByte(i)))
	}

	assert.Equal(t, mach.Byte(AppleSoft&0xFF), c.Main.Get(BootVector))
	assert.Equal(t, mach.Byte(AppleSoft>>8), c.Main.Get(BootVector+1))

	// Test one thing from the Reset() function just to make sure that
	// ran...
	assert.Equal(t, mach.Byte(0xFF), c.CPU.S)
}

func TestReset(t *testing.T) {
	c := NewComputer()
	defp := mos65c02.NEGATIVE | mos65c02.OVERFLOW | mos65c02.INTERRUPT | mos65c02.ZERO | mos65c02.CARRY

	// Note that Reset doesn't return an error, so we can only poke at c
	// state to see if something went wrong.
	c.Reset()

	assert.Equal(t, defp, c.CPU.P)
	assert.Equal(t, ResetPC, c.CPU.PC)
	assert.Equal(t, mach.Byte(0xFF), c.CPU.S)
	assert.Equal(t, MemDefault, c.MemMode)
	assert.Equal(t, BankDefault, c.BankMode)
	assert.Equal(t, PCSlotCxROM, c.PCMode)
}
