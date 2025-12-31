package a2

import (
	"fmt"
	"time"

	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/asm"
	"github.com/pevans/erc/internal/metrics"
	"github.com/pevans/erc/memory"
)

func diskUseDefaults(c *Computer) {
	c.State.SetAny(a2state.DiskComputer, c) // :cry:
}

func diskReadWrite(addr int, val *uint8, stm *memory.StateMap) {
	var (
		nib       = uint8(addr & 0xF)
		c         = stm.Any(a2state.DiskComputer).(*Computer)
		debugging = stm.Bool(a2state.DebuggerLookAhead)
	)

	switch nib {
	case 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7:
		if !debugging {
			c.SelectedDrive.SwitchPhase(int(nib))
			metrics.Increment(fmt.Sprintf("disk_switch_phase_%01x", nib), 1)
		}

		*val = c.SelectedDrive.RandomByte()

	case 0x8:
		// Turn both drives off
		if !debugging {
			c.Drive1.StopMotor()
			c.Drive2.StopMotor()
			c.ClockEmulator.FullSpeed = false
			metrics.Increment("disk_drives_off", 1)
		}

		*val = c.SelectedDrive.RandomByte()

	case 0x9:
		// Turn only the selected drive on
		if !debugging {
			c.SelectedDrive.StartMotor()

			// While the drive is on, we want to emulate without regard to
			// cycle-timing. This is so that any timing loops in software
			// (which would assume a disk is spinning and is trying to time
			// reads to when certain bytes should be found) can be quickly
			// executed, rather than require the user to wait.
			//
			// Another way of putting this is that it's not a goal of ours to
			// perfectly emulate disk spin, and because we can't help the fact
			// that timing loops exist in the software, this is our
			// compromise.
			c.ClockEmulator.FullSpeed = true

			metrics.Increment("disk_selected_drive_online", 1)
		}

		*val = c.SelectedDrive.RandomByte()

	case 0xA:
		// Set the selected drive to drive 1
		if !debugging {
			c.SelectedDrive = c.Drive1
			metrics.Increment("disk_drive_1", 1)
		}

		*val = c.SelectedDrive.RandomByte()

	case 0xB:
		// Set the selected drive to drive 2
		if !debugging {
			c.SelectedDrive = c.Drive2
			metrics.Increment("disk_drive_2", 1)
		}

		*val = c.SelectedDrive.RandomByte()

	case 0xC:
		// This is the SHIFT operation, which might write a byte, but might
		// also read a byte, depending on the drive state.

		if c.SelectedDrive.ReadMode() || c.SelectedDrive.WriteProtected() {
			// Record this now for the disk log because a read on the drive
			// will alter the sector pos
			sectorPos := c.SelectedDrive.sectorPos

			c.SelectedDrive.LoadLatch()
			*val = c.SelectedDrive.ReadLatch()
			c.SelectedDrive.Shift(1)

			if debugging {
				c.SelectedDrive.Shift(-1)
			}

			if c.diskLog != nil {
				c.diskLog.Add(&asm.DiskOp{
					Mode:        asm.DiskRead,
					Elapsed:     time.Since(c.BootTime),
					HalfTrack:   c.SelectedDrive.trackPos,
					Sector:      sectorPos,
					Byte:        *val,
					Instruction: c.CPU.ThisInstruction(),
				})
			}

			metrics.Increment("disk_read", 1)
		} else if c.SelectedDrive.WriteMode() {
			// Write the value currently in the latch
			if !debugging {
				sectorPos := c.SelectedDrive.sectorPos

				c.SelectedDrive.WriteLatch()
				c.SelectedDrive.Shift(1)

				if c.diskLog != nil {
					c.diskLog.Add(&asm.DiskOp{
						Mode:        asm.DiskWrite,
						Elapsed:     time.Since(c.BootTime),
						HalfTrack:   c.SelectedDrive.trackPos,
						Sector:      sectorPos,
						Byte:        c.SelectedDrive.latch,
						Instruction: c.CPU.ThisInstruction(),
					})
				}

				metrics.Increment("disk_write", 1)
			}
		} else {
			if !debugging {
				metrics.Increment("disk_failed_readwrites", 1)
			}
		}

	case 0xD:
		// Set the latch value (for writes) to val
		if !c.SelectedDrive.MotorOn() {
			break
		}

		if !debugging {
			if c.SelectedDrive.WriteMode() {
				c.SelectedDrive.latch = *val
				metrics.Increment("disk_write_latch", 1)
			} else {
				*val = c.SelectedDrive.latch
				metrics.Increment("disk_load_latch", 1)
			}
		}

	case 0xE:
		// Set the selected drive mode to read
		if !debugging {
			c.SelectedDrive.SetReadMode()
			metrics.Increment("disk_read_mode", 1)
		}

		*val = c.SelectedDrive.RandomByte()

		// We also need to return the state of write protection in bit 7
		if c.SelectedDrive.WriteProtected() {
			*val |= 0x80
		} else {
			*val &^= 0x80
		}

	case 0xF:
		// Set the selected drive mode to write
		if !debugging {
			c.SelectedDrive.SetWriteMode()
			metrics.Increment("disk_write_mode", 1)
		}

		*val = c.SelectedDrive.RandomByte()
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
