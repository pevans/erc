package a2

import (
	"fmt"

	"github.com/pevans/erc/input"
	"github.com/pevans/erc/internal/metrics"
)

const InstructionMapName = "./instruction_map.asm"

// Shutdown will execute whatever is necessary to basically cease operation of
// the computer.
func (c *Computer) Shutdown() error {
	// It'd be bad if we tried to shutdown more than once, and that is
	// possible if ebiten's Update() call issued many shutdown requests
	c.ShutdownMutex.Lock()
	defer c.ShutdownMutex.Unlock()

	// We already tried a shutdown, so don't do it again
	if c.WillShutDown {
		return nil
	}

	c.WillShutDown = true

	input.Shutdown()

	if c.MetricsFileName != "" {
		err := metrics.WriteToFile(c.MetricsFileName)
		if err != nil {
			return err
		}
	}

	if c.InstructionMap != nil {
		err := c.InstructionMap.WriteToFile(c.InstructionMapFileName)
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
		if err := c.diskLog.WriteToFile(c.diskLogFileName); err != nil {
			return err
		}
	}

	if c.screenLog != nil {
		if err := c.screenLog.WriteToFile(c.screenLogFileName); err != nil {
			return err
		}
	}

	if c.AudioLog != nil {
		if err := c.AudioLog.WriteToFile(c.audioLogFileName); err != nil {
			return err
		}
	}

	if c.instDiffMap != nil {
		if err := c.instDiffMap.WriteToFile(c.instDiffMapFileName); err != nil {
			return err
		}
	}

	if err := c.Drive1.Save(); err != nil {
		return fmt.Errorf("could not save image: %w", err)
	}

	if err := c.Drive2.Save(); err != nil {
		return fmt.Errorf("could not save image: %w", err)
	}

	return nil
}
