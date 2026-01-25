package a2

import (
	"time"

	"github.com/pevans/erc/a2/a2display"
	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/gfx"
	"github.com/pevans/erc/internal/metrics"
)

// Render will draw an updated picture of our graphics to the local
// framebuffer
func (c *Computer) Render() {
	if !c.State.Bool(a2state.DisplayRedraw) {
		return
	}

	metrics.Increment("renders", 1)

	// Snapshot display memory to prevent tearing during render. This copies
	// the current state so we render from a consistent view even if the CPU
	// modifies display memory mid-render. Use CopyFromState to respect page
	// switching and 80STORE settings.
	c.displaySnapshot.CopyFromState(c.Main, c.Aux, c.State)

	a2display.Render(c.displaySnapshot, c.Font40, c.State)

	// Handle screen capture logging for debugging
	if c.State.Bool(a2state.DisplayHires) {
		if c.screenLog != nil && time.Since(c.lastScreenCapture) >= time.Second {
			elapsed := time.Since(c.BootTime).Seconds()
			c.screenLog.CaptureFrame(gfx.Screen, elapsed)
			c.lastScreenCapture = time.Now()
		}
	}

	c.State.SetBool(a2state.DisplayRedraw, false)
}

// IsVerticalBlank returns true when the number of cycles we've emulated is
// during what the Apple would consider the period of "vertical blank" (when
// the screen would not have been drawn). It took the CRT gun 12,480 cycles to
// go from the top-left of the screen to the bottom-right, and 4,550 cycles to
// return to the top-left. Those 4,550 cycles are the vertical blank.
func (c *Computer) IsVerticalBlank() bool {
	return a2display.IsVerticalBlank(c.CPU.CycleCounter())
}
