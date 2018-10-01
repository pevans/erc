package mach

import (
	"fmt"
)

// A Segment is a block of memory divided into Bytes.
type Segment struct {
	Mem []Byte
}

// A Getter can return a byte from a given address.
type Getter interface {
	Get(Addressor) Byte
}

// A Setter can set the value at a given address to a given byte.
type Setter interface {
	Set(Addressor, Byte)
}

// NewSegment will return a new memory segment with for a given size.
func NewSegment(size int) *Segment {
	s := new(Segment)
	s.Mem = make([]Byte, size)

	return s
}

// CopySlice copies the contents of a slice of Bytes into a segment.
func (s *Segment) CopySlice(start, end int, bytes []Byte) error {
	if start < 0 || end > len(s.Mem) {
		return fmt.Errorf("Destination slice is out of bounds: %v, %v", start, end)
	}

	_ = copy(s.Mem[start:end], bytes)

	return nil
}

// Set will set the value at a given cell. If a write function is
// registered for this cell, then we will call that and exit.
func (s *Segment) Set(addr Addressor, val Byte) {
	offset := addr.Addr()

	if offset < 0 || offset > cap(s.Mem) {
		panic(fmt.Sprintf("Memory access fault: address %v", offset))
	}

	s.Mem[offset] = val
}

// Get will get the value from a given cell. If a read function is
// registered, we will return whatever that is; otherwise we will return
// the value directly.
func (s *Segment) Get(addr Addressor) Byte {
	offset := addr.Addr()

	if offset < 0 || offset > cap(s.Mem) {
		panic(fmt.Sprintf("Memory access fault: address %v", offset))
	}

	return s.Mem[offset]
}
