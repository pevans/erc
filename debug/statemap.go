package debug

import (
	"fmt"
	"slices"

	"github.com/pevans/erc/a2"
	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/memory"
	"golang.org/x/exp/maps"
)

func stateMap(comp *a2.Computer) {
	stateMap := comp.State.Map(a2state.KeyToString)
	stateKeys := maps.Keys(stateMap)
	slices.Sort(stateKeys)

	for _, key := range stateKeys {
		statefmt := "%20v | %v"
		val := stateMap[key]
		if _, isSeg := val.(*memory.Segment); isSeg {
			statefmt = "%20v | %p"
		}

		say(fmt.Sprintf(statefmt, key, val))
	}
}
