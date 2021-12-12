package a2

import (
	"github.com/pevans/erc/pkg/data"
)

type displaySwitcher struct{}

const (
	// Use the alternate character set if this is true (as opposed to the
	// primary set).
	displayAltChar = 500
	// Text display should show 80 columns, not the default 40
	displayCol80 = 501
	// If this is on, we will store page 2 data in aux memory.
	displayStore80 = 502
	// Page 2 will use that second page for graphics in some circumstances; in
	// others it might prefer page 1 but in auxiliary memory.
	displayPage2 = 503
	// Text controls whether we show text mode or not. This can be set in
	// addition to other modes, which is why this is treated as a bool and not a
	// const/enum for a resolution.
	displayText = 504
	// Controls whether we show a mix of low resolution and text mode in some situations.
	displayMixed = 505
	// If highRes is true, then we are in some form of high resolution mode;
	// otherwise we assume a low resolution mode.
	displayHires = 506
	// This enables "IOU" access for certain soft switches in the $C0 page. It
	// also enables double high resolution to be set.
	displayIou = 507
	// When this is true, then high resolution will be rendered as "double high"
	// resolution.
	displayDoubleHigh = 508

	displayRedraw     = 509
	displayAuxSegment = 510
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

func (ds *displaySwitcher) UseDefaults(c *Computer) {
	// Text mode should be enabled
	c.state.SetBool(displayText, true)

	// All other options should be disabled
	c.state.SetBool(displayAltChar, false)
	c.state.SetBool(displayCol80, false)
	c.state.SetBool(displayDoubleHigh, false)
	c.state.SetBool(displayHires, false)
	c.state.SetBool(displayIou, false)
	c.state.SetBool(displayMixed, false)
	c.state.SetBool(displayPage2, false)
	c.state.SetBool(displayStore80, false)
	c.state.SetBool(displayRedraw, true)
	c.state.SetSegment(displayAuxSegment, c.Aux)
}

func (ds *displaySwitcher) onOrOffReadWrite(c *Computer, a int) bool {
	switch a {
	case onPage2:
		c.state.SetBool(displayPage2, true)
		return true
	case onText:
		c.state.SetBool(displayText, true)
		return true
	case onMixed:
		c.state.SetBool(displayMixed, true)
		return true
	case onHires:
		c.state.SetBool(displayHires, true)
		return true
	case onDHires:
		if c.state.Bool(displayIou) {
			c.state.SetBool(displayDoubleHigh, true)
		}
		return true
	case offPage2:
		c.state.SetBool(displayPage2, false)
		return true
	case offText:
		c.state.SetBool(displayText, false)
		return true
	case offMixed:
		c.state.SetBool(displayMixed, false)
		return true
	case offHires:
		c.state.SetBool(displayHires, false)
		return true
	case offDHires:
		if c.state.Bool(displayIou) {
			c.state.SetBool(displayDoubleHigh, false)
		}
		return true
	}

	return false
}

func (ds *displaySwitcher) SwitchRead(c *Computer, a int) uint8 {
	var (
		hi uint8 = 0x80
		lo uint8 = 0x00
	)

	if ds.onOrOffReadWrite(c, a) {
		return lo
	}

	switch a {
	case rdAltChar:
		if c.state.Bool(displayAltChar) {
			return hi
		}
	case rd80Col:
		if c.state.Bool(displayCol80) {
			return hi
		}
	case rd80Store:
		if c.state.Bool(displayStore80) {
			return hi
		}
	case rdPage2:
		if c.state.Bool(displayPage2) {
			return hi
		}
	case rdText:
		if c.state.Bool(displayText) {
			return hi
		}
	case rdMixed:
		if c.state.Bool(displayMixed) {
			return hi
		}
	case rdHires:
		if c.state.Bool(displayHires) {
			return hi
		}
	case rdIOUDis:
		if c.state.Bool(displayIou) {
			return hi
		}
	case rdDHires:
		if c.state.Bool(displayDoubleHigh) {
			return hi
		}
	}

	return lo
}

func (ds *displaySwitcher) SwitchWrite(c *Computer, a int, val uint8) {
	if ds.onOrOffReadWrite(c, a) {
		return
	}

	switch a {
	case onAltChar:
		c.state.SetBool(displayAltChar, true)
	case on80Col:
		c.state.SetBool(displayCol80, true)
	case on80Store:
		c.state.SetBool(displayStore80, true)
	case onIOUDis:
		c.state.SetBool(displayIou, true)
	case offAltChar:
		c.state.SetBool(displayAltChar, false)
	case off80Col:
		c.state.SetBool(displayCol80, false)
	case off80Store:
		c.state.SetBool(displayStore80, false)
	case offIOUDis:
		c.state.SetBool(displayIou, false)
	}
}

func DisplaySegment(addr int, stm *data.StateMap) *data.Segment {
	if stm.Bool(displayStore80) {
		if addr >= 0x0400 && addr < 0x0800 && stm.Bool(displayHires) {
			return stm.Segment(displayAuxSegment)
		} else if addr >= 0x2000 && addr < 0x4000 &&
			stm.Bool(displayHires) &&
			stm.Bool(displayPage2) {
			return stm.Segment(displayAuxSegment)
		}
	}

	return ReadSegment(stm)
}

func DisplayRead(addr int, stm *data.StateMap) uint8 {
	return DisplaySegment(addr, stm).Get(int(addr))
}

func DisplayWrite(addr int, val uint8, stm *data.StateMap) {
	// Let the drawing routines we have know that it's time to re-render
	// the screen.
	stm.SetBool(displayRedraw, true)
	DisplaySegment(addr, stm).Set(int(addr), val)
}

// Render will draw an updated picture of our graphics to the local framebuffer
func (c *Computer) Render() {
	if !c.state.Bool(displayRedraw) {
		return
	}

	c.log.Debug("rendering...")

	// if it's text, do one thing
	// if it's lores, do another thing
	// if it's mixed, we need to do text + lores
	// if it's hires, do the hires thing
	// we also need to account for double text/lores/hires/mixed
	switch {
	case c.state.Bool(displayText):
		var (
			start int = 0x400
			end   int = 0x800
		)

		c.textRender(start, end)
	case c.state.Bool(displayHires):
		var (
			start int = 0x2000
			end   int = 0x4000
		)

		c.hiresRender(start, end)
	default:
		c.log.Debugf("i'm getting called with display mode %x", c.DisplayMode)
	}

	c.state.SetBool(displayRedraw, false)
}
