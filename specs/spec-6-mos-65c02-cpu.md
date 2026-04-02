---
Specification: 6
Category: CPU
Drafted At: 2026-03-12
Authors:
  - Peter Evans
---

# 1. Overview

This spec describes the behavior of the MOS 65C02 CPU as emulated in the `mos`
package. The 65C02 is an 8-bit processor with a 16-bit address space, a small
set of registers, and a fixed-length instruction encoding of 1-3 bytes per
instruction. The `mos` package is a self-contained emulation of this chip;
it does not own memory directly but operates on memory provided by the caller.

# 2. Registers

The CPU has six programmer-visible registers and two pieces of internal state
that are accessible for debugging.

## 2.1. Program Counter (PC)

A 16-bit register that holds the address of the next instruction to execute.
After each instruction, PC is advanced by the number of bytes the instruction
occupies (1, 2, or 3). Instructions that alter control flow (branches, jumps,
calls, returns) set PC directly rather than advancing it incrementally.

## 2.2. Accumulator (A)

An 8-bit general-purpose register used as the primary operand and result
register for arithmetic, logic, and most load/store operations.

## 2.3. Index Registers (X, Y)

Two 8-bit registers used as offsets in indexed addressing modes and as loop
counters. They can be loaded, stored, transferred, incremented, and decremented
independently.

## 2.4. Stack Pointer (S)

An 8-bit register that holds the low byte of the current stack address. The
stack is fixed at page 1 of memory ($0100-$01FF); the full stack address is
always $0100 + S. The stack grows downward: pushing a value writes to the
current address and decrements S; popping increments S and reads from the
resulting address.

## 2.5. Status Register (P)

An 8-bit register where each bit is an independent flag:

| Bit | Mask | Name      | Set when...                                            |
|-----|------|-----------|--------------------------------------------------------|
| 0   | $01  | Carry     | arithmetic produced a carry or borrow out              |
| 1   | $02  | Zero      | the result of an operation was zero                    |
| 2   | $04  | Interrupt | hardware interrupts are disabled                       |
| 3   | $08  | Decimal   | ADC and SBC operate in Binary-Coded Decimal mode       |
| 4   | $10  | Break     | a BRK instruction caused the most recent interrupt     |
| 5   | $20  | Unused    | always set; has no architectural meaning               |
| 6   | $40  | Overflow  | arithmetic produced a signed overflow                  |
| 7   | $80  | Negative  | the most-significant bit of the result is 1            |

## 2.6. Internal State (LastPC, Opcode)

The CPU records the address of the most recently executed instruction (LastPC)
and the opcode byte of the instruction currently being executed. These are
read-only to callers and are useful for debugging and instruction logging.

# 3. Memory Interface

The CPU does not own memory. Instead, the caller supplies two functions:

- **Read(addr uint16) byte** -- returns the byte at the given address.
- **Write(addr uint16, val byte)** -- stores a byte at the given address.

All memory access by the CPU goes through these two primitives. The CPU also
provides a 16-bit read and write that combine two consecutive 8-bit accesses
in little-endian order:

```
Read16(addr):
    lo = Read(addr)
    hi = Read(addr + 1)
    return (hi << 8) | lo

Write16(addr, val):
    Write(addr, val & $FF)
    Write(addr + 1, val >> 8)
```

# 4. Addressing Modes

Every instruction is paired with an addressing mode that determines how the
CPU computes the effective address (EffAddr) and effective value (EffVal)
before the instruction logic runs. The addressing mode also determines the
instruction's byte length.

## 4.1. Implied (IMP)

No operand. EffAddr and EffVal are both zero. The instruction operates entirely
on registers.

*Example:* `CLC`

## 4.2. Accumulator (ACC)

No memory operand. EffVal is the current value of A. Results are written back
to A.

*Example:* `ASL A`

## 4.3. Immediate (IMM)

The operand is a 1-byte literal value embedded in the instruction stream.

```
EffVal = Read(PC + 1)
```

*Example:* `LDA #$FF`

## 4.4. Zero Page (ZPG)

The operand is a 1-byte address in page zero ($0000-$00FF).

```
EffAddr = Read(PC + 1)
EffVal  = Read(EffAddr)
```

*Example:* `LDA $10`

## 4.5. Zero Page X-indexed (ZPX)

The effective address is computed by adding X to the zero-page operand,
wrapping within page zero (no carry into page 1).

```
EffAddr = (Read(PC + 1) + X) & $FF
EffVal  = Read(EffAddr)
```

*Example:* `LDA $10,X`

## 4.6. Zero Page Y-indexed (ZPY)

Same as ZPX but uses Y.

```
EffAddr = (Read(PC + 1) + Y) & $FF
EffVal  = Read(EffAddr)
```

*Example:* `LDX $10,Y`

## 4.7. Zero Page Indirect (ZPI)

