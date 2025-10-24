package a2

import (
	"fmt"
	"os"
	"time"

	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/internal/metrics"
	"github.com/pevans/erc/memory"
)

type DiskRead struct {
	Elapsed     time.Duration
	HalfTrack   int
	Sector      int
	Byte        uint8
	Instruction string
}

type DiskLog struct {
	Reads []DiskRead
	Name  string
}

func NewDiskLog(name string) *DiskLog {
	log := new(DiskLog)
	log.Name = name

	return log
}

func (l *DiskLog) Add(read *DiskRead) {
	l.Reads = append(l.Reads, *read)
}

func (l *DiskLog) WriteToFile() error {
	file := fmt.Sprintf("%v.disklog", l.Name)

	fp, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer fp.Close() //nolint:errcheck

	for _, read := range l.Reads {
		logLine := fmt.Sprintf(
			"[%-10v] T:%02X S:%01X P:%04X B:$%02X | %v\n",
			read.Elapsed.Round(time.Millisecond), // time since boot
			read.HalfTrack>>1,                    // track
			read.Sector/0x1A0, read.Sector,       // sect and pos
			read.Byte, read.Instruction,
		)

		if _, err := fp.WriteString(logLine); err != nil {
			return err
		}
	}

	return nil
}

func diskUseDefaults(c *Computer) {
	c.State.SetAny(a2state.DiskComputer, c) // :cry:
}

func diskReadWrite(addr int, val *uint8, stm *memory.StateMap) {
	var (
		nib       = uint8(addr & 0xF)
		c         = stm.Any(a2state.DiskComputer).(*Computer)
		lastCycle = stm.Int(a2state.DiskCycleOfLastAccess)
	)

	if lastCycle == 0 {
		lastCycle = c.CPU.CycleCount
	}

	switch nib {
	case 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7:
		// Set the drive phase, thus adjusting the track position
		if !c.SelectedDrive.Online {
			break
		}

		c.SelectedDrive.SwitchPhase(int(nib))

		*val = c.Drive1.RandomByte()

		metrics.Increment(fmt.Sprintf("disk_switch_phase_%01x", nib), 1)

	case 0x8:
		// Turn both drives off
		c.Drive1.Online = false
		c.Drive2.Online = false
		*val = c.Drive1.RandomByte()
		metrics.Increment("disk_drives_off", 1)

	case 0x9:
		// Turn only the selected drive on
		c.SelectedDrive.Online = true
		stm.SetInt(a2state.DiskCycleOfLastAccess, 0)
		*val = c.Drive1.RandomByte()
		metrics.Increment("disk_selected_drive_online", 1)

	case 0xA:
		// Set the selected drive to drive 1
		c.SelectedDrive = c.Drive1
		*val = 0
		metrics.Increment("disk_drive_1", 1)

	case 0xB:
		// Set the selected drive to drive 2
		c.SelectedDrive = c.Drive2
		*val = 0
		metrics.Increment("disk_drive_2", 1)

	case 0xC:
		// This is the SHIFT operation, which might write a byte, but might
		// also read a byte, depending on the drive state.

		// Default behavior is often that we return 0 -- if the drive isn't
		// on, etc.
		*val = 0

		if !c.SelectedDrive.Online {
			break
		}

		stm.SetInt(a2state.DiskCycleOfLastAccess, c.CPU.CycleCount)

		if c.SelectedDrive.Mode == ReadMode || c.SelectedDrive.WriteProtect {
			// Record this now for the disk log because a read on the drive
			// will alter the sector pos
			sectorPos := c.SelectedDrive.SectorPos

			*val = c.SelectedDrive.Read()

			if c.diskLog != nil {
				c.diskLog.Add(&DiskRead{
					Elapsed:     time.Since(c.BootTime),
					HalfTrack:   c.SelectedDrive.TrackPos,
					Sector:      sectorPos,
					Byte:        *val,
					Instruction: c.CPU.ThisInstruction(),
				})
			}

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
		if !c.SelectedDrive.Online {
			break
		}

		if c.SelectedDrive.Mode == WriteMode {
			c.SelectedDrive.Latch = *val
			metrics.Increment("disk_write_latch", 1)
		} else {
			metrics.Increment("disk_failed_latch", 1)
		}

	case 0xE:
		// Set the selected drive mode to read
		c.SelectedDrive.Mode = ReadMode

		*val = c.Drive1.RandomByte()

		// We also need to return the state of write protection
		if c.SelectedDrive.WriteProtect {
			*val = 0x80
		}

		metrics.Increment("disk_read_mode", 1)

	case 0xF:
		// Set the selected drive mode to write
		c.SelectedDrive.Mode = WriteMode
		*val = c.Drive1.RandomByte()
		metrics.Increment("disk_write_mode", 1)
	}

	if nib%2 == 0 {
		//*val = c.SelectedDrive.Latch
		//metrics.Increment("disk_read_latch", 1)
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
