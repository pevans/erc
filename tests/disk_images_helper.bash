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

# byte_at FILE OFFSET -- print the hex byte at OFFSET (two lowercase hex
# digits, no spaces).
byte_at() {
	od -An -tx1 -j "$2" -N 1 "$1" | tr -d ' \n'
}

# encode INPUT OUTPUT -- run erc encode, setting bats $status and $output.
encode() {
	run "$ERC_BIN" encode "$1" -o "$2"
}

# decode INPUT OUTPUT -- run erc decode, setting bats $status and $output.
decode() {
	run "$ERC_BIN" decode "$1" -o "$2"
}

# make_zeros FILE SIZE -- create a file of SIZE zero bytes.
make_zeros() {
	dd if=/dev/zero of="$1" bs="$2" count=1 2>/dev/null
}

# make_patterned FILE -- create a 143360-byte .dsk where each sector has
# distinct content (sector N of track T filled with byte ((T*16+N+1) & 0xFF)).
make_patterned() {
	local file="$1"
	: >"$file"
	for t in $(seq 0 34); do
		for s in $(seq 0 15); do
			local val=$(( ((t * 16 + s + 1) & 0xFF) ))
			printf "$(printf '\\x%02x' "$val")%.0s" {1..256} >>"$file"
		done
	done
}

# disk_run LINE [LINE...] -- assemble source lines, boot headless from $0801,
# watch zero-page $00-$04 and WriteProtect comp state.  State log lands in
# $OUT/state.log.
disk_run() {
	local steps="${DISK_STEPS:-100}"
	local src="$TMP/test.s"
	printf '%s\n' "$@" >"$src"
	if ! "$ASSEMBLER" -o "$TMP/test.dsk" "$src" 2>&1; then
		status=1
		return 1
	fi
	run "$ERC_BIN" headless \
		--output "$OUT" \
		--start-at 0801 \
		--steps "$steps" \
		--watch-mem 00-04 \
		--watch-comp WriteProtect \
		"$TMP/test.dsk"
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
