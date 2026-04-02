---
Specification: 23
Category: Tests
Drafted At: 2026-03-31
Authors:
  - Peter Evans
---

# 1. Overview

When the `--debug-image` flag is passed to the `run` command, the emulator
produces a collection of debug artifact files alongside the disk image being
run. These files capture detailed information about what happens during
emulation: instructions executed, disk operations, screen output, audio,
timing, and metrics.

All artifact filenames are derived from the disk image filename by appending an
extension. For an image named `game.dsk`, the artifacts would be `game.dsk.asm`,
`game.dsk.time`, and so on.

# 2. Activation

Debug artifacts are initialized when a disk image is loaded while the debug
image flag is set. If a new image is loaded during the same session (e.g. by
swapping disks), the previous disk log is flushed to disk, and a new set of
artifacts is initialized for the new image.

All artifacts except the physical image (section 3.7) are written at shutdown.
The physical image is written immediately at load time.

# 3. Artifact Files

## 3.1. Instruction Log (.asm)

**Extension**: `.asm`
**Format**: Text (MOS 6502 assembly notation)
**Written at**: Shutdown

The instruction map is a deduplicated record of every instruction the CPU
executed during emulation. Instructions are sorted by address. Each line
contains the program counter, raw opcode and operand bytes, an optional label,
the instruction mnemonic, a formatted operand, and an optional comment.

The format of a line is:

    ADDR:OP BBBBB | LABEL    MNEM OPERAND           COMMENT

The fields are:

- **ADDR**: 4-character hex address, zero-padded
- **OP**: 2-character hex opcode, zero-padded
- **BBBBB**: operand bytes, left-aligned in a 5-character field (e.g. `D4 0C`,
  `0C`, or blank)
- **LABEL**: left-aligned in an 8-character field
- **MNEM**: instruction mnemonic, variable width
- **OPERAND**: formatted operand, left-aligned in a 10-character field
- **COMMENT**: optional, variable width, preceded by 5 characters of spacing

Hex values are uppercase. The operand bytes are encoded least-significant byte
first, matching the machine's byte order.

For example:

    0801:4C 04 08 | MAIN     JMP $0804

Some lines are marked `(speculative)` in their comment field. These lines
represent instructions that were never executed but would have been reached if
a branch had been taken in the opposite direction. Speculation helps reveal
code that exists in memory near executed branches. It is not perfect: the
bytes after an untaken branch may be data rather than code. This may be
because a branch can be a cheaper form of jump, so long as the target is a
nearby address in memory.

Instructions that represent the end of a block (JMP, JMP indirect, JMP
absolute-indexed-indirect, RTS, RTI, or BRK) are followed by a blank line to
visually separate subroutines.

When the same address is executed with both a speculative and a real
instruction, the real instruction replaces the speculative one.

The file begins with a preamble of comments (prefixed with `*`) that explains
the format.

## 3.2. Instruction Timing (.time)

**Extension**: `.time`
**Format**: Text (columnar)
**Written at**: Shutdown

The timing file records aggregate execution statistics for each unique
instruction. An "instruction" here means a specific combination of address,
mnemonic, and operand -- the short-string form of a line from the instruction
map.

Each line has two fields separated by ` | `:

    INSTRUCTION                    | run COUNT cyc CYCLES spent DURATION

The fields are:

- **INSTRUCTION**: the short-string form of the instruction (address, mnemonic,
  and operand), left-aligned in a 30-character field
- **COUNT**: number of times this instruction was executed, variable width
- **CYCLES**: total CPU cycles consumed across all executions, variable width
- **DURATION**: estimated wall-clock time, computed as CYCLES multiplied by the
  time-per-cycle value of the clock emulator, variable width

The short-string form used in INSTRUCTION has the format `ADDR | MNEM OPERAND`,
where ADDR is 4-character hex, MNEM is variable width, and OPERAND is
left-aligned in a 10-character field.

Lines are sorted alphabetically by the INSTRUCTION field, which implies that
instructions are sorted by the address at which they are run (like the
instruction log) given the ADDR field is the first member of the field.

## 3.3. Metrics (.metrics)

**Extension**: `.metrics`
**Format**: Text (key-value pairs)
**Written at**: Shutdown

The metrics file contains global counters that track emulation events. Each
line is a key-value pair:

    key = value

Keys are sorted alphabetically. Values are integers representing cumulative
counts since boot. The set of keys is not fixed; various subsystems increment
counters as they operate. Examples of keys include:

- `instructions` -- total instructions executed
- `renders` -- total screen renders
- `disk_read`, `disk_write` -- disk I/O operations
- `bad_opcodes` -- unrecognized opcodes encountered
- Soft switch activations such as `soft_display_hires_on`,
  `soft_memory_read_aux_on`, `soft_read_speaker_toggle`, etc.

## 3.4. Disk Log (.disklog)

**Extension**: `.disklog`
**Format**: Text (timestamped log entries)
**Written at**: Shutdown (and when a new image is loaded mid-session)

The disk log records every individual byte read from or written to the disk
drive. Each line has the format:

    [ELAPSED   ] MODE T:TT S:S P:PPPP B:$BB | INSTRUCTION

