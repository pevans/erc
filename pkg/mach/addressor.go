package mach

// The Addressor interface allows you to abstract how to get an address
// (the way a slice might consider it) from another type.
type Addressor interface {
	Addr() int
}

// Addr produces an address for a slice, which is pretty simply done by
// casting to int.
func (b Byte) Addr() int {
	return int(b)
}

// Addr produces an address for a slice.
func (db DByte) Addr() int {
	return int(db)
}
