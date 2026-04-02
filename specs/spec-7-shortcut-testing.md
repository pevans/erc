---
Specification: 7
Category: Tests
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

All shortcuts defined in spec 24 are testable via `--keys` and observable
through `--watch-comp` and `--watch-reg`. Each shortcut's keys and behavior
are defined there; this spec only adds testing-specific notes.

- Shortcuts that modify observable state (speed, volume, pause, etc.) can be
  verified by watching the corresponding state key before and after injection.
- Save/load state can be round-trip tested by changing a register value
  between save and load, then confirming the register reverts to the saved
  value.
- Quit causes the headless run to exit cleanly (exit status 0), potentially
  producing fewer steps than `--steps` requested.
- Disk navigation shortcuts (`n`, `p`) require multiple disk images passed to
  `erc headless`.

# 6. Verification

1. `just lint` passes.
2. `go test ./...` passes.
3. All bats tests in `test/headless_shortcuts.bats` pass.
4. Existing bats tests in `test/headless_state.bats`, `test/headless_audio.bats`,
   and `test/headless_video.bats` continue to pass (no regressions).