The operand is a zero-page address that holds a 16-bit pointer. The CPU reads
the pointer and uses it as the effective address. The pointer read wraps within
page zero: if the operand is $FF, the low byte comes from $FF and the high byte
comes from $00.

```
ptr     = Read(PC + 1)
EffAddr = Read16_ZP(ptr)   ; wraps within page zero
EffVal  = Read(EffAddr)
```

*Example:* `LDA ($10)`

## 4.8. Absolute (ABS)

The operand is a 2-byte address.

```
EffAddr = Read16(PC + 1)
EffVal  = Read(EffAddr)
```

*Example:* `LDA $1234`

## 4.9. Absolute X-indexed (ABX)

The effective address is the 2-byte operand plus X.

```
EffAddr = Read16(PC + 1) + X
EffVal  = Read(EffAddr)
```

Crossing a page boundary (the high byte of the address changes) costs one
extra cycle on read instructions.

*Example:* `LDA $1234,X`

## 4.10. Absolute Y-indexed (ABY)

Same as ABX but uses Y.

```
EffAddr = Read16(PC + 1) + Y
EffVal  = Read(EffAddr)
```

*Example:* `LDA $1234,Y`

## 4.11. Indirect (IND)

The operand is a 2-byte address that holds a 16-bit pointer. Used only by JMP.

```
ptr    = Read16(PC + 1)
EffAddr = Read16(ptr)
```

*Example:* `JMP ($1234)`

## 4.12. X-indexed Indirect (IDX)

The operand is a zero-page address. X is added to it (wrapping within page
zero) to form a zero-page pointer, which is then read to produce EffAddr.

```
ptr     = (Read(PC + 1) + X) & $FF
EffAddr = Read16_ZP(ptr)
EffVal  = Read(EffAddr)
```

*Example:* `LDA ($10,X)`

## 4.13. Indirect Y-indexed (IDY)

The operand is a zero-page address holding a 16-bit pointer. Y is added to
the pointer value to form EffAddr.

```
ptr    = Read(PC + 1)
base   = Read16_ZP(ptr)
EffAddr = base + Y
EffVal  = Read(EffAddr)
```

Crossing a page boundary costs one extra cycle on read instructions.

*Example:* `LDA ($10),Y`

## 4.14. Relative (REL)

Used only by branch instructions. The operand is a signed 8-bit offset
relative to the address of the byte immediately following the branch
instruction (PC + 2).

```
offset  = Read(PC + 1)   ; interpreted as signed
EffAddr = PC + 2 + offset
```

## 4.15. Two-byte Placeholder (BY2)

Consumes a 1-byte operand but performs no action. Used for undefined opcodes
that are nonetheless 2 bytes long.

## 4.16. Three-byte Placeholder (BY3)

Consumes a 2-byte operand but performs no action. Used for undefined opcodes
that are nonetheless 3 bytes long.

# 5. Instruction Set

Instructions are grouped below by category. All 256 possible opcodes are
mapped; any opcode not in the 65C02 instruction set is treated as a NOP of
the appropriate byte length (1, 2, or 3 bytes).

## 5.1. Load and Store

| Mnemonic | Description                        |
|----------|------------------------------------|
| LDA      | Load accumulator from memory       |
| LDX      | Load X from memory                 |
| LDY      | Load Y from memory                 |
| STA      | Store accumulator to memory        |
| STX      | Store X to memory                  |
| STY      | Store Y to memory                  |
| STZ      | Store zero to memory               |

LDA, LDX, and LDY set the Negative and Zero flags based on the loaded value.

## 5.2. Register Transfers

| Mnemonic | Description                        |
|----------|------------------------------------|
| TAX      | Transfer A to X                    |
| TAY      | Transfer A to Y                    |
| TXA      | Transfer X to A                    |
| TYA      | Transfer Y to A                    |
| TSX      | Transfer S (stack pointer) to X    |
| TXS      | Transfer X to S (stack pointer)    |

All transfers except TXS set Negative and Zero based on the transferred value.

## 5.3. Stack Operations

| Mnemonic | Description                        |
|----------|------------------------------------|
| PHA      | Push A onto stack                  |
| PHP      | Push P (status) onto stack         |
| PHX      | Push X onto stack                  |
| PHY      | Push Y onto stack                  |
| PLA      | Pull A from stack                  |
| PLP      | Pull P (status) from stack         |
| PLX      | Pull X from stack                  |
| PLY      | Pull Y from stack                  |

PLA, PLX, and PLY set Negative and Zero based on the pulled value. PLP
restores all flags from the stack verbatim.

## 5.4. Arithmetic

