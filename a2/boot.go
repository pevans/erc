package a2

import (
	"time"

	"github.com/pevans/erc/mos"
	"github.com/pevans/erc/obj"
)

const (
	// AppleSoft holds the address of AppleSoft, which is the BASIC system
	// that is built into Apple 2.
	AppleSoft = 0xE000

	// ResetPC is the address that the processor jumps to when it is reset.
	ResetPC = 0xFFFC

	// BootVector is the location in memory that the operating system is
	// designed to jump to after the initial boot sequence occurs.
	BootVector = 0x03F2
)

// Boot steps through the boot procedures for the Apple II computer. This may
// also be called a cold start of the computer, and this occurs only when the
// computer is switched from a powered-off to a powered-on state.
func (c *Computer) Boot() error {
	_, err := c.ROM.CopySlice(0, obj.SystemROM())
	if err != nil {
		return err
	}

	_, err = c.ROM.CopySlice(len(obj.SystemROM()), obj.PeripheralROM())
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

// Reset will run through a warm start of the computer, which are procedures
// that will execute whenever the computer is reset but not powered off.
func (c *Computer) Reset() {
	// Set the initial status of the CPU
	c.CPU.P = mos.INTERRUPT | mos.BREAK | mos.UNUSED

	// When reset, the stack goes to its top (which is the end of the stack
	// page).
	c.CPU.S = 0xFF

	// Set our initial memory mode
	bankUseDefaults(c)
	displayUseDefaults(c)
	kbUseDefaults(c)
	memUseDefaults(c)
	pcUseDefaults(c)
	diskUseDefaults(c)
	speakerUseDefaults(c)

	c.BootTime = time.Now()

	// Jump to the reset PC address; note this must happen _after_ we set our
	// modes above, or else we might pull the PC value from the wrong place in
	// memory.
	c.CPU.PC = c.CPU.Get16(ResetPC)
}
