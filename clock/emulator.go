package clock

import (
	"time"
)

// Emulator is a device which you can use to simulate a slower
// computer's clockspeed.
type Emulator struct {
	// The rate at which we will emulate any process.
	hertz int64

	// For any given rate of hertz, this is the time we would expect to
	// be consumed by a single cycle.
	TimePerCycle time.Duration

	// The time that our wait period began. Periods are roughly a
	// second.
	lastPeriod time.Time

	// The number of cycles we've recorded within a period
	periodCycles int64
}

// NewEmulator returns a new emulator based on some number of hertz
func NewEmulator(hz int64) *Emulator {
	emu := &Emulator{
		hertz:        hz,
		TimePerCycle: (1 * time.Second) / time.Duration(hz),
	}

	return emu
}

// WaitForCycles will wait for a period of time that would allow the
// emulator to slow down the current thread based on its hertz. Whatever
// "waiting" means is up to the caller. Note that this function is not
// perfect; I've observed situations where the cycle rate can overclock
// a bit. The point is merely to get close to the correct rate.
func (e *Emulator) WaitForCycles(cycles int64, waitFunc func(d time.Duration)) {
	// If we've never ran this before, we should start a new period
	if e.lastPeriod.IsZero() {
		e.lastPeriod = time.Now()
	}

	e.periodCycles += cycles

	// idealCycles is the estimated number of cycles that should
	// have ran within a period
	idealCycles := int64(time.Since(e.lastPeriod) / e.TimePerCycle)

	// If we've executed more cycles in the period than are ideal, then
	// wait for some length of time that seems necessary to catch up to
	// the ideal rate
	if e.periodCycles > idealCycles {
		waitFunc(
			time.Since(e.lastPeriod) -
				(time.Duration(idealCycles) * e.TimePerCycle),
		)
	}

	// Reset the wait period
	if time.Since(e.lastPeriod) >= 1*time.Second {
		e.periodCycles = 0
		e.lastPeriod = time.Now()
	}
}
