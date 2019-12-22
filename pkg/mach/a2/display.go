package a2

import (
	"github.com/pevans/erc/pkg/data"
)

const (
	// DisplayDefault is the default state of the display settings.
	DisplayDefault = 0x0

	// DisplayAltCharset indicates that we should render characters in
	// Apple's alternate character set.
	DisplayAltCharset = 0x1

	// Display80Col directs us to show text in 80 columns, which is
	// double the normal width. The number of displayed rows is
	// unchanged.
	Display80Col = 0x2

	// DisplayText tells us to render the display buffer in text mode,
	// which means we should interpret the data there as text symbols
	// and not (for example) graphic cells.
	DisplayText = 0x4

	// DisplayMixed tells us to show both lores graphics and text. (It
	// is not possible to show hires graphics and text.) In this mode,
	// text is rendered at the bottom several rows; lores graphics,
	// above.
	DisplayMixed = 0x8

	// DisplayHires directs us to show high resolution graphics, rather
	// than low-resolution. The number of colors we can show decreases,
	// but the number of dots per inch increases.
	DisplayHires = 0x10

	// DisplayIOU enables IOU access for $C058 - $C05F.
	DisplayIOU = 0x20

	// DisplayDHires indicates that we will show double high-resolution
	// graphics. This mode requires the use of auxiliary memory.
	DisplayDHires = 0x40
)

func newDisplaySwitchCheck() *SwitchCheck {
	return &SwitchCheck{mode: displayMode, setMode: displaySetMode}
}

func displayMode(c *Computer) int {
	return c.DisplayMode
}

func displaySetMode(c *Computer, mode int) {
	c.DisplayMode = mode
}

func displayAuxSegment(c *Computer, addr data.DByte) *data.Segment {
	is80 := c.MemMode&Mem80Store > 0
	isHi := c.MemMode&MemHires > 0
	isP2 := c.MemMode&MemPage2 > 0

	if is80 {
		if addr >= 0x0400 && addr < 0x0800 && isHi {
			return c.Aux
		} else if addr >= 0x2000 && addr < 0x4000 && isHi && isP2 {
			return c.Aux
		}
	}

	return nil
}

func displayRead(c *Computer, addr data.Addressor) data.Byte {
	if seg := displayAuxSegment(c, data.DByte(addr.Addr())); seg != nil {
		return seg.Get(addr)
	}

	return c.ReadSegment().Get(addr)
}

func displayWrite(c *Computer, addr data.Addressor, val data.Byte) {
	// Let the drawing routines we have know that it's time to re-render
	// the screen.
	c.reDraw = true

	if seg := displayAuxSegment(c, data.DByte(addr.Addr())); seg != nil {
		seg.Set(addr, val)
		return
	}

	c.WriteSegment().Set(addr, val)
}
