package main

import (
	"fmt"
	"log/slog"

	"github.com/pevans/erc/a2"
	"github.com/pevans/erc/debug"
	"github.com/pevans/erc/statemap"
)

// processLoop executes a process loop whereby we simply execute instructions
// from the CPU endlessly. The only way this can end is if the CPU returns an
// error, or if some external process issues an Exit command to the OS.
func processLoop(comp *a2.Computer) {
	for {
		if comp.State.Bool(statemap.Debugger) {
			debug.Prompt(comp)
			continue
		}

		_, err := comp.Process()
		if err != nil {
			slog.Error(fmt.Sprintf("process execution failed: %v", err))
			return
		}
	}
}
