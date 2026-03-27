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

# _bank_run_impl EXTRA_ARGS LINE [LINE...] -- shared implementation for
# bank_run and bank_run_with_mem. EXTRA_ARGS is a string of additional erc
# flags (may be empty).
_bank_run_impl() {
	local extra="$1"; shift
	local steps="${BANK_STEPS:-100}"
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
		--watch-comp BankReadRAM,BankWriteRAM,BankDFBlockBank2,BankSysBlockAux \
		$extra \
		"$TMP/test.dsk"
}

# bank_run LINE [LINE...] -- assemble and run headless, watching bank comp
# states. State log lands in $OUT/state.log.
bank_run() {
	_bank_run_impl "" "$@"
}

# bank_run_with_mem LINE [LINE...] -- like bank_run but also watches zero-page
# memory $00-$04 for storing readback values.
bank_run_with_mem() {
	_bank_run_impl "--watch-mem 00-04" "$@"
}

# _last_comp NAME -- return the last "new" value logged for a comp state.
_last_comp() {
	local name="$1"
	awk -v n="$name" \
		'$3=="comp" && $4==n {v=$NF} END{print v}' \
		"$OUT/state.log" 2>/dev/null
}

# _first_comp NAME -- return the first "new" value logged for a comp state.
_first_comp() {
	local name="$1"
	awk -v n="$name" \
		'$3=="comp" && $4==n {print $NF; exit}' \
		"$OUT/state.log" 2>/dev/null
}

# _last_mem ADDR4HEX -- return the last "new" value logged for the memory
# address.
_last_mem() {
	local addr="$1"
	awk -v a="\$$addr" \
		'$3=="mem" && $4==a {v=$NF} END{print v}' \
		"$OUT/state.log" 2>/dev/null
}
