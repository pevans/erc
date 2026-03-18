---
Specification: 7
Drafted At: 2026-03-17
Authors:
  - Peter Evans
---

# 1. Overview

This spec describes a mechanism for injecting key events into the headless
execution loop, enabling black-box testing of the shortcut system defined in
the `shortcut` package. The approach adds a `--keys` flag to the existing `erc
headless` subcommand and uses the existing `--watch-comp` and `--watch-reg`
observers to verify shortcut effects.

# 2. Problem

The shortcut system processes `input.Event` values and modifies computer state
(pause, speed, volume, write-protect, etc.). Today, there is no way to inject
key events during a headless run, so none of this behavior can be tested
without a graphical window.

# 3. Key Injection

## 3.1. Flag Syntax

A new `--keys` flag accepts a comma-separated list of timed key events. Each
entry has the form `STEP:KEYSPEC`, where:

- `STEP` is the decimal step number at which the event is injected.
- `KEYSPEC` describes the key and optional modifier.

Examples:

```
--keys "100:ctrl-a,101:esc"
--keys "500:ctrl-a,501:+,800:ctrl-a,801:-"
```

## 3.2. Key Spec Format

A keyspec is one of:

- A single printable character: `a`, `q`, `+`, `-`, `[`, `]`, `1`-`9`, etc.
- A named key: `esc` (maps to rune `0x1B`).
- A modifier prefix followed by a hyphen and a key: `ctrl-a`, `ctrl-A`.

Only the `ctrl` modifier is needed for shortcut testing. The modifier maps to
`input.ModControl`.

Characters are case-sensitive: `q` produces `Event{Key: 'q', Modifier:
ModNone}`, while `Q` produces `Event{Key: 'Q', Modifier: ModNone}`. The
shortcut handler already accepts both cases for letter shortcuts.

## 3.3. Injection Point

Key events are injected immediately before the instruction at the given step
executes. This means step 0 injects before the first instruction. If multiple
events share the same step number, they are processed in the order they appear
in the flag value.

Each step number must be unique within the flag value. To send multiple keys
in sequence, use successive step numbers (e.g. `"100:ctrl-a,101:esc"`). The
command should fail with an error if duplicate step numbers are given.

Each injected event is passed through `shortcut.Check()` using the same code
path as the graphical `run` command. If `shortcut.Check()` returns `false`
(the event was not consumed by the shortcut system), the key is passed to
`comp.PressKey()` so that the emulated machine receives it.

# 4. New Observable State

Some shortcuts modify state that is not currently exposed via `--watch-comp`.
To make these testable, the following state keys should be added to
`a2state` and to the `headlessStateNameToKey` map in `cmd/headless.go`:

| State name     | Type | Source                           |
|----------------|------|----------------------------------|
| `Speed`        | int  | `Computer.speed`                 |
| `VolumeMuted`  | bool | `Computer.volumeMuted`           |
| `VolumeLevel`  | int  | `Computer.volumeLevel`           |
| `WriteProtect` | bool | `Computer.SelectedDrive().WriteProtected()` |
| `StateSlot`    | int  | `Computer.stateSlot`             |
| `DiskIndex`    | int  | `Computer.Disks.CurrentIndex()`  |

These join the existing set (e.g. `Paused`, `Debugger`, `KBLastKey`).

# 5. Testable Shortcuts

Each shortcut and its observable effects:

## 5.1. Pause and Resume

- Ctrl-A ESC pauses the computer. Observe `Paused` changing to `true`.
- While paused, ESC resumes. Observe `Paused` changing back to `false`.
- While paused, any non-ESC key keeps `Paused` as `true`.

## 5.2. Speed Up and Speed Down

- Ctrl-A `+` (or `=`) increases `Speed` by 1 (up to 5).
- Ctrl-A `-` (or `_`) decreases `Speed` by 1 (down to 1).
- Repeated speed-up at the maximum stays at 5.
- Repeated speed-down at the minimum stays at 1.

## 5.3. Volume

- Ctrl-A `v` toggles `VolumeMuted`.
- Ctrl-A `]` increases `VolumeLevel` by 10 and clears `VolumeMuted`.
- Ctrl-A `[` decreases `VolumeLevel` by 10. When the level reaches 0,
  `VolumeMuted` becomes `true`.

## 5.4. Write Protect Toggle

- Ctrl-A `w` toggles `WriteProtect` on the selected drive.

## 5.5. Debugger

