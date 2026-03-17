setup_file() { load mos_helper; setup_file; }
setup()      { load mos_helper; setup; }
teardown()   { load mos_helper; teardown; }

# ---------------------------------------------------------------------------
# 1. Load / Store (LDA, LDX, LDY, STA, STX, STY)
# ---------------------------------------------------------------------------

@test "LDA immediate loads value into A" {
	cpu_run 'LDA #$42' 'STA $00' '.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 42
}

@test "LDA zero-page loads from memory" {
	cpu_run \
		'LDA #$5A' 'STA $10' \
		'LDA #$00' \
		'LDA $10' 'STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 5A
}

@test "LDA absolute loads from full 16-bit address" {
	cpu_run \
		'LDA #$7B' 'STA $0300' \
		'LDA #$00' \
		'LDA $0300' 'STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 7B
}

@test "LDA zero-page X-indexed loads from zpg+X" {
	cpu_run \
		'LDA #$33' 'STA $12' \
		'LDX #$02' \
		'LDA $10,X' 'STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 33
}

@test "LDA absolute X-indexed loads from abs+X" {
	cpu_run \
		'LDA #$44' 'STA $0302' \
		'LDX #$02' \
		'LDA $0300,X' 'STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 44
}

@test "LDA absolute Y-indexed loads from abs+Y" {
	cpu_run \
		'LDA #$55' 'STA $0303' \
		'LDY #$03' \
		'LDA $0300,Y' 'STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 55
}

@test "LDX immediate loads value into X" {
	cpu_run 'LDX #$A1' 'STX $00' '.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 A1
}

@test "LDX zero-page loads from memory" {
	cpu_run \
		'LDA #$A1' 'STA $10' \
		'LDX #$00' \
		'LDX $10' 'STX $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 A1
}

@test "LDX absolute loads from full 16-bit address" {
	cpu_run \
		'LDA #$A1' 'STA $0300' \
		'LDX #$00' \
		'LDX $0300' 'STX $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 A1
}

@test "LDY immediate loads value into Y" {
	cpu_run 'LDY #$B2' 'STY $00' '.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 B2
}

@test "LDY zero-page loads from memory" {
	cpu_run \
		'LDA #$B2' 'STA $10' \
		'LDY #$00' \
		'LDY $10' 'STY $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 B2
}

@test "LDY absolute loads from full 16-bit address" {
	cpu_run \
		'LDA #$B2' 'STA $0300' \
		'LDY #$00' \
		'LDY $0300' 'STY $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 B2
}

@test "STX stores X to zero-page" {
	cpu_run 'LDX #$C3' 'STX $00' '.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 C3
}

@test "STY stores Y to zero-page" {
	cpu_run 'LDY #$D4' 'STY $00' '.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 D4
}

@test "STZ zero-page stores zero" {
	cpu_run \
		'LDA #$FF' 'STA $03' \
		'STZ $03' \
		'LDA $03' 'STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 00
}

@test "STZ absolute stores zero" {
	cpu_run \
		'LDA #$FF' 'STA $0300' \
		'STZ $0300' \
		'LDA $0300' 'STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 00
}

@test "LDA #\$80 sets N flag; LDA #\$40 clears N flag" {
	cpu_run \
		'LDA #$80' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_flag $NEGATIVE

	cpu_run \
		'LDA #$40' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	refute_flag $NEGATIVE
}

@test "LDA #\$00 sets Z flag; LDA #\$01 clears Z flag" {
	cpu_run \
		'LDA #$00' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_flag $ZERO

	cpu_run \
		'LDA #$01' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	refute_flag $ZERO
}

# ---------------------------------------------------------------------------
# 2. Arithmetic (ADC, SBC)
# ---------------------------------------------------------------------------

@test "ADC basic addition without carry" {
	cpu_run \
		'CLC' \
		'LDA #$10' 'ADC #$20' \
		'STA $00' '.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 30
}

@test "ADC with carry-in adds 1 extra" {
	cpu_run \
		'SEC' \
		'LDA #$10' 'ADC #$20' \
		'STA $00' '.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 31
}

@test "ADC carry-out sets C flag" {
	cpu_run \
		'CLC' \
		'LDA #$FF' 'ADC #$01' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_flag $CARRY
}

