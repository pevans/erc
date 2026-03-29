setup_file() {
	ERC="$BATS_FILE_TMPDIR/erc"
	ASSEMBLER="$BATS_FILE_TMPDIR/erc-assembler"
	(cd "$BATS_TEST_DIRNAME/.." && go build -o "$ERC" . && go build -o "$ASSEMBLER" ./cmd/erc-assembler)
}

setup() {
	ERC="$BATS_FILE_TMPDIR/erc"
	ASSEMBLER="$BATS_FILE_TMPDIR/erc-assembler"
	TMP="$BATS_TEST_TMPDIR"
	OUT="$TMP/out"
	FONTDIR="$BATS_TEST_DIRNAME/../specs/fontdata/apple2"
	mkdir -p "$OUT"
}

teardown() {
	rm -rf "$OUT"
}

# asm_video LINE [LINE...] -- assemble lines into a disk image, run headless
# with video capture at the last step.
asm_video() {
	local src="$TMP/test.s"
	printf '%s\n' "$@" >"$src"
	"$ASSEMBLER" -o "$TMP/test.dsk" "$src" || return 1
	run "$ERC" headless \
		--output "$OUT" \
		--steps 5000 \
		--start-at 0801 \
		--capture-video 4999 \
		"$TMP/test.dsk"
}

# pixel_row N -- print pixel row N (0-indexed) from the captured frame.
pixel_row() {
	sed -n "$((3 + $1))p" "$OUT/video.frame"
}

# fill_screen_40 HEX -- fill the 40-col screen with byte $HEX and capture.
fill_screen_40() {
	asm_video \
		"LDA #\$$1" 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' '.halt'
}

# fill_screen_40_alt HEX -- enable the alternate character set, then fill the
# 40-col screen with byte $HEX and capture.  Mousetext ($40-$5F) only appears
# when the alternate character set is active.
fill_screen_40_alt() {
	asm_video \
		'STA $C00F' \
		"LDA #\$$1" 'LDX #$00' \
		'fill: STA $0400,X' 'STA $0500,X' 'STA $0600,X' 'STA $0700,X' \
		'INX' 'BNE fill' '.halt'
}

# glyph_pattern_40 -- extract the 7x8 base glyph from cell (0,0) of a 40-col
# filled screen.  Outputs 8 lines of 7 chars: '#' for on, '.' for off.
# Background color is identified as the more frequent of the two colors in the
# first pixel row. Pixel (0,0) is not reliable because some glyphs (e.g.
# mousetext 0x5F) have a foreground pixel there.
glyph_pattern_40() {
	local first_row
	first_row=$(pixel_row 0)
	local ch1="${first_row:0:1}" count1=0 ch2="" count2=0
	local i ch
	for ((i = 0; i < ${#first_row}; i += 2)); do
		ch="${first_row:$i:1}"
		if [[ "$ch" == "$ch1" ]]; then
			((count1++))
		else
			ch2="$ch"
			((count2++))
		fi
	done
	local bg
	if ((count1 >= count2)); then
		bg="$ch1"
	else
		bg="$ch2"
	fi
	local r c
	for ((r = 0; r < 8; r++)); do
		local pr=$((r * 2))
		local line
		line=$(pixel_row "$pr")
		local glyph_row=""
		for ((c = 0; c < 7; c++)); do
			local pc=$((c * 2))
			local ch="${line:$pc:1}"
			if [[ "$ch" == "$bg" ]]; then
				glyph_row+="."
			else
				glyph_row+="#"
			fi
		done
		printf '%s\n' "$glyph_row"
	done
}

# bg_color -- print the hex RGB of the top-left pixel (the first palette entry).
bg_color() {
	sed -n '2s/^colors: .=\([0-9A-Fa-f]*\).*/\1/p' "$OUT/video.frame"
}

# reset_out -- clear and recreate $OUT between captures within a single test.
reset_out() {
	rm -rf "$OUT"
	mkdir -p "$OUT"
}

# expected_glyph FILE NAME -- print the 8-line glyph pattern from a font data
# file for the glyph named NAME.
expected_glyph() {
	local file="$1" name="$2"
	local line_num
	line_num=$(grep -nF -- "--- $name ---" "$file" | head -1 | cut -d: -f1)
	sed -n "$((line_num + 1)),$((line_num + 8))p" "$file"
}
