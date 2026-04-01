setup_file() { load debug_artifacts_helper; setup_file; }
setup()      { load debug_artifacts_helper; setup; }
teardown()   { load debug_artifacts_helper; teardown; }

# --- Section 2: Activation ---

@test "debug-image flag produces artifact files" {
	run_debug 50000
	[[ $status -eq 0 ]]
	[[ -f "$DISK.asm" ]]
	[[ -f "$DISK.time" ]]
	[[ -f "$DISK.metrics" ]]
	[[ -f "$DISK.disklog" ]]
	[[ -f "$DISK.screen" ]]
	[[ -f "$DISK.audio" ]]
	[[ -f "$DISK.physical" ]]
}

@test "no artifacts without debug-image flag" {
	run "$ERC_BIN" headless --steps 50000 "$DISK"
	[[ $status -eq 0 ]]
	[[ ! -f "$DISK.asm" ]]
	[[ ! -f "$DISK.time" ]]
	[[ ! -f "$DISK.metrics" ]]
	[[ ! -f "$DISK.physical" ]]
}

@test "disk swap flushes old artifacts and initializes new ones" {
	local session="erc-dbg-$$-$BATS_TEST_NUMBER"
	local DISK2="$TMP/test2.dsk"
	cp "$DISK" "$DISK2"

	tmux new-session -d -s "$session" \
		"$ERC_BIN" headless --start-in-debugger --steps 10000000 --debug-image "$DISK"

	# Wait for debugger prompt
	local i
	for (( i=0; i<15; i++ )); do
		if tmux capture-pane -p -S - -t "$session" 2>/dev/null | grep -q "debug>"; then
			break
		fi
		sleep 0.2
	done

	_send() {
		local before after
		before=$(tmux capture-pane -p -S - -t "$session" 2>/dev/null | grep -c "debug>" || true)
		tmux send-keys -t "$session" "$1" Enter
		for (( i=0; i<30; i++ )); do
			after=$(tmux capture-pane -p -S - -t "$session" 2>/dev/null | grep -c "debug>" || true)
			if (( after > before )); then return 0; fi
			sleep 0.2
		done
		return 1
	}

	# Execute some steps to generate disk activity for the first image
	_send "step 50000"

	# Load a second disk image mid-session
	_send "disk $DISK2"

	# Execute some steps on the new image
	_send "step 50000"

	tmux send-keys -t "$session" "quit" Enter
	sleep 1
	tmux kill-session -t "$session" 2>/dev/null || true

	# First image's disklog should have been flushed
	[[ -f "$DISK.disklog" ]]
	# Second image should have its own artifacts
	[[ -f "$DISK2.asm" ]]
	[[ -f "$DISK2.physical" ]]

	rm -f "$TMP"/test2.dsk*
}

# --- Section 3.1: Instruction Log (.asm) ---

@test "asm file begins with preamble comments" {
	run_debug 50000
	[[ $status -eq 0 ]]
	head -1 "$DISK.asm" | grep -q '^\*'
}

@test "asm preamble contains format explanation" {
	run_debug 50000
	[[ $status -eq 0 ]]
	grep -q 'PROGRAM COUNTER' "$DISK.asm"
	grep -q 'OPCODE' "$DISK.asm"
	grep -q 'OPERAND' "$DISK.asm"
}

@test "asm file contains instruction lines with expected format" {
	run_debug 50000
	[[ $status -eq 0 ]]
	# Match ADDR:OP pattern (e.g. "C100:4C 13 C2")
	grep -qE '^[0-9A-F]{4}:[0-9A-F]{2}' "$DISK.asm"
}

@test "asm hex values are uppercase" {
	run_debug 50000
	[[ $status -eq 0 ]]
	# Address and opcode/operand bytes should use uppercase hex
	# Extract just the "ADDR:OP BBBBB" prefix before the pipe
	! grep -E '^[0-9A-Fa-f]{4}:' "$DISK.asm" | sed 's/ |.*//' | grep -q '[a-f]'
}

@test "asm file contains speculative instructions" {
	run_debug 50000
	[[ $status -eq 0 ]]
	grep -q '(speculative)' "$DISK.asm"
}

@test "asm file has blank lines after block-ending instructions" {
	run_debug 50000
	[[ $status -eq 0 ]]
	# JMP, RTS, RTI, or BRK should be followed by a blank line
	grep -A1 '| .*JMP \|| .*RTS\|| .*RTI\|| .*BRK' "$DISK.asm" | grep -q '^$'
}

