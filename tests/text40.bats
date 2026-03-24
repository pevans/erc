setup_file() { load text40_helper; setup_file; }
setup()      { load text40_helper; setup; }
teardown()   { load text40_helper; teardown; }

# ---------------------------------------------------------------------------
# Screen geometry
# ---------------------------------------------------------------------------

@test "40-col text frame is 560x384" {
	asm_video \
		'LDA #$A0' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' '.halt'
	[[ $status -eq 0 ]]
	grep -q '^step 4999: video screen 560x384$' "$OUT/video.frame"
}

# ---------------------------------------------------------------------------
# Character rendering
# ---------------------------------------------------------------------------

@test "screen filled with spaces has one color" {
	asm_video \
		'LDA #$A0' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' '.halt'
	[[ $status -eq 0 ]]
	[[ $(color_count) -eq 1 ]]
}

@test "screen filled with normal A has two colors" {
	asm_video \
		'LDA #$C1' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' '.halt'
	[[ $status -eq 0 ]]
	[[ $(color_count) -eq 2 ]]
}

@test "screen filled with inverse A has two colors" {
	asm_video \
		'LDA #$01' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' '.halt'
	[[ $status -eq 0 ]]
	[[ $(color_count) -eq 2 ]]
}

@test "normal and inverse A use different color palettes" {
	asm_video \
		'LDA #$C1' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' '.halt'
	[[ $status -eq 0 ]]
	local normal_colors
	normal_colors=$(sed -n '2p' "$OUT/video.frame")

	rm -rf "$OUT"; mkdir -p "$OUT"
	asm_video \
		'LDA #$01' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' '.halt'
	[[ $status -eq 0 ]]
	local inverse_colors
	inverse_colors=$(sed -n '2p' "$OUT/video.frame")

	[[ "$normal_colors" != "$inverse_colors" ]]
}

# ---------------------------------------------------------------------------
# Interleaved address mapping
# ---------------------------------------------------------------------------

@test "address \$0400 maps to row 0: character in pixel rows 0-15" {
	asm_video \
		'LDA #$A0' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' \
		'LDA #$C1' 'STA $0400' \
		'.halt'
	[[ $status -eq 0 ]]
	any_row_has_mixed_colors 0 15
	row_is_uniform 16
}

@test "address \$0480 maps to row 1: character in pixel rows 16-31" {
	asm_video \
		'LDA #$A0' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' \
		'LDA #$C1' 'STA $0480' \
		'.halt'
	[[ $status -eq 0 ]]
	row_is_uniform 0
	any_row_has_mixed_colors 16 31
}

@test "address \$0428 maps to row 8: character in pixel rows 128-143" {
	asm_video \
		'LDA #$A0' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' \
		'LDA #$C1' 'STA $0428' \
		'.halt'
	[[ $status -eq 0 ]]
	row_is_uniform 0
	any_row_has_mixed_colors 128 143
}

@test "address \$0450 maps to row 16: character in pixel rows 256-271" {
	asm_video \
		'LDA #$A0' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' \
		'LDA #$C1' 'STA $0450' \
		'.halt'
	[[ $status -eq 0 ]]
	row_is_uniform 255
	any_row_has_mixed_colors 256 271
}

@test "address \$07D0 maps to row 23: character in pixel rows 368-383" {
	asm_video \
		'LDA #$A0' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' \
		'LDA #$C1' 'STA $07D0' \
		'.halt'
	[[ $status -eq 0 ]]
	row_is_uniform 367
	any_row_has_mixed_colors 368 383
}

@test "writing to hole byte \$0478 does not affect the display" {
	asm_video \
		'LDA #$A0' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' \
		'LDA #$C1' 'STA $0478' \
		'.halt'
	[[ $status -eq 0 ]]
	[[ $(color_count) -eq 1 ]]
}

# ---------------------------------------------------------------------------
# Page 2
# ---------------------------------------------------------------------------

