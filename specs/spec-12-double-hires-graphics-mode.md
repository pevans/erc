---
Specification: 12
Category: Graphics
Drafted At: 2026-03-25
Authors:
  - Peter Evans
---

# 1. Overview

This spec describes the implementation of double high-resolution (double
hires) graphics mode for the Apple //e. Double hires mode uses both main and
auxiliary memory to produce a 560-dot-wide by 192-row display. In color mode,
dots are grouped into 4-dot windows that map to 16 NTSC artifact colors,
yielding an effective resolution of 140x192 with 16 colors. In monochrome
mode, every dot is individually visible, giving a full 560x192 monochrome
display.

Double hires requires the DHIRES and 80COL soft switches to be active in
addition to HIRES, and uses twice as much display memory as standard hires by
drawing from both main and auxiliary memory simultaneously.

# 2. Screen Geometry

## 2.1. Logical Dimensions

The double hires display is 560 dots wide by 192 rows tall. Each row is stored
across 80 byte columns: 40 bytes from auxiliary memory interleaved with 40
bytes from main memory. Each byte contributes 7 dots (bits 0-6); bit 7 is
unused. This gives 80 x 7 = 560 dots per row.

In color mode, the 560 dots are interpreted through a sliding 4-dot window,
producing an effective color resolution of 140 columns by 192 rows with 16
colors.

## 2.2. Dot Dimensions

Each logical dot is rendered as a 1x2 pixel block on the 560x384 framebuffer:

- Width: 560 dots x 1 pixel = 560 pixels
- Height: 192 rows x 2 pixels = 384 pixels

Unlike standard hires (which renders each dot as 2x2 pixels), double hires
dots are single-pixel wide because the 560-dot row already fills the full
framebuffer width.

# 3. Display Memory

## 3.1. Memory Regions

Double hires uses both main and auxiliary memory in the hires address range.
Page 1 uses $2000-$3FFF in both main and aux (16 KB total). Page 2 uses
$4000-$5FFF in both main and aux. Only 7,680 bytes per segment are visible
(40 bytes x 192 rows); the remaining bytes in each 8 KB region are hole bytes
that exist in memory but do not appear on screen.

## 3.2. Address Interleaving

Double hires uses the same row address interleaving as standard hires. The 192
rows are organized into three groups of 64 rows, with the same nested
subdivision formula:

```
group    = r / 64          (0, 1, or 2)
section  = (r % 64) / 8    (0-7)
line     = r % 8           (0-7)

base = $2000 + (line * $0400) + (section * $0080) + (group * $0028)
```

For page 2, substitute $4000 for $2000. The address table is identical to the
one in spec 11, section 3.2.1.

## 3.3. Byte Interleaving Within a Row

For each of the 40 memory offsets (0-39) within a row, two bytes contribute
dots to the display: one from auxiliary memory and one from main memory. The
byte order alternates between aux and main across 80 screen columns:

```
Screen column 0:  aux byte at offset 0
Screen column 1:  main byte at offset 0
Screen column 2:  aux byte at offset 1
Screen column 3:  main byte at offset 1
...
Screen column 78: aux byte at offset 39
Screen column 79: main byte at offset 39
```

Each screen column produces 7 dots, giving 80 x 7 = 560 dots per row.

## 3.4. Byte Layout

Each byte contributes 7 dots from bits 0-6. Bit 7 is unused (unlike standard
hires, where bit 7 selects a color palette):

```
bit 7    bit 6   bit 5   bit 4   bit 3   bit 2   bit 1   bit 0
unused   dot 6   dot 5   dot 4   dot 3   dot 2   dot 1   dot 0
```

- **Bits 0-6**: each bit controls one dot. 1 = on, 0 = off. Bit 0 is the
  leftmost dot within the byte; bit 6 is the rightmost.
- **Bit 7**: ignored. Programs may set it to any value without affecting the
  display.

# 4. Color Generation

## 4.1. The 16 Double Hires Colors

Double hires produces 16 colors through NTSC artifact coloring. Each color
corresponds to a 4-bit pattern derived from a sliding window across the dot
stream:

| Pattern | Color        | RGB                     |
|---------|--------------|-------------------------|
| 0000    | Black        | (0, 0, 0)               |
| 0001    | Magenta      | (144, 23, 64)           |
| 0010    | Brown        | (64, 84, 0)             |
| 0011    | Orange       | (208, 106, 26)          |
| 0100    | Dark Green   | (0, 105, 64)            |
| 0101    | Gray 1       | (128, 128, 128)         |
| 0110    | Green        | (47, 188, 26)           |
| 0111    | Yellow       | (191, 211, 90)          |
| 1000    | Dark Blue    | (64, 44, 165)           |
| 1001    | Purple       | (208, 67, 229)          |
| 1010    | Gray 2       | (128, 128, 128)         |
| 1011    | Pink         | (255, 150, 191)         |
| 1100    | Medium Blue  | (47, 149, 229)          |
| 1101    | Light Blue   | (191, 171, 255)         |
| 1110    | Aqua         | (111, 232, 191)         |
| 1111    | White        | (255, 255, 255)         |

