package a2speaker

import (
	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/internal/metrics"
	"github.com/pevans/erc/memory"
)

const (
	speakerToggle int = 0xC030
)

// Computer is an interface for accessing the computer's speaker-related state
// and methods. This allows the speaker switches to work with the computer
// without creating a circular dependency.
type Computer interface {
	CycleCounter() uint64
	Speaker() Speaker
}

// ReadSwitches returns the list of speaker switch addresses that support
// reads.
func ReadSwitches() []int {
	return []int{speakerToggle}
}

// SwitchRead handles reads from speaker soft switches.
func SwitchRead(addr int, stm *memory.StateMap) uint8 {
	if addr != speakerToggle {
		return 0
	}

	metrics.Increment("soft_read_speaker_toggle", 1)

	comp := stm.Any(a2state.Computer).(Computer)

	// Toggle the speaker state
	currentState := stm.Bool(a2state.SpeakerState)
	newState := !currentState
	stm.SetBool(a2state.SpeakerState, newState)

	// Push event to the speaker buffer if available
	if comp.Speaker() != nil {
		cycle := comp.CycleCounter()
		comp.Speaker().Push(cycle, newState)
	}

	// Return floating bus value (we'll just return 0 for now)
	return 0
}

// UseDefaults sets up the default state for the speaker. Note: The computer
// reference is stored in the state map by memUseDefaults.
func UseDefaults(stm *memory.StateMap) {
	stm.SetBool(a2state.SpeakerState, false)
}
