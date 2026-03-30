---
Specification: 20
Drafted At: 2026-03-30
Authors:
  - Peter Evans
---

# 1. Overview

This spec describes the keyboard device of the Apple II. The keyboard is one
of the simplest peripherals on the machine: it presents a single 7-bit ASCII
character to software through a pair of soft switches and uses a strobe
mechanism to signal when new data is available.

Unlike slot-based peripherals, the keyboard does not occupy an expansion slot.
Its soft switches live in the fixed I/O page at $C000 and $C010.

# 2. Hardware Model

The Apple II keyboard produces 7-bit ASCII codes. When a key is pressed, the
keyboard encoder latches the corresponding ASCII value into a data register
and raises a strobe flag (bit 7) to indicate that new data is available.

The emulated keyboard maintains three pieces of state:

- **LastKey**: the 7-bit ASCII value of the most recently pressed key (bits
  0-6). Default: 0.
- **Strobe**: a single-bit flag stored in bit 7. Set to $80 when a key is
  pressed; cleared when software acknowledges the keypress. Default: 0.
- **KeyDown**: a single-bit flag stored in bit 7. Set to $80 while any key is
  physically held down; cleared to 0 when all keys are released. Default: 0.

All three values are initialized to 0 at power-on.

# 3. Soft Switches

## 3.1. $C000 -- Keyboard Data (Read)

A read of $C000 returns `LastKey | Strobe`. When strobe is set, the value
returned has bit 7 high, indicating that the data is valid and has not yet
been acknowledged. After software clears the strobe (by accessing $C010),
subsequent reads of $C000 still return LastKey, but with bit 7 low.

This address is read-only from the keyboard's perspective. Writes to $C000
are handled by other subsystems (display soft switches share the $C000-$C00F
range for writes).

## 3.2. $C010 -- Any Key Down / Strobe Clear (Read and Write)

A read of $C010 returns `KeyDown` (bit 7 indicates whether any key is
currently pressed) and, as a side effect, clears Strobe to 0. This is the
standard mechanism by which software acknowledges a keypress.

A write to $C010 (any value) also clears Strobe to 0. Some software uses
the write form instead of the read form to clear the strobe without needing
the KeyDown status.

    Address   Trigger   Return Value          Side Effect
    -------   -------   ------------          -----------
    $C000     Read      LastKey | Strobe      None
    $C010     Read      KeyDown               Clear Strobe to 0
    $C010     Write     N/A                   Clear Strobe to 0

# 4. Key Press

When the host environment detects that a key is pressed, the emulator
performs the following:

    LastKey  = ascii_value & $7F
    Strobe   = $80
    KeyDown  = $80

The `& $7F` mask ensures that only 7-bit ASCII values are stored. This
operation must be atomic with respect to concurrent soft switch reads.

# 5. Key Release

When all keys are released (no keys are held down), the emulator clears
KeyDown:

    KeyDown = 0

LastKey and Strobe are not affected by key release. The last key value
persists until the next keypress, and the strobe persists until software
clears it.

The host environment also resets its key repeat timers (initial delay and
repeat interval) when all keys are released. A subsequent keypress begins a
fresh repeat cycle.

# 6. Key Encoding

The Apple II keyboard produces 7-bit ASCII. The emulator must translate
host key events into the codes that Apple II software expects.

## 6.1. Printable Characters

Standard printable ASCII characters ($20-$7E) are passed through directly.
Unshifted letter keys produce lowercase (a-z, $61-$7A). Shift+letter produces
uppercase (A-Z, $41-$5A). Shift also modifies digit and symbol keys to produce
their shifted variants (e.g., Shift+1 produces '!', Shift+; produces ':').

## 6.2. Control Characters

When the Control modifier is held, the key value is masked with $1F:

    key = ascii_value & $1F

This produces control characters $01-$1A for the letters A-Z, matching
the standard ASCII control character mapping.

Note: when both Shift and Control are held, Control takes precedence. The `&
$1F` mask is applied to the unshifted key value, so Shift has no effect on
control characters.

More generally, when multiple modifiers are held simultaneously, exactly one
takes effect. The precedence order is:

1. Control -- produces a control character via `& $1F`
2. Shift -- produces the shifted variant of the key
3. No modifier -- produces the unshifted key value

The host environment must resolve modifier conflicts before delivering the
key event to the emulator. Each key event carries at most one effective
modifier.

## 6.3. Special Keys

Certain non-printable keys map to specific ASCII codes:

    Key           Code    Notes
    ---           ----    -----
    Return        $0D     Carriage return
    Escape        $1B
    Tab           $09
    Space         $20
    Backspace     $08     Same code as Left Arrow
    Delete        $7F

## 6.4. Arrow Keys

The Apple II maps arrow keys to control codes. These are the same codes
that cursor-movement software expects:

    Key           Code    Equivalent Control Key
    ---           ----    ----------------------
    Left Arrow    $08     Ctrl-H (Backspace)
    Right Arrow   $15     Ctrl-U
    Up Arrow      $0B     Ctrl-K
    Down Arrow    $0A     Ctrl-J (Line Feed)

Note that Left Arrow and Backspace share code $08. Software cannot
distinguish between them.

# 7. Key Repeat

The host environment is responsible for generating key repeat events. When
a key is held down:

1. The initial keypress is sent immediately.
2. After an initial delay of 500 ms, the key begins repeating.
3. Subsequent repeats occur every 100 ms.

Each repeat event goes through the same key press path described in section 4,
setting LastKey and Strobe as if the key had been freshly pressed. This allows
Apple II software to see repeated characters through the normal $C000/$C010
polling mechanism.

When the key is released, repeating stops and KeyDown is cleared per section
5.

# 8. Open Bus Behavior

Addresses $C001-$C00F are write-only switches for other subsystems (display
modes, etc.). Reads from these addresses are not driven by those subsystems, so
the data bus floats to whatever was last latched by the keyboard -- that is, the
value that would be returned by a read of $C000 (`LastKey | Strobe`). This is a
side effect of the Apple II's bus architecture, not keyboard functionality per
se, but software occasionally depends on it.

Addresses $C011-$C01F each have dedicated read handlers in other subsystems
(banking, memory mode, peripheral ROM, and display status switches) and do not
exhibit open bus behavior.

# 9. Thread Safety

Keyboard state may be written by the input-handling thread and read by the CPU
emulation thread concurrently. All mutations to LastKey, Strobe, and KeyDown
(in PressKey and ClearKeys) must be protected by a mutex or equivalent
synchronization mechanism to prevent data races.

# 10. Interaction with Shortcuts

The emulator intercepts certain key combinations before they reach the Apple
II keyboard path. If a key event is consumed by the shortcut system, it must
not be forwarded to the keyboard -- LastKey, Strobe, and KeyDown must not be
modified. The shortcut system is described in spec 7.
