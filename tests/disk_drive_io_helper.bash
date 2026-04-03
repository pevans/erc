ERC_BIN="$BATS_FILE_TMPDIR/erc"
ASSEMBLER="$BATS_FILE_TMPDIR/erc-assembler"

setup_file() {
	(cd "$BATS_TEST_DIRNAME/.." && go build -o "$ERC_BIN" .) &
	(cd "$BATS_TEST_DIRNAME/.." && go build -o "$ASSEMBLER" ./cmd/erc-assembler) &
	wait
}

setup() {
	TMP="$BATS_TEST_TMPDIR"
	OUT="$BATS_TEST_TMPDIR/out"
	mkdir -p "$OUT"
	export ERC_BIN ASSEMBLER TMP OUT
}

teardown() {
	rm -rf "$OUT"
}

# asm LINE [LINE...] -- assemble source lines to $TMP/test.dsk.
asm() {
	local src="$TMP/test.s"
	printf '%s\n' "$@" >"$src"
	"$ASSEMBLER" -o "$TMP/test.dsk" "$src"
}

# disk_run LINE [LINE...] -- assemble source lines, boot headless from $0801,
# watch zero-page $00-$04 and comp state (WriteProtect, FullSpeed).
# Set DISK_STEPS to override the default step count (100).
# Set DISK_KEYS to inject key events (e.g. "0:ctrl-a,1:w").
disk_run() {
	local steps="${DISK_STEPS:-100}"
	local src="$TMP/test.s"
	printf '%s\n' "$@" >"$src"
	if ! "$ASSEMBLER" -o "$TMP/test.dsk" "$src" 2>&1; then
		status=1
		return 1
	fi
	local args=(headless
		--output "$OUT"
		--start-at 0801
		--steps "$steps"
		--watch-mem 00-04
		--watch-comp WriteProtect,FullSpeed)
	if [[ -n "${DISK_KEYS:-}" ]]; then
		args+=(--keys "$DISK_KEYS")
	fi
	args+=("$TMP/test.dsk")
	run "$ERC_BIN" "${args[@]}"
}

# _last_mem ADDR4HEX -- return the last "new" value logged for the memory
# address.
_last_mem() {
	local addr="$1"
	awk -v a="\$$addr" \
		'$3=="mem" && $4==a {v=$NF} END{print v}' \
		"$OUT/state.log" 2>/dev/null
}

# _last_comp NAME -- return the last "new" value logged for a comp state.
_last_comp() {
	local name="$1"
	awk -v n="$name" \
		'$3=="comp" && $4==n {v=$NF} END{print v}' \
		"$OUT/state.log" 2>/dev/null
}

# decode_4and4 HIGH LOW -- decode a 4-and-4 encoded byte pair to its original
# value. HIGH and LOW are in $XX hex format (as returned by _last_mem).
decode_4and4() {
	local h l
	h=$((16#${1#\$}))
	l=$((16#${2#\$}))
	echo $(( ((h & 0x55) << 1) | (l & 0x55) ))
}
