package debug

import (
	"fmt"

	"github.com/pevans/erc/a2"
)

func reg(comp *a2.Computer, tokens []string) {
	if len(tokens) != 3 {
		say("invalid command: 'reg' requires a register and a value")
		return
	}

	val, err := hex(tokens[2], 16)
	if err != nil {
		say(fmt.Sprintf("invalid value: %v", err))
		return
	}

	switch tokens[1] {
	case "a":
		comp.CPU.A = uint8(val)
		status(comp)
	case "p":
		comp.CPU.P = uint8(val)
		status(comp)
	case "pc":
		comp.CPU.PC = uint16(val)
		status(comp)
	case "s":
		comp.CPU.S = uint8(val)
		status(comp)
	case "x":
		comp.CPU.X = uint8(val)
		status(comp)
	case "y":
		comp.CPU.Y = uint8(val)
		status(comp)
	default:
		say(fmt.Sprintf("invalid register: \"%v\"", tokens[1]))
	}
}
