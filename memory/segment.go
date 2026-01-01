package memory

import (
	"fmt"
	"os"
)

// Endian is a type that represents the endianness of a segment.
type Endian int

const (
	BigEndian    Endian = 0 // bytes [0x11, 0x22] are represented as 16-bit $1122
	LittleEndian Endian = 1 // bytes [0x11, 0x22] are represented as 16-bit $2211
)

// A Segment is a block of memory divided into uint8s.
type Segment struct {
	// mem is the byte buffer that contains all of the data in the segment.
	mem []uint8

	// smap is a softmap that may contain side-effectful (is that a word?)
	// behavior that can occur when data is retrieved or set at a certain
	// address.
	smap *SoftMap

	// endianness represents the endianness of the data in the segment.
	endianness Endian
}

// A SegmentReader is some type that implements the UseReadSegment method.
type SegmentReader interface {
	UseReadSegment(*Segment)
}

// A SegmentWriter is some type that implements the UseWriteSegment method.
type SegmentWriter interface {
	UseWriteSegment(*Segment)
}

// A Getter can return a byte from a given address.
type Getter interface {
	Get(int) uint8
	Get16(int) uint16
}

// A Setter can set the value at a given address to a given byte.
type Setter interface {
	Set(int, uint8)
	Set16(int, uint16)
}

// NewSegment will return a new memory segment with for a given size.
func NewSegment(size int) *Segment {
	s := new(Segment)
	s.mem = make([]uint8, size)

	// Segments default to LittleEndian because that's how the Apple II works.
	s.endianness = LittleEndian

	return s
}

// UseSoftMap uses the provided softmap when handling Set and Get methods.
func (s *Segment) UseSoftMap(sm *SoftMap) {
	s.smap = sm
}

// Size returns the size of the given segment.
func (s *Segment) Size() int {
	return len(s.mem)
}

// CopySlice copies the contents of a slice of uint8s into a segment.
func (s *Segment) CopySlice(start int, bytes []uint8) (int, error) {
	toWrite := len(bytes)
	end := start + toWrite

	if start < 0 || end > len(s.mem) {
		return 0, fmt.Errorf("destination slice is out of bounds: %v, %v", start, end)
	}

	_ = copy(s.mem[start:end], bytes)

	return toWrite, nil
}

// ExtractFrom takes the bytes from a given segment, at some position and for
// some length, and pull that into our the receiver segment. Think of this as
// taking a chunk of the from segment and making that its own segment.
func (s *Segment) ExtractFrom(from *Segment, start, end int) (int, error) {
	if start < 0 || end > len(from.mem) {
		return 0, fmt.Errorf("destination slice is out of bounds: %v, %v", start, end)
	}

	return s.CopySlice(0, from.mem[start:end])
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

	s.mem[addr] = val
}

// Set16 sets a 16-bit value at the given address.
func (s *Segment) Set16(addr int, val uint16) {
	if s.endianness == LittleEndian {
		s.set16LittleEndian(addr, val)
		return
	}

	s.set16BigEndian(addr, val)
}

// set16BigEndian sets the value at the given address in big endian order.
func (s *Segment) set16BigEndian(addr int, val uint16) {
	lsb := uint8(val & 0xFF)
	msb := uint8(val >> 8)

	s.Set(addr+1, lsb)
	s.Set(addr, msb)
}

// set16LittleEndian sets the value at the given address in little endian order.
func (s *Segment) set16LittleEndian(addr int, val uint16) {
	lsb := uint8(val & 0xFF)
	msb := uint8(val >> 8)

	s.Set(addr, lsb)
	s.Set(addr+1, msb)
}

// DirectSet will bypass the registered SoftMap to directly set the provided
// value at a given address.
func (s *Segment) DirectSet(addr int, val uint8) {
	s.mem[addr] = val
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

	return s.mem[addr]
}

// Get16 gets a 16-bit value at some provided address.
func (s *Segment) Get16(addr int) uint16 {
	if s.endianness == LittleEndian {
		return s.get16LittleEndian(addr)
	}

	return s.get16BigEndian(addr)
}

// get16BigEndian returns the value of an address in big-endian order.
func (s *Segment) get16BigEndian(addr int) uint16 {
	lsb := s.Get(addr + 1)
	msb := s.Get(addr)

	return (uint16(msb) << 8) | uint16(lsb)
}

// get16LittleEndian returns the value of an address in little-endian order.
func (s *Segment) get16LittleEndian(addr int) uint16 {
	lsb := s.Get(addr)
	msb := s.Get(addr + 1)

	return (uint16(msb) << 8) | uint16(lsb)
}

// DirectGet, unlike Get, bypasses the segment's softmap and returns exactly
// the value at the given address.
func (s *Segment) DirectGet(addr int) uint8 {
	return s.mem[addr]
}

// ReadFile will read the contents of a given file into the segment
// receiver.
func (s *Segment) ReadFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	s.mem = make([]uint8, len(data))
	for i, b := range data {
		s.mem[i] = uint8(b)
	}

	return nil
}

// WriteFile writes the contents of this segment out to a file.
func (s *Segment) WriteFile(path string) error {
	bytes := make([]byte, len(s.mem))

	for i, b := range s.mem {
		bytes[i] = byte(b)
	}

	return os.WriteFile(path, bytes, 0o644)
}
