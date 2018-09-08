package mach

// SegmentReadFn is a function signature for read mapper functions.
type SegmentReadFn func(s *Segment, addr Addressor) Byte

// SegmentWriteFn is a function signature for write mapper functions.
type SegmentWriteFn func(s *Segment, addr Addressor, val Byte)

// A Cell is one byte within a memory segment. Each Cell can also have a
// read and a write function mapped to it. When a Cell is read from, we
// return the value indicated by the read function if it is there; if
// not, we return val directly. When a Cell is written to, we
// execute the write function; if there is not one, then we write
// directly to val.
type Cell struct {
	Val     Byte
	ReadFn  SegmentReadFn
	WriteFn SegmentWriteFn
}

// A Segment is a block of memory divided into Bytes. Each Byte is
// stored inside of a Cell.
type Segment struct {
	Mem []Cell
}

// NewSegment will return a new memory segment with for a given size.
func NewSegment(size int) *Segment {
	s := new(Segment)
	s.Mem = make([]Cell, size)

	return s
}

// Set will set the value at a given cell. If a write function is
// registered for this cell, then we will call that and exit.
func (s *Segment) Set(addr Addressor, val Byte) {
	c := &s.Mem[addr.Addr()]
	if c.WriteFn != nil {
		c.WriteFn(s, addr, val)
		return
	}

	c.Val = val
}

// Get will get the value from a given cell. If a read function is
// registered, we will return whatever that is; otherwise we will return
// the value directly.
func (s *Segment) Get(addr Addressor) Byte {
	c := &s.Mem[addr.Addr()]
	if c.ReadFn != nil {
		return c.ReadFn(s, addr)
	}

	return c.Val
}

// SetReadFunc will set the read function for a given address
func (s *Segment) SetReadFunc(addr Addressor, fn SegmentReadFn) error {
	s.Mem[addr.Addr()].ReadFn = fn
	return nil
}

// SetWriteFunc will set the write function for a given address
func (s *Segment) SetWriteFunc(addr Addressor, fn SegmentWriteFn) error {
	s.Mem[addr.Addr()].WriteFn = fn
	return nil
}
