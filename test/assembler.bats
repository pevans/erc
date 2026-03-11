setup_file() { load assembler_helper; setup_file; }
setup()      { load assembler_helper; setup; }
teardown()   { load assembler_helper; teardown; }

# ---------------------------------------------------------------------------
# Output format
# ---------------------------------------------------------------------------

@test "valid program exits with status 0" {
	asm 'NOP'
	[[ $status -eq 0 ]]
}

@test "output is exactly 143360 bytes" {
	asm 'NOP'
	[[ $status -eq 0 ]]
	[[ $(wc -c <"$TMP/out.dsk" | tr -d ' ') -eq 143360 ]]
}

# ---------------------------------------------------------------------------
# Sector count byte (dsk[0])
# ---------------------------------------------------------------------------

@test "small program sets sector count to 1" {
	# 1 byte of code: total = 2, ceil(2/256) = 1
	asm 'NOP'
	[[ $status -eq 0 ]]
	[[ "$(byte_at "$TMP/out.dsk" 0)" == "01" ]]
}

@test "256-byte program sets sector count to 2" {
	# 256 NOPs = 256 bytes of code: total = 257, ceil(257/256) = 2
	local src="$TMP/test.s"
	for _ in $(seq 256); do echo NOP; done >"$src"
	run "$ASSEMBLER" -o "$TMP/out.dsk" "$src"
	[[ $status -eq 0 ]]
	[[ "$(byte_at "$TMP/out.dsk" 0)" == "02" ]]
}

# ---------------------------------------------------------------------------
# .halt directive
# ---------------------------------------------------------------------------

@test ".halt emits JMP to current address" {
	# .halt at $0801 should emit JMP $0801 = 4C 01 08
	asm '.halt'
	[[ $status -eq 0 ]]
	[[ "$(byte_at "$TMP/out.dsk" 1)" == "4c" ]]
	[[ "$(byte_at "$TMP/out.dsk" 2)" == "01" ]]
	[[ "$(byte_at "$TMP/out.dsk" 3)" == "08" ]]
}

# ---------------------------------------------------------------------------
# .byte directive
# ---------------------------------------------------------------------------

@test ".byte emits the specified bytes" {
	asm '.byte $AB, $CD'
	[[ $status -eq 0 ]]
	[[ "$(byte_at "$TMP/out.dsk" 1)" == "ab" ]]
	[[ "$(byte_at "$TMP/out.dsk" 2)" == "cd" ]]
}

# ---------------------------------------------------------------------------
# .word directive
# ---------------------------------------------------------------------------

@test ".word emits 16-bit values in little-endian order" {
	asm '.word $1234'
	[[ $status -eq 0 ]]
	[[ "$(byte_at "$TMP/out.dsk" 1)" == "34" ]]
	[[ "$(byte_at "$TMP/out.dsk" 2)" == "12" ]]
}

# ---------------------------------------------------------------------------
# .org directive
# ---------------------------------------------------------------------------

@test ".org changes the assembly origin" {
	# .halt at $0900 should emit JMP $0900 = 4C 00 09
	asm '.org $0900' '.halt'
	[[ $status -eq 0 ]]
	[[ "$(byte_at "$TMP/out.dsk" 1)" == "4c" ]]
	[[ "$(byte_at "$TMP/out.dsk" 2)" == "00" ]]
	[[ "$(byte_at "$TMP/out.dsk" 3)" == "09" ]]
}

# ---------------------------------------------------------------------------
# Addressing modes
# ---------------------------------------------------------------------------

@test "immediate mode: LDA #\$FF" {
	asm 'LDA #$FF'
	[[ $status -eq 0 ]]
	[[ "$(byte_at "$TMP/out.dsk" 1)" == "a9" ]]
	[[ "$(byte_at "$TMP/out.dsk" 2)" == "ff" ]]
}

@test "zero page mode: LDA \$10" {
	asm 'LDA $10'
	[[ $status -eq 0 ]]
	[[ "$(byte_at "$TMP/out.dsk" 1)" == "a5" ]]
	[[ "$(byte_at "$TMP/out.dsk" 2)" == "10" ]]
}

@test "four-digit hex forces absolute mode: LDA \$0010" {
	# $0010 fits in a byte but the four-digit form forces ABS (opcode $AD)
	asm 'LDA $0010'
	[[ $status -eq 0 ]]
	[[ "$(byte_at "$TMP/out.dsk" 1)" == "ad" ]]
	[[ "$(byte_at "$TMP/out.dsk" 2)" == "10" ]]
	[[ "$(byte_at "$TMP/out.dsk" 3)" == "00" ]]
}

