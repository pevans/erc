package mos

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/a2/a2sym"
	"github.com/pevans/erc/asm"
	"github.com/pevans/erc/internal/metrics"
)

// An Instruction is a function that performs an operation on the CPU.
type Instruction func(c *CPU)

// Below is a table of instructions that are mapped to opcodes. For
// corresponding address modes, see addr.go.
//
//	00   01   02   03   04   05   06   07   08   09   0A   0B   0C   0D   0E   0F
var instructions = [256]Instruction{
	Brk, Ora, Np2, Nop, Tsb, Ora, Asl, Nop, Php, Ora, Asl, Nop, Tsb, Ora, Asl, Nop, // 0x
	Bpl, Ora, Ora, Nop, Trb, Ora, Asl, Nop, Clc, Ora, Inc, Nop, Trb, Ora, Asl, Nop, // 1x
	Jsr, And, Np2, Nop, Bit, And, Rol, Nop, Plp, And, Rol, Nop, Bit, And, Rol, Nop, // 2x
	Bmi, And, And, Nop, Bit, And, Rol, Nop, Sec, And, Dec, Nop, Bit, And, Rol, Nop, // 3x
	Rti, Eor, Np2, Nop, Np2, Eor, Lsr, Nop, Pha, Eor, Lsr, Nop, Jmp, Eor, Lsr, Nop, // 4x
	Bvc, Eor, Eor, Nop, Np2, Eor, Lsr, Nop, Cli, Eor, Phy, Nop, Np3, Eor, Lsr, Nop, // 5x
	Rts, Adc, Np2, Nop, Stz, Adc, Ror, Nop, Pla, Adc, Ror, Nop, Jmp, Adc, Ror, Nop, // 6x
	Bvs, Adc, Adc, Nop, Stz, Adc, Ror, Nop, Sei, Adc, Ply, Nop, Jmp, Adc, Ror, Nop, // 7x
	Bra, Sta, Np2, Nop, Sty, Sta, Stx, Nop, Dey, Bim, Txa, Nop, Sty, Sta, Stx, Nop, // 8x
	Bcc, Sta, Sta, Nop, Sty, Sta, Stx, Nop, Tya, Sta, Txs, Nop, Stz, Sta, Stz, Nop, // 9x
	Ldy, Lda, Ldx, Nop, Ldy, Lda, Ldx, Nop, Tay, Lda, Tax, Nop, Ldy, Lda, Ldx, Nop, // Ax
	Bcs, Lda, Lda, Nop, Ldy, Lda, Ldx, Nop, Clv, Lda, Tsx, Nop, Ldy, Lda, Ldx, Nop, // Bx
	Cpy, Cmp, Np2, Nop, Cpy, Cmp, Dec, Nop, Iny, Cmp, Dex, Nop, Cpy, Cmp, Dec, Nop, // Cx
	Bne, Cmp, Cmp, Nop, Np2, Cmp, Dec, Nop, Cld, Cmp, Phx, Nop, Np3, Cmp, Dec, Nop, // Dx
	Cpx, Sbc, Np2, Nop, Cpx, Sbc, Inc, Nop, Inx, Sbc, Nop, Nop, Cpx, Sbc, Inc, Nop, // Ex
	Beq, Sbc, Sbc, Nop, Np2, Sbc, Inc, Nop, Sed, Sbc, Plx, Nop, Np3, Sbc, Inc, Nop, // Fx
}

// 0  1  2  3  4  5  6  7  8  9  A  B  C  D  E  F
var cycles = [256]uint8{
	7, 6, 2, 1, 5, 3, 5, 1, 3, 2, 2, 1, 6, 4, 6, 1, // 0x
	2, 5, 5, 1, 5, 4, 6, 1, 2, 4, 2, 1, 6, 4, 6, 1, // 1x
	6, 6, 2, 1, 3, 3, 5, 1, 4, 2, 2, 1, 4, 4, 6, 1, // 2x
	2, 5, 5, 1, 4, 4, 6, 1, 2, 4, 2, 1, 4, 4, 6, 1, // 3x
	6, 6, 2, 1, 3, 3, 5, 1, 3, 2, 2, 1, 3, 4, 6, 1, // 4x
	2, 5, 5, 1, 4, 4, 6, 1, 2, 4, 3, 1, 8, 4, 6, 1, // 5x
	6, 6, 2, 1, 3, 3, 5, 1, 4, 2, 2, 1, 6, 4, 6, 1, // 6x
	2, 5, 5, 1, 4, 4, 6, 1, 2, 4, 4, 1, 6, 4, 6, 1, // 7x
	2, 6, 2, 1, 3, 3, 3, 1, 2, 2, 2, 1, 4, 4, 4, 1, // 8x
	2, 6, 5, 1, 4, 4, 4, 1, 2, 5, 2, 1, 4, 5, 5, 1, // 9x
	2, 6, 2, 1, 3, 3, 3, 1, 2, 2, 2, 1, 4, 4, 4, 1, // Ax
	2, 5, 5, 1, 4, 4, 4, 1, 2, 4, 2, 1, 4, 4, 4, 1, // Bx
	2, 6, 2, 1, 3, 3, 5, 1, 2, 2, 2, 1, 4, 4, 6, 1, // Cx
	2, 5, 5, 1, 4, 4, 6, 1, 2, 4, 3, 1, 4, 4, 7, 1, // Dx
	2, 6, 2, 1, 3, 3, 5, 1, 2, 2, 2, 1, 4, 4, 6, 1, // Ex
	2, 5, 5, 1, 4, 4, 6, 1, 2, 4, 4, 1, 4, 4, 7, 1, // Fx
}

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
	case
		0xAE,
		0xA2,
		0xA6,
		0xB6,
		0xBE:
		return true // LDX
	case
		0xA0,
		0xA4,
		0xAC,
		0xB4,
		0xBC:
		return true // LDY
	}

	return false
}

