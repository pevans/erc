package shortcut

import (
	"fmt"
	"strconv"

	"github.com/pevans/erc/a2"
	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/gfx"
	"github.com/pevans/erc/input"
	"github.com/pevans/erc/obj"
)

const escapeKey rune = 0x1B

func Check(ev input.Event, comp *a2.Computer) (bool, error) {
	// If paused, ESC resumes; any other key flashes the pause graphic
	if comp.State.Bool(a2state.Paused) {
		if ev.Key == escapeKey {
			comp.ClearKeys()
			comp.State.SetBool(a2state.Paused, false)
			gfx.ShowStatus(obj.ResumePNG())
		} else {
			gfx.ShowStatus(obj.PausePNG())
		}

		return true, nil
	}

	// Check for prefix key (Control-A)
	if ev.Modifier == input.ModControl && (ev.Key == 'a' || ev.Key == 'A') {
		// If prefix is already active, Ctrl-A again sends a literal Ctrl-A
		if gfx.PrefixOverlay.IsActive() {
			gfx.PrefixOverlay.Hide()
			comp.PressKey(0x01) // Ctrl-A = ASCII 0x01
			return true, nil
		}

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
	case escapeKey:
		comp.ClearKeys()
		comp.State.SetBool(a2state.Paused, true)
		gfx.ShowStatus(obj.PausePNG())
		return true, nil

	case '-', '_':
		comp.SpeedDown()
		gfx.ShowStatus(obj.SpeedDownPNG())
		return true, nil

	case '+', '=':
		comp.SpeedUp()
		gfx.ShowStatus(obj.SpeedUpPNG())
		return true, nil

	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		num, _ := strconv.Atoi(string(ev.Key))
		comp.SetStateSlot(num)
		return true, nil

	case 'b', 'B':
		comp.State.SetBool(a2state.Debugger, true)
		gfx.ShowStatus(obj.DebugPNG())
		return true, nil

	case 'l', 'L':
		if err := comp.LoadStateSlot(); err != nil {
			comp.ShowText("could not load state")
		}
		return true, nil

	case 'n', 'N':
		// Both LoadNext and LoadPrevious will take care of showing the
		// correct disk's status graphic
		if err := comp.LoadNext(); err != nil {
			return false, fmt.Errorf("could not load file: %w", err)
		}
		return true, nil

	case 'p', 'P':
		if err := comp.LoadPrevious(); err != nil {
			return false, fmt.Errorf("could not load file: %w", err)
		}
		return true, nil

	case 's', 'S':
		if err := comp.SaveStateSlot(); err != nil {
			comp.ShowText("could not save state")
		}
		return true, nil

	case 'w', 'W':
		comp.SelectedDrive().ToggleWriteProtect()

		if comp.SelectedDrive().WriteProtected() {
			gfx.ShowStatus(obj.WriteProtectedPNG())
		} else {
			gfx.ShowStatus(obj.WriteablePNG())
		}

		return true, nil

	case 'q', 'Q':
		return true, comp.Shutdown()

	case 'v', 'V':
		comp.VolumeToggle()
		if comp.IsMuted() {
			gfx.ShowStatus(obj.VolumeOffPNG())
		}
		return true, nil

	case '[', '{':
		comp.VolumeDown(10)
		if comp.IsMuted() {
			gfx.ShowStatus(obj.VolumeOffPNG())
		} else {
			gfx.ShowStatus(obj.VolumeDownPNG())
		}
		return true, nil

	case ']', '}':
		comp.VolumeUp(10)
		gfx.ShowStatus(obj.VolumeUpPNG())
		return true, nil
	}

	return false, nil
}
