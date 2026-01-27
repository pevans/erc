package debug

import (
	"fmt"
	"strconv"
	"time"

	"github.com/pevans/erc/a2"
	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/gfx"
	"github.com/pevans/erc/obj"
)

func runfor(comp *a2.Computer, tokens []string) {
	if len(tokens) < 2 {
		say("you must provide a duration in seconds")
		return
	}

	seconds, err := strconv.ParseFloat(tokens[1], 64)
	if err != nil {
		say(fmt.Sprintf("invalid duration: %v", tokens[1]))
		return
	}

	if seconds <= 0 {
		say("duration must be positive")
		return
	}

	duration := time.Duration(seconds * float64(time.Second))

	comp.State.SetBool(a2state.Debugger, false)
	gfx.ShowStatus(obj.ResumePNG())
	say(fmt.Sprintf("running for %.2f seconds", seconds))

	// Start a timer to reenter the debugger
	go func() {
		time.Sleep(duration)
		comp.State.SetBool(a2state.Debugger, true)
		gfx.ShowStatus(obj.PausePNG())
	}()
}
