---
Specification: 8
Drafted At: 2026-03-18
Authors:
  - Peter Evans
---

# 1. Overview

This spec describes a mechanism for black-box testing the interactive debugger.
The debugger today runs inside the graphical `erc run` command and uses a
readline prompt (`debug> `) on the terminal. It reads input from stdin and
writes output (prefixed with ` --> `) to stdout. Testing it requires a way to
send commands and read responses without a human at the keyboard.

The approach uses `tmux` in headless mode to create a detached session that
runs `erc headless --start-in-debugger`, then sends keystrokes to the session
and captures pane output to verify debugger behavior. Using the headless
command avoids any need for a graphical display.

# 2. Problem

The debugger has 15+ commands (get, set, reg, status, step, until, help, quit,
resume, etc.) and none of them have any automated test coverage. Unit testing
the debugger is awkward because it is tightly coupled to a `liner.State` for
readline and to a live `a2.Computer`. Black-box testing via tmux sidesteps both
problems: the debugger runs in its real environment, and the tests interact
with it the same way a user would.

# 3. Approach: tmux Sessions

## 3.1. Why tmux

The debugger uses the `liner` library, which requires a real terminal (it reads
from `/dev/tty` and uses terminal escape sequences for line editing). Piping
stdin does not work with `liner`. A tmux session provides a real PTY that
`liner` is happy with, while still being fully scriptable.

tmux can run without a display server. A detached session (`tmux new-session
-d`) works in CI and on headless machines. The only requirement is that `tmux`
is installed.

## 3.2. Session Lifecycle

Each test follows this pattern:

1. **Create** a detached tmux session running `erc headless
   --start-in-debugger` with a disk image.
2. **Wait** for the `debug> ` prompt to appear in the pane.
3. **Send** one or more debugger commands using `tmux send-keys`.
4. **Capture** the pane contents using `tmux capture-pane -p`.
5. **Assert** that expected output strings appear in the captured text.
6. **Tear down** the session with `tmux kill-session`.

## 3.3. Waiting for the Prompt

After starting the session, tests must wait for the debugger to initialize
before sending commands. A polling loop checks `tmux capture-pane -p` for the
presence of `debug> `. The loop should have a timeout (e.g. 3 seconds) to
avoid hanging if the emulator fails to start.

Pseudocode:

```
for i in 1..MAX_RETRIES:
    output = capture_pane(session)
    if "debug> " in output:
        return success
    sleep POLL_INTERVAL
return failure("debugger prompt did not appear")
```

A reasonable default is MAX_RETRIES=15, POLL_INTERVAL=0.2s, giving a 3-second
window.

## 3.4. Sending Commands

`tmux send-keys -t SESSION "command text" Enter` types the command and presses
Enter. After sending a command, the test should wait for the next `debug> `
prompt to appear (indicating the command has completed) before capturing output
or sending the next command.

Since each command produces a new `debug> ` prompt, the wait-for-prompt helper
can count prompts: after sending the Nth command, wait until N+1 prompts are
visible.

## 3.5. Capturing Output

`tmux capture-pane -p -t SESSION` prints the current visible pane content to
stdout. The `-S -` and `-E -` flags can capture the entire scrollback buffer
if needed.

# 4. Test Helper

A new helper file, `test/debugger_helper.bash`, provides reusable functions for
the bats tests.

## 4.1. Helper Functions

- `setup_file`: builds the `erc` binary (same as `test_helper.bash`).
- `setup`: creates a unique tmux session name (e.g. `erc-dbg-$$-$BATS_TEST_NUMBER`)
  and starts a detached tmux session running `erc headless --start-in-debugger`
  with `$DISK` (inherited from `test_helper.bash`) and a generous `--steps`
  budget. Calls `wait_for_prompt` before returning.
- `teardown`: kills the tmux session unconditionally.
- `wait_for_prompt`: polls until a `debug> ` prompt appears, or fails with a
  timeout.
- `send_cmd CMD`: sends a command string to the tmux session and waits for the
  next prompt.