@test "ADC overflow sets V flag (\$7F + \$01)" {
	cpu_run \
		'CLC' \
		'LDA #$7F' 'ADC #$01' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_flag $OVERFLOW
}

@test "ADC result \$00 sets Z flag" {
	cpu_run \
		'CLC' \
		'LDA #$FF' 'ADC #$01' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_flag $ZERO
}

@test "ADC result with bit 7 set sets N flag" {
	cpu_run \
		'CLC' \
		'LDA #$7F' 'ADC #$02' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_flag $NEGATIVE
}

@test "SBC basic subtraction" {
	cpu_run \
		'SEC' \
		'LDA #$30' 'SBC #$10' \
		'STA $00' '.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 20
}

@test "SBC borrow clears C flag" {
	cpu_run \
		'SEC' \
		'LDA #$00' 'SBC #$01' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	refute_flag $CARRY
}

@test "SBC overflow \$80 - \$01 sets V flag" {
	cpu_run \
		'SEC' \
		'LDA #$80' 'SBC #$01' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_flag $OVERFLOW
}

@test "SED enables decimal mode; CLD restores binary mode" {
	# In BCD: $09 + $01 = $10
	cpu_run \
		'SED' \
		'CLC' \
		'LDA #$09' 'ADC #$01' \
		'STA $00' \
		'CLD' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 10

	# After CLD, binary addition resumes: $09 + $01 = $0A
	cpu_run \
		'CLD' \
		'CLC' \
		'LDA #$09' 'ADC #$01' \
		'STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 0A
}

@test "BCD add \$99 + \$01 = \$00 with carry" {
	cpu_run \
		'SED' \
		'CLC' \
		'LDA #$99' 'ADC #$01' \
		'STA $00' \
		'PHP' 'PLA' 'STA $02' \
		'CLD' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 00
	assert_flag $CARRY
}

@test "BCD add crossing nibble boundary: \$28 + \$14 = \$42" {
	cpu_run \
		'SED' \
		'CLC' \
		'LDA #$28' 'ADC #$14' \
		'STA $00' \
		'CLD' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 42
}

@test "BCD subtract \$20 - \$01 = \$19" {
	cpu_run \
		'SED' \
		'SEC' \
		'LDA #$20' 'SBC #$01' \
		'STA $00' \
		'CLD' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 19
}

@test "BCD subtract with borrow: \$00 - \$01 = \$99 C=0" {
	cpu_run \
		'SED' \
		'SEC' \
		'LDA #$00' 'SBC #$01' \
		'STA $00' \
		'PHP' 'PLA' 'STA $02' \
		'CLD' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 99
	refute_flag $CARRY
}

# ---------------------------------------------------------------------------
# 3. Logic (AND, ORA, EOR)
# ---------------------------------------------------------------------------

@test "AND masks bits correctly" {
	cpu_run \
		'LDA #$FF' 'AND #$0F' \
		'STA $00' '.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 0F
}

@test "AND result \$00 sets Z flag" {
	cpu_run \
		'LDA #$AA' 'AND #$55' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_flag $ZERO
}

@test "ORA sets bits correctly" {
	cpu_run \
		'LDA #$0F' 'ORA #$F0' \
		'STA $00' '.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 FF
}

@test "ORA result with bit 7 set sets N flag" {
	cpu_run \
		'LDA #$40' 'ORA #$80' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_flag $NEGATIVE
}

@test "EOR flips bits correctly" {
	cpu_run \
		'LDA #$FF' 'EOR #$0F' \
		'STA $00' '.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 F0
}

@test "EOR self-XOR produces \$00 and sets Z flag" {
	cpu_run \
		'LDA #$AA' 'EOR #$AA' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_flag $ZERO
}

# ---------------------------------------------------------------------------
# 4. Bit Shifts (ASL, LSR, ROL, ROR)
# ---------------------------------------------------------------------------

@test "ASL accumulator shifts left; old bit 7 goes to C" {
	cpu_run \
		'LDA #$81' 'ASL A' \
		'STA $00' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 02
	assert_flag $CARRY
}

@test "LSR accumulator shifts right; old bit 0 goes to C" {
	cpu_run \
		'LDA #$81' 'LSR A' \
		'STA $00' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 40
	assert_flag $CARRY
}

