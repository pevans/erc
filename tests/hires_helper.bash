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
	mkdir -p "$OUT"
}

teardown() {
	rm -rf "$OUT"
}

# asm_video LINE [LINE...] -- assemble lines into a disk image, run headless
# with video capture at the last step.  Sets bats $status and $output.
asm_video() {
	local src="$TMP/test.s"
	printf '%s\n' "$@" >"$src"
	"$ASSEMBLER" -o "$TMP/test.dsk" "$src" || return 1
	run "$ERC" headless \
		--output "$OUT" \
		--steps 50000 \
		--start-at 0801 \
		--capture-video 49999 \
		"$TMP/test.dsk"
}

# asm_video_long LINE [LINE...] -- like asm_video but with more steps,
# for tests that need two full hires page fills.
asm_video_long() {
	local src="$TMP/test.s"
	printf '%s\n' "$@" >"$src"
	"$ASSEMBLER" -o "$TMP/test.dsk" "$src" || return 1
	run "$ERC" headless \
		--output "$OUT" \
		--steps 80000 \
		--start-at 0801 \
		--capture-video 79999 \
		"$TMP/test.dsk"
}

# asm_state WATCH_COMP LINE [LINE...] -- assemble and run headless, observing
# the named computer state.  Sets bats $status and $output.
asm_state() {
	local watch="$1"
	shift
	local src="$TMP/test.s"
	printf '%s\n' "$@" >"$src"
	"$ASSEMBLER" -o "$TMP/test.dsk" "$src" || return 1
	run "$ERC" headless \
		--output "$OUT" \
		--steps 1000 \
		--start-at 0801 \
		--watch-comp "$watch" \
		"$TMP/test.dsk"
}

# asm_mem WATCH_ADDR LINE [LINE...] -- assemble and run headless, observing
# the named memory address.  Sets bats $status and $output.
asm_mem() {
	local watch="$1"
	shift
	local src="$TMP/test.s"
	printf '%s\n' "$@" >"$src"
	"$ASSEMBLER" -o "$TMP/test.dsk" "$src" || return 1
	run "$ERC" headless \
		--output "$OUT" \
		--steps 1000 \
		--start-at 0801 \
		--watch-mem "$watch" \
		"$TMP/test.dsk"
}

# asm_hires_fill VAL [EXTRA_LINES...] -- fill hires page 1 ($2000-$3FFF) with
# $VAL, enable hires mode, execute any extra lines, then halt.
# Uses zero-page $00-$01 as pointer, $02 as temp.
asm_hires_fill() {
	local val="$1"
	shift
	asm_video \
		"LDA #\$$val" 'STA $02' \
		'LDA #$00' 'STA $00' 'LDA #$20' 'STA $01' \
		'LDA $02' 'LDY #$00' \
		'hfill: STA ($00),Y' 'INY' 'BNE hfill' \
		'INC $01' 'LDX $01' 'CPX #$40' 'BNE hfill' \
		'STA $C050' 'STA $C057' \
		"$@" '.halt'
}

# asm_hires_fill_mono MODE VAL [EXTRA_LINES...] -- like asm_hires_fill but with
# monochrome rendering.  MODE is "green" or "amber".
asm_hires_fill_mono() {
	local mode="$1" val="$2"
	shift 2
	local src="$TMP/test.s"
	printf '%s\n' \
		"LDA #\$$val" 'STA $02' \
		'LDA #$00' 'STA $00' 'LDA #$20' 'STA $01' \
		'LDA $02' 'LDY #$00' \
		'hfill: STA ($00),Y' 'INY' 'BNE hfill' \
		'INC $01' 'LDX $01' 'CPX #$40' 'BNE hfill' \
		'STA $C050' 'STA $C057' \
		"$@" '.halt' >"$src"
	"$ASSEMBLER" -o "$TMP/test.dsk" "$src" || return 1
	run "$ERC" headless \
		--output "$OUT" \
		--steps 50000 \
		--start-at 0801 \
		--capture-video 49999 \
		--monochrome "$mode" \
		"$TMP/test.dsk"
}

# color_count -- print the number of distinct colors in the captured frame.
color_count() {
	sed -n '2p' "$OUT/video.frame" | grep -o '=' | wc -l | tr -d ' '
}

# pixel_row N -- print pixel row N (0-indexed) from the captured frame.
pixel_row() {
	sed -n "$((3 + $1))p" "$OUT/video.frame"
}

# row_is_uniform N -- succeed if every pixel in row N is the same color.
row_is_uniform() {
	local row first
	row=$(pixel_row "$1")
	first="${row:0:1}"
	[[ -z "$(printf '%s' "$row" | tr -d "$first")" ]]
}

# rows_same A B -- succeed if pixel rows A and B have identical content.
rows_same() {
	[[ "$(pixel_row "$1")" == "$(pixel_row "$2")" ]]
}

# rows_differ A B -- succeed if pixel rows A and B have different content.
rows_differ() {
	! rows_same "$@"
}

# legend_rgb -- print the RGB hex string(s) from the color legend.
legend_rgb() {
	sed -n '2p' "$OUT/video.frame" | grep -oE '[0-9A-F]{6}'
}

# col_range_is_uniform ROW START LEN -- succeed if all pixels in pixel row ROW
# from column START to START+LEN-1 are the same color.
col_range_is_uniform() {
	local row segment first
	row=$(pixel_row "$1")
	segment="${row:$2:$3}"
	first="${segment:0:1}"
	[[ -z "$(printf '%s' "$segment" | tr -d "$first")" ]]
}

# pixel_color ROW COL -- print the RGB hex of the pixel at (ROW, COL).
pixel_color() {
	local row legend char
	row=$(pixel_row "$1")
	legend=$(sed -n '2p' "$OUT/video.frame")
	char="${row:$2:1}"
	echo "$legend" | tr ' ' '\n' | grep "^${char}=" | cut -d= -f2
}

# last_mem_val ADDR -- print the last recorded value of memory address $ADDR
# from the state log.
last_mem_val() {
	local pattern="mem \$$1"
	awk -v p="$pattern" 'index($0, p) {v=$NF} END{print v}' "$OUT/state.log" | tr -d '$'
}
