package a2display

import (
	"github.com/pevans/erc/a2/a2dhires"
	"github.com/pevans/erc/a2/a2hires"
	"github.com/pevans/erc/a2/a2lores"
	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/a2/a2text"
	"github.com/pevans/erc/gfx"
	"github.com/pevans/erc/internal/metrics"
	"github.com/pevans/erc/memory"
)

// ComputerState provides the minimal interface needed by display soft
// switches to access computer state.
type ComputerState interface {
	IsVerticalBlank() bool
	GetScreen() *gfx.FrameBuffer
}

// ScanCycleCount is the number of cycles it takes to, theoretically, redraw
// the screen.
const ScanCycleCount uint64 = 17030

const (
	// These are R7 actions, meaning they are switches you read from that
	// return bit 7 high when the modes are on, and low if not.
	Rd80Store = int(0xC018)
	RdVBL     = int(0xC019) // VBL = Vertical Blank
	RdText    = int(0xC01A)
	RdMixed   = int(0xC01B)
	RdPage2   = int(0xC01C)
	RdHires   = int(0xC01D)
	RdAltChar = int(0xC01E)
	Rd80Col   = int(0xC01F)
	RdIOUDis  = int(0xC07E)
	RdDHires  = int(0xC07F)

	// Toggle switches
	Off80Store = int(0xC000) // W
	On80Store  = int(0xC001) // W
	Off80Col   = int(0xC00C) // W
	On80Col    = int(0xC00D) // W
	OffAltChar = int(0xC00E) // W
	OnAltChar  = int(0xC00F) // W
	OffText    = int(0xC050) // R/W
	OnText     = int(0xC051) // R/W
	OffMixed   = int(0xC052) // R/W
	OnMixed    = int(0xC053) // R/W
	OffPage2   = int(0xC054) // R/W
	OnPage2    = int(0xC055) // R/W
	OffHires   = int(0xC056) // R/W
	OnHires    = int(0xC057) // R/W
	OnDHires   = int(0xC05E) // R/W -- interesting the order is reversed
	OffDHires  = int(0xC05F) // R/W
	OffIOUDis  = int(0xC07E) // W
	OnIOUDis   = int(0xC07F) // W
)

func ReadSwitches() []int {
	return []int{
		OffDHires,
		OffHires,
		OffMixed,
		OffPage2,
		OffText,
		OnDHires,
		OnHires,
		OnMixed,
		OnPage2,
		OnText,
		Rd80Col,
		Rd80Store,
		RdAltChar,
		RdDHires,
		RdHires,
		RdIOUDis,
		RdMixed,
		RdPage2,
		RdText,
		RdVBL,
	}
}

func WriteSwitches() []int {
	return []int{
		Off80Col,
		Off80Store,
		OffAltChar,
		OffDHires,
		OffHires,
		OffIOUDis,
		OffMixed,
		OffPage2,
		OffText,
		On80Col,
		On80Store,
		OnAltChar,
		OnDHires,
		OnHires,
		OnIOUDis,
		OnMixed,
		OnPage2,
		OnText,
	}
}

// IsVerticalBlank returns true when the cycle count is during what the Apple
// would consider the period of "vertical blank" (when the screen would not
// have been drawn). It took the CRT gun 12,480 cycles to go from the top-left
// of the screen to the bottom-right, and 4,550 cycles to return to the
// top-left. Those 4,550 cycles are the vertical blank.
func IsVerticalBlank(cycleCounter uint64) bool {
	cycles := cycleCounter % ScanCycleCount

	return cycles >= 12480
}

// UseDefaults sets the display state to default values (text mode enabled,
// all other modes disabled).
func UseDefaults(state *memory.StateMap, auxSegment *memory.Segment) {
	state.SetBool(a2state.DisplayText, true)
	state.SetBool(a2state.DisplayAltChar, false)
	state.SetBool(a2state.DisplayCol80, false)
	state.SetBool(a2state.DisplayDoubleHigh, false)
	state.SetBool(a2state.DisplayHires, false)
	state.SetBool(a2state.DisplayIou, false)
	state.SetBool(a2state.DisplayMixed, false)
	state.SetBool(a2state.DisplayPage2, false)
	state.SetBool(a2state.DisplayStore80, false)
	state.SetBool(a2state.DisplayRedraw, true)
	state.SetSegment(a2state.DisplayAuxSegment, auxSegment)
}

