package asmrec

import (
	"io"

	"github.com/pevans/erc/clog"
)

// A Recorder is an interface which allows you to record assembly
// instructions that have been executed.
type Recorder interface {
	Record(*clog.Channel)
}

var (
	reclog   *clog.Channel
	shutdown = make(chan bool)
)

func Init(w io.Writer) {
	reclog = clog.NewChannel(w)
	go reclog.WriteLoop(shutdown)
}

func Shutdown() {
	if Available() {
		shutdown <- true
	}
}

func Record(s string) {
	reclog.Printf(s)
}

func Available() bool {
	return reclog != nil
}
