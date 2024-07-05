package debug

import (
	"fmt"

	"github.com/pevans/erc/a2"
)

// Simulate a keypres by registering one with the computer with an
// arbitrary ASCII value.
func keypress(comp *a2.Computer, tokens []string) {
	if len(tokens) != 2 {
		say("invalid command: 'keypress' requires a hex ascii value")
		return
	}

	keyCode, err := hex(tokens[1], 8)
	if err != nil {
		say(fmt.Sprintf("invalid value: '%v' (%v)", tokens[1], err))
		return
	}

	comp.PressKey(uint8(keyCode))
}
