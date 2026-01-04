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
	// hertz is the rate at which we will emulate any process.
	hertz int64

	// timePerCycle, for any given rate of hertz, is the time we would expect
	// to be consumed by a single cycle.
	timePerCycle time.Duration

	// startTime is the time when the emulator first began
	startTime time.Time

	// resumeTime is the time when execution was last resumed (which may equal
	// StartTime, and which may be later than StartTime)
	resumeTime time.Time

	// totalCycles is the number of cycles we've executed since the start of
	// emulation
	totalCycles int64

	// fullSpeed is true when the emulator is not emulating its clockspeed,
	// but is instead moving as fast as it can.
	fullSpeed bool

	// debuggerEntryFunc is a function we'll run whenever we have gone into
	// debugging mode.
	debuggerEntryFunc func()

	// breakpointCheckFunc is a function we'll run whenever we need to test
	// that we've hit a breakpoint.
	breakpointCheckFunc func()
}

// NewEmulator returns a new emulator based on some number of hertz
func NewEmulator(hz int64) *Emulator {
	emu := &Emulator{
		hertz:        hz,
		timePerCycle: (1 * time.Second) / time.Duration(hz),
	}

	return emu
}

// ChangeHertz resets the clock emulation's hertz value and expected time
// per cycle.
func (e *Emulator) ChangeHertz(hz int64) {
	e.hertz = hz
	e.timePerCycle = (1 * time.Second) / time.Duration(hz)

	// Reset timing so we don't think we're behind or ahead on cycles
	e.resumeTime = time.Now()
	e.totalCycles = 0
}

// ProcessLoop runs the execution process for the provided computer. It does
// not exit unless there was some problem.
func (e *Emulator) ProcessLoop(comp emu.Computer) {
	e.startTime = time.Now()
	e.resumeTime = e.startTime
	state := comp.StateMap()

	wasPaused := false

	for {
		e.breakpointCheckFunc()

		if state.Bool(a2state.Debugger) {
			e.debuggerEntryFunc()

			// Reset ResumeTime so that we don't think we're far behind on
			// cycle time just because we sat in the debugger for a while
			e.resumeTime = time.Now()
			e.totalCycles = 0

			continue
		}

		if state.Bool(a2state.Paused) {
			wasPaused = true
			time.Sleep(100 * time.Millisecond)
			continue
		}

		// If we get to this point, and wasPaused is true, then that means the
		// Paused status was unset.
		if wasPaused {
			wasPaused = false
			e.resumeTime = time.Now()
			e.totalCycles = 0
		}

		elapsed := time.Since(e.resumeTime)
		wantedCycles := int64(elapsed / e.timePerCycle)

		// There are times when we ignore cycle timing and want to emulate as
		// fast as possible -- should that happen, it's up to the instructions
		// we're processing to unset FullSpeed.
		for e.totalCycles < wantedCycles || e.fullSpeed {
			wereFullSpeed := e.fullSpeed

			cycles, err := comp.Process()
			if err != nil {
				slog.Error(fmt.Sprintf("process execution failed: %v", err))
				return
			}

			e.totalCycles += int64(cycles)

			// When transitioning out of full speed, reset the timing so we
			// don't think we're behind and need to catch up
			if wereFullSpeed && !e.fullSpeed {
				e.resumeTime = time.Now()
				e.totalCycles = 0
				break
			}
		}
	}
}

// SetBreakpointCheck will use the provided function, f, to test if a
// breakpoint has been hit during execution.
func (e *Emulator) SetBreakpointCheck(f func()) {
	e.breakpointCheckFunc = f
}

// SetDebuggerEntry will use the provided function to enter a debugger (the
// caller must provide whatever logic would be used to debug things).
func (e *Emulator) SetDebuggerEntry(f func()) {
	e.debuggerEntryFunc = f
}

// SetFullSpeed will change the full-speed state of the emulator to the
// provided status, telling us either to go at full speed (true) or to emulate
// software at the normal hertz.
func (e *Emulator) SetFullSpeed(status bool) {
	e.fullSpeed = status
}

// TimePerCycle returns the duration of time that would be spent per cycle by
// the emulator
func (e *Emulator) TimePerCycle() time.Duration {
	return e.timePerCycle
}
