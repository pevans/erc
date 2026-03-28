---
Specification: 16
Drafted At: 2026-03-27
Authors:
  - Peter Evans
---

# 1. Overview

This spec describes the auxiliary memory system of the Apple II, provided by
the extended 80-column card. Auxiliary memory is a second 64 KB RAM bank that
exists alongside main memory. Software can independently redirect reads and
writes to auxiliary memory using soft switches, and display modes can
selectively route specific address ranges to auxiliary memory for 80-column
text and double-hires graphics. The rendering of those display modes is
outside the scope of this spec; this spec covers only the memory routing.

The extended 80-column card provides 68 KB of RAM in total: 64 KB of
auxiliary memory (which includes its own bank 1 of $D000-$DFFF at the
natural offset) plus an additional 4 KB that serves as bank 2 of
$D000-$DFFF within the auxiliary segment. The bank-switching mechanism
itself is described in spec 15. While the RAM physically lives on the same
card, the two switching systems are controlled independently: bank switching
determines whether the high address range reads from RAM or ROM (and which
bank), while the auxiliary memory switches described here control which
64 KB segment is the source or destination of a given access.

# 2. Memory Segments

The Apple II has two primary RAM segments, each 0x11000 bytes (69,632
bytes). The extra 0x1000 bytes beyond 64 KB accommodate the second bank of
$D000-$DFFF for the bank-switching system (see spec 15).

- **Main memory**: the default segment for all reads and writes.
- **Auxiliary memory**: the expansion segment provided by the 80-column card.

Both segments have identical structure. Which segment is active for a given
access depends on the state of several soft switches described below.

# 3. State

The auxiliary memory system maintains two boolean flags:

- **ReadAux**: when true, general reads come from auxiliary memory. When
  false, reads come from main memory.
- **WriteAux**: when true, general writes go to auxiliary memory. When false,
  writes go to main memory.

Two segment pointers track the currently active segments:

- **ReadSegment**: points to whichever segment (main or aux) should be used
  for general reads.
- **WriteSegment**: points to whichever segment (main or aux) should be used
  for general writes.

## 3.1. Default State

On power-up or reset, all auxiliary memory and display override flags are
initialized to:

    ReadAux      = false   (read from main memory)
    WriteAux     = false   (write to main memory)
    ReadSegment  = main
    WriteSegment = main
    80STORE      = false
    PAGE2        = false
    HIRES        = false

# 4. Soft Switches

Four write-triggered soft switches control whether reads and writes are
directed to main or auxiliary memory. The CPU writes to the address to
activate the switch. The value written is ignored. Reads of these addresses
have no effect and return $00.

    Address   Name           Effect
    -------   ----           ------
    $C002     RAMRD off      ReadAux = false, ReadSegment = main
    $C003     RAMRD on       ReadAux = true, ReadSegment = aux
    $C004     RAMWRT off     WriteAux = false, WriteSegment = main
    $C005     RAMWRT on      WriteAux = true, WriteSegment = aux

## 4.1. Independence

RAMRD and RAMWRT are independent of each other. A program can read from
auxiliary memory while writing to main memory, or vice versa. This is useful
for copying data between the two segments.

# 5. Status Switches

Two read-only soft switches report the current state of auxiliary memory
selection. Each returns bit 7 set ($80) if the condition is true, or $00 if
false.

    Address   Name         Returns $80 when
    -------   ----         ----------------
    $C013     RDRAMRD      ReadAux is true (reads come from aux)
    $C014     RDRAMWRT     WriteAux is true (writes go to aux)

Writes to these addresses are ignored and have no effect. The same applies
to the display status addresses $C018 (RD80STORE), $C01C (RDPAGE2), and
$C01D (RDHIRES).

# 6. Interaction with Display Modes (80STORE)

When the 80STORE display switch is active, it overrides the RAMRD/RAMWRT
switches for specific address ranges used by display memory. The 80STORE
system is controlled by the display soft switches described below.

## 6.1. 80STORE Soft Switch

    Address   Name          Effect
    -------   ----          ------
    $C000     80STORE off   Disable 80STORE (write-triggered)
    $C001     80STORE on    Enable 80STORE (write-triggered)
    $C018     RD80STORE     Returns $80 if 80STORE is on (read-only)

Note: $C000 is also the keyboard data register on reads, and $C001 is
the keyboard strobe clear (KBDSTROBE) on reads. Both 80STORE switches only
respond to writes; reads of $C000 and $C001 retain their keyboard
functions.

When 80STORE is off, the RAMRD and RAMWRT switches fully control segment
selection for all addresses in their scope.

## 6.2. PAGE2 Soft Switch

    Address   Name          Effect
    -------   ----          ------
    $C054     PAGE2 off     Select page 1 (read/write-triggered, reads return $00)
    $C055     PAGE2 on      Select page 2 (read/write-triggered, reads return $00)
    $C01C     RDPAGE2       Returns $80 if PAGE2 is on (read-only)

When 80STORE is off, PAGE2 selects which display page the hardware renders
from within main memory:

- PAGE2 off: text page 1 ($0400-$07FF), hires page 1 ($2000-$3FFF).
- PAGE2 on: text page 2 ($0800-$0BFF), hires page 2 ($4000-$5FFF).

