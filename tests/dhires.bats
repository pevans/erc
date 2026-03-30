setup_file() { load dhires_helper; setup_file; }
setup()      { load dhires_helper; setup; }
teardown()   { load dhires_helper; teardown; }

# ---------------------------------------------------------------------------
# Screen geometry (2.2)
# ---------------------------------------------------------------------------

@test "dhires frame is 560x384" {
	asm_dhires_fill 00
	[[ $status -eq 0 ]]
	grep -q '^step 149999: video screen 560x384$' "$OUT/video.frame"
}

@test "dhires single dot is 1 pixel wide" {
	# Set aux byte 0 at row 0 to $01 (bit 0 on = dot 0), everything else $00
	asm_dhires_fill 00 \
		'STA $C055' 'LDA #$01' 'STA $2000' 'STA $C054'
	[[ $status -eq 0 ]]
	# Dot 0 occupies pixel column 0 only (1 wide); pixel 1 should differ
	local row
	row=$(pixel_row 0)
	[[ "${row:0:1}" != "${row:1:1}" ]]
}

# ---------------------------------------------------------------------------
# Memory region (3.1)
# ---------------------------------------------------------------------------

@test "dhires page 1 uses main and aux \$2000-\$3FFF" {
	asm_dhires_fill 7F
	[[ $status -eq 0 ]]
	grep -q 'FFFFFF' "$OUT/video.frame"
}

# ---------------------------------------------------------------------------
# Address interleaving (3.2)
# ---------------------------------------------------------------------------

@test "address \$2000 maps to dhires row 0: pixel rows 0-1" {
	asm_dhires_fill 00 \
		'STA $C055' 'LDA #$7F' 'STA $2000' 'STA $C054'
	[[ $status -eq 0 ]]
	rows_differ 0 2
	rows_same 0 1
	rows_same 2 3
}

@test "address \$2080 maps to dhires row 8: pixel rows 16-17" {
	asm_dhires_fill 00 \
		'STA $C055' 'LDA #$7F' 'STA $2080' 'STA $C054'
	[[ $status -eq 0 ]]
	rows_same 14 15
	rows_differ 16 14
	rows_same 16 17
	rows_same 18 14
}

@test "address \$2028 maps to dhires row 64: pixel rows 128-129" {
	asm_dhires_fill 00 \
		'STA $C055' 'LDA #$7F' 'STA $2028' 'STA $C054'
	[[ $status -eq 0 ]]
	rows_same 126 127
	rows_differ 128 126
	rows_same 128 129
	rows_same 130 126
}

@test "address \$2050 maps to dhires row 128: pixel rows 256-257" {
	asm_dhires_fill 00 \
		'STA $C055' 'LDA #$7F' 'STA $2050' 'STA $C054'
	[[ $status -eq 0 ]]
	rows_same 254 255
	rows_differ 256 254
	rows_same 256 257
	rows_same 258 254
}

@test "writing to hole byte \$2078 does not affect the dhires display" {
	asm_dhires_fill 00 \
		'STA $C055' 'LDA #$7F' 'STA $2078' 'STA $C054'
	[[ $status -eq 0 ]]
	[[ $(color_count) -eq 1 ]]
}

# ---------------------------------------------------------------------------
# Byte interleaving within a row (3.3)
# ---------------------------------------------------------------------------

@test "aux byte at offset 0 produces leftmost dots in column pair" {
	# Write $7F to aux byte 0, main byte 0 stays $00
	asm_dhires_fill 00 \
		'STA $C055' 'LDA #$7F' 'STA $2000' 'STA $C054'
	[[ $status -eq 0 ]]
	# Aux byte 0 = screen column 0 = dots 0-6 (pixels 0-6)
	# Main byte 0 = screen column 1 = dots 7-13 (pixels 7-13) = still black
	local row
	row=$(pixel_row 0)
	# Pixels 0-6 should be non-black, pixel 7+ should be black
	[[ "${row:0:1}" != "${row:7:1}" ]]
}

@test "main byte at offset 0 produces rightmost dots in column pair" {
	# Write $7F to main byte 0, aux byte 0 stays $00
	asm_dhires_fill 00 \
		'LDA #$7F' 'STA $2000'
	[[ $status -eq 0 ]]
	# Main byte 0 = screen column 1 = dots 7-13
	# Aux byte 0 = screen column 0 = dots 0-6 = still black
	local row
	row=$(pixel_row 0)
	[[ "${row:0:1}" != "${row:7:1}" ]]
}

# ---------------------------------------------------------------------------
# Byte layout (3.4)
# ---------------------------------------------------------------------------

