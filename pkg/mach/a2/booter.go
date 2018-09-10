package a2

import (
	"github.com/pevans/erc/pkg/mach"
	"github.com/pevans/erc/pkg/obj"
	"github.com/pevans/erc/pkg/proc/mos65c02"
)

const (
	// AppleSoft holds the address of AppleSoft, which is the BASIC
	// system that is built into Apple 2.
	AppleSoft = mach.DByte(0xE000)

	// ResetPC is the address that the processor jumps to when it is
	// reset.
	ResetPC = mach.DByte(0xFFFC)

	// BootVector is the location in memory that the operating system
	// is designed to jump to after the intial boot sequence occurs.
	BootVector = mach.DByte(0x03F2)
)

// Boot steps through the boot procedures for the Apple II computer.
// This may also be called a cold start of the computer, and this occurs
// only when the computer is switched from a powered-off to a powered-on
// state.
func (c *Computer) Boot() error {
	// Fetch the slice of bytes for system ROM and for peripheral ROM
	// (they go to together).
	rom, err := obj.Slice(0, RomMemorySize)
	if err != nil {
		return err
	}

	// Copy the ROM bytes into the ROM segment.
	err = c.ROM.CopySlice(0, RomMemorySize, rom)
	if err != nil {
		return err
	}

	// Set the initial reset vector to point to the AppleSoft BASIC system.
	c.Main.Set(BootVector, mach.Byte(AppleSoft&0xFF))
	c.Main.Set(BootVector+1, mach.Byte(AppleSoft>>8))

	// Now run the warm start code.
	c.Reset()

	return nil
}

// Reset will run through a warm start of the computer, which are
// procedures that will execute whenever the computer is reset but not
// powered off.
func (c *Computer) Reset() {
	// Set the initial status of the CPU
	c.CPU.P = mos65c02.NEGATIVE | mos65c02.OVERFLOW | mos65c02.INTERRUPT | mos65c02.ZERO | mos65c02.CARRY

	// Jump to the reset PC address
	c.CPU.PC = ResetPC

	// When reset, the stack goes to its top (which is the end of the
	// stack page).
	c.CPU.S = 0xFF

	// Set our initial memory mode
	c.MemMode = MemDefault | MemSlotCxROM
}
