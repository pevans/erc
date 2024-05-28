package a2

import (
	"github.com/pevans/erc/internal/metrics"
	"github.com/pevans/erc/memory"
	"github.com/pevans/erc/statemap"
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
	c.state.SetBool(statemap.PCExpansion, false)
	c.state.SetBool(statemap.PCSlotC3, false)
	c.state.SetBool(statemap.PCSlotCX, true)
	c.state.SetSegment(statemap.PCROMSegment, c.ROM)
}

// SwitchRead will return hi on bit 7 if slot c3 or cx is set to use peripheral
// rom; otherwise lo.
func pcSwitchRead(addr int, stm *memory.StateMap) uint8 {
	var (
		hi uint8 = 0x80
		lo uint8 = 0x00
	)

	switch addr {
	case rdSlotC3ROM:
		metrics.Increment("soft_pc_read_slot_c3_rom", 1)
		if stm.Bool(statemap.PCSlotC3) {
			return hi
		}
		// it _seems_ like this should return lo instead of hi...?

	case rdSlotCXROM:
		metrics.Increment("soft_pc_read_slot_cx_rom", 1)
		if stm.Bool(statemap.PCSlotCX) {
			return lo
		}

		return hi

	case offExpROM:
		metrics.Increment("soft_pc_exp_rom_off", 1)

		// This is kind of an unusual switch, though, in that calling it
		// produces a side effect while returning from ROM.
		val := PCRead(addr, stm)

		// Hitting this address will clear the IO SELECT' and IO STROBE' signals
		// in the hardware, which essentially means that expansion rom is turned
		// off. But only after we get the return value.
		stm.SetBool(statemap.PCExpansion, false)
		stm.SetInt(statemap.PCExpSlot, 0)

		return val
	}

	if slotXROM(addr) {
		metrics.Increment("soft_pc_xrom", 1)
		stm.SetInt(statemap.PCExpSlot, pcSlotFromAddr(addr))
		return PCRead(addr, stm)
	}

	if expROM(addr) && stm.Int(statemap.PCExpSlot) > 0 {
		metrics.Increment("soft_pc_exp_rom", 1)
		stm.SetBool(statemap.PCExpansion, true)
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
func pcSwitchWrite(addr int, val uint8, stm *memory.StateMap) {
	switch addr {
	case onSlotC3ROM:
		metrics.Increment("soft_pc_slot_c3_rom_on", 1)
		stm.SetBool(statemap.PCSlotC3, true)
	case offSlotC3ROM:
		metrics.Increment("soft_pc_slot_c3_rom_off", 1)
		stm.SetBool(statemap.PCSlotC3, false)
	case onSlotCXROM:
		metrics.Increment("soft_pc_slot_cx_c3_rom_on", 1)
		// Note that enabling slotcx rom _also_ enables slotc3 rom, and
		// disabling does the same.
		stm.SetBool(statemap.PCSlotCX, true)
		stm.SetBool(statemap.PCSlotC3, true)
	case offSlotCXROM:
		metrics.Increment("soft_pc_slot_cx_c3_rom_off", 1)
		// FIXME: the problem is that addresses aren't matching the
		// consts, even though they are equal values
		stm.SetBool(statemap.PCSlotCX, false)
		stm.SetBool(statemap.PCSlotC3, false)
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
func PCRead(addr int, stm *memory.StateMap) uint8 {
	var (
		intROM    = int(pcIROMAddr(addr))
		periphROM = int(pcPROMAddr(addr))
	)

	switch {
	case
		stm.Bool(statemap.PCExpansion) && expROM(addr),
		stm.Bool(statemap.PCSlotC3) && slot3ROM(addr),
		stm.Bool(statemap.PCSlotCX) && slotXROM(addr):
		metrics.Increment("soft_pc_get_periph_rom", 1)
		return stm.Segment(statemap.PCROMSegment).DirectGet(periphROM)
	}

	metrics.Increment("soft_pc_get_int_rom", 1)
	return stm.Segment(statemap.PCROMSegment).DirectGet(intROM)
}

// PCWrite is a stub which does nothing, since it handles writes into an
// explicitly read-only memory space.
func PCWrite(addr int, val uint8, stm *memory.StateMap) {
	metrics.Increment("soft_pc_failed_write", 1)
	// Do nothing
}