| Mnemonic | Description                                         |
|----------|-----------------------------------------------------|
| ADC      | Add with carry; sets N, V, Z, C                     |
| SBC      | Subtract with borrow; sets N, V, Z, C               |
| INC      | Increment memory by 1; sets N, Z                    |
| INX      | Increment X by 1; sets N, Z                         |
| INY      | Increment Y by 1; sets N, Z                         |
| DEC      | Decrement memory by 1; sets N, Z                    |
| DEX      | Decrement X by 1; sets N, Z                         |
| DEY      | Decrement Y by 1; sets N, Z                         |

ADC and SBC support both binary and BCD arithmetic (see section 8).

## 5.5. Compare

| Mnemonic | Description                                         |
|----------|-----------------------------------------------------|
| CMP      | Compare A with memory; sets N, Z, C                 |
| CPX      | Compare X with memory; sets N, Z, C                 |
| CPY      | Compare Y with memory; sets N, Z, C                 |

A compare subtracts the operand from the register but discards the result.
Carry is set if register >= operand (unsigned).

## 5.6. Bitwise and Shift

| Mnemonic | Description                                              |
|----------|----------------------------------------------------------|
| AND      | Bitwise AND with A; sets N, Z                            |
| ORA      | Bitwise OR with A; sets N, Z                             |
| EOR      | Bitwise XOR with A; sets N, Z                            |
| ASL      | Arithmetic shift left; MSB to Carry, 0 into LSB; N, Z, C|
| LSR      | Logical shift right; LSB to Carry, 0 into MSB; N, Z, C  |
| ROL      | Rotate left through Carry; sets N, Z, C                  |
| ROR      | Rotate right through Carry; sets N, Z, C                 |
| BIT      | Test bits: N = mem[7], V = mem[6], Z = (A & mem) == 0   |
| TSB      | Test and Set Bits: mem |= A; Z = (A & mem_before) == 0  |
| TRB      | Test and Reset Bits: mem &= ~A; Z = (A & mem) == 0      |

ASL, LSR, ROL, and ROR can operate on the accumulator (accumulator mode) or
on a memory location.

BIT with an immediate operand only sets Z; it does not affect N or V.

## 5.7. Branch

All branch instructions use relative addressing. A branch is taken (PC is
updated to EffAddr) if and only if the tested condition is true. A branch
that is not taken simply advances to the next instruction.

| Mnemonic | Condition           |
|----------|---------------------|
| BCC      | Carry clear         |
| BCS      | Carry set           |
| BEQ      | Zero set            |
| BNE      | Zero clear          |
| BMI      | Negative set        |
| BPL      | Negative clear      |
| BVC      | Overflow clear      |
| BVS      | Overflow set        |
| BRA      | Always (unconditional branch) |

## 5.8. Jump and Call

| Mnemonic | Description                                              |
|----------|----------------------------------------------------------|
| JMP      | Set PC to EffAddr (absolute or indirect)                 |
| JSR      | Push return address (PC + 2) onto stack, then jump       |
| RTS      | Pull return address from stack, add 1, set PC            |

JSR pushes the high byte of the return address first, then the low byte.
RTS pops in the reverse order and adds 1 to reconstruct the correct PC.

## 5.9. Interrupt and Break

| Mnemonic | Description                                              |
|----------|----------------------------------------------------------|
| BRK      | Software interrupt: push PC+2 and P, set I flag, clear D|
| RTI      | Return from interrupt: restore P then PC from stack      |

BRK pushes the high byte of PC+2, then the low byte, then P, and sets the
Interrupt flag. It does not set the Break flag in P before pushing; the Break
flag in the pushed value reflects the state of P at the time of the BRK.

RTI pops P first (unconditionally setting Unused and Break), then pops the
low byte of PC, then the high byte. Unlike RTS, RTI does not add 1 to the
popped address.

## 5.10. Status Flag Operations

| Mnemonic | Description                 |
|----------|-----------------------------|
| CLC      | Clear Carry flag            |
| SEC      | Set Carry flag              |
| CLD      | Clear Decimal flag          |
| SED      | Set Decimal flag            |
| CLI      | Clear Interrupt Disable flag|
| SEI      | Set Interrupt Disable flag  |
| CLV      | Clear Overflow flag         |

## 5.11. No Operation

| Mnemonic | Description                                        |
|----------|----------------------------------------------------|
| NOP      | No operation; 1 byte, 2 cycles                     |
| NP2      | Placeholder for undefined 2-byte opcodes           |
| NP3      | Placeholder for undefined 3-byte opcodes           |

NP2 and NP3 are internal names for undefined opcodes; they are not valid
assembly mnemonics.

# 6. Execution Model

## 6.1. Instruction Cycle

Each call to Execute performs one complete instruction:

```
Execute():
    LastPC  = PC
    opcode  = Read(PC)
    addrFn  = addrModeTable[opcode]
    instrFn = instructionTable[opcode]

    addrFn()    ; compute EffAddr and EffVal
    instrFn()   ; perform the operation
    PC += byteLength[opcode]

    P |= Unused | Break   ; convention: both flags always set after execution
    cycleCounter += OpcodeCycles()
```