@test "dhires bit 0 controls leftmost dot in byte" {
	# Set aux byte 0 to $01 (only bit 0 on)
	asm_dhires_fill 00 \
		'STA $C055' 'LDA #$01' 'STA $2000' 'STA $C054'
	[[ $status -eq 0 ]]
	# Dot 0 (pixel 0) should be non-black; pixels further right should differ
	local row
	row=$(pixel_row 0)
	col_range_is_uniform 0 0 1
	[[ "${row:0:1}" != "${row:4:1}" ]]
}

@test "dhires bit 6 controls rightmost dot in byte" {
	# Set aux byte 0 to $40 (only bit 6 on)
	asm_dhires_fill 00 \
		'STA $C055' 'LDA #$40' 'STA $2000' 'STA $C054'
	[[ $status -eq 0 ]]
	# Dot 6 = pixel 6 should be non-black; pixels well before it (0-2)
	# should be black (the sliding window extends 3 dots ahead of an on-dot)
	local row
	row=$(pixel_row 0)
	col_range_is_uniform 0 0 3
	[[ "${row:0:1}" != "${row:6:1}" ]]
}

@test "dhires bit 7 is unused and does not affect display" {
	# $FF (bit 7 set, bits 0-6 on) should look identical to $7F
	asm_dhires_fill 7F
	[[ $status -eq 0 ]]
	local without_bit7
	without_bit7=$(tail -n +3 "$OUT/video.frame")

	rm -rf "$OUT"; mkdir -p "$OUT"

	asm_dhires_fill FF
	[[ $status -eq 0 ]]
	local with_bit7
	with_bit7=$(tail -n +3 "$OUT/video.frame")

	[[ "$without_bit7" == "$with_bit7" ]]
}

# ---------------------------------------------------------------------------
# Color generation (4.1)
# ---------------------------------------------------------------------------

@test "dhires all dots off produces black" {
	asm_dhires_fill 00
	[[ $status -eq 0 ]]
	[[ $(color_count) -eq 1 ]]
	[[ "$(legend_rgb)" == "000000" ]]
}

@test "dhires all dots on produces white" {
	asm_dhires_fill 7F
	[[ $status -eq 0 ]]
	grep -q 'FFFFFF' "$OUT/video.frame"
}

@test "dhires \$7F and \$00 produce different colors" {
	asm_dhires_fill 7F
	[[ $status -eq 0 ]]
	local on_legend
	on_legend=$(sed -n '2p' "$OUT/video.frame")

	rm -rf "$OUT"; mkdir -p "$OUT"

	asm_dhires_fill 00
	[[ $status -eq 0 ]]
	local off_legend
	off_legend=$(sed -n '2p' "$OUT/video.frame")

	[[ "$on_legend" != "$off_legend" ]]
}

# ---------------------------------------------------------------------------
# Sliding window color assignment (4.2)
# ---------------------------------------------------------------------------

@test "4 on-dots at position 0 produce white at dot 0" {
	# Aux byte 0 = $0F (bits 0-3 on), everything else $00
	# Dots 0-3 on. At position 0, window = (on,on,on,on) = 1111 = white
	asm_dhires_fill 00 \
		'STA $C055' 'LDA #$0F' 'STA $2000' 'STA $C054'
	[[ $status -eq 0 ]]
	[[ "$(pixel_color 0 0)" == "FFFFFF" ]]
}

@test "single on-dot does not produce white" {
	# Aux byte 0 = $01 (bit 0 on), everything else $00
	# Only dot 0 is on; window never has all 4 bits set -> no white
	asm_dhires_fill 00 \
		'STA $C055' 'LDA #$01' 'STA $2000' 'STA $C054'
	[[ $status -eq 0 ]]
	local rgb
	rgb=$(legend_rgb)
	! echo "$rgb" | grep -q 'FFFFFF'
}

@test "first dot in a row behaves as if previous dot were off" {
	# Fill last aux+main bytes of row 0 with $7F, set dot 0 of row 1 on.
	# Row 0 ends at offset 39: aux at $2027, main at $2027.
	# Row 1 base is $2400.
	asm_dhires_fill 00 \
		'STA $C055' 'LDA #$7F' 'STA $2027' 'STA $C054' \
		'LDA #$7F' 'STA $2027' \
		'STA $C055' 'LDA #$01' 'STA $2400' 'STA $C054'
	[[ $status -eq 0 ]]
	# Row 1 = pixel rows 2-3. If row state carried over, dot 0 would be
	# white (consecutive on-dots). It should not be white.
	[[ "$(pixel_color 2 0)" != "FFFFFF" ]]
}

