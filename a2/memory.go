package a2

import (
	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/memory"
)

// Get will return the byte at addr, or will execute a read switch if one is
// present at the given address.
func (c *Computer) Get(addr int) uint8 {
	return ReadSegment(c.State).Get(addr)
}

func (c *Computer) Get16(addr int) uint16 {
	return ReadSegment(c.State).Get16(addr)
}

// Set will set the byte at addr to val, or will execute a write switch if one
// is present at the given address.
func (c *Computer) Set(addr int, val uint8) {
	WriteSegment(c.State).Set(addr, val)
}

func (c *Computer) Set16(addr int, val uint16) {
	WriteSegment(c.State).Set16(addr, val)
}

// MapRange will, given a range of addresses (from..to), set the read and
// write map functions to those given.
func (c *Computer) MapRange(from, to int, rfn memory.SoftRead, wfn memory.SoftWrite) {
	for addr := from; addr < to; addr++ {
		c.smap.SetRead(addr, rfn)
		c.smap.SetWrite(addr, wfn)
	}
}

// ReadSegment returns the segment that should be used for general reads,
// according to our current memory mode.
func ReadSegment(stm *memory.StateMap) *memory.Segment {
	return stm.Segment(a2state.MemReadSegment)
}

// WriteSegment returns the segment that should be used for general writes,
// according to our current memory mode.
func WriteSegment(stm *memory.StateMap) *memory.Segment {
	return stm.Segment(a2state.MemWriteSegment)
}
