// Package assembler provides a minimal 65C02 assembler that produces bootable
// Apple II DOS 3.3 disk images for black box testing.
package assembler

import (
	"fmt"
	"strings"

	"github.com/pevans/erc/a2/a2enc"
)

const (
	defaultOrigin = 0x0801
	maxCodeSize   = 4095 // track 0 capacity minus the one-byte sector count
)

// Assemble takes 65C02 assembly source and a filename (used in error
// messages) and returns a 143,360-byte DOS 3.3 disk image.
func Assemble(src []byte, filename string) ([]byte, error) {
	a := &assembler{
		filename: filename,
		origin:   defaultOrigin,
		labels:   make(map[string]uint16),
	}
	return a.run(src)
}

type assembler struct {
	filename string
	origin   uint16
	labels   map[string]uint16
}

// node holds information about one source line, populated across both passes.
type node struct {
	lineNum int
	label   string      // defined label (without colon), or ""
	mnem    string      // uppercase mnemonic, or ""
	oi      operandInfo // parsed operand
	dir     string      // directive name (lowercase), or ""
	dirArgs []string    // directive arguments
	pc      uint16      // address assigned during pass 1
	size    int         // byte count this line emits
	mode    int         // resolved addressing mode (instructions only)
}

func (a *assembler) errorf(lineNum int, format string, args ...any) error {
	return fmt.Errorf("%s:%d: "+format, append([]any{a.filename, lineNum}, args...)...)
}

func (a *assembler) run(src []byte) ([]byte, error) {
	rawLines := strings.Split(strings.ReplaceAll(string(src), "\r\n", "\n"), "\n")

	nodes, err := a.parseAll(rawLines)
	if err != nil {
		return nil, err
	}

	if err := a.pass1(nodes); err != nil {
		return nil, err
	}

	code, err := a.pass2(nodes)
	if err != nil {
		return nil, err
	}

	if len(code) > maxCodeSize {
		return nil, fmt.Errorf("assembled code exceeds track 0 capacity: %d bytes (max %d)", len(code), maxCodeSize)
	}

	return buildDiskImage(code), nil
}

func (a *assembler) parseAll(rawLines []string) ([]*node, error) {
	nodes := make([]*node, 0, len(rawLines))

	for i, raw := range rawLines {
		lineNum := i + 1
		p, err := parseLine(raw)
		if err != nil {
			return nil, a.errorf(lineNum, "%v", err)
		}

		n := &node{
			lineNum: lineNum,
			label:   p.label,
			mnem:    p.mnem,
			dir:     p.dir,
			dirArgs: p.dirArgs,
		}

		if p.mnem != "" {
			n.oi, err = parseOperand(p.operand)
			if err != nil {
				return nil, a.errorf(lineNum, "%v", err)
			}
		}

		nodes = append(nodes, n)
	}

	return nodes, nil
}

func (a *assembler) pass1(nodes []*node) error {
	pc := a.origin

	for _, n := range nodes {
		n.pc = pc

		if n.label != "" {
			if _, exists := a.labels[n.label]; exists {
				return a.errorf(n.lineNum, "duplicate label %q", n.label)
			}
			a.labels[n.label] = pc
		}

		if n.dir != "" {
			size, newPC, err := a.directiveInfo(n, pc)
			if err != nil {
				return err
			}
			n.size = size
			pc = newPC
			continue
		}

		if n.mnem != "" {
			mode, ok := resolveMode(n.mnem, n.oi)
			if !ok {
				return a.errorf(n.lineNum, "no valid addressing mode for %s %v", n.mnem, n.oi)
			}
			n.mode = mode
			n.size = 1 + operandSize[mode]
			pc += uint16(n.size)
		}
	}

	return nil
}

// directiveInfo returns (size, newPC) for a directive, applying any PC
// change.
func (a *assembler) directiveInfo(n *node, pc uint16) (size int, newPC uint16, err error) {
	switch n.dir {
	case "org":
		if len(n.dirArgs) != 1 {
			return 0, 0, a.errorf(n.lineNum, ".org requires exactly one argument")
		}
		v, err := parseHexValue(n.dirArgs[0])
		if err != nil {
			return 0, 0, a.errorf(n.lineNum, "invalid .org argument: %v", err)
		}
		return 0, uint16(v), nil

	case "byte":
		if len(n.dirArgs) == 0 {
			return 0, 0, a.errorf(n.lineNum, ".byte requires at least one argument")
		}
		for _, arg := range n.dirArgs {
			v, err := parseHexValue(arg)
			if err != nil {
				return 0, 0, a.errorf(n.lineNum, "invalid .byte argument %q: %v", arg, err)
			}
			if v > 0xFF {
				return 0, 0, a.errorf(n.lineNum, ".byte value out of range: %s", arg)
			}
		}
		return len(n.dirArgs), pc + uint16(len(n.dirArgs)), nil

	case "word":
		if len(n.dirArgs) == 0 {
			return 0, 0, a.errorf(n.lineNum, ".word requires at least one argument")
		}
		for _, arg := range n.dirArgs {
			if _, err := parseHexValue(arg); err != nil {
				return 0, 0, a.errorf(n.lineNum, "invalid .word argument %q: %v", arg, err)
			}
		}
		sz := len(n.dirArgs) * 2
		return sz, pc + uint16(sz), nil

	case "halt":
		return 3, pc + 3, nil // JMP * emits 3 bytes

	default:
		return 0, 0, a.errorf(n.lineNum, "unknown directive .%s", n.dir)
	}
}

