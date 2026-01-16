package a2

import (
	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/internal/metrics"
	"github.com/pevans/erc/memory"
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
	c.State.SetBool(a2state.PCExpansion, false)
	c.State.SetBool(a2state.PCIOSelect, false)
	c.State.SetBool(a2state.PCIOStrobe, false)
	c.State.SetBool(a2state.PCSlotC3, false)
	c.State.SetBool(a2state.PCSlotCX, true)
	c.State.SetSegment(a2state.PCROMSegment, c.ROM)
}

// SwitchRead will return hi on bit 7 if slot c3 or cx is set to use
// peripheral rom; otherwise lo.
func pcSwitchRead(addr int, stm *memory.StateMap) uint8 {
	var (
		hi uint8 = 0x80
		lo uint8 = 0x00
	)

	switch addr {
	case rdSlotC3ROM:
		metrics.Increment("soft_pc_read_slot_c3_rom", 1)
		if stm.Bool(a2state.PCSlotC3) {
			return hi
		}
		// it _seems_ like this should return lo instead of hi...?

	case rdSlotCXROM:
		metrics.Increment("soft_pc_read_slot_cx_rom", 1)
		if stm.Bool(a2state.PCSlotCX) {
			return lo
		}

		return hi

	case offExpROM:
		metrics.Increment("soft_pc_exp_rom_off", 1)

		// This is kind of an unusual switch, though, in that calling it
		// produces a side effect while returning from ROM.
		val := PCRead(addr, stm)

		// Hitting this address will clear the IO SELECT' and IO STROBE'
		// signals in the hardware, which essentially means that expansion rom
		// is turned off. But only after we get the return value.
		stm.SetBool(a2state.PCExpansion, false)
		stm.SetInt(a2state.PCExpSlot, 0)

		return val
	}

	if slotXROM(addr) {
		metrics.Increment("soft_pc_xrom", 1)
		stm.SetInt(a2state.PCExpSlot, pcSlotFromAddr(addr))
		return PCRead(addr, stm)
	}

	if expROM(addr) && stm.Int(a2state.PCExpSlot) > 0 {
		metrics.Increment("soft_pc_exp_rom", 1)
		stm.SetBool(a2state.PCExpansion, true)
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

// SwitchWrite will handle soft switch writes that, in our case, will enable
// or disable slot rom access.
func pcSwitchWrite(addr int, val uint8, stm *memory.StateMap) {
	switch addr {
	case onSlotC3ROM:
		metrics.Increment("soft_pc_slot_c3_rom_on", 1)
		stm.SetBool(a2state.PCSlotC3, true)
	case offSlotC3ROM:
		metrics.Increment("soft_pc_slot_c3_rom_off", 1)
		stm.SetBool(a2state.PCSlotC3, false)
	case onSlotCXROM:
		metrics.Increment("soft_pc_slot_cx_c3_rom_on", 1)
		stm.SetBool(a2state.PCSlotCX, true)
	case offSlotCXROM:
		metrics.Increment("soft_pc_slot_cx_c3_rom_off", 1)
		stm.SetBool(a2state.PCSlotCX, false)
	}
}

func pcIROMAddr(addr int) int {
	return addr - 0xC000
}

func pcPROMAddr(addr int) int {
	return addr - 0xC000 + 0x4000
}

// PCRead returns a byte from ROM within the peripheral card address space
// ($C1..$CF). Based on the contents of the computer's PC Switcher, this can
// be from internal ROM or from a dedicated peripheral ROM block.
func PCRead(addr int, stm *memory.StateMap) uint8 {
	var (
		intAddr    = int(pcIROMAddr(addr))
		periphAddr = int(pcPROMAddr(addr))
		pcrom      = stm.Segment(a2state.PCROMSegment)
	)

	// Regardless of circumstances, any read of the expansion disable address
	// should wipe our expansion state.
	if addr == offExpROM {
		disableExpansion(stm)
	}

	switch {
	case stm.Bool(a2state.PCSlotCX):
		// Special case #1: we should turn on IOSelect if it's a slotCX
		// address.
		if slotXROM(addr) {
			stm.SetBool(a2state.PCIOSelect, true)
		}

		// Even though we want to return peripheral ROM for Cxxx addresses, if
		// SLOTC3ROM is not active, we should obey that and return internal
		// ROM.
		if !stm.Bool(a2state.PCSlotC3) && slot3ROM(addr) {
			return pcrom.DirectGet(intAddr)
		}

		// Special case #2: we should turn on IOStrobe for expansion ROM
		// addresses.
		if expROM(addr) {
			stm.SetBool(a2state.PCIOStrobe, true)

			if stm.Bool(a2state.PCIOSelect) {
				// If both IOSelect and IOStrobe are true, we have enabled
				// expansion ROM.
				stm.SetBool(a2state.PCExpansion, true)

				return expansionROM(stm, addr)
			}
		}

		return pcrom.DirectGet(periphAddr)

	case stm.Bool(a2state.PCSlotC3) && slot3ROM(addr):
		metrics.Increment("soft_pc_get_periph_rom", 1)
		return pcrom.DirectGet(periphAddr)

	case stm.Bool(a2state.PCExpansion) && expROM(addr):
		return expansionROM(stm, addr)
	}

	metrics.Increment("soft_pc_get_int_rom", 1)
	return pcrom.DirectGet(intAddr)
}

// PCWrite is a stub which does nothing, since it handles writes into an
// explicitly read-only memory space.
func PCWrite(addr int, val uint8, stm *memory.StateMap) {
	metrics.Increment("soft_pc_failed_write", 1)

	// Even a write to the expansion rom disable address should cause us to
	// wipe all of our state.
	if addr == offExpROM {
		disableExpansion(stm)
	}
}

func disableExpansion(stm *memory.StateMap) {
	stm.SetBool(a2state.PCIOSelect, false)
	stm.SetBool(a2state.PCIOStrobe, false)
	stm.SetBool(a2state.PCExpansion, false)
	stm.SetInt(a2state.PCExpSlot, 0)
}

func expansionROM(stm *memory.StateMap, addr int) uint8 {
	// Since we don't support any peripherals that have dedicated ROM, we must
	// fall back to returning data from internal ROM.
	return stm.Segment(a2state.PCROMSegment).DirectGet(
		pcIROMAddr(addr),
	)
}
