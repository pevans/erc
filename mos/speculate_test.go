package mos

import (
	"testing"

	"github.com/pevans/erc/asm"
	"github.com/pevans/erc/memory"
	"github.com/stretchr/testify/assert"
)

func TestCPU_SpeculateInstruction(t *testing.T) {
	t.Run("speculates LDA immediate", func(t *testing.T) {
		cpu := new(CPU)
		seg := memory.NewSegment(0x10000)
		cpu.RMem = seg

		// LDA #$42
		seg.Set(0x1000, 0xA9)
		seg.Set(0x1001, 0x42)

		line, bytesRead := cpu.SpeculateInstuction(0x1000)

		assert.Equal(t, "LDA", line.Instruction)
		assert.Equal(t, uint16(0x42), line.Operand)
		assert.Equal(t, "#$42", line.PreparedOperand)
		assert.Equal(t, uint16(2), bytesRead)
		assert.NotNil(t, line.Address)
		assert.Equal(t, 0x1000, *line.Address)
	})

	t.Run("speculates JMP absolute", func(t *testing.T) {
		cpu := new(CPU)
		seg := memory.NewSegment(0x10000)
		cpu.RMem = seg

		// JMP $1234
		seg.Set(0x2000, 0x4C)
		seg.Set(0x2001, 0x34)
		seg.Set(0x2002, 0x12)

		line, bytesRead := cpu.SpeculateInstuction(0x2000)

		assert.Equal(t, "JMP", line.Instruction)
		assert.Equal(t, uint16(0x1234), line.Operand)
		assert.Equal(t, "$1234", line.PreparedOperand)
		assert.Equal(t, uint16(3), bytesRead)
	})

	t.Run("speculates INX implied", func(t *testing.T) {
		cpu := new(CPU)
		seg := memory.NewSegment(0x10000)
		cpu.RMem = seg

		// INX
		seg.Set(0x3000, 0xE8)

		line, bytesRead := cpu.SpeculateInstuction(0x3000)

		assert.Equal(t, "INX", line.Instruction)
		assert.Equal(t, uint16(0), line.Operand)
		assert.Equal(t, "", line.PreparedOperand)
		assert.Equal(t, uint16(1), bytesRead)
	})
}

