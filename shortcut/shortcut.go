package shortcut

import (
	"fmt"

	"github.com/pevans/erc/a2"
	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/gfx"
	"github.com/pevans/erc/input"
	"github.com/pevans/erc/obj"
)

func Check(ev input.Event, comp *a2.Computer) (bool, error) {
	// Check for prefix key (Control-A)
	if ev.Modifier == input.ModControl && (ev.Key == 'a' || ev.Key == 'A') {
		w, h := comp.Dimensions()
		gfx.PrefixOverlay.Show(int(w), int(h))
		return true, nil
	}

	// If prefix overlay is not active, pass through the event
	if !gfx.PrefixOverlay.IsActive() {
		return false, nil
	}

	// Prefix is active, process the shortcut key and hide the overlay
	gfx.PrefixOverlay.Hide()

	switch ev.Key {
	case 'b', 'B':
		comp.State.SetBool(a2state.Debugger, true)
		gfx.ShowStatus(obj.DebugPNG())
		return true, nil

	case 'w', 'W':
		comp.SelectedDrive.ToggleWriteProtect()
		fmt.Printf("write protect for drive 1 is now %v\n", comp.SelectedDrive.WriteProtected())
		return true, nil

	case 'd', 'D':
		if err := comp.LoadNext(); err != nil {
			return false, fmt.Errorf("could not load file: %w", err)
		}
		return true, nil
	}

	return false, nil
}