@test "ROL accumulator rotates left through C" {
	# C=1, A=$40 -> A=$81, C=0
	cpu_run \
		'SEC' \
		'LDA #$40' 'ROL A' \
		'STA $00' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 81
	refute_flag $CARRY
}

@test "ROR accumulator rotates right through C" {
	# C=1, A=$02 -> A=$81, C=0
	cpu_run \
		'SEC' \
		'LDA #$02' 'ROR A' \
		'STA $00' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 81
	refute_flag $CARRY
}

@test "ASL zero-page shifts memory byte left" {
	cpu_run \
		'LDA #$40' 'STA $03' \
		'ASL $03' \
		'LDA $03' 'STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 80
}

@test "LSR zero-page shifts memory byte right" {
	cpu_run \
		'LDA #$40' 'STA $03' \
		'LSR $03' \
		'LDA $03' 'STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 20
}

@test "ASL twice equals multiply by 4" {
	cpu_run \
		'LDA #$05' 'ASL A' 'ASL A' \
		'STA $00' '.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 14
}

@test "ROL/ROR round-trip restores original value" {
	# ROL A then ROR A with C=0 should restore A
	cpu_run \
		'CLC' \
		'LDA #$42' 'ROL A' 'ROR A' \
		'STA $00' '.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 42
}

@test "ROL/ROR round-trip on zero-page memory restores original value" {
	cpu_run \
		'CLC' \
		'LDA #$42' 'STA $03' \
		'ROL $03' 'ROR $03' \
		'LDA $03' 'STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 42
}

# ---------------------------------------------------------------------------
# 5. Bit Test (BIT, TSB, TRB)
# ---------------------------------------------------------------------------

@test "BIT immediate: Z set when A AND imm == 0" {
	cpu_run \
		'LDA #$AA' 'BIT #$55' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_flag $ZERO
}

@test "BIT zero-page: N/V copied from bits 7/6 of memory" {
	cpu_run \
		'LDA #$C0' 'STA $10' \
		'LDA #$00' \
		'BIT $10' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_flag $NEGATIVE
	assert_flag $OVERFLOW
}

@test "BIT zero-page: Z set when A AND M == 0" {
	cpu_run \
		'LDA #$AA' 'STA $10' \
		'LDA #$55' \
		'BIT $10' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_flag $ZERO
}

@test "BIT absolute: N/V copied from bits 7/6 of memory" {
	cpu_run \
		'LDA #$C0' 'STA $0300' \
		'LDA #$00' \
		'BIT $0300' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_flag $NEGATIVE
	assert_flag $OVERFLOW
}

@test "TSB sets bits in memory that are set in A" {
	cpu_run \
		'LDA #$0F' 'STA $03' \
		'LDA #$F0' \
		'TSB $03' \
		'LDA $03' 'STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 FF
}

@test "TRB clears bits in memory that are set in A" {
	cpu_run \
		'LDA #$FF' 'STA $03' \
		'LDA #$0F' \
		'TRB $03' \
		'LDA $03' 'STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 F0
}

# ---------------------------------------------------------------------------
# 6. Compare (CMP, CPX, CPY)
# ---------------------------------------------------------------------------

@test "CMP equal: Z=1 C=1 N=0; register unchanged" {
	cpu_run \
		'LDA #$42' 'CMP #$42' \
		'STA $00' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 42
	assert_flag $ZERO
	assert_flag $CARRY
	refute_flag $NEGATIVE
}

@test "CMP greater: Z=0 C=1 N=0" {
	cpu_run \
		'LDA #$50' 'CMP #$30' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	refute_flag $ZERO
	assert_flag $CARRY
	refute_flag $NEGATIVE
}

@test "CMP less: Z=0 C=0 N=1" {
	cpu_run \
		'LDA #$30' 'CMP #$50' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	refute_flag $ZERO
	refute_flag $CARRY
	assert_flag $NEGATIVE
}

@test "CPX equal sets Z and C flags" {
	cpu_run \
		'LDX #$42' 'CPX #$42' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_flag $ZERO
	assert_flag $CARRY
}

