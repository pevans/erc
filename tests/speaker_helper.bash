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

# spk_run LINE [LINE...] -- assemble source lines, boot headless from $0801,
# watch SpeakerState. Set SPK_STEPS to override the default step count (200).
spk_run() {
	local steps="${SPK_STEPS:-200}"
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
		--watch-comp SpeakerState)
	args+=("$TMP/test.dsk")
	run "$ERC_BIN" "${args[@]}"
}

# spk_run_audio LINE [LINE...] -- like spk_run but also records audio.
spk_run_audio() {
	local steps="${SPK_STEPS:-100000}"
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
		--record-audio
		--watch-comp SpeakerState)
	args+=("$TMP/test.dsk")
	run "$ERC_BIN" "${args[@]}"
}

# _last_comp NAME -- return the last "new" value logged for a comp state.
_last_comp() {
	local name="$1"
	awk -v n="$name" \
		'$3=="comp" && $4==n {v=$NF} END{print v}' \
		"$OUT/state.log" 2>/dev/null
}
