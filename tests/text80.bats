setup_file() { load text80_helper; setup_file; }
setup()      { load text80_helper; setup; }
teardown()   { load text80_helper; teardown; }

# ---------------------------------------------------------------------------
# Screen geometry (2.2)
# ---------------------------------------------------------------------------

@test "80-col text frame is 560x384" {
	asm_video \
		'STA $C001' 'STA $C00D' \
		'LDA #$A0' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' '.halt'
	[[ $status -eq 0 ]]
	grep -q '^step 4999: video screen 560x384$' "$OUT/video.frame"
}

# ---------------------------------------------------------------------------
# Mode dispatch (6.1)
# ---------------------------------------------------------------------------

@test "TEXT off with 80COL on does not render 80-col text" {
	# TEXT off + 80COL on: should not produce 80-column text output.
	# Fill aux and main with 'A'; with TEXT off the screen must not show text.
	asm_video \
		'STA $C001' 'STA $C00D' 'STA $C050' \
		'STA $C055' \
		'LDA #$C1' 'LDX #$00' \
		'fill_a: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_a' \
		'STA $C054' \
		'LDA #$C1' 'LDX #$00' \
		'fill_m: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_m' \
		'.halt'
	[[ $status -eq 0 ]]

	# Capture frame without TEXT on
	local no_text_pixels
	no_text_pixels=$(tail -n +3 "$OUT/video.frame")

	# Now capture with TEXT + 80COL both on
	rm -rf "$OUT"; mkdir -p "$OUT"
	asm_video \
		'STA $C001' 'STA $C00D' 'STA $C051' \
		'STA $C055' \
		'LDA #$C1' 'LDX #$00' \
		'fill_a2: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_a2' \
		'STA $C054' \
		'LDA #$C1' 'LDX #$00' \
		'fill_m2: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_m2' \
		'.halt'
	[[ $status -eq 0 ]]
	local text_on_pixels
	text_on_pixels=$(tail -n +3 "$OUT/video.frame")

	[[ "$no_text_pixels" != "$text_on_pixels" ]]
}

@test "80COL off does not display aux memory content" {
	# With 80STORE on but 80COL off: 40-col renderer reads main only.
	# Fill aux with 'A' and main with spaces; screen must show only spaces.
	asm_video \
		'STA $C001' \
		'STA $C055' \
		'LDA #$C1' 'LDX #$00' \
		'fill_a: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_a' \
		'STA $C054' \
		'LDA #$A0' 'LDX #$00' \
		'fill_m: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_m' \
		'.halt'
	[[ $status -eq 0 ]]
	[[ $(color_count) -eq 1 ]]
}

# ---------------------------------------------------------------------------
# Two-bank interleave: column placement (4.1, 4.2)
# ---------------------------------------------------------------------------

@test "aux memory character appears in even screen column" {
	# Aux 'A' renders at the left (even) column of each 14-pixel pair.
	# Main spaces keep the right (odd) column uniform.
	asm_video \
		'STA $C001' 'STA $C00D' \
		'STA $C055' \
		'LDA #$C1' 'LDX #$00' \
		'fill_a: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_a' \
		'STA $C054' \
		'LDA #$A0' 'LDX #$00' \
		'fill_m: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_m' \
		'.halt'
	[[ $status -eq 0 ]]
	any_col_range_has_mixed_colors 0 15 0 7
	col_range_is_uniform 0 7 7
}

@test "main memory character appears in odd screen column" {
	# Main 'A' renders at the right (odd) column of each 14-pixel pair.
	# Aux spaces keep the left (even) column uniform.
	asm_video \
		'STA $C001' 'STA $C00D' \
		'STA $C055' \
		'LDA #$A0' 'LDX #$00' \
		'fill_a: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_a' \
		'STA $C054' \
		'LDA #$C1' 'LDX #$00' \
		'fill_m: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_m' \
		'.halt'
	[[ $status -eq 0 ]]
	col_range_is_uniform 0 0 7
	any_col_range_has_mixed_colors 0 15 7 7
}

# ---------------------------------------------------------------------------
# Address-to-row mapping (4.2)
# ---------------------------------------------------------------------------

