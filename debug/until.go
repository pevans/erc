package debug

import (
	"fmt"
	"strings"

	"github.com/pevans/erc/a2"
)

func until(comp *a2.Computer, tokens []string) {
	const maxIterations = 100000

	if len(tokens) < 2 {
		say("you must provide an instruction to step until")
		return
	}

	instruction := tokens[1]

	if len(instruction) != 3 {
		say(fmt.Sprintf(
			"you must provide a valid instruction (\"%v\" given)", instruction,
		))
		return
	}

	for i := 0; i < maxIterations; i++ {
		if _, err := comp.Process(); err != nil {
			panic(fmt.Sprintf("failed execution while stepping over: %v", err))
		}

		say(comp.CPU.LastInstruction())

		if strings.HasPrefix(comp.CPU.LastInstruction(), instruction) {
			say(fmt.Sprintf("stepped over %v instructions", i+1))
			return
		}
	}

	// If we got here, that means we hit our max iteration limit
	say(fmt.Sprintf(
		"stopped after %v instructions without executing %v",
		maxIterations, instruction,
	))
}