@test "PAGE2 switch changes the displayed content" {
	# Run 1: fill page 1 with 'A', display page 1
	asm_video \
		'LDA #$C1' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' '.halt'
	[[ $status -eq 0 ]]
	cp "$OUT/video.frame" "$TMP/page1.frame"

	# Run 2: same page 1 fill, then switch to page 2
	rm -rf "$OUT"; mkdir -p "$OUT"
	asm_video \
		'LDA #$C1' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' \
		'STA $C055' '.halt'
	[[ $status -eq 0 ]]

	# Page 2 contains boot residue and program code, not 'A',
	# so the pixel patterns must differ
	local p1 p2
	p1=$(tail -n +3 "$TMP/page1.frame")
	p2=$(tail -n +3 "$OUT/video.frame")
	[[ "$p1" != "$p2" ]]
}

@test "PAGE2 with 80STORE on shows auxiliary memory" {
	# With 80STORE on and PAGE2 on, display reads from aux $0400-$07FF.
	# Fill aux with 'A', main with spaces, then switch to PAGE2.
	asm_video \
		'STA $C001' \
		'STA $C005' \
		'LDA #$C1' 'LDX #$00' \
		'fill1: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill1' \
		'STA $C004' \
		'LDA #$A0' 'LDX #$00' \
		'fill2: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill2' \
		'STA $C055' \
		'.halt'
	[[ $status -eq 0 ]]
	[[ $(color_count) -eq 2 ]]
}

# ---------------------------------------------------------------------------
# Soft switches
# ---------------------------------------------------------------------------

@test "toggling TEXT off logs DisplayText state change" {
	asm_state DisplayText \
		'STA $C050' '.halt'
	[[ $status -eq 0 ]]
	[[ -f "$OUT/state.log" ]]
	grep -q 'comp DisplayText' "$OUT/state.log"
}

@test "toggling PAGE2 on logs DisplayPage2 state change" {
	asm_state DisplayPage2 \
		'STA $C055' '.halt'
	[[ $status -eq 0 ]]
	[[ -f "$OUT/state.log" ]]
	grep -q 'comp DisplayPage2' "$OUT/state.log"
}

@test "toggling TEXT on logs DisplayText state change" {
	asm_state DisplayText \
		'STA $C050' 'STA $C051' '.halt'
	[[ $status -eq 0 ]]
	[[ -f "$OUT/state.log" ]]
	grep -q 'comp DisplayText' "$OUT/state.log"
}

@test "toggling MIXED on logs DisplayMixed state change" {
	asm_state DisplayMixed \
		'STA $C053' '.halt'
	[[ $status -eq 0 ]]
	[[ -f "$OUT/state.log" ]]
	grep -q 'comp DisplayMixed' "$OUT/state.log"
}

# ---------------------------------------------------------------------------
# Character encoding
# ---------------------------------------------------------------------------

@test "screen filled with flash A has two colors" {
	asm_video \
		'LDA #$41' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' '.halt'
	[[ $status -eq 0 ]]
	[[ $(color_count) -eq 2 ]]
}

@test "flash A and inverse A use the same color palette" {
	asm_video \
		'LDA #$41' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' '.halt'
	[[ $status -eq 0 ]]
	local flash_colors
	flash_colors=$(sed -n '2p' "$OUT/video.frame")

	rm -rf "$OUT"; mkdir -p "$OUT"
	asm_video \
		'LDA #$01' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' '.halt'
	[[ $status -eq 0 ]]
	local inverse_colors
	inverse_colors=$(sed -n '2p' "$OUT/video.frame")

	[[ "$flash_colors" == "$inverse_colors" ]]
}

@test "screen filled with lowercase a has two colors" {
	asm_video \
		'LDA #$E1' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' '.halt'
	[[ $status -eq 0 ]]
	[[ $(color_count) -eq 2 ]]
}

@test "lowercase and uppercase A use different pixel patterns" {
	asm_video \
		'LDA #$C1' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' '.halt'
	[[ $status -eq 0 ]]
	local upper_pixels
	upper_pixels=$(tail -n +3 "$OUT/video.frame")

	rm -rf "$OUT"; mkdir -p "$OUT"
	asm_video \
		'LDA #$E1' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' '.halt'
	[[ $status -eq 0 ]]
	local lower_pixels
	lower_pixels=$(tail -n +3 "$OUT/video.frame")

	[[ "$upper_pixels" != "$lower_pixels" ]]
}

