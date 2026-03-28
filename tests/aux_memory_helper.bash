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

_aux_run_impl() {
	local extra="$1"; shift
	local steps="${AUX_STEPS:-200}"
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
		--watch-comp MemReadAux,MemWriteAux,DisplayStore80,DisplayPage2,DisplayHires \
		$extra \
		"$TMP/test.dsk"
}

# aux_run LINE [LINE...] -- assemble and run headless, watching aux-related
# comp states. State log lands in $OUT/state.log.
aux_run() {
	_aux_run_impl "" "$@"
}

# aux_run_with_mem LINE [LINE...] -- like aux_run but also watches zero-page
# memory $00-$04 for storing readback values.
aux_run_with_mem() {
	_aux_run_impl "--watch-mem 00-04" "$@"
}

# _aux_run_80store_impl EXTRA LINE [LINE...] -- relocate user code to $0400
# with 80STORE enabled so instruction fetches stay in main memory even when
# ReadAux is on. A small relocator at $0801 copies the payload to $0400,
# enables 80STORE, and jumps there.
_aux_run_80store_impl() {
	local extra="$1"; shift
	local steps="${AUX_STEPS:-200}"
	local src="$TMP/test.s"
	{
		# Relocator: copy 128 bytes from $0812 to $0400, enable 80STORE,
		# then JMP $0400. The payload follows at code offset $0812.
		printf '%s\n' '    LDX #$80'
		printf '%s\n' 'reloc_loop:'
		printf '%s\n' '    LDA $0811,X'
		printf '%s\n' '    STA $03FF,X'
		printf '%s\n' '    DEX'
		printf '%s\n' '    BNE reloc_loop'
		printf '%s\n' '    STA $C001'
		printf '%s\n' '    JMP $0400'
		# Switch origin so labels in user code resolve to $0400-based
		# addresses. The bytes are still emitted contiguously after the
		# relocator, which is what LDA $0811,X reads from.
		printf '%s\n' '.org $0400'
		printf '%s\n' "$@"
	} >"$src"
	if ! "$ASSEMBLER" -o "$TMP/test.dsk" "$src" 2>&1; then
		status=1
		return 1
	fi

	run "$ERC_BIN" headless \
		--output "$OUT" \
		--start-at 0400 \
		--steps "$steps" \
		--watch-comp MemReadAux,MemWriteAux,DisplayStore80,DisplayPage2,DisplayHires \
		$extra \
		"$TMP/test.dsk"
}

# aux_run_80store LINE [LINE...] -- like aux_run but with code at $0400 and
# 80STORE on, so ReadAux doesn't affect instruction fetches.
aux_run_80store() {
	_aux_run_80store_impl "" "$@"
}

# aux_run_80store_with_mem LINE [LINE...] -- like aux_run_80store but also
# watches zero-page memory $00-$04.
aux_run_80store_with_mem() {
	_aux_run_80store_impl "--watch-mem 00-04" "$@"
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