The fields are:

- **ELAPSED**: wall-clock time since boot, rounded to the nearest millisecond,
  left-aligned in a 10-character field, enclosed in brackets
- **MODE**: `RD` for a read, `WR` for a write (2 characters)
- **T**: track number, zero-padded 2-character hex
- **S**: sector number, 1-character hex
- **P**: byte offset from the start of the track, zero-padded 4-character hex
- **B**: the byte value, zero-padded 2-character hex, prefixed with `$`
- **INSTRUCTION**: the instruction and operand that caused the disk operation,
  variable width

## 3.5. Screen Log (.screen)

**Extension**: `.screen`
**Format**: Text (ASCII-art frames)
**Written at**: Shutdown

The screen log captures periodic snapshots of the display. A frame is captured
at most once per second of wall-clock time.

Each frame begins with a header line:

    FRAME TIMESTAMP

where TIMESTAMP is a floating-point number (seconds since boot, six decimal
places). The header is followed by 192 lines of 280 characters each,
representing the Apple //e display at its native resolution of 280x192 pixels.

Each character represents one pixel. The character mapping is:

| Character | Color        |
|-----------|--------------|
| `W`       | White        |
| `B`       | Blue         |
| `O`       | Orange       |
| `G`       | Green        |
| `P`       | Purple       |
| (space)   | Black        |

Green and purple each match several specific RGB values to account for
different shades used in various graphics modes.

Frames are separated by a blank line.

## 3.6. Audio Log (.audio)

**Extension**: `.audio`
**Format**: Text (statistics and ASCII waveform)
**Written at**: Shutdown

The audio log captures audio data in one-second frames at a 44,100 Hz sample
rate (44,100 samples per frame). The file begins with a header:

    Audio Log - Sample Rate: 44100 Hz
    Each frame represents 1.0 second of audio

This is followed by a legend explaining the activity timeline characters.

Each frame contains:

1. **Header**: `FRAME TIMESTAMP` (seconds since boot, six decimal places)
2. **Sample statistics**: sample count, minimum and maximum amplitude, and
   average absolute amplitude. Amplitudes are 32-bit floats in the range
   -1.0 to 1.0.
3. **Analysis metrics**:
   - **Zero Crossings**: number of times the waveform crosses zero (an
     indicator of frequency content)
   - **Max Run**: longest consecutive sequence of identical sample values (an
     indicator of dropouts)
   - **Activity**: percentage of 100ms windows that contained more than 10
     unique sample values
4. **Activity timeline**: a string where each character represents a 100ms
   window of audio. A one-second frame produces 10 characters, padded with
   spaces to a fixed width of 80:
   - `█` = active audio (>10 unique values)
   - `▒` = moderate activity (6-10 unique values)
   - `░` = low activity (3-5 unique values)
   - `·` = idle or dropout (<=2 unique values)
5. **Waveform visualization**: a 20-line by 80-column ASCII art waveform. Each
   column represents a time slice; the height represents RMS amplitude,
   plotted symmetrically around the center line. Characters indicate amplitude
   level: `█` (>0.8), `▓` (>0.6), `▒` (>0.4), `░` (>0.2), `·` (<=0.2).

Frames are separated by a blank line.

## 3.7. Physical Disk Image (.physical)

**Extension**: `.physical`
**Format**: Binary (raw disk data)
**Written at**: Image load time

The physical disk image is a binary dump of the drive's data segment after
encoding. When a `.dsk` file is loaded, the emulator encodes it from logical
sector format into the physical nibble format that the drive hardware actually
reads. This file captures that encoded result.

This artifact is useful for inspecting the exact nibble-level encoding that the
emulated drive sees, which differs from the logical format stored in the `.dsk`
file.

## 3.8. Instruction Diff Map (.diff.asm)

**Extension**: `.diff.asm`
**Format**: Text (MOS 6502 assembly notation, same as section 3.1)
**Written at**: When a debug batch ends, or at shutdown if a batch is active

The instruction diff map uses the same format as the instruction map (section
3.1), but only records instructions that executed during a "debug batch"
session. A debug batch is started and stopped explicitly (e.g. via the
debugger). This allows a user to isolate the instructions executed during a
specific window of time, rather than for the entire emulation session.

The diff map is written and reset when the batch is stopped. If a batch is
still active at shutdown, it is written then. If no debug batch is started
during emulation, this file is not written.

# 4. Data Collection

Instruction and timing data are collected via a channel. The CPU sends each
executed instruction to an instruction channel, and a background goroutine
reads from this channel to update the instruction map, the time set, and (if a
debug batch is active) the instruction diff map.

Disk operations, screen frames, and audio samples are recorded directly by
their respective subsystems as they occur.

Metrics are incremented by various subsystems throughout execution using a
global, mutex-protected counter map.

# 5. Shutdown Sequence

At shutdown, the artifact files are written in the following order:

1. Metrics
2. Instruction map
3. Instruction timing
4. Disk log
5. Screen log
6. Audio log
7. Instruction diff map (if a debug batch is still active)

If any file fails to write, the shutdown sequence stops and returns the error.
The drive images are saved after all debug artifacts have been written.
