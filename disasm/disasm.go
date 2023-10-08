package disasm

import (
	"io"

	"github.com/pevans/erc/clog"
)

var (
	dislog    *clog.Channel
	shutdown  = make(chan bool)
	sourceMap = make(map[int]string)
)

// Init will initialize the log channel for writes.
func Init(w io.Writer) {
	dislog = clog.NewChannel(w)
	go dislog.WriteLoop(shutdown)
}

// Available returns true if we are able to log disassembly lines.
func Available() bool {
	return dislog != nil
}

// Shutdown turns off disassembly logging, and writes all of the
// combined disassembled addresses in one batch.
func Shutdown() {
	if Available() {
		for _, rec := range sourceMap {
			dislog.Printf(rec)
		}

		shutdown <- true
		dislog = nil
	}
}

// Map makes a record of the disassembly for a given address in memory.
func Map(addr int, s string) {
	if _, ok := sourceMap[addr]; ok {
		return
	}

	sourceMap[addr] = s
}
