package a2

import (
	"time"

	"github.com/pevans/erc/a2/a2dhires"
	"github.com/pevans/erc/a2/a2hires"
	"github.com/pevans/erc/a2/a2lores"
	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/a2/a2text"
	"github.com/pevans/erc/gfx"
	"github.com/pevans/erc/internal/metrics"
	"github.com/pevans/erc/memory"
)

// scanCycleCount is the number of cycles it takes to, theoretically, redraw
// the screen.
const scanCycleCount uint64 = 17030

const (
	// These are R7 actions, meaning they are switches you read from that
	// return bit 7 high when the modes are on, and low if not.
	rd80Store = int(0xC018)
	rdVBL     = int(0xC019) // VBL = Vertical Blank
	rdText    = int(0xC01A)
	rdMixed   = int(0xC01B)
	rdPage2   = int(0xC01C)
	rdHires   = int(0xC01D)
	rdAltChar = int(0xC01E)
	rd80Col   = int(0xC01F)
	rdIOUDis  = int(0xC07E)
	rdDHires  = int(0xC07F)

	// Toggle switches
	off80Store = int(0xC000) // W
	on80Store  = int(0xC001) // W
	off80Col   = int(0xC00C) // W
	on80Col    = int(0xC00D) // W
	offAltChar = int(0xC00E) // W
	onAltChar  = int(0xC00F) // W
	offText    = int(0xC050) // R/W
	onText     = int(0xC051) // R/W
	offMixed   = int(0xC052) // R/W
	onMixed    = int(0xC053) // R/W
	offPage2   = int(0xC054) // R/W
	onPage2    = int(0xC055) // R/W
	offHires   = int(0xC056) // R/W
	onHires    = int(0xC057) // R/W
	onDHires   = int(0xC05E) // R/W -- interesting the order is reversed
	offDHires  = int(0xC05F) // R/W
	offIOUDis  = int(0xC07E) // W
	onIOUDis   = int(0xC07F) // W
)

func displayReadSwitches() []int {
	return []int{
		offDHires,
		offHires,
		offMixed,
		offPage2,
		offText,
		onDHires,
		onHires,
		onMixed,
		onPage2,
		onText,
		rd80Col,
		rd80Store,
		rdAltChar,
		rdDHires,
		rdHires,
		rdIOUDis,
		rdMixed,
		rdPage2,
		rdText,
		rdVBL,
	}
}

func displayWriteSwitches() []int {
	return []int{
		off80Col,
		off80Store,
		offAltChar,
		offDHires,
		offHires,
		offIOUDis,
		offMixed,
		offPage2,
		offText,
		on80Col,
		on80Store,
		onAltChar,
		onDHires,
		onHires,
		onIOUDis,
		onMixed,
		onPage2,
		onText,
	}
}

func displayUseDefaults(c *Computer) {
	// Text mode should be enabled
	c.State.SetBool(a2state.DisplayText, true)

	// All other options should be disabled
	c.State.SetBool(a2state.DisplayAltChar, false)
	c.State.SetBool(a2state.DisplayCol80, false)
	c.State.SetBool(a2state.DisplayDoubleHigh, false)
	c.State.SetBool(a2state.DisplayHires, false)
	c.State.SetBool(a2state.DisplayIou, false)
	c.State.SetBool(a2state.DisplayMixed, false)
	c.State.SetBool(a2state.DisplayPage2, false)
	c.State.SetBool(a2state.DisplayStore80, false)
	c.State.SetBool(a2state.DisplayRedraw, true)
	c.State.SetSegment(a2state.DisplayAuxSegment, c.Aux)
}

