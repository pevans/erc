---
Specification: 17
Drafted At: 2026-03-28
Authors:
  - Peter Evans
---

# 1. Overview

This spec describes the peripheral support system of the Apple II enhanced.
Peripheral support encompasses two concerns: the I/O soft switches through
which the CPU communicates with peripheral hardware, and the ROM visibility
system that controls whether the CPU sees internal ROM or peripheral card ROM
in the $C100-$CFFF address range.

The Apple II has 7 numbered expansion slots (1-7) plus a special "slot 0"
reserved for internal hardware. Each slot is assigned a 256-byte ROM region in
$C100-$C7FF (one page per slot) and can share a 2 KB expansion ROM region at
$C800-$CFFF. Each slot is also assigned 16 I/O addresses in the $C080-$C0FF
range at $C080 + (slot * $10).

This spec describes the ROM visibility switching system in detail and
enumerates the peripheral subsystems that exist on the Apple II. The
detailed behavior of each peripheral subsystem is defined in its own spec.

# 2. Peripheral ROM Visibility

The address range $C100-$CFFF can show either internal ROM or peripheral card
ROM. Two independent switches control which ROM is visible, plus a mechanism
for cards to enable expansion ROM in the $C800-$CFFF range.

## 2.1. ROM Segments

The system maintains a single ROM segment that contains both internal and
peripheral ROM data:

- **Internal ROM**: the first 0x4000 bytes of the ROM segment. Internal ROM
  for the $C100-$CFFF range is at offset (address - $C000) within this region.
- **Peripheral ROM**: starts at offset $4000 in the ROM segment. Peripheral
  ROM for address $Cxxx is at offset (address - $C000 + $4000).

# 3. ROM Switching State

The ROM visibility system maintains these state flags:

- **SlotCX**: when true, reads of $C100-$C7FF (excluding $C300 in some cases)
  return peripheral ROM. When false, they return internal ROM. Default: true.
- **SlotC3**: when true, reads of $C300-$C3FF return peripheral ROM. When
  false, they return internal ROM. This flag is independent of SlotCX and only
  affects the $C300-$C3FF range. Default: false.
- **Expansion**: when true, reads of $C800-$CFFF return expansion ROM. When
  false, they return whichever ROM (internal or peripheral) the other flags
  select. Default: false.
- **IOSelect**: set to true when the CPU reads any address in $C100-$C7FF
  while SlotCX is true. Cleared when expansion is disabled. Default: false.
- **IOStrobe**: set to true when the CPU reads any address in $C800-$CFFE
  while SlotCX is true. A read of $CFFF does not set IOStrobe; instead it
  clears all expansion state. Cleared when expansion is disabled. Default:
  false.
- **ExpSlot**: tracks which slot's expansion ROM should be active, derived
  from the most recent read in the $C100-$C7FF range while SlotCX is true.
  Cleared when expansion is disabled. Default: 0.

# 4. ROM Switching Soft Switches

## 4.1. Control Switches

    Address   Trigger   Name           Effect
    -------   -------   ----           ------
    $C006     Write     SETSLOTCXROM   SlotCX = true (show peripheral ROM)
    $C007     Write     SETINTCXROM    SlotCX = false (show internal ROM)
    $C00A     Write     SETINTC3ROM    SlotC3 = false (show internal ROM for $C300)
    $C00B     Write     SETSLOTC3ROM   SlotC3 = true (show peripheral ROM for $C300)

Reads of these addresses have no effect on peripheral state. The return
value is open bus (typically the keyboard data latch).

## 4.2. Status Switches

    Address   Trigger   Name         Bit 7
    -------   -------   ----         -----
    $C015     Read      RDCXROM      0 if SlotCX is true; 1 if false
    $C017     Read      RDC3ROM      1 if SlotC3 is true; 0 if false

The status is reported in bit 7 of the return value. The remaining bits
reflect open bus (typically the keyboard data latch).

Writes to these addresses have no effect.

## 4.3. Expansion ROM Disable

    Address   Trigger   Name         Effect
    -------   -------   ----         ------
    $CFFF     Read      EXPROMOFF    Clears IOSelect, IOStrobe, Expansion, and ExpSlot

