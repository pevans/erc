package a2

import (
	"github.com/pevans/erc/pkg/data"
	"github.com/pkg/errors"
)

type bankSwitcher struct {
	// How many times have we tried to put write into bankRAM mode?
	writeAttempts int
}

const (
	bankRead            = 401
	bankWrite           = 402
	bankDFBlock         = 403
	bankSysBlock        = 404
	bankWriteAttempts   = 405
	bankSysBlockSegment = 406
	bankROMSegment      = 407
)

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
	offAltZP = int(0xC008)
	onAltZP  = int(0xC009)
	rdAltZP  = int(0xC016)
	rdBnk2   = int(0xC011)
	rdLCRAM  = int(0xC012)
)

func bankReadSwitches() []int {
	return []int{
		0xC080,
		0xC080,
		0xC081,
		0xC081,
		0xC082,
		0xC083,
		0xC083,
		0xC083,
		0xC088,
		0xC089,
		0xC08A,
		0xC08B,
		rdAltZP,
		rdBnk2,
		rdLCRAM,
	}
}

func bankWriteSwitches() []int {
	return []int{
		offAltZP,
		onAltZP,
	}
}

// SwitchRead manages reads from soft switches that mostly have to do with
// returning the state of bank-switching as well as, paradoxically, allowing
// callers to _modify_ said state.
func (bs *bankSwitcher) SwitchRead(c *Computer, addr int) uint8 {
	// In this set of addresses, it's possible that we might need to return a
	// value with bit 7 "checked" (which is to say, 1).
	switch addr {
	case rdBnk2:
		return bs.bit7(c.state.Int(bankDFBlock) == bank2)
	case rdLCRAM:
		return bs.bit7(c.state.Int(bankRead) == bankRAM)
	case rdAltZP:
		return bs.bit7(c.state.Int(bankSysBlock) == bankAux)
	}

	// Otherwise, we farm off the mode checks to other methods.
	c.state.SetInt(bankRead, bs.readMode(addr))
	c.state.SetInt(bankWrite, bs.writeMode(addr))
	c.state.SetInt(bankDFBlock, bs.dfBlockMode(addr))

	return 0x00
}

// SwitchWrite manages writes on soft switches that may modify bank-switch
// state, specifically that to do with the usage of main vs. auxilliary memory.
func (bs *bankSwitcher) SwitchWrite(c *Computer, addr int, val uint8) {
	origBlock := c.state.Int(bankSysBlock)

	switch addr {
	case offAltZP:
		c.state.SetInt(bankSysBlock, bankMain)
		c.state.SetSegment(bankSysBlockSegment, c.Main)
	case onAltZP:
		c.state.SetInt(bankSysBlock, bankAux)
		c.state.SetSegment(bankSysBlockSegment, c.Aux)
	}

	newBlock := c.state.Int(bankSysBlock)
	if origBlock != newBlock {
		if err := bankSyncPages(c, origBlock, newBlock); err != nil {
			panic(errors.Wrap(err, "could not copy bank memory between segments"))
		}
	}
}

// bit7 will return, given some boolean condition that has already been
// computed, either a value with bit 7 flagged on, or zero.
func (bs *bankSwitcher) bit7(cond bool) uint8 {
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
func (bs *bankSwitcher) UseDefaults(c *Computer) {
	// "When you turn power on or reset the Apple IIe, it initializes the bank
	// switches for reading the ROM and writing the RAM, using the second bank
	// of RAM."
	c.state.SetInt(bankRead, bankROM)
	c.state.SetInt(bankWrite, bankRAM)
	c.state.SetInt(bankDFBlock, bank2)
	c.state.SetInt(bankSysBlock, bankMain)
	c.state.SetSegment(bankSysBlockSegment, c.Main)
	c.state.SetSegment(bankROMSegment, c.ROM)
}

func bankSwitchRead(c *Computer, addr int) uint8 {
	return c.bank.SwitchRead(c, addr)
}

func bankSwitchWrite(c *Computer, addr int, val uint8) {
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
func BankSegment(stm *data.StateMap) *data.Segment {
	return stm.Segment(bankSysBlockSegment)
}

// BankDFRead implements logic for reads into the D0...FF pages of memory,
// taking into account the bank-switched states that the computer currently has.
func BankDFRead(addr int, stm *data.StateMap) uint8 {
	if stm.Int(bankDFBlock) == bank2 && addr < 0xE000 {
		return BankSegment(stm).DirectGet(int(addr) + 0x3000)
	}

	if stm.Int(bankRead) == bankROM {
		return stm.Segment(bankROMSegment).Get(int(addr) - SysRomOffset)
	}

	return BankSegment(stm).DirectGet(int(addr))
}

// BankDFWrite implements logic for writes into the D0...FF pages of memory,
// taking into account the bank-switched states that the computer currently has.
func BankDFWrite(addr int, val uint8, stm *data.StateMap) {
	if stm.Int(bankWrite) == bankNone {
		return
	}

	if stm.Int(bankDFBlock) == bank2 && addr < 0xE000 {
		BankSegment(stm).DirectSet(int(addr)+0x3000, val)
		return
	}

	BankSegment(stm).DirectSet(int(addr), val)
}

func BankZPRead(addr int, stm *data.StateMap) uint8 {
	return BankSegment(stm).DirectGet(int(addr))
}

func BankZPWrite(addr int, val uint8, stm *data.StateMap) {
	BankSegment(stm).DirectSet(int(addr), val)
}