func displayOnOrOffReadWrite(a int, stm *memory.StateMap) bool {
	comp := stm.Any(a2state.DiskComputer).(*Computer)

	switch a {
	case onPage2:
		metrics.Increment("soft_display_page_2_on", 1)
		stm.SetBool(a2state.DisplayPage2, true)
		stm.SetBool(a2state.DisplayRedraw, true)
		return true
	case onText:
		metrics.Increment("soft_display_text_on", 1)
		stm.SetBool(a2state.DisplayText, true)
		stm.SetBool(a2state.DisplayRedraw, true)
		return true
	case onMixed:
		metrics.Increment("soft_display_mixed_on", 1)
		stm.SetBool(a2state.DisplayMixed, true)
		stm.SetBool(a2state.DisplayRedraw, true)
		return true
	case onHires:
		metrics.Increment("soft_display_hires_on", 1)
		stm.SetBool(a2state.DisplayHires, true)
		stm.SetBool(a2state.DisplayRedraw, true)
		return true
	case onDHires:
		metrics.Increment("soft_display_dhires_on", 1)
		stm.SetBool(a2state.DisplayDoubleHigh, true)
		gfx.Screen = comp.Screen
		stm.SetBool(a2state.DisplayRedraw, true)
		return true
	case offPage2:
		metrics.Increment("soft_display_page_2_off", 1)
		stm.SetBool(a2state.DisplayPage2, false)
		stm.SetBool(a2state.DisplayRedraw, true)
		return true
	case offText:
		metrics.Increment("soft_display_text_off", 1)
		stm.SetBool(a2state.DisplayText, false)
		stm.SetBool(a2state.DisplayRedraw, true)
		return true
	case offMixed:
		metrics.Increment("soft_display_mixed_off", 1)
		stm.SetBool(a2state.DisplayMixed, false)
		stm.SetBool(a2state.DisplayRedraw, true)
		return true
	case offHires:
		metrics.Increment("soft_display_hires_off", 1)
		stm.SetBool(a2state.DisplayHires, false)
		stm.SetBool(a2state.DisplayRedraw, true)
		return true
	case offDHires:
		metrics.Increment("soft_display_dhires_off", 1)
		stm.SetBool(a2state.DisplayDoubleHigh, false)
		stm.SetBool(a2state.DisplayRedraw, true)
		return true
	}

	return false
}

func displaySwitchRead(a int, stm *memory.StateMap) uint8 {
	var (
		hi   uint8 = 0x80
		lo   uint8 = 0x00
		comp       = stm.Any(a2state.DiskComputer).(*Computer)
	)

	if displayOnOrOffReadWrite(a, stm) {
		return lo
	}

	switch a {
	case rdAltChar:
		if stm.Bool(a2state.DisplayAltChar) {
			return hi
		}
	case rd80Col:
		if stm.Bool(a2state.DisplayCol80) {
			return hi
		}
	case rd80Store:
		if stm.Bool(a2state.DisplayStore80) {
			return hi
		}
	case rdPage2:
		if stm.Bool(a2state.DisplayPage2) {
			return hi
		}
	case rdText:
		if stm.Bool(a2state.DisplayText) {
			return hi
		}
	case rdMixed:
		if stm.Bool(a2state.DisplayMixed) {
			return hi
		}
	case rdHires:
		if stm.Bool(a2state.DisplayHires) {
			return hi
		}
	case rdIOUDis:
		if stm.Bool(a2state.DisplayIou) {
			return hi
		}
	case rdDHires:
		if stm.Bool(a2state.DisplayDoubleHigh) {
			return hi
		}
	case rdVBL:
		if comp.IsVerticalBlank() {
			return hi
		}
	}

	return lo
}

func displaySwitchWrite(a int, val uint8, stm *memory.StateMap) {
	if displayOnOrOffReadWrite(a, stm) {
		return
	}

	switch a {
	case onAltChar:
		metrics.Increment("soft_display_altchar_on", 1)
		stm.SetBool(a2state.DisplayAltChar, true)
		stm.SetBool(a2state.DisplayRedraw, true)
	case on80Col:
		metrics.Increment("soft_display_80col_on", 1)
		stm.SetBool(a2state.DisplayCol80, true)
		stm.SetBool(a2state.DisplayRedraw, true)
	case on80Store:
		metrics.Increment("soft_display_80store_on", 1)
		stm.SetBool(a2state.DisplayStore80, true)
		stm.SetBool(a2state.DisplayRedraw, true)
	case onIOUDis:
		metrics.Increment("soft_display_ioudis_on", 1)
		stm.SetBool(a2state.DisplayIou, true)
		stm.SetBool(a2state.DisplayRedraw, true)
	case offAltChar:
		metrics.Increment("soft_display_altchar_off", 1)
		stm.SetBool(a2state.DisplayAltChar, false)
		stm.SetBool(a2state.DisplayRedraw, true)
	case off80Col:
		metrics.Increment("soft_display_80col_off", 1)
		stm.SetBool(a2state.DisplayCol80, false)
		stm.SetBool(a2state.DisplayRedraw, true)
	case off80Store:
		metrics.Increment("soft_display_80store_off", 1)
		stm.SetBool(a2state.DisplayStore80, false)
		stm.SetBool(a2state.DisplayRedraw, true)
	case offIOUDis:
		metrics.Increment("soft_display_ioudis_off", 1)
		stm.SetBool(a2state.DisplayIou, false)
		stm.SetBool(a2state.DisplayRedraw, true)
	}
}

