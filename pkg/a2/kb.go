package a2

import "github.com/pevans/erc/pkg/data"

type kbSwitcher struct {
	lastKey data.Byte
	strobe  data.Byte
	keyDown data.Byte
}

const (
	kbDataAndStrobe data.DByte = 0xC000
	kbAnyKeyDown    data.DByte = 0xC010
)

func kbReadSwitches() []data.DByte {
	return []data.DByte{
		kbDataAndStrobe,
		kbAnyKeyDown,
	}
}

func kbWriteSwitches() []data.DByte {
	return []data.DByte{
		kbAnyKeyDown,
	}
}

func (ks *kbSwitcher) SwitchRead(c *Computer, addr data.DByte) data.Byte {
	switch addr {
	case kbDataAndStrobe:
		return ks.lastKey | ks.strobe
	case kbAnyKeyDown:
		ks.strobe = 0
		return ks.keyDown
	}

	// Nothing else can really come in here, but if something did...
	return 0
}

func (ks *kbSwitcher) SwitchWrite(c *Computer, addr data.DByte, val data.Byte) {
	if addr == kbAnyKeyDown {
		ks.strobe = 0
	}
}

func (ks *kbSwitcher) UseDefaults() {
	ks.keyDown = 0
	ks.lastKey = 0
	ks.strobe = 0
}

func kbSwitchRead(c *Computer, addr data.DByte) data.Byte {
	return c.kb.SwitchRead(c, addr)
}

func kbSwitchWrite(c *Computer, addr data.DByte, val data.Byte) {
	c.kb.SwitchWrite(c, addr, val)
}

func (c *Computer) PressKey(key data.Byte) {
	// There can only be 7-bit ASCII in an Apple II, so we explicitly
	// take off the high bit.
	c.kb.lastKey = key & 0x7F

	// We need to set the strobe bit, which (when returned) is always
	// with the high bit at 1.
	c.kb.strobe = 0x80

	// This flag (again with the high bit set to 1) is set _while_ a key
	// is pressed.
	c.kb.keyDown = 0x80
}

func (c *Computer) ClearKeys() {
	c.kb.keyDown = 0
}
