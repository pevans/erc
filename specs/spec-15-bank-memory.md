---
Specification: 15
Category: Memory
Drafted At: 2026-03-26
Authors:
  - Peter Evans
---

# 1. Overview

This spec describes the bank-switched memory system of the Apple IIe. Bank
memory allows the CPU to access RAM in the address range $D000-$FFFF, which
is normally occupied by system ROM. It also provides two independent 4 KB
banks of RAM in the $D000-$DFFF range, giving the computer more usable RAM
than its 16-bit address space would otherwise allow.

This spec does not cover auxiliary memory (the 80-column card's 64 KB
expansion), which is a separate switching mechanism.

# 2. Memory Map

The bank-switched region spans addresses $D000 through $FFFF (12 KB total).
It is divided into two subregions:

- **$D000-$DFFF** (4 KB): the "DF block." This range has two independent RAM
  banks (bank 1 and bank 2), only one of which is accessible at a time.
- **$E000-$FFFF** (8 KB): a single RAM region shared by both bank
  configurations.

When bank-switching is set to read ROM, the system ROM at the equivalent
offset is returned instead of RAM. System ROM is mapped starting at $C000, so
a read of address $D000 in ROM mode returns the byte at ROM offset $1000
($D000 - $C000).

# 3. State

The bank-switching system maintains four boolean flags:

- **ReadRAM**: when true, reads in $D000-$FFFF return RAM. When false, reads
  return system ROM.
- **WriteRAM**: when true, writes in $D000-$FFFF store to RAM. When false,
  writes are silently discarded.
- **DFBlockBank2**: when true, the $D000-$DFFF range maps to bank 2. When
  false, it maps to bank 1. This flag has no effect on $E000-$FFFF.
- **ReadAttempts**: a counter used to implement the double-read protection for
  write-enable switches (see section 5).

## 3.1. Default State

On power-up or reset, the bank switches are initialized to:

    ReadRAM       = false   (read from ROM)
    WriteRAM      = true    (writes go to RAM)
    DFBlockBank2  = true    (use bank 2)
    ReadAttempts  = 0

This means the computer boots reading ROM and writing RAM in bank 2. Bank 2 is
the original language card bank and the default on the Apple IIe. This is the
expected state for the system monitor and Applesoft BASIC ROM to be accessible
on startup.

# 4. RAM Layout

Main memory is 0x11000 bytes (69,632 bytes). The extra 0x1000 bytes beyond
the normal 64 KB address space accommodate bank 2:

- **Bank 1** data for $D000-$DFFF lives at its natural offset in the segment
  (addresses $D000-$DFFF).
- **Bank 2** data for $D000-$DFFF lives at an offset of address + $3000. A
  read of $D000 in bank 2 mode reads from segment offset $10000. A read of
  $DFFF in bank 2 mode reads from segment offset $10FFF.
- **$E000-$FFFF** data always lives at its natural offset in the segment,
  regardless of bank selection.

# 5. Soft Switches

The bank-switching soft switches occupy addresses $C080-$C08F. All 16 are
read-triggered -- the CPU reads from the address to activate the switch. The
side effect is the state change. The value returned by the read is $00. On
real hardware, these switches return the floating bus value, but erc does not
currently emulate the floating bus.

The 16 switches are organized as two groups of 8, with identical behavior
except for the bank selection:

- **$C080-$C087**: select bank 2 (DFBlockBank2 = true)
- **$C088-$C08F**: select bank 1 (DFBlockBank2 = false)

Within each group, addresses $C084-$C087 and $C08C-$C08F are duplicates of
$C080-$C083 and $C088-$C08B respectively. This duplication exists in the
hardware and must be emulated.

Writes to $C080-$C08F have no effect. Only reads trigger state changes.

## 5.1. Switch Table

    Address   Read RAM   Write RAM   Bank    Notes
    -------   --------   ---------   ----    -----
    $C080     yes        no          2
    $C081     no         yes*        2       write requires double read
    $C082     no         no          2
    $C083     yes        yes*        2       write requires double read
    $C084     yes        no          2       duplicate of $C080
    $C085     no         yes*        2       duplicate of $C081
    $C086     no         no          2       duplicate of $C082
    $C087     yes        yes*        2       duplicate of $C083
    $C088     yes        no          1
    $C089     no         yes*        1       write requires double read
    $C08A     no         no          1
    $C08B     yes        yes*        1       write requires double read
    $C08C     yes        no          1       duplicate of $C088
    $C08D     no         yes*        1       duplicate of $C089
    $C08E     no         no          1       duplicate of $C08A
    $C08F     yes        yes*        1       duplicate of $C08B

