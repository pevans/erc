package clock

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/emu"
)

// Emulator is a device which you can use to simulate a slower
// computer's clockspeed.
type Emulator struct {
	// The rate at which we will emulate any process.
	hertz int64

	// For any given rate of hertz, this is the time we would expect to
	// be consumed by a single cycle.
	TimePerCycle time.Duration

	// When the emulator first began
	StartTime time.Time

	// When execution was last resumed (which may equal StartTime, and which
	// may be later than StartTime)
	ResumeTime time.Time

	// The number of cycles we've executed since the start of emulation
	TotalCycles int64

	EnterDebuggerFunc   func()
	CheckBreakpointFunc func()
}

// NewEmulator returns a new emulator based on some number of hertz
func NewEmulator(hz int64) *Emulator {
	emu := &Emulator{
		hertz:        hz,
		TimePerCycle: (1 * time.Second) / time.Duration(hz),
	}

	return emu
}

func (e *Emulator) ProcessLoop(comp emu.Computer) {
	e.StartTime = time.Now()
	e.ResumeTime = e.StartTime
	state := comp.StateMap()

	for {
		e.CheckBreakpointFunc()

		if state.Bool(a2state.Debugger) {
			e.EnterDebuggerFunc()

			// Reset ResumeTime so that we don't think we're far behind on
			// cycle time just because we sat in the debugger for a while
			e.ResumeTime = time.Now()

			continue
		}

		elapsed := time.Since(e.ResumeTime)
		wantedCycles := int64(elapsed / e.TimePerCycle)

		// We know how much time has elapsed, so we need to execute the cycles
		// necessary for the speed at which we're operating
		for e.TotalCycles < wantedCycles {
			cycles, err := comp.Process()
			if err != nil {
				slog.Error(fmt.Sprintf("process execution failed: %v", err))
				return
			}

			e.TotalCycles += int64(cycles)
		}
	}
}
