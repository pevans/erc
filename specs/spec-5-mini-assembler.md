---
Specification: 5
Category: Tools
Drafted At: 2026-03-09
Authors:
  - Peter Evans
---

# 1. Overview

This spec describes a minimal 65C02 assembler whose sole purpose is to produce
bootable Apple II disk images for black box testing. The assembler accepts
source files written in standard 65C02 assembly syntax and emits a
143,360-byte `.dsk` image with the assembled program placed in track 0. The
boot ROM loads track 0 into memory at `$0800` and jumps to `$0801`, so the
assembled code executes immediately after boot with no DOS or operating system
involved.

# 2. Scope

The assembler is intentionally limited. It supports the 65C02 instruction set,
labels, a handful of directives, and nothing else. There are no macros, no
includes, no segments, no linker, and no relocations. The output is always a
single `.dsk` file.

# 3. Source Format

## 3.1. Lines

A source file is a sequence of lines. Each line is one of:

- A **label definition** (an identifier followed by a colon)
- An **instruction** (a mnemonic, optionally followed by an operand)
- A **directive** (a dot-prefixed keyword, optionally followed by arguments)
- A **blank line** or **comment-only line**

Labels, instructions, and directives may appear on the same line:

```
loop:   LDA #$00    ; clear accumulator
```

## 3.2. Comments

A semicolon (`;`) begins a comment that extends to the end of the line.

## 3.3. Labels

A label is an identifier matching `[A-Za-z_][A-Za-z0-9_]*` followed by a
colon. Labels represent the address of the next assembled byte. A label may
appear on its own line or before an instruction or directive on the same line.

Labels are case-sensitive.

## 3.4. Identifiers as Operands

When an identifier appears as an operand (without a colon), it refers to a
previously or forward-declared label. The assembler resolves all label
references after the first pass.

# 4. Instructions

The assembler supports the 65C02 mnemonics defined in the `instructionNames`
table in `mos/tables.go`. Each mnemonic is three uppercase letters. The
assembler should accept lowercase mnemonics as well but canonicalize them to
uppercase internally. Internal names like NP2 and NP3 (undefined opcode
placeholders) are not valid assembler mnemonics.

## 4.1. Addressing Modes

The addressing mode is determined by the operand syntax:

| Syntax            | Mode                  | Example           |
|-------------------|-----------------------|-------------------|
| *(none)*          | Implied               | `CLC`             |
| `A`               | Accumulator           | `ASL A`           |
| `#$NN`            | Immediate             | `LDA #$FF`        |
| `$NNNN`           | Absolute              | `STA $0400`       |
| `$NNNN,X`         | Absolute X-indexed    | `LDA $0400,X`     |
| `$NNNN,Y`         | Absolute Y-indexed    | `LDA $0400,Y`     |
| `$NN`             | Zero page             | `LDA $10`         |
| `$NN,X`           | Zero page X-indexed   | `LDA $10,X`       |
| `$NN,Y`           | Zero page Y-indexed   | `LDX $10,Y`       |
| `($NNNN)`         | Indirect              | `JMP ($1000)`     |
| `($NN,X)`         | X-indexed indirect    | `LDA ($10,X)`     |
| `($NN),Y`         | Indirect Y-indexed    | `LDA ($10),Y`     |
| `($NN)`           | Zero page indirect    | `LDA ($10)`       |
| *label*           | Relative or absolute  | `BEQ loop`        |

Numeric literals use a `$` prefix for hexadecimal. The assembler distinguishes
zero page from absolute by the value: if the operand fits in one byte and the
instruction supports zero page mode for that opcode, zero page is used;
otherwise absolute is used. An operand written with four hex digits (e.g.,
`$0010`) forces absolute mode even if the value fits in one byte.

For branch instructions, a label operand is assembled as a relative offset.
The assembler reports an error if the target is out of the signed 8-bit range
(-128 to +127 bytes from the instruction following the branch).

## 4.2. Accumulator Disambiguation

Some instructions (ASL, LSR, ROL, ROR, INC, DEC) can operate on the
accumulator or on a memory address. When no operand is given, the assembler
uses accumulator mode. An explicit `A` operand also selects accumulator mode.

# 5. Directives

## 5.1. `.byte`

Emits one or more literal bytes.