@test "aux address \$0400 maps to row 0: character in pixel rows 0-15" {
	asm_video \
		'STA $C001' 'STA $C00D' \
		'STA $C055' \
		'LDA #$A0' 'LDX #$00' \
		'fill_a: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_a' \
		'LDA #$C1' 'STA $0400' \
		'STA $C054' \
		'LDA #$A0' 'LDX #$00' \
		'fill_m: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_m' \
		'.halt'
	[[ $status -eq 0 ]]
	any_col_range_has_mixed_colors 0 15 0 7
	col_range_is_uniform 16 0 7
}

@test "aux address \$0480 maps to row 1: character in pixel rows 16-31" {
	asm_video \
		'STA $C001' 'STA $C00D' \
		'STA $C055' \
		'LDA #$A0' 'LDX #$00' \
		'fill_a: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_a' \
		'LDA #$C1' 'STA $0480' \
		'STA $C054' \
		'LDA #$A0' 'LDX #$00' \
		'fill_m: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_m' \
		'.halt'
	[[ $status -eq 0 ]]
	col_range_is_uniform 0 0 7
	any_col_range_has_mixed_colors 16 31 0 7
}

@test "aux address \$07D0 maps to row 23: character in pixel rows 368-383" {
	asm_video \
		'STA $C001' 'STA $C00D' \
		'STA $C055' \
		'LDA #$A0' 'LDX #$00' \
		'fill_a: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_a' \
		'LDA #$C1' 'STA $07D0' \
		'STA $C054' \
		'LDA #$A0' 'LDX #$00' \
		'fill_m: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_m' \
		'.halt'
	[[ $status -eq 0 ]]
	col_range_is_uniform 367 0 7
	any_col_range_has_mixed_colors 368 383 0 7
}

@test "main address \$0480 maps to row 1: character in pixel rows 16-31" {
	asm_video \
		'STA $C001' 'STA $C00D' \
		'STA $C055' \
		'LDA #$A0' 'LDX #$00' \
		'fill_a: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_a' \
		'STA $C054' \
		'LDA #$A0' 'LDX #$00' \
		'fill_m: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_m' \
		'LDA #$C1' 'STA $0480' \
		'.halt'
	[[ $status -eq 0 ]]
	col_range_is_uniform 0 7 7
	any_col_range_has_mixed_colors 16 31 7 7
}

# ---------------------------------------------------------------------------
# Address-to-column mapping (4.2)
# ---------------------------------------------------------------------------

@test "aux address \$0401 maps to column 2: character in pixel columns 14-20" {
	# Aux offset 1 -> 40-col position 1 -> screen column 2 -> x=14-20.
	asm_video \
		'STA $C001' 'STA $C00D' \
		'STA $C055' \
		'LDA #$A0' 'LDX #$00' \
		'fill_a: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_a' \
		'LDA #$C1' 'STA $0401' \
		'STA $C054' \
		'LDA #$A0' 'LDX #$00' \
		'fill_m: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_m' \
		'.halt'
	[[ $status -eq 0 ]]
	col_range_is_uniform 0 0 14
	any_col_range_has_mixed_colors 0 15 14 7
}

# ---------------------------------------------------------------------------
# Hole bytes (3.1 via spec 9)
# ---------------------------------------------------------------------------

@test "hole byte \$0478 in main does not affect 80-col display" {
	asm_video \
		'STA $C001' 'STA $C00D' \
		'STA $C055' \
		'LDA #$A0' 'LDX #$00' \
		'fill_a: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_a' \
		'STA $C054' \
		'LDA #$A0' 'LDX #$00' \
		'fill_m: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_m' \
		'LDA #$C1' 'STA $0478' \
		'.halt'
	[[ $status -eq 0 ]]
	[[ $(color_count) -eq 1 ]]
}

@test "hole byte \$0478 in aux does not affect 80-col display" {
	asm_video \
		'STA $C001' 'STA $C00D' \
		'STA $C055' \
		'LDA #$A0' 'LDX #$00' \
		'fill_a: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_a' \
		'LDA #$C1' 'STA $0478' \
		'STA $C054' \
		'LDA #$A0' 'LDX #$00' \
		'fill_m: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_m' \
		'.halt'
	[[ $status -eq 0 ]]
	[[ $(color_count) -eq 1 ]]
}

# ---------------------------------------------------------------------------
# Character encoding (6.2, 6.3)
# ---------------------------------------------------------------------------