- Ctrl-A `b` sets `Debugger` to `true`.

## 5.6. Load Next and Load Previous

- Ctrl-A `n` loads the next disk image from the disk set. Observe `DiskIndex`
  incrementing.
- Ctrl-A `p` loads the previous disk image. Observe `DiskIndex` decrementing.
- These require multiple disk images passed to `erc headless`.

## 5.7. Save and Load State

- Ctrl-A `s` saves state to the current slot. The save file is written to the
  disk image's directory using the naming convention
  `IMAGENAME.SLOT.state`.
- Ctrl-A `l` loads state from the current slot. If the file does not exist,
  the computer shows an error message but does not crash.
- A round-trip test can verify save/load by changing a register value between
  save and load, then confirming the register reverts to the saved value.

## 5.8. State Slot Selection

- Ctrl-A followed by a digit `1`-`9` selects that state slot. Observe
  `StateSlot` changing to the given digit.

## 5.9. Quit

- Ctrl-A `q` causes the headless run to exit cleanly (exit status 0). The run
  may produce fewer steps than `--steps` requested.

## 5.10. Prefix Pass-Through (Double Ctrl-A)

- Ctrl-A followed immediately by another Ctrl-A sends a literal `0x01` to the
  machine. Observe `KBLastKey` changing to `1`.

## 5.11. Unrecognized Key After Prefix

- Ctrl-A followed by a key that is not a recognized shortcut (e.g. `x`)
  should produce no state changes. The prefix is consumed and the key is
  discarded.

# 6. Bats Test Structure

Tests live in `test/headless_shortcuts.bats` and use the existing
`test_helper.bash` (which builds erc and provides `run_headless`).

Each test:

1. Calls `run_headless` with `--keys`, `--steps`, and the appropriate
   `--watch-comp` flag.
2. Asserts exit status.
3. Greps `state.log` for the expected state transitions.

## 6.1. Example: Pause and Resume

```bash
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
    # Paused goes true then back to false
    grep -q 'comp Paused .* -> true' "$OUT/state.log"
    grep -q 'comp Paused .* -> false' "$OUT/state.log"
}
```

## 6.2. Example: Speed

```bash
@test "ctrl-a + increases speed" {
    run_headless --steps 1000 \
        --keys "100:ctrl-a,101:+" \
        --watch-comp Speed \
        "$DISK"
    [[ $status -eq 0 ]]
    grep -q 'comp Speed .* -> 2' "$OUT/state.log"
}
```

## 6.3. Example: Double Ctrl-A

```bash
@test "double ctrl-a sends literal ctrl-a to machine" {
    run_headless --steps 1000 \
        --keys "100:ctrl-a,101:ctrl-a" \
        --watch-comp KBLastKey \
        "$DISK"
    [[ $status -eq 0 ]]
    grep -q 'comp KBLastKey .* -> 1' "$OUT/state.log"
}
```

## 6.4. Example: Volume

```bash
@test "ctrl-a v mutes audio" {
    run_headless --steps 1000 \
        --keys "100:ctrl-a,101:v" \
        --watch-comp VolumeMuted \
        "$DISK"
    [[ $status -eq 0 ]]
    grep -q 'comp VolumeMuted .* -> true' "$OUT/state.log"
}

@test "ctrl-a ] increases volume" {
    run_headless --steps 1000 \
        --keys "100:ctrl-a,101:]" \
        --watch-comp VolumeLevel \
        "$DISK"
    [[ $status -eq 0 ]]
    grep -q 'comp VolumeLevel .* -> 60' "$OUT/state.log"
}

@test "ctrl-a [ decreases volume" {
    run_headless --steps 1000 \
        --keys "100:ctrl-a,101:[" \
        --watch-comp VolumeLevel,VolumeMuted \
        "$DISK"
    [[ $status -eq 0 ]]
    grep -q 'comp VolumeLevel .* -> 40' "$OUT/state.log"
}
```

## 6.5. Example: Write Protect

```bash
@test "ctrl-a w toggles write protect" {
    run_headless --steps 1000 \
        --keys "100:ctrl-a,101:w" \
        --watch-comp WriteProtect \
        "$DISK"
    [[ $status -eq 0 ]]
    grep -q 'comp WriteProtect .* -> true' "$OUT/state.log"
}
```

## 6.6. Example: Debugger

