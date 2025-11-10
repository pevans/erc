package memory

import (
	"fmt"
	"os"
)

// You could make an argument that this should be Endianness, and I would nod
// my head
type Endian int

const (
	BigEndian    Endian = 0 // bytes [0x11, 0x22] are represented as 16-bit $1122
	LittleEndian Endian = 1 // bytes [0x11, 0x22] are represented as 16-bit $2211
)

// A Segment is a block of memory divided into uint8s.
type Segment struct {
	Mem        []uint8
	smap       *SoftMap
	Endianness Endian
}

type SegmentReader interface {
	UseReadSegment(*Segment)
}

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
	s.Mem = make([]uint8, size)

	// Segments default to LittleEndian because that's how the Apple II works.
	s.Endianness = LittleEndian

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

// Take the bytes from a given segment, at some position and for some length,
// and pull that into our the receiver segment. Think of this as taking a
// chunk of the from segment and making that its own segment.
func (s *Segment) ExtractFrom(from *Segment, start, end int) (int, error) {
	if start < 0 || end > len(from.Mem) {
		return 0, fmt.Errorf("destination slice is out of bounds: %v, %v", start, end)
	}

	return s.CopySlice(0, from.Mem[start:end])
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

// Sets a 16-bit value with respect given to the endianness of the segment
func (s *Segment) Set16(addr int, val uint16) {
	if s.Endianness == LittleEndian {
		s.Set16LittleEndian(addr, val)
		return
	}

	s.Set16BigEndian(addr, val)
}

func (s *Segment) Set16BigEndian(addr int, val uint16) {
	lsb := uint8(val & 0xFF)
	msb := uint8(val >> 8)

	s.Set(addr+1, lsb)
	s.Set(addr, msb)
}

func (s *Segment) Set16LittleEndian(addr int, val uint16) {
	lsb := uint8(val & 0xFF)
	msb := uint8(val >> 8)

	s.Set(addr, lsb)
	s.Set(addr+1, msb)
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

// Gets a 16-bit value with respect given to the endianness of the segment
func (s *Segment) Get16(addr int) uint16 {
	if s.Endianness == LittleEndian {
		return s.Get16LittleEndian(addr)
	}

	return s.Get16BigEndian(addr)
}

func (s *Segment) Get16BigEndian(addr int) uint16 {
	lsb := s.Get(addr + 1)
	msb := s.Get(addr)

	return (uint16(msb) << 8) | uint16(lsb)
}

func (s *Segment) Get16LittleEndian(addr int) uint16 {
	lsb := s.Get(addr)
	msb := s.Get(addr + 1)

	return (uint16(msb) << 8) | uint16(lsb)
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

	return os.WriteFile(path, bytes, 0o644)
}

// ReadFile will read the contents of a given file into the segment
// receiver.
func (s *Segment) ReadFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	s.Mem = make([]uint8, len(data))
	for i, b := range data {
		s.Mem[i] = uint8(b)
	}

	return nil
}