- `capture`: captures the pane output into a variable (e.g. `$PANE`).

## 4.2. Session Naming

Each test needs its own session to avoid collisions when bats runs tests in
parallel. The session name should incorporate the PID and test number, e.g.
`erc-dbg-$$-$BATS_TEST_NUMBER`.

# 5. Testable Commands

## 5.1. help

Send `help` and verify that the output contains the list of commands. Assert
that key strings like `get <addr>`, `step <times>`, `dbatch start`, `runfor
<seconds>`, and `resume` appear.

## 5.2. status

Send `status` and verify that the output contains `registers`, `last
instruction`, and `next instruction` lines.

## 5.3. get

Send `get 0000` and verify the output contains `address $0000:` followed by a
hex value.

## 5.4. set and get round-trip

Send `set 0050 42`, then `get 0050`, and verify the output contains `$42`.

## 5.5. reg and status round-trip

Send `reg a ff`, then `status`, and verify the status output contains `A:$FF`.
The `status` command prints registers in the format `A:$XX`, so matching
`A:$FF` avoids false positives from hex values elsewhere in the output.

Also test `reg pc 1234`, then `status`, and verify the output contains
`PC:$1234`. This exercises a 16-bit register (PC) in addition to the 8-bit
register (A) tested above.

## 5.6. step

Send `step` (single step) and verify the output contains `executed 1 times`.
Send `step 5` and verify the output contains `executed 5 times`.

## 5.7. until

Send `until LDA` and verify the output contains `stepped over`. LDA is one of
the most common instructions and will be hit very early in boot.

The `until` command has a built-in cap of 100,000,000 iterations. If the cap is
reached, it prints `stopped after ... instructions without executing`. This
does not need a dedicated test since it would be extremely slow.

## 5.8. state

Send `state` and verify the output contains known state keys (e.g. `Debugger`,
`Paused`).

## 5.9. keypress

Send `keypress C1` (uppercase 'A' with high bit set, as Apple II expects) and
verify the command completes without error.

## 5.10. resume and re-entry via --debug-break

Use `--debug-break` with address `FA62`. This is deep inside the Apple //e boot
ROM reset vector and is hit very early during startup. Start the session with
`erc headless --start-in-debugger --debug-break FA62`, then send `resume`. The
emulator runs until it hits the breakpoint, at which point the debugger
re-enters. Verify by waiting for a new `debug> ` prompt to appear after resume.

This test needs its own setup since it passes extra flags. Override setup to
create the tmux session with `--debug-break FA62`.

## 5.11. quit

Send `quit` and verify the tmux pane shows the session has ended (the process
exits). This can be detected by checking that `tmux has-session -t SESSION`
fails after a short wait.

## 5.12. dbatch

Send `dbatch start` and verify the output contains `debug batch started`. Send
`dbatch stop` and verify the output contains `debug batch stopped`.

## 5.13. Error handling

- Send an empty line (just Enter) and verify the output contains `no command
  given`.
- Send an unknown command (e.g. `foobar`) and verify the output contains
  `unknown command`.
- Send `get` with no address and verify the output contains `requires an
  address`.
- Send `set` with no arguments and verify the output contains `requires an
  address and value`.
- Send `reg q ff` (invalid register name) and verify the output contains
  `invalid register`.
- Send `reg a zz` and verify the output contains `invalid`.
- Send `step -1` and verify the output contains `must be positive`.
- Send `until NO` (fewer than 3 characters) and verify the output contains
  `you must provide a valid instruction`.
- Send `get zzzz` and verify the output contains `invalid address`.
- Send `set 0050` (missing value) and verify the output contains `requires an
  address and value`.
- Send `set 0050 zz` and verify the output contains `invalid value`.
- Send `keypress` with no argument and verify the output contains `requires a
  hex ascii value`.
- Send `keypress zz` and verify the output contains `invalid value`.
- Send `step 0` and verify the output contains `executed 0 times` (not an
  error; zero is non-negative).
- Send `dbatch` with no subcommand and verify the output contains `usage`.
- Send `dbatch foo` and verify the output contains `unknown dbatch command`.

