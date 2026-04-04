---
Specification: 25
Category: Persistence
Drafted At: 2026-04-03
Authors:
  - Peter Evans
---

# 1. Overview

This spec describes the save/load state system, which allows the user to
capture a snapshot of the entire emulator state to a file and restore it
later. The system uses numbered slots so that multiple snapshots can coexist
for the same disk image.

# 2. State File Creation

## 2.1. Naming Convention

State files are named using the pattern:

    DISKNAME.SLOT.state

where DISKNAME is the full path of the first disk image in the disk set and
SLOT is the current slot number (an integer). For example, if the disk is
`/path/to/software.dsk` and the slot is 3, the state file is
`/path/to/software.dsk.3.state`.

## 2.2. Slot Selection

The slot number defaults to 0. The user can select a different slot (0-9)
using the Ctrl-A prefix followed by a digit key. The selected slot persists
until changed again or until the emulator exits.

## 2.3. Save Operation

When the user triggers a save (Ctrl-A S), the emulator:

1. Pauses the process loop and waits for it to stop.
2. Captures a snapshot of all state (section 3).
3. Serializes the snapshot to a file using Go's gob encoding.
4. Resumes the process loop if it was running before the save.

If the file already exists, it is overwritten.

## 2.4. Load Operation

When the user triggers a load (Ctrl-A L), the emulator:

1. Pauses the process loop and waits for it to stop.
2. Reads and deserializes the state file.
3. Restores all state (section 3).
4. Rebuilds internal references that cannot be serialized (section 3.7).
5. Clears keyboard state to avoid stuck keys from the loaded snapshot.
6. Forces a display redraw.
7. Resumes the process loop if it was running before the load.

If the state file does not exist or cannot be decoded, the emulator shows an
error message and continues running without changes.

# 3. Round-Trip Preservation

The following state is captured during save and restored during load. After a
save followed by a load, the emulator should be indistinguishable from its
state at the time of the save.

## 3.1. CPU Registers

All CPU registers are preserved: PC, A, X, Y, S (stack pointer), and P
(processor flags). Internal CPU bookkeeping (cycle counter, current opcode,
operand, effective address, addressing mode) is also preserved.

## 3.2. Keyboard State

The keyboard state flags are preserved: caps lock, key-down indicator, last
key value, and strobe value. After a load, the keyboard state is then cleared
(section 2.4, step 6) to prevent keys that were held at save time from
appearing stuck.

## 3.3. Speed Setting

The emulation speed setting is preserved. After loading, the emulator runs at
the same speed as when the state was saved.

## 3.4. Memory Contents

Both main memory (64 KB) and auxiliary memory (64 KB) are preserved in their
entirety.

## 3.5. Display State Flags

All display-related boolean flags are preserved: AltChar, Col80, DoubleHigh,
Hires, Iou, Mixed, Page2, Store80, and Text. This ensures the display mode
at load time matches the mode at save time.

## 3.6. Drive State

Both disk drives are preserved independently. For each drive, the snapshot
includes:

- Motor state and stepper phase
- Track and sector position
- Latch value, read/write mode, and shift register state
- The disk image data and metadata (name, type, write protection)

The selected drive (1 or 2) is also preserved.

## 3.7. Segment References

Memory segment pointers (which segment is the current read target, write
target, etc.) cannot be serialized directly. After restoring all state, the
emulator rebuilds these pointers based on the restored boolean flags
(MemReadAux, MemWriteAux, BankSysBlockAux). The core segment references
(main, aux, ROM) are always restored to the new computer's segments.

# 4. Slot Isolation

Each slot is a completely independent save file. Loading from slot N reads
only the file for slot N. Saving to slot N writes only the file for slot N.
Operations on one slot have no effect on other slots' files.

# 5. Error Handling

- If a save fails (e.g., filesystem error), the emulator shows an error
  message and continues running. The emulator state is not affected.
- If a load fails (file not found, corrupt data), the
  emulator shows an error message and continues running with its current
  state unchanged.
