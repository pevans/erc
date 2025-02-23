package asm

import "fmt"

// Line is a representation of some line of assembly to output. There
// are many kinds of assembly; this is intended to model that of a
// 6502-style system.
type Line struct {
	Address     int
	Instruction string
	Operand     string
	Opcode      uint8
	Comment     string
}

// String returns some representation of a line of assembly. There's no
// single grammar for assembly -- it's usually a notation that works for
// a specific assembler. As long as it "looks right", that's good enough
// for now.
func (ln Line) String() string {
	linefmt := "$%04X" + // address
		"%3s" +
		"%02X" +
		"%3s" + // spacing
		"%s " + // instruction
		"%-10s" + // operand
		"%5s" // spacing

	str := fmt.Sprintf(
		linefmt,
		ln.Address, " ",
		ln.Opcode, " ",
		ln.Instruction, ln.Operand, " ",
	)

	if ln.Comment != "" {
		return str + "; " + ln.Comment
	}

	return str
}
