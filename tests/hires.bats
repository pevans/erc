setup_file() { load hires_helper; setup_file; }
setup()      { load hires_helper; setup; }
teardown()   { load hires_helper; teardown; }

# ---------------------------------------------------------------------------
# Screen geometry (2.2)
# ---------------------------------------------------------------------------

@test "hires frame is 560x384" {
	asm_hires_fill 00
	[[ $status -eq 0 ]]
	grep -q '^step 49999: video screen 560x384$' "$OUT/video.frame"
}

@test "single dot is 2 pixels wide" {
	asm_hires_fill 00 'LDA #$01' 'STA $2000'
	[[ $status -eq 0 ]]
	# Dot 0 occupies pixels 0-1 (2 wide), dot 2 onward is black
	col_range_is_uniform 0 0 2
	local row
	row=$(pixel_row 0)
	[[ "${row:0:1}" != "${row:4:1}" ]]
}

# ---------------------------------------------------------------------------
# Memory region (3.1)
# ---------------------------------------------------------------------------

@test "hires page 1 uses \$2000-\$3FFF" {
	asm_hires_fill 7F
	[[ $status -eq 0 ]]
	grep -q 'FFFFFF' "$OUT/video.frame"
}

# ---------------------------------------------------------------------------
# Address interleaving (3.2, 3.2.1)
# ---------------------------------------------------------------------------

@test "address \$2000 maps to hires row 0: pixel rows 0-1" {
	asm_hires_fill 00 'LDA #$7F' 'STA $2000'
	[[ $status -eq 0 ]]
	rows_differ 0 2    # row 0 changed vs row 2 (black)
	rows_same 0 1      # 2x vertical scale
	rows_same 2 3      # row 2 unchanged
}

@test "address \$2080 maps to hires row 8: pixel rows 16-17" {
	asm_hires_fill 00 'LDA #$7F' 'STA $2080'
	[[ $status -eq 0 ]]
	rows_same 14 15     # row before unchanged
	rows_differ 16 14   # row 16 changed
	rows_same 16 17     # 2x scale
	rows_same 18 14     # row after unchanged
}

@test "address \$2028 maps to hires row 64: pixel rows 128-129" {
	asm_hires_fill 00 'LDA #$7F' 'STA $2028'
	[[ $status -eq 0 ]]
	rows_same 126 127   # row before unchanged
	rows_differ 128 126  # row 128 changed
	rows_same 128 129    # 2x scale
	rows_same 130 126    # row after unchanged
}

@test "address \$2050 maps to hires row 128: pixel rows 256-257" {
	asm_hires_fill 00 'LDA #$7F' 'STA $2050'
	[[ $status -eq 0 ]]
	rows_same 254 255   # row before unchanged
	rows_differ 256 254  # row 256 changed
	rows_same 256 257    # 2x scale
	rows_same 258 254    # row after unchanged
}

@test "writing to hole byte \$2078 does not affect the hires display" {
	asm_hires_fill 00 'LDA #$7F' 'STA $2078'
	[[ $status -eq 0 ]]
	[[ $(color_count) -eq 1 ]]
}

# ---------------------------------------------------------------------------
# Byte layout (3.3)
# ---------------------------------------------------------------------------

@test "bit 0 controls leftmost dot in byte" {
	asm_hires_fill 00 'LDA #$01' 'STA $2000'
	[[ $status -eq 0 ]]
	# Dot 0 (pixels 0-1) should be colored; dot 2+ (pixels 4+) should be black
	col_range_is_uniform 0 0 2
	local row
	row=$(pixel_row 0)
	[[ "${row:0:1}" != "${row:4:1}" ]]
}

@test "bit 6 controls rightmost dot in byte" {
	# $40 = palette 0, only bit 6 (dot 6) on
	asm_hires_fill 00 'LDA #$40' 'STA $2000'
	[[ $status -eq 0 ]]
	# Dot 6 = pixels 12-13; dots 0-5 = pixels 0-11 should be black
	local row
	row=$(pixel_row 0)
	col_range_is_uniform 0 0 12
	[[ "${row:0:1}" != "${row:12:1}" ]]
}

