---
Specification: 19
Drafted At: 2026-03-29
Authors:
  - Peter Evans
---

# 1. Overview

This spec describes the implementation of 80-column text mode for the Apple
//e Enhanced. 80-column text mode displays 80 characters per row instead of
40, using the same 24 rows. It is provided by the extended 80-column card,
which also supplies the auxiliary memory bank. Each row of 80 characters is
split across two memory segments: auxiliary memory provides one character per
position and main memory provides the other, with the two interleaved on
screen to form 80 visible columns.

This spec builds on three prior specs:

- **Spec 9** (40-column text mode): memory layout, address mapping, character
  encoding, rendering pipeline, and flash behavior. 80-column text reuses the
  same interleaved address scheme and character encoding.
- **Spec 16** (auxiliary memory): the 80STORE mechanism that routes accesses
  in the text page address range to main or auxiliary memory based on the
  PAGE2 flag.
- **Spec 18** (fonts): the 80-column font objects, which use 7-pixel-wide
  by 16-pixel-tall glyphs (height-doubled only, no width doubling).

# 2. Screen Geometry

## 2.1. Logical Dimensions

The 80-column text display is a grid of 80 columns by 24 rows, for a total
of 1920 character cells.

## 2.2. Pixel Dimensions

Each character cell is 7 pixels wide and 16 pixels tall (the base 7x8 glyph
doubled in height only, as described in spec 18 section 5.2).

The full screen resolution is 560 x 384 pixels:

- Width: 80 columns x 7 pixels = 560
- Height: 24 rows x 16 pixels = 384

This is the same framebuffer used by 40-column text mode (40 x 14 = 560,
24 x 16 = 384) and all other display modes.

# 3. Activation

## 3.1. Required State

80-column text mode is active when all of the following are true:

- **TEXT** is on (soft switch `$C051`)
- **80COL** is on (soft switch `$C00D`)

The 80STORE switch (`$C001`) must also be enabled for the memory routing
described in section 4 to work correctly. Software that uses 80-column text
mode invariably enables 80STORE, since without it there is no way to
independently address the auxiliary half of the text page. However, the
display mode itself is selected by the TEXT and 80COL flags; 80STORE affects
only memory routing, not mode selection.

## 3.2. Typical Activation Sequence

A program enabling 80-column text mode will typically write to these soft
switches in order:

1. `$C001` -- 80STORE on (enable display memory routing)
2. `$C00D` -- 80COL on (enable 80-column display)
3. `$C051` -- TEXT on (enable text mode, if not already on)

To return to 40-column text mode:

1. `$C00C` -- 80COL off
2. `$C000` -- 80STORE off (optional, but standard practice)

# 4. Memory Layout

## 4.1. Two-Bank Interleave

In 80-column mode, each of the 40 byte positions in the text page produces
two screen characters instead of one. For a given byte offset within the text
page, the auxiliary memory byte appears at the left screen column and the
main memory byte appears at the right screen column.

Both banks use the same address range (`$0400`-`$07FF`) with the same
interleaved address mapping defined in spec 9 section 3.3. The difference is
which memory segment the byte is read from.

## 4.2. Column Mapping

For a byte at offset `N` within the text page (where `N` ranges from 0 to
1023), the 40-column address tables from spec 9 give a row `R` and a
40-column position `C` (where `R` and `C` are -1 for hole bytes). In
80-column mode, each such position maps to two screen columns:

- **Screen column `2 * C`**: the byte at offset `N` in auxiliary memory
- **Screen column `2 * C + 1`**: the byte at offset `N` in main memory

The auxiliary byte is always the left (even) column of the pair and the main
byte is always the right (odd) column.

## 4.3. Addressing via 80STORE and PAGE2

Software writes to the 80-column text page using the 80STORE mechanism
described in spec 16. When 80STORE is on:

- **PAGE2 off**: writes to `$0400`-`$07FF` go to main memory.
- **PAGE2 on**: writes to `$0400`-`$07FF` go to auxiliary memory.

A typical sequence to write character `A` at 80-column screen column 0 of
row 0 and character `B` at screen column 1 of the same row:

1. Write `$C055` (PAGE2 on) -- route text page to auxiliary memory
2. Store `A` at `$0400` -- goes to aux, appears at screen column 0
3. Write `$C054` (PAGE2 off) -- route text page back to main memory
4. Store `B` at `$0400` -- goes to main, appears at screen column 1

The same offset in both segments maps to the same row; the only difference is
which column of the pair the character occupies.

# 5. Display Memory Snapshot

## 5.1. Snapshot Contents

The existing snapshot system (spec 9 section 7.2) copies a single 1 KB
region for text rendering. For 80-column mode, the snapshot must capture two
1 KB regions:

- **Main text**: `$0400`-`$07FF` from main memory
- **Auxiliary text**: `$0400`-`$07FF` from auxiliary memory

Both regions are always captured from the page 1 address range. In 80-column
mode with 80STORE, PAGE2 controls which segment is written to, but both
segments contribute to the display simultaneously. The display always renders
from the page 1 addresses in both banks.

## 5.2. Snapshot Interface

The snapshot must provide access to both banks for 80-column rendering. The
existing `Get(addr)` method returns the appropriate byte for 40-column mode
(from whichever single bank was captured). For 80-column mode, the renderer
needs two accessor methods:

- `GetMain(addr)` -- returns the byte from the main memory snapshot at the
  given text page address
- `GetAux(addr)` -- returns the byte from the auxiliary memory snapshot at
  the given text page address

These are analogous to the existing `GetMain`/`GetAux` methods used for
double hi-res, extended to cover the text page range.

# 6. Rendering Pipeline

## 6.1. Mode Dispatch

The display renderer checks display state flags to determine the rendering
mode. The updated priority is:

