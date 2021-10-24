package a2

import (
	"github.com/pevans/erc/pkg/disasm"
	"github.com/pevans/erc/pkg/mos65c02"
	"github.com/pevans/erc/pkg/obj"
)

const (
	// AppleSoft holds the address of AppleSoft, which is the BASIC
	// system that is built into Apple 2.
	AppleSoft = 0xE000

	// ResetPC is the address that the processor jumps to when it is
	// reset.
	ResetPC = 0xFFFC

	// BootVector is the location in memory that the operating system
	// is designed to jump to after the initial boot sequence occurs.
	BootVector = 0x03F2
)

// Boot steps through the boot procedures for the Apple II computer.
// This may also be called a cold start of the computer, and this occurs
// only when the computer is switched from a powered-off to a powered-on
// state.
func (c *Computer) Boot(disFile string) error {
	c.CPU.SMap = disasm.NewSourceMap(disFile)

	// Fetch the slice of bytes for system ROM and for peripheral ROM
	// (they go to together).
	rom, err := obj.Slice(4, RomMemorySize+4)
	if err != nil {
		return err
	}

	// Copy the ROM bytes into the ROM segment. This copies not only
	// system ROM; it also copies in peripheral ROM.
	_, err = c.ROM.CopySlice(0, rom)
	if err != nil {
		return err
	}

	// Set the initial reset vector to point to the AppleSoft BASIC system.
	c.Main.Set(BootVector, uint8(AppleSoft&0xFF))
	c.Main.Set(BootVector+1, uint8(AppleSoft>>8))

	// Set up all the soft switches we'll need
	c.MapSoftSwitches()

	// Now run the warm start code.
	c.Reset()

	return nil
}

// Reset will run through a warm start of the computer, which are
// procedures that will execute whenever the computer is reset but not
// powered off.
func (c *Computer) Reset() {
	// Set the initial status of the CPU
	//c.CPU.P = mos65c02.NEGATIVE | mos65c02.OVERFLOW | mos65c02.INTERRUPT | mos65c02.ZERO | mos65c02.CARRY
	c.CPU.P = mos65c02.INTERRUPT | mos65c02.BREAK | mos65c02.UNUSED

	// When reset, the stack goes to its top (which is the end of the
	// stack page).
	c.CPU.S = 0xFF

	// Set our initial memory mode
	c.bank.UseDefaults(c)
	c.disp.UseDefaults(c)
	c.kb.UseDefaults(c)
	c.mem.UseDefaults(c)
	c.pc.UseDefaults(c)

	// Jump to the reset PC address; note this must happen _after_ we
	// set our modes above, or else we might pull the PC value from the
	// wrong place in memory.
	c.CPU.PC = c.CPU.Get16(ResetPC)
}