@test "implied mode: NOP" {
	asm 'NOP'
	[[ $status -eq 0 ]]
	[[ "$(byte_at "$TMP/out.dsk" 1)" == "ea" ]]
}

@test "implied mode: CLC" {
	asm 'CLC'
	[[ $status -eq 0 ]]
	[[ "$(byte_at "$TMP/out.dsk" 1)" == "18" ]]
}

@test "accumulator mode: ASL A" {
	asm 'ASL A'
	[[ $status -eq 0 ]]
	[[ "$(byte_at "$TMP/out.dsk" 1)" == "0a" ]]
}

@test "accumulator mode: ASL with no operand" {
	asm 'ASL'
	[[ $status -eq 0 ]]
	[[ "$(byte_at "$TMP/out.dsk" 1)" == "0a" ]]
}

@test "absolute X-indexed mode: LDA \$0400,X" {
	asm 'LDA $0400,X'
	[[ $status -eq 0 ]]
	[[ "$(byte_at "$TMP/out.dsk" 1)" == "bd" ]]
	[[ "$(byte_at "$TMP/out.dsk" 2)" == "00" ]]
	[[ "$(byte_at "$TMP/out.dsk" 3)" == "04" ]]
}

@test "zero page X-indexed mode: LDA \$10,X" {
	asm 'LDA $10,X'
	[[ $status -eq 0 ]]
	[[ "$(byte_at "$TMP/out.dsk" 1)" == "b5" ]]
	[[ "$(byte_at "$TMP/out.dsk" 2)" == "10" ]]
}

@test "absolute indirect mode: JMP (\$FFFC)" {
	asm 'JMP ($FFFC)'
	[[ $status -eq 0 ]]
	[[ "$(byte_at "$TMP/out.dsk" 1)" == "6c" ]]
	[[ "$(byte_at "$TMP/out.dsk" 2)" == "fc" ]]
	[[ "$(byte_at "$TMP/out.dsk" 3)" == "ff" ]]
}

@test "X-indexed indirect mode: LDA (\$10,X)" {
	asm 'LDA ($10,X)'
	[[ $status -eq 0 ]]
	[[ "$(byte_at "$TMP/out.dsk" 1)" == "a1" ]]
	[[ "$(byte_at "$TMP/out.dsk" 2)" == "10" ]]
}

@test "indirect Y-indexed mode: LDA (\$10),Y" {
	asm 'LDA ($10),Y'
	[[ $status -eq 0 ]]
	[[ "$(byte_at "$TMP/out.dsk" 1)" == "b1" ]]
	[[ "$(byte_at "$TMP/out.dsk" 2)" == "10" ]]
}

@test "zero page indirect mode: LDA (\$10)" {
	asm 'LDA ($10)'
	[[ $status -eq 0 ]]
	[[ "$(byte_at "$TMP/out.dsk" 1)" == "b2" ]]
	[[ "$(byte_at "$TMP/out.dsk" 2)" == "10" ]]
}

# ---------------------------------------------------------------------------
# Labels and branches
# ---------------------------------------------------------------------------

@test "forward branch label resolves to correct relative offset" {
	# BEQ at $0801 (next PC = $0803), NOP at $0803, done: at $0804
	# offset = $0804 - $0803 = 1
	asm 'BEQ done' 'NOP' 'done: NOP'
	[[ $status -eq 0 ]]
	[[ "$(byte_at "$TMP/out.dsk" 1)" == "f0" ]]
	[[ "$(byte_at "$TMP/out.dsk" 2)" == "01" ]]
}

@test "backward branch label resolves to correct negative relative offset" {
	# start: at $0801, two NOPs, BEQ at $0803 (next PC = $0805)
	# offset = $0801 - $0805 = -4 = 0xfc
	asm 'start: NOP' 'NOP' 'BEQ start'
	[[ $status -eq 0 ]]
	[[ "$(byte_at "$TMP/out.dsk" 3)" == "f0" ]]
	[[ "$(byte_at "$TMP/out.dsk" 4)" == "fc" ]]
}

@test "JMP to label uses absolute addressing" {
	# start: NOP at $0801, JMP start at $0802 => 4C 01 08
	asm 'start: NOP' 'JMP start'
	[[ $status -eq 0 ]]
	[[ "$(byte_at "$TMP/out.dsk" 2)" == "4c" ]]
	[[ "$(byte_at "$TMP/out.dsk" 3)" == "01" ]]
	[[ "$(byte_at "$TMP/out.dsk" 4)" == "08" ]]
}

