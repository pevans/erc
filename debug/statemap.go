package debug

import (
	"fmt"
	"maps"
	"slices"

	"github.com/pevans/erc/a2"
	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/memory"
)

func stateMap(comp *a2.Computer) {
	stateMap := comp.State.Map(a2state.KeyToString)
	stateKeys := slices.Sorted(maps.Keys(stateMap))

	for _, key := range stateKeys {
		statefmt := "%20v | %v"
		val := stateMap[key]
		if _, isSeg := val.(*memory.Segment); isSeg {
			statefmt = "%20v | %p"
		}

		say(fmt.Sprintf(statefmt, key, val))
	}
}