```
.byte $FF, $00, $42
```

Arguments are comma-separated hexadecimal values (using the `$` prefix), each
in the range `$00`--`$FF`.

## 5.2. `.word`

Emits one or more 16-bit values in little-endian order.

```
.word $0800, $FFFC
```

## 5.3. `.org`

Sets the current assembly address. This is only valid as the first directive
in the file (before any instructions or data). The default origin is `$0801`.

```
.org $0801
```

If `.org` is not specified, the origin defaults to `$0801`.

## 5.4. `.halt`

A convenience directive that emits `JMP *` -- a jump to the current address,
creating an infinite loop. This is the standard way to end a test program.

```
.halt
```

# 6. Assembly Process

## 6.1. Two-Pass Assembly

The assembler uses two passes:

1. **Pass 1**: Scan all lines, recording label addresses and computing the size
   of each instruction and directive. After this pass, every label has a known
   address.
2. **Pass 2**: Emit bytes for each instruction and directive, resolving label
   references to concrete addresses or offsets.

## 6.2. Error Reporting

Errors include the source file name, line number, and a description. The
assembler halts on the first error. Examples of errors:

- Unknown mnemonic
- Invalid operand syntax
- Duplicate label
- Undefined label reference
- Branch target out of range
- Code exceeds track 0 capacity

# 7. Disk Image Output

## 7.1. Layout

The output is a standard 143,360-byte DOS 3.3 disk image. Only track 0 is
populated; all other tracks are zero-filled.

Track 0 layout:

- **Byte 0** (`$0800`): sector count -- the number of 256-byte sectors the boot
  ROM should load from track 0. This is computed from the assembled code size:
  `ceil(code_size / 256)`, capped at 15.
- **Bytes 1+** (`$0801`+): the assembled program.

The remaining bytes through the end of track 0 (16 sectors x 256 bytes = 4096
bytes total) are zero-padded.

## 7.2. Sector Interleaving

DOS 3.3 `.dsk` images store sectors in logical order, but the Apple II boot
ROM reads physical sectors sequentially. Because of the DOS 3.3 sector
interleave, the physical-to-logical mapping is not 1:1. The assembler must
account for this when writing track 0: assembled bytes that should appear
contiguously in memory must be placed into the correct logical sectors so that
the boot ROM's physical-order reads reconstruct the data correctly.

The assembler uses the existing `a2enc` package to encode track data into the
`.dsk` image. The `a2enc` encoder already handles the DOS 3.3 sector
interleave, so the assembler writes its assembled bytes as a flat 4096-byte
track buffer and lets `a2enc` take care of the physical-to-logical mapping.

## 7.3. Size Limit

The assembled program must fit within track 0 minus the sector count byte:
4095 bytes maximum. The assembler reports an error if the program exceeds
this limit.

# 8. Interface

The assembler is a standalone executable (not a subcommand of `erc`) that
reads one source file and writes one `.dsk` file:

```
erc-assembler input.s -o output.dsk
```

The binary lives in `cmd/erc-assembler/main.go` and is built and installed
independently of the main `erc` binary.

If `-o` is omitted, the assembler writes to stdout. If the input is `-`, the
assembler reads from stdin.

The tool exits with status 0 on success and non-zero on any error.

# 9. Design Considerations

## 9.1. Opcode Table

The assembler needs a table mapping (mnemonic, address mode) pairs to opcodes.
This table can be derived from the existing instruction and address mode
tables in the emulator -- those tables map opcodes to instruction functions
and address modes, so inverting them yields the assembler's lookup table.

## 9.2. No Macro System

Macros, conditional assembly, and includes are explicitly out of scope. Test
programs are short and self-contained. If a pattern recurs across many test
files, it is better handled by the test harness (e.g., a shared preamble
concatenated before assembly) than by adding complexity to the assembler.

## 9.3. Relationship to Existing Packages

The `mos` package defines the opcode-to-instruction and opcode-to-address-mode
tables. The assembler inverts these to build its own mnemonic-to-opcode
lookup. The assembler does not import or depend on the emulator's runtime --
it only needs the static tables.

## 9.4. Test Disk Workflow

The intended workflow for CPU instruction testing is to produce some disk
image that can be used by `erc headless` for testing. The details of that are
left out of scope for this spec.