func TestCPU_Speculate(t *testing.T) {
	t.Run("speculates simple sequence", func(t *testing.T) {
		cpu := new(CPU)
		cpu.InstructionLog = asm.NewCallMap()
		seg := memory.NewSegment(0x10000)
		cpu.RMem = seg

		// 1000: LDA #$42
		// 1002: STA $20
		// 1004: RTS
		seg.Set(0x1000, 0xA9) // LDA #
		seg.Set(0x1001, 0x42)
		seg.Set(0x1002, 0x85) // STA zp
		seg.Set(0x1003, 0x20)
		seg.Set(0x1004, 0x60) // RTS

		cpu.Speculate(0x1000)

		// Debug: check if lines at our addresses exist
		lines := cpu.InstructionLog.Lines()
		found1000 := false
		found1002 := false
		found1004 := false
		for _, line := range lines {
			if len(line) >= 4 && line[:4] == "1000" {
				t.Logf("Found 1000 line: %q", line)
				found1000 = true
			}
			if len(line) >= 4 && line[:4] == "1002" {
				t.Logf("Found 1002 line: %q", line)
				found1002 = true
			}
			if len(line) >= 4 && line[:4] == "1004" {
				t.Logf("Found 1004 line: %q", line)
				found1004 = true
			}
		}

		assert.True(t, found1000, "Should have line at 1000")
		assert.True(t, found1002, "Should have line at 1002")
		assert.True(t, found1004, "Should have line at 1004")
	})

	t.Run("stops at previously executed code", func(t *testing.T) {
		cpu := new(CPU)
		cpu.InstructionLog = asm.NewCallMap()
		seg := memory.NewSegment(0x10000)
		cpu.RMem = seg

		// 1000: LDA #$42
		// 1002: STA $20
		seg.Set(0x1000, 0xA9) // LDA #
		seg.Set(0x1001, 0x42)
		seg.Set(0x1002, 0x85) // STA zp
		seg.Set(0x1003, 0x20)

		// Mark STA as already executed (non-speculative)
		addr1002 := 0x1002
		lsb := uint8(0x20)
		staLine := &asm.Line{
			Address:     &addr1002,
			Instruction: "STA",
			Opcode:      0x85,
			OperandLSB:  &lsb,
			Operand:     0x20,
		}
		cpu.InstructionLog.Add(staLine)

		cpu.Speculate(0x1000)

		// LDA should be added as speculative
		lines := cpu.InstructionLog.Lines()
		foundLDA := false
		foundSTA := false
		for _, line := range lines {
			if len(line) >= 4 && line[:4] == "1000" {
				foundLDA = true
			}
			if len(line) >= 4 && line[:4] == "1002" {
				foundSTA = true
			}
		}

		assert.True(t, foundLDA, "Should have speculative LDA at 1000")
		assert.True(t, foundSTA, "Should have non-speculative STA at 1002")
	})

	t.Run("speculates on branch target", func(t *testing.T) {
		cpu := new(CPU)
		cpu.InstructionLog = asm.NewCallMap()
		seg := memory.NewSegment(0x10000)
		cpu.RMem = seg

		// 1000: LDA #$00
		// 1002: BEQ $1006  (branch ahead 2 bytes from PC=1004)
		// 1004: LDA #$FF   (not taken path)
		// 1006: RTS        (branch target)
		seg.Set(0x1000, 0xA9) // LDA #
		seg.Set(0x1001, 0x00)
		seg.Set(0x1002, 0xF0) // BEQ
		seg.Set(0x1003, 0x02) // relative offset +2
		seg.Set(0x1004, 0xA9) // LDA #
		seg.Set(0x1005, 0xFF)
		seg.Set(0x1006, 0x60) // RTS

		cpu.Speculate(0x1000)

		// Check what got speculated
		lines := cpu.InstructionLog.Lines()
		found1000 := false
		found1002 := false

		for _, line := range lines {
			if len(line) >= 4 && line[:4] == "1000" {
				found1000 = true
			}
			if len(line) >= 4 && line[:4] == "1002" {
				found1002 = true
			}
		}

		assert.True(t, found1000, "Should speculate LDA at 1000")
		assert.True(t, found1002, "Should speculate BEQ at 1002")
	})

	t.Run("handles backward branches", func(t *testing.T) {
		cpu := new(CPU)
		cpu.InstructionLog = asm.NewCallMap()
		seg := memory.NewSegment(0x10000)
		cpu.RMem = seg

		// 1000: INX
		// 1001: CPX #$10
		// 1003: BNE $1000  (branch back -3 bytes)
		// 1005: RTS
		seg.Set(0x1000, 0xE8) // INX
		seg.Set(0x1001, 0xE0) // CPX
		seg.Set(0x1002, 0x10)
		seg.Set(0x1003, 0xD0) // BNE
		seg.Set(0x1004, 0xFB) // relative offset -5
		seg.Set(0x1005, 0x60) // RTS

		cpu.Speculate(0x1000)

		lines := cpu.InstructionLog.Lines()
		found1000 := false
		found1001 := false
		found1003 := false

		for _, line := range lines {
			if len(line) >= 4 && line[:4] == "1000" {
				found1000 = true
			}
			if len(line) >= 4 && line[:4] == "1001" {
				found1001 = true
			}
			if len(line) >= 4 && line[:4] == "1003" {
				found1003 = true
			}
		}

		assert.True(t, found1000, "Should speculate INX at 1000")
		assert.True(t, found1001, "Should speculate CPX at 1001")
		assert.True(t, found1003, "Should speculate BNE at 1003")
	})

	t.Run("avoids infinite recursion on loops", func(t *testing.T) {
		cpu := new(CPU)
		cpu.InstructionLog = asm.NewCallMap()
		seg := memory.NewSegment(0x10000)
		cpu.RMem = seg

		// 1000: JMP $1000
		seg.Set(0x1000, 0x4C) // JMP
		seg.Set(0x1001, 0x00)
		seg.Set(0x1002, 0x10)

		// Add the first instruction as already seen
		addr1000 := 0x1000
		lsbJmp := uint8(0x00)
		msbJmp := uint8(0x10)
		jmpLine := &asm.Line{
			Address:     &addr1000,
			Instruction: "JMP",
			Opcode:      0x4C,
			OperandLSB:  &lsbJmp,
			OperandMSB:  &msbJmp,
			Operand:     0x1000,
		}
		cpu.InstructionLog.Add(jmpLine)

		// This should not hang -- speculation should stop immediately
		cpu.Speculate(0x1000)
	})

	t.Run("handles nested branches", func(t *testing.T) {
		cpu := new(CPU)
		cpu.InstructionLog = asm.NewCallMap()
		seg := memory.NewSegment(0x10000)
		cpu.RMem = seg

		// 1000: LDA #$00
		// 1002: BEQ $1006
		// 1004: LDA #$01
		// 1006: BEQ $100A
		// 1008: LDA #$02
		// 100A: RTS
		seg.Set(0x1000, 0xA9) // LDA
		seg.Set(0x1001, 0x00)
		seg.Set(0x1002, 0xF0) // BEQ
		seg.Set(0x1003, 0x02) // +2
		seg.Set(0x1004, 0xA9) // LDA
		seg.Set(0x1005, 0x01)
		seg.Set(0x1006, 0xF0) // BEQ
		seg.Set(0x1007, 0x02) // +2
		seg.Set(0x1008, 0xA9) // LDA
		seg.Set(0x1009, 0x02)
		seg.Set(0x100A, 0x60) // RTS

		cpu.Speculate(0x1000)

		// Check what got speculated
		lines := cpu.InstructionLog.Lines()
		found1000 := false
		found1002 := false

		for _, line := range lines {
			if len(line) >= 4 && line[:4] == "1000" {
				found1000 = true
			}
			if len(line) >= 4 && line[:4] == "1002" {
				found1002 = true
			}
		}

		assert.True(t, found1000, "Should speculate first LDA at 1000")
		assert.True(t, found1002, "Should speculate first BEQ at 1002")
	})
}

