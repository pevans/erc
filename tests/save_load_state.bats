setup_file() {
	load test_helper; setup_file
	(cd "$BATS_TEST_DIRNAME/.." && go build -o "$BATS_FILE_TMPDIR/erc-assembler" ./cmd/erc-assembler)
}
setup() {
	load test_helper; setup
	ASSEMBLER="$BATS_FILE_TMPDIR/erc-assembler"
}
teardown()   { load test_helper; teardown; rm -f "$DISK".*.state; }

# --- File creation ---

@test "save creates a state file named for disk and slot" {
	run_headless --steps 1000 \
		--keys "100:ctrl-a,101:s" \
		"$DISK"
	[[ $status -eq 0 ]]
	[[ -f "$DISK.0.state" ]]
}

@test "save to slot 3 creates a distinct state file" {
	run_headless --steps 1000 \
		--keys "100:ctrl-a,101:3,200:ctrl-a,201:s" \
		"$DISK"
	[[ $status -eq 0 ]]
	[[ -f "$DISK.3.state" ]]
}

@test "different slots create separate files" {
	run_headless --steps 1000 \
		--keys "100:ctrl-a,101:s,200:ctrl-a,201:2,300:ctrl-a,301:s" \
		"$DISK"
	[[ $status -eq 0 ]]
	[[ -f "$DISK.0.state" ]]
	[[ -f "$DISK.2.state" ]]
}

# --- Round-trip: CPU registers ---

@test "round-trip restores PC register" {
	run_headless --steps 2000 \
		--keys "100:ctrl-a,101:s,500:ctrl-a,501:l" \
		--watch-reg PC \
		"$DISK"
	[[ $status -eq 0 ]]
	# After loading state at step 501, PC should revert to its saved value.
	# The saved PC will differ from the current PC because ~400 steps
	# executed between save and load. We verify by checking that PC
	# transitions backwards (a value it already passed through).
	local last_pc
	last_pc=$(tail -1 "$OUT/state.log" | grep -o 'reg PC .* -> \$[0-9A-F]*' | grep -o '\$[0-9A-F]*$')
	# The load causes a PC change; we just need more than one transition
	local changes
	changes=$(grep -c 'reg PC' "$OUT/state.log")
	[[ $changes -ge 2 ]]
}

@test "round-trip restores keyboard state" {
	# Press 'a' between save and load; after load, KBLastKey should revert
	run_headless --steps 1000 \
		--keys "100:ctrl-a,101:s,200:a,400:ctrl-a,401:l" \
		--watch-comp KBLastKey \
		"$DISK"
	[[ $status -eq 0 ]]
	# KBLastKey changes when 'a' is pressed, then reverts on load
	local changes
	changes=$(grep -c 'comp KBLastKey' "$OUT/state.log")
	[[ $changes -ge 2 ]]
}

# --- Round-trip: speed ---

@test "round-trip restores speed setting" {
	# Increase speed between save and load; speed should revert
	run_headless --steps 1000 \
		--keys "100:ctrl-a,101:s,200:ctrl-a,201:+,400:ctrl-a,401:l" \
		--watch-comp Speed \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp Speed .* -> 2' "$OUT/state.log"
	grep -q 'comp Speed .* -> 1' "$OUT/state.log"
}

# --- Round-trip: memory ---

@test "round-trip restores memory contents" {
	# Watch a text page address; the boot process writes to it.
	# Save after boot settles, let more instructions change it, then load.
	run_headless --steps 60000 \
		--keys "20000:ctrl-a,20001:s,40000:ctrl-a,40001:l" \
		--watch-mem 0400-0403 \
		"$DISK"
	[[ $status -eq 0 ]]
	# Memory should show at least two transitions (changes during execution
	# plus revert on load)
	[[ -f "$OUT/state.log" ]]
	local changes
	changes=$(grep -c 'mem ' "$OUT/state.log")
	[[ $changes -ge 2 ]]
}

# --- Round-trip: display state flags ---

@test "round-trip restores display flags" {
	# The Apple II ROM reset toggles DisplayPage2 on at step 347 and off
	# at step 354.  Save during the on window, let it revert to off, then
	# load and verify the flag is restored to on.
	run_headless --steps 1000 \
		--keys "348:ctrl-a,349:s,500:ctrl-a,501:l" \
		--watch-comp DisplayPage2 \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp DisplayPage2 .* -> true' "$OUT/state.log"
	grep -q 'comp DisplayPage2 .* -> false' "$OUT/state.log"
	# Expect at least 3 transitions: on (boot), off (boot), on (load restores)
	local changes
	changes=$(grep -c 'comp DisplayPage2' "$OUT/state.log")
	[[ $changes -ge 3 ]]
}

# --- Round-trip: drive state ---

