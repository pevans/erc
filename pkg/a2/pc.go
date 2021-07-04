package a2

import (
	"github.com/pevans/erc/pkg/data"
)

type pcSwitcher struct {
	expansion bool
	slotC3    bool
	slotCX    bool
	expSlot   int
}

const (
	offExpROM    = data.DByte(0xCFFF)
	offSlotC3ROM = data.DByte(0xC00A)
	offSlotCXROM = data.DByte(0xC007)
	onSlotC3ROM  = data.DByte(0xC00B)
	onSlotCXROM  = data.DByte(0xC006)
	rdSlotC3ROM  = data.DByte(0xC017)
	rdSlotCXROM  = data.DByte(0xC015)
)

func pcReadSwitches() []data.DByte {
	return []data.DByte{
		offExpROM,
		rdSlotC3ROM,
		rdSlotCXROM,
	}
}

func pcWriteSwitches() []data.DByte {
	return []data.DByte{
		offSlotC3ROM,
		offSlotCXROM,
		onSlotC3ROM,
		onSlotCXROM,
	}
}

// UseDefaults sets the state of the pc switcher to that which it should have
// after a cold or warm boot.
func (ps *pcSwitcher) UseDefaults() {
	ps.expansion = false
	ps.slotC3 = false
	ps.slotCX = true
}

// SwitchRead will return hi on bit 7 if slot c3 or cx is set to use peripheral
// rom; otherwise lo.
func (ps *pcSwitcher) SwitchRead(c *Computer, addr data.DByte) data.Byte {
	var (
		hi      data.Byte = 0x80
		lo      data.Byte = 0x00
		addrInt           = addr.Int()
	)

	switch addr {
	case rdSlotC3ROM:
		if ps.slotC3 {
			return hi
		}
		// it _seems_ like this should return lo instead of hi...?
	case rdSlotCXROM:
		if ps.slotCX {
			return lo
		}
	case offExpROM:
		// This is kind of an unusual switch, though, in that calling it
		// produces a side effect while returning from ROM.
		val := PCRead(c, addr)

		// Hitting this address will clear the IO SELECT' and IO STROBE' signals
		// in the hardware, which essentially means that expansion rom is turned
		// off. But only after we get the return value.
		ps.expansion = false
		ps.expSlot = 0

		return val
	}

	if ps.slotXROM(addrInt) {
		ps.expSlot = ps.slotFromAddr(addrInt)
		return PCRead(c, addr)
	}

	if ps.expROM(addrInt) && ps.expSlot > 0 {
		ps.expansion = true
		return PCRead(c, addr)
	}

	return lo
}

func (ps *pcSwitcher) slotXROM(addr int) bool {
	return addr >= 0xC100 && addr < 0xC800
}

func (ps *pcSwitcher) expROM(addr int) bool {
	return addr >= 0xC800 && addr < 0xD000
}

func (ps *pcSwitcher) slot3ROM(addr int) bool {
	return addr >= 0xC300 && addr < 0xC400
}

// slotFromAddr returns the effective slot number from a given CnXX address.
// While this can theoretically scale to any of sixteen slots, in practice the
// `n` will be between 1-7.
func (ps *pcSwitcher) slotFromAddr(addr int) int {
	return (addr >> 8) & 0xf
}

// SwitchWrite will handle soft switch writes that, in our case, will enable or
// disable slot rom access.
func (ps *pcSwitcher) SwitchWrite(c *Computer, addr data.DByte, val data.Byte) {
	switch addr {
	case onSlotC3ROM:
		ps.slotC3 = true
	case offSlotC3ROM:
		ps.slotC3 = false
	case onSlotCXROM:
		// Note that enabling slotcx rom _also_ enables slotc3 rom, and
		// disabling does the same.
		ps.slotCX = true
		ps.slotC3 = true
	case offSlotCXROM:
		// FIXME: the problem is that addresses aren't matching the
		// consts, even though they are equal values
		ps.slotCX = false
		ps.slotC3 = false
	}
}

func pcIROMAddr(addr int) int {
	return addr - 0xC000
}

func pcPROMAddr(addr int) int {
	return addr - 0xC000 + 0x4000
}

// PCRead returns a byte from ROM within the peripheral card address space
// ($C1..$CF). Based on the contents of the computer's PC Switcher, this can be
// from internal ROM or from a dedicated peripheral ROM block.
func PCRead(c *Computer, addr data.DByte) data.Byte {
	var (
		addrInt   = addr.Int()
		intROM    = data.DByte(pcIROMAddr(addrInt))
		periphROM = data.DByte(pcPROMAddr(addrInt))
	)

	switch {
	case
		c.pc.expansion && c.pc.expROM(addrInt),
		c.pc.slotC3 && c.pc.slot3ROM(addrInt),
		c.pc.slotCX && c.pc.slotXROM(addrInt):
		return c.ROM.Get(periphROM.Int())
	}

	return c.ROM.Get(intROM.Int())
}

// PCWrite is a stub which does nothing, since it handles writes into an
// explicitly read-only memory space.
func PCWrite(c *Computer, addr data.DByte, val data.Byte) {
	// Do nothing
}
