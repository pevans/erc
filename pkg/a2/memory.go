package a2

import (
	"github.com/pevans/erc/pkg/data"
)

type memSwitcher struct {
	read  int
	write int
}

const (
	memMain = iota
	memAux
)

func (ms *memSwitcher) UseDefaults() {
	ms.read = memMain
	ms.write = memMain
}

func (ms *memSwitcher) SwitchRead(c *Computer, addr data.Addressor) data.Byte {
	var (
		hi data.Byte = 0x80
		lo data.Byte = 0x00
	)

	switch addr.Addr() {
	case 0xC013:
		if ms.read == memAux {
			return hi
		}

	case 0xC014:
		if ms.write == memAux {
			return hi
		}
	}

	return lo
}

func (ms *memSwitcher) SwitchWrite(c *Computer, addr data.Addressor, val data.Byte) {
	switch addr.Addr() {
	case 0xC003:
		ms.read = memAux
	case 0xC002:
		ms.read = memMain
	case 0xC005:
		ms.write = memAux
	case 0xC004:
		ms.write = memMain
	}
}

func memSwitchRead(c *Computer, addr data.Addressor) data.Byte {
	return c.mem.SwitchRead(c, addr)
}

func memSwitchWrite(c *Computer, addr data.Addressor, val data.Byte) {
	c.mem.SwitchWrite(c, addr, val)
}

// Get will return the byte at addr, or will execute a read switch if
// one is present at the given address.
func (c *Computer) Get(addr data.Addressor) data.Byte {
	if fn, ok := c.RMap[addr.Addr()]; ok {
		return fn(c, addr)
	}

	return c.ReadSegment().Get(addr)
}

// Set will set the byte at addr to val, or will execute a write switch
// if one is present at the given address.
func (c *Computer) Set(addr data.Addressor, val data.Byte) {
	if fn, ok := c.WMap[addr.Addr()]; ok {
		fn(c, addr, val)
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
	if c.mem.read == memAux {
		return c.Aux
	}

	return c.Main
}

// WriteSegment returns the segment that should be used for general
// writes, according to our current memory mode.
func (c *Computer) WriteSegment() *data.Segment {
	if c.mem.write == memAux {
		return c.Aux
	}

	return c.Main
}
