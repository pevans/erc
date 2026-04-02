---
Specification: 18
Category: Graphics
Drafted At: 2026-03-28
Authors:
  - Peter Evans
---

# 1. Overview

This spec defines the bitmap fonts used by erc to render text for the Apple
//e. The Apple //e has a built-in character generator that produces glyphs on
a 7-dot-wide by 8-dot-tall grid. These glyphs are used in both 40-column and
80-column text modes, with different scaling applied to fit the shared 560x384
framebuffer.

The spec covers the base glyph format, the three glyph sets (uppercase,
special/punctuation, and lowercase), how those sets are mapped across the full
256-entry character space, display modes (normal, inverse, and flash), and the
scaling transformations used for 40-column and 80-column rendering.

The precise bitmap data for each glyph is recorded in companion files under
`specs/fontdata/apple2/`.

# 2. Base Glyph Format

## 2.1. Dimensions

Every glyph is defined on a grid of 7 dots wide by 8 dots tall, for a total of
56 dots per glyph.

## 2.2. Dot Values

Each dot is either on (foreground) or off (background). In the bitmap data, a
value of `1` means on and `0` means off. In the companion font data files, `#`
represents on and `.` represents off.

## 2.3. Storage Order

Dots are stored in row-major order: the first 7 values are the top row (row
0), the next 7 values are row 1, and so on through row 7. Within each row,
values are ordered left to right (column 0 through column 6).

## 2.4. Border Convention

All glyphs in the Apple //e font share a common border structure. Column 0
(the leftmost dot) is off in every row, with the sole exception of `{`, whose
row 3 extends into column 0. Similarly, column 6 (the rightmost dot) is off in
every row, with the sole exception of `}`, whose row 3 extends into column 6.
This provides a 1-dot gap between adjacent characters when rendered. Most
glyphs also have row 7 (the bottom row) entirely off, providing vertical
separation. The glyphs that use row 7 are the lowercase descenders (g, j, p,
q, y) and the characters !, |, and _.

# 3. Glyph Sets

The font is organized into four glyph sets of 32 characters each, covering 128
unique glyph shapes. The remaining 128 entries in the 256-entry character
space are filled by reusing the uppercase, special, and lowercase sets with
different display modes applied.

## 3.1. Uppercase Set

32 glyphs covering ASCII `$40`-`$5F`:

| Index | Character |
|-------|-----------|
| 0     | @         |
| 1-26  | A-Z       |
| 27    | [         |
| 28    | \         |
| 29    | ]         |
| 30    | ^         |
| 31    | _         |

The precise bitmaps are recorded in `specs/fontdata/apple2/uppercase.txt`.

## 3.2. Special/Punctuation Set

32 glyphs covering ASCII `$20`-`$3F`:

| Index | Character |
|-------|-----------|
| 0     | SP (space)|
| 1     | !         |
| 2     | "         |
| 3     | #         |
| 4     | $         |
| 5     | %         |
| 6     | &         |
| 7     | '         |
| 8     | (         |
| 9     | )         |
| 10    | *         |
| 11    | +         |
| 12    | ,         |
| 13    | -         |
| 14    | .         |
| 15    | /         |
| 16-25 | 0-9       |
| 26    | :         |
| 27    | ;         |
| 28    | <         |
| 29    | =         |
| 30    | >         |
| 31    | ?         |

The precise bitmaps are recorded in `specs/fontdata/apple2/special.txt`.

## 3.3. Lowercase Set

32 glyphs covering ASCII `$60`-`$7F`:

