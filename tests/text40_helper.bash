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
		--steps 5000 \
		--start-at 0801 \
		--capture-video 4999 \
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

# asm_video_mono MODE LINE [LINE...] -- like asm_video but with monochrome
# rendering.  MODE is "green" or "amber".
asm_video_mono() {
	local mode="$1"
	shift
	local src="$TMP/test.s"
	printf '%s\n' "$@" >"$src"
	"$ASSEMBLER" -o "$TMP/test.dsk" "$src" || return 1
	run "$ERC" headless \
		--output "$OUT" \
		--steps 5000 \
		--start-at 0801 \
		--capture-video 4999 \
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

# any_row_has_mixed_colors START END -- succeed if at least one pixel row in
# the range [START, END] contains more than one color.
any_row_has_mixed_colors() {
	local r
	for ((r = $1; r <= $2; r++)); do
		if ! row_is_uniform "$r"; then
			return 0
		fi
	done
	return 1
}

# asm_mem WATCH_ADDR LINE [LINE...] -- assemble and run headless, watching
# the given memory address for changes.  Sets bats $status and $output.
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

# col_range_has_mixed_colors ROW START LEN -- succeed if the given pixel range
# contains more than one color.
col_range_has_mixed_colors() {
	! col_range_is_uniform "$@"
}
