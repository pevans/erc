setup_file() { load peripheral_helper; setup_file; }
setup()      { load peripheral_helper; setup; }
teardown()   { load peripheral_helper; teardown; }

# Run an assembly program watching keyboard comp state and zero-page memory.
# Set KB_KEYS to inject key events (e.g. KB_KEYS="0:a").
kb_run() {
	local src="$TMP/test.s"
	local steps="${KB_STEPS:-100}"
	local -a keys_args=()
	if [[ -n "${KB_KEYS:-}" ]]; then
		keys_args=(--keys "$KB_KEYS")
	fi
	printf '%s\n' "$@" >"$src"
	if ! "$ASSEMBLER" -o "$TMP/test.dsk" "$src" 2>&1; then
		status=1
		return 1
	fi
	run "$ERC_BIN" headless \
		--output "$OUT" \
		--start-at 0801 \
		--steps "$steps" \
		--watch-comp KBStrobe,KBKeyDown \
		--watch-mem 00-04 \
		"${keys_args[@]}" \
		"$TMP/test.dsk"
}

# --- Section 3.2: $C010 -- Any Key Down / Strobe Clear ---

@test "reading C010 clears strobe" {
	KB_KEYS="0:a" kb_run \
		'NOP' \
		'LDA $C010' \
		'.halt'
	[[ $status -eq 0 ]]
	# NOP executes at step 0 (key pressed, strobe set to $80, recorded)
	# LDA $C010 executes at step 1 (SwitchRead clears strobe to $00)
	[[ "$(_last_comp KBStrobe)" == '$00' ]]
}

@test "writing C010 clears strobe" {
	KB_KEYS="0:a" kb_run \
		'NOP' \
		'STA $C010' \
		'.halt'
	[[ $status -eq 0 ]]
	[[ "$(_last_comp KBStrobe)" == '$00' ]]
}

@test "reading C010 returns KeyDown" {
	KB_KEYS="0:a" kb_run \
		'LDA $C010' \
		'STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	# Key is down at step 0 when LDA $C010 executes; KeyDown is $80
	[[ "$(_last_mem 0000)" == '$80' ]]
}
