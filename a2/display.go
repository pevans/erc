package a2

import (
	"github.com/pevans/erc/internal/metrics"
	"github.com/pevans/erc/memory"
	"github.com/pevans/erc/statemap"
)

const (
	// These are R7 actions, meaning they are switches you read from that return
	// bit 7 high when the modes are on, and low if not.
	rd80Col   = int(0xC01F)
	rd80Store = int(0xC018)
	rdAltChar = int(0xC01E)
	rdDHires  = int(0xC07F)
	rdHires   = int(0xC01D)
	rdIOUDis  = int(0xC07E)
	rdMixed   = int(0xC01B)
	rdPage2   = int(0xC01C)
	rdText    = int(0xC01A)

	// These switches turn on modes
	on80Col   = int(0xC00D) // W
	on80Store = int(0xC001) // W
	onAltChar = int(0xC00F) // W
	onDHires  = int(0xC05F) // R/W
	onHires   = int(0xC057) // R/W
	onIOUDis  = int(0xC07F) // W
	onMixed   = int(0xC053) // R/W
	onPage2   = int(0xC055) // R/W
	onText    = int(0xC051) // R/W

	// And these switches turn them off.
	off80Col   = int(0xC00C) // W
	off80Store = int(0xC000) // W
	offAltChar = int(0xC00E) // W
	offDHires  = int(0xC05E) // R/W
	offHires   = int(0xC056) // R/W
	offIOUDis  = int(0xC07E) // W
	offMixed   = int(0xC052) // R/W
	offPage2   = int(0xC054) // R/W
	offText    = int(0xC050) // R/W
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
	c.state.SetBool(statemap.DisplayText, true)

	// All other options should be disabled
	c.state.SetBool(statemap.DisplayAltChar, false)
	c.state.SetBool(statemap.DisplayCol80, false)
	c.state.SetBool(statemap.DisplayDoubleHigh, false)
	c.state.SetBool(statemap.DisplayHires, false)
	c.state.SetBool(statemap.DisplayIou, false)
	c.state.SetBool(statemap.DisplayMixed, false)
	c.state.SetBool(statemap.DisplayPage2, false)
	c.state.SetBool(statemap.DisplayStore80, false)
	c.state.SetBool(statemap.DisplayRedraw, true)
	c.state.SetSegment(statemap.DisplayAuxSegment, c.Aux)
}

func displayOnOrOffReadWrite(a int, stm *memory.StateMap) bool {
	switch a {
	case onPage2:
		metrics.Increment("soft_display_page_2_on", 1)
		stm.SetBool(statemap.DisplayPage2, true)
		return true
	case onText:
		metrics.Increment("soft_display_text_on", 1)
		stm.SetBool(statemap.DisplayText, true)
		return true
	case onMixed:
		metrics.Increment("soft_display_mixed_on", 1)
		stm.SetBool(statemap.DisplayMixed, true)
		return true
	case onHires:
		metrics.Increment("soft_display_hires_on", 1)
		stm.SetBool(statemap.DisplayHires, true)
		return true
	case onDHires:
		metrics.Increment("soft_display_dhires_on", 1)
		if stm.Bool(statemap.DisplayIou) {
			stm.SetBool(statemap.DisplayDoubleHigh, true)
		}
		return true
	case offPage2:
		metrics.Increment("soft_display_page_2_off", 1)
		stm.SetBool(statemap.DisplayPage2, false)
		return true
	case offText:
		metrics.Increment("soft_display_text_off", 1)
		stm.SetBool(statemap.DisplayText, false)
		return true
	case offMixed:
		metrics.Increment("soft_display_mixed_off", 1)
		stm.SetBool(statemap.DisplayMixed, false)
		return true
	case offHires:
		metrics.Increment("soft_display_hires_off", 1)
		stm.SetBool(statemap.DisplayHires, false)
		return true
	case offDHires:
		metrics.Increment("soft_display_dhires_off", 1)
		if stm.Bool(statemap.DisplayIou) {
			stm.SetBool(statemap.DisplayDoubleHigh, false)
		}
		return true
	}

	return false
}

