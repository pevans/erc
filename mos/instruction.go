package mos

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/a2/a2sym"
	"github.com/pevans/erc/asm"
	"github.com/pevans/erc/internal/metrics"
)

// An Instruction is a function that performs an operation on the CPU.
type Instruction func(c *CPU)

// String composes an instruction function into a string and returns
// that
func (i Instruction) String() string {
	var (
		funcName = runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
		parts    = strings.Split(funcName, ".")
	)

	return strings.ToUpper(parts[len(parts)-1])
}

// OpcodeReadsMemory returns true for all opcodes we would consider as
// "reading" data in memory. This is useful, for instance, for soft
// switches which care about whether data is being read from somewhere
// or written somewhere.
func OpcodeReadsMemory(opcode uint8) bool {
	switch opcode {
	case
		0xA1,
		0xA5,
		0xA9,
		0xAD,
		0xB1,
		0xB5,
		0xB9,
		0xBD:
		return true // LDA
	case 0xAE, 0xA2, 0xA6, 0xB6, 0xBE:
		return true // LDX
	case 0xA0, 0xA4, 0xAC, 0xB4, 0xBC:
		return true // LDY
	}

	return false
}

func (c *CPU) OpcodeCycles() int {
	cyc := int(cycles[c.opcode])
	effPage := c.EffAddr & 0xFF00

	switch addrModes[c.opcode] {
	case AmABX, AmABY:
		// We may be crossing page boundaries; if so, we need to add a 1-cycle
		// penalty
		basePage := c.Operand & 0xFF00

		if basePage != effPage {
			cyc++
		}

	case AmIDY:
		// Similar to ABX/ABY, we may need to add a 1-cycle penalty for
		// crossing boundaries. The logic is slightly different because of how
		// IDY works.
		baseAddr := c.EffAddr - uint16(c.Y)
		basePage := baseAddr & 0xFF00

		if basePage != effPage {
			cyc++
		}

	case AmREL:
		// The number of cycles consumed by a branch are variable based on its
		// outcomes, which may also factor in a page-cross penalty.
		nextPC := c.LastPC + 2
		if c.EffAddr != nextPC {
			cyc++

			if (nextPC & 0xFF00) != (c.EffAddr & 0xFF00) {
				cyc++
			}
		}
	}

	// Setting the D flag causes instructions to take one cycle longer than
	// normal.
	//
	// NOTE: I'm not sure if this applies universally or only for arithmetic
	// operations. In practice, code typically sets decimal only for
	// arithmetic operations, then unsets.
	if c.P&DECIMAL > 0 {
		cyc++
	}

	return cyc
}

// Execute will process through one instruction and return. While doing
// so, the CPU state will update such that it moves to the next
// instruction. Note that the MOS 65C02 processor can execute
// indefinitely; while there are definitely parts of memory that don't
// really house opcodes (the zero page being one such part), the
// processor would absolutely try to execute those if the PC register
// pointed at those parts. And technically, if PC incremented beyond the
// 0xFFFF address, it would simply overflow back to the zero page.
func (c *CPU) Execute() error {
	var (
		inst Instruction
		mode AddrMode
	)

	metrics.Increment("instructions", 1)

	// We want to record the current PC before it might change as the
	// result of any instruction we execute
	c.LastPC = c.PC

	c.opcode = c.Get(c.PC)
	mode = addrModeFuncs[c.opcode]
	inst = instructions[c.opcode]

	c.State.SetBool(
		a2state.InstructionReadOp,
		OpcodeReadsMemory(c.opcode),
	)

	// NOTE: neither the address mode resolver nor the instruction
	// handler have any error conditions. This is by design: they DO NOT
	// error out. They handle whatever situation comes up.

	// Resolve the values of EffAddr and EffVal by executing the address
	// mode handler.
	mode(c)

	// Now execute the instruction
	inst(c)

	// Adjust the program counter to beyond the expected instruction
	// sequence (1 byte for the opcode, + N bytes for the operand, based
	// on address mode).
	c.PC += offsets[c.opcode]

	// We always apply BREAK and UNUSED after each execution, mostly in
	// observance for how other emulators have handled this step.
	c.P |= UNUSED | BREAK

	if c.State.Bool(a2state.DebugImage) && c.InstructionChannel != nil {
		select {
		case c.InstructionChannel <- c.LastInstructionLine(int(c.OpcodeCycles())):
		// Sent successfully
		default:
			// Channel full, drop this instruction log
		}
	}

	c.cycleCounter += uint64(c.OpcodeCycles())

	return nil
}

