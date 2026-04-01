ERC_BIN="$BATS_FILE_TMPDIR/erc"
SRC_DISK="$BATS_TEST_DIRNAME/../data/memreg.dsk"

setup_file() {
	if [[ ! -f "$SRC_DISK" ]]; then
		skip "disk image not found: $SRC_DISK"
	fi

	(cd "$BATS_TEST_DIRNAME/.." && go build -o "$ERC_BIN" .)
}

setup() {
	TMP="$BATS_TEST_TMPDIR"
	DISK="$TMP/test.dsk"
	cp "$SRC_DISK" "$DISK"
	export ERC_BIN TMP DISK
}

teardown() {
	rm -f "$TMP"/test.dsk*
}

# run_debug STEPS [EXTRA_FLAGS...] -- run headless with --debug-image on a
# copy of the disk image. Artifact files appear alongside $DISK.
run_debug() {
	local steps="$1"; shift
	run "$ERC_BIN" headless --steps "$steps" --debug-image "$@" "$DISK"
}
