package debug

import (
	"fmt"
	"slices"

	"github.com/pevans/erc/a2"
	"github.com/pevans/erc/memory"
	"golang.org/x/exp/maps"
)

func status(comp *a2.Computer) {
	var (
		regfmt  = "registers .......... %s"
		nextfmt = "next instruction ... %s"
	)

	say(fmt.Sprintf(regfmt, comp.CPU.Status()))
	say(fmt.Sprintf(nextfmt, comp.CPU.NextInstruction()))

	stateMap := comp.State()
	stateKeys := maps.Keys(stateMap)
	slices.Sort(stateKeys)

	say("--- state map ---")
	for _, key := range stateKeys {
		statefmt := "%v = %v"
		val := stateMap[key]
		if _, isSeg := val.(*memory.Segment); isSeg {
			statefmt = "%v = %p"
		}

		say(fmt.Sprintf(statefmt, key, val))
	}
}
