package a2

import (
	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/internal/metrics"
	"github.com/pevans/erc/memory"
)

const (
	offMemReadAux  = int(0xC002)
	offMemWriteAux = int(0xC004)
	onMemReadAux   = int(0xC003)
	onMemWriteAux  = int(0xC005)
	rdMemReadAux   = int(0xC013)
	rdMemWriteAux  = int(0xC014)
)

func memReadSwitches() []int {
	return []int{
		rdMemReadAux,
		rdMemWriteAux,
	}
}

func memWriteSwitches() []int {
	return []int{
		offMemReadAux,
		offMemWriteAux,
		onMemReadAux,
		onMemWriteAux,
	}
}

func memUseDefaults(c *Computer) {
	c.State.SetBool(a2state.MemReadAux, false)
	c.State.SetBool(a2state.MemWriteAux, false)
	c.State.SetSegment(a2state.MemReadSegment, c.Main)
	c.State.SetSegment(a2state.MemWriteSegment, c.Main)
	c.State.SetSegment(a2state.MemAuxSegment, c.Aux)
	c.State.SetSegment(a2state.MemMainSegment, c.Main)
}

func memSwitchRead(addr int, stm *memory.StateMap) uint8 {
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

func memSwitchWrite(addr int, val uint8, stm *memory.StateMap) {
	switch addr {
	case onMemReadAux:
		metrics.Increment("soft_memory_read_aux_on", 1)
		stm.SetBool(a2state.MemReadAux, true)
		stm.SetSegment(a2state.MemReadSegment, stm.Segment(a2state.MemAuxSegment))
	case offMemReadAux:
		metrics.Increment("soft_memory_read_aux_off", 1)
		stm.SetBool(a2state.MemReadAux, false)
		stm.SetSegment(a2state.MemReadSegment, stm.Segment(a2state.MemMainSegment))
	case onMemWriteAux:
		metrics.Increment("soft_memory_write_aux_on", 1)
		stm.SetBool(a2state.MemWriteAux, true)
		stm.SetSegment(a2state.MemWriteSegment, stm.Segment(a2state.MemAuxSegment))
	case offMemWriteAux:
		metrics.Increment("soft_memory_write_aux_off", 1)
		stm.SetBool(a2state.MemWriteAux, false)
		stm.SetSegment(a2state.MemWriteSegment, stm.Segment(a2state.MemMainSegment))
	}
}

// Get will return the byte at addr, or will execute a read switch if
// one is present at the given address.
func (c *Computer) Get(addr int) uint8 {
	return ReadSegment(c.State).Get(addr)
}

// Set will set the byte at addr to val, or will execute a write switch
// if one is present at the given address.
func (c *Computer) Set(addr int, val uint8) {
	WriteSegment(c.State).Set(addr, val)
}

// MapRange will, given a range of addresses (from..to), set the read
// and write map functions to those given.
func (c *Computer) MapRange(from, to int, rfn memory.SoftRead, wfn memory.SoftWrite) {
	for addr := from; addr < to; addr++ {
		c.smap.SetRead(addr, rfn)
		c.smap.SetWrite(addr, wfn)
	}
}

// ReadSegment returns the segment that should be used for general
// reads, according to our current memory mode.
func ReadSegment(stm *memory.StateMap) *memory.Segment {
	return stm.Segment(a2state.MemReadSegment)
}

// WriteSegment returns the segment that should be used for general
// writes, according to our current memory mode.
func WriteSegment(stm *memory.StateMap) *memory.Segment {
	return stm.Segment(a2state.MemWriteSegment)
}
