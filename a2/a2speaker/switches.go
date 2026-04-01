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

// WriteSwitches returns the list of speaker switch addresses that support
// writes.
func WriteSwitches() []int {
	return []int{speakerToggle}
}

// SwitchRead handles reads from speaker soft switches.
func SwitchRead(addr int, stm *memory.StateMap) uint8 {
	if addr != speakerToggle {
		return 0
	}

	toggle(stm)
	return 0
}

// SwitchWrite handles writes to speaker soft switches.
func SwitchWrite(addr int, val uint8, stm *memory.StateMap) {
	if addr != speakerToggle {
		return
	}

	toggle(stm)
}

func toggle(stm *memory.StateMap) {
	metrics.Increment("soft_read_speaker_toggle", 1)

	comp := stm.Any(a2state.Computer).(Computer)

	currentState := stm.Bool(a2state.SpeakerState)
	newState := !currentState
	stm.SetBool(a2state.SpeakerState, newState)

	if spk := comp.Speaker(); spk != nil {
		cycle := comp.CycleCounter()
		spk.Push(cycle, newState)
	}
}

// UseDefaults sets up the default state for the speaker. Note: The computer
// reference is stored in the state map by memUseDefaults.
func UseDefaults(stm *memory.StateMap) {
	stm.SetBool(a2state.SpeakerState, false)
}
