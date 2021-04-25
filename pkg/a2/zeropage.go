package a2

import "github.com/pevans/erc/pkg/data"

func zeroPageRead(c *Computer, addr data.Addressor) data.Byte {
	seg := c.Main
	if c.bank.sysBlock == bankAux {
		seg = c.Aux
	}

	return seg.Get(addr)
}

func zeroPageWrite(c *Computer, addr data.Addressor, val data.Byte) {
	seg := c.Main
	if c.bank.sysBlock == bankAux {
		seg = c.Aux
	}

	seg.Set(addr, val)
}