@test "asm instructions are sorted by address" {
	run_debug 50000
	[[ $status -eq 0 ]]
	# Extract addresses from instruction lines, check they're sorted
	local addrs
	addrs=$(grep -oE '^[0-9A-F]{4}:' "$DISK.asm" | sed 's/://')
	local sorted
	sorted=$(echo "$addrs" | sort)
	[[ "$addrs" == "$sorted" ]]
}

@test "asm preamble mentions MAIN label convention" {
	run_debug 50000
	[[ $status -eq 0 ]]
	# The preamble explains the MAIN label convention for address 0801
	grep -q 'MAIN' "$DISK.asm"
}

@test "asm real instruction replaces speculative at same address" {
	run_debug 500000
	[[ $status -eq 0 ]]
	# Collect addresses that appear in instruction lines (not comments/blanks)
	local dupes
	dupes=$(grep -oE '^[0-9A-F]{4}:' "$DISK.asm" | sort | uniq -d)
	# No address should appear twice -- real replaces speculative
	[[ -z "$dupes" ]]
	# And we know speculation occurred (tested separately), so replacement
	# must have happened for any overlapping addresses
}

# --- Section 3.2: Instruction Timing (.time) ---

@test "time file contains timing entries" {
	run_debug 50000
	[[ $status -eq 0 ]]
	[[ -s "$DISK.time" ]]
}

@test "time entries have expected format" {
	run_debug 50000
	[[ $status -eq 0 ]]
	# Format: INSTRUCTION | run COUNT cyc CYCLES spent DURATION
	grep -qE '\| run [0-9]+ cyc [0-9]+ spent ' "$DISK.time"
}

@test "time entries are sorted alphabetically" {
	run_debug 50000
	[[ $status -eq 0 ]]
	local lines sorted
	lines=$(cat "$DISK.time")
	sorted=$(sort "$DISK.time")
	[[ "$lines" == "$sorted" ]]
}

@test "time instruction field matches short-string format" {
	run_debug 50000
	[[ $status -eq 0 ]]
	# Short-string: "ADDR | MNEM OPERAND"
	grep -qE '^[0-9A-F]{4} \| [A-Z]{3}' "$DISK.time"
}

# --- Section 3.3: Metrics (.metrics) ---

@test "metrics file contains key-value pairs" {
	run_debug 50000
	[[ $status -eq 0 ]]
	[[ -s "$DISK.metrics" ]]
	# Every non-empty line should be "key = value"
	! grep -v '^$' "$DISK.metrics" | grep -qvE '^[a-zA-Z0-9_]+ = [0-9]+$'
}

@test "metrics file contains instructions count" {
	run_debug 50000
	[[ $status -eq 0 ]]
	grep -q '^instructions = 50000$' "$DISK.metrics"
}

@test "metrics keys are sorted alphabetically" {
	run_debug 50000
	[[ $status -eq 0 ]]
	local keys sorted
	keys=$(awk -F' = ' '{print $1}' "$DISK.metrics")
	sorted=$(echo "$keys" | sort)
	[[ "$keys" == "$sorted" ]]
}

# --- Section 3.4: Disk Log (.disklog) ---

@test "disklog file is created" {
	run_debug 50000
	[[ $status -eq 0 ]]
	[[ -f "$DISK.disklog" ]]
}

@test "disklog file is non-empty" {
	run_debug 500000
	[[ $status -eq 0 ]]
	[[ -s "$DISK.disklog" ]]
}

@test "disklog entries match expected format" {
	run_debug 500000
	[[ $status -eq 0 ]]
	# Format: [ELAPSED   ] MODE T:TT S:S P:PPPP B:$BB | INSTRUCTION
	local pattern='^\[.{10}\] (RD|WR) T:[0-9A-F]{2} S:[0-9A-F] P:[0-9A-F]{4} B:\$[0-9A-F]{2} \|'
	# Every line must match
	! grep -vE "$pattern" "$DISK.disklog" | grep -q .
}

@test "disklog entries contain RD mode" {
	run_debug 500000
	[[ $status -eq 0 ]]
	grep -q '] RD ' "$DISK.disklog"
}

@test "disklog elapsed time is left-aligned in brackets" {
	run_debug 500000
	[[ $status -eq 0 ]]
	# Elapsed field is in brackets, left-aligned in 10-char field
	head -1 "$DISK.disklog" | grep -qE '^\[[0-9]'
}