func (c *CPU) Status() string {
	return fmt.Sprintf(
		"A:$%02X X:$%02X Y:$%02X S:$%02X P:$%02X (%s) PC:$%04X EA:$%04X EV:$%02X",
		c.A, c.X, c.Y, c.S, c.P, formatStatus(c.P), c.PC, c.EffAddr, c.EffVal,
	)
}

func (c *CPU) ThisInstruction() string {
	pc := int(c.PC)

	line := &asm.Line{
		Address:     &pc,
		Instruction: instructions[c.opcode].String(),
		Operand:     c.Operand,
		Opcode:      c.opcode,
	}

	PrepareOperand(line, c.PC)

	return fmt.Sprintf(
		"%04X:%v %v",
		*line.Address, line.Instruction, line.PreparedOperand,
	)
}

func (c *CPU) LastInstructionLine(cycles int) *asm.Line {
	lastPC := int(c.LastPC)
	line := &asm.Line{
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

func (c *CPU) LastInstruction() string {
	ln := c.LastInstructionLine(0)

	return ln.String()
}

// NextInstruction returns a string representing the next opcode that
// would be executed
func (c *CPU) NextInstruction() string {
	opcode := c.Get(c.PC)
	mode := addrModeFuncs[opcode]

	// Copy the CPU so we don't alter our own operand, effective
	// address, etc.
	copyOfCPU := c

	// There are some cases where resolving the address mode may mutate
	// the CPU or state map, so we use DebuggerLookAhead to let the
	// address mode code know what's about to happen.
	copyOfCPU.State.SetBool(a2state.DebuggerLookAhead, true)
	mode(copyOfCPU)
	copyOfCPU.State.SetBool(a2state.DebuggerLookAhead, false)

	pc := int(c.PC)
	ln := &asm.Line{
		Address:     &pc,
		Instruction: instructions[opcode].String(),
		Operand:     c.Operand,
	}

	PrepareOperand(ln, c.PC)
	ExplainInstruction(ln, c.PC, c.EffAddr)

	return ln.String()
}

func ExplainInstruction(line *asm.Line, pc uint16, effAddr uint16) {
	addr := int(effAddr)
	addrMode := addrModes[line.Opcode]

	line.EndOfBlock = endsBlock(line.Opcode)

	if maybeRoutine(line.Opcode) {
		if routine := a2sym.Subroutine(addr); routine != "" {
			line.PreparedOperand = routine
			return
		}
	}

	// FIXME: this is a pretty bad hack so we don't have to test the
	// InstructionReadOp state in the CPU
	if addrMode == AmABX || addrMode == AmABY {
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

func maybeRoutine(opcode uint8) bool {
	return opcode == 0x20 || // JSR
		opcode == 0x4C || // JMP (ABS)
		opcode == 0x6C || // JMP (IND)
		opcode == 0x7C || // JMP (ABX)
		addrModes[opcode] == AmREL // any branch
}

func PrepareOperand(line *asm.Line, pc uint16) {
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

func OpcodeInstruction(opcode uint8) string {
	return instructions[opcode].String()
}

func endsBlock(opcode uint8) bool {
	return opcode == 0x40 || // RTI
		opcode == 0x60 || // RTS
		opcode == 0x4C || // JMP (ABS)
		opcode == 0x6C || // JMP (IND)
		opcode == 0x7C // JMP (ABX)
}
