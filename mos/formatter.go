package mos

import (
	"fmt"

	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/elog"
)

// Status returns a formatted string that concisely captures the state of the
// CPU chip (its registers, position in memory for execution, and the most
// recent effective address and value on which it operated).
func (c *CPU) Status() string {
	return fmt.Sprintf(
		"A:$%02X X:$%02X Y:$%02X S:$%02X P:$%02X (%s) PC:$%04X EA:$%04X EV:$%02X",
		c.A, c.X, c.Y, c.S, c.P, formatStatus(c.P), c.PC, c.EffAddr, c.EffVal,
	)
}

// CurrentInstructionShort returns the current instruction that we would
// execute, if asked, in the form of a formatted string (including address in
// memory, instruction, and formatted operand).
func (c *CPU) CurrentInstructionShort() string {
	pc := int(c.PC)

	line := &elog.Instruction{
		Address:     &pc,
		Instruction: instructions[c.opcode].String(),
		Operand:     c.Operand,
		Opcode:      c.opcode,
	}

	PrepareOperand(line, c.PC)

	return line.ShortString()
}

// LastInstructionLine returns the last instruction that was executed in the
// form of an elog.Instruction object.
func (c *CPU) LastInstructionLine(cycles int) *elog.Instruction {
	lastPC := int(c.LastPC)
	line := &elog.Instruction{
		Address:     &lastPC,
		Instruction: instructions[c.opcode].String(),
		Opcode:      c.opcode,
		Operand:     c.Operand,
		Cycles:      cycles,
	}

	PrepareOperand(line, c.LastPC)
	ExplainInstruction(line, c.LastPC, c.EffAddr)

	return line
}

// LastInstruction returns the last instruction that was executed in the form
// of a formatted string.
func (c *CPU) LastInstruction() string {
	ln := c.LastInstructionLine(0)

	return ln.String()
}

// NextInstruction returns a string representing the next opcode that would be
// executed. There's some speculative inference happening -- we'll actually
// execute the address mode code to see what the effective address and value
// will look like.
func (c *CPU) NextInstruction() string {
	opcode := c.Get(c.PC)
	mode := addrModeFuncs[opcode]

	// Copy the CPU so we don't alter our own operand, effective address, etc.
	// Note that this won't copy memory segments, etc. Notably the statemap
	// (State) used below will be shared between the original and the copy
	// CPU.
	copyOfCPU := *c

	// There are some cases where resolving the address mode may mutate the
	// CPU or state map, so we use DebuggerLookAhead to let the address mode
	// code know what's about to happen.
	copyOfCPU.State.SetBool(a2state.DebuggerLookAhead, true)
	mode(&copyOfCPU)
	copyOfCPU.State.SetBool(a2state.DebuggerLookAhead, false)

	pc := int(c.PC)
	ln := &elog.Instruction{
		Address:     &pc,
		Instruction: instructions[opcode].String(),
		Operand:     c.Operand,
	}

	PrepareOperand(ln, c.PC)
	ExplainInstruction(ln, c.PC, c.EffAddr)

	return ln.String()
}

// PrepareOperand will fill in a provided elog.Instruction object with a
// formatted operand (specifically modifying the PreparedOperand, OperandLSB,
// and OperandMSB fields). The provided pc (program counter) value is used to
// calculate the branch address, given that branch operands are relative
// values.
func PrepareOperand(line *elog.Instruction, pc uint16) {
	addrMode := OpcodeAddrMode(line.Opcode)
	lsb := uint8(line.Operand & 0xFF)
	msb := uint8(line.Operand >> 8)

	switch addrMode {
	case AmACC, AmIMP, AmBY2, AmBY3:
		break
	case AmABS:
		line.PreparedOperand = fmt.Sprintf("$%04X", line.Operand)
		line.OperandLSB = &lsb
		line.OperandMSB = &msb
	case AmABX:
		line.PreparedOperand = fmt.Sprintf("$%04X,X", line.Operand)
		line.OperandLSB = &lsb
		line.OperandMSB = &msb
	case AmABY:
		line.PreparedOperand = fmt.Sprintf("$%04X,Y", line.Operand)
		line.OperandLSB = &lsb
		line.OperandMSB = &msb
	case AmIDX:
		line.PreparedOperand = fmt.Sprintf("($%02X,X)", line.Operand)
		line.OperandLSB = &lsb
	case AmIDY:
		line.PreparedOperand = fmt.Sprintf("($%02X),Y", line.Operand)
		line.OperandLSB = &lsb
	case AmIND:
		line.PreparedOperand = fmt.Sprintf("($%04X)", line.Operand)
		line.OperandLSB = &lsb
		line.OperandMSB = &msb
	case AmIMM:
		line.PreparedOperand = fmt.Sprintf("#$%02X", line.Operand)
		line.OperandLSB = &lsb
	case AmREL:
		newAddr := pc + line.Operand + 2

		// It's signed, so the effect of the operand should be negative w/r/t
		// two's complement.
		if line.Operand >= 0x80 {
			newAddr -= 256
		}

		line.PreparedOperand = fmt.Sprintf("$%04X", newAddr)
		line.OperandLSB = &lsb
	case AmZPG:
		line.PreparedOperand = fmt.Sprintf("$%02X", line.Operand)
		line.OperandLSB = &lsb
	case AmZPX:
		line.PreparedOperand = fmt.Sprintf("$%02X,X", line.Operand)
		line.OperandLSB = &lsb
	case AmZPY:
		line.PreparedOperand = fmt.Sprintf("$%02X,Y", line.Operand)
		line.OperandLSB = &lsb
	}
}

// formatStatus returns a string form of the status flags available in the P
// register (one character per flag). Each flag is ordered based on where the
// bit is located in the register.
func formatStatus(p uint8) string {
	pstatus := []rune("NVUBDIZC")

	for i := 7; i >= 0; i-- {
		bit := (p >> uint(i)) & 1
		if bit == 0 {
			pstatus[7-i] = '.'
		}
	}

	return string(pstatus)
}

// OpcodeInstruction returns the string label of an opcode's instruction.
func OpcodeInstruction(opcode uint8) string {
	return instructions[opcode].String()
}