```bash
@test "ctrl-a b enables debugger" {
    run_headless --steps 1000 \
        --keys "100:ctrl-a,101:b" \
        --watch-comp Debugger \
        "$DISK"
    [[ $status -eq 0 ]]
    grep -q 'comp Debugger .* -> true' "$OUT/state.log"
}
```

## 6.7. Example: State Slot

```bash
@test "ctrl-a 3 selects state slot 3" {
    run_headless --steps 1000 \
        --keys "100:ctrl-a,101:3" \
        --watch-comp StateSlot \
        "$DISK"
    [[ $status -eq 0 ]]
    grep -q 'comp StateSlot .* -> 3' "$OUT/state.log"
}
```

## 6.8. Example: Load Next Disk

```bash
@test "ctrl-a n loads next disk" {
    run_headless --steps 1000 \
        --keys "100:ctrl-a,101:n" \
        --watch-comp DiskIndex \
        "$DISK" "$DISK2"
    [[ $status -eq 0 ]]
    grep -q 'comp DiskIndex .* -> 1' "$OUT/state.log"
}
```

## 6.9. Example: Quit

```bash
@test "ctrl-a q exits cleanly" {
    run_headless --steps 100000 \
        --keys "100:ctrl-a,101:q" \
        "$DISK"
    [[ $status -eq 0 ]]
}
```

## 6.10. Example: Non-ESC Key While Paused

```bash
@test "non-esc key while paused stays paused" {
    run_headless --steps 1000 \
        --keys "100:ctrl-a,101:esc,200:a" \
        --watch-comp Paused \
        "$DISK"
    [[ $status -eq 0 ]]
    grep -q 'comp Paused .* -> true' "$OUT/state.log"
    # Paused never goes back to false
    ! grep -q 'comp Paused .* -> false' "$OUT/state.log"
}
```

## 6.11. Example: Unrecognized Key After Prefix

```bash
@test "unrecognized key after prefix produces no state change" {
    run_headless --steps 1000 \
        --keys "100:ctrl-a,101:x" \
        --watch-comp Paused,Speed \
        "$DISK"
    [[ $status -eq 0 ]]
    # state.log should either not exist or contain no entries
    [[ ! -f "$OUT/state.log" ]] || [[ ! -s "$OUT/state.log" ]]
}
```

## 6.12. Example: Parse Errors

```bash
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
```

# 7. Implementation Notes

## 7.1. Parsing --keys

The flag value is split on commas. Each entry is split on the first colon to
separate the step number from the keyspec. The keyspec is parsed into an
`input.Event`:

- If the keyspec starts with `ctrl-`, set `Modifier = input.ModControl` and
  parse the remainder as the key rune.
- If the keyspec is `esc`, set `Key = 0x1B`.
- Otherwise, the keyspec is a single character used as the key rune.

Invalid entries (bad step number, empty keyspec) cause the command to fail
with an error before execution begins.

## 7.2. Event Dispatch in the Headless Loop

The headless loop should be extended to check for scheduled key events at each
step. A sorted slice or map of step-to-events works. Before calling
`comp.Process()`, any events for the current step are dispatched through
`shortcut.Check()`. If not consumed, they fall through to `comp.PressKey()`.

## 7.3. Quit Handling

When `shortcut.Check()` returns an error from `comp.Shutdown()`, the headless
loop should break early and exit cleanly rather than calling `fail()`.

# 8. Files

| File | Change |
|------|--------|
| `cmd/headless.go` | Add `--keys` flag, parse keyspecs, dispatch events in loop |
| `a2/a2state/keys.go` | Add `Speed`, `VolumeMuted`, `VolumeLevel`, `WriteProtect`, `StateSlot`, `DiskIndex` keys |
| `a2/computer.go` | Expose speed/volume as observable state (getters or state map entries) |
| `test/headless_shortcuts.bats` | New -- black-box shortcut tests |

# 9. What Does Not Change

- `shortcut/shortcut.go` -- the shortcut handler is tested as-is.
- `input/event.go` -- the `Event` type is reused without modification.
- `gfx/` -- overlay show/hide calls are harmless in headless mode (they set
  struct fields but never draw).
- Existing headless tests and flags.

# 10. Verification

1. `just lint` passes.
2. `go test ./...` passes.
3. All bats tests in `test/headless_shortcuts.bats` pass.
4. Existing bats tests in `test/headless_state.bats`, `test/headless_audio.bats`,
   and `test/headless_video.bats` continue to pass (no regressions).