Branch and jump instructions set PC directly during `instrFn()`. The
`byteLength[opcode]` increment at the end is 0 for those instructions
(their own logic has already set PC to the correct target).

## 6.2. Flag Maintenance

The Unused flag (bit 5 of P) is always set after every instruction. The Break
flag (bit 4 of P) is also always set after every instruction. These two bits
are set unconditionally in P at the end of Execute, not by individual
instructions.

# 7. Stack

The stack occupies memory addresses $0100-$01FF (page 1). S holds the low
byte of the current stack pointer. The full address for stack access is always
$0100 + S.

```
Push(val):
    Write($0100 + S, val)
    S -= 1

Pop():
    S += 1
    return Read($0100 + S)
```

The stack has no overflow or underflow detection; S wraps if it exceeds the
range $00-$FF. JSR and BRK push two and three bytes respectively; RTS and
RTI pop the same quantities.

# 8. Decimal (BCD) Mode

When the Decimal flag is set, ADC and SBC operate on Binary-Coded Decimal
values rather than plain binary. In BCD each byte encodes two decimal digits
as two 4-bit nibbles in the range 0-9 each (e.g., $42 represents the decimal
value 42).

```
BCD_Add(a, b, carry):
    result = decimal(a) + decimal(b) + carry
    carry  = result > 99 ? 1 : 0
    return (binary(result % 100), carry)

BCD_Subtract(a, b, carry):
    result = decimal(a) - decimal(b) - (1 - carry)
    borrow = result < 0 ? 1 : 0
    return (binary((result + 100) % 100), 1 - borrow)
```

The Zero and Negative flags are set based on the binary result. The Overflow
flag behavior in decimal mode is undefined on the real hardware; the emulator
sets it based on the binary operation.

Decimal mode adds 1 extra cycle to every instruction when the Decimal flag
is set.

SED enables decimal mode; CLD restores binary mode.

# 9. Cycle Counting

Every instruction has a base cycle count derived from a static table. Several
conditions add extra cycles at runtime:

## 9.1. Page-Boundary Crossing

For read instructions using ABX, ABY, or IDY addressing, if the addition of
the index register causes the high byte of the address to change (i.e., a
page boundary is crossed), one extra cycle is charged.

```
if (base_addr & $FF00) != (EffAddr & $FF00):
    cycles += 1
```

## 9.2. Branch Instructions

A branch instruction that is not taken costs its base cycle count (2 cycles).
A branch that is taken adds 1 cycle. If taking the branch also crosses a page
boundary, an additional cycle is added.

```
if branch_taken:
    cycles += 1
    if (PC & $FF00) != (EffAddr & $FF00):
        cycles += 1
```

## 9.3. Decimal Mode

When the Decimal flag is set, 1 extra cycle is added to every instruction.

## 9.4. Cycle Accumulation

The CPU maintains a running total of all cycles consumed since reset. This
counter is monotonically increasing and is the primary mechanism the emulator
uses to synchronize the CPU with other chips (e.g., timing for video and
audio).

# 10. State Serialization

The CPU supports saving and restoring its complete state. A snapshot captures:
PC, LastPC, A, X, Y, P, S, the current cycle counter, and any internal
per-instruction state needed to resume correctly. Restoring a snapshot returns
the CPU to exactly the state it was in when the snapshot was taken.

# 11. Exported Interface

The package exposes the following to callers:

**Types**
- `CPU` -- the processor; callers create one and attach memory read/write
  functions before calling Execute.

**CPU construction fields** (set before first Execute call):
- `RMem` -- the memory read function
- `WMem` -- the memory write function

**Execution**
- `Execute() error` -- run one instruction; returns an error only on an
  unrecoverable internal fault.

**Registers** (direct field access)
- `PC`, `A`, `X`, `Y`, `P`, `S`

**Introspection**
- `CycleCounter() uint64` -- total cycles elapsed
- `Opcode() uint8` -- opcode of the instruction most recently fetched
- `OpcodeInstructionName(opcode) string` -- mnemonic string for an opcode
- `OpcodeAddrMode(opcode) int` -- addressing mode constant for an opcode
- `OperandSize(opcode) uint16` -- number of operand bytes for an opcode
- `AddrModeName(mode) string` -- human-readable name for a mode constant

**Debugging**
- `Status() string` -- one-line formatted CPU state
- `CurrentInstructionShort() string` -- disassembly of instruction at PC
- `LastInstructionLine(cycles) *Instruction` -- structured record of last
  executed instruction
- `Speculate(addr)` -- speculatively execute from an address without
  committing changes, for look-ahead debugging

**State save/restore**
- `Snapshot() *CPUState` -- capture current state
- `Restore(state *CPUState)` -- restore a previously captured state