@test "80-col screen filled with spaces has one color" {
	asm_video \
		'STA $C001' 'STA $C00D' \
		'STA $C055' \
		'LDA #$A0' 'LDX #$00' \
		'fill_a: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_a' \
		'STA $C054' \
		'LDA #$A0' 'LDX #$00' \
		'fill_m: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_m' \
		'.halt'
	[[ $status -eq 0 ]]
	[[ $(color_count) -eq 1 ]]
}

@test "80-col screen filled with normal A has two colors" {
	asm_video \
		'STA $C001' 'STA $C00D' \
		'STA $C055' \
		'LDA #$C1' 'LDX #$00' \
		'fill_a: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_a' \
		'STA $C054' \
		'LDA #$C1' 'LDX #$00' \
		'fill_m: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_m' \
		'.halt'
	[[ $status -eq 0 ]]
	[[ $(color_count) -eq 2 ]]
}

@test "80-col ALTCHAR on renders differently from ALTCHAR off for \$40-\$5F range" {
	# With ALTCHAR on, $40-$5F are MouseText glyphs instead of inverse uppercase.
	asm_video \
		'STA $C001' 'STA $C00D' \
		'STA $C055' \
		'LDA #$40' 'LDX #$00' \
		'fill_a: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_a' \
		'STA $C054' \
		'LDA #$40' 'LDX #$00' \
		'fill_m: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_m' \
		'.halt'
	[[ $status -eq 0 ]]
	local altchar_off_pixels
	altchar_off_pixels=$(tail -n +3 "$OUT/video.frame")

	rm -rf "$OUT"; mkdir -p "$OUT"
	asm_video \
		'STA $C001' 'STA $C00D' 'STA $C00F' \
		'STA $C055' \
		'LDA #$40' 'LDX #$00' \
		'fill_a2: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_a2' \
		'STA $C054' \
		'LDA #$40' 'LDX #$00' \
		'fill_m2: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_m2' \
		'.halt'
	[[ $status -eq 0 ]]
	local altchar_on_pixels
	altchar_on_pixels=$(tail -n +3 "$OUT/video.frame")

	[[ "$altchar_off_pixels" != "$altchar_on_pixels" ]]
}

@test "80-col screen filled with inverse A has two colors" {
	asm_video \
		'STA $C001' 'STA $C00D' \
		'STA $C055' \
		'LDA #$01' 'LDX #$00' \
		'fill_a: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_a' \
		'STA $C054' \
		'LDA #$01' 'LDX #$00' \
		'fill_m: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_m' \
		'.halt'
	[[ $status -eq 0 ]]
	[[ $(color_count) -eq 2 ]]
}

# ---------------------------------------------------------------------------
# Flash (6.4)
#
# The flash transition occurs at cycle 272480.  The nested Y/X delay loop
# burns ~329000 cycles total; capturing at step 115000 lands at roughly
# 290000 cycles, safely inside the flash-off window [272480, 544959].
# ---------------------------------------------------------------------------

@test "80-col flash A pixel pattern during flash-on matches inverse A pixel pattern" {
	asm_video \
		'STA $C001' 'STA $C00D' \
		'STA $C055' \
		'LDA #$41' 'LDX #$00' \
		'fill_a: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_a' \
		'STA $C054' \
		'LDA #$41' 'LDX #$00' \
		'fill_m: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_m' \
		'.halt'
	[[ $status -eq 0 ]]
	local flash_on_pixels
	flash_on_pixels=$(tail -n +3 "$OUT/video.frame")

	rm -rf "$OUT"; mkdir -p "$OUT"
	asm_video \
		'STA $C001' 'STA $C00D' \
		'STA $C055' \
		'LDA #$01' 'LDX #$00' \
		'fill_a: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_a' \
		'STA $C054' \
		'LDA #$01' 'LDX #$00' \
		'fill_m: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_m' \
		'.halt'
	[[ $status -eq 0 ]]
	local inverse_pixels
	inverse_pixels=$(tail -n +3 "$OUT/video.frame")

	[[ "$flash_on_pixels" == "$inverse_pixels" ]]
}

