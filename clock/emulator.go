package clock

import (
	"time"
)

// Emulator is a device which you can use to simulate a slower
// computer's clockspeed.
type Emulator struct {
	hertz        int64
	timePerCycle time.Duration
	lastWaitTime time.Time
}

// NewEmulator returns a new emulator based on some number of hertz
func NewEmulator(hz int64) *Emulator {
	// We intentionally leave lastWaitTime at the zero value
	emu := &Emulator{
		hertz:        hz,
		timePerCycle: (1 * time.Second) / time.Duration(hz),
	}

	return emu
}

// Override any past wait time with the given time. Useful when we start
// the clockspeed emulator later than when it was first initialized.
func (e *Emulator) SetWaitTime(t time.Time) {
	e.lastWaitTime = t
}

// WaitForCycles will wait for a period of time that would allow the
// emulator to slow down the current thread based on its hertz.
// Whatever "waiting" means is up to the caller.
func (e *Emulator) WaitForCycles(cycles int64, waitFunc func(d time.Duration)) {
	elapsedTime := time.Since(e.lastWaitTime)
	cycleTime := time.Duration(cycles) * e.timePerCycle

	waitFunc(cycleTime - elapsedTime)

	e.lastWaitTime = time.Now()
}
