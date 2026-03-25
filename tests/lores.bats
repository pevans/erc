setup_file() { load lores_helper; setup_file; }
setup()      { load lores_helper; setup; }
teardown()   { load lores_helper; teardown; }

# ---------------------------------------------------------------------------
# Screen geometry
# ---------------------------------------------------------------------------

@test "lores frame is 560x384" {
	asm_video \
		'LDA #$00' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' \
		'STA $C050' '.halt'
	[[ $status -eq 0 ]]
	grep -q '^step 4999: video screen 560x384$' "$OUT/video.frame"
}

@test "single block is 14 pixels wide" {
	# Black background, one white block at row 0 col 0
	asm_video \
		'LDA #$00' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' \
		'LDA #$FF' 'STA $0400' \
		'STA $C050' '.halt'
	[[ $status -eq 0 ]]
	# $0400 = row 0, col 0. White block occupies pixels 0-13 (14 wide).
	col_range_is_uniform 0 0 14
	col_range_is_uniform 0 14 546
	local row
	row=$(pixel_row 0)
	[[ "${row:0:1}" != "${row:14:1}" ]]
}

# ---------------------------------------------------------------------------
# Mode activation
# ---------------------------------------------------------------------------

@test "TEXT off changes display from text to lores" {
	# Run 1: text mode (default) -- 'A' glyphs rendered
	asm_video \
		'LDA #$C1' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' '.halt'
	[[ $status -eq 0 ]]
	local text_pixels
	text_pixels=$(tail -n +3 "$OUT/video.frame")

	rm -rf "$OUT"; mkdir -p "$OUT"

	# Run 2: lores mode -- same memory, but TEXT off
	asm_video \
		'LDA #$C1' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' \
		'STA $C050' '.halt'
	[[ $status -eq 0 ]]
	local lores_pixels
	lores_pixels=$(tail -n +3 "$OUT/video.frame")

	[[ "$text_pixels" != "$lores_pixels" ]]
}

@test "TEXT on returns from lores to text mode" {
	# Enter lores (TEXT off), then back to text (TEXT on)
	asm_video \
		'LDA #$C1' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' \
		'STA $C050' 'STA $C051' '.halt'
	[[ $status -eq 0 ]]
	local roundtrip_pixels
	roundtrip_pixels=$(tail -n +3 "$OUT/video.frame")

	rm -rf "$OUT"; mkdir -p "$OUT"

	# Reference: plain text mode with same memory
	asm_video \
		'LDA #$C1' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' '.halt'
	[[ $status -eq 0 ]]
	local text_pixels
	text_pixels=$(tail -n +3 "$OUT/video.frame")

	[[ "$roundtrip_pixels" == "$text_pixels" ]]
}

@test "HIRES on overrides lores mode" {
	# Run 1: lores mode (TEXT off, HIRES off) -- page 1 filled with $FF = all white
	asm_video \
		'LDA #$FF' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' \
		'STA $C050' '.halt'
	[[ $status -eq 0 ]]
	local lores_frame
	lores_frame=$(tail -n +2 "$OUT/video.frame")

	rm -rf "$OUT"; mkdir -p "$OUT"

	# Run 2: hires mode (TEXT off, HIRES on) -- hires uses $2000-$3FFF, not page 1
	asm_video \
		'LDA #$FF' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' \
		'STA $C050' 'STA $C057' '.halt'
	[[ $status -eq 0 ]]
	local hires_frame
	hires_frame=$(tail -n +2 "$OUT/video.frame")

	[[ "$lores_frame" != "$hires_frame" ]]
}

# ---------------------------------------------------------------------------
# Color rendering
# ---------------------------------------------------------------------------

@test "all 16 color indices produce correct RGB values" {
	local -a expected=(
		000000 901740 402CA5 D043E5 006940 808080 2F95E5 BFABFF
		405400 D06A1A 808080 FF96BF 2FBC1A BFD35A 6FE8BF FFFFFF
	)

	local idx fill_hex rgb
	for idx in {0..15}; do
		rm -rf "$OUT"; mkdir -p "$OUT"
		printf -v fill_hex '%X%X' "$idx" "$idx"
		asm_video \
			"LDA #\$$fill_hex" 'LDX #$00' \
			'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
			'INX' 'BNE fill' \
			'STA $C050' '.halt'
		[[ $status -eq 0 ]]
		rgb=$(legend_rgb)
		[[ "$rgb" == "${expected[$idx]}" ]] || {
			echo "color $idx: expected ${expected[$idx]}, got $rgb" >&2
			return 1
		}
	done
}

