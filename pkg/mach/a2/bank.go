package a2

import (
	"github.com/pevans/erc/pkg/mach"
)

func bankRead(c *Computer, addr mach.Addressor) mach.Byte {
	if ^c.BankMode&BankRAM > 0 {
		return c.ROM.Get(mach.Plus(addr, -SysRomOffset))
	}

	if addr.Addr() < 0xE000 && c.BankMode&BankRAM2 > 0 {
		return c.ReadSegment().Get(mach.Plus(addr, 0x3000))
	}

	return c.ReadSegment().Get(addr)
}

func bankWrite(c *Computer, addr mach.Addressor, val mach.Byte) {
	if ^c.BankMode&BankWrite > 0 {
		return
	}

	if addr.Addr() < 0xE000 && c.BankMode&BankRAM2 > 0 {
		c.WriteSegment().Set(mach.Plus(addr, 0x3000), val)
		return
	}

	c.WriteSegment().Set(addr, val)
}

func (c *Computer) bankSwitchSetR(flag int) ReadMapFn {
	return func(c *Computer, addr mach.Addressor) mach.Byte {
		// You'll note that this assigns to BankMode, rather than
		// OR'ing. That is intentional.
		c.BankMode = flag
		return 0x80
	}
}

func (c *Computer) bankSwitchIsSetR(flag int) ReadMapFn {
	return func(c *Computer, addr mach.Addressor) mach.Byte {
		if c.BankMode&flag > 0 {
			return 0x80
		}

		return 0x0
	}
}

func (c *Computer) bankSwitchUnsetW(flag int) WriteMapFn {
	return func(c *Computer, addr mach.Addressor, val mach.Byte) {
		c.BankMode &= ^flag
	}
}

func (c *Computer) bankSwitchSetW(flag int) WriteMapFn {
	return func(c *Computer, addr mach.Addressor, val mach.Byte) {
		c.BankMode |= flag
	}
}
