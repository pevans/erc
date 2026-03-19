DISK="$BATS_TEST_DIRNAME/../data/memreg.dsk"

setup_file() {
	if [[ ! -f "$DISK" ]]; then
		skip "disk image not found: $DISK"
	fi

	ERC="$BATS_FILE_TMPDIR/erc"
	export ERC
	(cd "$BATS_TEST_DIRNAME/.." && go build -o "$ERC" .)
}

setup() {
	ERC="$BATS_FILE_TMPDIR/erc"
	OUT="$BATS_TEST_TMPDIR/out"
	mkdir -p "$OUT"
}

teardown() {
	rm -rf "$OUT"
}

run_headless() {
	run "$ERC" headless --output "$OUT" "$@"
}
