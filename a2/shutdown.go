package a2

import (
	"fmt"
	"sort"

	"github.com/pevans/erc/input"
	"github.com/pevans/erc/internal/metrics"
)

// Shutdown will execute whatever is necessary to basically cease operation of
// the computer.
func (c *Computer) Shutdown() error {
	input.Shutdown()

	fmt.Println("--- METRICS ---")

	mets := metrics.Export()
	keys := []string{}

	for name, _ := range mets {
		keys = append(keys, name)
	}

	sort.Strings(keys)
	for _, key := range keys {
		fmt.Printf("%v = %v\n", key, mets[key])
	}

	if c.diskLog != nil {
		c.diskLog.WriteToFile()
	}

	return nil
}
