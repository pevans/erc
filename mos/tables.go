package mos

// Below is a table of instructions that are mapped to opcodes. For
// corresponding address modes, see addr.go.
//
//	00   01   02   03   04   05   06   07   08   09   0A   0B   0C   0D   0E   0F
var instructions = [256]Instruction{
	Brk, Ora, Np2, Nop, Tsb, Ora, Asl, Nop, Php, Ora, Asl, Nop, Tsb, Ora, Asl, Nop, // 0x
	Bpl, Ora, Ora, Nop, Trb, Ora, Asl, Nop, Clc, Ora, Inc, Nop, Trb, Ora, Asl, Nop, // 1x
	Jsr, And, Np2, Nop, Bit, And, Rol, Nop, Plp, And, Rol, Nop, Bit, And, Rol, Nop, // 2x
	Bmi, And, And, Nop, Bit, And, Rol, Nop, Sec, And, Dec, Nop, Bit, And, Rol, Nop, // 3x
	Rti, Eor, Np2, Nop, Np2, Eor, Lsr, Nop, Pha, Eor, Lsr, Nop, Jmp, Eor, Lsr, Nop, // 4x
	Bvc, Eor, Eor, Nop, Np2, Eor, Lsr, Nop, Cli, Eor, Phy, Nop, Np3, Eor, Lsr, Nop, // 5x
	Rts, Adc, Np2, Nop, Stz, Adc, Ror, Nop, Pla, Adc, Ror, Nop, Jmp, Adc, Ror, Nop, // 6x
	Bvs, Adc, Adc, Nop, Stz, Adc, Ror, Nop, Sei, Adc, Ply, Nop, Jmp, Adc, Ror, Nop, // 7x
	Bra, Sta, Np2, Nop, Sty, Sta, Stx, Nop, Dey, Bim, Txa, Nop, Sty, Sta, Stx, Nop, // 8x
	Bcc, Sta, Sta, Nop, Sty, Sta, Stx, Nop, Tya, Sta, Txs, Nop, Stz, Sta, Stz, Nop, // 9x
	Ldy, Lda, Ldx, Nop, Ldy, Lda, Ldx, Nop, Tay, Lda, Tax, Nop, Ldy, Lda, Ldx, Nop, // Ax
	Bcs, Lda, Lda, Nop, Ldy, Lda, Ldx, Nop, Clv, Lda, Tsx, Nop, Ldy, Lda, Ldx, Nop, // Bx
	Cpy, Cmp, Np2, Nop, Cpy, Cmp, Dec, Nop, Iny, Cmp, Dex, Nop, Cpy, Cmp, Dec, Nop, // Cx
	Bne, Cmp, Cmp, Nop, Np2, Cmp, Dec, Nop, Cld, Cmp, Phx, Nop, Np3, Cmp, Dec, Nop, // Dx
	Cpx, Sbc, Np2, Nop, Cpx, Sbc, Inc, Nop, Inx, Sbc, Nop, Nop, Cpx, Sbc, Inc, Nop, // Ex
	Beq, Sbc, Sbc, Nop, Np2, Sbc, Inc, Nop, Sed, Sbc, Plx, Nop, Np3, Sbc, Inc, Nop, // Fx
}

