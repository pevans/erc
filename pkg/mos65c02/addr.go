package mos65c02

import (
	"reflect"
	"runtime"
	"strings"
)

// An AddrMode is a function which resolves what the effective address
// (EffAddr) is given the current state of the CPU.
type AddrMode func(c *CPU)

// AddrSpace is the total number of addressable values that an MOS 65c02
// processor can work with. It's possible for a computer to have more
// RAM than there is addressable space--and the Apple II is an example
// of this--but for that to work, the computer must introduce some kind
// of bank-switch mechanism that can swap our segments of RAM.
const AddrSpace = 0x10000

const (
	_     = iota // no address mode
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

// Below is an address mode table that maps mode functions to specific
// opcodes.
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
var offsets = [256]uint16{
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

// String will figure out what address mode function this is and return
// it in string form.
func (m AddrMode) String() string {
	var (
		funcName = runtime.FuncForPC(reflect.ValueOf(m).Pointer()).Name()
		parts    = strings.Split(funcName, ".")
	)

	return strings.ToUpper(parts[len(parts)-1])
}

// Acc will resolve the Accumulator address mode, which is very simple:
// the effective value is the data within the A register.
//
// Ex. INC increments the A register
func Acc(c *CPU) {
	c.EffVal = c.A
	c.EffAddr = 0
	c.Operand = 0
	c.AddrMode = amAcc
}

// Abs resolves the Absolute address mode. Given a 16-bit operand, we
// dereference that for our effective value.
//
// Ex. INC $1234 increments the byte at $1234
func Abs(c *CPU) {
	c.EffAddr = c.Get16(c.PC + 1)
	c.Operand = c.EffAddr
	c.EffVal = c.Get(c.EffAddr)
	c.AddrMode = amAbs
}

// Abx resolves Absolute X address mode, which is like Absolute mode but
// adds the X register content to the operand.
//
// Ex. INC $1234,X increments the byte at $1234 + X
func Abx(c *CPU) {
	c.Operand = c.Get16(c.PC + 1)
	c.EffAddr = c.Operand + uint16(c.X)
	c.EffVal = c.Get(c.EffAddr)
	c.AddrMode = amAbx
}

// Aby is much like ABX, except that it adds the Y register content.
//
// Ex. INC $1234,Y increments the byte at $1234 + Y
func Aby(c *CPU) {
	c.Operand = c.Get16(c.PC + 1)
	c.EffAddr = c.Operand + uint16(c.Y)
	c.EffVal = c.Get(c.EffAddr)
	c.AddrMode = amAby
}

// By2 is a placeholder mode
func By2(c *CPU) {
	c.EffAddr = 0
	c.EffVal = 0
	c.Operand = 0
	c.AddrMode = amBy2
}

// By3 is a placeholder mode
func By3(c *CPU) {
	c.EffAddr = 0
	c.EffVal = 0
	c.Operand = 0
	c.AddrMode = amBy3
}

// Imm resolves Immediate address mode. The operand is the literal
// effective value we will observe.
//
// Ex. ADC #$12 adds $12 to the A register
func Imm(c *CPU) {
	c.EffAddr = 0
	c.EffVal = c.Get(c.PC + 1)
	c.Operand = uint16(c.EffVal)
	c.AddrMode = amImm
}

// Imp resolves the implied address mode. IMP is used to handle cases
// where whatever the opcode does has an implied, singular purpose, and
// cannot be modified by any operand.
func Imp(c *CPU) {
	c.EffVal = 0
	c.EffAddr = 0
	c.Operand = 0
	c.AddrMode = amImp
}

// Ind resolves the indirect address mode. If you can imagine that the
// ABS mode describes a kind of pointer, then we can say that IND
// describes a double pointer: that is, a pointer to a pointer.
//
// Ex. JMP ($NNNN) resolves the address at $NNNN, then jumps to that
// address.
func Ind(c *CPU) {
	// The inner part of the operand `$NNNN` is the address of... yet
	// another address; so we derefence that `($NNNN)` to get the value.
	c.Operand = c.Get16(c.PC + 1)
	c.EffAddr = c.Get16(c.Operand)
	c.EffVal = c.Get(c.EffAddr)
	c.AddrMode = amInd
}

// Idx resolves the indexed indirect address mode, which resolves to a
// zero page address that points to another address that is our final
// destination. In a practical sense, IDX can be used to loop over a
// table of pointers to other things that is located in the zero
// page--which may not come up too often because you don't have a ton of
// space in the zero page to work with. But I can imagine operating
// system code that could use this quite often.
//
// Ex. INC ($NN,X) resolves $NN + X as <addr1>; then looks up <addr1>
// and resolves that to <addr2>; the effective address becomes <addr2>,
// and the effective value is the byte at <addr2>.
func Idx(c *CPU) {
	operand := c.Get(c.PC + 1)
	c.Operand = uint16(operand)

	// Our effective address is the dereferenced value found at the base
	// address.
	c.EffAddr = c.Get16(uint16(operand + c.X))
	c.EffVal = c.Get(c.EffAddr)
	c.AddrMode = amIdx
}

// Idy resolves the indirect indexed address mode, which is essentially
// a zero-page pointer that addresses something else and let's us loop
// on the _something else_.
//
// Ex. INC ($NN),Y looks up $NN and resolves that to <addr1>; then looks
// up <addr1>, and resolves that + Y as <addr2>; then saves <addr2> as
// the effective address and the looked-up byte at <addr2>.
func Idy(c *CPU) {
	// The base address for the instruction; the `$NN` part of the
	// operand.
	c.Operand = uint16(c.Get(c.PC + 1))

	// This dereferences the base address, essentially resolving the
	// `()` part of the operand.
	effAddr := c.Get16(c.Operand)

	// And here we account for the `,Y` part; Y is added to the
	// dereferenced address.
	c.EffAddr = effAddr + uint16(c.Y)
	c.EffVal = c.Get(c.EffAddr)
	c.AddrMode = amIdy
}

// Rel resolves the relative address mode. This is only used by branch
// instructions, and lets you define an offset of -128..127 that the
// branch instructions can jump to. Note the given operand is, uniquely
// among other address modes, a _signed_ 8-bit operand. If you imagine
// that JSR and JMP enable long-range jumps, then the branches (via REL)
// enable short-range jumps both forwards and backwards.
//
// Ex. BEQ COUNTER (where "COUNTER" is some label that compiles to a
// signed byte offset) will jump to COUNTER if the Z status is set;
// otherwise, no action is taken besides stepping past the branch
// instruction.
func Rel(c *CPU) {
	// The next byte is the signed offset of where we're going; positive
	// = forward, negative = backward.
	c.Operand = uint16(c.Get(c.PC + 1))

	// But we don't want to convert change (or addr) into a valid
	// address yet. We want the uint16-ness of addresses in the MOS 6502
	// to force overflow to work as expected: going from FFFF to 0000,
	// or 0000 to FFFF. We add 2 more bytes to account for the (fixed)
	// size of all branch instruction sequences, which is one byte for
	// the opcode and one byte for the operand.
	addr := c.PC + c.Operand + 2

	// Because negative numbers in the MOS 6502 are encoded with
	// twos-complement, if change has its eighth bit set to 1, then we
	// need to perform a subtraction to get the desired value.
	if c.Operand > 127 {
		addr -= 256
	}

	c.EffAddr = addr
	c.EffVal = 0
	c.AddrMode = amRel
}

// Zpg resolves the zero page address mode. This is most analogous to ABS,
// except that instead of a two-byte operand, it takes a one-byte
// operand that can only be in the zero page.
//
// Ex. INC $12 increments the byte at $12 by one.
func Zpg(c *CPU) {
	c.Operand = uint16(c.Get(c.PC + 1))
	c.EffAddr = c.Operand
	c.EffVal = c.Get(c.EffAddr)
	c.AddrMode = amZpg
}

// Zpx resolves the zero page x address mode. This is analogous to ABX,
// except it takes a one-byte operand.
//
// Ex. INC $12,X increments the byte at $12 + X by one.
func Zpx(c *CPU) {
	operand := c.Get(c.PC + 1)
	c.Operand = uint16(operand)
	c.EffAddr = uint16(operand + c.X)
	c.EffVal = c.Get(c.EffAddr)
	c.AddrMode = amZpx
}

// Zpy resolves the zero page y address mode. This is analogous to ABY,
// except it takes a one-byte operand.
//
// Ex. INC $12,Y increments the byte at $12 + Y by one.
func Zpy(c *CPU) {
	operand := c.Get(c.PC + 1)
	c.Operand = uint16(operand)
	c.EffAddr = uint16(operand + c.Y)
	c.EffVal = c.Get(c.EffAddr)
	c.AddrMode = amZpy
}
