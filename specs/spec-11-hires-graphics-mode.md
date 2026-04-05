---
Specification: 11
Category: Graphics
Drafted At: 2026-03-24
Authors:
  - Peter Evans
---

# 1. Overview

This spec describes the implementation of high-resolution (hires) graphics
mode for the Apple //e. Hires mode displays a 280x192 grid of individually
addressable dots, where each dot can produce one of six colors (black, white,
purple, green, blue, or orange) depending on its bit position, its neighbors,
and a per-byte palette selector. The hires display occupies 8 KB of memory
($2000-$3FFF for page 1, $4000-$5FFF for page 2).

# 2. Screen Geometry

## 2.1. Logical Dimensions

The hires display is 280 dots wide by 192 rows tall. Each row is stored in 40
consecutive bytes, with each byte contributing 7 dots (the low 7 bits). The
high bit of each byte selects a color palette rather than a dot.

## 2.2. Dot Dimensions

Each logical dot is rendered as a 2x2 pixel block on the 560x384 framebuffer:

- Width: 280 dots x 2 pixels = 560 pixels
- Height: 192 rows x 2 pixels = 384 pixels

This 2x scale matches the framebuffer dimensions shared with text and lores
modes.

# 3. Display Memory

## 3.1. Memory Region

Page 1 occupies $2000-$3FFF (8,192 bytes). Page 2 occupies $4000-$5FFF. Only
7,680 bytes per page are visible (40 bytes x 192 rows); the remaining 512
bytes are "hole" bytes that exist in memory but do not appear on screen.

## 3.2. Address Interleaving

Like text and lores modes, hires row addresses are interleaved rather than
sequential. The 192 rows are organized into three groups of 64 rows, and
within each group the addresses follow a pattern based on three nested
subdivisions.

The base address for row `r` on page 1 is:

```
group    = r / 64          (0, 1, or 2)
section  = (r % 64) / 8    (0-7)
line     = r % 8           (0-7)

base = $2000 + (line * $0400) + (section * $0080) + (group * $0028)
```

Each row occupies 40 consecutive bytes starting from its base address. The 40
bytes at offsets 40-127 from each section start are hole bytes -- they exist in
the address space but are not displayed.

For page 2, the same formula applies with $4000 as the starting address
instead of $2000.

### 3.2.1. Address Table

The full set of 192 row base addresses (page 1) is precomputed at package
initialization:

| Rows    | Base Addresses                                              |
|---------|-------------------------------------------------------------|
| 0-7     | $2000, $2400, $2800, $2C00, $3000, $3400, $3800, $3C00     |
| 8-15    | $2080, $2480, $2880, $2C80, $3080, $3480, $3880, $3C80     |
| 16-23   | $2100, $2500, $2900, $2D00, $3100, $3500, $3900, $3D00     |
| 24-31   | $2180, $2580, $2980, $2D80, $3180, $3580, $3980, $3D80     |
| 32-39   | $2200, $2600, $2A00, $2E00, $3200, $3600, $3A00, $3E00     |
| 40-47   | $2280, $2680, $2A80, $2E80, $3280, $3680, $3A80, $3E80     |
| 48-55   | $2300, $2700, $2B00, $2F00, $3300, $3700, $3B00, $3F00     |
| 56-63   | $2380, $2780, $2B80, $2F80, $3380, $3780, $3B80, $3F80     |
| 64-71   | $2028, $2428, $2828, $2C28, $3028, $3428, $3828, $3C28     |
| 72-79   | $20A8, $24A8, $28A8, $2CA8, $30A8, $34A8, $38A8, $3CA8     |
| 80-87   | $2128, $2528, $2928, $2D28, $3128, $3528, $3928, $3D28     |
| 88-95   | $21A8, $25A8, $29A8, $2DA8, $31A8, $35A8, $39A8, $3DA8     |
| 96-103  | $2228, $2628, $2A28, $2E28, $3228, $3628, $3A28, $3E28     |
| 104-111 | $22A8, $26A8, $2AA8, $2EA8, $32A8, $36A8, $3AA8, $3EA8     |
| 112-119 | $2328, $2728, $2B28, $2F28, $3328, $3728, $3B28, $3F28     |
| 120-127 | $23A8, $27A8, $2BA8, $2FA8, $33A8, $37A8, $3BA8, $3FA8     |
| 128-135 | $2050, $2450, $2850, $2C50, $3050, $3450, $3850, $3C50     |
| 136-143 | $20D0, $24D0, $28D0, $2CD0, $30D0, $34D0, $38D0, $3CD0     |
| 144-151 | $2150, $2550, $2950, $2D50, $3150, $3550, $3950, $3D50     |
| 152-159 | $21D0, $25D0, $29D0, $2DD0, $31D0, $35D0, $39D0, $3DD0     |
| 160-167 | $2250, $2650, $2A50, $2E50, $3250, $3650, $3A50, $3E50     |
| 168-175 | $22D0, $26D0, $2AD0, $2ED0, $32D0, $36D0, $3AD0, $3ED0     |
| 176-183 | $2350, $2750, $2B50, $2F50, $3350, $3750, $3B50, $3F50     |
| 184-191 | $23D0, $27D0, $2BD0, $2FD0, $33D0, $37D0, $3BD0, $3FD0     |