// The names of the instructions according to their opcodes
var opcodeNames = [256]string{
	"BRK", "ORA", "NP2", "NOP", "TSB", "ORA", "ASL", "NOP", "PHP", "ORA", "ASL", "NOP", "TSB", "ORA", "ASL", "NOP", // 0x
	"BPL", "ORA", "ORA", "NOP", "TRB", "ORA", "ASL", "NOP", "CLC", "ORA", "INC", "NOP", "TRB", "ORA", "ASL", "NOP", // 1x
	"JSR", "AND", "NP2", "NOP", "BIT", "AND", "ROL", "NOP", "PLP", "AND", "ROL", "NOP", "BIT", "AND", "ROL", "NOP", // 2x
	"BMI", "AND", "AND", "NOP", "BIT", "AND", "ROL", "NOP", "SEC", "AND", "DEC", "NOP", "BIT", "AND", "ROL", "NOP", // 3x
	"RTI", "EOR", "NP2", "NOP", "NP2", "EOR", "LSR", "NOP", "PHA", "EOR", "LSR", "NOP", "JMP", "EOR", "LSR", "NOP", // 4x
	"BVC", "EOR", "EOR", "NOP", "NP2", "EOR", "LSR", "NOP", "CLI", "EOR", "PHY", "NOP", "NP3", "EOR", "LSR", "NOP", // 5x
	"RTS", "ADC", "NP2", "NOP", "STZ", "ADC", "ROR", "NOP", "PLA", "ADC", "ROR", "NOP", "JMP", "ADC", "ROR", "NOP", // 6x
	"BVS", "ADC", "ADC", "NOP", "STZ", "ADC", "ROR", "NOP", "SEI", "ADC", "PLY", "NOP", "JMP", "ADC", "ROR", "NOP", // 7x
	"BRA", "STA", "NP2", "NOP", "STY", "STA", "STX", "NOP", "DEY", "BIM", "TXA", "NOP", "STY", "STA", "STX", "NOP", // 8x
	"BCC", "STA", "STA", "NOP", "STY", "STA", "STX", "NOP", "TYA", "STA", "TXS", "NOP", "STZ", "STA", "STZ", "NOP", // 9x
	"LDY", "LDA", "LDX", "NOP", "LDY", "LDA", "LDX", "NOP", "TAY", "LDA", "TAX", "NOP", "LDY", "LDA", "LDX", "NOP", // Ax
	"BCS", "LDA", "LDA", "NOP", "LDY", "LDA", "LDX", "NOP", "CLV", "LDA", "TSX", "NOP", "LDY", "LDA", "LDX", "NOP", // Bx
	"CPY", "CMP", "NP2", "NOP", "CPY", "CMP", "DEC", "NOP", "INY", "CMP", "DEX", "NOP", "CPY", "CMP", "DEC", "NOP", // Cx
	"BNE", "CMP", "CMP", "NOP", "NP2", "CMP", "DEC", "NOP", "CLD", "CMP", "PHX", "NOP", "NP3", "CMP", "DEC", "NOP", // Dx
	"CPX", "SBC", "NP2", "NOP", "CPX", "SBC", "INC", "NOP", "INX", "SBC", "NOP", "NOP", "CPX", "SBC", "INC", "NOP", // Ex
	"BEQ", "SBC", "SBC", "NOP", "NP2", "SBC", "INC", "NOP", "SED", "SBC", "PLX", "NOP", "NP3", "SBC", "INC", "NOP", // Fx
}

// 0  1  2  3  4  5  6  7  8  9  A  B  C  D  E  F
var cycles = [256]uint8{
	7, 6, 2, 1, 5, 3, 5, 1, 3, 2, 2, 1, 6, 4, 6, 1, // 0x
	2, 5, 5, 1, 5, 4, 6, 1, 2, 4, 2, 1, 6, 4, 6, 1, // 1x
	6, 6, 2, 1, 3, 3, 5, 1, 4, 2, 2, 1, 4, 4, 6, 1, // 2x
	2, 5, 5, 1, 4, 4, 6, 1, 2, 4, 2, 1, 4, 4, 6, 1, // 3x
	6, 6, 2, 1, 3, 3, 5, 1, 3, 2, 2, 1, 3, 4, 6, 1, // 4x
	2, 5, 5, 1, 4, 4, 6, 1, 2, 4, 3, 1, 8, 4, 6, 1, // 5x
	6, 6, 2, 1, 3, 3, 5, 1, 4, 2, 2, 1, 6, 4, 6, 1, // 6x
	2, 5, 5, 1, 4, 4, 6, 1, 2, 4, 4, 1, 6, 4, 6, 1, // 7x
	2, 6, 2, 1, 3, 3, 3, 1, 2, 2, 2, 1, 4, 4, 4, 1, // 8x
	2, 6, 5, 1, 4, 4, 4, 1, 2, 5, 2, 1, 4, 5, 5, 1, // 9x
	2, 6, 2, 1, 3, 3, 3, 1, 2, 2, 2, 1, 4, 4, 4, 1, // Ax
	2, 5, 5, 1, 4, 4, 4, 1, 2, 4, 2, 1, 4, 4, 4, 1, // Bx
	2, 6, 2, 1, 3, 3, 5, 1, 2, 2, 2, 1, 4, 4, 6, 1, // Cx
	2, 5, 5, 1, 4, 4, 6, 1, 2, 4, 3, 1, 4, 4, 7, 1, // Dx
	2, 6, 2, 1, 3, 3, 5, 1, 2, 2, 2, 1, 4, 4, 6, 1, // Ex
	2, 5, 5, 1, 4, 4, 6, 1, 2, 4, 4, 1, 4, 4, 7, 1, // Fx
}