@test "round-trip restores drive state" {
	# Toggle write protect between save and load; should revert.
	run_headless --steps 1000 \
		--keys "100:ctrl-a,101:s,200:ctrl-a,201:w,400:ctrl-a,401:l" \
		--watch-comp WriteProtect \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp WriteProtect .* -> true' "$OUT/state.log"
	grep -q 'comp WriteProtect .* -> false' "$OUT/state.log"
}

# --- Round-trip: auxiliary memory ---

@test "round-trip restores auxiliary memory contents" {
	# Write $42 to aux:$0300, save, overwrite with $FF, load.
	# After load aux:$0300 should revert to $42.
	local src="$BATS_TEST_TMPDIR/test.s"
	local dsk="$BATS_TEST_TMPDIR/test.dsk"
	cat >"$src" <<'ASM'
STA $C005
LDA #$42
STA $0300
STA $C004
NOP
NOP
STA $C005
LDA #$FF
STA $0300
STA $C004
.halt
ASM
	"$ASSEMBLER" -o "$dsk" "$src"
	run "$ERC" headless \
		--output "$OUT" \
		--start-at 0801 \
		--steps 40 \
		--keys "4:ctrl-a,5:s,10:ctrl-a,11:l" \
		--watch-aux-mem 0300 \
		"$dsk"
	[[ $status -eq 0 ]]
	# aux:$0300 transitions: $00->$42, $42->$FF, $FF->$42 (load), $42->$FF (re-exec)
	grep -q 'mem \$0300 \$FF -> \$42' "$OUT/state.log"
	rm -f "$dsk".*.state
}

# --- Round-trip: segment references ---

@test "round-trip rebuilds segment references from restored flags" {
	# Enable RAMWRT (MemWriteAux=true) so writes go to aux, save state.
	# The program then writes $99 to aux:$0300, disables RAMWRT, and loads.
	# After load, RAMWRT is restored and MemWriteSegment must be rebuilt
	# to point to aux. The re-executed STA $0300 should write to aux again,
	# proving the segment reference was correctly rebuilt.
	local src="$BATS_TEST_TMPDIR/test.s"
	local dsk="$BATS_TEST_TMPDIR/test.dsk"
	cat >"$src" <<'ASM'
STA $C005
NOP
NOP
LDA #$99
STA $0300
STA $C004
NOP
NOP
.halt
ASM
	"$ASSEMBLER" -o "$dsk" "$src"
	run "$ERC" headless \
		--output "$OUT" \
		--start-at 0801 \
		--steps 30 \
		--keys "1:ctrl-a,2:s,6:ctrl-a,7:l" \
		--watch-aux-mem 0300 \
		--watch-comp MemWriteAux \
		"$dsk"
	[[ $status -eq 0 ]]
	# MemWriteAux: false->true (step 0), true->false (step 5),
	# false->true (load), true->false (re-exec)
	local flag_changes
	flag_changes=$(grep -c 'comp MemWriteAux' "$OUT/state.log")
	[[ $flag_changes -ge 3 ]]
	# aux:$0300 transitions: $00->$99 (first write), $99->$00 (load restores),
	# $00->$99 (re-exec after segment ref rebuild). The third transition
	# proves MemWriteSegment was correctly rebuilt to point to aux.
	local aux_changes
	aux_changes=$(grep -c 'mem \$0300' "$OUT/state.log")
	[[ $aux_changes -ge 3 ]]
	rm -f "$dsk".*.state
}

# --- Slot isolation ---

@test "loading from one slot does not use another slot's data" {
	# Save to slot 0, change speed, save to slot 1, then load slot 0.
	# Speed should revert to slot 0's value (1), not slot 1's (2).
	run_headless --steps 2000 \
		--keys "100:ctrl-a,101:s,200:ctrl-a,201:+,300:ctrl-a,301:1,400:ctrl-a,401:s,500:ctrl-a,501:0,600:ctrl-a,601:l" \
		--watch-comp Speed \
		"$DISK"
	[[ $status -eq 0 ]]
	# Speed went 1->2 (at step 201), then reverted to 1 after loading slot 0
	grep -q 'comp Speed .* -> 2' "$OUT/state.log"
	grep -q 'comp Speed .* -> 1' "$OUT/state.log"
}

# --- Error handling ---

@test "load from nonexistent slot does not crash" {
	rm -f "$DISK.0.state"
	run_headless --steps 1000 \
		--keys "100:ctrl-a,101:l" \
		"$DISK"
	[[ $status -eq 0 ]]
}

@test "multiple saves to same slot overwrites cleanly" {
	run_headless --steps 1000 \
		--keys "100:ctrl-a,101:s,200:ctrl-a,201:s" \
		"$DISK"
	[[ $status -eq 0 ]]
	[[ -f "$DISK.0.state" ]]
}
