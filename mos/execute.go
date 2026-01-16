package mos

import (
	"reflect"
	"runtime"
	"strings"

	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/internal/metrics"
)

// An Instruction is a function that performs an operation on the CPU.
type Instruction func(c *CPU)

// String composes an instruction function into a string and returns that
func (i Instruction) String() string {
	var (
		funcName = runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
		parts    = strings.Split(funcName, ".")
	)

	return strings.ToUpper(parts[len(parts)-1])
}

// OpcodeReadsMemory returns true for all opcodes we would consider as
// "reading" data in memory. This is useful, for instance, for soft switches
// which care about whether data is being read from somewhere or written
// somewhere.
func OpcodeReadsMemory(opcode uint8) bool {
	switch opcode {
	case 0xA1, 0xA5, 0xA9, 0xAD, 0xB1, 0xB2, 0xB5, 0xB9, 0xBD:
		return true // LDA
	case 0xA2, 0xA6, 0xAE, 0xB6, 0xBE:
		return true // LDX
	case 0xA0, 0xA4, 0xAC, 0xB4, 0xBC:
		return true // LDY
	case 0x24, 0x2C, 0x34, 0x3C:
		return true // BIT (not IMM)
	case 0x01, 0x05, 0x0D, 0x11, 0x12, 0x15, 0x19, 0x1D:
		return true // ORA (not IMM)
	case 0x21, 0x25, 0x2D, 0x31, 0x32, 0x35, 0x39, 0x3D:
		return true // AND (not IMM)
	case 0x41, 0x45, 0x4D, 0x51, 0x52, 0x55, 0x59, 0x5D:
		return true // EOR (not IMM)
	case 0x61, 0x65, 0x6D, 0x71, 0x72, 0x75, 0x79, 0x7D:
		return true // ADC (not IMM)
	case 0xE1, 0xE5, 0xED, 0xF1, 0xF2, 0xF5, 0xF9, 0xFD:
		return true // SBC (not IMM)
	case 0xC1, 0xC5, 0xCD, 0xD1, 0xD2, 0xD5, 0xD9, 0xDD:
		return true // CMP (not IMM)
	case 0xE4, 0xEC:
		return true // CPX (not IMM)
	case 0xC4, 0xCC:
		return true // CPY (not IMM)
	case 0xE6, 0xEE, 0xF6, 0xFE:
		return true // INC (not ACC)
	case 0xC6, 0xCE, 0xD6, 0xDE:
		return true // DEC (not ACC)
	case 0x06, 0x0E, 0x16, 0x1E:
		return true // ASL (not ACC)
	case 0x46, 0x4E, 0x56, 0x5E:
		return true // LSR (not ACC)
	case 0x26, 0x2E, 0x36, 0x3E:
		return true // ROL (not ACC)
	case 0x66, 0x6E, 0x76, 0x7E:
		return true // ROR (not ACC)
	case 0x14, 0x1C:
		return true // TRB
	case 0x04, 0x0C:
		return true // TSB
	}

	return false
}

// OpcodeCycles returns the number of cycles consumed by a given opcode. This
// number can be influenced by a few factors: if we're crossing page
// boundaries for memory access, including when branching; if we're performing
// some decimal operation.
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

// Execute will process through one instruction and return. While doing so,
// the CPU state will update such that it moves to the next instruction. Note
// that the MOS 65C02 processor can execute indefinitely; while there are
// definitely parts of memory that don't really house opcodes (the zero page
// being one such part), the processor would absolutely try to execute those
// if the PC register pointed at those parts. And technically, if PC
// incremented beyond the 0xFFFF address, it would simply overflow back to the
// zero page.
func (c *CPU) Execute() error {
	var (
		inst Instruction
		mode AddrMode
	)

	metrics.Increment("instructions", 1)

	// We want to record the current PC before it might change as the result
	// of any instruction we execute
	c.LastPC = c.PC

	c.opcode = c.Get(c.PC)
	mode = addrModeFuncs[c.opcode]
	inst = instructions[c.opcode]

	c.State.SetBool(
		a2state.InstructionReadOp,
		OpcodeReadsMemory(c.opcode),
	)

	// NOTE: neither the address mode resolver nor the instruction handler
	// have any error conditions. This is by design: they DO NOT error out.
	// They handle whatever situation comes up.

	// Resolve the values of EffAddr and EffVal by executing the address mode
	// handler.
	mode(c)

	// Now execute the instruction
	inst(c)

	// Adjust the program counter to beyond the expected instruction sequence
	// (1 byte for the opcode, + N bytes for the operand, based on address
	// mode).
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
