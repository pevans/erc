package debug

import (
	"fmt"

	"github.com/pevans/erc/a2"
)

func get(comp *a2.Computer, tokens []string) {
	if len(tokens) != 2 {
		say("invalid command: 'get' requires an address")
		return
	}

	addr, err := hex(tokens[1], 16)
	if err != nil {
		say(fmt.Sprintf("invalid address: %v", err))
		return
	}

	val := comp.Get(addr)
	say(fmt.Sprintf("address $%04x: $%02x", addr, val))
}
