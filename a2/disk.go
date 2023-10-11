package a2

import (
	"github.com/pevans/erc/internal/metrics"
	"github.com/pevans/erc/memory"
)

const (
	diskComputer = 600
)

func diskUseDefaults(c *Computer) {
	c.state.SetAny(diskComputer, c) // :cry:
}

func diskReadWrite(addr int, val *uint8, stm *memory.StateMap) {
	var (
		nib = addr & 0xF
		c   = stm.Any(diskComputer).(*Computer)
	)

	switch nib {
	case 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7:
		// Set the drive phase, thus adjusting the track position
		c.SelectedDrive.SwitchPhase(nib)

	case 0x8:
		// Turn both drives on
		c.Drive1.Online = true
		c.Drive2.Online = true

	case 0x9:
		// Turn only the selected drive on
		c.SelectedDrive.Online = true

	case 0xA:
		// Set the selected drive to drive 1
		c.SelectedDrive = c.Drive1

	case 0xB:
		// Set the selected drive to drive 2
		c.SelectedDrive = c.Drive2

	case 0xC:
		if c.SelectedDrive.Mode == ReadMode || c.SelectedDrive.WriteProtect {
			*val = c.SelectedDrive.Read()
			metrics.Increment("disk_reads", 1)
		} else if c.SelectedDrive.Mode == WriteMode {
			// Write the value currently in the latch
			c.SelectedDrive.Write()
			metrics.Increment("disk_writes", 1)
		} else {
			metrics.Increment("failed_disk_readwrites", 1)
		}

	case 0xD:
		// Set the latch value (for writes) to val
		if c.SelectedDrive.Mode == WriteMode {
			c.SelectedDrive.Latch = *val
			metrics.Increment("disk_latches", 1)
		} else {
			metrics.Increment("failed_disk_latches", 1)
		}

	case 0xE:
		// Set the selected drive mode to read
		c.SelectedDrive.Mode = ReadMode

	case 0xF:
		// Set the selected drive mode to write
		c.SelectedDrive.Mode = WriteMode
	}

	if nib%2 == 0 {
		*val = c.SelectedDrive.Latch
	}
}

func diskRead(addr int, stm *memory.StateMap) uint8 {
	// With reads, we pass a byte value for the ReadWrite function to
	// modify.
	val := uint8(0)

	diskReadWrite(addr, &val, stm)

	return val
}

func diskWrite(addr int, val uint8, stm *memory.StateMap) {
	// Compared to Read, we pass the val exactly as it comes in.
	diskReadWrite(addr, &val, stm)
}