| Index | Character |
|-------|-----------|
| 0     | `         |
| 1-26  | a-z       |
| 27    | {         |
| 28    | \|        |
| 29    | }         |
| 30    | ~         |
| 31    | DEL       |

The DEL glyph is entirely blank (all dots off).

The precise bitmaps are recorded in `specs/fontdata/apple2/lowercase.txt`.

## 3.4. MouseText Set

32 unique graphic glyphs used by the Apple //e Enhanced when the alternate
character set is active (soft switch `$C00F`, state `DisplayAltChar = true`).
They occupy character codes `$40`-`$5F` in display memory.

When the primary character set is active (`DisplayAltChar = false`), that same
`$40`-`$5F` byte range uses the uppercase set with flash display mode instead
(see section 4.1).

MouseText glyphs are rendered in normal display mode (white on black). Unlike
the uppercase, special, and lowercase sets, the MouseText set is not reused
elsewhere in the character space.

The glyph index within the set is the low 5 bits of the byte value. For
example, byte `$4C` selects MouseText glyph index 12.

The precise bitmaps are recorded in `specs/fontdata/apple2/mousetext.txt`.
The header for each glyph uses its hex code (e.g., `--- 0x4C ---`) rather than
a printable character name.

# 4. Character Space Mapping

## 4.1. The 256-Entry Font

A complete Apple //e text font has 256 entries, one for each possible byte
value that can appear in display memory. Each entry is built from one of the
three glyph sets (section 3) plus an optional display-mode transformation.

| Byte Range  | Glyph Source (primary charset) | Glyph Source (alt charset) | Display Mode         |
|-------------|-------------------------------|---------------------------|----------------------|
| `$00`-`$1F` | Uppercase                     | Uppercase                 | Inverse              |
| `$20`-`$3F` | Special                       | Special                   | Inverse              |
| `$40`-`$5F` | Uppercase                     | MouseText                 | Flash / Normal (*)   |
| `$60`-`$7F` | Special                       | Special                   | Flash                |
| `$80`-`$9F` | Uppercase                     | Uppercase                 | Normal               |
| `$A0`-`$BF` | Special                       | Special                   | Normal               |
| `$C0`-`$DF` | Uppercase                     | Uppercase                 | Normal               |
| `$E0`-`$FF` | Lowercase                     | Lowercase                 | Normal               |

(*) In the primary character set, `$40`-`$5F` uses the uppercase set with flash
display mode. In the alternate character set, `$40`-`$5F` uses the MouseText
set with normal display mode.

Within each range, the glyph index is the low 5 bits of the byte value. For
example, byte `$C1` uses glyph index 1 from the uppercase set, which is `A`.

The active character set is controlled by soft switches `$C00E` (off, primary)
and `$C00F` (on, alternate), reflected in the `DisplayAltChar` state flag.

## 4.2. Display Modes

### 4.2.1. Normal

The glyph is rendered as-is: on-dots become white foreground pixels, off-dots
become black background pixels.

### 4.2.2. Inverse

Every dot in the glyph is flipped: on becomes off and off becomes on. This
produces a white field with a black character shape, giving the classic
"highlighted" appearance. The inversion is applied by XORing each dot value
with 1.

### 4.2.3. Flash

Flashing characters alternate between their normal and inverse renderings at
approximately 1.9 Hz (toggling roughly every 16 vertical blank periods). The
flash state is global -- all flashing characters on screen toggle in unison.

Note: the current implementation renders flash characters as static inverse.
See section 7 for details on the unimplemented flash behavior.

# 5. Scaling for Display Modes

The base 7x8 glyphs must be scaled to fit the 560x384 framebuffer. The scaling
is different for 40-column and 80-column text modes but is applied once at
font construction time, not at render time.

## 5.1. 40-Column Scaling

Each base glyph is doubled in both dimensions, producing a 14-pixel-wide by
16-pixel-tall rendered glyph. The doubling works as follows:

```
for each row (0 through 7):
    for pass in (0, 1):           // emit the row twice
        for each col (0 through 6):
            dot = base[row * 7 + col]
            emit dot, dot         // emit the column twice
```

This produces 40 columns x 14 pixels = 560 pixels wide and 24 rows x 16
pixels = 384 pixels tall, exactly filling the framebuffer.

## 5.2. 80-Column Scaling

Each base glyph is doubled in height only, producing a 7-pixel-wide by
16-pixel-tall rendered glyph. The doubling works as follows:

```
for each row (0 through 7):
    for pass in (0, 1):           // emit the row twice
        for each col (0 through 6):
            dot = base[row * 7 + col]
            emit dot              // no horizontal doubling