@test "bit 7 selects palette without adding a dot" {
	# $80 = palette 1, all dots off; $00 = palette 0, all dots off.
	# Both should be all black -- palette bit does not create visible dots.
	asm_hires_fill 80
	[[ $status -eq 0 ]]
	[[ $(color_count) -eq 1 ]]
	[[ "$(legend_rgb)" == "000000" ]]

	rm -rf "$OUT"; mkdir -p "$OUT"

	asm_hires_fill 00
	[[ $status -eq 0 ]]
	[[ $(color_count) -eq 1 ]]
	[[ "$(legend_rgb)" == "000000" ]]
}

# ---------------------------------------------------------------------------
# Color generation (4.1, 4.2)
# ---------------------------------------------------------------------------

@test "all dots off produces black" {
	asm_hires_fill 00
	[[ $status -eq 0 ]]
	[[ $(color_count) -eq 1 ]]
	[[ "$(legend_rgb)" == "000000" ]]
}

@test "all dots on with palette 0 produces white" {
	asm_hires_fill 7F
	[[ $status -eq 0 ]]
	grep -q 'FFFFFF' "$OUT/video.frame"
}

@test "\$7F and \$00 produce different colors" {
	asm_hires_fill 7F
	[[ $status -eq 0 ]]
	local on_legend
	on_legend=$(sed -n '2p' "$OUT/video.frame")

	rm -rf "$OUT"; mkdir -p "$OUT"

	asm_hires_fill 00
	[[ $status -eq 0 ]]
	local off_legend
	off_legend=$(sed -n '2p' "$OUT/video.frame")

	[[ "$on_legend" != "$off_legend" ]]
}

@test "palette 0 isolated dot produces purple and green" {
	# $01 per byte: dot 0 on, dots 1-6 off, palette 0
	# Dot 0 (on, prev off) -> purple; dot 1 (off, prev on) -> green (fringing)
	asm_hires_fill 01
	[[ $status -eq 0 ]]
	local rgb
	rgb=$(legend_rgb)
	echo "$rgb" | grep -q 'D043E5'  # purple
	echo "$rgb" | grep -q '2FBC1A'  # green
}

@test "palette 1 isolated dot produces blue and orange" {
	# $81 per byte: dot 0 on, dots 1-6 off, palette 1
	# Dot 0 (on, prev off) -> blue; dot 1 (off, prev on) -> orange (fringing)
	asm_hires_fill 81
	[[ $status -eq 0 ]]
	local rgb
	rgb=$(legend_rgb)
	echo "$rgb" | grep -q '2F95E5'  # blue
	echo "$rgb" | grep -q 'D06A1A'  # orange
}

@test "even-position dot gets first palette color, odd gets second" {
	# Dot 0 (position 0, even) on in palette 0 -> purple
	asm_hires_fill 00 'LDA #$01' 'STA $2000'
	[[ $status -eq 0 ]]
	[[ "$(pixel_color 0 0)" == "D043E5" ]]

	rm -rf "$OUT"; mkdir -p "$OUT"

	# Dot 1 (position 1, odd) on in palette 0 -> green
	asm_hires_fill 00 'LDA #$02' 'STA $2000'
	[[ $status -eq 0 ]]
	[[ "$(pixel_color 0 2)" == "2FBC1A" ]]
}

@test "palette 0 and palette 1 produce different colors" {
	asm_hires_fill 01
	[[ $status -eq 0 ]]
	local pal0_legend
	pal0_legend=$(sed -n '2p' "$OUT/video.frame")

	rm -rf "$OUT"; mkdir -p "$OUT"

	asm_hires_fill 81
	[[ $status -eq 0 ]]
	local pal1_legend
	pal1_legend=$(sed -n '2p' "$OUT/video.frame")

	[[ "$pal0_legend" != "$pal1_legend" ]]
}

# ---------------------------------------------------------------------------
# Color rules (4.3, 4.4)
# ---------------------------------------------------------------------------

@test "consecutive on bits produce white" {
	# $03 = palette 0, dots 0 and 1 on. Dot 1 (on, prev on) -> white.
	asm_hires_fill 00 'LDA #$03' 'STA $2000'
	[[ $status -eq 0 ]]
	local rgb
	rgb=$(legend_rgb)
	echo "$rgb" | grep -q 'FFFFFF'
}

