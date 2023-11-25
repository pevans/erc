package debug

import (
	"fmt"

	"github.com/pevans/erc/a2"
)

func set(comp *a2.Computer, tokens []string) {
	if len(tokens) != 3 {
		say("invalid command: 'set' requires an address and value")
		return
	}

	addr, err := hex(tokens[1], 16)
	if err != nil {
		say(fmt.Sprintf("invalid address: %v", err))
		return
	}

	val, err := hex(tokens[2], 8)
	if err != nil {
		say(fmt.Sprintf("invalid value: %v", err))
		return
	}

	comp.Set(addr, uint8(val))
	say(fmt.Sprintf("address $%04x: $%02x written", addr, val))
}
