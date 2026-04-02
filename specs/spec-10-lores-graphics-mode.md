---
Specification: 10
Category: Graphics
Drafted At: 2026-03-24
Authors:
  - Peter Evans
---

# 1. Overview

This spec describes the implementation of low-resolution (lores) graphics mode
for the Apple //e. Lores mode displays a 40x48 grid of colored blocks, where
each block can be one of 16 colors. It is the default graphics mode: when TEXT
is off and HIRES is off, the display renders lores graphics. Lores shares the
same memory region as text mode ($0400-$07FF for page 1), but interprets each
byte as two vertically stacked color blocks rather than a character glyph.

# 2. Screen Geometry

## 2.1. Logical Dimensions

The lores display is a grid of 40 columns by 48 rows, for a total of 1,920
color blocks. Each memory byte encodes two vertically adjacent blocks, so the
960 visible bytes of the text page produce 1,920 blocks.

## 2.2. Block Dimensions

Each color block is 14 pixels wide and 8 pixels tall. These dimensions are
derived from the shared framebuffer geometry:

- Width: the 560-pixel-wide framebuffer divided by 40 columns = 14 pixels per
  column. This matches the 40-column text cell width.
- Height: the 384-pixel-tall framebuffer divided by 48 rows = 8 pixels per
  row. This is half the text cell height (16 pixels), since each text row
  corresponds to two lores block rows.

The full screen resolution is 560 x 384 pixels:

- Width: 40 columns x 14 pixels = 560
- Height: 48 rows x 8 pixels = 384

# 3. Display Memory

## 3.1. Shared Memory with Text Mode

Lores graphics uses the same memory region and address mapping as 40-column
text mode. Page 1 occupies $0400-$07FF (1024 bytes), page 2 occupies
$0800-$0BFF, and the interleaved address-to-position mapping is identical.
See spec 9, sections 3.1 through 3.5, for the full description of page
layout, address interleaving, and the base address formula.

## 3.2. Nibble Encoding

The difference from text mode is in interpretation. In text mode, each byte
selects a character glyph and display mode. In lores mode, each byte encodes
two color blocks using its two nibbles:

- **Low nibble** (bits 0-3): color of the **top** block
- **High nibble** (bits 4-7): color of the **bottom** block

For a byte at text row `r`, column `c`, the two blocks occupy lores rows
`r * 2` (top) and `r * 2 + 1` (bottom).

For example, a byte value of $D4 at row 0, column 0 produces:

- Top block (row 0): color index 4 (dark green), from the low nibble ($D4 &
  $0F = $04)
- Bottom block (row 1): color index 13 (yellow), from the high nibble ($D4 >>
  4 = $0D)

# 4. Color Palette

## 4.1. The 16 Colors

Each nibble value (0-15) selects one of 16 colors. The palette matches the
standard Apple II lores color assignments:

| Index | Name        | RGB                     |
|-------|-------------|-------------------------|
| 0     | Black       | (0, 0, 0)               |
| 1     | Magenta     | (144, 23, 64)           |
| 2     | Dark Blue   | (64, 44, 165)           |
| 3     | Purple      | (208, 67, 229)          |
| 4     | Dark Green  | (0, 105, 64)            |
| 5     | Gray 1      | (128, 128, 128)         |
| 6     | Medium Blue | (47, 149, 229)          |
| 7     | Light Blue  | (191, 171, 255)         |
| 8     | Brown       | (64, 84, 0)             |
| 9     | Orange      | (208, 106, 26)          |
| 10    | Gray 2      | (128, 128, 128)         |
| 11    | Pink        | (255, 150, 191)         |
| 12    | Light Green | (47, 188, 26)           |
| 13    | Yellow      | (191, 211, 90)          |
| 14    | Aquamarine  | (111, 232, 191)         |
| 15    | White       | (255, 255, 255)         |

Gray 1 (index 5) and Gray 2 (index 10) are identical in color.

## 4.2. Color Block Construction

Each color block is a solid rectangle of 14 x 8 pixels, filled entirely with
the RGBA value from the palette. The 16 color blocks are constructed once at
initialization and reused for every render. There is no per-pixel variation
within a block -- each block is a flat fill.