@test "80-col flash chars render as normal during flash-off phase" {
	asm_video_long 115000 \
		'STA $C001' 'STA $C00D' \
		'STA $C055' \
		'LDA #$41' 'LDX #$00' \
		'fill_a: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_a' \
		'STA $C054' \
		'LDA #$41' 'LDX #$00' \
		'fill_m: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_m' \
		'LDY #$00' \
		'outer: LDX #$00' \
		'inner: DEX' \
		'BNE inner' \
		'DEY' \
		'BNE outer' \
		'.halt'
	[[ $status -eq 0 ]]
	local flash_off_pixels
	flash_off_pixels=$(tail -n +3 "$OUT/video.frame")

	rm -rf "$OUT"; mkdir -p "$OUT"
	asm_video \
		'STA $C001' 'STA $C00D' \
		'STA $C055' \
		'LDA #$C1' 'LDX #$00' \
		'fill_a: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_a' \
		'STA $C054' \
		'LDA #$C1' 'LDX #$00' \
		'fill_m: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_m' \
		'.halt'
	[[ $status -eq 0 ]]
	local normal_pixels
	normal_pixels=$(tail -n +3 "$OUT/video.frame")

	[[ "$flash_off_pixels" == "$normal_pixels" ]]
}

@test "80-col flash chars pixel pattern differs between flash-on and flash-off phases" {
	asm_video \
		'STA $C001' 'STA $C00D' \
		'STA $C055' \
		'LDA #$41' 'LDX #$00' \
		'fill_a: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_a' \
		'STA $C054' \
		'LDA #$41' 'LDX #$00' \
		'fill_m: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_m' \
		'.halt'
	[[ $status -eq 0 ]]
	local flash_on_frame
	flash_on_frame=$(tail -n +2 "$OUT/video.frame")

	rm -rf "$OUT"; mkdir -p "$OUT"
	asm_video_long 115000 \
		'STA $C001' 'STA $C00D' \
		'STA $C055' \
		'LDA #$41' 'LDX #$00' \
		'fill_a: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_a' \
		'STA $C054' \
		'LDA #$41' 'LDX #$00' \
		'fill_m: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_m' \
		'LDY #$00' \
		'outer: LDX #$00' \
		'inner: DEX' \
		'BNE inner' \
		'DEY' \
		'BNE outer' \
		'.halt'
	[[ $status -eq 0 ]]
	local flash_off_frame
	flash_off_frame=$(tail -n +2 "$OUT/video.frame")

	[[ "$flash_on_frame" != "$flash_off_frame" ]]
}

# ---------------------------------------------------------------------------
# Main-bank column mapping (4.2)
# ---------------------------------------------------------------------------

@test "main address \$0401 maps to column 3: character in pixel columns 21-27" {
	# Main offset 1 -> 40-col position 1 -> screen column 3 (2*1+1) -> x=21-27.
	asm_video \
		'STA $C001' 'STA $C00D' \
		'STA $C055' \
		'LDA #$A0' 'LDX #$00' \
		'fill_a: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_a' \
		'STA $C054' \
		'LDA #$A0' 'LDX #$00' \
		'fill_m: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_m' \
		'LDA #$C1' 'STA $0401' \
		'.halt'
	[[ $status -eq 0 ]]
	col_range_is_uniform 0 0 21
	any_col_range_has_mixed_colors 0 15 21 7
}

# ---------------------------------------------------------------------------
# ALTCHAR $60-$7F range (6.3)
# ---------------------------------------------------------------------------

@test "80-col ALTCHAR on renders differently from ALTCHAR off for \$60-\$7F range" {
	# With ALTCHAR off, $60-$7F render as inverse lowercase (flash chars).
	# With ALTCHAR on, $60-$7F render as normal lowercase.
	asm_video \
		'STA $C001' 'STA $C00D' \
		'STA $C055' \
		'LDA #$60' 'LDX #$00' \
		'fill_a: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_a' \
		'STA $C054' \
		'LDA #$60' 'LDX #$00' \
		'fill_m: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_m' \
		'.halt'
	[[ $status -eq 0 ]]
	local altchar_off_pixels
	altchar_off_pixels=$(tail -n +3 "$OUT/video.frame")

	rm -rf "$OUT"; mkdir -p "$OUT"
	asm_video \
		'STA $C001' 'STA $C00D' 'STA $C00F' \
		'STA $C055' \
		'LDA #$60' 'LDX #$00' \
		'fill_a2: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_a2' \
		'STA $C054' \
		'LDA #$60' 'LDX #$00' \
		'fill_m2: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_m2' \
		'.halt'
	[[ $status -eq 0 ]]
	local altchar_on_pixels
	altchar_on_pixels=$(tail -n +3 "$OUT/video.frame")

	[[ "$altchar_off_pixels" != "$altchar_on_pixels" ]]
}