func displaySwitchRead(a int, stm *memory.StateMap) uint8 {
	var (
		hi uint8 = 0x80
		lo uint8 = 0x00
	)

	if displayOnOrOffReadWrite(a, stm) {
		return lo
	}

	switch a {
	case rdAltChar:
		if stm.Bool(statemap.DisplayAltChar) {
			return hi
		}
	case rd80Col:
		if stm.Bool(statemap.DisplayCol80) {
			return hi
		}
	case rd80Store:
		if stm.Bool(statemap.DisplayStore80) {
			return hi
		}
	case rdPage2:
		if stm.Bool(statemap.DisplayPage2) {
			return hi
		}
	case rdText:
		if stm.Bool(statemap.DisplayText) {
			return hi
		}
	case rdMixed:
		if stm.Bool(statemap.DisplayMixed) {
			return hi
		}
	case rdHires:
		if stm.Bool(statemap.DisplayHires) {
			return hi
		}
	case rdIOUDis:
		if stm.Bool(statemap.DisplayIou) {
			return hi
		}
	case rdDHires:
		if stm.Bool(statemap.DisplayDoubleHigh) {
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
		stm.SetBool(statemap.DisplayAltChar, true)
	case on80Col:
		stm.SetBool(statemap.DisplayCol80, true)
	case on80Store:
		stm.SetBool(statemap.DisplayStore80, true)
	case onIOUDis:
		stm.SetBool(statemap.DisplayIou, true)
	case offAltChar:
		stm.SetBool(statemap.DisplayAltChar, false)
	case off80Col:
		stm.SetBool(statemap.DisplayCol80, false)
	case off80Store:
		stm.SetBool(statemap.DisplayStore80, false)
	case offIOUDis:
		stm.SetBool(statemap.DisplayIou, false)
	}
}

func DisplaySegment(
	addr int,
	stm *memory.StateMap,
	segfunc func(*memory.StateMap) *memory.Segment,
) *memory.Segment {
	main := stm.Segment(statemap.MemMainSegment)
	displayAux := stm.Segment(statemap.DisplayAuxSegment)

	if stm.Bool(statemap.DisplayStore80) {
		if addr >= 0x0400 && addr < 0x0800 {
			if stm.Bool(statemap.DisplayPage2) {
				return displayAux
			}

			return main
		}

		if addr >= 0x2000 && addr < 0x4000 {
			if stm.Bool(statemap.DisplayHires) {
				if stm.Bool(statemap.DisplayPage2) {
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
	// Let the drawing routines we have know that it's time to re-render
	// the screen.
	stm.SetBool(statemap.DisplayRedraw, true)
	DisplaySegment(
		addr,
		stm,
		WriteSegment,
	).DirectSet(int(addr), val)
}

// Render will draw an updated picture of our graphics to the local framebuffer
func (c *Computer) Render() {
	if !c.state.Bool(statemap.DisplayRedraw) {
		return
	}

	metrics.Increment("renders", 1)

	// if it's text, do one thing
	// if it's lores, do another thing
	// if it's mixed, we need to do text + lores
	// if it's hires, do the hires thing
	// we also need to account for double text/lores/hires/mixed
	switch {
	case c.state.Bool(statemap.DisplayText):
		var (
			start int = 0x400
			end   int = 0x800
		)

		c.textRender(start, end)
	case c.state.Bool(statemap.DisplayHires):
		var (
			start int = 0x2000
			end   int = 0x4000
		)

		c.hiresRender(start, end)

	default:
		var (
			start int = 0x400
			end   int = 0x800
		)

		c.loresRender(start, end)
	}

	c.state.SetBool(statemap.DisplayRedraw, false)
}
