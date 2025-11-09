package a2disk_test

import (
	"testing"

	"github.com/pevans/erc/a2/a2disk"
	"github.com/pevans/erc/a2/a2enc"
	"github.com/pevans/erc/memory"
	"github.com/stretchr/testify/assert"
)

func TestNewImage(t *testing.T) {
	assert.NotNil(t, a2disk.NewImage())
}

func TestImage_Parse(t *testing.T) {
	cases := []struct {
		name  string
		size  int
		setup func(*memory.Segment)
		errfn assert.ErrorAssertionFunc
	}{
		{
			name: "valid full DOS 3.3 disk image",
			size: a2enc.DosSize,
			setup: func(seg *memory.Segment) {
				for i := range seg.Size() {
					seg.Set(i, uint8(i%256))
				}
			},
			errfn: assert.NoError,
		},
		{
			name: "segment too small",
			size: a2enc.LogTrackLen - 1,
			setup: func(seg *memory.Segment) {
				for i := range seg.Size() {
					seg.Set(i, 0xFF)
				}
			},
			errfn: assert.Error,
		},
		{
			name: "empty segment",
			size: 0,
			setup: func(seg *memory.Segment) {
			},
			errfn: assert.Error,
		},
		{
			name: "single track",
			size: a2enc.LogTrackLen,
			setup: func(seg *memory.Segment) {
				for i := range seg.Size() {
					seg.Set(i, 0xAA)
				}
			},
			errfn: assert.Error,
		},
		{
			name: "partial disk image",
			size: a2enc.LogTrackLen * 20,
			setup: func(seg *memory.Segment) {
				for i := range seg.Size() {
					seg.Set(i, uint8(i&0xFF))
				}
			},
			errfn: assert.Error,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			seg := memory.NewSegment(c.size)
			c.setup(seg)

			img := a2disk.NewImage()
			err := img.Parse(seg)
			c.errfn(t, err)
		})
	}
}

func TestImage_DisassembleNextInstruction(t *testing.T) {
	cases := []struct {
		name              string
		bytes             []byte
		offset            int
		expectedInstr     string
		expectedOperand   uint16
		expectedPrepared  string
		expectedBytesRead int
	}{
		{
			name:              "LDA immediate",
			bytes:             []byte{0xA9, 0x42}, // LDA #$42
			offset:            0,
			expectedInstr:     "LDA",
			expectedOperand:   0x42,
			expectedPrepared:  "#$42",
			expectedBytesRead: 2,
		},
		{
			name:              "STA absolute",
			bytes:             []byte{0x8D, 0x00, 0x04}, // STA $0400
			offset:            0,
			expectedInstr:     "STA",
			expectedOperand:   0x0400,
			expectedPrepared:  "$0400",
			expectedBytesRead: 3,
		},
		{
			name:              "JMP absolute",
			bytes:             []byte{0x4C, 0x34, 0x12}, // JMP $1234
			offset:            0,
			expectedInstr:     "JMP",
			expectedOperand:   0x1234,
			expectedPrepared:  "$1234",
			expectedBytesRead: 3,
		},
		{
			name:              "INX implied",
			bytes:             []byte{0xE8}, // INX
			offset:            0,
			expectedInstr:     "INX",
			expectedOperand:   0,
			expectedPrepared:  "",
			expectedBytesRead: 1,
		},
		{
			name:              "LDA zero page",
			bytes:             []byte{0xA5, 0x20}, // LDA $20
			offset:            0,
			expectedInstr:     "LDA",
			expectedOperand:   0x20,
			expectedPrepared:  "$20",
			expectedBytesRead: 2,
		},
		{
			name:              "LDA absolute,X",
			bytes:             []byte{0xBD, 0x00, 0x20}, // LDA $2000,X
			offset:            0,
			expectedInstr:     "LDA",
			expectedOperand:   0x2000,
			expectedPrepared:  "$2000,X",
			expectedBytesRead: 3,
		},
		{
			name:              "BNE relative",
			bytes:             []byte{0xD0, 0xFE}, // BNE $FE (relative -2)
			offset:            0,
			expectedInstr:     "BNE",
			expectedOperand:   0xFE,
			expectedPrepared:  "$0000", // Will be calculated based on PC=0
			expectedBytesRead: 2,
		},
		{
			name:              "instruction at offset",
			bytes:             []byte{0x00, 0x00, 0xA9, 0x10}, // LDA #$10 at offset 2
			offset:            2,
			expectedInstr:     "LDA",
			expectedOperand:   0x10,
			expectedPrepared:  "#$10",
			expectedBytesRead: 2,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			seg := memory.NewSegment(len(c.bytes))
			for i, b := range c.bytes {
				seg.Set(i, b)
			}

			img := a2disk.NewImage()
			line, bytesRead, err := img.DisassembleNextInstruction(seg, c.offset)

			assert.NoError(t, err)
			assert.Equal(t, c.expectedInstr, line.Instruction)
			assert.Equal(t, c.expectedOperand, line.Operand)
			assert.Equal(t, c.expectedPrepared, line.PreparedOperand)
			assert.Equal(t, c.expectedBytesRead, bytesRead)
		})
	}
}

