package debug

import (
	"fmt"

	"github.com/pevans/erc/a2"
)

func dbatch(comp *a2.Computer, tokens []string) {
	if len(tokens) < 2 {
		say("usage: dbatch <start|stop>")
		return
	}

	switch tokens[1] {
	case "start":
		dbatchStart(comp)
	case "stop":
		dbatchStop(comp)
	default:
		say(fmt.Sprintf("unknown dbatch command: %v", tokens[1]))
	}
}

func dbatchStart(comp *a2.Computer) {
	comp.StartDbatch()
	say("debug batch started")
}

func dbatchStop(comp *a2.Computer) {
	if err := comp.StopDbatch(); err != nil {
		say(fmt.Sprintf("error writing debug batch: %v", err))
		return
	}
	say("debug batch stopped")
}
