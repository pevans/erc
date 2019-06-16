package mos65c02

import (
	"os"
	"reflect"
	"runtime"
	"strings"

	"github.com/pevans/erc/pkg/asmrec/a2rec"
	"github.com/pevans/erc/pkg/data"
)

// A CPU is an implementation of an MOS 65c02 processor.
type CPU struct {
	// RMem and WMem are the segments from which we will read or write
	// whenever it is necessary.
	RMem data.Getter
	WMem data.Setter

	// This is the current address mode that the CPU is operating
	// within. The address mode affects how the CPU will determine the
	// effective address for an instruction.
	AddrMode int

	// The Opcode is the byte which indicates both the instruction and
	// the address mode of the instruction that we must carry out.
	Opcode data.Byte

	// The Operand is the one or two bytes which is an argument to the
	// opcode.
	Operand data.DByte

	// This is the effective address for the current operation. The
	// effective address is the one computed by the address mode, taking
	// into account the current state of the CPU and the current operand
	// for the instruction.
	EffAddr data.DByte

	// The effective value is data that the instruction wants after the
	// effective address is dereferenced. In some cases, the instruction
	// only cares about an address, and this may be zero; in other
	// cases, the instruction does not take an address, and EffVal is
	// all it cares about. In yet other cases, both this and the EffAddr
	// may be zero because the behavior of the instruction is implied
	// and cannot be modified by any operand.
	EffVal data.Byte

	// PC is the Program Counter. It is where the processor
	// will look to execute its next instruction.
	PC data.DByte

	// The A register is the Accumulator. You can think of the
	// accumulator as similar to how old calculators work; arithmetic
	// operations will add to, subtract from, etc., this register.
	A data.Byte

	// The X and Y registers are most often treated as indexes for
	// loops, but can also be treated as general-purpose registers to
	// hold onto numbers.
	X, Y data.Byte

	// The P register doesn't seem to have a formal name, but I like to
	// think of it as the Predicator. Its bits are used to indicate
	// several statuses the CPU can have; 1 to mean the status is on, 0
	// to mean it is off.
	P data.Byte

	// The S register is the Stack pointer. The stack in the MOS 6502
	// processor is in memory page 1 ($100 - $1FF); the S register
	// value is treated as an offset from $100. S will begin at $FF and
	// decrease as the stack depth increases.
	S data.Byte
}

// An Instruction is a function that performs an operation on the CPU.
type Instruction func(c *CPU)

// An AddrMode is a function which resolves what the effective address
// (EffAddr) is given the current state of the CPU.
type AddrMode func(c *CPU)

// This block defines the flags that we recognize within the status
// register.
const (
	CARRY     = data.Byte(1)
	ZERO      = data.Byte(2)
	INTERRUPT = data.Byte(4)
	DECIMAL   = data.Byte(8)
	BREAK     = data.Byte(16)
	UNUSED    = data.Byte(32)
	OVERFLOW  = data.Byte(64)
	NEGATIVE  = data.Byte(128)
)

// While here we define the address modes that we can work with.
const (
	amNoa = iota // no address mode
	amAcc        // accumulator
	amAbs        // absolute
	amAbx        // absolute x-index
	amAby        // absolute y-index
	amBy2        // Consume 2 bytes (for NP2)
	amBy3        // Consume 3 bytes (for NP3)
	amImm        // immediate
	amImp        // implied
	amInd        // indirect
	amIdx        // x-index indirect
	amIdy        // indirect y-index
	amRel        // relative
	amZpg        // zero page
	amZpx        // zero page x-index
	amZpy        // zero page y-index
)

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