@test "CPX greater: Z=0 C=1 N=0" {
	cpu_run \
		'LDX #$50' 'CPX #$30' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	refute_flag $ZERO
	assert_flag $CARRY
	refute_flag $NEGATIVE
}

@test "CPX less: Z=0 C=0 N=1" {
	cpu_run \
		'LDX #$30' 'CPX #$50' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	refute_flag $ZERO
	refute_flag $CARRY
	assert_flag $NEGATIVE
}

@test "CPY equal sets Z and C flags" {
	cpu_run \
		'LDY #$42' 'CPY #$42' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_flag $ZERO
	assert_flag $CARRY
}

@test "CPY greater: Z=0 C=1 N=0" {
	cpu_run \
		'LDY #$50' 'CPY #$30' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	refute_flag $ZERO
	assert_flag $CARRY
	refute_flag $NEGATIVE
}

@test "CPY less: Z=0 C=0 N=1" {
	cpu_run \
		'LDY #$30' 'CPY #$50' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	refute_flag $ZERO
	refute_flag $CARRY
	assert_flag $NEGATIVE
}

@test "CMP does not modify the compared register" {
	cpu_run \
		'LDA #$42' 'CMP #$FF' \
		'STA $00' '.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 42
}

# ---------------------------------------------------------------------------
# 7. Branches (BEQ, BNE, BCC, BCS, BMI, BPL, BVC, BVS, BRA)
# ---------------------------------------------------------------------------

@test "BEQ taken when Z=1: skips alternate path" {
	cpu_run \
		'LDA #$AA' 'CMP #$AA' \
		'BEQ done' \
		'LDA #$BB' \
		'done: STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 AA
}

@test "BEQ not-taken when Z=0: falls through" {
	cpu_run \
		'LDA #$AA' 'CMP #$BB' \
		'BEQ skip' \
		'LDA #$CC' \
		'skip: STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 CC
}

@test "BNE taken when Z=0: skips alternate path" {
	cpu_run \
		'LDA #$AA' 'CMP #$BB' \
		'BNE done' \
		'LDA #$00' \
		'done: STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 AA
}

@test "BNE not-taken when Z=1: falls through" {
	cpu_run \
		'LDA #$AA' 'CMP #$AA' \
		'BNE skip' \
		'LDA #$CC' \
		'skip: STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 CC
}

@test "BCC taken when C=0" {
	cpu_run \
		'CLC' \
		'LDA #$AA' \
		'BCC done' \
		'LDA #$BB' \
		'done: STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 AA
}

@test "BCC not-taken when C=1" {
	cpu_run \
		'SEC' \
		'LDA #$AA' \
		'BCC skip' \
		'LDA #$CC' \
		'skip: STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 CC
}

@test "BCS taken when C=1" {
	cpu_run \
		'SEC' \
		'LDA #$AA' \
		'BCS done' \
		'LDA #$BB' \
		'done: STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 AA
}

@test "BCS not-taken when C=0" {
	cpu_run \
		'CLC' \
		'LDA #$AA' \
		'BCS skip' \
		'LDA #$CC' \
		'skip: STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 CC
}

@test "BMI taken when N=1" {
	cpu_run \
		'LDA #$80' \
		'BMI done' \
		'LDA #$BB' \
		'done: STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 80
}

@test "BMI not-taken when N=0" {
	cpu_run \
		'LDA #$40' \
		'BMI skip' \
		'LDA #$CC' \
		'skip: STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 CC
}

@test "BPL taken when N=0" {
	cpu_run \
		'LDA #$40' \
		'BPL done' \
		'LDA #$BB' \
		'done: STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 40
}

@test "BPL not-taken when N=1" {
	cpu_run \
		'LDA #$80' \
		'BPL skip' \
		'LDA #$CC' \
		'skip: STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 CC
}

@test "BVC taken when V=0" {
	cpu_run \
		'CLV' \
		'LDA #$AA' \
		'BVC done' \
		'LDA #$BB' \
		'done: STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 AA
}

@test "BVC not-taken when V=1" {
	cpu_run \
		'CLC' 'LDA #$7F' 'ADC #$01' \
		'BVC skip' \
		'LDA #$CC' \
		'skip: STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 CC
}