@test "first dot in a row behaves as if previous dot were off" {
	# Last byte of row 0 ($2027) all on, first byte of row 1 ($2400) dot 0 on.
	# If row state carried over, dot 0 of row 1 would be white (consecutive on).
	# Row-initial behavior means it should be colored, not white.
	asm_hires_fill 00 \
		'LDA #$7F' 'STA $2027' \
		'LDA #$01' 'STA $2400'
	[[ $status -eq 0 ]]
	# Row 1 = pixel rows 2-3. Dot 0 = pixel column 0.
	[[ "$(pixel_color 2 0)" != "FFFFFF" ]]
}

@test "consecutive on bits produce white with palette 1" {
	# $83 = palette 1, dots 0+1 on -> white regardless of palette
	asm_hires_fill 00 'LDA #$83' 'STA $2000'
	[[ $status -eq 0 ]]
	local rgb
	rgb=$(legend_rgb)
	echo "$rgb" | grep -q 'FFFFFF'
}

@test "isolated on bit produces colored dot, not white" {
	# $01 = palette 0, only dot 0 on. No consecutive on bits -> no white.
	asm_hires_fill 00 'LDA #$01' 'STA $2000'
	[[ $status -eq 0 ]]
	local rgb
	rgb=$(legend_rgb)
	! echo "$rgb" | grep -q 'FFFFFF'
}

# ---------------------------------------------------------------------------
# Byte boundary effects (5.2, 5.3)
# ---------------------------------------------------------------------------

@test "boundary shift produces dark colors when right byte is palette 0" {
	# Byte 0: palette 1 all off ($80), byte 1: palette 0 dot 0 on ($01)
	# Right dot at position 7 (odd) -> green, shifts to dark green
	asm_hires_fill 00 'LDA #$80' 'STA $2000' 'LDA #$01' 'STA $2001'
	[[ $status -eq 0 ]]
	local rgb
	rgb=$(legend_rgb)
	echo "$rgb" | grep -q '3F4C12'
}

@test "boundary shift produces light colors when right byte is palette 1 and left is non-black" {
	# Byte 0: palette 0 dot 6 on ($40, purple at boundary)
	# Byte 1: palette 1 all off ($80, dot 0 gets blue from fringing)
	# Blue shifts to light purple since left is non-black
	asm_hires_fill 00 'LDA #$40' 'STA $2000' 'LDA #$80' 'STA $2001'
	[[ $status -eq 0 ]]
	local rgb
	rgb=$(legend_rgb)
	echo "$rgb" | grep -q 'BBAFF6'
}

@test "no boundary shift when left boundary dot is black and right is palette 1" {
	# Byte 0: palette 0 all off ($00, black boundary)
	# Byte 1: palette 1 dot 0 on ($81, orange at position 7)
	# Left is black -> no shift; should see regular orange, not light green
	asm_hires_fill 00 'LDA #$00' 'STA $2000' 'LDA #$81' 'STA $2001'
	[[ $status -eq 0 ]]
	local rgb
	rgb=$(legend_rgb)
	echo "$rgb" | grep -q 'D06A1A'
	! echo "$rgb" | grep -q 'BDEA86'
}

@test "adjacent bytes with different palettes produce boundary shift colors" {
	# Baseline: both bytes palette 0
	asm_hires_fill 00 \
		'LDA #$55' 'STA $2000' 'STA $2001'
	[[ $status -eq 0 ]]
	local same_palette
	same_palette=$(pixel_row 0)

	rm -rf "$OUT"; mkdir -p "$OUT"

	# Different palettes at byte boundary: palette 0 then palette 1
	asm_hires_fill 00 \
		'LDA #$55' 'STA $2000' 'LDA #$D5' 'STA $2001'
	[[ $status -eq 0 ]]
	local diff_palette
	diff_palette=$(pixel_row 0)

	[[ "$same_palette" != "$diff_palette" ]]
}

# ---------------------------------------------------------------------------
# Monochrome colors (6.2)
# ---------------------------------------------------------------------------

@test "green screen renders on-dots as 98FF98" {
	asm_hires_fill_mono green 7F
	[[ $status -eq 0 ]]
	local rgb
	rgb=$(legend_rgb)
	echo "$rgb" | grep -q '98FF98'
}

@test "amber screen renders on-dots as FFBF00" {
	asm_hires_fill_mono amber 7F
	[[ $status -eq 0 ]]
	local rgb
	rgb=$(legend_rgb)
	echo "$rgb" | grep -q 'FFBF00'
}

# ---------------------------------------------------------------------------
# Monochrome rendering behavior (6.3)
# ---------------------------------------------------------------------------

