package a2

import "github.com/pevans/erc/pkg/data"

const (
	pcExpansion  = 300
	pcSlotC3     = 301
	pcSlotCX     = 302
	pcExpSlot    = 303
	pcROMSegment = 304
)

const (
	offExpROM    = int(0xCFFF)
	offSlotC3ROM = int(0xC00A)
	offSlotCXROM = int(0xC007)
	onSlotC3ROM  = int(0xC00B)
	onSlotCXROM  = int(0xC006)
	rdSlotC3ROM  = int(0xC017)
	rdSlotCXROM  = int(0xC015)
)

func pcReadSwitches() []int {
	return []int{
		offExpROM,
		rdSlotC3ROM,
		rdSlotCXROM,
	}
}

func pcWriteSwitches() []int {
	return []int{
		offSlotC3ROM,
		offSlotCXROM,
		onSlotC3ROM,
		onSlotCXROM,
	}
}

// UseDefaults sets the state of the pc switcher to that which it should have
// after a cold or warm boot.
func pcUseDefaults(c *Computer) {
	c.state.SetBool(pcExpansion, false)
	c.state.SetBool(pcSlotC3, false)
	c.state.SetBool(pcSlotCX, true)
	c.state.SetSegment(pcROMSegment, c.ROM)
}

// SwitchRead will return hi on bit 7 if slot c3 or cx is set to use peripheral
// rom; otherwise lo.
func pcSwitchRead(addr int, stm *data.StateMap) uint8 {
	var (
		hi uint8 = 0x80
		lo uint8 = 0x00
	)

	switch addr {
	case rdSlotC3ROM:
		if stm.Bool(pcSlotC3) {
			return hi
		}
		// it _seems_ like this should return lo instead of hi...?
	case rdSlotCXROM:
		if stm.Bool(pcSlotCX) {
			return lo
		}
	case offExpROM:
		// This is kind of an unusual switch, though, in that calling it
		// produces a side effect while returning from ROM.
		val := PCRead(addr, stm)

		// Hitting this address will clear the IO SELECT' and IO STROBE' signals
		// in the hardware, which essentially means that expansion rom is turned
		// off. But only after we get the return value.
		stm.SetBool(pcExpansion, false)
		stm.SetInt(pcExpSlot, 0)

		return val
	}

	if slotXROM(addr) {
		stm.SetInt(pcExpSlot, pcSlotFromAddr(addr))
		return PCRead(addr, stm)
	}

	if expROM(addr) && stm.Int(pcExpSlot) > 0 {
		stm.SetBool(pcExpansion, true)
		return PCRead(addr, stm)
	}

	return lo
}

func slotXROM(addr int) bool {
	return addr >= 0xC100 && addr < 0xC800
}

func expROM(addr int) bool {
	return addr >= 0xC800 && addr < 0xD000
}

func slot3ROM(addr int) bool {
	return addr >= 0xC300 && addr < 0xC400
}

// pcSlotFromAddr returns the effective slot number from a given CnXX address.
// While this can theoretically scale to any of sixteen slots, in practice the
// `n` will be between 1-7.
func pcSlotFromAddr(addr int) int {
	return (addr >> 8) & 0xf
}

// SwitchWrite will handle soft switch writes that, in our case, will enable or
// disable slot rom access.
func pcSwitchWrite(addr int, val uint8, stm *data.StateMap) {
	switch addr {
	case onSlotC3ROM:
		stm.SetBool(pcSlotC3, true)
	case offSlotC3ROM:
		stm.SetBool(pcSlotC3, false)
	case onSlotCXROM:
		// Note that enabling slotcx rom _also_ enables slotc3 rom, and
		// disabling does the same.
		stm.SetBool(pcSlotCX, true)
		stm.SetBool(pcSlotC3, true)
	case offSlotCXROM:
		// FIXME: the problem is that addresses aren't matching the
		// consts, even though they are equal values
		stm.SetBool(pcSlotCX, false)
		stm.SetBool(pcSlotC3, false)
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
func PCRead(addr int, stm *data.StateMap) uint8 {
	var (
		intROM    = int(pcIROMAddr(addr))
		periphROM = int(pcPROMAddr(addr))
	)

	switch {
	case
		stm.Bool(pcExpansion) && expROM(addr),
		stm.Bool(pcSlotC3) && slot3ROM(addr),
		stm.Bool(pcSlotCX) && slotXROM(addr):
		return stm.Segment(pcROMSegment).Get(periphROM)
	}

	return stm.Segment(pcROMSegment).Get(intROM)
}

// PCWrite is a stub which does nothing, since it handles writes into an
// explicitly read-only memory space.
func PCWrite(addr int, val uint8, stm *data.StateMap) {
	// Do nothing
}