# ---------------------------------------------------------------------------
# Soft switches (continued)
# ---------------------------------------------------------------------------

@test "toggling PAGE2 off logs DisplayPage2 state change" {
	asm_state DisplayPage2 \
		'STA $C055' 'STA $C054' '.halt'
	[[ $status -eq 0 ]]
	[[ -f "$OUT/state.log" ]]
	grep -q 'comp DisplayPage2' "$OUT/state.log"
}

@test "toggling MIXED off logs DisplayMixed state change" {
	asm_state DisplayMixed \
		'STA $C053' 'STA $C052' '.halt'
	[[ $status -eq 0 ]]
	[[ -f "$OUT/state.log" ]]
	grep -q 'comp DisplayMixed' "$OUT/state.log"
}

@test "toggling 80COL on logs DisplayCol80 state change" {
	asm_state DisplayCol80 \
		'STA $C00D' '.halt'
	[[ $status -eq 0 ]]
	[[ -f "$OUT/state.log" ]]
	grep -q 'comp DisplayCol80' "$OUT/state.log"
}

@test "toggling 80COL off logs DisplayCol80 state change" {
	asm_state DisplayCol80 \
		'STA $C00D' 'STA $C00C' '.halt'
	[[ $status -eq 0 ]]
	[[ -f "$OUT/state.log" ]]
	grep -q 'comp DisplayCol80' "$OUT/state.log"
}

@test "toggling 80STORE on logs DisplayStore80 state change" {
	asm_state DisplayStore80 \
		'STA $C001' '.halt'
	[[ $status -eq 0 ]]
	[[ -f "$OUT/state.log" ]]
	grep -q 'comp DisplayStore80' "$OUT/state.log"
}

@test "toggling 80STORE off logs DisplayStore80 state change" {
	asm_state DisplayStore80 \
		'STA $C001' 'STA $C000' '.halt'
	[[ $status -eq 0 ]]
	[[ -f "$OUT/state.log" ]]
	grep -q 'comp DisplayStore80' "$OUT/state.log"
}

@test "toggling ALTCHAR on logs DisplayAltChar state change" {
	asm_state DisplayAltChar \
		'STA $C00F' '.halt'
	[[ $status -eq 0 ]]
	[[ -f "$OUT/state.log" ]]
	grep -q 'comp DisplayAltChar' "$OUT/state.log"
}

@test "toggling ALTCHAR off logs DisplayAltChar state change" {
	asm_state DisplayAltChar \
		'STA $C00F' 'STA $C00E' '.halt'
	[[ $status -eq 0 ]]
	[[ -f "$OUT/state.log" ]]
	grep -q 'comp DisplayAltChar' "$OUT/state.log"
}

# ---------------------------------------------------------------------------
# Soft switch status reads
# ---------------------------------------------------------------------------

@test "reading \$C01A returns \$80 when TEXT is on" {
	asm_mem 0300 \
		'LDA #$00' 'STA $0300' \
		'LDA $C01A' 'STA $0300' \
		'.halt'
	[[ $status -eq 0 ]]
	grep -q 'mem $0300 $00 -> $80' "$OUT/state.log"
}

@test "reading \$C01A returns \$00 when TEXT is off" {
	asm_mem 0300 \
		'STA $C050' \
		'LDA #$FF' 'STA $0300' \
		'LDA $C01A' 'STA $0300' \
		'.halt'
	[[ $status -eq 0 ]]
	grep -q 'mem $0300 $FF -> $00' "$OUT/state.log"
}

@test "reading \$C01C returns \$80 when PAGE2 is on" {
	asm_mem 0300 \
		'STA $C055' \
		'LDA #$00' 'STA $0300' \
		'LDA $C01C' 'STA $0300' \
		'.halt'
	[[ $status -eq 0 ]]
	grep -q 'mem $0300 $00 -> $80' "$OUT/state.log"
}

@test "reading \$C01C returns \$00 when PAGE2 is off" {
	asm_mem 0300 \
		'LDA #$FF' 'STA $0300' \
		'LDA $C01C' 'STA $0300' \
		'.halt'
	[[ $status -eq 0 ]]
	grep -q 'mem $0300 $FF -> $00' "$OUT/state.log"
}

