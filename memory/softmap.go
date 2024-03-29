package memory

// SoftRead is a type of function that takes an integer address and
// returns an 8-bit value.
type SoftRead func(int, *StateMap) uint8

// SoftWrite is a type of function that takes an integer address and an
// 8-bit value and returns nothing.
type SoftWrite func(int, uint8, *StateMap)

// A SoftMap is a set of read and write function maps, keyed by address.
// It is intended to model soft switches in Apple IIs and similar
// constructs across other architectures.
type SoftMap struct {
	reads  []SoftRead
	writes []SoftWrite
	state  *StateMap
}

// NewSoftMap returns a newly allocated softmap with valid maps for
// reads and writes.
func NewSoftMap(size int) *SoftMap {
	sm := new(SoftMap)
	sm.reads = make([]SoftRead, size)
	sm.writes = make([]SoftWrite, size)
	return sm
}

func (sm *SoftMap) UseState(st *StateMap) {
	sm.state = st
}

// SetRead will assign a read function to a given address in the
// softmap.
func (sm *SoftMap) SetRead(addr int, fn SoftRead) {
	sm.reads[addr] = fn
}

// SetWrite will assign a write function to a given address in the
// softmap.
func (sm *SoftMap) SetWrite(addr int, fn SoftWrite) {
	sm.writes[addr] = fn
}

// Read executes a read call against the softmap. If no entry for an
// address exists, (0, false) is returned. Otherwise, the resulting
// value from the call and true.
func (sm *SoftMap) Read(addr int) (uint8, bool) {
	fn := sm.reads[addr]
	if fn == nil {
		return 0, false
	}

	return fn(addr, sm.state), true
}

// Write will execute a write call against the softmap; if no entry for
// an address exists, false is returned. Otherwise, true.
func (sm *SoftMap) Write(addr int, val uint8) bool {
	fn := sm.writes[addr]
	if fn == nil {
		return false
	}

	fn(addr, val, sm.state)
	return true
}
