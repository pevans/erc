package data

// A Byte is simply an abstraction of what a byte would be to a generic
// machine.
type Byte uint8

// A DByte is simply a double-sized Byte.
type DByte uint16

func (b Byte) Int() int {
	return int(b)
}

func (db DByte) Int() int {
	return int(db)
}
