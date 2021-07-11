package a2

import "github.com/pevans/erc/pkg/data"

// A Switcher is a type which provides a way to handle soft switch reads and
// writes in a relatively generic way.
type Switcher interface {
	SwitchRead(c *Computer, addr data.DByte) data.Byte
	SwitchWrite(c *Computer, addr data.DByte, val data.Byte)
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

	rfn := func(s Switcher) func(*Computer, data.DByte) data.Byte {
		return func(c *Computer, addr data.DByte) data.Byte {
			return s.SwitchRead(c, addr)
		}
	}

	wfn := func(s Switcher) func(*Computer, data.DByte, data.Byte) {
		return func(c *Computer, addr data.DByte, val data.Byte) {
			s.SwitchWrite(c, addr, val)
		}
	}

	for _, a := range kbReadSwitches() {
		c.RMap[a.Int()] = rfn(&c.kb)
	}

	for _, a := range kbWriteSwitches() {
		c.WMap[a.Int()] = wfn(&c.kb)
	}

	for _, a := range memReadSwitches() {
		c.RMap[a.Int()] = rfn(&c.mem)
	}

	for _, a := range memWriteSwitches() {
		c.WMap[a.Int()] = wfn(&c.mem)
	}

	for _, addr := range bankReadSwitches() {
		c.RMap[addr.Int()] = rfn(&c.bank)
	}

	for _, addr := range bankWriteSwitches() {
		c.WMap[addr.Int()] = wfn(&c.bank)
	}

	for _, a := range pcReadSwitches() {
		c.RMap[a.Int()] = rfn(&c.pc)
	}

	for _, a := range pcWriteSwitches() {
		c.WMap[a.Int()] = func(c *Computer, a data.DByte, val data.Byte) {
			c.pc.SwitchWrite(c, a, val)
		}
	}

	for _, a := range displayReadSwitches() {
		c.RMap[a.Int()] = rfn(&c.disp)
	}

	for _, a := range displayWriteSwitches() {
		c.WMap[a.Int()] = wfn(&c.disp)
	}
}
