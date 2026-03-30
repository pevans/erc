setup_file() { load test_helper; setup_file; }
setup()      { load test_helper; setup; }
teardown()   { load test_helper; teardown; }

# --- Section 4: Key Press ---

@test "key press sets LastKey" {
	run_headless --steps 1000 \
		--keys "100:a" \
		--watch-comp KBLastKey \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp KBLastKey .* -> \$61' "$OUT/state.log"
}

@test "key press sets Strobe" {
	run_headless --steps 1000 \
		--keys "100:a" \
		--watch-comp KBStrobe \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp KBStrobe .* -> \$80' "$OUT/state.log"
}

@test "key press sets KeyDown" {
	run_headless --steps 1000 \
		--keys "100:a" \
		--watch-comp KBKeyDown \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp KBKeyDown .* -> \$80' "$OUT/state.log"
}

# --- Section 6.1: Printable Characters ---

@test "lowercase letter produces lowercase code" {
	run_headless --steps 1000 \
		--keys "100:z" \
		--watch-comp KBLastKey \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp KBLastKey .* -> \$7A' "$OUT/state.log"
}

@test "uppercase letter produces uppercase code" {
	run_headless --steps 1000 \
		--keys "100:A" \
		--watch-comp KBLastKey \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp KBLastKey .* -> \$41' "$OUT/state.log"
}

@test "digit key produces digit code" {
	run_headless --steps 1000 \
		--keys "100:1" \
		--watch-comp KBLastKey \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp KBLastKey .* -> \$31' "$OUT/state.log"
}

# --- Section 5: Key Release ---

@test "key release clears KeyDown" {
	run_headless --steps 500 \
		--keys "100:a,200:@release" \
		--watch-comp KBKeyDown \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp KBKeyDown .* -> \$00' "$OUT/state.log"
}

# --- Section 6.2: Control Characters ---

@test "ctrl key produces control character" {
	run_headless --steps 1000 \
		--keys "100:ctrl-b" \
		--watch-comp KBLastKey \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp KBLastKey .* -> \$02' "$OUT/state.log"
}

@test "ctrl with shifted letter produces same control character" {
	run_headless --steps 1000 \
		--keys "100:ctrl-B" \
		--watch-comp KBLastKey \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp KBLastKey .* -> \$02' "$OUT/state.log"
}

# --- Section 6.3: Special Keys ---

@test "space key produces space code" {
	run_headless --steps 1000 \
		--keys "100:space" \
		--watch-comp KBLastKey \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp KBLastKey .* -> \$20' "$OUT/state.log"
}

@test "return key produces carriage return code" {
	run_headless --steps 1000 \
		--keys "100:return" \
		--watch-comp KBLastKey \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp KBLastKey .* -> \$0D' "$OUT/state.log"
}

@test "escape key produces escape code" {
	run_headless --steps 1000 \
		--keys "100:esc" \
		--watch-comp KBLastKey \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp KBLastKey .* -> \$1B' "$OUT/state.log"
}

@test "tab key produces tab code" {
	run_headless --steps 1000 \
		--keys "100:tab" \
		--watch-comp KBLastKey \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp KBLastKey .* -> \$09' "$OUT/state.log"
}

@test "backspace key produces backspace code" {
	run_headless --steps 1000 \
		--keys "100:backspace" \
		--watch-comp KBLastKey \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp KBLastKey .* -> \$08' "$OUT/state.log"
}

@test "delete key produces delete code" {
	run_headless --steps 1000 \
		--keys "100:delete" \
		--watch-comp KBLastKey \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp KBLastKey .* -> \$7F' "$OUT/state.log"
}

# --- Section 6.4: Arrow Keys ---

@test "left arrow produces backspace code" {
	run_headless --steps 1000 \
		--keys "100:left" \
		--watch-comp KBLastKey \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp KBLastKey .* -> \$08' "$OUT/state.log"
}

@test "right arrow produces ctrl-u code" {
	run_headless --steps 1000 \
		--keys "100:right" \
		--watch-comp KBLastKey \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp KBLastKey .* -> \$15' "$OUT/state.log"
}

@test "up arrow produces ctrl-k code" {
	run_headless --steps 1000 \
		--keys "100:up" \
		--watch-comp KBLastKey \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp KBLastKey .* -> \$0B' "$OUT/state.log"
}

@test "down arrow produces ctrl-j code" {
	run_headless --steps 1000 \
		--keys "100:down" \
		--watch-comp KBLastKey \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp KBLastKey .* -> \$0A' "$OUT/state.log"
}

# --- Section 10: Interaction with Shortcuts ---

@test "ctrl-a prefix does not modify keyboard state" {
	run_headless --steps 1000 \
		--keys "100:ctrl-a" \
		--watch-comp KBLastKey,KBStrobe,KBKeyDown \
		"$DISK"
	[[ $status -eq 0 ]]
	[[ ! -f "$OUT/state.log" ]] || ! grep -q 'comp KB' "$OUT/state.log"
}
