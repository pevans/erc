package mos65c02

import (
	"reflect"
	"runtime"
	"strings"

	"github.com/pevans/erc/pkg/asmrec/a2rec"
	"github.com/pevans/erc/pkg/data"
)

// An Instruction is a function that performs an operation on the CPU.
type Instruction func(c *CPU)

// Below is a table of instructions that are mapped to opcodes. For
// corresponding address modes, see addr.go.
//   00   01   02   03   04   05   06   07   08   09   0A   0B   0C   0D   0E   0F
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

/*
//  0  1  2  3  4  5  6  7  8  9  A  B  C  D  E  F
var cycles = [256]data.Byte{
	7, 6, 2, 1, 5, 3, 5, 1, 3, 2, 2, 1, 6, 4, 6, 1, // 0x
	2, 5, 5, 1, 5, 4, 6, 1, 2, 4, 2, 1, 6, 4, 6, 1, // 1x
	6, 6, 2, 1, 3, 3, 5, 1, 4, 2, 2, 1, 4, 4, 6, 1, // 2x
	2, 5, 5, 1, 4, 4, 6, 1, 2, 4, 2, 1, 4, 4, 6, 1, // 3x
	6, 6, 2, 1, 3, 3, 5, 1, 3, 2, 2, 1, 3, 4, 6, 1, // 4x
	2, 5, 5, 1, 4, 4, 6, 1, 2, 4, 3, 1, 8, 4, 6, 1, // 5x
	6, 6, 2, 1, 3, 3, 5, 1, 4, 2, 2, 1, 5, 4, 6, 1, // 6x
	2, 5, 5, 1, 4, 4, 6, 1, 2, 4, 4, 1, 6, 4, 6, 1, // 7x
	3, 6, 2, 1, 3, 3, 3, 1, 2, 2, 2, 1, 4, 4, 4, 1, // 8x
	2, 6, 5, 1, 4, 4, 4, 1, 2, 5, 2, 1, 4, 5, 5, 1, // 9x
	2, 6, 2, 1, 3, 3, 3, 1, 2, 2, 2, 1, 4, 4, 4, 1, // Ax
	2, 5, 5, 1, 4, 4, 4, 1, 2, 4, 2, 1, 4, 4, 4, 1, // Bx
	2, 6, 2, 1, 3, 3, 5, 1, 2, 2, 2, 1, 4, 4, 3, 1, // Cx
	2, 5, 5, 1, 4, 4, 6, 1, 2, 4, 3, 1, 4, 4, 7, 1, // Dx
	2, 6, 2, 1, 3, 3, 5, 1, 2, 2, 2, 1, 4, 4, 6, 1, // Ex
	2, 5, 5, 1, 4, 4, 6, 1, 2, 4, 4, 1, 4, 4, 7, 1, // Fx
}
*/

// String composes an instruction function into a string and returns
// that
func (i Instruction) String() string {
	var (
		funcName = runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
		parts    = strings.Split(funcName, ".")
	)

	return strings.ToUpper(parts[len(parts)-1])
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
		inst   Instruction
		mode   AddrMode
		opcode data.Byte
		rec    a2rec.Recorder
	)

	opcode = c.Get(c.PC)
	mode = addrModes[opcode]
	inst = instructions[opcode]

	rec.PC = c.PC
	rec.PrintState = true
	rec.Opcode = opcode
	rec.A = c.A
	rec.X = c.X
	rec.Y = c.Y
	rec.S = c.S
	rec.P = c.P
	rec.Inst = inst.String()
	rec.Mode = mode.String()

	// NOTE: neither the address mode resolver nor the instruction
	// handler have any error conditions. This is by design: they DO NOT
	// error out. They handle whatever situation comes up.

	// Resolve the values of EffAddr and EffVal by executing the address
	// mode handler.
	mode(c)

	rec.Operand = c.Operand
	rec.EffAddr = c.EffAddr
	rec.EffVal = c.EffVal

	// Now execute the instruction
	inst(c)

	// Record the operation, but let the rest of the func complete even
	// if this errors
	var err error
	if c.RecWriter != nil {
		err = rec.Record(c.RecWriter)
	}

	srec := rec
	srec.PrintState = false
	_ = c.SMap.Map(int(srec.PC), &srec)

	// Adjust the program counter to beyond the expected instruction
	// sequence (1 byte for the opcode, + N bytes for the operand, based
	// on address mode).
	c.PC += offsets[opcode]

	// We always apply BREAK and UNUSED after each execution, mostly in
	// observance for how other emulators have handled this step.
	c.P |= UNUSED | BREAK

	return err
}
