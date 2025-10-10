package mos

import (
	"math"

	"github.com/pevans/erc/asm"
	"github.com/pevans/erc/clock"
	"github.com/pevans/erc/memory"
)

// A CPU is an implementation of an MOS 65c02 processor.
type CPU struct {
	RMem  memory.Getter
	WMem  memory.Setter
	State *memory.StateMap

	// A map of instructions that we have executed. This is only used when
	// we're debugging an image.
	InstructionLog *asm.CallMap

	ClockEmulator *clock.Emulator

	// The total number of cycles executed by the CPU. This can overflow, so
	// code that cares about the cycle count needs to bear that in mind.
	CycleCount int

	// This is the current address mode that the CPU is operating
	// within. The address mode affects how the CPU will determine the
	// effective address for an instruction.
	AddrMode int

	// The Opcode is the byte which indicates both the instruction and
	// the address mode of the instruction that we must carry out.
	Opcode uint8

	// ReadOp will be true if the last operation was a "read", as
	// opposed to a "write". This is important for some soft switch
	// logic.
	ReadOp bool

	// The Operand is the one or two bytes which is an argument to the
	// opcode.
	Operand uint16

	// This is the effective address for the current operation. The
	// effective address is the one computed by the address mode, taking
	// into account the current state of the CPU and the current operand
	// for the instruction.
	EffAddr uint16

	// The effective value is data that the instruction wants after the
	// effective address is dereferenced. In some cases, the instruction
	// only cares about an address, and this may be zero; in other
	// cases, the instruction does not take an address, and EffVal is
	// all it cares about. In yet other cases, both this and the EffAddr
	// may be zero because the behavior of the instruction is implied
	// and cannot be modified by any operand.
	EffVal uint8

	// PC is the Program Counter. It is where the processor
	// will look to execute its next instruction. LastPC is the address
	// of the last instruction that was executed.
	PC     uint16
	LastPC uint16

	// The A register is the Accumulator. You can think of the
	// accumulator as similar to how old calculators work; arithmetic
	// operations will add to, subtract from, etc., this register.
	A uint8

	// The X and Y registers are most often treated as indexes for
	// loops, but can also be treated as general-purpose registers to
	// hold onto numbers.
	X, Y uint8

	// The P register doesn't seem to have a formal name, but I like to
	// think of it as the Predicator. Its bits are used to indicate
	// several statuses the CPU can have; 1 to mean the status is on, 0
	// to mean it is off.
	P uint8

	// The S register is the Stack pointer. The stack in the MOS 6502
	// processor is in memory page 1 ($100 - $1FF); the S register
	// value is treated as an offset from $100. S will begin at $FF and
	// decrease as the stack depth increases.
	S uint8
}

func (c *CPU) CyclesSince(lastCycle int) int {
	// if CycleCount is less than lastCycle, we probably overflowed
	if c.CycleCount < lastCycle {
		return c.CycleCount + (math.MaxInt - lastCycle)
	}

	return c.CycleCount - lastCycle
}