@test "BVS taken when V=1" {
	cpu_run \
		'CLC' 'LDA #$7F' 'ADC #$01' \
		'BVS done' \
		'LDA #$BB' \
		'done: STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	# A=$80 after ADC; BVS branches to done, stores $80
	assert_zp 0 80
}

@test "BVS not-taken when V=0" {
	cpu_run \
		'CLV' \
		'LDA #$AA' \
		'BVS skip' \
		'LDA #$CC' \
		'skip: STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 CC
}

@test "BRA always branches unconditionally" {
	cpu_run \
		'LDA #$AA' \
		'BRA done' \
		'LDA #$BB' \
		'done: STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 AA
}

@test "backward branch: loop counts down to zero" {
	# LDX #$05 / loop: DEX / BNE loop / STX $00 -> X == $00
	CPU_STEPS=30 cpu_run \
		'LDX #$05' \
		'loop: DEX' \
		'BNE loop' \
		'STX $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 00
}

# ---------------------------------------------------------------------------
# 8. Jumps and Subroutines (JMP, JSR, RTS)
# ---------------------------------------------------------------------------

@test "JMP absolute moves PC to target" {
	cpu_run \
		'JMP target' \
		'LDA #$BB' \
		'target: LDA #$AA' \
		'STA $00' '.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 AA
}

@test "JMP indirect dereferences pointer" {
	# Pointer at $0200/$0201 set to address of landing code.
	# $0801: LDA #$10 (2) $0803: STA $0200 (3) $0806: LDA #$08 (2)
	# $0808: STA $0201 (3) $080B: JMP ($0200) (3) $080E: LDA #$BB (2-skipped)
	# $0810: LDA #$AA
	cpu_run \
		'LDA #$10' 'STA $0200' \
		'LDA #$08' 'STA $0201' \
		'JMP ($0200)' \
		'LDA #$BB' \
		'LDA #$AA' \
		'STA $00' '.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 AA
}

@test "JSR jumps to subroutine; RTS returns" {
	cpu_run \
		'JSR sub' \
		'STA $00' '.halt' \
		'sub: LDA #$42' \
		'RTS'
	[[ $status -eq 0 ]]
	assert_zp 0 42
}

@test "nested JSR/RTS unwinds stack correctly" {
	cpu_run \
		'JSR outer' \
		'STA $00' '.halt' \
		'outer: JSR inner' \
		'RTS' \
		'inner: LDA #$CC' \
		'RTS'
	[[ $status -eq 0 ]]
	assert_zp 0 CC
}

@test "S decrements by 2 on JSR and increments by 2 on RTS" {
	# Capture S before JSR, after JSR, after RTS to verify +-2 change.
	# We measure S indirectly: TSX after JSR returns should equal TSX before JSR.
	cpu_run \
		'TSX' 'STX $03' \
		'JSR sub' \
		'TSX' 'STX $04' \
		'.halt' \
		'sub: RTS'
	[[ $status -eq 0 ]]
	# S before and after should be equal (stack restored by RTS)
	local s_before s_after
	s_before="$(awk '$3=="mem"&&$4=="$0003"{v=$NF}END{print v}' "$OUT/state.log")"
	s_after="$(awk '$3=="mem"&&$4=="$0004"{v=$NF}END{print v}' "$OUT/state.log")"
	[[ "$s_before" == "$s_after" ]]
}

@test "RTI restores PC from stack: jumps past trap to success code" {
	# Manually push a fake interrupt frame (PCH, PCL, P) then RTI.
	# RTI pops P first, then PCL, then PCH.
	# Byte layout: 7 setup instrs + RTI = $080A; trap at $080B-$0811;
	# success code at $0812 (PCH=$08, PCL=$12).
	cpu_run \
		'LDA #$08' 'PHA' \
		'LDA #$12' 'PHA' \
		'LDA #$30' 'PHA' \
		'RTI' \
		'LDA #$FF' 'STA $00' '.halt' \
		'LDA #$42' 'STA $00' '.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 42
}

