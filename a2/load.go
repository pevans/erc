package a2

import (
	"fmt"
	"io"
	"path"

	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/elog"
	"github.com/pevans/erc/gfx"
	"github.com/pevans/erc/obj"
	"github.com/pkg/errors"
)

// Load will load a disk image from the given reader. The given filename is
// not strictly important -- we're not reading from a filesystem, we already
// have a reader -- but it will be used to help us determine what _kind_ of
// image it is (nibble, dos).
func (c *Computer) Load(r io.Reader, fileName string) error {
	if c.diskLog != nil {
		_ = c.diskLog.WriteToFile(c.diskLogFileName)
	}

	if err := c.SelectedDrive.Save(); err != nil {
		return fmt.Errorf("could not save previous image: %w", err)
	}

	if err := c.SelectedDrive.Load(r, fileName); err != nil {
		return errors.Wrapf(err, "could not read file: %s", fileName)
	}

	c.diskLog = nil

	if c.State.Bool(a2state.DebugImage) {
		c.InstructionLog = elog.NewInstructionMap()
		c.InstructionLogFileName = fmt.Sprintf("%v.asm", fileName)

		// Share the instruction log with the CPU in case it needs to access
		// it for some reason (e.g. for speculation). We could alternatively
		// put this into the state map.
		c.CPU.InstructionLog = c.InstructionLog

		c.TimeSet = elog.NewTimeset(c.ClockEmulator.TimePerCycle())
		c.TimeSetFileName = fmt.Sprintf("%v.time", fileName)

		c.MetricsFileName = fmt.Sprintf("%v.metrics", fileName)

		c.CPU.InstructionChannel = make(chan *elog.Instruction, 100)
		go MaybeLogInstructions(c)

		c.diskLogFileName = fmt.Sprintf("%v.disklog", fileName)
		c.diskLog = elog.NewDiskLog()

		c.screenLog = elog.NewScreenLog()
		c.screenLogFileName = fmt.Sprintf("%v.screen", fileName)
		c.lastScreenCapture = c.BootTime

		c.AudioLog = elog.NewAudioLog()
		c.audioLogFileName = fmt.Sprintf("%v.audio", fileName)

		return c.SelectedDrive.WriteDataToFile(fmt.Sprintf("%v.physical", fileName))
	}

	return nil
}

// LoadFirst will load the first disk in the diskset, regardless of the
// diskset's current index. Note that this will _alter_ the index to become
// that of the first disk. An error is returned if the disk can't be loaded.
func (c *Computer) LoadFirst() error {
	data, filename, err := c.Disks.First()
	if err != nil {
		return fmt.Errorf("could not load next disk: %w", err)
	}

	defer data.Close() //nolint:errcheck

	return c.Load(data, filename)
}

// LoadNext will load the next disk in the computer's diskset. If that can't
// be done, an error is returned. A status graphic will be displayed to
// indicate the new disk that has been loaded.
func (c *Computer) LoadNext() error {
	data, filename, err := c.Disks.Next()
	if err != nil {
		return fmt.Errorf("could not load next disk: %w", err)
	}

	defer data.Close() //nolint:errcheck

	if png := diskPNG(c.Disks.CurrentIndex()); png != nil {
		gfx.ShowStatus(png)
	}

	c.ShowText(path.Base(filename))

	return c.Load(data, filename)
}

// LoadPrevious will load the previous disk in the computer's diskset. If that
// can't be done, an error is returned. A status graphic will be displayed to
// indicate the new disk that has been loaded.
func (c *Computer) LoadPrevious() error {
	data, filename, err := c.Disks.Previous()
	if err != nil {
		return fmt.Errorf("could not load previous disk: %w", err)
	}

	defer data.Close() //nolint:errcheck

	if png := diskPNG(c.Disks.CurrentIndex()); png != nil {
		gfx.ShowStatus(png)
	}

	c.ShowText(path.Base(filename))

	return c.Load(data, filename)
}

// diskPNG returns the Disk#PNG status graphic, where # corresponds to some
// disk image. Since the disk images are numbered 1-10, we have to adjust what
// we return to align to the 0-based index from the disket.
func diskPNG(index int) []byte {
	switch index {
	case 0:
		return obj.Disk1PNG()
	case 1:
		return obj.Disk2PNG()
	case 2:
		return obj.Disk3PNG()
	case 3:
		return obj.Disk4PNG()
	case 4:
		return obj.Disk5PNG()
	case 5:
		return obj.Disk6PNG()
	case 6:
		return obj.Disk7PNG()
	case 7:
		return obj.Disk8PNG()
	case 8:
		return obj.Disk9PNG()
	case 9:
		return obj.Disk10PNG()
	}
	return nil
}

// MaybeLogInstructions will, if an InstructionLog is available and if
// messages are available on the CPU's InstructionChannel, record those
// instructions in the instruction log. If a TimeSet is also available, we
// will also record a timeset entry there.
func MaybeLogInstructions(c *Computer) {
	for line := range c.CPU.InstructionChannel {
		if c.InstructionLog != nil {
			c.InstructionLog.Add(line)
		}
		if c.TimeSet != nil {
			c.TimeSet.Record(line.ShortString(), line.Cycles)
		}
	}
}