Gray 1 and Gray 2 are visually identical but arise from different dot
patterns.

Although lores and double hires share the same set of 16 RGB color values, the
mapping from bit patterns to colors differs between the two modes. Lores
assigns colors by a 4-bit nibble index (e.g., index 2 = dark blue), while
double hires assigns colors by the NTSC phase pattern of the sliding window
(e.g., pattern 0010 = brown). The same 16 colors appear in both modes, but a
given numeric value does not produce the same color in both.

## 4.2. Sliding Window Color Assignment

Color is determined by a per-dot sliding window of 4 dots. For each of the 560
dot positions, a 4-dot window starting at that position is examined. Each dot in
the window maps to a bit based on its phase in the NTSC color cycle:

```
phase = dot_position % 4

phase 0 -> bit 3
phase 1 -> bit 2
phase 2 -> bit 1
phase 3 -> bit 0
```

If a dot within the window is on, its corresponding bit in the 4-bit pattern is
set. If the dot position falls beyond the end of the row (past dot 559), it does
not contribute to the pattern. The resulting 4-bit pattern indexes into the
16-color table to produce the color for that dot position.

## 4.3. Color vs. Hires Color

Double hires color generation differs fundamentally from standard hires:

- Standard hires uses per-byte palette bits and neighbor-based rules to
  produce 6 colors.
- Double hires ignores bit 7 entirely and uses a sliding window across the
  full dot stream to produce 16 colors from NTSC phase relationships.
- Standard hires has byte boundary effects where adjacent bytes with different
  palettes cause color shifts. Double hires has no such boundary effects
  because there is no palette bit.

# 5. Monochrome Modes

## 5.1. Overview

When a monochrome display mode is active (green screen or amber screen), the
color generation logic is bypassed entirely. Each dot is rendered individually
as either the monochrome color (if the dot's bit is on) or black (if the bit
is off). This provides a true 560x192 resolution display since every dot is
independently visible without the 4-dot color grouping.

## 5.2. Monochrome Colors

| Mode         | On Color                | Off Color       |
|--------------|-------------------------|-----------------|
| Green Screen | (152, 255, 152)         | (0, 0, 0)       |
| Amber Screen | (255, 191, 0)           | (0, 0, 0)       |

## 5.3. Rendering Behavior

In monochrome mode, all 560 dots per row are rendered independently. For each
dot, the renderer reads the corresponding bit from the aux or main byte and
writes the monochrome on or off color into a 1x2 pixel block on the
framebuffer (one dot wide, two rows tall for vertical doubling).

# 6. Rendering Pipeline

## 6.1. Mode Dispatch

Double hires mode is active when all three conditions are met:

1. TEXT is off
2. HIRES is on
3. Both DHIRES and 80COL are on

The display mode priority is:

1. If TEXT is on: render text
2. If HIRES is on and DHIRES is on and 80COL is on: render double hires
3. If HIRES is on: render standard hires
4. Otherwise: render lores

## 6.2. Display Memory Snapshot

Before rendering, the current contents of both main and auxiliary memory in
the hires region are copied into a snapshot buffer. The snapshot provides two
accessors:

- `GetMain(addr)`: returns the byte at `addr` from main memory
- `GetAux(addr)`: returns the byte at `addr` from auxiliary memory

Both accessors operate on addresses in the range $2000-$3FFF (page 1) or
$4000-$5FFF (page 2). The snapshot respects page switching: when PAGE2 is on,
the snapshot copies from $4000-$5FFF instead of $2000-$3FFF, but the accessors
always use $2000-$3FFF addressing.

## 6.3. Per-Row Rendering (Monochrome)

For each of the 192 rows:

1. Look up the row's base address from the precomputed address table.
2. For each of the 40 memory offsets:
   a. Read the aux byte and main byte at `base + offset`.
   b. The aux byte maps to screen column `offset * 2`; the main byte maps to
      screen column `offset * 2 + 1`.
   c. For each byte, extract 7 dots from bits 0-6.
   d. For each dot, write the monochrome on/off color into a 1x2 pixel block at
      position `(screenCol * 7 + bit, row * 2)`.

## 6.4. Per-Row Rendering (Color)

For each of the 192 rows:

1. **Build dot array**: read all 80 bytes (40 aux + 40 main, interleaved) and
   extract 7 dots from each, producing a 560-element boolean array.

2. **Assign colors**: for each of the 560 dot positions, compute a 4-bit
   pattern from a 4-dot window starting at that position (see section 4.2).
   Look up the color from the 16-color table and write it into a 1x2 pixel
   block at position `(dot, row * 2)`.