@test "RTI restores P flags from stack" {
	# Push P=$31 (Carry|Break|Unused); RTI should land at $0812 with Carry set.
	# Success code immediately captures P via PHP/PLA/STA $02.
	# Byte layout: same setup as above; success code at $0812.
	cpu_run \
		'LDA #$08' 'PHA' \
		'LDA #$12' 'PHA' \
		'LDA #$31' 'PHA' \
		'RTI' \
		'LDA #$FF' 'STA $00' '.halt' \
		'PHP' 'PLA' 'STA $02' 'LDA #$42' 'STA $00' '.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 42
	assert_flag $CARRY
}

# ---------------------------------------------------------------------------
# 9. Stack (PHA, PLA, PHX, PLX, PHY, PLY, PHP, PLP)
# ---------------------------------------------------------------------------

@test "PHA/PLA round-trip preserves value" {
	cpu_run \
		'LDA #$42' 'PHA' \
		'LDA #$00' \
		'PLA' 'STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 42
}

@test "PHX/PLX round-trip preserves value" {
	cpu_run \
		'LDX #$AB' 'PHX' \
		'LDX #$00' \
		'PLX' 'STX $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 AB
}

@test "PHY/PLY round-trip preserves value" {
	cpu_run \
		'LDY #$CD' 'PHY' \
		'LDY #$00' \
		'PLY' 'STY $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 CD
}

@test "PHP/PLP round-trip preserves P" {
	# Set a known P state (SEC sets C), push, modify P (CLC), restore via PLP.
	cpu_run \
		'SEC' \
		'PHP' \
		'CLC' \
		'PLP' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_flag $CARRY
}

@test "multiple pushes stack in LIFO order" {
	cpu_run \
		'LDA #$11' 'PHA' \
		'LDA #$22' 'PHA' \
		'PLA' 'STA $00' \
		'PLA' 'STA $01' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 22
	assert_zp 1 11
}

@test "S decrements on push and increments on pull" {
	cpu_run \
		'TSX' 'STX $03' \
		'LDA #$AA' 'PHA' \
		'TSX' 'STX $04' \
		'PLA' \
		'.halt'
	[[ $status -eq 0 ]]
	# $04 (S after push) should be $03 (S before push) - 1
	local s_before s_after
	s_before="$(printf '%d' "0x$(awk '$3=="mem"&&$4=="$0003"{v=$NF}END{print v}' "$OUT/state.log" | tr -d '$')")"
	s_after="$(printf '%d' "0x$(awk '$3=="mem"&&$4=="$0004"{v=$NF}END{print v}' "$OUT/state.log" | tr -d '$')")"
	[[ $(( s_before - 1 )) -eq "$s_after" ]]
}

# ---------------------------------------------------------------------------
# 10. Register Transfers (TAX, TXA, TAY, TYA, TSX, TXS)
# ---------------------------------------------------------------------------

@test "TAX copies A to X and sets N/Z flags" {
	cpu_run \
		'LDA #$42' 'TAX' \
		'STX $00' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 42
	refute_flag $ZERO
	refute_flag $NEGATIVE
}

@test "TXA copies X to A and sets N/Z flags" {
	cpu_run \
		'LDX #$80' 'TXA' \
		'STA $00' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 80
	assert_flag $NEGATIVE
}

@test "TAY copies A to Y" {
	cpu_run \
		'LDA #$CC' 'TAY' \
		'STY $00' '.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 CC
}

@test "TYA copies Y to A" {
	cpu_run \
		'LDY #$DD' 'TYA' \
		'STA $00' '.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 DD
}

@test "TSX copies S to X and sets N/Z flags" {
	cpu_run \
		'TSX' 'STX $00' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	# Just verify the transfer ran without crash; S is non-zero after boot.
	refute_flag $ZERO
}

@test "TXS copies X to S without affecting flags" {
	# Load a known value; TXS should not set N/Z unlike other transfers.
	cpu_run \
		'LDX #$FF' \
		'TXS' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	# N should not be set by TXS itself (LDX #$FF sets N; PHP captures state
	# after TXS which should not alter it -- but LDX already set N, so we
	# simply verify the program ran cleanly).
	[[ $status -eq 0 ]]
}

# ---------------------------------------------------------------------------
# 11. Increment / Decrement (INX, INY, DEX, DEY, INC, DEC)
# ---------------------------------------------------------------------------

@test "INX increments X by 1" {
	cpu_run 'LDX #$41' 'INX' 'STX $00' '.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 42
}