## 5.14. disk

Send `disk` with `$DISK` (the boot disk path from the test helper) and verify
the output contains `loaded` and `into drive`. The actual output is
`loaded FILE into drive`.

## 5.15. writeprotect

Send `writeprotect` and verify the output contains `write protect on drive 1 is
ON`. Send it again and verify the output contains `write protect on drive 1 is
OFF`.

## 5.16. runfor

The `runfor` command is not tested. It resumes emulation for a timed duration
using a background goroutine and calls `gfx.ShowStatus`, which makes it
difficult to exercise meaningfully in headless/tmux mode.

# 6. Bats Test Structure

Tests live in `test/debugger.bats` and load `debugger_helper.bash`.

## 6.1. Example: help

```bash
@test "help lists available commands" {
    send_cmd "help"
    capture
    [[ "$PANE" == *"get <addr>"* ]]
    [[ "$PANE" == *"step <times>"* ]]
    [[ "$PANE" == *"resume"* ]]
}
```

## 6.2. Example: get

```bash
@test "get prints value at address" {
    send_cmd "get 0000"
    capture
    [[ "$PANE" == *"address \$0000"* ]]
}
```

## 6.3. Example: set/get round-trip

```bash
@test "set then get returns written value" {
    send_cmd "set 0050 42"
    send_cmd "get 0050"
    capture
    [[ "$PANE" == *'$42'* ]]
}
```

## 6.4. Example: reg/status round-trip

```bash
@test "reg writes register visible in status" {
    send_cmd "reg a ff"
    send_cmd "status"
    capture
    [[ "$PANE" == *'A:$FF'* ]]
}
```

## 6.5. Example: reg pc/status round-trip

```bash
@test "reg pc writes PC visible in status" {
    send_cmd "reg pc 1234"
    send_cmd "status"
    capture
    [[ "$PANE" == *'PC:$1234'* ]]
}
```

## 6.6. Example: step

```bash
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
```

## 6.7. Example: dbatch

```bash
@test "dbatch start and stop" {
    send_cmd "dbatch start"
    capture
    [[ "$PANE" == *"debug batch started"* ]]
    send_cmd "dbatch stop"
    capture
    [[ "$PANE" == *"debug batch stopped"* ]]
}
```

## 6.8. Example: unknown command

```bash
@test "unknown command shows error and help" {
    send_cmd "foobar"
    capture
    [[ "$PANE" == *'unknown command: "foobar"'* ]]
}
```

## 6.9. Example: error cases

```bash
@test "empty command shows error" {
    send_cmd ""
    capture
    [[ "$PANE" == *"no command given"* ]]
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

@test "dbatch with no subcommand shows usage" {
    send_cmd "dbatch"
    capture
    [[ "$PANE" == *"usage"* ]]
}
```

## 6.10. Example: quit

```bash
@test "quit exits the emulator" {
    tmux send-keys -t "$SESSION" "quit" Enter
    sleep 1
    ! tmux has-session -t "$SESSION" 2>/dev/null
}
```

## 6.11. Example: additional error cases

```bash
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
```

## 6.12. Example: writeprotect toggle

```bash
@test "writeprotect toggles on and off" {
    send_cmd "writeprotect"
    capture
    [[ "$PANE" == *"write protect on drive 1 is ON"* ]]
    send_cmd "writeprotect"
    capture
    [[ "$PANE" == *"write protect on drive 1 is OFF"* ]]
}
```

## 6.13. Example: disk

```bash
@test "disk loads image into drive" {
    send_cmd "disk $DISK"
    capture
    [[ "$PANE" == *"loaded"* ]]
    [[ "$PANE" == *"into drive"* ]]
}
```

# 7. Implementation Notes

## 7.1. Adding Debugger Support to Headless

The `erc headless` command currently runs a fixed step loop and exits. To
support the debugger, two new flags are added:

- `--start-in-debugger`: enter the debugger prompt before executing any steps.
- `--debug-break ADDRS`: comma-separated hex addresses; when the PC hits one,
  the debugger is entered.

