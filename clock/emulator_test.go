package clock

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewEmulator(t *testing.T) {
	assert.NotNil(t, NewEmulator(1))
}

func TestWaitForCycles(t *testing.T) {
	var (
		emu = NewEmulator(12345)

		wantDiff   time.Duration = emu.timePerCycle * 8
		actualDiff time.Duration = 0
	)

	emu.WaitForCycles(8, func(d time.Duration) {
		actualDiff = d
	})

	assert.GreaterOrEqual(t, wantDiff, actualDiff)
}
