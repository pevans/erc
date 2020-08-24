package a2

import (
	"github.com/pevans/erc/pkg/data"
	"github.com/pkg/errors"
)

const (
	// BankDefault is the default bank-switching scheme: reads in
	// bs-memory go to ROM; writes to RAM are disallowed; bank 1 memory
	// is used.
	BankDefault = 0x00

	// BankRAM indicates that reads are from RAM rather than ROM.
	BankRAM = 0x01

	// BankWrite tells us that we can write to RAM in bs-memory.
	BankWrite = 0x02

	// BankRAM2 tells us to read from bank 2 memory for $D000..$DFFF.
	BankRAM2 = 0x04

	// BankAuxiliary indicates that we should reads and writes in the
	// zero page AND stack page will be done in auxiliary memory rather
	// than main memory. This flag ALSO indicates that reads and/or
	// writes to bs-memory are done in auxiliary memory.
	BankAuxiliary = 0x08
)

func newBankSwitchCheck() *SwitchCheck {
	return &SwitchCheck{mode: bankMode, setMode: bankSetMode}
}

func bankMode(c *Computer) int {
	return c.BankMode
}

func bankSetMode(c *Computer, mode int) {
	wasBank := c.BankMode&BankAuxiliary > 0
	nowBank := mode&BankAuxiliary > 0

	c.BankMode = mode

	var err error

	// We need to copy the zero page and stack to the Main segment from
	// the Aux segment, or vice versa
	switch {
	case wasBank && !nowBank:
		err = bankSyncPagesFromAux(c)
	case !wasBank && nowBank:
		err = bankSyncPagesToAux(c)
	}

	if err != nil {
		panic(errors.Wrap(err, "could not copy bank memory between segments"))
	}
}

func bankSyncPagesToAux(c *Computer) error {
	_, err := c.Aux.CopySlice(0, c.Main.Mem[0:0x200])
	return err
}

func bankSyncPagesFromAux(c *Computer) error {
	_, err := c.Main.CopySlice(0, c.Aux.Mem[0:0x200])
	return err
}

func bankRead(c *Computer, addr data.Addressor) data.Byte {
	if ^c.BankMode&BankRAM > 0 {
		return c.ROM.Get(data.Plus(addr, -SysRomOffset))
	}

	if addr.Addr() < 0xE000 && c.BankMode&BankRAM2 > 0 {
		return c.ReadSegment().Get(data.Plus(addr, 0x3000))
	}

	return c.ReadSegment().Get(addr)
}

func bankWrite(c *Computer, addr data.Addressor, val data.Byte) {
	if ^c.BankMode&BankWrite > 0 {
		return
	}

	if addr.Addr() < 0xE000 && c.BankMode&BankRAM2 > 0 {
		c.WriteSegment().Set(data.Plus(addr, 0x3000), val)
		return
	}

	c.WriteSegment().Set(addr, val)
}
