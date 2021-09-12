package a2

import (
	"github.com/pevans/erc/pkg/data"
)

type displaySwitcher struct {
	// Use the alternate character set if this is true (as opposed to the
	// primary set).
	altChar bool

	// Text display should show 80 columns, not the default 40
	col80 bool

	// If this is on, we will store page 2 data in aux memory.
	store80 bool

	// Page 2 will use that second page for graphics in some circumstances; in
	// others it might prefer page 1 but in auxiliary memory.
	page2 bool

	// Text controls whether we show text mode or not. This can be set in
	// addition to other modes, which is why this is treated as a bool and not a
	// const/enum for a resolution.
	text bool

	// Controls whether we show a mix of low resolution and text mode in some situations.
	mixed bool

	// If highRes is true, then we are in some form of high resolution mode;
	// otherwise we assume a low resolution mode.
	highRes bool

	// This enables "IOU" access for certain soft switches in the $C0 page. It
	// also enables double high resolution to be set.
	iou bool

	// When this is true, then high resolution will be rendered as "double high"
	// resolution.
	doubleHigh bool
}

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

func (ds *displaySwitcher) UseDefaults() {
	// Text mode should be enabled
	ds.text = true

	// All other options should be disabled
	ds.altChar = false
	ds.col80 = false
	ds.doubleHigh = false
	ds.highRes = false
	ds.iou = false
	ds.mixed = false
	ds.page2 = false
	ds.store80 = false
}

func (ds *displaySwitcher) onOrOffReadWrite(a int) bool {
	switch a {
	case onPage2:
		ds.page2 = true
		return true
	case onText:
		ds.text = true
		return true
	case onMixed:
		ds.mixed = true
		return true
	case onHires:
		ds.highRes = true
		return true
	case onDHires:
		if ds.iou {
			ds.doubleHigh = true
		}
		return true
	case offPage2:
		ds.page2 = false
		return true
	case offText:
		ds.text = false
		return true
	case offMixed:
		ds.mixed = false
		return true
	case offHires:
		ds.highRes = false
		return true
	case offDHires:
		if ds.iou {
			ds.doubleHigh = false
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

	if ds.onOrOffReadWrite(a) {
		return lo
	}

	switch a {
	case rdAltChar:
		if ds.altChar {
			return hi
		}
	case rd80Col:
		if ds.col80 {
			return hi
		}
	case rd80Store:
		if ds.store80 {
			return hi
		}
	case rdPage2:
		if ds.page2 {
			return hi
		}
	case rdText:
		if ds.text {
			return hi
		}
	case rdMixed:
		if ds.mixed {
			return hi
		}
	case rdHires:
		if ds.highRes {
			return hi
		}
	case rdIOUDis:
		if ds.iou {
			return hi
		}
	case rdDHires:
		if ds.doubleHigh {
			return hi
		}
	}

	return lo
}

func (ds *displaySwitcher) SwitchWrite(c *Computer, a int, val uint8) {
	if ds.onOrOffReadWrite(a) {
		return
	}

	switch a {
	case onAltChar:
		ds.altChar = true
	case on80Col:
		ds.col80 = true
	case on80Store:
		ds.store80 = true
	case onIOUDis:
		ds.iou = true
	case offAltChar:
		ds.altChar = false
	case off80Col:
		ds.col80 = false
	case off80Store:
		ds.store80 = false
	case offIOUDis:
		ds.iou = false
	}
}

func (c *Computer) DisplaySegment(addr int) *data.Segment {
	if c.disp.store80 {
		if addr >= 0x0400 && addr < 0x0800 && c.disp.highRes {
			return c.Aux
		} else if addr >= 0x2000 && addr < 0x4000 && c.disp.highRes && c.disp.page2 {
			return c.Aux
		}
	}

	return c.ReadSegment()
}

func DisplayRead(c *Computer, addr int) uint8 {
	return c.DisplaySegment(addr).Get(int(addr))
}

func DisplayWrite(c *Computer, addr int, val uint8) {
	// Let the drawing routines we have know that it's time to re-render
	// the screen.
	c.reDraw = true
	c.DisplaySegment(addr).Set(int(addr), val)
}

// Render will draw an updated picture of our graphics to the local framebuffer
func (c *Computer) Render() {
	if !c.reDraw {
		return
	}

	c.log.Debug("rendering...")

	// if it's text, do one thing
	// if it's lores, do another thing
	// if it's mixed, we need to do text + lores
	// if it's hires, do the hires thing
	// we also need to account for double text/lores/hires/mixed
	switch {
	case c.disp.text:
		var (
			start int = 0x400
			end   int = 0x800
		)

		c.textRender(start, end)
	case c.disp.highRes:
		var (
			start int = 0x2000
			end   int = 0x4000
		)

		c.hiresRender(start, end)
	default:
		c.log.Debugf("i'm getting called with display mode %x", c.DisplayMode)
	}

	c.reDraw = false
}
