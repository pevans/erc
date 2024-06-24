package a2

import (
	"fmt"

	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/internal/metrics"
	"github.com/pevans/erc/memory"
)

func diskUseDefaults(c *Computer) {
	c.State.SetAny(a2state.DiskComputer, c) // :cry:
}

func diskReadWrite(addr int, val *uint8, stm *memory.StateMap) {
	var (
		nib = addr & 0xF
		c   = stm.Any(a2state.DiskComputer).(*Computer)
	)

	switch nib {
	case 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7:
		// Set the drive phase, thus adjusting the track position
		c.SelectedDrive.SwitchPhase(nib)
		metrics.Increment(fmt.Sprintf("disk_switch_phase_%01x", nib), 1)

	case 0x8:
		// Turn both drives on
		c.Drive1.Online = true
		c.Drive2.Online = true
		metrics.Increment("disk_drive1_drive2_online", 1)

	case 0x9:
		// Turn only the selected drive on
		c.SelectedDrive.Online = true
		metrics.Increment("disk_selected_drive_online", 1)

	case 0xA:
		// Set the selected drive to drive 1
		c.SelectedDrive = c.Drive1
		metrics.Increment("disk_drive_1", 1)

	case 0xB:
		// Set the selected drive to drive 2
		c.SelectedDrive = c.Drive2
		metrics.Increment("disk_drive_2", 1)

	case 0xC:
		if c.SelectedDrive.Mode == ReadMode || c.SelectedDrive.WriteProtect {
			*val = c.SelectedDrive.Read()
			metrics.Increment("disk_read", 1)
		} else if c.SelectedDrive.Mode == WriteMode {
			// Write the value currently in the latch
			c.SelectedDrive.Write()
			metrics.Increment("disk_write", 1)
		} else {
			metrics.Increment("disk_failed_readwrites", 1)
		}

	case 0xD:
		// Set the latch value (for writes) to val
		if c.SelectedDrive.Mode == WriteMode {
			c.SelectedDrive.Latch = *val
			metrics.Increment("disk_write_latch", 1)
		} else {
			metrics.Increment("disk_failed_latch", 1)
		}

	case 0xE:
		// Set the selected drive mode to read
		c.SelectedDrive.Mode = ReadMode
		metrics.Increment("disk_read_mode", 1)

	case 0xF:
		// Set the selected drive mode to write
		c.SelectedDrive.Mode = WriteMode
		metrics.Increment("disk_write_mode", 1)
	}

	if nib%2 == 0 {
		*val = c.SelectedDrive.Latch
		metrics.Increment("disk_read_latch", 1)
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