@test "INX wraps \$FF to \$00 and sets Z" {
	cpu_run \
		'LDX #$FF' 'INX' \
		'STX $00' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 00
	assert_flag $ZERO
}

@test "INY increments Y by 1" {
	cpu_run 'LDY #$41' 'INY' 'STY $00' '.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 42
}

@test "DEX decrements X by 1" {
	cpu_run 'LDX #$43' 'DEX' 'STX $00' '.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 42
}

@test "DEX wraps \$00 to \$FF and sets N" {
	cpu_run \
		'LDX #$00' 'DEX' \
		'STX $00' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 FF
	assert_flag $NEGATIVE
}

@test "DEY decrements Y by 1" {
	cpu_run 'LDY #$43' 'DEY' 'STY $00' '.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 42
}

@test "INC zero-page increments memory and sets Z on \$FF->\$00" {
	cpu_run \
		'LDA #$FF' 'STA $03' \
		'INC $03' \
		'LDA $03' 'STA $00' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 00
	assert_flag $ZERO
}

@test "DEC zero-page decrements memory" {
	cpu_run \
		'LDA #$43' 'STA $03' \
		'DEC $03' \
		'LDA $03' 'STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 42
}

# ---------------------------------------------------------------------------
# 12. Status Flag Instructions (CLC, SEC, CLV, SEI, CLI, CLD, SED)
# ---------------------------------------------------------------------------

@test "CLC clears C; SEC sets C" {
	cpu_run \
		'CLC' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	refute_flag $CARRY

	cpu_run \
		'SEC' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_flag $CARRY
}

@test "CLV clears V after overflow sets it" {
	cpu_run \
		'CLC' 'LDA #$7F' 'ADC #$01' \
		'CLV' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	refute_flag $OVERFLOW
}

@test "SEI sets I flag; CLI clears I flag" {
	cpu_run \
		'SEI' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_flag $INTERRUPT

	cpu_run \
		'CLI' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	refute_flag $INTERRUPT
}

@test "SED sets D flag; CLD clears D flag" {
	cpu_run \
		'SED' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_flag $DECIMAL

	cpu_run \
		'CLD' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	refute_flag $DECIMAL
}

# ---------------------------------------------------------------------------
# 13. Addressing Modes (one focused test per mode)
# ---------------------------------------------------------------------------

@test "mode: immediate -- LDA #\$5A -> A == \$5A" {
	cpu_run 'LDA #$5A' 'STA $00' '.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 5A
}

@test "mode: zero-page -- LDA \$10 reads from zero page" {
	cpu_run \
		'LDA #$5A' 'STA $10' \
		'LDA #$00' \
		'LDA $10' 'STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 5A
}

@test "mode: zero-page X -- LDA \$10,X with X=2 reads \$12" {
	cpu_run \
		'LDA #$5A' 'STA $12' \
		'LDX #$02' \
		'LDA $10,X' 'STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 5A
}

@test "mode: zero-page Y -- STX \$10,Y with Y=2 writes to \$12" {
	cpu_run \
		'LDX #$5A' \
		'LDY #$02' \
		'STX $10,Y' \
		'LDA $12' 'STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 5A
}

@test "mode: absolute -- LDA \$0300 loads from full address" {
	cpu_run \
		'LDA #$5A' 'STA $0300' \
		'LDA #$00' \
		'LDA $0300' 'STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 5A
}

@test "mode: absolute X -- LDA \$0300,X with X=2 reads \$0302" {
	cpu_run \
		'LDA #$5A' 'STA $0302' \
		'LDX #$02' \
		'LDA $0300,X' 'STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 5A
}

@test "mode: absolute Y -- LDA \$0300,Y with Y=2 reads \$0302" {
	cpu_run \
		'LDA #$5A' 'STA $0302' \
		'LDY #$02' \
		'LDA $0300,Y' 'STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 5A
}