//   00   01   02   03   04   05   06   07   08   09   0A   0B   0C   0D   0E   0F
var addrModes = [256]AddrMode{
	Imp, Idx, By2, Imp, Zpg, Zpg, Zpg, Imp, Imp, Imm, Acc, Imp, Abs, Abs, Abs, Imp, // 0x
	Rel, Idy, Zpg, Imp, Zpg, Zpx, Zpx, Imp, Imp, Aby, Acc, Imp, Abs, Abx, Abx, Imp, // 1x
	Abs, Idx, By2, Imp, Zpg, Zpg, Zpg, Imp, Imp, Imm, Acc, Imp, Abs, Abs, Abs, Imp, // 2x
	Rel, Idy, Zpg, Imp, Zpx, Zpx, Zpx, Imp, Imp, Aby, Acc, Imp, Abx, Abx, Abx, Imp, // 3x
	Imp, Idx, By2, Imp, By2, Zpg, Zpg, Imp, Imp, Imm, Acc, Imp, Abs, Abs, Abs, Imp, // 4x
	Rel, Idy, Zpg, Imp, By2, Zpx, Zpx, Imp, Imp, Aby, Imp, Imp, By3, Abx, Abx, Imp, // 5x
	Imp, Idx, By2, Imp, Zpg, Zpg, Zpg, Imp, Imp, Imm, Acc, Imp, Ind, Abs, Abs, Imp, // 6x
	Rel, Idy, Zpg, Imp, Zpx, Zpx, Zpx, Imp, Imp, Aby, Imp, Imp, Abx, Abx, Abx, Imp, // 7x
	Rel, Idx, By2, Imp, Zpg, Zpg, Zpg, Imp, Imp, Imm, Imp, Imp, Abs, Abs, Abs, Imp, // 8x
	Rel, Idy, Zpg, Imp, Zpx, Zpx, Zpy, Imp, Imp, Aby, Imp, Imp, Abs, Abx, Abx, Imp, // 9x
	Imm, Idx, Imm, Imp, Zpg, Zpg, Zpg, Imp, Imp, Imm, Imp, Imp, Abs, Abs, Abs, Imp, // Ax
	Rel, Idy, Zpg, Imp, Zpx, Zpx, Zpy, Imp, Imp, Aby, Imp, Imp, Abx, Abx, Aby, Imp, // Bx
	Imm, Idx, By2, Imp, Zpg, Zpg, Zpg, Imp, Imp, Imm, Imp, Imp, Abs, Abs, Abs, Imp, // Cx
	Rel, Idy, Zpg, Imp, By2, Zpx, Zpx, Imp, Imp, Aby, Imp, Imp, By3, Abx, Abx, Imp, // Dx
	Imm, Idx, By2, Imp, Zpg, Zpg, Zpg, Imp, Imp, Imm, Imp, Imp, Abs, Abs, Abs, Imp, // Ex
	Rel, Idy, Zpg, Imp, By2, Zpx, Zpx, Imp, Imp, Aby, Imp, Imp, By3, Abx, Abx, Imp, // Fx
}

// The offsets table defines the number of bytes we must increment the
// PC register after a given instruction. The bytes vary based on
// address mode, rather than the specific instruction. In cases where
// the instruction would change the PC due to its defined behavior, the
// offset is given as zero.
//
//  0  1  2  3  4  5  6  7  8  9  A  B  C  D  E  F
var offsets = [256]data.DByte{
	1, 2, 3, 1, 2, 2, 2, 1, 1, 2, 1, 1, 3, 3, 3, 1, // 0x
	0, 2, 2, 1, 2, 2, 2, 1, 1, 3, 1, 1, 3, 3, 3, 1, // 1x
	0, 2, 3, 1, 2, 2, 2, 1, 1, 2, 1, 1, 3, 3, 3, 1, // 2x
	0, 2, 2, 1, 2, 2, 2, 1, 1, 3, 1, 1, 3, 3, 3, 1, // 3x
	0, 2, 3, 1, 3, 2, 2, 1, 1, 2, 1, 1, 0, 3, 3, 1, // 4x
	0, 2, 2, 1, 3, 2, 2, 1, 1, 3, 1, 1, 4, 3, 3, 1, // 5x
	0, 2, 3, 1, 2, 2, 2, 1, 1, 2, 1, 1, 0, 3, 3, 1, // 6x
	0, 2, 2, 1, 2, 2, 2, 1, 1, 3, 1, 1, 0, 3, 3, 1, // 7x
	0, 2, 3, 1, 2, 2, 2, 1, 1, 2, 1, 1, 3, 3, 3, 1, // 8x
	0, 2, 2, 1, 2, 2, 2, 1, 1, 3, 1, 1, 3, 3, 3, 1, // 9x
	2, 2, 2, 1, 2, 2, 2, 1, 1, 2, 1, 1, 3, 3, 3, 1, // Ax
	0, 2, 2, 1, 2, 2, 2, 1, 1, 3, 1, 1, 3, 3, 3, 1, // Bx
	2, 2, 3, 1, 2, 2, 2, 1, 1, 2, 1, 1, 3, 3, 3, 1, // Cx
	0, 2, 2, 1, 3, 2, 2, 1, 1, 3, 1, 1, 4, 3, 3, 1, // Dx
	2, 2, 3, 1, 2, 2, 2, 1, 1, 2, 1, 1, 3, 3, 3, 1, // Ex
	0, 2, 2, 1, 3, 2, 2, 1, 1, 3, 1, 1, 4, 3, 3, 1, // Fx
}

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