// onOrOffReadWrite handles soft switches that can be both read and written to
// toggle display modes on/off.
func onOrOffReadWrite(addr int, stm *memory.StateMap) bool {
	switch addr {
	case OnPage2:
		metrics.Increment("soft_display_page_2_on", 1)
		stm.SetBool(a2state.DisplayPage2, true)
		stm.SetBool(a2state.DisplayRedraw, true)
		return true
	case OnText:
		metrics.Increment("soft_display_text_on", 1)
		stm.SetBool(a2state.DisplayText, true)
		stm.SetBool(a2state.DisplayRedraw, true)
		return true
	case OnMixed:
		metrics.Increment("soft_display_mixed_on", 1)
		stm.SetBool(a2state.DisplayMixed, true)
		stm.SetBool(a2state.DisplayRedraw, true)
		return true
	case OnHires:
		metrics.Increment("soft_display_hires_on", 1)
		stm.SetBool(a2state.DisplayHires, true)
		stm.SetBool(a2state.DisplayRedraw, true)
		return true
	case OnDHires:
		metrics.Increment("soft_display_dhires_on", 1)
		stm.SetBool(a2state.DisplayDoubleHigh, true)
		comp := stm.Any(a2state.Computer).(ComputerState)
		gfx.Screen = comp.GetScreen()
		stm.SetBool(a2state.DisplayRedraw, true)
		return true
	case OffPage2:
		metrics.Increment("soft_display_page_2_off", 1)
		stm.SetBool(a2state.DisplayPage2, false)
		stm.SetBool(a2state.DisplayRedraw, true)
		return true
	case OffText:
		metrics.Increment("soft_display_text_off", 1)
		stm.SetBool(a2state.DisplayText, false)
		stm.SetBool(a2state.DisplayRedraw, true)
		return true
	case OffMixed:
		metrics.Increment("soft_display_mixed_off", 1)
		stm.SetBool(a2state.DisplayMixed, false)
		stm.SetBool(a2state.DisplayRedraw, true)
		return true
	case OffHires:
		metrics.Increment("soft_display_hires_off", 1)
		stm.SetBool(a2state.DisplayHires, false)
		stm.SetBool(a2state.DisplayRedraw, true)
		return true
	case OffDHires:
		metrics.Increment("soft_display_dhires_off", 1)
		stm.SetBool(a2state.DisplayDoubleHigh, false)
		stm.SetBool(a2state.DisplayRedraw, true)
		return true
	}

	return false
}

// SwitchRead handles reads from display soft switches.
func SwitchRead(addr int, stm *memory.StateMap) uint8 {
	var (
		hi uint8 = 0x80
		lo uint8 = 0x00
	)

	if onOrOffReadWrite(addr, stm) {
		return lo
	}

	comp := stm.Any(a2state.Computer).(ComputerState)

	switch addr {
	case RdAltChar:
		if stm.Bool(a2state.DisplayAltChar) {
			return hi
		}
	case Rd80Col:
		if stm.Bool(a2state.DisplayCol80) {
			return hi
		}
	case Rd80Store:
		if stm.Bool(a2state.DisplayStore80) {
			return hi
		}
	case RdPage2:
		if stm.Bool(a2state.DisplayPage2) {
			return hi
		}
	case RdText:
		if stm.Bool(a2state.DisplayText) {
			return hi
		}
	case RdMixed:
		if stm.Bool(a2state.DisplayMixed) {
			return hi
		}
	case RdHires:
		if stm.Bool(a2state.DisplayHires) {
			return hi
		}
	case RdIOUDis:
		if stm.Bool(a2state.DisplayIou) {
			return hi
		}
	case RdDHires:
		if stm.Bool(a2state.DisplayDoubleHigh) {
			return hi
		}
	case RdVBL:
		if comp.IsVerticalBlank() {
			return hi
		}
	}

	return lo
}