@test "monochrome ignores palette bit" {
	# Palette 0 ($7F) and palette 1 ($FF) should look identical in mono
	asm_hires_fill_mono green 7F
	[[ $status -eq 0 ]]
	local pal0_pixels
	pal0_pixels=$(tail -n +3 "$OUT/video.frame")

	rm -rf "$OUT"; mkdir -p "$OUT"

	asm_hires_fill_mono green FF
	[[ $status -eq 0 ]]
	local pal1_pixels
	pal1_pixels=$(tail -n +3 "$OUT/video.frame")

	[[ "$pal0_pixels" == "$pal1_pixels" ]]
}

@test "monochrome isolated dot has no color fringing" {
	# $01 in color mode produces two colors (dot + fringe).
	# In monochrome, only on-dots light up -- no fringing artifact.
	asm_hires_fill_mono green 01
	[[ $status -eq 0 ]]
	local rgb
	rgb=$(legend_rgb)
	# Should have exactly black and the green mono color -- no third color
	[[ $(echo "$rgb" | wc -w) -eq 2 ]]
	echo "$rgb" | grep -q '000000'
	echo "$rgb" | grep -q '98FF98'
}

# ---------------------------------------------------------------------------
# Mode dispatch (7.1)
# ---------------------------------------------------------------------------

@test "TEXT off and HIRES on activates hires mode" {
	# Fill hires with $7F (white), clear text page
	asm_video \
		'LDA #$00' 'LDX #$00' \
		'clrt: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE clrt' \
		"LDA #\$7F" 'STA $02' \
		'LDA #$00' 'STA $00' 'LDA #$20' 'STA $01' \
		'LDA $02' 'LDY #$00' \
		'hfill: STA ($00),Y' 'INY' 'BNE hfill' \
		'INC $01' 'LDX $01' 'CPX #$40' 'BNE hfill' \
		'STA $C050' 'STA $C057' '.halt'
	[[ $status -eq 0 ]]
	grep -q 'FFFFFF' "$OUT/video.frame"
}

@test "TEXT on returns from hires to text mode" {
	# Enter hires, then TEXT on -> text mode
	asm_hires_fill 7F 'STA $C051'
	[[ $status -eq 0 ]]
	local roundtrip_pixels
	roundtrip_pixels=$(tail -n +3 "$OUT/video.frame")

	rm -rf "$OUT"; mkdir -p "$OUT"

	# Reference: plain text mode
	asm_video '.halt'
	[[ $status -eq 0 ]]
	local text_pixels
	text_pixels=$(tail -n +3 "$OUT/video.frame")

	[[ "$roundtrip_pixels" == "$text_pixels" ]]
}

@test "HIRES off returns from hires to lores mode" {
	# Enter hires, then HIRES off -> lores (TEXT still off)
	asm_hires_fill 7F 'STA $C056'
	[[ $status -eq 0 ]]
	local hires_off_pixels
	hires_off_pixels=$(tail -n +3 "$OUT/video.frame")

	rm -rf "$OUT"; mkdir -p "$OUT"

	# Reference: lores mode (TEXT off, HIRES off)
	asm_video 'STA $C050' '.halt'
	[[ $status -eq 0 ]]
	local lores_pixels
	lores_pixels=$(tail -n +3 "$OUT/video.frame")

	[[ "$hires_off_pixels" == "$lores_pixels" ]]
}

# ---------------------------------------------------------------------------
# Soft switches (8.1)
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

# ---------------------------------------------------------------------------
# Read switches (8.2)
# ---------------------------------------------------------------------------

