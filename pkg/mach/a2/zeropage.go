package a2

import "github.com/pevans/erc/pkg/mach"

func zeroPageRead(c *Computer, addr mach.Addressor) mach.Byte {
	seg := c.Main
	if c.BankMode&BankAuxiliary > 0 {
		seg = c.Aux
	}

	return seg.Get(addr)
}

func zeroPageWrite(c *Computer, addr mach.Addressor, val mach.Byte) {
	seg := c.Main
	if c.BankMode&BankAuxiliary > 0 {
		seg = c.Aux
	}

	seg.Set(addr, val)
}
