package a2

// A Switcher is a type which provides a way to handle soft switch reads and
// writes in a relatively generic way.
type Switcher interface {
	SwitchRead(c *Computer, addr uint16) uint8
	SwitchWrite(c *Computer, addr uint16, val uint8)
}

// MapSoftSwitches will add several mappings for the soft switches that our
// computer uses.
func (c *Computer) MapSoftSwitches() {
	c.MapRange(0x0, 0x200, BankZPRead, BankZPWrite)
	c.MapRange(0x0400, 0x0800, DisplayRead, DisplayWrite)
	c.MapRange(0x2000, 0x4000, DisplayRead, DisplayWrite)
	// Note that there are other peripheral slots beginning with $C090, all the
	// way until $C100. We just don't emulate them right now.
	c.MapRange(0xC0E0, 0xC100, diskRead, diskWrite)
	c.MapRange(0xC100, 0xD000, PCRead, PCWrite)
	c.MapRange(0xD000, 0x10000, BankDFRead, BankDFWrite)

	rfn := func(s Switcher) func(*Computer, uint16) uint8 {
		return func(c *Computer, addr uint16) uint8 {
			return s.SwitchRead(c, addr)
		}
	}

	wfn := func(s Switcher) func(*Computer, uint16, uint8) {
		return func(c *Computer, addr uint16, val uint8) {
			s.SwitchWrite(c, addr, val)
		}
	}

	for _, a := range kbReadSwitches() {
		c.RMap[a] = rfn(&c.kb)
	}

	for _, a := range kbWriteSwitches() {
		c.WMap[a] = wfn(&c.kb)
	}

	for _, a := range memReadSwitches() {
		c.RMap[a] = rfn(&c.mem)
	}

	for _, a := range memWriteSwitches() {
		c.WMap[a] = wfn(&c.mem)
	}

	for _, addr := range bankReadSwitches() {
		c.RMap[addr] = rfn(&c.bank)
	}

	for _, addr := range bankWriteSwitches() {
		c.WMap[addr] = wfn(&c.bank)
	}

	for _, a := range pcReadSwitches() {
		c.RMap[a] = rfn(&c.pc)
	}

	for _, a := range pcWriteSwitches() {
		c.WMap[a] = func(c *Computer, a uint16, val uint8) {
			c.pc.SwitchWrite(c, a, val)
		}
	}

	for _, a := range displayReadSwitches() {
		c.RMap[a] = rfn(&c.disp)
	}

	for _, a := range displayWriteSwitches() {
		c.WMap[a] = wfn(&c.disp)
	}
}