func TestCPU_Speculate_Integration(t *testing.T) {
	t.Run("realistic branch scenario", func(t *testing.T) {
		cpu := new(CPU)
		cpu.InstructionLog = asm.NewCallMap()
		seg := memory.NewSegment(0x10000)
		cpu.RMem = seg

		// Realistic code pattern:
		// 0800: LDA $20       load a value
		// 0802: CMP #$42      compare to $42
		// 0804: BEQ success   branch if equal ($080A)
		// 0806: LDA #$00      return 0 (failure path)
		// 0808: RTS
		// 080A: LDA #$01      return 1 (success path)
		// 080C: RTS

		seg.Set(0x0800, 0xA5) // LDA
		seg.Set(0x0801, 0x20)
		seg.Set(0x0802, 0xC9) // CMP
		seg.Set(0x0803, 0x42)
		seg.Set(0x0804, 0xF0) // BEQ
		seg.Set(0x0805, 0x04) // +4
		seg.Set(0x0806, 0xA9) // LDA
		seg.Set(0x0807, 0x00)
		seg.Set(0x0808, 0x60) // RTS
		seg.Set(0x080A, 0xA9) // LDA
		seg.Set(0x080B, 0x01)
		seg.Set(0x080C, 0x60) // RTS

		cpu.Speculate(0x0800)

		// Check what got speculated
		lines := cpu.InstructionLog.Lines()
		found0800 := false
		found0802 := false
		found0804 := false

		for _, line := range lines {
			if len(line) >= 4 && line[:4] == "0800" {
				found0800 = true
			}
			if len(line) >= 4 && line[:4] == "0802" {
				found0802 = true
			}
			if len(line) >= 4 && line[:4] == "0804" {
				found0804 = true
			}
		}

		assert.True(t, found0800, "Should speculate LDA at 0800")
		assert.True(t, found0802, "Should speculate CMP at 0802")
		assert.True(t, found0804, "Should speculate BEQ at 0804")
	})

	t.Run("stops at terminating instructions", func(t *testing.T) {
		cpu := new(CPU)
		cpu.InstructionLog = asm.NewCallMap()
		seg := memory.NewSegment(0x10000)
		cpu.RMem = seg

		// 1000: LDA #$42
		// 1002: RTS
		// 1003: BRK
		seg.Set(0x1000, 0xA9) // LDA
		seg.Set(0x1001, 0x42)
		seg.Set(0x1002, 0x60) // RTS
		seg.Set(0x1003, 0x00) // BRK

		cpu.Speculate(0x1000)

		lines := cpu.InstructionLog.Lines()
		found1000 := false
		found1002 := false
		found1003 := false

		for _, line := range lines {
			if len(line) >= 4 && line[:4] == "1000" {
				found1000 = true
			}
			if len(line) >= 4 && line[:4] == "1002" {
				found1002 = true
			}
			if len(line) >= 4 && line[:4] == "1003" {
				found1003 = true
			}
		}

		assert.True(t, found1000, "Should speculate LDA at 1000")
		assert.True(t, found1002, "Should speculate RTS at 1002")
		assert.False(t, found1003, "Should NOT speculate past RTS at 1003")
	})

	t.Run("stops at JMP instruction", func(t *testing.T) {
		cpu := new(CPU)
		cpu.InstructionLog = asm.NewCallMap()
		seg := memory.NewSegment(0x10000)
		cpu.RMem = seg

		// 2000: LDA #$FF
		// 2002: JMP $3000
		// 2005: BRK
		seg.Set(0x2000, 0xA9) // LDA
		seg.Set(0x2001, 0xFF)
		seg.Set(0x2002, 0x4C) // JMP
		seg.Set(0x2003, 0x00)
		seg.Set(0x2004, 0x30)
		seg.Set(0x2005, 0x00) // BRK

		cpu.Speculate(0x2000)

		lines := cpu.InstructionLog.Lines()
		found2000 := false
		found2002 := false
		found2005 := false

		for _, line := range lines {
			if len(line) >= 4 && line[:4] == "2000" {
				found2000 = true
			}
			if len(line) >= 4 && line[:4] == "2002" {
				found2002 = true
			}
			if len(line) >= 4 && line[:4] == "2005" {
				found2005 = true
			}
		}

		assert.True(t, found2000, "Should speculate LDA at 2000")
		assert.True(t, found2002, "Should speculate JMP at 2002")
		assert.False(t, found2005, "Should NOT speculate past JMP at 2005")
	})
}