When 80STORE is on, PAGE2 is repurposed: instead of selecting between page 1
and page 2 in main memory, it determines whether display memory accesses in
the page 1 address range go to main or auxiliary memory. In this mode, the
page 2 address ranges are not affected by 80STORE and follow RAMRD/RAMWRT
normally.

## 6.3. 80STORE Override Rules

When 80STORE is enabled, the following address ranges are affected:

**Text page ($0400-$07FF)**: both reads and writes are routed based on the
PAGE2 flag, regardless of the RAMRD/RAMWRT settings.

- PAGE2 off: text page accesses go to main memory.
- PAGE2 on: text page accesses go to auxiliary memory.

**Hires page 1 ($2000-$3FFF)**: both reads and writes are routed based on
the PAGE2 flag, but only when the HIRES display mode is also enabled.

- HIRES off: no override; RAMRD/RAMWRT apply normally.
- HIRES on, PAGE2 off: hires accesses go to main memory.
- HIRES on, PAGE2 on: hires accesses go to auxiliary memory.

The HIRES soft switch is:

    Address   Name          Effect
    -------   ----          ------
    $C056     HIRES off     Disable hires mode (read/write-triggered, reads return $00)
    $C057     HIRES on      Enable hires mode (read/write-triggered, reads return $00)
    $C01D     RDHIRES       Returns $80 if HIRES is on (read-only)

All other addresses are unaffected by 80STORE and continue to follow the
RAMRD/RAMWRT switches.

# 7. Interaction with Zero Page Switching

The ALTZP switch (described in spec 15, section 7) controls which segment is
used for zero page ($0000-$00FF) and the stack page ($0100-$01FF). When
ALTZP is active, these two pages are read from and written to auxiliary memory
regardless of the RAMRD/RAMWRT settings.

ALTZP also determines which segment is used for bank-switched memory in
$D000-$FFFF. See spec 15 for details.

# 8. Read Logic

When the CPU reads an address in the general address space (outside of
bank-switched and zero-page/stack ranges):

    if 80STORE is on:
        if address is in $0400-$07FF:
            if PAGE2 is on:
                return aux[address]
            return main[address]

        if address is in $2000-$3FFF AND HIRES is on:
            if PAGE2 is on:
                return aux[address]
            return main[address]

    return ReadSegment[address]

Where ReadSegment is main or auxiliary memory, as determined by the RAMRD
switch.

# 9. Write Logic

When the CPU writes to an address in the general address space (outside of
bank-switched and zero-page/stack ranges):

    if 80STORE is on:
        if address is in $0400-$07FF:
            if PAGE2 is on:
                aux[address] = value
                return
            main[address] = value
            return

        if address is in $2000-$3FFF AND HIRES is on:
            if PAGE2 is on:
                aux[address] = value
                return
            main[address] = value
            return

    WriteSegment[address] = value

Where WriteSegment is main or auxiliary memory, as determined by the RAMWRT
switch.

# 10. Mapping Ranges

The soft switch addresses for the auxiliary memory system are:

- $C000-$C001: 80STORE control (write-triggered)
- $C002-$C005: RAMRD/RAMWRT control (write-triggered)
- $C013-$C014: RAMRD/RAMWRT status (read-only)
- $C018: 80STORE status (read-only)

The 80STORE override applies to memory routing in these address ranges:

- $0400-$07FF: text display page 1 (always affected when 80STORE is on)
- $2000-$3FFF: hires display page 1 (affected when both 80STORE and HIRES
  are on)

The page 2 display ranges ($0800-$0BFF for text, $4000-$5FFF for hires)
are not affected by 80STORE. They always follow RAMRD/RAMWRT for memory
routing. PAGE2 selects which page the display hardware renders from; that
rendering behavior is described in the display mode specs (specs 9, 11,
and 12).

# 11. Summary of Address Space Routing

For any given address, the segment used for access depends on multiple
switches. The following table summarizes the resolution order:

    Address Range    Condition                       Segment
    -------------    ---------                       -------
    $0000-$01FF      ALTZP on                        aux
    $0000-$01FF      ALTZP off                       main
    $0200-$03FF      (always)                        RAMRD/RAMWRT
    $0400-$07FF      80STORE on, PAGE2 on            aux
    $0400-$07FF      80STORE on, PAGE2 off           main
    $0400-$07FF      80STORE off                     RAMRD/RAMWRT
    $0800-$1FFF      (always)                        RAMRD/RAMWRT
    $2000-$3FFF      80STORE on, HIRES on, PAGE2 on  aux
    $2000-$3FFF      80STORE on, HIRES on, PAGE2 off main
    $2000-$3FFF      80STORE on, HIRES off           RAMRD/RAMWRT
    $2000-$3FFF      80STORE off                     RAMRD/RAMWRT
    $4000-$BFFF      (always)                        RAMRD/RAMWRT
    $C000-$C0FF      (I/O space, soft switches)      n/a
    $C100-$CFFF      (peripheral ROM, separate spec)  n/a
    $D000-$FFFF      (bank-switched, see spec 15)    ALTZP segment
