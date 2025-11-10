package mos

import (
	"github.com/pevans/erc/asm"
)

// Speculate takes an address of some code we _could_ execute -- for example,
// from a branch -- and walks through the code in memory at that point,
// recording in our instruction channel. It'll keep recursing through other
// branch points, and will only stop if it reaches code we've previously
// recorded.
func (c *CPU) Speculate(addr uint16) {
	for {
		line, read := c.SpeculateInstuction(addr)

		// We should stop speculating if we've hit a line that we've seen
		// before. This will work whether the line previously recorded was
		// speculative or not.
		if c.InstructionLog.Exists(line) {
			return
		}

		c.InstructionLog.Add(line)

		// If this is a branch, we want to speculate on what might happen if
		// the branch is taken, but not if the branch is ignored.
		if isBranch(line.Opcode) {
			// This is pretty funky. OperandLSB is an 8-bit number that is
			// meant to be signed, so we convert to int8 to allow the most
			// significant bit to retain its signed-ness.
			offset := int8(*line.OperandLSB)

			// We want to apply the potentially-negative offset to the address
			// so we know where we really ought to branch. We add two so that
			// we account for the branch opcode and operand.
			branchAddr := int16(addr) + int16(offset) + 2

			c.Speculate(uint16(branchAddr))
			return
		}

		// There are several opcodes which signal we should go no further
		if shouldEndSpeculation(line.Opcode) {
			return
		}

		// Keep on loopin'!
		addr += read
	}
}

func isBranch(opcode uint8) bool {
	return addrModes[opcode] == AmREL
}

func shouldEndSpeculation(opcode uint8) bool {
	switch opcode {
	case 0x00:
		// These would be BRK, which is unusual, and _probably_ a bad opcode
		return true

	case 0x40, 0x60:
		// These would return control to something on the stack (RTI, RTS)
		return true

	case 0x20, 0x4C, 0x6C, 0x7C:
		// Any JMPs or JSRs are calls which we can't know would return control
		// back to the caller
		return true
	}

	// These opcodes aren't used. They are encoded as NOPs but the "true" NOP
	// opcode is 0xEA.
	switch opcode & 0xF {
	case 0x3, 0x7, 0xB, 0xF:
		return true
	}

	return false
}

// This will record what _would_ be executed by the instruction located at
// `addr`, and return the asm.Line along with the number of bytes that compose
// that instruction (i.e. the opcode plus its operand, if one exists).
func (c *CPU) SpeculateInstuction(addr uint16) (*asm.Line, uint16) {
	line := &asm.Line{
		Speculative: true,
	}

	opcode := c.Get(addr)

	line.EndOfBlock = endsBlock(opcode)

	// Since we need to point to an integer, we need to make a copy of addr,
	// then reference it
	lineAddress := int(addr)
	line.Address = &lineAddress

	line.Opcode = opcode
	line.Instruction = instructionNames[opcode]

	width := OperandSize(opcode)

	switch width {
	case 2:
		line.Operand = c.Get16(addr + 1)
		lsb, msb := uint8(line.Operand&0xFF), uint8(line.Operand>>8)
		line.OperandLSB = &lsb
		line.OperandMSB = &msb
	case 1:
		line.Operand = uint16(c.Get(addr + 1))
		lsb := uint8(line.Operand)
		line.OperandLSB = &lsb
	}

	PrepareOperand(line, addr)
	// ExplainInstruction(line, addr, addr)

	return line, width + 1
}