## 3.3. Byte Layout

Each byte in the hires region encodes 7 dots plus a palette selector:

```
bit 7    bit 6   bit 5   bit 4   bit 3   bit 2   bit 1   bit 0
palette  dot 6   dot 5   dot 4   dot 3   dot 2   dot 1   dot 0
```

- **Bits 0-6**: each bit controls one dot. 1 = on, 0 = off. Bit 0 is the
  leftmost dot within the byte; bit 6 is the rightmost.
- **Bit 7**: selects the color palette for all 7 dots in this byte. 0 =
  purple/green, 1 = blue/orange.

The 40 bytes per row produce 40 x 7 = 280 dots per row.

# 4. Color Generation

## 4.1. The Six Hires Colors

The hires display produces six colors:

| Color  | RGB                     |
|--------|-------------------------|
| Black  | (0, 0, 0)               |
| White  | (255, 255, 255)         |
| Purple | (208, 67, 229)          |
| Green  | (47, 188, 26)           |
| Blue   | (47, 149, 229)          |
| Orange | (208, 106, 26)          |

## 4.2. Color Palettes

The two palettes each contain a pair of complementary colors:

- **Palette 0** (bit 7 = 0): purple and green
- **Palette 1** (bit 7 = 1): blue and orange

Within each palette, the two colors alternate by dot position. Even-numbered
dots (0, 2, 4, ...) use the first color in the pair; odd-numbered dots (1, 3,
5, ...) use the second. Dot positions are counted across the full 280-dot row,
not within each byte.

## 4.3. Color Rules

The color of each dot is determined by examining the dot itself and its
left neighbor:

| This Dot | Previous Dot | Result                                |
|----------|--------------|---------------------------------------|
| on       | on           | White                                 |
| on       | off          | Palette color for this dot's position |
| off      | on           | Palette color for the *opposite* position (the color the previous dot did not get) |
| off      | off          | Black                                 |

The "opposite position" rule for off-after-on means: if the current dot is at
an even position (which would normally produce the first palette color), an
off-after-on dot instead produces the second palette color, and vice versa.
This creates the characteristic hires color fringing where a single on bit
affects both its own dot and the following off dot.

The first dot in a row (position 0) has no previous dot, so it behaves as if
the previous dot were off.

## 4.4. White and Black Generation

Two consecutive on bits always produce white, regardless of palette. This is
the only way to produce white in hires mode. Similarly, two consecutive off
bits always produce black.

Because white requires consecutive on bits, a single isolated on bit never
produces white -- it produces a colored dot instead. This is the fundamental
mechanism behind Apple II hires color: the hardware generates color artifacts
from the timing of individual dots relative to the NTSC color reference
signal.

# 5. Byte Boundary Effects

## 5.1. Overview

When adjacent bytes use different palettes (one has bit 7 = 0, the other has
bit 7 = 1), the dots at the byte boundary may shift color. This is an
approximation of NTSC artifact color behavior, where the color reference phase
changes when the palette bit changes.

## 5.2. Boundary Detection

