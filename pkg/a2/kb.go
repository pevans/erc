package a2

type kbSwitcher struct{}

const (
	kbLastKey = 100
	kbStrobe  = 101
	kbKeyDown = 102
)

const (
	kbDataAndStrobe int = 0xC000
	kbAnyKeyDown    int = 0xC010
)

func kbReadSwitches() []int {
	return []int{
		kbDataAndStrobe,
		kbAnyKeyDown,
	}
}

func kbWriteSwitches() []int {
	return []int{
		kbAnyKeyDown,
	}
}

func (ks *kbSwitcher) SwitchRead(c *Computer, addr int) uint8 {
	switch addr {
	case kbDataAndStrobe:
		return c.state.Uint8(kbLastKey) | c.state.Uint8(kbStrobe)
	case kbAnyKeyDown:
		c.state.SetUint8(kbStrobe, 0)
		return c.state.Uint8(kbKeyDown)
	}

	// Nothing else can really come in here, but if something did...
	return 0
}

func (ks *kbSwitcher) SwitchWrite(c *Computer, addr int, val uint8) {
	if addr == kbAnyKeyDown {
		c.state.SetUint8(kbStrobe, 0)
	}
}

func (ks *kbSwitcher) UseDefaults(c *Computer) {
	c.state.SetUint8(kbLastKey, 0)
	c.state.SetUint8(kbKeyDown, 0)
	c.state.SetUint8(kbStrobe, 0)
}

func kbSwitchRead(c *Computer, addr int) uint8 {
	return c.kb.SwitchRead(c, addr)
}

func kbSwitchWrite(c *Computer, addr int, val uint8) {
	c.kb.SwitchWrite(c, addr, val)
}

func (c *Computer) PressKey(key uint8) {
	// There can only be 7-bit ASCII in an Apple II, so we explicitly
	// take off the high bit.
	c.state.SetUint8(kbLastKey, key&0x7F)

	// We need to set the strobe bit, which (when returned) is always
	// with the high bit at 1.
	c.state.SetUint8(kbStrobe, 0x80)

	// This flag (again with the high bit set to 1) is set _while_ a key
	// is pressed.
	c.state.SetUint8(kbKeyDown, 0x80)
}

func (c *Computer) ClearKeys() {
	c.state.SetUint8(kbKeyDown, 0)
}