func (i Instruction) String() string {
	var (
		funcName = runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
		parts    = strings.Split(funcName, ".")
	)

	return strings.ToUpper(parts[len(parts)-1])
}

func (m AddrMode) String() string {
	var (
		funcName = runtime.FuncForPC(reflect.ValueOf(m).Pointer()).Name()
		parts    = strings.Split(funcName, ".")
	)

	return strings.ToUpper(parts[len(parts)-1])
}

// ApplyStatus will make a status update for the given flag based upon
// cond being true or not.
func (c *CPU) ApplyStatus(cond bool, flag data.Byte) {
	c.P &= ^flag
	if cond {
		c.P |= flag
	}
}

// ApplyN will apply the normal negative status check (which is whether
// the eighth bit is high or not).
func (c *CPU) ApplyN(val data.Byte) {
	c.ApplyStatus(val&0x80 > 0, NEGATIVE)
}

// ApplyZ will apply the normal zero status check, which is literally if
// val is zero or not.
func (c *CPU) ApplyZ(val data.Byte) {
	c.ApplyStatus(val == 0, ZERO)
}

// ApplyNZ will apply both the normal negative and zero checks.
func (c *CPU) ApplyNZ(val data.Byte) {
	c.ApplyN(val)
	c.ApplyZ(val)
}

// Compare will compute the difference between the given base and the
// current EffVal value of c. ApplyNZ is called on the result. CARRY is
// set if the result is greater than zero.
func Compare(c *CPU, base data.Byte) {
	res := base - c.EffVal

	c.ApplyNZ(res)
	c.ApplyStatus(base > c.EffVal, CARRY)
}

// Get will return the byte at a given address.
func (c *CPU) Get(addr data.DByte) data.Byte {
	return c.RMem.Get(addr)
}

// Set will set the byte at a given address to the given value.
func (c *CPU) Set(addr data.DByte, val data.Byte) {
	c.WMem.Set(addr, val)
}

// Get16 returns a 16-bit value at a given address, which is read in
// little-endian order.
func (c *CPU) Get16(addr data.DByte) data.DByte {
	lsb := c.RMem.Get(addr)
	msb := c.RMem.Get(addr + 1)

	return (data.DByte(msb) << 8) | data.DByte(lsb)
}

// Set16 sets the two bytes beginning at the given address to the given
// value. The bytes are set in little-endian order.
func (c *CPU) Set16(addr data.DByte, val data.DByte) {
	lsb := data.Byte(val & 0xFF)
	msb := data.Byte(val >> 8)

	c.WMem.Set(addr, lsb)
	c.WMem.Set(addr+1, msb)
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

	rec.Record(os.Stdout)

	// Adjust the program counter to beyond the expected instruction
	// sequence (1 byte for the opcode, + N bytes for the operand, based
	// on address mode).
	c.PC += offsets[opcode]

	// We always apply BREAK and UNUSED after each execution, mostly in
	// observance for how other emulators have handled this step.
	c.P |= UNUSED | BREAK

	return nil
}