The last dot of each 7-dot byte group (dot position 6 within the byte) is
marked as a boundary dot. After the initial color assignment pass, boundary
dots are re-examined: if the palette of the left byte differs from the
palette of the right byte, a color shift is applied.

## 5.3. Color Shift Rules

The shift depends on the right-hand dot's color:

**When the right dot is purple or green** (right byte is palette 0, left byte
is palette 1): both the left boundary dot and the right dot shift to darker
variants:

| Original | Shifted       | RGB                     |
|----------|---------------|-------------------------|
| Purple   | Dark Purple   | (62, 49, 121)           |
| Green    | Dark Green    | (63, 76, 18)            |

**When the right dot is blue or orange** (right byte is palette 1, left byte
is palette 0): only the right dot shifts, and only if the left boundary dot
is not black:

| Original | Shifted       | RGB                     |
|----------|---------------|-------------------------|
| Blue     | Light Purple  | (187, 175, 246)         |
| Orange   | Light Green   | (189, 234, 134)         |

If the left boundary dot is black when the right dot is blue or orange, no
shift occurs.

# 6. Monochrome Modes

## 6.1. Overview

When a monochrome display mode is active (green screen or amber screen), the
color generation logic is bypassed entirely. Each dot is rendered as either the
monochrome color (if the dot's bit is on) or black (if the bit is off).

## 6.2. Monochrome Colors

| Mode         | On Color                | Off Color       |
|--------------|-------------------------|-----------------|
| Green Screen | (152, 255, 152)         | (0, 0, 0)       |
| Amber Screen | (255, 191, 0)           | (0, 0, 0)       |

## 6.3. Rendering Behavior

In monochrome mode, the palette bit (bit 7), dot position, and neighbor
state are all irrelevant. The only input is whether each dot's bit is set.
This means monochrome mode produces no color fringing, no boundary shifting,
and no white-from-consecutive-bits behavior -- just a 1-bit-per-dot display.

# 7. Rendering Pipeline

## 7.1. Mode Dispatch

The rendering pipeline uses the same mode dispatch as text and lores modes.
Hires mode is active when TEXT is off and HIRES is on. The display mode
priority is:

1. If TEXT is on: render text
2. If HIRES is on: render hi-res graphics
3. Otherwise: render lo-res graphics

## 7.2. Display Memory Snapshot

Before rendering, the current contents of the hires memory region are copied
into a snapshot buffer. The snapshot respects page switching and 80STORE
settings (see section 8.3). This is the same snapshot mechanism used by text
and lores modes, extended to cover the $2000-$3FFF region.

## 7.3. Per-Row Rendering

Rendering proceeds one row at a time across all 192 rows. For each row:

1. **Fill dots**: read 40 bytes from the snapshot at the row's base address.
   For each byte, extract 7 dot on/off states from bits 0-6 and the palette
   from bit 7. This produces a 280-element array of Dot structs, each
   containing its on/off state and palette.

2. **Assign colors**: if monochrome mode is active, assign the monochrome
   color (on) or black (off) to each dot and skip to step 4. Otherwise, scan
   left to right applying the color rules from section 4.3.

3. **Apply boundary shifts**: after color assignment, re-examine each byte
   boundary. If the palettes differ, apply the shifts described in section 5.3.

4. **Blit to framebuffer**: for each of the 280 dots, write its color into a
   2x2 pixel block on the framebuffer at position (dot * 2, row * 2).

## 7.4. Dot Struct

Each dot is represented as a struct containing:

- `on`: whether the dot's bit is set
- `boundary`: whether this is the last dot in a 7-dot byte group
- `palette`: which palette this byte uses (purple/green or blue/orange)
- `color`: the final RGBA color after all rules are applied

# 8. Soft Switches

## 8.1. Mode Selection

The following soft switches control whether hires mode is active:

| Switch    | Address | Access | Effect                       |
|-----------|---------|--------|------------------------------|
| TEXT off  | `$C050` | R/W    | Disable text mode (graphics) |
| TEXT on   | `$C051` | R/W    | Enable text mode             |
| MIXED off | `$C052` | R/W    | Disable mixed mode           |
| MIXED on  | `$C053` | R/W    | Enable mixed mode            |
| PAGE2 off | `$C054` | R/W    | Select display page 1        |
| PAGE2 on  | `$C055` | R/W    | Select display page 2        |
| HIRES off | `$C056` | R/W    | Disable hires mode           |
| HIRES on  | `$C057` | R/W    | Enable hires mode            |

For hires mode to be active, TEXT must be off and HIRES must be on.

## 8.2. Read Switches

| Switch   | Address | Returns                              |
|----------|---------|--------------------------------------|
| RDHIRES  | `$C01D` | Bit 7 high if HIRES is on            |
| RDPAGE2  | `$C01C` | Bit 7 high if PAGE2 is on            |
| RDMIXED  | `$C01B` | Bit 7 high if MIXED is on            |

## 8.3. Page Selection and 80STORE

PAGE2 selects between page 1 ($2000-$3FFF) and page 2 ($4000-$5FFF). When
80STORE is active, the interaction between PAGE2 and HIRES affects memory
routing:

- **80STORE off**: PAGE2 selects page 1 or page 2 in main memory as normal.
- **80STORE on, HIRES on**: PAGE2 controls whether writes to $2000-$3FFF go
  to main memory (PAGE2 off) or auxiliary memory (PAGE2 on). The display
  always reads from the selected segment. This is used by double hires mode
  and by programs that use auxiliary memory for page-flipping.

## 8.4. Display Redraw

Every soft switch toggle that changes display state also sets a redraw flag.
The render loop checks this flag to determine whether the framebuffer needs to
be updated. Writes to the hires memory region ($2000-$3FFF) also set this
flag.

# 9. Mixed Mode

## 9.1. Overview

When MIXED is on and TEXT is off, the display shows a split screen: graphics
in the upper portion and 4 rows of text at the bottom. For hires mode, the
top 160 rows are rendered as hires graphics, and the bottom 32 rows (4 text
rows x 8 scan lines each) are rendered as text.

## 9.2. Screen Layout

In mixed mode, the 280x192 hires display is reduced to 280x160. The top 160
rows are rendered as hires graphics. The bottom 32 rows (4 text rows x 8 scan
lines each) are rendered as text characters instead of graphics.

The visible area is:

| Region         | Hires Rows | Text Rows | Pixel Rows | Content        |
|----------------|------------|-----------|------------|----------------|
| Graphics       | 0-159      | 0-19      | 0-319      | Hires dots     |
| Text           | 160-191    | 20-23     | 320-383    | Text glyphs    |

## 9.3. Rendering Behavior

The hires renderer only iterates over rows 0-159 when MIXED is on, instead of
the full 0-191.

The text portion is rendered by the normal text renderer using the text/lores
memory region ($0400-$07FF). The display dispatch is responsible for calling
both the hires renderer (which stops at row 160) and the text renderer (which
renders only the bottom 4 rows).

## 9.4. Mode Dispatch

When MIXED is on, TEXT is off, and HIRES is on, the display dispatch must:

1. Call the hires renderer, which renders only rows 0-159.
2. Call the text renderer for text rows 20-23.

The text renderer uses the same font, flash state, and monochrome settings as
normal text mode. The only difference is that it renders just the bottom 4 rows
rather than all 24.

# 10. Design Considerations

## 10.1. Two-Pass Color Assignment

The renderer performs two passes over each row: one to fill dot on/off states
from memory bytes, and a second to assign colors. This separation keeps the
bit-extraction logic simple and independent of the color rules. On modern
hardware, the overhead of a second pass over 280 dots is negligible.

## 10.2. Precomputed Row Addresses

The 192 row base addresses are stored in a precomputed lookup table rather
than calculated per frame. This avoids the division and modulo arithmetic of
the interleaved address formula on every render.

## 10.3. Boundary Shift as Post-Processing

Byte boundary color shifts are applied as a post-processing step after the
main color assignment, rather than being integrated into the color rules. This
keeps the main color logic (section 4.3) clean and makes the boundary behavior
easy to adjust independently.

## 10.4. Double Hires (Out of Scope)

The Apple //e also supports double high-resolution graphics, which uses both
main and auxiliary memory to produce 560 dots per row with 16 colors. Double
hires requires both the DHIRES and 80COL soft switches to be active in
addition to HIRES. This mode is not covered by this spec.
