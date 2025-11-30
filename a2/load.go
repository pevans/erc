package a2

import (
	"fmt"
	"io"

	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/asm"
	"github.com/pkg/errors"
)

// Load will load a disk image from the given reader. The given filename is not
// strictly important -- we're not reading from a filesystem, we already have a
// reader -- but it will be used to help us determine what _kind_ of image it is
// (nibble, dos).
func (c *Computer) Load(r io.Reader, fileName string) error {
	if c.diskLog != nil {
		_ = c.diskLog.WriteToFile()
	}

	// We need to save whatever was the previous file
	if c.SelectedDrive.ImageName != "" {
		err := c.SelectedDrive.Save()
		if err != nil {
			return fmt.Errorf("could not save previous image: %w", err)
		}
	}

	if err := c.SelectedDrive.Load(r, fileName); err != nil {
		return errors.Wrapf(err, "could not read file: %s", fileName)
	}

	c.diskLog = nil

	if c.State.Bool(a2state.DebugImage) {
		c.InstructionLog = asm.NewCallMap()
		c.InstructionLogFileName = fmt.Sprintf("%v.asm", fileName)

		// Share the instruction log with the CPU in case it needs to access
		// it for some reason (e.g. for speculation). We could alternatively
		// put this into the state map.
		c.CPU.InstructionLog = c.InstructionLog

		c.TimeSet = asm.NewTimeset(c.ClockEmulator.TimePerCycle)
		c.TimeSetFileName = fmt.Sprintf("%v.time", fileName)

		c.MetricsFileName = fmt.Sprintf("%v.metrics", fileName)

		c.CPU.InstructionChannel = make(chan *asm.Line, 100)
		go MaybeLogInstructions(c)

		c.diskLog = asm.NewDiskLog(fileName)
		return c.SelectedDrive.Data.WriteFile(fmt.Sprintf("%v.physical", fileName))
	}

	return nil
}

func (c *Computer) LoadFirst() error {
	data, filename, err := c.Disks.First()
	if err != nil {
		return fmt.Errorf("could not load next disk: %w", err)
	}

	defer data.Close()

	return c.Load(data, filename)
}

func (c *Computer) LoadNext() error {
	data, filename, err := c.Disks.Next()
	if err != nil {
		return fmt.Errorf("could not load next disk: %w", err)
	}

	defer data.Close()

	return c.Load(data, filename)
}

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
