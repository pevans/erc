package disasm

import (
	"io"

	"github.com/pevans/erc/pkg/clog"
)

var (
	dislog    *clog.Channel
	shutdown  = make(chan bool)
	sourceMap = make(map[int]string)
)

func Init(w io.Writer) {
	dislog = clog.NewChannel(w)
	go dislog.WriteLoop(shutdown)
}

func Available() bool {
	return dislog != nil
}

func Shutdown() {
	if Available() {
		for _, rec := range sourceMap {
			dislog.Printf(rec)
		}

		shutdown <- true
	}
}

func Map(addr int, s string) {
	if _, ok := sourceMap[addr]; ok {
		return
	}

	sourceMap[addr] = s
}