@test "RDPAGE2 returns bit 7 high when PAGE2 is on" {
	asm_mem 0300 \
		'STA $C055' \
		'LDA $C01C' 'STA $0300' \
		'.halt'
	[[ $status -eq 0 ]]
	local val
	val=$(last_mem_val 0300)
	[[ $((16#$val)) -ge 128 ]]
}

@test "RDPAGE2 returns bit 7 low when PAGE2 is off" {
	asm_mem 0300 \
		'LDA #$FF' 'STA $0300' \
		'LDA $C01C' 'STA $0300' \
		'.halt'
	[[ $status -eq 0 ]]
	local val
	val=$(last_mem_val 0300)
	[[ $((16#$val)) -lt 128 ]]
}

@test "RDMIXED returns bit 7 high when MIXED is on" {
	asm_mem 0300 \
		'STA $C053' \
		'LDA $C01B' 'STA $0300' \
		'.halt'
	[[ $status -eq 0 ]]
	local val
	val=$(last_mem_val 0300)
	[[ $((16#$val)) -ge 128 ]]
}

@test "RDMIXED returns bit 7 low when MIXED is off" {
	asm_mem 0300 \
		'LDA #$FF' 'STA $0300' \
		'LDA $C01B' 'STA $0300' \
		'.halt'
	[[ $status -eq 0 ]]
	local val
	val=$(last_mem_val 0300)
	[[ $((16#$val)) -lt 128 ]]
}

@test "RDHIRES returns bit 7 high when HIRES is on" {
	asm_mem 0300 \
		'STA $C057' \
		'LDA $C01D' 'STA $0300' \
		'.halt'
	[[ $status -eq 0 ]]
	local val
	val=$(last_mem_val 0300)
	[[ $((16#$val)) -ge 128 ]]
}

@test "RDHIRES returns bit 7 low when HIRES is off" {
	asm_mem 0300 \
		'LDA #$FF' 'STA $0300' \
		'LDA $C01D' 'STA $0300' \
		'.halt'
	[[ $status -eq 0 ]]
	local val
	val=$(last_mem_val 0300)
	[[ $((16#$val)) -lt 128 ]]
}

# ---------------------------------------------------------------------------
# Page selection (8.3)
# ---------------------------------------------------------------------------

@test "80STORE routes hires writes to aux when PAGE2 is on" {
	# With 80STORE+HIRES, PAGE2 controls main vs aux for $2000-$3FFF.
	# Write $00 to main, $7F to aux, display with PAGE2 off (main) -> black
	asm_video \
		'STA $C001' \
		'STA $C057' \
		'STA $C050' \
		'STA $C054' \
		'LDA #$00' 'STA $2000' \
		'STA $C055' \
		'LDA #$7F' 'STA $2000' \
		'STA $C054' \
		'.halt'
	[[ $status -eq 0 ]]
	local main_row
	main_row=$(pixel_row 0)

	rm -rf "$OUT"; mkdir -p "$OUT"

	# Same writes, but display with PAGE2 on (aux) -> white at row 0
	asm_video \
		'STA $C001' \
		'STA $C057' \
		'STA $C050' \
		'STA $C054' \
		'LDA #$00' 'STA $2000' \
		'STA $C055' \
		'LDA #$7F' 'STA $2000' \
		'.halt'
	[[ $status -eq 0 ]]
	local aux_row
	aux_row=$(pixel_row 0)

	[[ "$main_row" != "$aux_row" ]]
}

@test "PAGE2 displays hires page 2 content" {
	# Fill page 1 ($2000-$3FFF) with $00, page 2 ($4000-$5FFF) with $7F
	asm_video_long \
		'LDA #$00' 'STA $02' \
		'LDA #$00' 'STA $00' 'LDA #$20' 'STA $01' \
		'LDA $02' 'LDY #$00' \
		'fp1: STA ($00),Y' 'INY' 'BNE fp1' \
		'INC $01' 'LDX $01' 'CPX #$40' 'BNE fp1' \
		'LDA #$7F' 'STA $02' \
		'LDA #$00' 'STA $00' 'LDA #$40' 'STA $01' \
		'LDA $02' 'LDY #$00' \
		'fp2: STA ($00),Y' 'INY' 'BNE fp2' \
		'INC $01' 'LDX $01' 'CPX #$60' 'BNE fp2' \
		'STA $C050' 'STA $C057' 'STA $C055' '.halt'
	[[ $status -eq 0 ]]
	grep -q 'FFFFFF' "$OUT/video.frame"
}

@test "PAGE2 switch changes hires displayed content" {
	# Run 1: hires page 1 with $7F
	asm_hires_fill 7F
	[[ $status -eq 0 ]]
	local p1_pixels
	p1_pixels=$(tail -n +3 "$OUT/video.frame")

	rm -rf "$OUT"; mkdir -p "$OUT"

	# Run 2: same page 1 fill but switch to page 2 (different content)
	asm_hires_fill 7F 'STA $C055'
	[[ $status -eq 0 ]]
	local p2_pixels
	p2_pixels=$(tail -n +3 "$OUT/video.frame")

	[[ "$p1_pixels" != "$p2_pixels" ]]
}