func TestImage_Disassemble(t *testing.T) {
	t.Run("disassembles simple program", func(t *testing.T) {
		// Create a simple program:
		// LDA #$42
		// STA $0400
		// INX
		// RTS
		program := []byte{
			0xA9, 0x42, // LDA #$42
			0x8D, 0x00, 0x04, // STA $0400
			0xE8, // INX
			0x60, // RTS
		}

		// Create a disk image with one track containing the program
		// Note: Track 0, Sector 0 is loaded to $0800, but execution starts at $0801
		// So we place the program starting at offset 1
		seg := memory.NewSegment(a2enc.DosSize)
		seg.Set(0, 0x00)
		for i, b := range program {
			seg.Set(i+1, b)
		}

		img := a2disk.NewImage()
		err := img.Parse(seg)
		assert.NoError(t, err)

		err = img.Disassemble()
		assert.NoError(t, err)

		assert.NotEmpty(t, img.Code)

		assert.Equal(t, "LDA", img.Code[0].Instruction)
		assert.Equal(t, uint16(0x42), img.Code[0].Operand)

		assert.Equal(t, "STA", img.Code[1].Instruction)
		assert.Equal(t, uint16(0x0400), img.Code[1].Operand)

		assert.Equal(t, "INX", img.Code[2].Instruction)

		assert.Equal(t, "RTS", img.Code[3].Instruction)
	})

	t.Run("disassembles multiple tracks", func(t *testing.T) {
		seg := memory.NewSegment(a2enc.DosSize)

		// Put a JMP instruction at the start of track 0, but not at the first
		// byte
		seg.Set(0, 0x00) // First byte is skipped (placeholder)
		seg.Set(1, 0x4C) // JMP $1234
		seg.Set(2, 0x34)
		seg.Set(3, 0x12)

		// Put an LDA instruction at the start of track 1
		seg.Set(a2enc.LogTrackLen, 0xA9) // LDA #$FF
		seg.Set(a2enc.LogTrackLen+1, 0xFF)

		img := a2disk.NewImage()
		err := img.Parse(seg)
		assert.NoError(t, err)

		err = img.Disassemble()
		assert.NoError(t, err)

		// First instruction should be JMP
		assert.Equal(t, "JMP", img.Code[0].Instruction)
		assert.Equal(t, uint16(0x1234), img.Code[0].Operand)

		// Find the LDA instruction in track 1
		foundLDA := false
		for _, line := range img.Code {
			if line.Instruction == "LDA" && line.Operand == 0xFF {
				foundLDA = true
				break
			}
		}
		assert.True(t, foundLDA, "Should find LDA #$FF instruction from track 1")
	})
}
