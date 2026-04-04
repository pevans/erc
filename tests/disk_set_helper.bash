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

# make_disk NAME -- create a minimal disk image at $TMP/NAME.dsk.
# Each disk just halts immediately; the content doesn't matter,
# only that the files are distinct paths the disk set can load.
make_disk() {
	local name="$1"
	local src="$TMP/${name}.s"
	printf '%s\n' '.halt' >"$src"
	"$ASSEMBLER" -o "$TMP/${name}.dsk" "$src"
}

# diskset_run STEPS KEYS DISK [DISK...] -- run headless with the given disks,
# watching DiskIndex state changes.
diskset_run() {
	local steps="$1"; shift
	local keys="$1"; shift
	local args=(headless
		--output "$OUT"
		--start-at 0801
		--steps "$steps"
		--watch-comp DiskIndex)
	if [[ -n "$keys" ]]; then
		args+=(--keys "$keys")
	fi
	args+=("$@")
	run "$ERC_BIN" "${args[@]}"
}