```

This produces 80 columns x 7 pixels = 560 pixels wide and 24 rows x 16
pixels = 384 pixels tall, also exactly filling the framebuffer.

## 5.3. Transformation Order

For each glyph entry in the font, the transformations are applied in this
order:

1. Start with the base 7x8 bitmap from the appropriate glyph set.
2. Apply the scaling transformation (40-column or 80-column doubling).
3. Apply the display mode mask (inverse XOR), if applicable.

The result is a pre-rendered pixel buffer that the renderer can blit directly
to the framebuffer without any per-pixel logic at render time.

# 6. Font Data Files

## 6.1. Location

The glyph bitmap data files are stored in `specs/fontdata/apple2/`:

- `uppercase.txt` -- the 32 uppercase glyphs (section 3.1)
- `special.txt` -- the 32 special/punctuation glyphs (section 3.2)
- `lowercase.txt` -- the 32 lowercase glyphs (section 3.3)
- `mousetext.txt` -- the 32 MouseText glyphs (section 3.4)

## 6.2. File Format

Each file is a plain text file containing glyph definitions in sequence. The
format is:

- Lines beginning with `#` are comments.
- Each glyph begins with a header line: `--- <name> ---` where `<name>` is the
  character or a mnemonic (e.g., `SP` for space, `DEL` for delete).
- Following the header are exactly 8 lines of 7 characters each, representing
  the 8 rows of the glyph from top to bottom.
- Within each row, `#` represents an on-dot and `.` represents an off-dot.
- A blank line separates consecutive glyphs.

## 6.3. Normative Status

The bitmap data in these files is normative. An implementation must produce
identical glyph bitmaps to pass conformance. The data was derived from
examination of the Apple //e character generator ROM and represents the
expected output of the emulator's font system.

# 7. Unimplemented: Flash Animation

Characters in the `$60`-`$7F` range should alternate between normal and
inverse display at approximately 1.9 Hz. The current implementation renders
them as static inverse. A correct implementation would maintain a global flash
state that toggles based on a VBL counter and select between the normal and
inverse glyph renderings accordingly. When the flash state changes, the redraw
flag must be set so the screen updates even if no display memory has been
written.

Note: the `$40`-`$5F` range is not subject to flash -- those positions hold
the MouseText glyphs (section 3.4), which are always rendered in normal mode.

# 8. Design Considerations

## 8.1. Pre-Rendered Glyphs

Scaling and masking are applied at font construction time rather than at
render time. The font object stores 256 pre-rendered glyph buffers (one per
possible byte value), so the render loop is a simple lookup-and-blit with no
per-pixel arithmetic. This keeps the per-frame cost proportional to the number
of character cells (960 for 40-column, 1920 for 80-column) with minimal work
per cell.

## 8.2. Shared Glyph Definitions

The uppercase, special, and lowercase glyph sets are defined once and reused
across multiple byte ranges with different masks. The uppercase set, for
example, is registered three times (at offsets `$00`, `$80`, and `$C0`) with
either the inverse mask or no mask. The MouseText set is the exception: it
occupies exactly one byte range (`$40`-`$5F`) and is not reused elsewhere.
This avoids duplicating bitmap data and ensures that all instances of a given
character shape are pixel-identical (modulo the display mode transformation).

## 8.3. Separate Font Objects

The 40-column and 80-column fonts are separate font objects with different
glyph dimensions, even though they share the same base bitmaps. This avoids
branching on the column mode inside the render loop -- the renderer simply
uses whichever font object matches the current mode.

## 8.4. External Bitmap Data

Storing the precise bitmap data in separate files under `specs/fontdata/`
rather than inline in the spec keeps the spec readable while still making the
data normative and version-controlled. The data files use a human-readable
visual format (`#` and `.`) so that glyphs can be inspected and compared
without tooling.
