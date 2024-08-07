package main

import (
	"fmt"
	"log/slog"

	"github.com/peterh/liner"
	"github.com/pevans/erc/a2"
	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/debug"
)

// processLoop executes a process loop whereby we simply execute instructions
// from the CPU endlessly. The only way this can end is if the CPU returns an
// error, or if some external process issues an Exit command to the OS.
func processLoop(comp *a2.Computer) {
	line := liner.NewLiner()
	defer line.Close()

	for {
		if comp.State.Bool(a2state.Debugger) {
			debug.Prompt(comp, line)
			continue
		}

		_, err := comp.Process()
		if err != nil {
			slog.Error(fmt.Sprintf("process execution failed: %v", err))
			return
		}
	}
}
