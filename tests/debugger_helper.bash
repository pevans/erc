DISK="$BATS_TEST_DIRNAME/../data/memreg.dsk"

setup_file() {
	if [[ ! -f "$DISK" ]]; then
		skip "disk image not found: $DISK"
	fi

	ERC="$BATS_FILE_TMPDIR/erc"
	export ERC
	(cd "$BATS_TEST_DIRNAME/.." && go build -o "$ERC" .)
}

setup() {
	ERC="$BATS_FILE_TMPDIR/erc"
	SESSION="erc-dbg-$$-$BATS_TEST_NUMBER"
	export SESSION

	tmux new-session -d -s "$SESSION" \
		"$ERC" headless --start-in-debugger --steps 10000000 "$DISK"

	wait_for_prompt
}

teardown() {
	tmux kill-session -t "$SESSION" 2>/dev/null || true
}

# wait_for_prompt waits until "debug> " appears in the pane, or fails on timeout.
wait_for_prompt() {
	local max_retries=15
	local i
	for (( i=0; i<max_retries; i++ )); do
		if tmux capture-pane -p -S - -t "$SESSION" 2>/dev/null | grep -q "debug>"; then
			return 0
		fi
		sleep 0.2
	done
	echo "timeout: debugger prompt did not appear" >&2
	return 1
}

# send_cmd sends a command to the debugger and waits for the next prompt.
send_cmd() {
	local cmd="$1"
	local before after max_retries=30 i

	# Count current prompts before sending
	before=$(tmux capture-pane -p -S - -t "$SESSION" 2>/dev/null | grep -c "debug>" || true)

	tmux send-keys -t "$SESSION" "$cmd" Enter

	# Wait until prompt count increases
	for (( i=0; i<max_retries; i++ )); do
		after=$(tmux capture-pane -p -S - -t "$SESSION" 2>/dev/null | grep -c "debug>" || true)
		if (( after > before )); then
			return 0
		fi
		sleep 0.2
	done
	echo "timeout: debugger did not return prompt after command: $cmd" >&2
	return 1
}

# capture captures the full pane scrollback into $PANE.
capture() {
	PANE=$(tmux capture-pane -p -S - -t "$SESSION" 2>/dev/null)
}
