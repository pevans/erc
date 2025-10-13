package a2

import (
	"fmt"
	"sort"

	"github.com/pevans/erc/input"
	"github.com/pevans/erc/internal/metrics"
)

const InstructionLogName = "./instruction_log.asm"

// Shutdown will execute whatever is necessary to basically cease operation of
// the computer.
func (c *Computer) Shutdown() error {
	input.Shutdown()

	fmt.Println("--- METRICS ---")

	mets := metrics.Export()
	keys := []string{}

	for name := range mets {
		keys = append(keys, name)
	}

	sort.Strings(keys)
	for _, key := range keys {
		fmt.Printf("%v = %v\n", key, mets[key])
	}

	if c.InstructionLog != nil {
		err := c.InstructionLog.WriteToFile(c.InstructionLogFileName)
		if err != nil {
			return err
		}
	}

	if c.diskLog != nil {
		return c.diskLog.WriteToFile()
	}

	return nil
}