# 7. Soft Switches

## 7.1. Mode Selection

The following soft switches control double hires mode. All standard hires
switches (TEXT, MIXED, PAGE2, HIRES) apply as in spec 11. The additional
switches specific to double hires are:

| Switch      | Address | Access | Effect                        |
|-------------|---------|--------|-------------------------------|
| DHIRES on   | `$C05E` | R/W    | Enable double hires mode      |
| DHIRES off  | `$C05F` | R/W    | Disable double hires mode     |
| 80COL on    | `$C00D` | W      | Enable 80-column mode         |
| 80COL off   | `$C00C` | W      | Disable 80-column mode        |

Note that the DHIRES switch addresses are reversed from the typical on/off
ordering: the lower address ($C05E) enables the mode, and the higher address
($C05F) disables it.

For double hires to be active, TEXT must be off, HIRES must be on, and both
DHIRES and 80COL must be on.

## 7.2. Read Switches

| Switch    | Address | Returns                              |
|-----------|---------|--------------------------------------|
| RDDHIRES  | `$C07F` | Bit 7 high if DHIRES is on           |
| RD80COL   | `$C01F` | Bit 7 high if 80COL is on            |
| RDHIRES   | `$C01D` | Bit 7 high if HIRES is on            |
| RDPAGE2   | `$C01C` | Bit 7 high if PAGE2 is on            |

## 7.3. Page Selection

PAGE2 selects between page 1 ($2000-$3FFF) and page 2 ($4000-$5FFF) in both
main and auxiliary memory simultaneously. When PAGE2 is on, the snapshot
copies from $4000-$5FFF in both main and aux; when PAGE2 is off, it copies
from $2000-$3FFF.

When 80STORE is active, PAGE2 controls memory routing rather than page
selection: writes to $2000-$3FFF go to main memory (PAGE2 off) or auxiliary
memory (PAGE2 on). This allows programs to update double hires display memory
through the normal address space.

## 7.4. Display Redraw

Every soft switch toggle that changes display state sets a redraw flag. Writes
to the hires memory region also set this flag. The render loop checks this
flag to determine whether the framebuffer needs to be updated.

# 8. Mixed Mode

## 8.1. Current Status

Mixed mode is not currently implemented in the double hires renderer. When
MIXED is on, the entire screen is rendered as double hires graphics, including
the bottom 32 rows that should show text.

A correct implementation would:

1. Check the MIXED flag before rendering.
2. If MIXED is on, render only double hires rows 0-159.
3. Render text rows 20-23 using the text renderer for the bottom portion.

# 9. Comparison with Standard Hires

| Property              | Standard Hires        | Double Hires             |
|-----------------------|-----------------------|--------------------------|
| Dot resolution        | 280x192               | 560x192                  |
| Color resolution      | 280x192 (6 colors)    | 140x192 (16 colors)      |
| Memory per page       | 8 KB (main only)      | 16 KB (8 KB main + 8 KB aux) |
| Bit 7 usage           | Palette selector      | Unused                   |
| Color mechanism       | Neighbor + palette    | 4-dot sliding window     |
| Byte boundary effects | Yes (palette shifts)  | No                       |
| Monochrome resolution | 280x192               | 560x192                  |
| Required switches     | TEXT off, HIRES on    | TEXT off, HIRES on, DHIRES on, 80COL on |

# 10. Design Considerations

## 10.1. Sliding Window vs. Neighbor Rules

Standard hires determines color from a dot and its immediate neighbor, with a
per-byte palette bit selecting between two color pairs. Double hires instead
uses a sliding 4-dot window with no palette bit. This more closely models NTSC
artifact coloring, where 4 phases of the color subcarrier cycle combine to
produce a composite color signal. The sliding window naturally produces 16
colors (2^4 patterns) rather than the 6 colors of the neighbor-based approach.

## 10.2. Shared Address Table

The row address table is identical between standard hires and double hires.
The `a2dhires` package maintains its own copy of this table rather than
importing it from `a2hires`, keeping the packages independent.

## 10.3. Snapshot Isolation

The display snapshot captures both main and auxiliary memory before rendering
begins. This prevents tearing that could occur if the CPU modifies display
memory mid-render. The snapshot always normalizes addresses to the $2000-$3FFF
range regardless of which page is active, simplifying the rendering code.

## 10.4. 1x2 Pixel Blocks

Double hires dots are rendered as 1x2 pixel blocks (1 pixel wide, 2 pixels
tall) rather than the 2x2 blocks used by standard hires. The 560 dots per row
already fill the full 560-pixel framebuffer width, so no horizontal doubling
is needed. Vertical doubling is still applied to fill the 384-pixel
framebuffer height from 192 logical rows.
