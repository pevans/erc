package a2disk

import (
	"fmt"

	"github.com/pevans/erc/a2/a2enc"
	"github.com/pevans/erc/asm"
	"github.com/pevans/erc/memory"
	"github.com/pevans/erc/mos"
)

// This is an admittedly experimental structure. There are several ways to
// consider a disk image. One is to strictly carve it up into tracks and
// sectors. Another is to read the VTOC, then the catalog, and discover
// discrete files. Although we can begin with the former for now, it's my goal
// to support the latter over time.
type Image struct {
	Tracks []*memory.Segment
	Code   []asm.Line
}

func NewImage() *Image {
	img := new(Image)
	img.Tracks = make([]*memory.Segment, a2enc.MaxSteps/2)

	return img
}

func (img *Image) Parse(seg *memory.Segment) error {
	maxTracks := a2enc.MaxSteps / 2

	for track := range maxTracks {
		tseg := memory.NewSegment(a2enc.LogTrackLen)

		count, err := tseg.ExtractFrom(seg, track*a2enc.LogTrackLen, (track+1)*a2enc.LogTrackLen)
		if err != nil {
			return fmt.Errorf("failed to extract data from disk image: %w", err)
		}

		if count != a2enc.LogTrackLen {
			return fmt.Errorf("did not extract the number of expected bytes")
		}

		img.Tracks[track] = tseg
	}

	return nil
}

func (img *Image) Disassemble() error {
	for trackNum, track := range img.Tracks {
		offset := 0

		// The first sector of track 0 is always loaded into $0800 by the
		// Apple's bootstrap disk code. But it jumps to $0801 afterward; the
		// first byte is never executed. We should assume this is the point we
		// should be interpreting machine code.
		if trackNum == 0 {
			offset = 1
		}

		for offset < a2enc.LogTrackLen {
			line, read, err := img.DisassembleNextInstruction(track, offset)
			if err != nil {
				break
			}

			img.Code = append(img.Code, line)
			offset += read
		}
	}

	return nil
}

func (img *Image) DisassembleNextInstruction(track *memory.Segment, offset int) (asm.Line, int, error) {
	line := asm.Line{}
	read := 0
	zeroaddr := 0

	if offset+read >= track.Size() {
		return line, 0, fmt.Errorf("offset %d is beyond track size %d", offset, track.Size())
	}

	opcode := track.Get(offset + read)
	line.Address = &zeroaddr
	line.Opcode = opcode
	line.Instruction = mos.OpcodeInstruction(opcode)
	read++

	width := mos.OperandSize(opcode)

	if offset+read+width > track.Size() {
		return line, read, fmt.Errorf("instruction should have an operand but no data left")
	}

	switch width {
	case 2:
		line.Operand = track.Get16(offset + read)
		read += 2

	case 1:
		line.Operand = uint16(track.Get(offset + read))
		read++
	}

	mos.PrepareOperand(&line, uint16(offset))

	return line, read, nil
}
