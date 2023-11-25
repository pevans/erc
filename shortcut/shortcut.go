package shortcut

import (
	"github.com/pevans/erc/a2"
	"github.com/pevans/erc/input"
)

func Check(ev input.Event, comp *a2.Computer) bool {
	if ev.Modifier != input.ModControl {
		return false
	}

	if ev.Key == 'd' {
		comp.Debugger = true
		return true
	}

	return false
}