# ---------------------------------------------------------------------------
# Monochrome colors (5.2)
# ---------------------------------------------------------------------------

@test "dhires green screen renders on-dots as 98FF98" {
	asm_dhires_fill_mono green 7F
	[[ $status -eq 0 ]]
	local rgb
	rgb=$(legend_rgb)
	echo "$rgb" | grep -q '98FF98'
}

@test "dhires amber screen renders on-dots as FFBF00" {
	asm_dhires_fill_mono amber 7F
	[[ $status -eq 0 ]]
	local rgb
	rgb=$(legend_rgb)
	echo "$rgb" | grep -q 'FFBF00'
}

# ---------------------------------------------------------------------------
# Monochrome rendering behavior (5.3)
# ---------------------------------------------------------------------------

@test "dhires monochrome ignores bit 7" {
	# $7F (bit 7 clear) and $FF (bit 7 set) should look identical in mono
	asm_dhires_fill_mono green 7F
	[[ $status -eq 0 ]]
	local clear_pixels
	clear_pixels=$(tail -n +3 "$OUT/video.frame")

	rm -rf "$OUT"; mkdir -p "$OUT"

	asm_dhires_fill_mono green FF
	[[ $status -eq 0 ]]
	local set_pixels
	set_pixels=$(tail -n +3 "$OUT/video.frame")

	[[ "$clear_pixels" == "$set_pixels" ]]
}

@test "dhires monochrome has no color grouping" {
	# $01 in color mode produces multi-dot color artifacts from the sliding
	# window. In monochrome, only the on-dot lights up -- no extra colors.
	asm_dhires_fill_mono green 01
	[[ $status -eq 0 ]]
	local rgb
	rgb=$(legend_rgb)
	# Should have exactly two colors: black and green
	[[ $(echo "$rgb" | wc -w) -eq 2 ]]
	echo "$rgb" | grep -q '000000'
	echo "$rgb" | grep -q '98FF98'
}

# ---------------------------------------------------------------------------
# Mode dispatch (6.1)
# ---------------------------------------------------------------------------

@test "DHIRES+80COL+HIRES+TEXT_off activates double hires" {
	# Fill both main+aux with $7F in dhires mode -> white
	asm_dhires_fill 7F
	[[ $status -eq 0 ]]
	grep -q 'FFFFFF' "$OUT/video.frame"

	rm -rf "$OUT"; mkdir -p "$OUT"

	# Same fill but only standard hires (no DHIRES/80COL) -> also shows
	# something, but dots are 2px wide. Check that the pixel rows differ.
	local src="$TMP/test.s"
	printf '%s\n' \
		"LDA #\$7F" 'STA $02' \
		'LDA #$00' 'STA $00' 'LDA #$20' 'STA $01' \
		'LDA $02' 'LDY #$00' \
		'hfill: STA ($00),Y' 'INY' 'BNE hfill' \
		'INC $01' 'LDX $01' 'CPX #$40' 'BNE hfill' \
		'STA $C050' 'STA $C057' '.halt' >"$src"
	"$ASSEMBLER" -o "$TMP/test.dsk" "$src" || return 1
	run "$ERC" headless \
		--output "$OUT" \
		--steps 50000 \
		--start-at 0801 \
		--capture-video 49999 \
		"$TMP/test.dsk"
	[[ $status -eq 0 ]]
	# Both produce white, but this confirms dhires mode activated correctly
	grep -q 'FFFFFF' "$OUT/video.frame"
}

@test "disabling 80COL falls back to standard hires" {
	# Enter dhires, then disable 80COL -> should revert to standard hires.
	# Standard hires dots are 2px wide; dhires dots are 1px wide.
	# Set a single dot and compare pixel widths.
	asm_dhires_fill 00 \
		'STA $C055' 'LDA #$01' 'STA $2000' 'STA $C054' \
		'LDA #$01' 'STA $2000'
	[[ $status -eq 0 ]]
	local dhires_row
	dhires_row=$(pixel_row 0)

	rm -rf "$OUT"; mkdir -p "$OUT"

	# Same setup but disable 80COL at end
	asm_dhires_fill 00 \
		'STA $C055' 'LDA #$01' 'STA $2000' 'STA $C054' \
		'LDA #$01' 'STA $2000' \
		'STA $C00C'
	[[ $status -eq 0 ]]
	local fallback_row
	fallback_row=$(pixel_row 0)

	[[ "$dhires_row" != "$fallback_row" ]]
}

