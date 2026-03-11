package assembler

import "github.com/pevans/erc/mos"

// Addressing mode constants used internally by the assembler.
const (
	modeIMP = iota + 1 // implied
	modeACC            // accumulator
	modeIMM            // immediate
	modeABS            // absolute
	modeABX            // absolute x-indexed
	modeABY            // absolute y-indexed
	modeZPG            // zero page
	modeZPX            // zero page x-indexed
	modeZPY            // zero page y-indexed
	modeIND            // indirect
	modeIDX            // x-indexed indirect
	modeIDY            // indirect y-indexed
	modeREL            // relative
	modeZPI            // zero page indirect (65C02)
)

// operandSize maps each mode to its operand byte count (not counting the
// opcode byte). Indexed by the mode constants (modeIMP=1 .. modeZPI=14).
var operandSize = [15]int{
	0,    // 0: unused
	0, 0, // modeIMP, modeACC
	1,       // modeIMM
	2, 2, 2, // modeABS, modeABX, modeABY
	1, 1, 1, // modeZPG, modeZPX, modeZPY
	2,    // modeIND
	1, 1, // modeIDX, modeIDY
	1, // modeREL
	1, // modeZPI
}

// mosToAsm maps mos address mode constants to assembler mode constants.
var mosToAsm = map[int]int{
	mos.AmIMP: modeIMP,
	mos.AmACC: modeACC,
	mos.AmIMM: modeIMM,
	mos.AmABS: modeABS,
	mos.AmABX: modeABX,
	mos.AmABY: modeABY,
	mos.AmZPG: modeZPG,
	mos.AmZPX: modeZPX,
	mos.AmZPY: modeZPY,
	mos.AmIND: modeIND,
	mos.AmIDX: modeIDX,
	mos.AmIDY: modeIDY,
	mos.AmREL: modeREL,
}

type opcodeKey struct {
	mnem string
	mode int
}

var (
	opcodeTable    map[opcodeKey]byte
	validMnemonics map[string]struct{}
)

// mnemNormalize maps mos-internal names to standard 65C02 assembler
// mnemonics.
var mnemNormalize = map[string]string{
	"BIM": "BIT", // $89 is BIT immediate; mos names it BIM internally
}

// zpiTable maps mnemonics to their 65C02 zero-page-indirect opcodes. The mos
// package maps these opcodes to modeZPG (a simplification), so we maintain
// them separately to support the ($NN) syntax.
var zpiTable = map[string]byte{
	"ORA": 0x12,
	"AND": 0x32,
	"EOR": 0x52,
	"ADC": 0x72,
	"STA": 0x92,
	"LDA": 0xB2,
	"CMP": 0xD2,
	"SBC": 0xF2,
}

func init() {
	opcodeTable = buildOpcodeTable()
	validMnemonics = make(map[string]struct{}, len(opcodeTable))
	for key := range opcodeTable {
		validMnemonics[key.mnem] = struct{}{}
	}
}

func buildOpcodeTable() map[opcodeKey]byte {
	table := make(map[opcodeKey]byte)

	for i := range 256 {
		opcode := uint8(i)
		name := mos.OpcodeInstructionName(opcode)

		// Skip placeholder/undefined/filler opcodes.
		if name == "NP2" || name == "NP3" || name == "NOP" {
			continue
		}

		if norm, ok := mnemNormalize[name]; ok {
			name = norm
		}

		asmMode, ok := mosToAsm[mos.OpcodeAddrMode(opcode)]
		if !ok {
			continue // skip BY2/BY3 placeholder modes
		}

		key := opcodeKey{name, asmMode}
		if _, exists := table[key]; !exists {
			table[key] = opcode
		}
	}

	// NOP is always $EA (implied); earlier opcodes labeled NOP are undefined.
	table[opcodeKey{"NOP", modeIMP}] = 0xEA

	// ZPI (zero-page indirect) entries use ($NN) syntax.
	for mnem, opcode := range zpiTable {
		table[opcodeKey{mnem, modeZPI}] = opcode
	}

	return table
}

func lookupOpcode(mnem string, mode int) (byte, bool) {
	op, ok := opcodeTable[opcodeKey{mnem, mode}]
	return op, ok
}

func isValidMnem(mnem string) bool {
	_, ok := validMnemonics[mnem]
	return ok
}