@test "reading \$C01B returns \$80 when MIXED is on" {
	asm_mem 0300 \
		'STA $C053' \
		'LDA #$00' 'STA $0300' \
		'LDA $C01B' 'STA $0300' \
		'.halt'
	[[ $status -eq 0 ]]
	grep -q 'mem $0300 $00 -> $80' "$OUT/state.log"
}

@test "reading \$C018 returns \$00 when 80STORE is off" {
	asm_mem 0300 \
		'LDA #$FF' 'STA $0300' \
		'LDA $C018' 'STA $0300' \
		'.halt'
	[[ $status -eq 0 ]]
	grep -q 'mem $0300 $FF -> $00' "$OUT/state.log"
}

@test "reading \$C018 returns \$80 when 80STORE is on" {
	asm_mem 0300 \
		'STA $C001' \
		'LDA #$00' 'STA $0300' \
		'LDA $C018' 'STA $0300' \
		'.halt'
	[[ $status -eq 0 ]]
	grep -q 'mem $0300 $00 -> $80' "$OUT/state.log"
}

@test "reading \$C01B returns \$00 when MIXED is off" {
	asm_mem 0300 \
		'LDA #$FF' 'STA $0300' \
		'LDA $C01B' 'STA $0300' \
		'.halt'
	[[ $status -eq 0 ]]
	grep -q 'mem $0300 $FF -> $00' "$OUT/state.log"
}

@test "reading \$C01E returns \$00 when ALTCHAR is off" {
	asm_mem 0300 \
		'LDA #$FF' 'STA $0300' \
		'LDA $C01E' 'STA $0300' \
		'.halt'
	[[ $status -eq 0 ]]
	grep -q 'mem $0300 $FF -> $00' "$OUT/state.log"
}

@test "reading \$C01E returns \$80 when ALTCHAR is on" {
	asm_mem 0300 \
		'STA $C00F' \
		'LDA #$00' 'STA $0300' \
		'LDA $C01E' 'STA $0300' \
		'.halt'
	[[ $status -eq 0 ]]
	grep -q 'mem $0300 $00 -> $80' "$OUT/state.log"
}

@test "reading \$C01F returns \$00 when 80COL is off" {
	asm_mem 0300 \
		'LDA #$FF' 'STA $0300' \
		'LDA $C01F' 'STA $0300' \
		'.halt'
	[[ $status -eq 0 ]]
	grep -q 'mem $0300 $FF -> $00' "$OUT/state.log"
}

@test "reading \$C01F returns \$80 when 80COL is on" {
	asm_mem 0300 \
		'STA $C00D' \
		'LDA #$00' 'STA $0300' \
		'LDA $C01F' 'STA $0300' \
		'.halt'
	[[ $status -eq 0 ]]
	grep -q 'mem $0300 $00 -> $80' "$OUT/state.log"
}

# ---------------------------------------------------------------------------
# Vertical blank
# ---------------------------------------------------------------------------

@test "reading \$C019 returns \$80 during vertical blank" {
	local src="$TMP/test.s"
	printf '%s\n' \
		'LDA #$00' \
		'STA $0300' \
		'loop: LDA $C019' \
		'BPL loop' \
		'STA $0300' \
		'.halt' \
		>"$src"
	"$ASSEMBLER" -o "$TMP/test.dsk" "$src" || return 1
	run "$ERC" headless \
		--output "$OUT" \
		--steps 10000 \
		--start-at 0801 \
		--watch-mem 0300 \
		"$TMP/test.dsk"
	[[ $status -eq 0 ]]
	grep -q 'mem $0300 $00 -> $80' "$OUT/state.log"
}

# ---------------------------------------------------------------------------
# Column position
# ---------------------------------------------------------------------------

@test "address \$0401 maps to column 1: character in pixel columns 14-27" {
	asm_video \
		'LDA #$A0' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' \
		'LDA #$C1' 'STA $0401' \
		'.halt'
	[[ $status -eq 0 ]]
	col_range_is_uniform 0 0 14
	col_range_has_mixed_colors 0 14 14
}