@test "screen filled with \$00 has one color" {
	asm_video \
		'LDA #$00' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' \
		'STA $C050' '.halt'
	[[ $status -eq 0 ]]
	[[ $(color_count) -eq 1 ]]
}

@test "screen filled with \$FF has one color" {
	asm_video \
		'LDA #$FF' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' \
		'STA $C050' '.halt'
	[[ $status -eq 0 ]]
	[[ $(color_count) -eq 1 ]]
}

@test "\$00 and \$FF screens produce different colors" {
	# The legend line encodes actual RGB values, so single-color fills of
	# different colors are distinguishable even though both produce all-'!'
	# pixel rows.
	asm_video \
		'LDA #$00' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' \
		'STA $C050' '.halt'
	[[ $status -eq 0 ]]
	local black_legend
	black_legend=$(sed -n '2p' "$OUT/video.frame")

	rm -rf "$OUT"; mkdir -p "$OUT"

	asm_video \
		'LDA #$FF' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' \
		'STA $C050' '.halt'
	[[ $status -eq 0 ]]
	local white_legend
	white_legend=$(sed -n '2p' "$OUT/video.frame")

	[[ "$black_legend" != "$white_legend" ]]
}

# ---------------------------------------------------------------------------
# Nibble encoding
# ---------------------------------------------------------------------------

@test "byte \$F0 produces two colors (top black, bottom white)" {
	# low nibble = 0 (black), high nibble = 15 (white)
	asm_video \
		'LDA #$F0' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' \
		'STA $C050' '.halt'
	[[ $status -eq 0 ]]
	[[ $(color_count) -eq 2 ]]
}

@test "low nibble sets top block color" {
	# $F0: low nibble=0=black (top), high nibble=15=white (bottom)
	# Rows 0-7 (top block of first text row) should be black.
	# Rows 8-15 (bottom block) should be white.
	asm_video \
		'LDA #$F0' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' \
		'STA $C050' '.halt'
	[[ $status -eq 0 ]]
	rows_same 0 7     # entire top block is one color (black)
	rows_differ 7 8   # boundary between top and bottom block
	rows_same 8 15    # entire bottom block is one color (white)
}

@test "high nibble sets bottom block color" {
	# $0F: low nibble=15=white (top), high nibble=0=black (bottom)
	asm_video \
		'LDA #$0F' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' \
		'STA $C050' '.halt'
	[[ $status -eq 0 ]]
	rows_differ 7 8   # boundary: row 7 is white (top), row 8 is black (bottom)
}

@test "\$F0 and \$0F produce different pixel patterns" {
	# $F0 has top=black/bottom=white; $0F has top=white/bottom=black.
	# The full frame (legend + pixels) differs between the two because the
	# color-to-character assignment reflects scan order, so the legends differ.
	asm_video \
		'LDA #$F0' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' \
		'STA $C050' '.halt'
	[[ $status -eq 0 ]]
	local f0_frame
	f0_frame=$(tail -n +2 "$OUT/video.frame")

	rm -rf "$OUT"; mkdir -p "$OUT"

	asm_video \
		'LDA #$0F' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' \
		'STA $C050' '.halt'
	[[ $status -eq 0 ]]
	local of_frame
	of_frame=$(tail -n +2 "$OUT/video.frame")

	[[ "$f0_frame" != "$of_frame" ]]
}

# ---------------------------------------------------------------------------
# Interleaved address mapping
# ---------------------------------------------------------------------------

@test "address \$0400 maps to lores rows 0-1: pixel rows 0-15" {
	asm_video \
		'LDA #$00' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' \
		'LDA #$FF' 'STA $0400' \
		'STA $C050' '.halt'
	[[ $status -eq 0 ]]
	rows_differ 0 383    # row 0 changed (white vs black)
	rows_differ 15 383   # row 15 changed (white vs black)
	rows_same 16 383     # row 16 unchanged
}

