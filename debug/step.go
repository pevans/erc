package debug

import (
	"fmt"

	"github.com/pevans/erc/a2"
	"github.com/pevans/erc/asm"
)

func step(comp *a2.Computer, tokens []string) {
	var (
		step = 1
		err  error
	)

	if len(tokens) >= 2 {
		step, err = integer(tokens[1])
		if err != nil {
			say(fmt.Sprintf("invalid command: %v", err))
			return
		}
	}

	if step < 0 {
		say("invalid command: step must be positive")
		return
	}

	callmap := make(asm.CallMap)

	for i := 0; i < step; i++ {
		if _, err := comp.Process(); err != nil {
			panic(fmt.Sprintf("could not step instruction: %v", err))
		}

		callmap.Add(comp.CPU.LastInstruction())
	}

	callmap.WriteToFile("callmap.txt")

	say(fmt.Sprintf("executed %v times, current state is now", step))
	status(comp)
}
