package asm

import "fmt"

// Line is a representation of some line of assembly to output. There
// are many kinds of assembly; this is intended to model that of a
// 6502-style system.
type Line struct {
	// Address is the address at which this instruction was executed.
	Address *int

	// Instruction is some string representation of the instruction being
	// executed in the line.
	Instruction string

	// Label is an optional feature of a line of code. If provided, it will
	// mark that this line may be jumped or branched with that label rather
	// than the raw address in memory at which the instruction was executed.
	Label string

	// PreparedOperand is a formatted representation of an operand, which may
	// include information about its address mode, etc.
	PreparedOperand string

	// OperandMSB and OperandLSB are the most and least signficant bytes that
	// comprise an operand. If these are nil, then we treat the instruction
	// has not having had an operand.
	OperandMSB *uint8
	OperandLSB *uint8

	// Operand is the full 16-bit operand for some instruction. This will be
	// zero even if the instruction does not technically have an operand.
	Operand uint16

	// Opcode is a numeric representation for the precise instruction and
	// address mode that was executed. (That is, one instruction may have many
	// opcodes, one for each address mode in which it may be run.)
	Opcode uint8

	// Comment is some descriptive comment that is appended to the end of the
	// line.
	Comment string

	// Cycles is the number of CPU cycles consumed by this instruction.
	Cycles int

	// Speculative is true when this line should be considered as
	// "speculative" execution: an instruction that did not run, but _would
	// have,_ had the conditions been right. Some branches are not taken in
	// the code, but had they been, this line would represent an instruction
	// from that block.
	Speculative bool

	// EndOfBlock is true if we regard this line of execution as being the
	// end of some subroutine.
	EndOfBlock bool
}

// ShortString returns a shortened version of String. This version includes
// only the address, instruction, and operand.
func (ln Line) ShortString() string {
	linefmt := "%s" + // address
		" | " + // spacing
		"%s " + // instruction
		"%-10s" // operand

	str := fmt.Sprintf(
		linefmt,
		ln.onlyAddress(),
		ln.Instruction, ln.PreparedOperand,
	)

	return str
}

// String returns some representation of a line of assembly. There's no
// single grammar for assembly -- it's usually a notation that works for
// a specific assembler. As long as it "looks right", that's good enough
// for now.
func (ln Line) String() string {
	linefmt := "%s" + // address
		"%-8s " + // label
		"%s " + // instruction
		"%-10s" + // operand
		"%5s" + // spacing
		"%s" // comment

	str := fmt.Sprintf(
		linefmt,
		ln.fullAddress(),
		ln.Label, ln.Instruction, ln.PreparedOperand,
		" ", ln.Comment,
	)

	return str
}

// onlyAddress returns a string form of the address in memory.
func (ln Line) onlyAddress() string {
	if ln.Address == nil {
		return ""
	}

	return fmt.Sprintf("%04X | ", *ln.Address)
}

// fullAddress returns a string form of the address and bytes in memory at
// which this instruction operates. It includes the opcode and operand bytes
// with the address where those are located.
func (ln Line) fullAddress() string {
	if ln.Address == nil {
		return ""
	}

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

	addrfmt := "%04X" + // address
		":%02X" + // opcode
		" %-5s | " // operand

	return fmt.Sprintf(addrfmt, *ln.Address, ln.Opcode, fmtOper)
}
