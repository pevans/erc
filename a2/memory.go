package a2

import (
	"github.com/pevans/erc/memory"
)

const (
	memRead         = 200
	memWrite        = 201
	memReadSegment  = 202
	memWriteSegment = 203
	memAuxSegment   = 204
	memMainSegment  = 205
)

const (
	memMain = iota
	memAux
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
	c.state.SetInt(memRead, memMain)
	c.state.SetInt(memWrite, memMain)
	c.state.SetSegment(memReadSegment, c.Main)
	c.state.SetSegment(memWriteSegment, c.Main)
	c.state.SetSegment(memAuxSegment, c.Aux)
	c.state.SetSegment(memMainSegment, c.Main)
}

func memSwitchRead(addr int, stm *memory.StateMap) uint8 {
	var (
		hi uint8 = 0x80
		lo uint8 = 0x00
	)

	switch addr {
	case rdMemReadAux:
		if stm.Int(memRead) == memAux {
			return hi
		}

	case rdMemWriteAux:
		if stm.Int(memWrite) == memAux {
			return hi
		}
	}

	return lo
}

func memSwitchWrite(addr int, val uint8, stm *memory.StateMap) {
	switch addr {
	case onMemReadAux:
		stm.SetInt(memRead, memAux)
		stm.SetSegment(memReadSegment, stm.Segment(memAuxSegment))
	case offMemReadAux:
		stm.SetInt(memRead, memMain)
		stm.SetSegment(memReadSegment, stm.Segment(memMainSegment))
	case onMemWriteAux:
		stm.SetInt(memWrite, memAux)
		stm.SetSegment(memWriteSegment, stm.Segment(memAuxSegment))
	case offMemWriteAux:
		stm.SetInt(memWrite, memMain)
		stm.SetSegment(memWriteSegment, stm.Segment(memMainSegment))
	}
}

// Get will return the byte at addr, or will execute a read switch if
// one is present at the given address.
func (c *Computer) Get(addr int) uint8 {
	return ReadSegment(c.state).Get(addr)
}

// Set will set the byte at addr to val, or will execute a write switch
// if one is present at the given address.
func (c *Computer) Set(addr int, val uint8) {
	WriteSegment(c.state).Set(addr, val)
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
	return stm.Segment(memReadSegment)
}

// WriteSegment returns the segment that should be used for general
// writes, according to our current memory mode.
func WriteSegment(stm *memory.StateMap) *memory.Segment {
	return stm.Segment(memWriteSegment)
}
