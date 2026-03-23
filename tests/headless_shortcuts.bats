setup_file() { load test_helper; setup_file; }
setup()      { load test_helper; setup; }
teardown()   { load test_helper; teardown; }

# --- Pause and Resume ---

@test "ctrl-a esc pauses the computer" {
	run_headless --steps 1000 \
		--keys "100:ctrl-a,101:esc" \
		--watch-comp Paused \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp Paused .* -> true' "$OUT/state.log"
}

@test "esc while paused resumes" {
	run_headless --steps 1000 \
		--keys "100:ctrl-a,101:esc,200:esc" \
		--watch-comp Paused \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp Paused .* -> true' "$OUT/state.log"
	grep -q 'comp Paused .* -> false' "$OUT/state.log"
}

@test "non-esc key while paused stays paused" {
	run_headless --steps 1000 \
		--keys "100:ctrl-a,101:esc,200:a" \
		--watch-comp Paused \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp Paused .* -> true' "$OUT/state.log"
	! grep -q 'comp Paused .* -> false' "$OUT/state.log"
}

# --- Speed ---

@test "ctrl-a + increases speed" {
	run_headless --steps 1000 \
		--keys "100:ctrl-a,101:+" \
		--watch-comp Speed \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp Speed .* -> 2' "$OUT/state.log"
}

@test "ctrl-a = increases speed" {
	run_headless --steps 1000 \
		--keys "100:ctrl-a,101:=" \
		--watch-comp Speed \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp Speed .* -> 2' "$OUT/state.log"
}

@test "ctrl-a - at minimum speed stays at 1" {
	run_headless --steps 1000 \
		--keys "100:ctrl-a,101:-" \
		--watch-comp Speed \
		"$DISK"
	[[ $status -eq 0 ]]
	[[ ! -f "$OUT/state.log" ]] || ! grep -q 'comp Speed' "$OUT/state.log"
}

@test "ctrl-a _ decreases speed" {
	# First increase to 2, then decrease with _ back to 1
	run_headless --steps 1000 \
		--keys "100:ctrl-a,101:+,200:ctrl-a,201:_" \
		--watch-comp Speed \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp Speed .* -> 2' "$OUT/state.log"
	grep -q 'comp Speed .* -> 1' "$OUT/state.log"
}

@test "speed clamps at maximum of 5" {
	# Increase 5 times (1->2->3->4->5->5), then check no transition beyond 5
	run_headless --steps 2000 \
		--keys "100:ctrl-a,101:+,200:ctrl-a,201:+,300:ctrl-a,301:+,400:ctrl-a,401:+,500:ctrl-a,501:+" \
		--watch-comp Speed \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp Speed .* -> 5' "$OUT/state.log"
	! grep -q 'comp Speed .* -> 6' "$OUT/state.log"
}

# --- Volume ---

@test "ctrl-a v mutes audio" {
	run_headless --steps 1000 \
		--keys "100:ctrl-a,101:v" \
		--watch-comp VolumeMuted \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp VolumeMuted .* -> true' "$OUT/state.log"
}

@test "ctrl-a ] increases volume level" {
	run_headless --steps 1000 \
		--keys "100:ctrl-a,101:]" \
		--watch-comp VolumeLevel \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp VolumeLevel .* -> 60' "$OUT/state.log"
}

@test "ctrl-a [ decreases volume level" {
	run_headless --steps 1000 \
		--keys "100:ctrl-a,101:[" \
		--watch-comp VolumeLevel \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp VolumeLevel .* -> 40' "$OUT/state.log"
}

@test "ctrl-a v toggles mute off after muting" {
	run_headless --steps 1000 \
		--keys "100:ctrl-a,101:v,200:ctrl-a,201:v" \
		--watch-comp VolumeMuted \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp VolumeMuted .* -> true' "$OUT/state.log"
	grep -q 'comp VolumeMuted .* -> false' "$OUT/state.log"
}

@test "volume decrease to zero sets muted" {
	# Default volume is 50; five decreases of 10 reaches 0
	run_headless --steps 2000 \
		--keys "100:ctrl-a,101:[,200:ctrl-a,201:[,300:ctrl-a,301:[,400:ctrl-a,401:[,500:ctrl-a,501:[" \
		--watch-comp VolumeLevel,VolumeMuted \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp VolumeMuted .* -> true' "$OUT/state.log"
}

@test "volume increase clears muted" {
	# Mute first, then increase volume -- should clear muted
	run_headless --steps 1000 \
		--keys "100:ctrl-a,101:v,200:ctrl-a,201:]" \
		--watch-comp VolumeMuted,VolumeLevel \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp VolumeMuted .* -> true' "$OUT/state.log"
	grep -q 'comp VolumeMuted .* -> false' "$OUT/state.log"
}

# --- Write Protect ---

