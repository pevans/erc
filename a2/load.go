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

	if err := c.Drive1.Load(r, fileName); err != nil {
		return errors.Wrapf(err, "could not read file: %s", fileName)
	}

	c.diskLog = nil

	if c.State.Bool(a2state.DebugImage) {
		c.InstructionLog = asm.NewCallMap()
		c.InstructionLogFileName = fmt.Sprintf("%v.asm", fileName)

		c.TimeSet = asm.NewTimeset(c.ClockEmulator.TimePerCycle)
		c.TimeSetFileName = fmt.Sprintf("%v.time", fileName)

		c.MetricsFileName = fmt.Sprintf("%v.metrics", fileName)

		go MaybeLogInstructions(c)

		c.diskLog = asm.NewDiskLog(fileName)
		return c.Drive1.Data.WriteFile(fmt.Sprintf("%v.physical", fileName))
	}

	return nil
}

func MaybeLogInstructions(c *Computer) {
	for line := range c.CPU.InstructionChannel {
		if c.InstructionLog != nil {
			c.InstructionLog.Add(line.String())
		}
		if c.TimeSet != nil {
			c.TimeSet.Record(line.ShortString(), line.Cycles)
		}
	}
}
