package a2

import "github.com/pevans/erc/pkg/mach"

const (
	DisplayDefault = 0x0

	// Display text in the "alternate" character set
	DisplayAltCharset = 0x1

	// Show text in 80 columns, rather than the default 40 columns
	Display80Col = 0x2

	// Display only text. By default, we display lo-res graphics and
	// perhaps mixed graphics and text if the MIXED bit is high.
	DisplayText = 0x4

	// If TEXT is not high, then we are directed to display both text
	// and graphics.
	DisplayMixed = 0x8

	// If this is high, we will show high-resolution graphics; if not,
	// low-resolution. This bit is overridden by TEXT; if TEXT is high,
	// we will only show text.
	DisplayHires = 0x10

	// Enable IOU access for $C058..$C05F when this bit is on; NOTE: the
	// tech ref says that this is left on by the firmware
	DisplayIOU = 0x20

	// Display double-high-resolution graphics
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

func displayAuxSegment(c *Computer, addr mach.DByte) *mach.Segment {
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

func displayRead(c *Computer, addr mach.Addressor) mach.Byte {
	if seg := displayAuxSegment(c, mach.DByte(addr.Addr())); seg != nil {
		return seg.Get(addr)
	}

	return c.ReadSegment().Get(addr)
}

func displayWrite(c *Computer, addr mach.Addressor, val mach.Byte) {
	if seg := displayAuxSegment(c, mach.DByte(addr.Addr())); seg != nil {
		seg.Set(addr, val)
	}

	c.WriteSegment().Set(addr, val)
}
