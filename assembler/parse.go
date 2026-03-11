package assembler

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// parsedLine holds the decoded fields of one source line.
type parsedLine struct {
	label   string   // label name (without colon), or ""
	mnem    string   // uppercase mnemonic, or ""
	operand string   // raw operand text (trimmed), or ""
	dir     string   // directive name (without dot, lowercase), or ""
	dirArgs []string // directive arguments
}

var (
	labelRe = regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_]*):\s*`)
	identRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
)

func parseLine(raw string) (*parsedLine, error) {
	p := &parsedLine{}

	// Strip comment.
	if idx := strings.IndexByte(raw, ';'); idx >= 0 {
		raw = raw[:idx]
	}
	line := strings.TrimSpace(raw)

	if line == "" {
		return p, nil
	}

	// Check for a label prefix.
	if m := labelRe.FindStringSubmatchIndex(line); m != nil {
		p.label = line[m[2]:m[3]]
		line = strings.TrimSpace(line[m[1]:])
	}

	if line == "" {
		return p, nil
	}

	// Directive.
	if line[0] == '.' {
		rest := strings.TrimSpace(line[1:])
		if rest == "" {
			return nil, fmt.Errorf("empty directive")
		}
		fields := strings.Fields(rest)
		p.dir = strings.ToLower(fields[0])
		if len(fields) > 1 {
			argStr := strings.Join(fields[1:], " ")
			for arg := range strings.SplitSeq(argStr, ",") {
				arg = strings.TrimSpace(arg)
				if arg != "" {
					p.dirArgs = append(p.dirArgs, arg)
				}
			}
		}
		return p, nil
	}

	// Instruction: first token is the mnemonic, rest is the operand.
	idx := strings.IndexAny(line, " \t")
	var mnemRaw, operand string
	if idx < 0 {
		mnemRaw = line
	} else {
		mnemRaw = line[:idx]
		operand = strings.TrimSpace(line[idx+1:])
		// Strip any trailing comment from the operand.
		if ci := strings.IndexByte(operand, ';'); ci >= 0 {
			operand = strings.TrimSpace(operand[:ci])
		}
	}

	mnem := strings.ToUpper(mnemRaw)
	if !isValidMnem(mnem) {
		return nil, fmt.Errorf("unknown mnemonic %q", mnem)
	}

	p.mnem = mnem
	p.operand = operand
	return p, nil
}

// operandInfo describes a parsed operand and the addressing modes to try.
type operandInfo struct {
	modes   []int  // modes to try in preference order
	value   uint32 // numeric operand value (for non-label operands)
	isLabel bool
	label   string
}

func parseOperand(operand string) (operandInfo, error) {
	operand = strings.TrimSpace(operand)

	// No operand: try ACC then IMP (handles both accumulator-implicit and
	// truly implied instructions).
	if operand == "" {
		return operandInfo{modes: []int{modeACC, modeIMP}}, nil
	}

	// Explicit accumulator.
	if strings.EqualFold(operand, "A") {
		return operandInfo{modes: []int{modeACC}}, nil
	}

	// Immediate: #$NN
	if strings.HasPrefix(operand, "#") {
		v, err := parseHexValue(operand[1:])
		if err != nil {
			return operandInfo{}, fmt.Errorf("invalid immediate operand: %v", err)
		}
		if v > 0xFF {
			return operandInfo{}, fmt.Errorf("immediate value out of range: %s", operand[1:])
		}
		return operandInfo{modes: []int{modeIMM}, value: v}, nil
	}

	// Indirect and indexed-indirect modes: operand starts with '('.
	if strings.HasPrefix(operand, "(") {
		return parseIndirectOperand(operand)
	}

	// Indexed: contains a comma (e.g. $NN,X or $NNNN,Y).
	if strings.Contains(operand, ",") {
		return parseIndexedOperand(operand)
	}

	// Label reference (identifier without colon).
	if identRe.MatchString(operand) {
		return operandInfo{isLabel: true, label: operand, modes: []int{modeREL, modeABS}}, nil
	}

	// Plain address: $NN or $NNNN.
	return parsePlainAddr(operand)
}

func parseIndirectOperand(operand string) (operandInfo, error) {
	upper := strings.ToUpper(operand)

	// ($NN),Y -- indirect Y-indexed
	if strings.HasSuffix(upper, "),Y") {
		inner := operand[:len(operand)-2] // strip ,Y
		inner = strings.TrimSuffix(inner, ")")
		if !strings.HasPrefix(inner, "(") {
			return operandInfo{}, fmt.Errorf("invalid indirect Y operand: %s", operand)
		}
		inner = inner[1:]
		v, err := parseHexValue(inner)
		if err != nil {
			return operandInfo{}, fmt.Errorf("invalid indirect Y operand: %v", err)
		}
		if v > 0xFF {
			return operandInfo{}, fmt.Errorf("indirect Y operand out of range: %s", inner)
		}
		return operandInfo{modes: []int{modeIDY}, value: v}, nil
	}

	if !strings.HasPrefix(operand, "(") || !strings.HasSuffix(operand, ")") {
		return operandInfo{}, fmt.Errorf("invalid indirect operand: %s", operand)
	}
	inner := operand[1 : len(operand)-1]

	// ($NN,X) -- x-indexed indirect
	if strings.HasSuffix(strings.ToUpper(inner), ",X") {
		addrStr := inner[:len(inner)-2]
		v, err := parseHexValue(addrStr)
		if err != nil {
			return operandInfo{}, fmt.Errorf("invalid x-indexed indirect operand: %v", err)
		}
		if v > 0xFF {
			return operandInfo{}, fmt.Errorf("x-indexed indirect operand out of range: %s", addrStr)
		}
		return operandInfo{modes: []int{modeIDX}, value: v}, nil
	}

	v, err := parseHexValue(inner)
	if err != nil {
		return operandInfo{}, fmt.Errorf("invalid indirect operand: %v", err)
	}

	if isHex4(inner) || v > 0xFF {
		// ($NNNN) -- absolute indirect (JMP only)
		return operandInfo{modes: []int{modeIND}, value: v}, nil
	}

	// ($NN) -- zero page indirect (65C02)
	return operandInfo{modes: []int{modeZPI}, value: v}, nil
}

func parseIndexedOperand(operand string) (operandInfo, error) {
	idx := strings.LastIndex(operand, ",")
	if idx < 0 {
		return operandInfo{}, fmt.Errorf("invalid indexed operand: %s", operand)
	}
	addrStr := strings.TrimSpace(operand[:idx])
	regStr := strings.ToUpper(strings.TrimSpace(operand[idx+1:]))

	v, err := parseHexValue(addrStr)
	if err != nil {
		return operandInfo{}, fmt.Errorf("invalid indexed operand: %v", err)
	}

	forced4 := isHex4(addrStr)

	switch regStr {
	case "X":
		if !forced4 && v <= 0xFF {
			return operandInfo{modes: []int{modeZPX, modeABX}, value: v}, nil
		}
		return operandInfo{modes: []int{modeABX}, value: v}, nil
	case "Y":
		if !forced4 && v <= 0xFF {
			return operandInfo{modes: []int{modeZPY, modeABY}, value: v}, nil
		}
		return operandInfo{modes: []int{modeABY}, value: v}, nil
	default:
		return operandInfo{}, fmt.Errorf("unknown index register %q in operand %s", regStr, operand)
	}
}

func parsePlainAddr(operand string) (operandInfo, error) {
	v, err := parseHexValue(operand)
	if err != nil {
		return operandInfo{}, fmt.Errorf("invalid operand %q: %v", operand, err)
	}

	if isHex4(operand) || v > 0xFF {
		return operandInfo{modes: []int{modeABS}, value: v}, nil
	}
	return operandInfo{modes: []int{modeZPG, modeABS}, value: v}, nil
}

func parseHexValue(s string) (uint32, error) {
	if !strings.HasPrefix(s, "$") {
		return 0, fmt.Errorf("expected $ prefix, got %q", s)
	}
	v, err := strconv.ParseUint(s[1:], 16, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid hex value %q", s)
	}
	return uint32(v), nil
}

// isHex4 returns true if s is a $-prefixed 4-digit hex literal (forces
// absolute mode).
func isHex4(s string) bool {
	return strings.HasPrefix(s, "$") && len(s[1:]) == 4
}