1. If TEXT is on and 80COL is on: render 80-column text
2. If TEXT is on and 80COL is off: render 40-column text
3. If HIRES is on: render hi-res or double hi-res graphics
4. Otherwise: render lo-res graphics

This adds a check for the 80COL flag before the existing TEXT-mode rendering
path.

## 6.2. Rendering Loop

The 80-column text renderer iterates over the same 1024-byte address range
as the 40-column renderer, but produces two characters per visible byte
position:

```
for offset = 0 to 1023:
    row = addressRows[offset]
    col40 = addressCols[offset]

    if row == -1 or col40 == -1:
        continue    // hole byte, skip

    auxChar = auxSnapshot.get($0400 + offset)
    mainChar = mainSnapshot.get($0400 + offset)

    auxGlyph = font.glyph(auxChar)
    mainGlyph = font.glyph(mainChar)

    leftX = (col40 * 2) * 7       // even column (aux)
    rightX = (col40 * 2 + 1) * 7  // odd column (main)
    y = row * 16

    screen.blit(leftX, y, auxGlyph)
    screen.blit(rightX, y, mainGlyph)
```

The address lookup tables from spec 9 are reused without modification. The
only change is that each table entry produces two screen columns instead of
one, and the glyph width is 7 instead of 14.

## 6.3. Font Selection

The 80-column renderer uses the 80-column font objects defined in spec 18:

- **Primary font** (`SystemFont80`): 7x16 glyphs, `$40`-`$7F` rendered as
  inverse (for flash-on state).
- **Flash-alternate font** (`SystemFont80FlashAlt`): 7x16 glyphs,
  `$40`-`$7F` rendered as normal (for flash-off state).
- **Alternate charset font** (`SystemFont80Alt`): 7x16 glyphs, `$40`-`$5F`
  holds MouseText.

Font selection follows the same logic as 40-column mode:

- If ALTCHAR is on, use `SystemFont80Alt` (no flash behavior).
- If ALTCHAR is off and `flashOn` is true, use `SystemFont80`.
- If ALTCHAR is off and `flashOn` is false, use `SystemFont80FlashAlt`.

## 6.4. Flash

Flash behavior is identical to 40-column mode (spec 9 section 8). The flash
state is global and derived from the CPU cycle counter. The same dual-font
swap approach is used. Both the auxiliary and main characters in a given frame
use the same font, since flash state is global.

## 6.5. Monochrome Modes

Monochrome recoloring (green screen, amber screen) applies identically to
80-column mode. Each glyph is post-processed in the same way as described in
spec 9 section 7.5.

# 7. Mixed Mode

When MIXED is on and TEXT is off, the bottom 4 rows of the screen display
text while the upper portion displays graphics. In 80-column mode (80COL on),
these bottom 4 rows render as 80-column text, not 40-column. The mixed-mode
text rows use the same two-bank interleave and 80-column font as full-screen
80-column text.

The mixed-mode text region covers rows 20-23, which correspond to the
following address offsets within the text page (using the interleaved mapping
from spec 9 section 3.3):

| Row | Base Offset (hex) | Offset Range (relative to `$0400`) |
|-----|-------------------|------------------------------------|
| 20  | `$250`            | `$250`-`$277`                      |
| 21  | `$2D0`            | `$2D0`-`$2F7`                      |
| 22  | `$350`            | `$350`-`$377`                      |
| 23  | `$3D0`            | `$3D0`-`$3F7`                      |

Each row spans 40 bytes from its base offset. These are the third row within
each 128-byte group (rows 16-23 occupy the third position in groups 0-7).
The graphics renderers are responsible for invoking the text renderer for
these rows when mixed mode is active; the text renderer itself does not need
to know about mixed mode.

# 8. Page 2

## 8.1. Without 80STORE

When 80STORE is off and PAGE2 is on, 80-column text mode reads from page 2:

- Main characters: `$0800`-`$0BFF` in main memory
- Auxiliary characters: `$0800`-`$0BFF` in auxiliary memory

This is the straightforward case where PAGE2 selects a different 1 KB region
in each segment.

## 8.2. With 80STORE

When 80STORE is on, PAGE2 is repurposed for memory routing (spec 16 section
6.2). The display always renders from the page 1 address range
(`$0400`-`$07FF`) in both main and auxiliary memory. PAGE2 controls which
segment CPU writes target, not which page the display reads from.

This is the normal configuration for 80-column text mode.

# 9. Design Considerations

## 9.1. Reuse of Address Tables

The 40-column address lookup tables (spec 9 section 3.5) map each of the
1024 byte offsets to a row and a 40-column position. The 80-column renderer
reuses these tables directly: it simply multiplies the 40-column position by
2 to obtain the screen column for the auxiliary character and adds 1 for the
main character. No separate 80-column address tables are needed.

## 9.2. Single Rendering Function

A unified text rendering function could accept both 40-column and 80-column
fonts and branch on the mode. However, the two modes differ in how they
source characters (one bank vs. two banks) and how they compute screen
coordinates (single column vs. column pair). A separate 80-column rendering
function is clearer and avoids per-character branching in the inner loop. The
40-column renderer remains unchanged.

## 9.3. Snapshot Expansion

The snapshot must capture auxiliary text memory in addition to main text
memory for 80-column rendering. This adds 1 KB of data to the snapshot copy
(auxiliary `$0400`-`$07FF`). The cost is negligible and follows the same
pattern already established for double hi-res, which captures both main and
auxiliary hi-res pages.

## 9.4. Framebuffer Compatibility

The 560x384 framebuffer accommodates 80-column text without any changes.
80 columns x 7 pixels = 560 pixels wide, and 24 rows x 16 pixels = 384
pixels tall. The framebuffer was designed with this dual-use in mind (spec 9
section 2.2).
