package a2disk

import (
	"fmt"
	"time"

	"github.com/pevans/erc/a2/a2drive"
	"github.com/pevans/erc/a2/a2speaker"
	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/clock"
	"github.com/pevans/erc/elog"
	"github.com/pevans/erc/internal/metrics"
	"github.com/pevans/erc/memory"
)

const (
	diskBase = 0xC0E0
)

// Computer is an interface for accessing the computer's disk-related state
// and methods. This allows the disk switches to work with the computer
// without creating a circular dependency.
type Computer interface {
	Drive(n int) *a2drive.Drive
	SelectedDrive() *a2drive.Drive
	SelectDrive(n int)
	ClockEmulator() *clock.Emulator
	Speaker() a2speaker.Speaker
	LogDiskOp(op *elog.DiskOp)
	CPUCurrentInstructionShort() string
	StartTime() time.Time
}

// ReadSwitches returns the list of disk switch addresses that support reads.
func ReadSwitches() []int {
	switches := make([]int, 0x10)
	for i := range 0x10 {
		switches[i] = diskBase + i
	}
	return switches
}

// WriteSwitches returns the list of disk switch addresses that support
// writes.
func WriteSwitches() []int {
	switches := make([]int, 0x10)
	for i := range 0x10 {
		switches[i] = diskBase + i
	}
	return switches
}

func readWrite(addr int, val *uint8, stm *memory.StateMap) {
	var (
		nib       = uint8(addr & 0xF)
		c         = stm.Any(a2state.Computer).(Computer)
		debugging = stm.Bool(a2state.DebuggerLookAhead)
	)

	switch nib {
	case 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7:
		if !debugging {
			c.SelectedDrive().SwitchPhase(int(nib))
			metrics.Increment(fmt.Sprintf("disk_switch_phase_%01x", nib), 1)
		}

		*val = c.SelectedDrive().RandomByte()

	case 0x8:
		// Turn both drives off
		if !debugging {
			c.Drive(1).StopMotor()
			c.Drive(2).StopMotor()
			c.ClockEmulator().SetFullSpeed(false)
			metrics.Increment("disk_drives_off", 1)
		}

		*val = c.SelectedDrive().RandomByte()

	case 0x9:
		// Turn only the selected drive on
		if !debugging {
			c.SelectedDrive().StartMotor()

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
			//
			// To avoid weirdness with any sound being generated, if the
			// speaker is or was recently toggled, we'll not set fullspeed.
			if c.Speaker() == nil || !c.Speaker().IsActive() {
				c.ClockEmulator().SetFullSpeed(true)
			}

			metrics.Increment("disk_selected_drive_online", 1)
		}

		*val = c.SelectedDrive().RandomByte()

	case 0xA:
		// Set the selected drive to drive 1
		if !debugging {
			c.SelectDrive(1)
			metrics.Increment("disk_drive_1", 1)
		}

		*val = c.SelectedDrive().RandomByte()

	case 0xB:
		// Set the selected drive to drive 2
		if !debugging {
			c.SelectDrive(2)
			metrics.Increment("disk_drive_2", 1)
		}

		*val = c.SelectedDrive().RandomByte()

	case 0xC:
		// This is the SHIFT operation, which might write a byte, but might
		// also read a byte, depending on the drive state.

		if c.SelectedDrive().ReadMode() || c.SelectedDrive().WriteProtected() {
			c.SelectedDrive().LoadLatch()
			*val = c.SelectedDrive().ReadLatch()

			c.LogDiskOp(&elog.DiskOp{
				Mode:           elog.DiskRead,
				Elapsed:        time.Since(c.StartTime()),
				Track:          c.SelectedDrive().Track(),
				Sector:         c.SelectedDrive().Sector(),
				SectorPosition: c.SelectedDrive().SectorPosition(),
				Byte:           *val,
				Instruction:    c.CPUCurrentInstructionShort(),
			})

			c.SelectedDrive().Shift(1)

			if debugging {
				c.SelectedDrive().Shift(-1)
			}

			metrics.Increment("disk_read", 1)
		} else if c.SelectedDrive().WriteMode() {
			// Write the value currently in the latch
			if !debugging {
				c.SelectedDrive().WriteLatch()

				c.LogDiskOp(&elog.DiskOp{
					Mode:           elog.DiskWrite,
					Elapsed:        time.Since(c.StartTime()),
					Track:          c.SelectedDrive().Track(),
					Sector:         c.SelectedDrive().Sector(),
					SectorPosition: c.SelectedDrive().SectorPosition(),
					Byte:           c.SelectedDrive().PeekLatch(),
					Instruction:    c.CPUCurrentInstructionShort(),
				})

				c.SelectedDrive().Shift(1)

				metrics.Increment("disk_write", 1)
			}
		} else {
			if !debugging {
				metrics.Increment("disk_failed_readwrites", 1)
			}
		}

	case 0xD:
		// Set the latch value (for writes) to val
		if !c.SelectedDrive().MotorOn() {
			break
		}

		if !debugging {
			if c.SelectedDrive().WriteMode() {
				c.SelectedDrive().SetLatch(*val)
				metrics.Increment("disk_write_latch", 1)
			} else {
				*val = c.SelectedDrive().PeekLatch()
				metrics.Increment("disk_load_latch", 1)
			}
		}

	case 0xE:
		// Set the selected drive mode to read
		if !debugging {
			c.SelectedDrive().SetReadMode()
			metrics.Increment("disk_read_mode", 1)
		}

		*val = c.SelectedDrive().RandomByte()

		// We also need to return the state of write protection in bit 7
		if c.SelectedDrive().WriteProtected() {
			*val |= 0x80
		} else {
			*val &^= 0x80
		}

	case 0xF:
		// Set the selected drive mode to write
		if !debugging {
			c.SelectedDrive().SetWriteMode()
			metrics.Increment("disk_write_mode", 1)
		}

		*val = c.SelectedDrive().RandomByte()
	}
}

// SwitchRead handles reads from disk controller soft switches.
func SwitchRead(addr int, stm *memory.StateMap) uint8 {
	// With reads, we pass a byte value for the ReadWrite function to modify.
	val := uint8(0)

	readWrite(addr, &val, stm)

	return val
}

// SwitchWrite handles writes to disk controller soft switches.
func SwitchWrite(addr int, val uint8, stm *memory.StateMap) {
	// Compared to Read, we pass the val exactly as it comes in.
	readWrite(addr, &val, stm)
}

// UseDefaults sets up the default state for the disk controller. Note: The
// computer reference is stored in the state map by memUseDefaults.
func UseDefaults(_ *memory.StateMap) {
	// Nothing to initialize for disk controller
}
