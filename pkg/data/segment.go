package data

import (
	"fmt"
	"io/ioutil"
)

// A Segment is a block of memory divided into uint8s.
type Segment struct {
	Mem  []uint8
	smap *SoftMap
}

// A Getter can return a byte from a given address.
type Getter interface {
	Get(int) uint8
}

// A Setter can set the value at a given address to a given byte.
type Setter interface {
	Set(int, uint8)
}

// NewSegment will return a new memory segment with for a given size.
func NewSegment(size int) *Segment {
	s := new(Segment)
	s.Mem = make([]uint8, size)

	return s
}

func (s *Segment) UseSoftMap(sm *SoftMap) {
	s.smap = sm
}

// Size returns the size of the given segment.
func (s *Segment) Size() int {
	return len(s.Mem)
}

// CopySlice copies the contents of a slice of uint8s into a segment.
func (s *Segment) CopySlice(start int, bytes []uint8) (int, error) {
	toWrite := len(bytes)
	end := start + toWrite

	if start < 0 || end > len(s.Mem) {
		return 0, fmt.Errorf("destination slice is out of bounds: %v, %v", start, end)
	}

	_ = copy(s.Mem[start:end], bytes)

	return toWrite, nil
}

// Set will set the value at a given cell. If a write function is
// registered for this cell, then we will call that and exit.
func (s *Segment) Set(addr int, val uint8) {
	if s.smap != nil {
		ok := s.smap.Write(addr, val)
		if ok {
			return
		}
	}

	s.Mem[addr] = val
}

func (s *Segment) DirectSet(addr int, val uint8) {
	s.Mem[addr] = val
}

// Get will get the value from a given cell. If a read function is
// registered, we will return whatever that is; otherwise we will return
// the value directly.
func (s *Segment) Get(addr int) uint8 {
	if s.smap != nil {
		val, ok := s.smap.Read(addr)
		if ok {
			return val
		}
	}

	return s.Mem[addr]
}

func (s *Segment) DirectGet(addr int) uint8 {
	return s.Mem[addr]
}

// WriteFile writes the contents of this segment out to a file.
func (s *Segment) WriteFile(path string) error {
	bytes := make([]byte, len(s.Mem))

	for i, b := range s.Mem {
		bytes[i] = byte(b)
	}

	return ioutil.WriteFile(path, bytes, 0644)
}

// ReadFile will read the contents of a given file into the segment
// receiver.
func (s *Segment) ReadFile(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	s.Mem = make([]uint8, len(data))
	for i, b := range data {
		s.Mem[i] = uint8(b)
	}

	return nil
}
