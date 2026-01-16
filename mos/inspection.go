package mos

import (
	"fmt"

	"github.com/pevans/erc/a2/a2sym"
	"github.com/pevans/erc/elog"
)

// ExplainInstruction modifies a given elog.Instruction in an effort to better
// explain what's happening in the code. It will add comments, if possible; it
// will swap operands with labels (e.g. for subroutines) and variable names
// (such as those defined with EQU).
func ExplainInstruction(line *elog.Instruction, pc uint16, effAddr uint16) {
	addr := int(effAddr)
	addrMode := addrModes[line.Opcode]

	line.EndOfBlock = endsBlock(line.Opcode)

	if maybeRoutine(line.Opcode) {
		if routine := a2sym.Subroutine(addr); routine != "" {
			line.PreparedOperand = routine
			return
		}
	}

	if OpcodeReadsMemory(line.Opcode) {
		if rs := a2sym.ReadSwitch(addr); rs.Mode != a2sym.ModeNone {
			line.Comment = rs.String()
		}
	}

	if ws := a2sym.WriteSwitch(addr); ws.Mode != a2sym.ModeNone {
		line.Comment = ws.String()
	}

	if addrMode == AmZPG || addrMode == AmABS {
		if variable := a2sym.Variable(addr); variable != "" {
			line.PreparedOperand = variable
		}
	}

	if variable := a2sym.Variable(int(line.Operand)); variable != "" {
		switch addrMode {
		case AmIDX:
			line.PreparedOperand = fmt.Sprintf("(%v,X)", variable)
		case AmIDY:
			line.PreparedOperand = fmt.Sprintf("(%v),Y", variable)
		case AmZPX:
			line.PreparedOperand = fmt.Sprintf("%v,X", variable)
		case AmZPY:
			line.PreparedOperand = fmt.Sprintf("%v,Y", variable)
		case AmIND:
			line.PreparedOperand = fmt.Sprintf("(%v)", variable)
		}
	}

	if routine := a2sym.Subroutine(int(pc)); routine != "" {
		line.Label = routine
	}
}

// maybeRoutine returns true if the opcode represents a jump in control flow
func maybeRoutine(opcode uint8) bool {
	return opcode == 0x20 || // JSR
		opcode == 0x4C || // JMP (ABS)
		opcode == 0x6C || // JMP (IND)
		opcode == 0x7C || // JMP (ABX)
		addrModes[opcode] == AmREL // any branch
}

// endsBlock returns true if we think this instruction represents the logical
// "end" of some block of code. Branches don't count, but returns do (RTI,
// RTS), as do JMPs.
func endsBlock(opcode uint8) bool {
	return opcode == 0x40 || // RTI
		opcode == 0x60 || // RTS
		opcode == 0x4C || // JMP (ABS)
		opcode == 0x6C || // JMP (IND)
		opcode == 0x7C // JMP (ABX)
}
