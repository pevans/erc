package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/pevans/erc/a2"
	"github.com/pevans/erc/debug"
)

// processLoop executes a process loop whereby we simply execute instructions
// from the CPU endlessly. The only way this can end is if the CPU returns an
// error, or if some external process issues an Exit command to the OS.
func processLoop(comp *a2.Computer, delay time.Duration) {
	for {
		if comp.Debugger {
			debug.Prompt(comp)
			continue
		}

		if err := comp.Process(); err != nil {
			slog.Error(fmt.Sprintf("process execution failed: %v", err))
			return
		}

		time.Sleep(delay)
	}
}
