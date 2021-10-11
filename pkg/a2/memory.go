package a2

import (
	"github.com/pevans/erc/pkg/data"
)

type memSwitcher struct{}

const (
	memRead  = 200
	memWrite = 201
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

func (ms *memSwitcher) UseDefaults(c *Computer) {
	c.state.SetInt(memRead, memMain)
	c.state.SetInt(memWrite, memMain)
}

func (ms *memSwitcher) SwitchRead(c *Computer, addr int) uint8 {
	var (
		hi uint8 = 0x80
		lo uint8 = 0x00
	)

	switch addr {
	case rdMemReadAux:
		if c.state.Int(memRead) == memAux {
			return hi
		}

	case rdMemWriteAux:
		if c.state.Int(memWrite) == memAux {
			return hi
		}
	}

	return lo
}

func (ms *memSwitcher) SwitchWrite(c *Computer, addr int, val uint8) {
	switch addr {
	case onMemReadAux:
		c.state.SetInt(memRead, memAux)
	case offMemReadAux:
		c.state.SetInt(memRead, memMain)
	case onMemWriteAux:
		c.state.SetInt(memWrite, memAux)
	case offMemWriteAux:
		c.state.SetInt(memWrite, memMain)
	}
}

// Get will return the byte at addr, or will execute a read switch if
// one is present at the given address.
func (c *Computer) Get(addr int) uint8 {
	uaddr := int(addr)
	if fn, ok := c.RMap[uaddr]; ok {
		return fn(c, uaddr)
	}

	return c.ReadSegment().Get(addr)
}

// Set will set the byte at addr to val, or will execute a write switch
// if one is present at the given address.
func (c *Computer) Set(addr int, val uint8) {
	uaddr := int(addr)
	if fn, ok := c.WMap[uaddr]; ok {
		fn(c, uaddr, val)
		return
	}

	c.WriteSegment().Set(addr, val)
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
	if c.state.Int(memRead) == memAux {
		return c.Aux
	}

	return c.Main
}

// WriteSegment returns the segment that should be used for general
// writes, according to our current memory mode.
func (c *Computer) WriteSegment() *data.Segment {
	if c.state.Int(memWrite) == memAux {
		return c.Aux
	}

	return c.Main
}