@test "label on its own line sets address for next instruction" {
	# loop: at $0801 (no code), LDA #$00 at $0801 (2 bytes), JMP loop at $0803
	# JMP loop => 4C 01 08 at dsk[3..5]
	asm 'loop:' 'LDA #$00' 'JMP loop'
	[[ $status -eq 0 ]]
	[[ "$(byte_at "$TMP/out.dsk" 3)" == "4c" ]]
	[[ "$(byte_at "$TMP/out.dsk" 4)" == "01" ]]
	[[ "$(byte_at "$TMP/out.dsk" 5)" == "08" ]]
}

# ---------------------------------------------------------------------------
# Syntax features
# ---------------------------------------------------------------------------

@test "lowercase mnemonics are accepted" {
	asm 'lda #$ff'
	[[ $status -eq 0 ]]
	[[ "$(byte_at "$TMP/out.dsk" 1)" == "a9" ]]
	[[ "$(byte_at "$TMP/out.dsk" 2)" == "ff" ]]
}

@test "comments are ignored" {
	asm '; this is a comment' 'NOP ; inline comment'
	[[ $status -eq 0 ]]
	[[ "$(byte_at "$TMP/out.dsk" 1)" == "ea" ]]
}

@test "blank lines are ignored" {
	asm '' 'NOP' ''
	[[ $status -eq 0 ]]
	[[ "$(byte_at "$TMP/out.dsk" 1)" == "ea" ]]
}

# ---------------------------------------------------------------------------
# CLI interface
# ---------------------------------------------------------------------------

@test "stdin input via - produces correct output" {
	run bash -c "echo 'NOP' | \"$ASSEMBLER\" -o \"$TMP/stdin.dsk\" -"
	[[ $status -eq 0 ]]
	[[ $(wc -c <"$TMP/stdin.dsk" | tr -d ' ') -eq 143360 ]]
}

@test "stdout output when -o is omitted" {
	local src="$TMP/test.s"
	echo 'NOP' >"$src"
	run bash -c "\"$ASSEMBLER\" \"$src\" | wc -c | tr -d ' \t'"
	[[ $status -eq 0 ]]
	[[ "$output" -eq 143360 ]]
}

# ---------------------------------------------------------------------------
# Error cases
# ---------------------------------------------------------------------------

@test "unknown mnemonic exits with non-zero status" {
	local src="$TMP/test.s"
	echo 'XYZ' >"$src"
	run "$ASSEMBLER" -o /dev/null "$src"
	[[ $status -ne 0 ]]
}

@test "unknown mnemonic error message includes filename and line number" {
	local src="$TMP/test.s"
	printf 'NOP\nXYZ\n' >"$src"
	msg=$("$ASSEMBLER" -o /dev/null "$src" 2>&1 || true)
	[[ "$msg" =~ "test.s:2" ]]
}

@test "duplicate label exits with non-zero status" {
	local src="$TMP/test.s"
	printf 'foo: NOP\nfoo: NOP\n' >"$src"
	run "$ASSEMBLER" -o /dev/null "$src"
	[[ $status -ne 0 ]]
}

@test "duplicate label error message mentions duplicate" {
	local src="$TMP/test.s"
	printf 'foo: NOP\nfoo: NOP\n' >"$src"
	msg=$("$ASSEMBLER" -o /dev/null "$src" 2>&1 || true)
	[[ "$msg" =~ "duplicate" ]]
}

@test "undefined label reference exits with non-zero status" {
	local src="$TMP/test.s"
	echo 'JMP nowhere' >"$src"
	run "$ASSEMBLER" -o /dev/null "$src"
	[[ $status -ne 0 ]]
}

@test "undefined label error message mentions undefined" {
	local src="$TMP/test.s"
	echo 'JMP nowhere' >"$src"
	msg=$("$ASSEMBLER" -o /dev/null "$src" 2>&1 || true)
	[[ "$msg" =~ "undefined" ]]
}

@test "branch target out of range exits with non-zero status" {
	# BEQ at $0801, 128 NOPs, target at $0803+128=$0883
	# offset = $0883 - $0803 = 128 > 127: out of range
	local src="$TMP/test.s"
	{
		echo 'BEQ target'
		for _ in $(seq 128); do echo NOP; done
		echo 'target: NOP'
	} >"$src"
	run "$ASSEMBLER" -o /dev/null "$src"
	[[ $status -ne 0 ]]
}

@test "branch target out of range error message mentions range" {
	local src="$TMP/test.s"
	{
		echo 'BEQ target'
		for _ in $(seq 128); do echo NOP; done
		echo 'target: NOP'
	} >"$src"
	msg=$("$ASSEMBLER" -o /dev/null "$src" 2>&1 || true)
	[[ "$msg" =~ "range" ]]
}