Since `erc headless` never opens a graphical window, these tests work in CI
and under tmux without any display-related workarounds.

The `--steps` flag remains required. If the step budget is exhausted while the
debugger is not active, the process exits normally. If the debugger is active
when steps run out, the debugger remains open until the user issues `resume` or
`quit`.

### 7.1.1. Step Loop Integration

The graphical `erc run` command uses the `clock.Emulator` process loop, which
calls `breakpointCheckFunc` on every iteration and enters the debugger whenever
the `a2state.Debugger` state is true. The headless command does not use
`clock.Emulator` -- it has its own `for i := range headlessStepsFlag` loop.

To integrate the debugger, the headless step loop must add two checks on each
iteration:

1. **Breakpoint check.** If `--debug-break` was provided, check whether
   `debug.HasBreakpoint(int(comp.CPU.PC))` is true before executing each step.
   If so, set `comp.State.SetBool(a2state.Debugger, true)`.

2. **Debugger entry.** If `comp.State.Bool(a2state.Debugger)` is true, call
   `debug.Prompt(comp, line)` in a loop until the debugger state is set back to
   false (i.e., the user runs `resume`). While the debugger is active, the step
   counter does not advance.

This mirrors the logic in `clock.Emulator.ProcessLoop` (lines 94-105 of
`clock/emulator.go`) but adapted for the headless step-counted loop rather than
the clock-timed loop.

If `--start-in-debugger` is set, the headless command sets
`comp.State.SetBool(a2state.Debugger, true)` before the step loop begins, so
the debugger entry check triggers on the first iteration.

### 7.1.2. Liner Integration

The headless command must create a `liner.State` to provide the `debug> `
readline prompt, just as `erc run` does. The liner is created once before the
step loop:

```go
line := liner.NewLiner()
defer line.Close()
```

`debug.Prompt(comp, line)` reads one line of input from the terminal, appends
it to history, and calls `execute`. The headless step loop calls `Prompt` in a
loop while the debugger is active:

```go
for comp.State.Bool(a2state.Debugger) {
    debug.Prompt(comp, line)
}
```

This works under tmux because the tmux PTY provides a real terminal that
`liner` can read from. It does not work with plain stdin piping (liner requires
a TTY), which is why the tests use tmux rather than shell pipes.

## 7.2. Scrollback vs. Visible Pane

By default, `tmux capture-pane -p` only captures the visible portion of the
pane (typically 24 lines). For commands that produce long output (like `state`
or `until`), use `tmux capture-pane -p -S -` to include scrollback. The helper
should use scrollback capture by default.

## 7.3. Prompt Counting

A robust `send_cmd` implementation counts the number of `debug> ` strings in
the captured output. Before sending a command, record the current count. After
sending, wait until the count increments. This avoids races where the test
reads stale output.

## 7.4. Parallel Execution

bats can run test files in parallel, but tests within a file run sequentially.
Since each test gets its own tmux session (via unique session names), parallel
file execution is safe. Tests within `debugger.bats` should be independent --
each test starts a fresh session.

## 7.5. CI Requirements

The CI environment must have `tmux` installed. Most Linux CI images include it
or can install it with `apt-get install tmux`. macOS runners have it available
via Homebrew.

# 8. Files

| File | Change |
|------|--------|
| `cmd/headless.go` | Add `--start-in-debugger` and `--debug-break` flags; integrate debugger prompt into step loop |
| `test/debugger_helper.bash` | New -- tmux session management and helper functions |
| `test/debugger.bats` | New -- black-box debugger tests |

# 9. What Does Not Change

- `debug/*.go` -- the debugger package is tested as-is, without modification.
- `cmd/run.go` -- the graphical run command is not modified.
- `render/` -- the draw loop is not involved; headless never opens a window.
- Existing bats tests and helpers.

# 10. Verification

1. `just lint` passes.
2. `go test ./...` passes.
3. `tmux` is available on the test machine.
4. All bats tests in `test/debugger.bats` pass.
5. Existing bats tests continue to pass (no regressions).
