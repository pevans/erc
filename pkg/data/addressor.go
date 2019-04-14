package data

// The Addressor interface allows you to abstract how to get an address
// (the way a slice might consider it) from another type.
type Addressor interface {
	Addr() int
}

// Int is a simple type wrapper over int so that we can define methods
// for it.
type Int int

// Addr simply returns the underlying integer form of i.
func (i Int) Addr() int {
	return int(i)
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

// Plus returns an address that is n more than the a. You can, of
// course, have subtraction rather than addition if n is negative.
func Plus(a Addressor, n int) Addressor {
	addr := a.Addr()
	return Int(addr + n)
}