@test "address \$0480 maps to lores rows 2-3: pixel rows 16-31" {
	asm_video \
		'LDA #$00' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' \
		'LDA #$FF' 'STA $0480' \
		'STA $C050' '.halt'
	[[ $status -eq 0 ]]
	rows_same 0 383      # row 0 unchanged
	rows_differ 16 383   # row 16 changed (white vs black)
	rows_same 32 383     # row 32 unchanged
}

@test "address \$0428 maps to lores rows 16-17: pixel rows 128-143" {
	asm_video \
		'LDA #$00' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' \
		'LDA #$FF' 'STA $0428' \
		'STA $C050' '.halt'
	[[ $status -eq 0 ]]
	rows_same 127 383    # row just before range unchanged
	rows_differ 128 383  # first row of range changed
	rows_same 144 383    # row just after range unchanged
}

@test "address \$07D0 maps to lores rows 46-47: pixel rows 368-383" {
	asm_video \
		'LDA #$00' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' \
		'LDA #$FF' 'STA $07D0' \
		'STA $C050' '.halt'
	[[ $status -eq 0 ]]
	rows_same 367 0      # row just before range unchanged (use row 0 as ref)
	rows_differ 368 0    # first row of range changed
	rows_differ 383 0    # last row of range changed
}

@test "writing to hole byte \$0478 does not affect the lores display" {
	asm_video \
		'LDA #$00' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' \
		'LDA #$FF' 'STA $0478' \
		'STA $C050' '.halt'
	[[ $status -eq 0 ]]
	[[ $(color_count) -eq 1 ]]
}

# ---------------------------------------------------------------------------
# Page 2
# ---------------------------------------------------------------------------

@test "PAGE2 displays explicitly filled page 2 content" {
	# Page 1 = all black. Page 2 $0900-$0BFF = all white.
	# Avoids filling $0800 range where program code resides at $0801+.
	asm_video \
		'LDA #$00' 'LDX #$00' \
		'clr1: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE clr1' \
		'LDA #$FF' 'LDX #$00' \
		'fill2: STA $0900,X' 'STA $0A00,X' 'STA $0B00,X' \
		'INX' 'BNE fill2' \
		'STA $C050' 'STA $C055' '.halt'
	[[ $status -eq 0 ]]
	# White on screen proves page 2 is rendered, not page 1 (all black).
	grep -q 'FFFFFF' "$OUT/video.frame"
	# $0900 = text row 2 = pixel rows 32-39; should be uniformly white.
	row_is_uniform 32
}

@test "PAGE2 switch changes lores displayed content" {
	# Run 1: page 1 -- fill with $FF (white), lores mode
	asm_video \
		'LDA #$FF' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' \
		'STA $C050' '.halt'
	[[ $status -eq 0 ]]
	local p1_pixels
	p1_pixels=$(tail -n +3 "$OUT/video.frame")

	rm -rf "$OUT"; mkdir -p "$OUT"

	# Run 2: switch to page 2 -- page 2 has boot residue, not $FF
	asm_video \
		'LDA #$FF' 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' \
		'STA $C050' 'STA $C055' '.halt'
	[[ $status -eq 0 ]]
	local p2_pixels
	p2_pixels=$(tail -n +3 "$OUT/video.frame")

	[[ "$p1_pixels" != "$p2_pixels" ]]
}

# ---------------------------------------------------------------------------
# Soft switches
# ---------------------------------------------------------------------------

@test "toggling HIRES on logs DisplayHires state change" {
	asm_state DisplayHires \
		'STA $C057' '.halt'
	[[ $status -eq 0 ]]
	[[ -f "$OUT/state.log" ]]
	grep -q 'comp DisplayHires' "$OUT/state.log"
}

@test "toggling HIRES off logs DisplayHires state change" {
	asm_state DisplayHires \
		'STA $C057' 'STA $C056' '.halt'
	[[ $status -eq 0 ]]
	[[ -f "$OUT/state.log" ]]
	grep -q 'comp DisplayHires' "$OUT/state.log"
}