# Note: WR mode entries in disklog cannot be tested with memreg.dsk because
# the image never initiates disk writes. Write mode is triggered by the
# running software (via soft switch $C0EF), not by emulator flags.

# --- Section 3.5: Screen Log (.screen) ---

@test "screen file is created" {
	run_debug 50000
	[[ $status -eq 0 ]]
	[[ -f "$DISK.screen" ]]
}

@test "screen file contains FRAME header" {
	run_debug 500000
	[[ $status -eq 0 ]]
	grep -qE '^FRAME [0-9]+\.[0-9]{6}$' "$DISK.screen"
}

@test "screen frame has 192 rows" {
	run_debug 500000
	[[ $status -eq 0 ]]
	# Extract lines between FRAME header and blank line (one frame)
	local count
	count=$(sed -n '/^FRAME/,/^$/{ /^FRAME/d; /^$/d; p; }' "$DISK.screen" | head -192 | wc -l)
	[[ $count -eq 192 ]]
}

@test "screen rows are 280 characters wide" {
	run_debug 500000
	[[ $status -eq 0 ]]
	# Check that content rows (non-FRAME, non-blank) are 280 chars
	local bad
	bad=$(sed -n '/^FRAME/,/^$/{ /^FRAME/d; /^$/d; p; }' "$DISK.screen" | \
		head -192 | awk '{ if (length($0) != 280) print NR }')
	[[ -z "$bad" ]]
}

@test "screen pixels use only valid characters" {
	run_debug 500000
	[[ $status -eq 0 ]]
	# Valid chars: W B O G P and space
	! sed -n '/^FRAME/,/^$/{ /^FRAME/d; /^$/d; p; }' "$DISK.screen" | \
		grep -q '[^WBOGP ]'
}

@test "screen frames are separated by blank lines" {
	run_debug 500000
	[[ $status -eq 0 ]]
	# A FRAME line should be preceded by a blank line (except the first)
	local count
	count=$(grep -c '^FRAME' "$DISK.screen")
	if (( count > 1 )); then
		# Every FRAME after the first should be preceded by blank line
		grep -B1 '^FRAME' "$DISK.screen" | grep -c '^$' | \
			grep -q "$(( count - 1 ))"
	fi
}

# --- Section 3.6: Audio Log (.audio) ---

@test "audio file begins with header" {
	run_debug 50000
	[[ $status -eq 0 ]]
	head -1 "$DISK.audio" | grep -q 'Audio Log - Sample Rate: 44100 Hz'
}

@test "audio file contains frame duration line" {
	run_debug 50000
	[[ $status -eq 0 ]]
	grep -q 'Each frame represents 1.0 second of audio' "$DISK.audio"
}

@test "audio file contains activity timeline legend" {
	run_debug 50000
	[[ $status -eq 0 ]]
	grep -q 'Activity Timeline Legend' "$DISK.audio"
}

@test "audio file contains FRAME header" {
	run_debug 500000
	[[ $status -eq 0 ]]
	grep -qE '^FRAME [0-9]+\.[0-9]{6}$' "$DISK.audio"
}

@test "audio frame contains sample statistics" {
	run_debug 500000
	[[ $status -eq 0 ]]
	grep -qE '^\s+Samples: [0-9]+, Min: ' "$DISK.audio"
}

@test "audio frame contains analysis metrics" {
	run_debug 500000
	[[ $status -eq 0 ]]
	grep -qE '^\s+Zero Crossings: [0-9]+, Max Run: [0-9]+ samples, Activity: ' "$DISK.audio"
}

@test "audio frame contains activity timeline" {
	run_debug 500000
	[[ $status -eq 0 ]]
	grep -qE '^\s+Timeline: ' "$DISK.audio"
}

@test "audio frame sample count is 44100" {
	run_debug 500000
	[[ $status -eq 0 ]]
	grep -q 'Samples: 44100,' "$DISK.audio"
}

@test "audio frame contains waveform visualization" {
	run_debug 500000
	[[ $status -eq 0 ]]
	# Waveform is 20 lines of 80 columns, each indented with 2 spaces.
	# Lines contain only spaces and waveform characters.
	local count
	count=$(sed -n '/^FRAME/,/^$/{
		/^FRAME/d
		/^$/d
		/Samples:/d
		/Zero Crossings:/d
		/Timeline:/d
		/^  .\{80\}$/p
	}' "$DISK.audio" | head -20 | wc -l)
	[[ $count -eq 20 ]]
}

