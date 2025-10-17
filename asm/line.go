package asm

import "fmt"

// Line is a representation of some line of assembly to output. There
// are many kinds of assembly; this is intended to model that of a
// 6502-style system.
type Line struct {
	Address     int
	Instruction string
	Label       string

	// A formatted representation of an operand, which may include information
	// about its address mode, etc.
	PreparedOperand string

	// If an operand is provided, this will be non-nil, and will contain the
	// value of the operand.
	OperandMSB *uint8
	OperandLSB *uint8

	Opcode  uint8
	Comment string
}

// String returns some representation of a line of assembly. There's no
// single grammar for assembly -- it's usually a notation that works for
// a specific assembler. As long as it "looks right", that's good enough
// for now.
func (ln Line) String() string {
	fmtOper := " "

	switch {
	case ln.OperandLSB != nil && ln.OperandMSB != nil:
		fmtOper = fmt.Sprintf(
			"%02X %02X", *ln.OperandLSB, *ln.OperandMSB,
		)
	case ln.OperandLSB != nil:
		fmtOper = fmt.Sprintf(
			"%02X", *ln.OperandLSB,
		)
	}

	linefmt := "%04X" + // address
		":%02X" + // opcode
		" %-5s" + // operand
		" | " + // spacing
		"%-8s " + // label
		"%s " + // instruction
		"%-10s" + // operand
		"%5s" // spacing

	str := fmt.Sprintf(
		linefmt,
		ln.Address, ln.Opcode, fmtOper,
		ln.Label, ln.Instruction, ln.PreparedOperand, " ",
	)

	if ln.Comment != "" {
		return str + ln.Comment
	}

	// Add line padding for what seem to be the end of a subroutine.
	if ln.Instruction == "JMP" || ln.Instruction == "RTS" {
		str += "\n"
	}

	return str
}