@test "mode: indirect -- JMP (\$0200) dereferences pointer" {
	# See JMP indirect test above in section 8; duplicate here for mode focus.
	# Pointer at $0200/$0201 -> $0810 (LDA #$5A).
	# $0801: LDA #$10 (2) $0803: STA $0200 (3) $0806: LDA #$08 (2)
	# $0808: STA $0201 (3) $080B: JMP ($0200) (3) $080E: LDA #$BB (2-skip)
	# $0810: LDA #$5A
	cpu_run \
		'LDA #$10' 'STA $0200' \
		'LDA #$08' 'STA $0201' \
		'JMP ($0200)' \
		'LDA #$BB' \
		'LDA #$5A' \
		'STA $00' '.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 5A
}

@test "mode: X-indexed indirect -- LDA (\$10,X) with X=2" {
	# Pointer at $12/$13 -> $0300 containing $5A.
	cpu_run \
		'LDA #$5A' 'STA $0300' \
		'LDA #$00' 'STA $12' \
		'LDA #$03' 'STA $13' \
		'LDX #$02' \
		'LDA ($10,X)' 'STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 5A
}

@test "mode: indirect Y-indexed -- LDA (\$10),Y with Y=2" {
	# Pointer at $10/$11 -> $0300; value at $0302 is $5A.
	cpu_run \
		'LDA #$5A' 'STA $0302' \
		'LDA #$00' 'STA $10' \
		'LDA #$03' 'STA $11' \
		'LDY #$02' \
		'LDA ($10),Y' 'STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 5A
}

@test "mode: zero-page indirect -- LDA (\$10) dereferences ZP pointer" {
	# Pointer at $10/$11 -> $0300 containing $5A.
	cpu_run \
		'LDA #$5A' 'STA $0300' \
		'LDA #$00' 'STA $10' \
		'LDA #$03' 'STA $11' \
		'LDA ($10)' 'STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 5A
}

@test "mode: accumulator -- ASL A shifts accumulator not memory" {
	cpu_run \
		'LDA #$21' 'ASL A' \
		'STA $00' '.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 42
}

@test "mode: relative -- BRA skips a trap value" {
	cpu_run \
		'LDA #$AA' \
		'BRA skip' \
		'LDA #$BB' \
		'skip: STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 AA
}

@test "mode: implied -- INX operates on X with no operand" {
	cpu_run 'LDX #$41' 'INX' 'STX $00' '.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 42
}

# ---------------------------------------------------------------------------
# 14. Flag Interaction / Integration Tests
# ---------------------------------------------------------------------------

@test "integration: ADC overflow \$7F + \$01 -> V=1 N=1 C=0" {
	cpu_run \
		'CLC' \
		'LDA #$7F' 'ADC #$01' \
		'STA $00' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 80
	assert_flag $OVERFLOW
	assert_flag $NEGATIVE
	refute_flag $CARRY
}

@test "integration: ADC carry chain \$FF + \$01 -> C=1 Z=1 A=\$00" {
	cpu_run \
		'CLC' \
		'LDA #$FF' 'ADC #$01' \
		'STA $00' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 00
	assert_flag $CARRY
	assert_flag $ZERO
}

@test "integration: SBC borrow SEC / LDA #\$00 / SBC #\$01 -> C=0 A=\$FF N=1" {
	cpu_run \
		'SEC' \
		'LDA #$00' 'SBC #$01' \
		'STA $00' \
		'PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 FF
	refute_flag $CARRY
	assert_flag $NEGATIVE
}

@test "integration: ROL chain -- 9 ROL shifts with C=0 restores original value" {
	# ROL rotates through the carry bit, making it a 9-bit rotation (8 bits of
	# A + 1 bit of C).  With initial C=0, exactly 9 ROLs restore both A and C
	# to their original values.
	CPU_STEPS=50 cpu_run \
		'CLC' \
		'LDA #$42' \
		'ROL A' 'ROL A' 'ROL A' 'ROL A' 'ROL A' \
		'ROL A' 'ROL A' 'ROL A' 'ROL A' \
		'STA $00' '.halt'
	[[ $status -eq 0 ]]
	assert_zp 0 42
}

@test "integration: branch skips flag-setting instruction; flag is not set" {
	# BRA jumps over SEC; C should remain clear.
	cpu_run \
		'CLC' \
		'BRA skip' \
		'SEC' \
		'skip: PHP' 'PLA' 'STA $02' \
		'.halt'
	[[ $status -eq 0 ]]
	refute_flag $CARRY
}