# ---------------------------------------------------------------------------
# Page 2 without 80STORE (8.1)
# ---------------------------------------------------------------------------

@test "80-col without 80STORE and PAGE2 on displays page 2 content" {
	# 80STORE off, PAGE2 on: display reads from $0800-$0BFF in both banks.
	# Write 'A' to page 2 in both banks (using RAMWRT to reach aux $0800).
	# Leave page 1 as spaces. Display should show 'A'.
	asm_video \
		'STA $C005' \
		'LDA #$C1' 'LDX #$00' \
		'fill_a: STA $0800,X' 'STA $0900,X' 'STA $0A00,X' 'STA $0B00,X' \
		'INX' 'BNE fill_a' \
		'STA $C004' \
		'LDA #$C1' 'LDX #$00' \
		'fill_m: STA $0800,X' 'STA $0900,X' 'STA $0A00,X' 'STA $0B00,X' \
		'INX' 'BNE fill_m' \
		'LDA #$A0' 'LDX #$00' \
		'fill_p1: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_p1' \
		'STA $C00D' 'STA $C055' \
		'.halt'
	[[ $status -eq 0 ]]
	[[ $(color_count) -eq 2 ]]
}

# ---------------------------------------------------------------------------
# 80STORE + PAGE2 display source (8.2)
# ---------------------------------------------------------------------------

@test "80-col with 80STORE displays page 1 regardless of PAGE2 state" {
	# With 80STORE on, PAGE2 controls memory routing but the display always
	# reads from page 1.  Write 'A' to aux page 1, then leave PAGE2 on;
	# display should still show the page 1 content.
	asm_video \
		'STA $C001' 'STA $C00D' \
		'STA $C055' \
		'LDA #$C1' 'LDX #$00' \
		'fill_a: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_a' \
		'STA $C054' \
		'LDA #$C1' 'LDX #$00' \
		'fill_m: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_m' \
		'.halt'
	[[ $status -eq 0 ]]
	local page2_off_pixels
	page2_off_pixels=$(tail -n +3 "$OUT/video.frame")

	rm -rf "$OUT"; mkdir -p "$OUT"
	# Same fill but leave PAGE2 on at capture time
	asm_video \
		'STA $C001' 'STA $C00D' \
		'STA $C055' \
		'LDA #$C1' 'LDX #$00' \
		'fill_a2: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_a2' \
		'STA $C054' \
		'LDA #$C1' 'LDX #$00' \
		'fill_m2: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_m2' \
		'STA $C055' \
		'.halt'
	[[ $status -eq 0 ]]
	local page2_on_pixels
	page2_on_pixels=$(tail -n +3 "$OUT/video.frame")

	[[ "$page2_off_pixels" == "$page2_on_pixels" ]]
}

# ---------------------------------------------------------------------------
# Monochrome modes (6.5)
# ---------------------------------------------------------------------------

@test "80-col green screen renders white pixels as 98FF98" {
	asm_video_mono green \
		'STA $C001' 'STA $C00D' \
		'STA $C055' \
		'LDA #$C1' 'LDX #$00' \
		'fill_a: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_a' \
		'STA $C054' \
		'LDA #$C1' 'LDX #$00' \
		'fill_m: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_m' \
		'.halt'
	[[ $status -eq 0 ]]
	local rgb
	rgb=$(legend_rgb)
	echo "$rgb" | grep -q '98FF98'
}

@test "80-col amber screen renders white pixels as FFBF00" {
	asm_video_mono amber \
		'STA $C001' 'STA $C00D' \
		'STA $C055' \
		'LDA #$C1' 'LDX #$00' \
		'fill_a: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_a' \
		'STA $C054' \
		'LDA #$C1' 'LDX #$00' \
		'fill_m: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill_m' \
		'.halt'
	[[ $status -eq 0 ]]
	local rgb
	rgb=$(legend_rgb)
	echo "$rgb" | grep -q 'FFBF00'
}