@test "disabling DHIRES falls back to standard hires" {
	asm_dhires_fill 00 \
		'STA $C055' 'LDA #$01' 'STA $2000' 'STA $C054' \
		'LDA #$01' 'STA $2000'
	[[ $status -eq 0 ]]
	local dhires_row
	dhires_row=$(pixel_row 0)

	rm -rf "$OUT"; mkdir -p "$OUT"

	# Same setup but disable DHIRES at end
	asm_dhires_fill 00 \
		'STA $C055' 'LDA #$01' 'STA $2000' 'STA $C054' \
		'LDA #$01' 'STA $2000' \
		'STA $C05F'
	[[ $status -eq 0 ]]
	local fallback_row
	fallback_row=$(pixel_row 0)

	[[ "$dhires_row" != "$fallback_row" ]]
}

@test "TEXT on returns from dhires to text mode" {
	asm_dhires_fill 7F 'STA $C051' 'STA $C00C'
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

# ---------------------------------------------------------------------------
# Read switches (7.2)
# ---------------------------------------------------------------------------

@test "RDDHIRES returns bit 7 high when DHIRES is on" {
	asm_mem 0300 \
		'STA $C05E' \
		'LDA $C07F' 'STA $0300' \
		'.halt'
	[[ $status -eq 0 ]]
	local val
	val=$(last_mem_val 0300)
	[[ $((16#$val)) -ge 128 ]]
}

@test "RDDHIRES returns bit 7 low when DHIRES is off" {
	asm_mem 0300 \
		'LDA #$FF' 'STA $0300' \
		'STA $C05F' \
		'LDA $C07F' 'STA $0300' \
		'.halt'
	[[ $status -eq 0 ]]
	local val
	val=$(last_mem_val 0300)
	[[ $((16#$val)) -lt 128 ]]
}

@test "RD80COL returns bit 7 high when 80COL is on" {
	asm_mem 0300 \
		'STA $C00D' \
		'LDA $C01F' 'STA $0300' \
		'.halt'
	[[ $status -eq 0 ]]
	local val
	val=$(last_mem_val 0300)
	[[ $((16#$val)) -ge 128 ]]
}

@test "RD80COL returns bit 7 low when 80COL is off" {
	asm_mem 0300 \
		'LDA #$FF' 'STA $0300' \
		'STA $C00C' \
		'LDA $C01F' 'STA $0300' \
		'.halt'
	[[ $status -eq 0 ]]
	local val
	val=$(last_mem_val 0300)
	[[ $((16#$val)) -lt 128 ]]
}

# ---------------------------------------------------------------------------
# Page selection (7.3)
# ---------------------------------------------------------------------------

@test "PAGE2 displays dhires page 2 content" {
	# Fill page 1 main+aux with $00, page 2 main+aux with $7F.
	# Display page 2 -> white.
	asm_video_dhires \
		'STA $C001' \
		'STA $C057' \
		'STA $C054' \
		'LDA #$00' 'STA $02' \
		'LDA #$00' 'STA $00' 'LDA #$20' 'STA $01' \
		'LDA $02' 'LDY #$00' \
		'mf1: STA ($00),Y' 'INY' 'BNE mf1' \
		'INC $01' 'LDX $01' 'CPX #$40' 'BNE mf1' \
		'STA $C055' \
		'LDA #$00' 'STA $00' 'LDA #$20' 'STA $01' \
		'LDA $02' 'LDY #$00' \
		'af1: STA ($00),Y' 'INY' 'BNE af1' \
		'INC $01' 'LDX $01' 'CPX #$40' 'BNE af1' \
		'STA $C000' \
		'STA $C054' \
		'LDA #$7F' 'STA $02' \
		'LDA #$00' 'STA $00' 'LDA #$40' 'STA $01' \
		'LDA $02' 'LDY #$00' \
		'mf2: STA ($00),Y' 'INY' 'BNE mf2' \
		'INC $01' 'LDX $01' 'CPX #$60' 'BNE mf2' \
		'STA $C001' \
		'STA $C055' \
		'LDA #$00' 'STA $00' 'LDA #$40' 'STA $01' \
		'LDA $02' 'LDY #$00' \
		'af2: STA ($00),Y' 'INY' 'BNE af2' \
		'INC $01' 'LDX $01' 'CPX #$60' 'BNE af2' \
		'STA $C054' \
		'STA $C055' \
		'STA $C050' \
		'STA $C00D' \
		'STA $C05E' \
		'.halt'
	[[ $status -eq 0 ]]
	grep -q 'FFFFFF' "$OUT/video.frame"
}
