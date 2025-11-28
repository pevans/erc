package a2

import (
	"fmt"

	"github.com/pevans/erc/input"
	"github.com/pevans/erc/internal/metrics"
)

const InstructionLogName = "./instruction_log.asm"

// Shutdown will execute whatever is necessary to basically cease operation of
// the computer.
func (c *Computer) Shutdown() error {
	input.Shutdown()

	if err := c.Drive1.Save(); err != nil {
		return fmt.Errorf("could not save image: %w", err)
	}

	if c.MetricsFileName != "" {
		err := metrics.WriteToFile(c.MetricsFileName)
		if err != nil {
			return err
		}
	}

	if c.InstructionLog != nil {
		err := c.InstructionLog.WriteToFile(c.InstructionLogFileName)
		if err != nil {
			return err
		}
	}

	if c.TimeSet != nil {
		err := c.TimeSet.WriteToFile(c.TimeSetFileName)
		if err != nil {
			return err
		}
	}

	if c.diskLog != nil {
		return c.diskLog.WriteToFile()
	}

	return nil
}