# 5. Monochrome Modes

## 5.1. Overview

When a monochrome display mode is active (green screen or amber screen), the
16-color palette is replaced with shaded variations of a single monochrome
color. This simulates the appearance of lores graphics on a monochrome monitor,
where different colors appear as different brightness levels.

## 5.2. Shade Mapping

Each of the 16 color indices is assigned a shade intensity level:

| Shade   | Intensity | Colors                                                        |
|---------|-----------|---------------------------------------------------------------|
| Light   | 75%       | Light Blue (7), Light Green (12)                              |
| Medium  | 50%       | Magenta (1), Purple (3), Gray 1 (5), Medium Blue (6),         |
|         |           | Brown (8), Orange (9), Gray 2 (10), Pink (11), Yellow (13),   |
|         |           | Aquamarine (14)                                               |
| Dark    | 25%       | Dark Blue (2), Dark Green (4)                                 |

Two special cases are handled separately:

- **Black (0)**: always rendered as pure black (0, 0, 0), regardless of
  monochrome mode.
- **White (15)**: always rendered as the full monochrome color with no shading.

## 5.3. Shade Computation

For a given monochrome base color (green or amber) and intensity percentage,
each RGB component is scaled independently:

```
shaded.R = baseColor.R * intensity
shaded.G = baseColor.G * intensity
shaded.B = baseColor.B * intensity
```

The alpha channel is always $FF.

For example, green screen (base color 152, 255, 152) at 50% intensity produces
(76, 127, 76).

## 5.4. Precomputation

The 16 monochrome blocks for each mode (green and amber) are computed once at
package initialization. At render time, the appropriate precomputed block is
selected by color index and monochrome mode with no per-pixel calculation.

# 6. Rendering Pipeline

## 6.1. Mode Dispatch

The rendering pipeline uses the same mode dispatch described in spec 9,
section 7.3. Lores mode is the fallback: it is active when both TEXT and HIRES
are off. The display mode priority is:

1. If TEXT is on: render text
2. If HIRES is on: render hi-res graphics
3. Otherwise: render lo-res graphics

The mode dispatch, snapshot, and framebuffer output stages are shared across
all display modes. Only the inner rendering loop (section 6.3) is specific to
lores.

## 6.2. Display Memory Snapshot

Before rendering, the current contents of the text/lores memory region are
copied into a snapshot buffer. This is the same snapshot mechanism described in
spec 9, section 7.2. The snapshot respects page switching and 80STORE settings.

## 6.3. Rendering Loop

The renderer iterates over every address in the 1024-byte snapshot buffer
(offsets $000 through $3FF). The snapshot already contains the correct page's
data (page 1 or page 2), so the rendering loop indexes directly into the
snapshot without adding a page base address:

```
for offset = 0 to 1023:
    row = addressRows[offset]
    col = addressCols[offset]

    if row < 0 or col < 0:
        continue    // hole byte, skip

    x = col * blockWidth     // blockWidth = 14
    y = row * blockHeight    // blockHeight = 8

    byte = snapshot.get(offset)

    blit(x, y,     colorBlock[byte & $0F])    // top block (low nibble)
    blit(x, y + 8, colorBlock[byte >> 4])     // bottom block (high nibble)
```

## 6.4. Address Lookup Tables

The address-to-position mapping uses the same precomputed lookup tables as
text mode, with one difference: the row values are doubled. Where text mode
maps to screen rows 0-23, the lores tables map to half-rows 0, 2, 4, ..., 46.
Multiplying a half-row value by the 8-pixel block height gives the correct
pixel Y coordinate, and the `y + 8` offset for the bottom block fills the
second half of each text-row-height cell.

The column mapping is identical to text mode (0-39). Hole bytes are marked
with -1 in both tables, same as text mode.

# 7. Soft Switches

## 7.1. Mode Selection

The following soft switches control whether lores mode is active and which
page is displayed:

| Switch     | Address | Access | Effect                          |
|------------|---------|--------|---------------------------------|
| TEXT off   | `$C050` | R/W    | Disable text mode (graphics)    |
| TEXT on    | `$C051` | R/W    | Enable text mode                |
| MIXED off  | `$C052` | R/W    | Disable mixed mode              |
| MIXED on   | `$C053` | R/W    | Enable mixed mode (4 text rows) |
| PAGE2 off  | `$C054` | R/W    | Select display page 1           |
| PAGE2 on   | `$C055` | R/W    | Select display page 2           |
| HIRES off  | `$C056` | R/W    | Disable hi-res mode (lo-res)    |
| HIRES on   | `$C057` | R/W    | Enable hi-res mode              |

For lores mode to be active, TEXT must be off and HIRES must be off. When
both conditions are met, the display renders lores graphics.

## 7.2. Irrelevant Switches

The 80COL and ALTCHARSET switches have no effect on lores rendering. 80COL
controls 80-column text mode, and ALTCHARSET selects an alternate character
set -- both are text-only features. The 80STORE switch affects page selection
logic (see spec 9, section 7.2) but does not change lores rendering behavior.

## 7.3. Page Selection

PAGE2 selects between display page 1 ($0400-$07FF) and display page 2
($0800-$0BFF). The snapshot mechanism resolves which page to copy before
rendering begins (see section 6.2). The rendering loop itself is identical
regardless of which page is active -- it always iterates over the 1024-byte
snapshot buffer.

# 8. Mixed Mode

## 8.1. Overview

When MIXED is on and TEXT is off, the display shows a split screen: graphics in
the upper portion and 4 rows of text at the bottom. For lores mode, this means
the top 40 lores rows (20 text rows) are rendered as colored blocks, and the
bottom 4 text rows (rows 20-23) are rendered as text characters.

## 8.2. Current Status

Mixed mode is not currently implemented in the lores renderer. The MIXED soft
switch state is tracked (spec 9, section 6.1), but the lores rendering loop
does not check it. When mixed mode is active, the entire screen is rendered as
lores blocks, including the bottom 4 rows that should show text.

A correct implementation would:

1. Check the MIXED flag before rendering.
2. If MIXED is on, skip lores rendering for text rows 20-23 (lores rows 40-47).
3. Render those 4 rows using the text renderer instead.

# 9. Design Considerations

## 9.1. Shared Infrastructure with Text Mode

Lores graphics reuses almost all of the text mode infrastructure: the same
memory region, the same interleaved address mapping, the same snapshot
mechanism, and the same soft switch handling. The only differences are the
interpretation of each byte (nibble pair vs. character code) and the rendered
output (solid color blocks vs. character glyphs).

## 9.2. Precomputed Color Blocks

All 16 color blocks (and 32 monochrome blocks -- 16 for green, 16 for amber)
are allocated once at initialization. The renderer simply indexes into the
appropriate array by nibble value. This avoids any per-pixel color computation
in the render loop.

## 9.3. Double Lores (Out of Scope)

The Apple //e also supports double lores graphics, which uses auxiliary
memory to double the horizontal resolution to 80x48. This mode is not covered
by this spec.

## 9.4. Monochrome Shade Rationale

On real hardware, lores colors are generated by filling each pixel with a
repeating 4-bit binary pattern timed to the Colorburst reference signal. A
color display interprets these patterns as color; a monochrome display (or one
with Colorburst disabled) reveals the underlying bit patterns directly. The
brightness on a monochrome monitor therefore depends on the density of "on"
bits in the pattern.

The shade mapping approximates this bit density:

- **Dark (25%)**: patterns with 1 out of 4 bits on (e.g., Dark Blue, Dark
  Green)
- **Medium (50%)**: patterns with 2 out of 4 bits on (e.g., Magenta, Purple,
  Orange). Gray 1 (0101) and Gray 2 (1010) also fall here -- their "on" bits
  are polar opposites on the quadrature color signal, canceling out to produce
  identical gray on a color display.
- **Light (75%)**: patterns with 3 out of 4 bits on (e.g., Light Blue, Light
  Green)

Black (0 bits on) and white (4 bits on) are special-cased to pure black and
the full monochrome color respectively.
