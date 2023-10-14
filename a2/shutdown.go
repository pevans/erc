package a2

import (
	"fmt"

	"github.com/pevans/erc/input"
	"github.com/pevans/erc/internal/metrics"
)

// Shutdown will execute whatever is necessary to basically cease operation of
// the computer.
func (c *Computer) Shutdown() error {
	input.Shutdown()

	fmt.Println("--- METRICS ---")
	mets := metrics.Export()
	for name, counter := range mets {
		fmt.Printf("%v = %v\n", name, counter)
	}

	return nil
}