@test "ctrl-a w toggles write protect" {
	run_headless --steps 1000 \
		--keys "100:ctrl-a,101:w" \
		--watch-comp WriteProtect \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp WriteProtect .* -> true' "$OUT/state.log"
}

# --- Debugger ---

@test "ctrl-a b enables debugger" {
	run_headless --steps 1000 \
		--keys "100:ctrl-a,101:b" \
		--watch-comp Debugger \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp Debugger .* -> true' "$OUT/state.log"
}

# --- State Slot ---

@test "ctrl-a 3 selects state slot 3" {
	run_headless --steps 1000 \
		--keys "100:ctrl-a,101:3" \
		--watch-comp StateSlot \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp StateSlot .* -> 3' "$OUT/state.log"
}

# --- Quit ---

@test "ctrl-a q exits cleanly" {
	run_headless --steps 100000 \
		--keys "100:ctrl-a,101:q" \
		"$DISK"
	[[ $status -eq 0 ]]
}

# --- Double Ctrl-A ---

@test "double ctrl-a sends literal ctrl-a to machine" {
	run_headless --steps 1000 \
		--keys "100:ctrl-a,101:ctrl-a" \
		--watch-comp KBLastKey \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp KBLastKey .* -> \$01' "$OUT/state.log"
}

# --- Unrecognized key after prefix ---

@test "unrecognized key after prefix produces no state change" {
	run_headless --steps 1000 \
		--keys "100:ctrl-a,101:x" \
		--watch-comp Paused,Speed \
		"$DISK"
	[[ $status -eq 0 ]]
	[[ ! -f "$OUT/state.log" ]] || [[ ! -s "$OUT/state.log" ]]
}

# --- Load Next / Previous Disk ---

@test "ctrl-a n loads next disk" {
	DISK2="$BATS_TEST_DIRNAME/../work/bt1_char.dsk"
	[[ -f "$DISK2" ]] || skip "second disk image not found: $DISK2"
	run_headless --steps 1000 \
		--keys "100:ctrl-a,101:n" \
		--watch-comp DiskIndex \
		"$DISK" "$DISK2"
	[[ $status -eq 0 ]]
	grep -q 'comp DiskIndex .* -> 1' "$OUT/state.log"
}

@test "ctrl-a p loads previous disk" {
	DISK2="$BATS_TEST_DIRNAME/../work/bt1_char.dsk"
	[[ -f "$DISK2" ]] || skip "second disk image not found: $DISK2"
	run_headless --steps 1000 \
		--keys "100:ctrl-a,101:n,200:ctrl-a,201:p" \
		--watch-comp DiskIndex \
		"$DISK" "$DISK2"
	[[ $status -eq 0 ]]
	grep -q 'comp DiskIndex .* -> 1' "$OUT/state.log"
	grep -q 'comp DiskIndex .* -> 0' "$OUT/state.log"
}

# --- Save and Load State ---

@test "ctrl-a s saves state" {
	run_headless --steps 1000 \
		--keys "100:ctrl-a,101:s" \
		"$DISK"
	[[ $status -eq 0 ]]
	[[ -f "$DISK.1.state" ]]
	rm -f "$DISK.1.state"
}

@test "ctrl-a l loads state without crashing when no state file exists" {
	rm -f "$DISK.1.state"
	run_headless --steps 1000 \
		--keys "100:ctrl-a,101:l" \
		"$DISK"
	[[ $status -eq 0 ]]
}

@test "save and load state round-trip preserves register" {
	# Save state at step 100, then load it at step 300.
	# Between save and load, the CPU will have executed ~200 steps,
	# changing register values. After load, KBLastKey should revert
	# to its saved value.
	run_headless --steps 1000 \
		--keys "100:ctrl-a,101:s,200:a,300:ctrl-a,301:l" \
		--watch-comp KBLastKey \
		"$DISK"
	[[ $status -eq 0 ]]
	# After pressing 'a' at step 200, KBLastKey changes.
	# After loading state at step 301, KBLastKey reverts to pre-'a' value.
	# We check that KBLastKey changed at least twice (set by 'a', then reverted by load).
	local changes
	changes=$(grep -c 'comp KBLastKey' "$OUT/state.log")
	[[ $changes -ge 2 ]]
	rm -f "$DISK.1.state"
}

# --- Parse errors ---

@test "invalid step number fails" {
	run_headless --steps 1000 \
		--keys "abc:ctrl-a" \
		"$DISK"
	[[ $status -ne 0 ]]
}

@test "empty keyspec fails" {
	run_headless --steps 1000 \
		--keys "100:" \
		"$DISK"
	[[ $status -ne 0 ]]
}

@test "duplicate step numbers fail" {
	run_headless --steps 1000 \
		--keys "100:ctrl-a,100:esc" \
		"$DISK"
	[[ $status -ne 0 ]]
}