Entries marked with * require the double-read mechanism described in section
5.3 to enable WriteRAM.

## 5.2. Non-Write-Enable Switches

Switches that do not enable WriteRAM ($C080, $C082, $C084, $C086, $C088,
$C08A, $C08C, $C08E) take effect immediately on a single read. ReadRAM,
WriteRAM, and DFBlockBank2 are all set according to the switch table. In
particular, these switches set WriteRAM to false unconditionally.

Reading a non-write-enable switch also resets the ReadAttempts counter to zero
(see section 5.3).

## 5.3. Double-Read Write Protection

Switches that enable WriteRAM ($C081, $C083, $C085, $C087, $C089, $C08B,
$C08D, $C08F) do not enable it on the first access. The CPU must read a
write-enable switch twice in a row for WriteRAM to become true. The two reads
do not need to be the same switch -- any two consecutive reads of write-enable
switches will satisfy the requirement. Specifically:

1. On the first read, ReadRAM and DFBlockBank2 are set according to the
   switch, but WriteRAM remains unchanged.
2. On the second consecutive read of a write-enable switch, if the access
   is an instruction-initiated read operation, WriteRAM is set to true.

The ReadAttempts counter tracks how many consecutive instruction-initiated
reads have been made to write-enable switches. After each instruction, if the
CPU's effective address was a write-enable switch and the access was an
instruction-initiated read, ReadAttempts is incremented. Otherwise,
ReadAttempts resets to zero. WriteRAM is enabled when the counter is at least
1 at the time the switch is read.

This mechanism prevents accidental write-enable by stray reads. A single
`LDA $C083` will not enable writing; the program must execute two consecutive
reads of a write-enable switch.

## 5.4. Instruction Read Requirement

The double-read mechanism only counts reads that are part of an instruction's
operand fetch -- that is, reads initiated by the CPU as part of executing an
instruction (such as `LDA $C083`). Reads caused by other mechanisms (such as
DMA or incidental bus activity) do not count toward the double-read
requirement.

# 6. Status Switches

Three read-only soft switches report the current state of bank switching.
Each returns bit 7 set ($80) if the condition is true, or $00 if false.

    Address   Name      Returns $80 when
    -------   ----      ----------------
    $C011     RDBNK2    DFBlockBank2 is true (bank 2 selected for $D000-$DFFF)
    $C012     RDLCRAM   ReadRAM is true ($D000-$FFFF reads from RAM)
    $C016     RDALTZP   zero page is using auxiliary memory

Note: $C016 (RDALTZP) reports the state of the zero-page/stack-page memory
source selection, which is described in section 7. It is listed here because
it is implemented alongside the other bank status switches.

# 7. Zero Page and Stack Page Switching

Two soft switches control which memory segment is used for page zero
($0000-$00FF) and page one ($0100-$01FF, the hardware stack):

    Address   Name       Effect
    -------   ----       ------
    $C008     SETSTDZP   Use main memory for pages 0 and 1
    $C009     SETALTZP   Use auxiliary memory for pages 0 and 1

These switches are write-triggered -- the CPU writes to the address to activate
the switch. Reads of these addresses have no effect and return $00.

When SETALTZP is active, reads and writes to $0000-$01FF are directed to the
auxiliary memory segment instead of main memory. This also affects which
segment is used for bank-switched reads and writes in $D000-$FFFF -- the
bank-switching logic operates on whichever segment is currently selected by
the zero-page switch.

The default state is SETSTDZP (main memory).

# 8. Read Logic

When the CPU reads an address in the range $D000-$FFFF:

    if ReadRAM is false:
        return ROM[address - $C000]

    if DFBlockBank2 is true AND address < $E000:
        return segment[address + $3000]

    return segment[address]

Where `segment` is either the main or auxiliary memory segment, as determined
by the zero-page switch state.

# 9. Write Logic

When the CPU writes to an address in the range $D000-$FFFF:

    if WriteRAM is false:
        discard the write (no-op)
        return

    if DFBlockBank2 is true AND address < $E000:
        segment[address + $3000] = value
        return

    segment[address] = value

Where `segment` is either the main or auxiliary memory segment, as determined
by the zero-page switch state.

# 10. Mapping Range

The bank-switched read and write logic is mapped to the address range
$D000-$FFFF (exclusive of $10000). The zero-page switching logic is mapped to
$0000-$01FF.

All other addresses in the 64 KB space are unaffected by bank switching.
