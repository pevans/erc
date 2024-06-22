package asm

import "fmt"

type Line struct {
	Address     int
	Instruction string
	Operand     string
	Comment     string
}

// String returns some representation of a line of assembly. There's no
// single grammar for assembly -- it's usually a notation that works for
// a specific assembler. As long as it "looks right", that's good enough
// for now.
func (ln Line) String() string {
	linefmt := "$%04X" + // address
		"%3s" + // spacing
		"%s " + // instruction
		"%-10s" + // operand
		"%5s" // spacing

	str := fmt.Sprintf(
		linefmt,
		ln.Address, " ",
		ln.Instruction, ln.Operand, " ",
	)

	if ln.Comment != "" {
		return str + "; " + ln.Comment
	}

	return str
}
