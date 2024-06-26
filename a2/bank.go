package a2

import (
	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/internal/metrics"
	"github.com/pevans/erc/memory"
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
		0xC081,
		0xC082,
		0xC083,
		0xC084,
		0xC085,
		0xC086,
		0xC087,
		0xC088,
		0xC089,
		0xC08A,
		0xC08B,
		0xC08C,
		0xC08D,
		0xC08E,
		0xC08F,
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
func bankSwitchRead(addr int, stm *memory.StateMap) uint8 {
	// In this set of addresses, it's possible that we might need to return a
	// value with bit 7 "checked" (which is to say, 1).
	switch addr {
	case rdBnk2:
		metrics.Increment("soft_read_bank_2", 1)
		return bankBit7(stm.Bool(a2state.BankDFBlockBank2))
	case rdLCRAM:
		metrics.Increment("soft_read_bank_lang_card", 1)
		return bankBit7(stm.Bool(a2state.BankReadRAM))
	case rdAltZP:
		metrics.Increment("soft_read_bank_alt_zp", 1)
		return bankBit7(stm.Bool(a2state.BankSysBlockAux))

	case 0xC080, 0xC084:
		stm.SetBool(a2state.BankReadRAM, true)
		stm.SetBool(a2state.BankWriteRAM, false)
		stm.SetBool(a2state.BankDFBlockBank2, true)

	case 0xC081, 0xC085:
		if stm.Int(a2state.BankReadAttempts) >= 1 && stm.Bool(a2state.InstructionReadOp) {
			stm.SetBool(a2state.BankWriteRAM, true)
		}
		stm.SetBool(a2state.BankReadRAM, false)
		stm.SetBool(a2state.BankDFBlockBank2, true)

	case 0xC082, 0xC086:
		stm.SetBool(a2state.BankReadRAM, false)
		stm.SetBool(a2state.BankWriteRAM, false)
		stm.SetBool(a2state.BankDFBlockBank2, true)

	case 0xC083, 0xC087:
		if stm.Int(a2state.BankReadAttempts) >= 1 && stm.Bool(a2state.InstructionReadOp) {
			stm.SetBool(a2state.BankWriteRAM, true)
		}
		stm.SetBool(a2state.BankReadRAM, true)
		stm.SetBool(a2state.BankDFBlockBank2, true)

	case 0xC088, 0xC08C:
		stm.SetBool(a2state.BankReadRAM, true)
		stm.SetBool(a2state.BankWriteRAM, false)
		stm.SetBool(a2state.BankDFBlockBank2, false)

	case 0xC089, 0xC08D:
		if stm.Int(a2state.BankReadAttempts) >= 1 && stm.Bool(a2state.InstructionReadOp) {
			stm.SetBool(a2state.BankWriteRAM, true)
		}
		stm.SetBool(a2state.BankReadRAM, false)
		stm.SetBool(a2state.BankDFBlockBank2, false)

	case 0xC08A, 0xC08E:
		stm.SetBool(a2state.BankReadRAM, false)
		stm.SetBool(a2state.BankWriteRAM, false)
		stm.SetBool(a2state.BankDFBlockBank2, false)

	case 0xC08B, 0xC08F:
		if stm.Int(a2state.BankReadAttempts) >= 1 && stm.Bool(a2state.InstructionReadOp) {
			stm.SetBool(a2state.BankWriteRAM, true)
		}
		stm.SetBool(a2state.BankReadRAM, true)
		stm.SetBool(a2state.BankDFBlockBank2, false)

	}

	return 0x00
}

// SwitchWrite manages writes on soft switches that may modify bank-switch
// state, specifically that to do with the usage of main vs. auxilliary memory.
func bankSwitchWrite(addr int, val uint8, stm *memory.StateMap) {
	switch addr {
	case offAltZP:
		metrics.Increment("soft_write_bank_alt_zp_off", 1)
		stm.SetBool(a2state.BankSysBlockAux, false)
		stm.SetSegment(a2state.BankSysBlockSegment, stm.Segment(a2state.MemMainSegment))
	case onAltZP:
		metrics.Increment("soft_write_bank_alt_zp_on", 1)
		stm.SetBool(a2state.BankSysBlockAux, true)
		stm.SetSegment(a2state.BankSysBlockSegment, stm.Segment(a2state.MemAuxSegment))
	}
}

// bankBit7 will return, given some boolean condition that has already been
// computed, either a value with bit 7 flagged on, or zero.
func bankBit7(cond bool) uint8 {
	metrics.Increment("soft_read_bank_bit_7", 1)
	if cond {
		return 0x80
	}

	return 0x00
}

// UseDefaults will set the state of the bank switcher to use what the computer
// would have if you cold- or warm-booted.
func bankUseDefaults(c *Computer) {
	// "When you turn power on or reset the Apple IIe, it initializes the bank
	// switches for reading the ROM and writing the RAM, using the second bank
	// of RAM."
	c.State.SetBool(a2state.BankReadRAM, false)
	c.State.SetBool(a2state.BankWriteRAM, true)
	c.State.SetBool(a2state.BankDFBlockBank2, false)
	c.State.SetBool(a2state.BankSysBlockAux, false)
	c.State.SetSegment(a2state.BankSysBlockSegment, c.Main)
	c.State.SetSegment(a2state.BankROMSegment, c.ROM)
}

// BankSegment returns the memory segment that should be used with respect to
// bank-switched auxiliary memory.
func BankSegment(stm *memory.StateMap) *memory.Segment {
	return stm.Segment(a2state.BankSysBlockSegment)
}

// BankDFRead implements logic for reads into the D0...FF pages of memory,
// taking into account the bank-switched states that the computer currently has.
func BankDFRead(addr int, stm *memory.StateMap) uint8 {
	metrics.Increment("soft_read_bank_df_block", 1)

	if !stm.Bool(a2state.BankReadRAM) {
		return stm.Segment(a2state.BankROMSegment).DirectGet(int(addr) - SysRomOffset)
	}

	if stm.Bool(a2state.BankDFBlockBank2) && addr < 0xE000 {
		return BankSegment(stm).DirectGet(int(addr) + 0x3000)
	}

	return BankSegment(stm).DirectGet(int(addr))
}

// BankDFWrite implements logic for writes into the D0...FF pages of memory,
// taking into account the bank-switched states that the computer currently has.
func BankDFWrite(addr int, val uint8, stm *memory.StateMap) {
	metrics.Increment("soft_write_bank_df_block", 1)
	if !stm.Bool(a2state.BankWriteRAM) {
		return
	}

	if stm.Bool(a2state.BankDFBlockBank2) && addr < 0xE000 {
		BankSegment(stm).DirectSet(int(addr)+0x3000, val)
		return
	}

	BankSegment(stm).DirectSet(int(addr), val)
}

func BankZPRead(addr int, stm *memory.StateMap) uint8 {
	metrics.Increment("soft_read_bank_zp", 1)
	return BankSegment(stm).DirectGet(int(addr))
}

func BankZPWrite(addr int, val uint8, stm *memory.StateMap) {
	metrics.Increment("soft_write_bank_zp", 1)
	BankSegment(stm).DirectSet(int(addr), val)
}
