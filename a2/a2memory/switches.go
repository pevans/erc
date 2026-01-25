package a2memory

import (
	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/internal/metrics"
	"github.com/pevans/erc/memory"
)

const (
	setMemReadMain  = int(0xC002)
	setMemReadAux   = int(0xC003)
	setMemWriteMain = int(0xC004)
	setMemWriteAux  = int(0xC005)
	rdMemReadAux    = int(0xC013)
	rdMemWriteAux   = int(0xC014)
)

// ReadSwitches returns the list of memory switch addresses that support
// reads.
func ReadSwitches() []int {
	return []int{
		rdMemReadAux,
		rdMemWriteAux,
	}
}

// WriteSwitches returns the list of memory switch addresses that support
// writes.
func WriteSwitches() []int {
	return []int{
		setMemReadMain,
		setMemWriteMain,
		setMemReadAux,
		setMemWriteAux,
	}
}

// SwitchRead handles reads from memory mode soft switches.
func SwitchRead(addr int, stm *memory.StateMap) uint8 {
	var (
		hi uint8 = 0x80
		lo uint8 = 0x00
	)

	switch addr {
	case rdMemReadAux:
		if stm.Bool(a2state.MemReadAux) {
			return hi
		}

	case rdMemWriteAux:
		if stm.Bool(a2state.MemWriteAux) {
			return hi
		}
	}

	return lo
}

// SwitchWrite handles writes to memory mode soft switches.
func SwitchWrite(addr int, val uint8, stm *memory.StateMap) {
	switch addr {
	case setMemReadAux:
		metrics.Increment("soft_memory_read_aux_on", 1)
		stm.SetBool(a2state.MemReadAux, true)
		stm.SetSegment(a2state.MemReadSegment, stm.Segment(a2state.MemAuxSegment))
	case setMemReadMain:
		metrics.Increment("soft_memory_read_aux_off", 1)
		stm.SetBool(a2state.MemReadAux, false)
		stm.SetSegment(a2state.MemReadSegment, stm.Segment(a2state.MemMainSegment))
	case setMemWriteAux:
		metrics.Increment("soft_memory_write_aux_on", 1)
		stm.SetBool(a2state.MemWriteAux, true)
		stm.SetSegment(a2state.MemWriteSegment, stm.Segment(a2state.MemAuxSegment))
	case setMemWriteMain:
		metrics.Increment("soft_memory_write_aux_off", 1)
		stm.SetBool(a2state.MemWriteAux, false)
		stm.SetSegment(a2state.MemWriteSegment, stm.Segment(a2state.MemMainSegment))
	}
}

// UseDefaults sets up the default state for memory modes.
func UseDefaults(stm *memory.StateMap, main, aux *memory.Segment) {
	stm.SetBool(a2state.MemReadAux, false)
	stm.SetBool(a2state.MemWriteAux, false)
	stm.SetSegment(a2state.MemReadSegment, main)
	stm.SetSegment(a2state.MemWriteSegment, main)
	stm.SetSegment(a2state.MemAuxSegment, aux)
	stm.SetSegment(a2state.MemMainSegment, main)
}
