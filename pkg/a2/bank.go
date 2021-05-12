package a2

import (
	"github.com/pevans/erc/pkg/data"
	"github.com/pkg/errors"
)

type bankSwitcher struct {
	read     int
	write    int
	dfBlock  int
	sysBlock int

	// How many times have we tried to put write into bankRAM mode?
	writeAttempts int
}

// This const block defines some modes that our bank switcher can have.
const (
	bankRAM = iota
	bankROM
	bankNone
	bank1
	bank2
	bankMain
	bankAux
)

const (
	offAltZP = data.DByte(0xC008)
	onAltZP  = data.DByte(0xC009)
	rdAltZP  = data.DByte(0xC016)
	rdBnk2   = data.DByte(0xC011)
	rdLCRAM  = data.DByte(0xC012)
)

func bankReadSwitches() []data.Addressor {
	return []data.Addressor{
		data.DByte(0xC080),
		data.DByte(0xC080),
		data.DByte(0xC081),
		data.DByte(0xC081),
		data.DByte(0xC082),
		data.DByte(0xC083),
		data.DByte(0xC083),
		data.DByte(0xC083),
		data.DByte(0xC088),
		data.DByte(0xC089),
		data.DByte(0xC08A),
		data.DByte(0xC08B),
		rdAltZP,
		rdBnk2,
		rdLCRAM,
	}
}

func bankWriteSwitches() []data.Addressor {
	return []data.Addressor{
		offAltZP,
		onAltZP,
	}
}

// SwitchRead manages reads from soft switches that mostly have to do with
// returning the state of bank-switching as well as, paradoxically, allowing
// callers to _modify_ said state.
func (bs *bankSwitcher) SwitchRead(c *Computer, addr data.Addressor) data.Byte {
	// In this set of addresses, it's possible that we might need to return a
	// value with bit 7 "checked" (which is to say, 1).
	switch addr {
	case rdBnk2:
		return bs.bit7(bs.dfBlock == bank2)
	case rdLCRAM:
		return bs.bit7(bs.read == bankRAM)
	case rdAltZP:
		return bs.bit7(bs.sysBlock == bankAux)
	}

	// Otherwise, we farm off the mode checks to other methods.
	bs.read = bs.readMode(addr.Addr())
	bs.write = bs.writeMode(addr.Addr())
	bs.dfBlock = bs.dfBlockMode(addr.Addr())

	return 0x00
}

// SwitchWrite manages writes on soft switches that may modify bank-switch
// state, specifically that to do with the usage of main vs. auxilliary memory.
func (bs *bankSwitcher) SwitchWrite(c *Computer, addr data.Addressor, val data.Byte) {
	origBlock := bs.sysBlock

	switch addr {
	case offAltZP:
		bs.sysBlock = bankMain
	case onAltZP:
		bs.sysBlock = bankAux
	}

	if origBlock != bs.sysBlock {
		if err := bankSyncPages(c, origBlock, bs.sysBlock); err != nil {
			panic(errors.Wrap(err, "could not copy bank memory between segments"))
		}
	}
}

// bit7 will return, given some boolean condition that has already been
// computed, either a value with bit 7 flagged on, or zero.
func (bs *bankSwitcher) bit7(cond bool) data.Byte {
	if cond {
		return 0x80
	}

	return 0x00
}

func (bs *bankSwitcher) readMode(addr int) int {
	switch addr {
	case 0xC080, 0xC083, 0xC088, 0xC08B:
		return bankRAM
	}

	return bankROM
}

func (bs *bankSwitcher) writeMode(addr int) int {
	switch addr {
	case 0xC081, 0xC083, 0xC089, 0xC08B:
		bs.writeAttempts++
	}

	if bs.writeAttempts > 1 {
		bs.writeAttempts = 0
		return bankRAM
	}

	return bankNone
}

func (bs *bankSwitcher) dfBlockMode(addr int) int {
	switch addr {
	case 0xC080, 0xC081, 0xC082, 0xC083:
		return bank2
	}

	return bank1
}

// UseDefaults will set the state of the bank switcher to use what the computer
// would have if you cold- or warm-booted.
func (bs *bankSwitcher) UseDefaults() {
	// "When you turn power on or reset the Apple IIe, it initializes the bank
	// switches for reading the ROM and writing the RAM, using the second bank
	// of RAM."
	bs.read = bankROM
	bs.write = bankRAM
	bs.dfBlock = bank2
	bs.sysBlock = bankMain
}

func bankSwitchRead(c *Computer, addr data.Addressor) data.Byte {
	return c.bank.SwitchRead(c, addr)
}

func bankSwitchWrite(c *Computer, addr data.Addressor, val data.Byte) {
	c.bank.SwitchWrite(c, addr, val)
}

func bankSyncPages(c *Computer, oldmode, newmode int) error {
	if oldmode == newmode {
		return nil
	}

	if oldmode == bankMain {
		return bankSyncPagesToAux(c)
	}

	return bankSyncPagesFromAux(c)
}

func bankSyncPagesToAux(c *Computer) error {
	_, err := c.Aux.CopySlice(0, c.Main.Mem[0:0x200])
	return err
}

func bankSyncPagesFromAux(c *Computer) error {
	_, err := c.Main.CopySlice(0, c.Aux.Mem[0:0x200])
	return err
}

// BankSegment returns the memory segment that should be used with respect to
// bank-switched auxiliary memory.
func (c *Computer) BankSegment() *data.Segment {
	if c.bank.sysBlock == bankAux {
		return c.Aux
	}

	return c.Main
}

// BankDFRead implements logic for reads into the D0...FF pages of memory,
// taking into account the bank-switched states that the computer currently has.
func BankDFRead(c *Computer, addr data.Addressor) data.Byte {
	if c.bank.dfBlock == bank2 && addr.Addr() < 0xE000 {
		return c.BankSegment().Get(data.Plus(addr, 0x3000))
	}

	if c.bank.read == bankROM {
		return c.ROM.Get(data.Plus(addr, -SysRomOffset))
	}

	return c.BankSegment().Get(addr)
}

// BankDFWrite implements logic for writes into the D0...FF pages of memory,
// taking into account the bank-switched states that the computer currently has.
func BankDFWrite(c *Computer, addr data.Addressor, val data.Byte) {
	if c.bank.write == bankNone {
		return
	}

	if c.bank.dfBlock == bank2 && addr.Addr() < 0xE000 {
		c.BankSegment().Set(data.Plus(addr, 0x3000), val)
		return
	}

	c.BankSegment().Set(addr, val)
}

func BankZPRead(c *Computer, addr data.Addressor) data.Byte {
	return c.BankSegment().Get(addr)
}

func BankZPWrite(c *Computer, addr data.Addressor, val data.Byte) {
	c.BankSegment().Set(addr, val)
}