# --- Section 3.7: Physical Disk Image (.physical) ---

@test "physical file is non-empty binary" {
	run_debug 50000
	[[ $status -eq 0 ]]
	[[ -s "$DISK.physical" ]]
}

@test "physical file differs from original dsk" {
	run_debug 50000
	[[ $status -eq 0 ]]
	# The physical (nibblized) image should differ from the logical .dsk
	! cmp -s "$DISK" "$DISK.physical"
}

@test "physical file is written at load time" {
	# Run zero steps -- the physical image should still be written
	run_debug 0
	[[ $status -eq 0 ]]
	[[ -s "$DISK.physical" ]]
}

# --- Section 3.8: Instruction Diff Map (.diff.asm) ---

@test "diff.asm is not created without debug batch" {
	run_debug 50000
	[[ $status -eq 0 ]]
	[[ ! -f "$DISK.diff.asm" ]]
}

@test "diff.asm is written at shutdown when batch still active" {
	local session="erc-dbg-$$-$BATS_TEST_NUMBER"

	tmux new-session -d -s "$session" \
		"$ERC_BIN" headless --start-in-debugger --steps 10000000 --debug-image "$DISK"

	local i
	for (( i=0; i<15; i++ )); do
		if tmux capture-pane -p -S - -t "$session" 2>/dev/null | grep -q "debug>"; then
			break
		fi
		sleep 0.2
	done

	_send() {
		local before after
		before=$(tmux capture-pane -p -S - -t "$session" 2>/dev/null | grep -c "debug>" || true)
		tmux send-keys -t "$session" "$1" Enter
		for (( i=0; i<30; i++ )); do
			after=$(tmux capture-pane -p -S - -t "$session" 2>/dev/null | grep -c "debug>" || true)
			if (( after > before )); then return 0; fi
			sleep 0.2
		done
		return 1
	}

	# Start a batch but do NOT stop it before quitting
	_send "dbatch start"
	_send "step 1000"

	tmux send-keys -t "$session" "quit" Enter
	sleep 1
	tmux kill-session -t "$session" 2>/dev/null || true

	# Should still be written at shutdown
	[[ -f "$DISK.diff.asm" ]]
	[[ -s "$DISK.diff.asm" ]]
	grep -qE '^[0-9A-F]{4}:[0-9A-F]{2}' "$DISK.diff.asm"
}

@test "diff.asm is created after debug batch" {
	local session="erc-dbg-$$-$BATS_TEST_NUMBER"

	tmux new-session -d -s "$session" \
		"$ERC_BIN" headless --start-in-debugger --steps 10000000 --debug-image "$DISK"

	# Wait for debugger prompt
	local i
	for (( i=0; i<15; i++ )); do
		if tmux capture-pane -p -S - -t "$session" 2>/dev/null | grep -q "debug>"; then
			break
		fi
		sleep 0.2
	done

	# Helper to send a command and wait for the next prompt
	_send() {
		local before after
		before=$(tmux capture-pane -p -S - -t "$session" 2>/dev/null | grep -c "debug>" || true)
		tmux send-keys -t "$session" "$1" Enter
		for (( i=0; i<30; i++ )); do
			after=$(tmux capture-pane -p -S - -t "$session" 2>/dev/null | grep -c "debug>" || true)
			if (( after > before )); then return 0; fi
			sleep 0.2
		done
		return 1
	}

	_send "dbatch start"
	_send "step 1000"
	_send "dbatch stop"

	# quit calls os.Exit, so don't wait for a prompt
	tmux send-keys -t "$session" "quit" Enter
	sleep 1
	tmux kill-session -t "$session" 2>/dev/null || true

	[[ -f "$DISK.diff.asm" ]]
	[[ -s "$DISK.diff.asm" ]]
	# Should contain instruction lines in the same format as .asm
	grep -qE '^[0-9A-F]{4}:[0-9A-F]{2}' "$DISK.diff.asm"
}

# --- Section 5: Shutdown Sequence ---

@test "all artifacts are written in a single run" {
	run_debug 50000
	[[ $status -eq 0 ]]
	# metrics, asm, time are all written at shutdown
	[[ -s "$DISK.metrics" ]]
	[[ -s "$DISK.asm" ]]
	[[ -s "$DISK.time" ]]
}
