package main

import (
	"time"

	"github.com/pevans/erc/pkg/a2"
	"github.com/pevans/erc/pkg/boot"
	"github.com/pkg/errors"
)

// processLoop executes a process loop whereby we simply execute instructions
// from the CPU endlessly. The only way this can end is if the CPU returns an
// error, or if some external process issues an Exit command to the OS.
func processLoop(comp *a2.Computer, log *boot.Logger, delay time.Duration) {
	for {
		if err := comp.Process(); err != nil {
			log.Fatal(errors.Wrap(err, "process execution failed"))
			return
		}

		time.Sleep(delay)
	}
}
