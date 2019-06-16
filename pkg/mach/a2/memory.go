package a2

import (
	"github.com/pevans/erc/pkg/data"
)

const (
	// MemDefault tells us to read and write only to main memory.
	MemDefault = 0x00

	// MemReadAux will tell us to read the core first 48k memory from
	// auxiliary memory.
	MemReadAux = 0x01

	// MemWriteAux is the switch that tells us write to auxiliary memory
	// in the core 48k memory range.
	MemWriteAux = 0x02

	// Mem80Store is an "enabling" switch for MemPage2 and MemHires
	// below.  If this bit is not on, then those two other bits don't do
	// anything, and all aux memory access is governed by MemWriteAux
	// and MemReadAux above.
	Mem80Store = 0x04

	// MemPage2 allows access to auxiliary memory for the display page,
	// which is $0400..$07FF. This switch only works if Mem80Store is
	// also enabled.
	MemPage2 = 0x08

	// MemHires allows auxiliary memory access for $2000..$3FFF, as long
	// as MemPage2 and Mem80Store are also enabled.
	MemHires = 0x10
)

// Get will return the byte at addr, or will execute a read switch if
// one is present at the given address.
func (c *Computer) Get(addr data.Addressor) data.Byte {
	if fn, ok := c.RMap[addr.Addr()]; ok {
		return fn(c, addr)
	}

	if c.MemMode&MemReadAux > 0 {
		return c.Aux.Get(addr)
	}

	return c.Main.Get(addr)
}

// Set will set the byte at addr to val, or will execute a write switch
// if one is present at the given address.
func (c *Computer) Set(addr data.Addressor, val data.Byte) {
	if fn, ok := c.WMap[addr.Addr()]; ok {
		fn(c, addr, val)
		return
	}

	if c.MemMode&MemWriteAux > 0 {
		c.Aux.Set(addr, val)
		return
	}

	c.Main.Set(addr, val)
}

// MapRange will, given a range of addresses (from..to), set the read
// and write map functions to those given.
func (c *Computer) MapRange(from, to int, rfn ReadMapFn, wfn WriteMapFn) {
	for addr := from; addr < to; addr++ {
		c.RMap[addr] = rfn
		c.WMap[addr] = wfn
	}
}

// ReadSegment returns the segment that should be used for general
// reads, according to our current memory mode.
func (c *Computer) ReadSegment() *data.Segment {
	if c.MemMode&MemReadAux > 0 {
		return c.Aux
	}

	return c.Main
}

// WriteSegment returns the segment that should be used for general
// writes, according to our current memory mode.
func (c *Computer) WriteSegment() *data.Segment {
	if c.MemMode&MemWriteAux > 0 {
		return c.Aux
	}

	return c.Main
}

func newMemorySwitchCheck() *SwitchCheck {
	return &SwitchCheck{mode: memoryMode, setMode: memorySetMode}
}

func memoryMode(c *Computer) int {
	return c.MemMode
}

func memorySetMode(c *Computer, mode int) {
	c.MemMode = mode
}
