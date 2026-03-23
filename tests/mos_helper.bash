ERC_BIN="$BATS_FILE_TMPDIR/erc"
ASSEMBLER="$BATS_FILE_TMPDIR/erc-assembler"

# Flag masks matching the mos package bit positions.
CARRY=1
ZERO=2
INTERRUPT=4
DECIMAL=8
BREAK=16
UNUSED=32
OVERFLOW=64
NEGATIVE=128

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
	export CARRY ZERO INTERRUPT DECIMAL BREAK UNUSED OVERFLOW NEGATIVE
}

teardown() {
	rm -rf "$OUT"
}

# cpu_run LINE [LINE...] -- assemble the given source lines and run them
# headless from $0801, watching zero-page $00-$04 and registers A,X,Y,S,P.
# After the call, $status reflects the headless exit code and state.log is
# available at $OUT/state.log.
#
# Memory layout convention:
#   $00 -- primary result (store A, X, Y, or tested memory value here)
#   $01 -- secondary result
#   $02 -- P register capture (via PHP / PLA / STA $02 before .halt)
#   $03 -- extra
#   $04 -- extra
cpu_run() {
	local steps="${CPU_STEPS:-50}"
	local src="$TMP/test.s"
	printf '%s\n' "$@" >"$src"
	if ! "$ASSEMBLER" -o "$TMP/out.dsk" "$src" 2>&1; then
		status=1
		return 1
	fi
	run "$ERC_BIN" headless \
		--output "$OUT" \
		--start-at 0801 \
		--steps "$steps" \
		--watch-mem 00-04 \
		--watch-reg A,X,Y,S,P \
		"$TMP/out.dsk"
}

# _last_mem ADDR4HEX -- return the last "new" value logged for the memory
# address (e.g. "0000", "0002").  Returns empty string if no entry exists.
_last_mem() {
	local addr="$1"
	awk -v a="\$$addr" \
		'$3=="mem" && $4==a {v=$NF} END{print v}' \
		"$OUT/state.log" 2>/dev/null
}

# _last_reg NAME -- return the last "new" value logged for a register.
# Returns empty string if no entry exists.
_last_reg() {
	local name="$1"
	awk -v n="$name" \
		'$3=="reg" && $4==n {v=$NF} END{print v}' \
		"$OUT/state.log" 2>/dev/null
}

# _assert_hex_eq LABEL EXPECTED_HEX ACTUAL -- compare ACTUAL (in "$XX" format
# or empty) to EXPECTED_HEX (e.g. "42").  Prints a message and returns 1 on
# mismatch.
_assert_hex_eq() {
	local label="$1" expected_raw="$2" actual="$3"
	local expected actual_hex
	expected="$(printf '%02X' "0x${expected_raw}")"
	if [[ -z "$actual" ]]; then
		actual_hex="00"
	else
		actual_hex="$(printf '%02X' "0x${actual#\$}")"
	fi
	if [[ "$actual_hex" != "$expected" ]]; then
		echo "$label: expected \$$expected, got \$$actual_hex" >&2
		return 1
	fi
}

# assert_zp DECIMAL_OFFSET EXPECTED_HEX -- assert that the zero-page byte at
# OFFSET (0-4) equals EXPECTED_HEX (e.g. "42", "FF").
assert_zp() {
	local offset="$1"
	local addr
	addr="$(printf '%04X' "$offset")"
	_assert_hex_eq "assert_zp $offset" "$2" "$(_last_mem "$addr")"
}

# assert_reg NAME EXPECTED_HEX -- assert that a register's last observed value
# equals EXPECTED_HEX (e.g. assert_reg A 42).
assert_reg() {
	_assert_hex_eq "assert_reg $1" "$2" "$(_last_reg "$1")"
}

# _get_p -- read P from memory $02, where tests capture it via PHP/PLA/STA $02.
# Returns the value as a decimal integer.
_get_p() {
	local val
	val="$(_last_mem "0002")"
	if [[ -z "$val" ]]; then
		printf '0'
		return
	fi
	printf '%d' "0x${val#\$}"
}

# assert_flag DECIMAL_MASK -- assert that all bits in MASK are set in P.
assert_flag() {
	local mask="$1"
	local p result
	p="$(_get_p)"
	result=$(( p & mask ))
	if [[ "$result" -ne "$mask" ]]; then
		echo "assert_flag $mask: P=$(printf '%02X' "$p") does not have mask $(printf '%02X' "$mask") set" >&2
		return 1
	fi
}

# refute_flag DECIMAL_MASK -- assert that all bits in MASK are clear in P.
refute_flag() {
	local mask="$1"
	local p result
	p="$(_get_p)"
	result=$(( p & mask ))
	if [[ "$result" -ne 0 ]]; then
		echo "refute_flag $mask: P=$(printf '%02X' "$p") has mask $(printf '%02X' "$mask") set (should be clear)" >&2
		return 1
	fi
}
