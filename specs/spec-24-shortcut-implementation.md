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

## 2.11. Help

- `?` or `h` or `H`: Show the help modal (see section 6).

## 2.12. Unrecognized Key

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

# 6. Help Modal

When the user presses `Ctrl-A ?` or `Ctrl-A h`, the emulator pauses and
displays a modal overlay listing all available shortcuts. The modal remains
visible until the user presses ESC, which dismisses it and resumes the
emulator.

## 6.1. Appearance

The modal is drawn as a filled black rectangle with a 1-pixel white border,
centered on the screen. It should have internal padding of 16 pixels on all
sides. Text is rendered in white using the same font used by `TextOverlay`.

The title "Keyboard Shortcuts" is displayed at the top of the modal, centered
horizontally. Below the title, each shortcut is listed as a line of the form:

```
Ctrl-A ESC    Pause / Resume
Ctrl-A +/-    Speed up / down
Ctrl-A [/]    Volume down / up
Ctrl-A V      Mute / unmute
Ctrl-A W      Write protect
Ctrl-A B      Debugger
Ctrl-A C      Caps lock
Ctrl-A 0-9    Select state slot
Ctrl-A S      Save state
Ctrl-A L      Load state
Ctrl-A N/P    Next / previous disk
Ctrl-A Q      Quit
Ctrl-A ?/H    This help screen
```

The key column is left-aligned and the description column is left-aligned,
separated by enough space to be clearly readable.

## 6.2. Pagination

If the list of shortcuts does not fit within the available height of the modal,
the entries are split across pages. A footer at the bottom of the modal shows
the current page and total pages (e.g. "Page 1 of 2") along with navigation
hints: "LEFT/RIGHT to page, ESC to close".

LEFT arrow and RIGHT arrow move between pages. The arrow keys use the same
rune values defined in the `--keys` spec (LEFT = 0x08, RIGHT = 0x15).

## 6.3. Behavior

When the help modal is activated:

1. The computer is paused (`Paused` state is set to true).
2. The prefix overlay is hidden.
3. The help modal overlay becomes active and is drawn on top of the screen.

While the help modal is active, all key events are routed to the modal
rather than to the shortcut system or the computer. The modal handles:

- ESC: Dismiss the modal and resume the computer.
- LEFT (0x08): Go to the previous page if one exists.
- RIGHT (0x15): Go to the next page if one exists.
- Any other key: Ignored.

The help modal state should be tracked with a new boolean state key
(`HelpModal`) so that `shortcut.Check` can detect when the modal is active and
route events to it.

## 6.4. Rendering

The help modal is implemented as a new overlay type in the `gfx` package
(`HelpOverlay`). A global instance (`gfx.HelpModal`) is added alongside the
existing overlays. The render pipeline in `render/ebiten.go` draws it after all
other overlays so it appears on top.

The overlay does not fade. It is either fully visible or hidden.

# 7. Files

| File | Change |
|------|--------|
| `cmd/headless.go` | Add `--keys` flag, parse keyspecs, dispatch events in loop |
| `a2/a2state/keys.go` | Add `Speed`, `VolumeMuted`, `VolumeLevel`, `WriteProtect`, `StateSlot`, `DiskIndex`, `HelpModal` keys |
| `a2/computer.go` | Expose speed/volume/caps lock as observable state |
| `gfx/helpoverlay.go` | New file: `HelpOverlay` type with Show/Hide/Draw/Update and pagination |
| `shortcut/shortcut.go` | Add `?`/`h` case; route keys to help modal when active |
| `render/ebiten.go` | Draw `gfx.HelpModal` in the render pipeline |