func DisplaySegment(
	addr int,
	stm *memory.StateMap,
	segfunc func(*memory.StateMap) *memory.Segment,
) *memory.Segment {
	main := stm.Segment(a2state.MemMainSegment)
	displayAux := stm.Segment(a2state.DisplayAuxSegment)

	if stm.Bool(a2state.DisplayStore80) {
		if addr >= 0x0400 && addr < 0x0800 {
			if stm.Bool(a2state.DisplayPage2) {
				return displayAux
			}

			return main
		}

		if addr >= 0x2000 && addr < 0x4000 {
			if stm.Bool(a2state.DisplayHires) {
				if stm.Bool(a2state.DisplayPage2) {
					return displayAux
				}

				return main
			}
		}
	}

	return segfunc(stm)
}

func DisplayRead(addr int, stm *memory.StateMap) uint8 {
	return DisplaySegment(
		addr,
		stm,
		ReadSegment,
	).DirectGet(int(addr))
}

func DisplayWrite(addr int, val uint8, stm *memory.StateMap) {
	// Let the drawing routines we have know that it's time to re-render the
	// screen.
	stm.SetBool(a2state.DisplayRedraw, true)

	// Track display memory writes for debugging
	if addr >= 0x0400 && addr < 0x0800 {
		metrics.Increment("display_write_text", 1)
	} else if addr >= 0x2000 && addr < 0x4000 {
		metrics.Increment("display_write_hires", 1)
	}

	DisplaySegment(
		addr,
		stm,
		WriteSegment,
	).DirectSet(int(addr), val)
}

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

	// - if it's text, do one thing
	// - if it's lores, do another thing
	// - if it's mixed, we need to do text + lores
	// - if it's hires, do the hires thing
	// we also need to account for double text/lores/hires/mixed
	switch {
	case c.State.Bool(a2state.DisplayText):
		var (
			start = 0x400
			end   = 0x800
		)

		monochromeMode := c.State.Int(a2state.DisplayMonochrome)

		a2text.Render(c.displaySnapshot, c.Font40, start, end, monochromeMode)

	case c.State.Bool(a2state.DisplayHires):
		monochromeMode := c.State.Int(a2state.DisplayMonochrome)

		if c.State.Bool(a2state.DisplayDoubleHigh) && c.State.Bool(a2state.DisplayCol80) {
			a2dhires.Render(c.displaySnapshot, monochromeMode)
		} else {
			var (
				start = 0x2000
				end   = 0x4000
			)

			a2hires.Render(c.displaySnapshot, start, end, monochromeMode)
		}

		if c.screenLog != nil && time.Since(c.lastScreenCapture) >= time.Second {
			elapsed := time.Since(c.BootTime).Seconds()
			c.screenLog.CaptureFrame(gfx.Screen, elapsed)
			c.lastScreenCapture = time.Now()
		}

	default:
		var (
			start = 0x400
			end   = 0x800
		)

		monochromeMode := c.State.Int(a2state.DisplayMonochrome)

		a2lores.Render(c.displaySnapshot, start, end, monochromeMode)
	}

	c.State.SetBool(a2state.DisplayRedraw, false)
}

// IsVerticalBlank returns true when the number of cycles we've emulated is
// during what the Apple would consider the period of "vertical blank" (when
// the screen would not have been drawn). It took the CRT gun 12,480 cycles to
// go from the top-left of the screen to the bottom-right, and 4,550 cycles to
// return to the top-left. Those 4,550 cycles are the vertical blank.
func (c *Computer) IsVerticalBlank() bool {
	cycles := c.CPU.CycleCounter() % scanCycleCount

	return cycles >= 12480
}