Any read of $CFFF disables expansion ROM and clears all expansion-related
state. This is a side-effect read: it returns the ROM byte at address $CFFF
(from whichever ROM was active before the side effect takes place), and then
clears the expansion flags. A write to $CFFF also disables expansion ROM.

# 5. ROM Read Logic

When the CPU reads an address in the range $C100-$CFFF, the following logic
determines the return value:

    if address == $CFFF:
        if SlotCX is true:
            if IOSelect is true:
                value = expansion ROM at $CFFF
            else:
                value = peripheral ROM at $CFFF
        else if Expansion is true:
            value = expansion ROM at $CFFF
        else:
            value = internal ROM at $CFFF
        disable expansion (clear IOSelect, IOStrobe, Expansion, ExpSlot)
        return value

    if SlotCX is true:
        if address is in $C100-$C7FF:
            set IOSelect = true
            set ExpSlot = slot number from address (bits 8-11)

            if SlotC3 is false AND address is in $C300-$C3FF:
                return internal ROM

            return peripheral ROM

        if address is in $C800-$CFFF:
            set IOStrobe = true
            if IOSelect is true:
                set Expansion = true
                return expansion ROM
            return peripheral ROM

    if SlotC3 is true AND address is in $C300-$C3FF:
        return peripheral ROM

    if Expansion is true AND address is in $C800-$CFFF:
        return expansion ROM

    return internal ROM

## 5.1. Expansion ROM Content

Erc does not currently support peripheral cards that supply their own ROM.
When expansion ROM is read, the system falls back to returning internal ROM
data. This is correct behavior for a system with no physical expansion cards
installed; if peripheral card ROM support is added in the future, this
fallback would be replaced with ROM data from the card in the active
expansion slot.

# 6. ROM Write Logic

Writes to the $C100-$CFFF address range are ignored. ROM is read-only. The
only side effect of a write is that a write to $CFFF disables expansion ROM,
the same as a read.

# 7. Scope

This spec covers only slot-based peripherals -- cards that occupy one of the
Apple II's 7 expansion slots and are accessed through the slot's I/O range
($C080 + slot * $10) and ROM space ($C100 + slot * $100). Other subsystems
that use soft switches in the $C0 page but do not occupy a slot (keyboard,
speaker, display, bank switching, auxiliary memory) are described in their
own specs.

The only slot-based peripheral currently emulated is the Disk II controller
in slot 6. Full details are in spec 13.

# 8. Subsystems Not Currently Emulated

The following peripherals existed for the Apple II but are not currently
emulated. They are listed here for completeness and as guidance for future
work.

## 8.1. Game I/O (Joystick / Paddles)

The Apple II provides game I/O through addresses $C061-$C067 (button
inputs) and $C070 (paddle trigger). The annunciator outputs are at
$C058-$C05F. Game I/O supports:

- 3 pushbutton inputs (buttons 0-2) at $C061-$C063 (bit 7 high when pressed)
- 4 analog paddle inputs read via a timing loop triggered by $C070
- 4 annunciator outputs toggled by $C058-$C05F

## 8.2. Printer Card

A common slot 1 card. The printer interface typically uses one I/O address at
$C090 (slot 1 base) for data output, plus a small ROM at $C100-$C1FF for the
driver.

## 8.3. Serial / Super Serial Card

Typically installed in slot 2. Provides RS-232 serial communication using
ACIA (6551) registers mapped to the slot's I/O range. Used for modems,
printers, and other serial devices.

## 8.4. 80-Column Firmware

Slot 3 is reserved for the built-in 80-column firmware on the Apple II.
The SlotC3 switch controls whether the internal 80-column ROM or a physical
card's ROM appears at $C300-$C3FF. The memory provided by the extended
80-column card is emulated (see spec 16), but no additional slot 3 card
firmware beyond the internal ROM is supported.

## 8.5. Mouse Card

Typically installed in slot 4. Provides mouse position and button state
through the slot's I/O addresses and a firmware ROM that includes interrupt
handling routines.

## 8.6. Clock Card

Various real-time clock cards existed for the Apple II, commonly in slot 5
or slot 7. They provide date and time data through the slot's I/O addresses.

## 8.7. Slot 7

Slot 7 is often used for a RAM disk (such as the RAMWorks card) or
additional drive controllers. Its I/O range is $C0F0-$C0FF.
