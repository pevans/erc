ASSEMBLER="$BATS_FILE_TMPDIR/erc-assembler"

setup_file() {
	(cd "$BATS_TEST_DIRNAME/.." && go build -o "$ASSEMBLER" ./cmd/erc-assembler)
}

setup() {
	TMP="$BATS_TEST_TMPDIR"
	export ASSEMBLER TMP
}

teardown() {
	:
}

# asm LINE [LINE...] -- write each argument as a source line to $TMP/test.s
# and assemble it to $TMP/out.dsk.  Sets bats $status and $output.
asm() {
	local src="$TMP/test.s"
	printf '%s\n' "$@" >"$src"
	run "$ASSEMBLER" -o "$TMP/out.dsk" "$src"
}

# byte_at FILE OFFSET -- print the hex byte at OFFSET in FILE (two lowercase
# hex digits, no spaces).
byte_at() {
	od -An -tx1 -j "$2" -N 1 "$1" | tr -d ' \n'
}