// Below is an address mode table that maps mode functions to specific
// opcodes.
//
//	00   01   02   03   04   05   06   07   08   09   0A   0B   0C   0D   0E   0F
var addrModeFuncs = [256]AddrMode{
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

// Like the above table, only it maps opcodes to the symbolic constants.
//
//	00     01     02     03     04     05     06     07     08     09     0A     0B     0C     0D     0E     0F
var addrModes = [256]int{
	AmIMP, AmIDX, AmBY2, AmIMP, AmZPG, AmZPG, AmZPG, AmIMP, AmIMP, AmIMM, AmACC, AmIMP, AmABS, AmABS, AmABS, AmIMP, // 0x
	AmREL, AmIDY, AmZPG, AmIMP, AmZPG, AmZPX, AmZPX, AmIMP, AmIMP, AmABY, AmACC, AmIMP, AmABS, AmABX, AmABX, AmIMP, // 1x
	AmABS, AmIDX, AmBY2, AmIMP, AmZPG, AmZPG, AmZPG, AmIMP, AmIMP, AmIMM, AmACC, AmIMP, AmABS, AmABS, AmABS, AmIMP, // 2x
	AmREL, AmIDY, AmZPG, AmIMP, AmZPX, AmZPX, AmZPX, AmIMP, AmIMP, AmABY, AmACC, AmIMP, AmABX, AmABX, AmABX, AmIMP, // 3x
	AmIMP, AmIDX, AmBY2, AmIMP, AmBY2, AmZPG, AmZPG, AmIMP, AmIMP, AmIMM, AmACC, AmIMP, AmABS, AmABS, AmABS, AmIMP, // 4x
	AmREL, AmIDY, AmZPG, AmIMP, AmBY2, AmZPX, AmZPX, AmIMP, AmIMP, AmABY, AmIMP, AmIMP, AmBY3, AmABX, AmABX, AmIMP, // 5x
	AmIMP, AmIDX, AmBY2, AmIMP, AmZPG, AmZPG, AmZPG, AmIMP, AmIMP, AmIMM, AmACC, AmIMP, AmIND, AmABS, AmABS, AmIMP, // 6x
	AmREL, AmIDY, AmZPG, AmIMP, AmZPX, AmZPX, AmZPX, AmIMP, AmIMP, AmABY, AmIMP, AmIMP, AmABX, AmABX, AmABX, AmIMP, // 7x
	AmREL, AmIDX, AmBY2, AmIMP, AmZPG, AmZPG, AmZPG, AmIMP, AmIMP, AmIMM, AmIMP, AmIMP, AmABS, AmABS, AmABS, AmIMP, // 8x
	AmREL, AmIDY, AmZPG, AmIMP, AmZPX, AmZPX, AmZPY, AmIMP, AmIMP, AmABY, AmIMP, AmIMP, AmABS, AmABX, AmABX, AmIMP, // 9x
	AmIMM, AmIDX, AmIMM, AmIMP, AmZPG, AmZPG, AmZPG, AmIMP, AmIMP, AmIMM, AmIMP, AmIMP, AmABS, AmABS, AmABS, AmIMP, // Ax
	AmREL, AmIDY, AmZPG, AmIMP, AmZPX, AmZPX, AmZPY, AmIMP, AmIMP, AmABY, AmIMP, AmIMP, AmABX, AmABX, AmABY, AmIMP, // Bx
	AmIMM, AmIDX, AmBY2, AmIMP, AmZPG, AmZPG, AmZPG, AmIMP, AmIMP, AmIMM, AmIMP, AmIMP, AmABS, AmABS, AmABS, AmIMP, // Cx
	AmREL, AmIDY, AmZPG, AmIMP, AmBY2, AmZPX, AmZPX, AmIMP, AmIMP, AmABY, AmIMP, AmIMP, AmBY3, AmABX, AmABX, AmIMP, // Dx
	AmIMM, AmIDX, AmBY2, AmIMP, AmZPG, AmZPG, AmZPG, AmIMP, AmIMP, AmIMM, AmIMP, AmIMP, AmABS, AmABS, AmABS, AmIMP, // Ex
	AmREL, AmIDY, AmZPG, AmIMP, AmBY2, AmZPX, AmZPX, AmIMP, AmIMP, AmABY, AmIMP, AmIMP, AmBY3, AmABX, AmABX, AmIMP, // Fx
}

// The offsets table defines the number of bytes we must increment the
// PC register after a given instruction. The bytes vary based on
// address mode, rather than the specific instruction. In cases where
// the instruction would change the PC due to its defined behavior, the
// offset is given as zero.
//
//	0  1  2  3  4  5  6  7  8  9  A  B  C  D  E  F
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
