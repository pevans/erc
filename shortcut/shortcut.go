package shortcut

import (
	"fmt"

	"github.com/pevans/erc/a2"
	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/input"
)

func Check(ev input.Event, comp *a2.Computer) (bool, error) {
	if ev.Modifier != input.ModControl {
		return false, nil
	}

	if ev.Key == 's' || ev.Key == 'S' {
		return true, comp.Drive1.Save()
	}

	if ev.Key == 'd' || ev.Key == 'D' {
		comp.State.SetBool(a2state.Debugger, true)
		return true, nil
	}

	if ev.Key == 'o' || ev.Key == 'O' {
		fmt.Println("hi!")
		if err := comp.LoadNext(); err != nil {
			return false, fmt.Errorf("could not load file: %w", err)
		}
	}

	return false, nil
}