func (c *CPU) Cycles() int {
	return int(cycles[c.Opcode])
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

	c.Opcode = c.Get(c.PC)
	mode = addrModes[c.Opcode]
	inst = instructions[c.Opcode]

	c.State.SetBool(
		a2state.InstructionReadOp,
		OpcodeReadsMemory(c.Opcode),
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
	c.PC += offsets[c.Opcode]

	// We always apply BREAK and UNUSED after each execution, mostly in
	// observance for how other emulators have handled this step.
	c.P |= UNUSED | BREAK

	if c.ClockEmulator != nil {
		c.ClockEmulator.WaitForCycles(int64(cycles[c.Opcode]), time.Sleep)
	}

	return nil
}

func (c *CPU) Status() string {
	return fmt.Sprintf(
		"A:$%02X X:$%02X Y:$%02X S:$%02X P:$%02X (%s) PC:$%04X EA:$%04X EV:$%02X",
		c.A, c.X, c.Y, c.P, c.S, formatStatus(c.P), c.PC, c.EffAddr, c.EffVal,
	)
}

func (c *CPU) LastInstruction() string {
	ln := asm.Line{
		Address:     int(c.LastPC),
		Instruction: instructions[c.Opcode].String(),
		Opcode:      c.Opcode,
		Operand: formatOperand(
			c.AddrMode, c.Operand, c.LastPC,
		),
		Comment: c.explainInstruction(c.Opcode),
	}

	return ln.String()
}

// NextInstruction returns a string representing the next opcode that
// would be executed
func (c *CPU) NextInstruction() string {
	opcode := c.Get(c.PC)
	mode := addrModes[opcode]

	// Copy the CPU so we don't alter our own operand, effective
	// address, etc.
	copyOfCPU := c

	// There are some cases where resolving the address mode may mutate
	// the CPU or state map, so we use DebuggerLookAhead to let the
	// address mode code know what's about to happen.
	copyOfCPU.State.SetBool(a2state.DebuggerLookAhead, true)
	mode(copyOfCPU)
	copyOfCPU.State.SetBool(a2state.DebuggerLookAhead, false)

	ln := asm.Line{
		Address:     int(c.PC),
		Instruction: instructions[opcode].String(),
		Operand: formatOperand(
			copyOfCPU.AddrMode, copyOfCPU.Operand, c.PC,
		),
		Comment: c.explainInstruction(opcode),
	}

	return ln.String()
}

func (c *CPU) explainInstruction(opcode uint8) string {
	addr := int(c.EffAddr)

	if isJSR(opcode) {
		if routine := a2sym.Subroutine(addr); routine != "" {
			return fmt.Sprintf("subroutine %v", routine)
		}
	}

	if c.State.Bool(a2state.InstructionReadOp) {
		if rs := a2sym.ReadSwitch(addr); rs.Mode != a2sym.ModeNone {
			return rs.String()
		}
	}

	if ws := a2sym.WriteSwitch(addr); ws.Mode != a2sym.ModeNone {
		return ws.String()
	}

	return ""
}

func isJSR(opcode uint8) bool {
	return opcode == 0x20
}

func formatOperand(
	mode int,
	operand uint16,
	pc uint16,
) string {
	str := ""

	switch mode {
	case AmACC, AmIMP, AmBY2, AmBY3:
		break
	case AmABS:
		str = fmt.Sprintf("$%04X", operand)
	case AmABX:
		str = fmt.Sprintf("$%04X,X", operand)
	case AmABY:
		str = fmt.Sprintf("$%04X,Y", operand)
	case AmIDX:
		str = fmt.Sprintf("($%02X,X)", operand)
	case AmIDY:
		str = fmt.Sprintf("($%02X),Y", operand)
	case AmIND:
		str = fmt.Sprintf("($%04X)", operand)
	case AmIMM:
		str = fmt.Sprintf("#$%02X", operand)
	case AmREL:
		newAddr := pc + operand + 2

		// It's signed, so the effect of the operand should be negative w/r/t
		// two's complement.
		if operand >= 0x80 {
			newAddr -= 256
		}

		str = fmt.Sprintf("$%04X", newAddr)
	case AmZPG:
		str = fmt.Sprintf("$%02X", operand)
	case AmZPX:
		str = fmt.Sprintf("$%02X,X", operand)
	case AmZPY:
		str = fmt.Sprintf("$%02X,Y", operand)
	}

	if str == "" {
		return fmt.Sprintf("mode=%v operand=%x", mode, operand)
	}

	return str
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