// SwitchWrite handles writes to display soft switches.
func SwitchWrite(addr int, val uint8, stm *memory.StateMap) {
	if onOrOffReadWrite(addr, stm) {
		return
	}

	switch addr {
	case OnAltChar:
		metrics.Increment("soft_display_altchar_on", 1)
		stm.SetBool(a2state.DisplayAltChar, true)
		stm.SetBool(a2state.DisplayRedraw, true)
	case On80Col:
		metrics.Increment("soft_display_80col_on", 1)
		stm.SetBool(a2state.DisplayCol80, true)
		stm.SetBool(a2state.DisplayRedraw, true)
	case On80Store:
		metrics.Increment("soft_display_80store_on", 1)
		stm.SetBool(a2state.DisplayStore80, true)
		stm.SetBool(a2state.DisplayRedraw, true)
	case OnIOUDis:
		metrics.Increment("soft_display_ioudis_on", 1)
		stm.SetBool(a2state.DisplayIou, true)
		stm.SetBool(a2state.DisplayRedraw, true)
	case OffAltChar:
		metrics.Increment("soft_display_altchar_off", 1)
		stm.SetBool(a2state.DisplayAltChar, false)
		stm.SetBool(a2state.DisplayRedraw, true)
	case Off80Col:
		metrics.Increment("soft_display_80col_off", 1)
		stm.SetBool(a2state.DisplayCol80, false)
		stm.SetBool(a2state.DisplayRedraw, true)
	case Off80Store:
		metrics.Increment("soft_display_80store_off", 1)
		stm.SetBool(a2state.DisplayStore80, false)
		stm.SetBool(a2state.DisplayRedraw, true)
	case OffIOUDis:
		metrics.Increment("soft_display_ioudis_off", 1)
		stm.SetBool(a2state.DisplayIou, false)
		stm.SetBool(a2state.DisplayRedraw, true)
	}
}

// Segment returns the appropriate memory segment for the given address based
// on display state (80STORE, page switching, etc).
func Segment(
	addr int,
	stm *memory.StateMap,
	segmentKey int,
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

	return stm.Segment(segmentKey)
}

// Read reads a byte from display memory at the given address.
func Read(addr int, stm *memory.StateMap) uint8 {
	return Segment(addr, stm, a2state.MemReadSegment).DirectGet(int(addr))
}

// Write writes a byte to display memory at the given address.
func Write(addr int, val uint8, stm *memory.StateMap) {
	stm.SetBool(a2state.DisplayRedraw, true)

	// Track display memory writes for debugging
	if addr >= 0x0400 && addr < 0x0800 {
		metrics.Increment("display_write_text", 1)
	} else if addr >= 0x2000 && addr < 0x4000 {
		metrics.Increment("display_write_hires", 1)
	}

	Segment(addr, stm, a2state.MemWriteSegment).DirectSet(int(addr), val)
}

// Render draws the display based on the current state and mode.
func Render(
	snapshot *Snapshot,
	font40 *gfx.Font,
	state *memory.StateMap,
) {
	switch {
	case state.Bool(a2state.DisplayText):
		var (
			start = 0x400
			end   = 0x800
		)

		monochromeMode := state.Int(a2state.DisplayMonochrome)

		a2text.Render(snapshot, font40, start, end, monochromeMode)

	case state.Bool(a2state.DisplayHires):
		monochromeMode := state.Int(a2state.DisplayMonochrome)

		if state.Bool(a2state.DisplayDoubleHigh) && state.Bool(a2state.DisplayCol80) {
			a2dhires.Render(snapshot, monochromeMode)
		} else {
			var (
				start = 0x2000
				end   = 0x4000
			)

			a2hires.Render(snapshot, start, end, monochromeMode)
		}

	default:
		var (
			start = 0x400
			end   = 0x800
		)

		monochromeMode := state.Int(a2state.DisplayMonochrome)

		a2lores.Render(snapshot, start, end, monochromeMode)
	}
}