func (a *assembler) pass2(nodes []*node) ([]byte, error) {
	total := 0
	for _, n := range nodes {
		total += n.size
	}
	code := make([]byte, 0, total)

	for _, n := range nodes {
		bytes, err := a.emitNode(n)
		if err != nil {
			return nil, err
		}
		code = append(code, bytes...)
	}

	return code, nil
}

func (a *assembler) emitNode(n *node) ([]byte, error) {
	if n.dir != "" {
		return a.emitDirective(n)
	}
	if n.mnem != "" {
		return a.emitInstruction(n)
	}
	return nil, nil
}

func (a *assembler) emitDirective(n *node) ([]byte, error) {
	switch n.dir {
	case "org":
		return nil, nil

	case "byte":
		out := make([]byte, len(n.dirArgs))
		for i, arg := range n.dirArgs {
			v, _ := parseHexValue(arg)
			out[i] = byte(v)
		}
		return out, nil

	case "word":
		out := make([]byte, len(n.dirArgs)*2)
		for i, arg := range n.dirArgs {
			v, _ := parseHexValue(arg)
			out[i*2] = byte(v & 0xFF)
			out[i*2+1] = byte(v >> 8)
		}
		return out, nil

	case "halt":
		opcode, _ := lookupOpcode("JMP", modeABS)
		pc := n.pc
		return []byte{opcode, byte(pc & 0xFF), byte(pc >> 8)}, nil

	default:
		return nil, a.errorf(n.lineNum, "unknown directive .%s", n.dir)
	}
}

func (a *assembler) emitInstruction(n *node) ([]byte, error) {
	opcode, _ := lookupOpcode(n.mnem, n.mode)
	out := []byte{opcode}

	if n.oi.isLabel {
		addr, ok := a.labels[n.oi.label]
		if !ok {
			return nil, a.errorf(n.lineNum, "undefined label %q", n.oi.label)
		}
		switch n.mode {
		case modeREL:
			next := int(n.pc) + 2
			offset := int(addr) - next
			if offset < -128 || offset > 127 {
				return nil, a.errorf(n.lineNum, "branch target out of range: offset %d", offset)
			}
			out = append(out, byte(int8(offset)))
		case modeABS:
			out = append(out, byte(addr&0xFF), byte(addr>>8))
		default:
			return nil, a.errorf(n.lineNum, "unexpected mode for label operand")
		}
	} else {
		v := n.oi.value
		switch operandSize[n.mode] {
		case 1:
			out = append(out, byte(v))
		case 2:
			out = append(out, byte(v&0xFF), byte(v>>8))
		}
	}

	return out, nil
}

// resolveMode returns the first mode from oi.modes that the instruction
// supports.
func resolveMode(mnem string, oi operandInfo) (int, bool) {
	for _, m := range oi.modes {
		if _, ok := lookupOpcode(mnem, m); ok {
			return m, true
		}
	}
	return 0, false
}

// buildDiskImage wraps assembled code bytes into a 143,360-byte DOS 3.3 .dsk
// image. The code is placed in track 0 with the correct sector interleave so
// that the Apple II boot ROM loads it contiguously starting at $0800.
func buildDiskImage(code []byte) []byte {
	// Sector count: number of 256-byte sectors needed to hold 1 (sector count
	// byte) + len(code) bytes.
	total := len(code) + 1
	sectorCount := min((total+a2enc.LogSectorLen-1)/a2enc.LogSectorLen, 15)

	// Build a flat 4096-byte track buffer in physical/memory order: flat[0] =
	// sector count byte (loaded to $0800) flat[1..] = assembled code (loaded
	// to $0801 onward)
	flat := make([]byte, a2enc.LogTrackLen)
	flat[0] = byte(sectorCount)
	copy(flat[1:], code)

	// Distribute the flat track sectors into the .dsk logical sectors using
	// the DOS 3.3 physical-to-logical interleave. When erc loads the .dsk it
	// calls a2enc.Encode, which reads logical sector LogicalSector(DOS33,
	// physSect) for each physical sector -- so placing
	// flat[physSect*256:(physSect+1)*256] at logical sector
	// LogicalSector(DOS33, physSect) ensures the data is fetched in the right
	// order and loaded contiguously to memory.
	dsk := make([]byte, a2enc.DosSize)
	for physSect := range a2enc.NumSectors {
		logSect := a2enc.LogicalSector(a2enc.DOS33, physSect)
		srcOff := physSect * a2enc.LogSectorLen
		dstOff := logSect * a2enc.LogSectorLen
		copy(dsk[dstOff:dstOff+a2enc.LogSectorLen], flat[srcOff:srcOff+a2enc.LogSectorLen])
	}

	return dsk
}
