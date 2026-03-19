#!/usr/bin/env bats

load debugger_helper

# 5.1 help
@test "help lists available commands" {
	send_cmd "help"
	capture
	[[ "$PANE" == *"get <addr>"* ]]
	[[ "$PANE" == *"step <times>"* ]]
	[[ "$PANE" == *"dbatch start"* ]]
	[[ "$PANE" == *"runfor <seconds>"* ]]
	[[ "$PANE" == *"resume"* ]]
}

# 5.2 status
@test "status shows registers and instructions" {
	send_cmd "status"
	capture
	[[ "$PANE" == *"registers"* ]]
	[[ "$PANE" == *"last instruction"* ]]
	[[ "$PANE" == *"next instruction"* ]]
}

# 5.3 get
@test "get prints value at address" {
	send_cmd "get 0000"
	capture
	[[ "$PANE" == *'address $0000:'* ]]
}

# 5.4 set/get round-trip
@test "set then get returns written value" {
	send_cmd "set 0050 42"
	send_cmd "get 0050"
	capture
	[[ "$PANE" == *'$42'* ]]
}

# 5.5 reg/status round-trip (8-bit register)
@test "reg writes register visible in status" {
	send_cmd "reg a ff"
	send_cmd "status"
	capture
	[[ "$PANE" == *'A:$FF'* ]]
}

# 5.5 reg/status round-trip (16-bit PC)
@test "reg pc writes PC visible in status" {
	send_cmd "reg pc 1234"
	send_cmd "status"
	capture
	[[ "$PANE" == *'PC:$1234'* ]]
}

# 5.6 step
@test "step executes one instruction" {
	send_cmd "step"
	capture
	[[ "$PANE" == *"executed 1 times"* ]]
}

@test "step N executes N instructions" {
	send_cmd "step 5"
	capture
	[[ "$PANE" == *"executed 5 times"* ]]
}

# 5.7 until
@test "until steps to matching instruction" {
	send_cmd "until LDA"
	capture
	[[ "$PANE" == *"stepped over"* ]]
}

# 5.8 state
@test "state shows known state keys" {
	send_cmd "state"
	capture
	[[ "$PANE" == *"Debugger"* ]]
	[[ "$PANE" == *"DisplayText"* ]]
}

# 5.9 keypress
@test "keypress completes without error" {
	send_cmd "keypress C1"
	capture
	# No error output expected; just ensure another prompt appeared (send_cmd handled it)
	[[ "$PANE" == *"debug>"* ]]
}

# 5.10 resume and re-entry via --debug-break
@test "resume runs until breakpoint and re-enters debugger" {
	# Replace session with one that includes --debug-break FA62
	tmux kill-session -t "$SESSION" 2>/dev/null || true
	SESSION="erc-dbg-break-$$-$BATS_TEST_NUMBER"
	export SESSION
	tmux new-session -d -s "$SESSION" \
		"$ERC" headless --start-in-debugger --debug-break FA62 --steps 10000000 "$DISK"
	wait_for_prompt

	send_cmd "resume"
}

# 5.11 quit
@test "quit exits the emulator" {
	tmux send-keys -t "$SESSION" "quit" Enter
	sleep 1
	! tmux has-session -t "$SESSION" 2>/dev/null
}

# 5.12 dbatch
@test "dbatch start and stop" {
	send_cmd "dbatch start"
	capture
	[[ "$PANE" == *"debug batch started"* ]]
	send_cmd "dbatch stop"
	capture
	[[ "$PANE" == *"debug batch stopped"* ]]
}

# 5.13 error handling
@test "empty command shows error" {
	send_cmd ""
	capture
	[[ "$PANE" == *"no command given"* ]]
}

@test "unknown command shows error" {
	send_cmd "foobar"
	capture
	[[ "$PANE" == *'unknown command'* ]]
}

@test "get with no address shows error" {
	send_cmd "get"
	capture
	[[ "$PANE" == *"requires an address"* ]]
}

@test "set with no arguments shows error" {
	send_cmd "set"
	capture
	[[ "$PANE" == *"requires an address and value"* ]]
}

@test "reg with invalid register shows error" {
	send_cmd "reg q ff"
	capture
	[[ "$PANE" == *"invalid register"* ]]
}

@test "reg with invalid value shows error" {
	send_cmd "reg a zz"
	capture
	[[ "$PANE" == *"invalid"* ]]
}

@test "step with negative value shows error" {
	send_cmd "step -1"
	capture
	[[ "$PANE" == *"invalid"* ]]
}

@test "until with short mnemonic shows error" {
	send_cmd "until NO"
	capture
	[[ "$PANE" == *"you must provide a valid instruction"* ]]
}

@test "get with invalid address shows error" {
	send_cmd "get zzzz"
	capture
	[[ "$PANE" == *"invalid address"* ]]
}

@test "set with missing value shows error" {
	send_cmd "set 0050"
	capture
	[[ "$PANE" == *"requires an address and value"* ]]
}

@test "set with invalid value shows error" {
	send_cmd "set 0050 zz"
	capture
	[[ "$PANE" == *"invalid value"* ]]
}

@test "keypress with no argument shows error" {
	send_cmd "keypress"
	capture
	[[ "$PANE" == *"requires a hex ascii value"* ]]
}

@test "keypress with invalid hex shows error" {
	send_cmd "keypress zz"
	capture
	[[ "$PANE" == *"invalid value"* ]]
}

@test "step 0 executes zero times" {
	send_cmd "step 0"
	capture
	[[ "$PANE" == *"executed 0 times"* ]]
}

@test "dbatch with no subcommand shows usage" {
	send_cmd "dbatch"
	capture
	[[ "$PANE" == *"usage"* ]]
}

@test "dbatch with unknown subcommand shows error" {
	send_cmd "dbatch foo"
	capture
	[[ "$PANE" == *"unknown dbatch command"* ]]
}

# 5.14 disk
@test "disk loads image into drive" {
	send_cmd "disk $DISK"
	capture
	[[ "$PANE" == *"loaded"* ]]
	[[ "$PANE" == *"into drive"* ]]
}

# 5.15 writeprotect
@test "writeprotect toggles on and off" {
	send_cmd "writeprotect"
	capture
	[[ "$PANE" == *"write protect on drive 1 is ON"* ]]
	send_cmd "writeprotect"
	capture
	[[ "$PANE" == *"write protect on drive 1 is OFF"* ]]
}
