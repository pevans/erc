package a2

import (
	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/memory"
)

// DisplaySnapshot holds a point-in-time copy of display memory. This prevents
// tearing/flickering during rendering by ensuring the render functions see a
// consistent state even if the CPU modifies display memory mid-render.
type DisplaySnapshot struct {
	// text/lores region: 0x400-0x800
	textLores [0x400]byte

	// hires region: 0x2000-0x4000
	hires [0x2000]byte
}

// NewDisplaySnapshot creates a new empty display snapshot.
func NewDisplaySnapshot() *DisplaySnapshot {
	return &DisplaySnapshot{}
}

// CopyFrom copies display memory from the given segment into the snapshot.
func (s *DisplaySnapshot) CopyFrom(seg *memory.Segment) {
	for i := range 0x400 {
		s.textLores[i] = seg.DirectGet(0x400 + i)
	}

	for i := range 0x2000 {
		s.hires[i] = seg.DirectGet(0x2000 + i)
	}
}

// CopyFromState copies display memory into the snapshot, respecting the
// current display page and memory configuration (80STORE, page 1/2, etc).
func (s *DisplaySnapshot) CopyFromState(main, aux *memory.Segment, stm *memory.StateMap) {
	// Determine which segment and address range to use for text/lores
	textSeg := main
	if stm.Bool(a2state.DisplayStore80) && stm.Bool(a2state.DisplayPage2) {
		textSeg = aux
	}

	for i := range 0x400 {
		s.textLores[i] = textSeg.DirectGet(0x400 + i)
	}

	// Determine which segment and address range to use for hi-res
	hiresSeg := main
	hiresStart := 0x2000

	if stm.Bool(a2state.DisplayPage2) {
		if stm.Bool(a2state.DisplayStore80) {
			// Page 2 with 80STORE: use aux memory at 0x2000-0x3FFF
			hiresSeg = aux
		} else {
			// Page 2 without 80STORE: use main memory at 0x4000-0x5FFF
			hiresStart = 0x4000
		}
	}

	for i := range 0x2000 {
		s.hires[i] = hiresSeg.DirectGet(hiresStart + i)
	}
}

// Get returns the byte at the given address from the snapshot.
func (s *DisplaySnapshot) Get(addr int) uint8 {
	if addr >= 0x400 && addr < 0x800 {
		return s.textLores[addr-0x400]
	}

	if addr >= 0x2000 && addr < 0x4000 {
		return s.hires[addr-0x2000]
	}

	return 0
}

// Get16 returns a 16-bit value at the given address from the snapshot.
func (s *DisplaySnapshot) Get16(addr int) uint16 {
	lo := uint16(s.Get(addr))
	hi := uint16(s.Get(addr + 1))

	return (hi << 8) | lo
}
