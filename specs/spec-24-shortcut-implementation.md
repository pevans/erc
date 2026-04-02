---
Specification: 24
Category: User Interface
Drafted At: 2026-04-01
Authors:
  - Peter Evans
---

# 1. Overview

This spec describes the implementation of the shortcut system and the details
for adding key injection support to the headless execution loop. It
accompanies spec 7, which defines the testing interface and expected shortcut
behaviors.

# 2. Shortcut System

All shortcuts use a Ctrl-A prefix. The user presses Ctrl-A, which activates a
prefix overlay, then presses a second key to trigger the shortcut. While the
prefix is active, the following keys are recognized:

## 2.1. Pause and Resume

- ESC: Pause the computer and show a pause graphic. While paused, ESC
  resumes and any other key re-flashes the pause graphic.

## 2.2. Speed

- `+` or `=`: Increase emulation speed by one step (max 5).
- `-` or `_`: Decrease emulation speed by one step (min 1).

## 2.3. Volume

- `v` or `V`: Toggle mute.
- `]` or `}`: Increase volume by 10 and unmute.
- `[` or `{`: Decrease volume by 10. When the level reaches 0, mute.

## 2.4. Write Protect

- `w` or `W`: Toggle write protection on the selected drive.

## 2.5. Debugger

- `b` or `B`: Activate the debugger.

## 2.6. Caps Lock

- `c` or `C`: Toggle caps lock. Shows a text status message ("caps lock:
  on" or "caps lock: off").

## 2.7. State Slots

- `0`-`9`: Select the given state slot number.
- `s` or `S`: Save state to the current slot. The save file is written using
  the naming convention `IMAGENAME.SLOT.state`.
- `l` or `L`: Load state from the current slot. If the file does not exist,
  show an error message.

## 2.8. Disk Navigation

- `n` or `N`: Load the next disk image from the disk set.
- `p` or `P`: Load the previous disk image from the disk set.

## 2.9. Quit

- `q` or `Q`: Shut down the computer and exit.

## 2.10. Double Prefix (Ctrl-A Ctrl-A)

If the prefix overlay is already active and Ctrl-A is pressed again, the
overlay is hidden and a literal Ctrl-A (ASCII 0x01) is sent to the machine.

## 2.11. Unrecognized Key

Any key not listed above is silently discarded. The prefix is consumed and no
state changes occur.

# 3. Parsing --keys

The flag value is split on commas. Each entry is split on the first colon to
separate the step number from the keyspec. The keyspec is parsed into an
`input.Event`:

- If the keyspec is `@release`, release all currently held keys instead of
  pressing a new one.
- If the keyspec starts with `ctrl-`, set `Modifier = input.ModControl` and
  parse the remainder as the key rune.
- Named keys: `esc` (0x1B), `return` (0x0D), `tab` (0x09), `space` (0x20),
  `backspace` (0x08), `delete` (0x7F), `left` (0x08), `right` (0x15),
  `up` (0x0B), `down` (0x0A).
- Otherwise, the keyspec is a single character used as the key rune.

Invalid entries (bad step number, empty keyspec, or duplicate step numbers)
cause the command to fail with an error before execution begins.

# 4. Event Dispatch in the Headless Loop

The headless loop should be extended to check for scheduled key events at each
step. A sorted slice or map of step-to-events works. Before calling
`comp.Process()`, any events for the current step are dispatched through
`shortcut.Check()`. If not consumed, they fall through to `comp.PressKey()`.
When the event has the control modifier, the key rune is masked with `& 0x1F`
before being passed to `PressKey()`.

# 5. Quit Handling

When `shortcut.Check()` returns an error from `comp.Shutdown()`, the headless
loop should break early and exit cleanly rather than calling `fail()`.

# 6. Files

| File | Change |
|------|--------|
| `cmd/headless.go` | Add `--keys` flag, parse keyspecs, dispatch events in loop |
| `a2/a2state/keys.go` | Add `Speed`, `VolumeMuted`, `VolumeLevel`, `WriteProtect`, `StateSlot`, `DiskIndex` keys |
| `a2/computer.go` | Expose speed/volume/caps lock as observable state |
