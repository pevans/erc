package a2

import "github.com/pevans/erc/pkg/data"

// A Switcher is a type which provides a way to handle soft switch reads and
// writes in a relatively generic way.
type Switcher interface {
	SwitchRead(c *Computer, addr data.DByte) data.Byte
	SwitchWrite(c *Computer, addr data.DByte, val data.Byte)
}

var bankReadSwitches = []int{
	0xC011,
	0xC012,
	0xC016,
	0xC080,
	0xC081,
	0xC082,
	0xC083,
	0xC088,
	0xC089,
	0xC08A,
	0xC08B,
}

var bankWriteSwitches = []int{
	0xC008,
	0xC009,
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

	for _, addr := range []int{0xC013, 0xC014} {
		c.RMap[addr] = func(c *Computer, addr data.Addressor) data.Byte {
			return c.mem.SwitchRead(c, addr)
		}
	}

	for _, addr := range []int{0xC002, 0xC003, 0xC004, 0xC005} {
		c.WMap[addr] = func(c *Computer, addr data.Addressor, val data.Byte) {
			c.mem.SwitchWrite(c, addr, val)
		}
	}

	for _, addr := range bankReadSwitches {
		c.RMap[addr] = bankSwitchRead
	}

	for _, addr := range bankWriteSwitches {
		c.WMap[addr] = bankSwitchWrite
	}

	for _, addr := range []int{0xC015, 0xC017} {
		c.RMap[addr] = func(c *Computer, addr data.Addressor) data.Byte {
			return c.pc.SwitchRead(c, addr)
		}
	}

	for _, addr := range []int{0xC006, 0xC007, 0xC00A, 0xC00B} {
		c.WMap[addr] = func(c *Computer, addr data.Addressor, val data.Byte) {
			c.pc.SwitchWrite(c, addr, val)
		}
	}

	for _, a := range displayReadSwitches() {
		c.RMap[a.Addr()] = func(c *Computer, addr data.Addressor) data.Byte {
			return c.disp.SwitchRead(c, addr)
		}
	}

	for _, a := range displayWriteSwitches() {
		c.WMap[a.Addr()] = func(c *Computer, addr data.Addressor, val data.Byte) {
			c.disp.SwitchWrite(c, addr, val)
		}
	}
}
