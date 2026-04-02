---
Specification: 9
Category: Graphics
Drafted At: 2026-03-23
Authors:
  - Peter Evans
---

# 1. Overview

This spec describes the implementation of 40-column text mode for the Apple
//e. Text mode is the default display mode of the machine: on boot, the screen
shows 24 rows of 40 characters each, rendered in white on black using the
built-in character generator ROM. The implementation covers memory layout,
address decoding, character set encoding, glyph rendering, soft switch control,
and the end-to-end rendering pipeline.

# 2. Screen Geometry

## 2.1. Logical Dimensions

The 40-column text display is a grid of 40 columns by 24 rows, for a total of
960 character cells.

## 2.2. Pixel Dimensions

Each character cell is 7 dots wide and 8 dots tall in the base glyph
definition. For 40-column mode, each dot is doubled in both dimensions,
producing a 14-pixel-wide by 16-pixel-tall cell on screen.

The full screen resolution is 560 x 384 pixels:

- Width: 40 columns x 14 pixels = 560
- Height: 24 rows x 16 pixels = 384

This resolution matches the Apple //e's effective output when accounting for
the double-width pixels of 40-column mode. It also conveniently accommodates
80-column mode (80 x 7 = 560) in the same framebuffer.

# 3. Text Page Memory

## 3.1. Page 1

Text page 1 occupies main memory from `$0400` to `$07FF` (1024 bytes). Only
960 of these bytes correspond to visible character cells; the remaining 64
bytes are "holes" that do not map to any screen position.

## 3.2. Page 2

Text page 2 is selected via the PAGE2 soft switch. Its location depends on
the 80STORE state:

- **80STORE off**: page 2 is at `$0800`-`$0BFF` in main memory.
- **80STORE on**: page 2 is at `$0400`-`$07FF` in auxiliary memory.

The 80STORE behavior exists to support double-hires and 80-column modes, which
use the auxiliary bank for the "other half" of the display. In pure 40-column
mode without 80STORE, page 2 simply uses a different 1 KB region of main
memory.

## 3.3. Interleaved Address Mapping

Characters are not stored sequentially in memory. The Apple II uses an
interleaved addressing scheme inherited from the hardware's display generation
circuitry. The 1024-byte text page is divided into groups of 128 bytes, and
each group contains parts of three different screen rows, plus an 8-byte hole
at the end.

The mapping from memory offset (relative to `$0400`) to screen position
follows this pattern:

| Offset Range  | Row Group | Bytes per Row | Holes        |
|---------------|-----------|---------------|--------------|
| `$000`-`$027` | Row 0     | 40            |              |
| `$028`-`$04F` | Row 8     | 40            |              |
| `$050`-`$077` | Row 16    | 40            |              |
| `$078`-`$07F` | *(hole)*  |               | 8 bytes      |
| `$080`-`$0A7` | Row 1     | 40            |              |
| `$0A8`-`$0CF` | Row 9     | 40            |              |
| `$0D0`-`$0F7` | Row 17    | 40            |              |
| `$0F8`-`$0FF` | *(hole)*  |               | 8 bytes      |
| ...           | ...       |               |              |
| `$380`-`$3A7` | Row 7     | 40            |              |
| `$3A8`-`$3CF` | Row 15    | 40            |              |
| `$3D0`-`$3F7` | Row 23    | 40            |              |
| `$3F8`-`$3FF` | *(hole)*  |               | 8 bytes      |

Within each 128-byte group, three rows of 40 bytes are stored consecutively,
followed by an 8-byte hole. The rows within each group are spaced 8 rows
apart on screen (rows 0/8/16, then 1/9/17, then 2/10/18, etc.).

More precisely, for a given 128-byte group index `g` (0-7):

- Bytes `0`-`39`: screen row `g`
- Bytes `40`-`79`: screen row `g + 8`
- Bytes `80`-`119`: screen row `g + 16`
- Bytes `120`-`127`: hole (not displayed)

The column within a row is simply the byte position modulo 40.

## 3.4. Base Address Formula

The base address of a screen row can be computed as:

```
base(row) = $0400 + (row % 8) * 128 + (row / 8) * 40
```

The address of the character at column `col` in row `row` is:

```
addr(row, col) = base(row) + col
```

## 3.5. Address Lookup Tables

Rather than computing the row/column from an address at runtime, the
implementation uses a pair of precomputed lookup tables -- one for rows and one
for columns -- indexed by the memory offset (`addr - $0400`). Each table has
1024 entries. Entries corresponding to hole bytes contain the sentinel value
`-1`, which signals the renderer to skip that address.

# 4. Character Encoding

## 4.1. Character Byte Layout

Each byte in the text page encodes both the character identity and its display
mode. The 256 possible byte values are divided into four ranges of 64:

| Byte Range    | Display Mode | Characters               |
|---------------|--------------|--------------------------|
| `$00`-`$1F`   | Inverse      | Uppercase letters (@-_)  |
| `$20`-`$3F`   | Inverse      | Special/punctuation      |
| `$40`-`$5F`   | Flash        | Uppercase letters (@-_)  |
| `$60`-`$7F`   | Flash        | Special/punctuation      |
| `$80`-`$9F`   | Normal       | Uppercase letters (@-_)  |
| `$A0`-`$BF`   | Normal       | Special/punctuation      |
| `$C0`-`$DF`   | Normal       | Uppercase letters (@-_)  |
| `$E0`-`$FF`   | Normal       | Lowercase letters        |

The "special/punctuation" set covers the 32 characters from space (`$20`)
through `?` (`$3F`) in ASCII: space, `!`, `"`, `#`, `$`, `%`, `&`, `'`, `(`,
`)`, `*`, `+`, `,`, `-`, `.`, `/`, `0`-`9`, `:`, `;`, `<`, `=`, `>`, `?`.

The "uppercase letters" set covers the 32 characters from `@` (`$40`) through
`_` (`$5F`) in ASCII: `@`, `A`-`Z`, `[`, `\`, `]`, `^`, `_`.

## 4.2. Display Modes

- **Normal**: white foreground on black background (the default for most text).
- **Inverse**: black foreground on white background (the glyph dots are
  XOR-flipped).
- **Flash**: the character alternates between normal and inverse at a rate tied
  to the vertical blank period. On real hardware, the flash rate is
  approximately 1.9 Hz (toggling every ~16 NTSC scans at ~60 Hz).

## 4.3. Alternate Character Set

The ALTCHAR soft switch selects between the primary and alternate character
sets. In the primary set, bytes `$40`-`$7F` display as flashing characters. In
the alternate character set, those same byte values display as normal
(non-flashing) characters with the MouseText glyphs or, on the enhanced //e,
with additional symbols. The alternate character set is out of scope for this
spec.

# 5. Glyph Definitions

## 5.1. Base Glyph Format

Each character glyph is defined as a 7x8 grid of dots. A dot value of `1`
means the pixel is "on" (foreground); `0` means "off" (background). The glyph
data is stored as a flat byte slice of length 56 (7 * 8), row-major.

## 5.2. Character Sets

The glyphs are organized into three definition groups:

- **Uppercase**: 32 glyphs for `@`, `A`-`Z`, `[`, `\`, `]`, `^`, `_`
- **Special**: 32 glyphs for space, `!`-`/`, `0`-`9`, `:`-`?`
- **Lowercase**: 32 glyphs for `` ` ``, `a`-`z`, `{`, `|`, `}`, `~`, DEL

Each group is a function that defines glyphs at a given offset in the font,
applying an optional mask function.

## 5.3. Doubling for 40-Column Display

The framebuffer is 560 pixels wide to accommodate both 40- and 80-column
modes. In 80-column mode, each 7-dot-wide glyph maps directly to 7 pixels
(80 x 7 = 560). In 40-column mode, each character must span 14 pixels
(40 x 14 = 560), so the 7x8 base glyphs are doubled. A doubling function
expands each glyph by replicating every dot into a 2x2 block:

```
for each row (0-7):
    for pass in (0, 1):          // duplicate the row
        for each col (0-6):
            dot = base[row * 7 + col]
            output two copies of dot  // duplicate the column
```

The result is a 14x16 byte array suitable for the 40-column font.

## 5.4. Inverse Mask

For inverse characters (bytes `$00`-`$3F`), a mask function is applied after
doubling. The mask XORs every byte in the glyph data with 1, flipping `0` to
`1` and `1` to `0`. This swaps foreground and background, producing the
inverse video effect.

## 5.5. Font Construction

The 40-column system font is assembled by defining all 256 character entries:

| Offset | Glyph Source | Mask    |
|--------|-------------|---------|
| `$00`  | Uppercase   | Inverse |
| `$20`  | Special     | Inverse |
| `$40`  | Uppercase   | Inverse |
| `$60`  | Special     | Inverse |
| `$80`  | Uppercase   | None    |
| `$A0`  | Special     | None    |
| `$C0`  | Uppercase   | None    |
| `$E0`  | Lowercase   | None    |

In the primary font, bytes `$40`-`$7F` are rendered as inverse. The flash
behavior that alternates these between inverse and normal display is described
in section 8.

# 6. Soft Switches

## 6.1. Mode Selection

The following soft switches control whether 40-column text mode is active and
which page is displayed:

| Switch      | Address       | Access | Effect                          |
|-------------|---------------|--------|---------------------------------|
| TEXT on     | `$C051`       | R/W    | Enable text mode                |
| TEXT off    | `$C050`       | R/W    | Disable text mode (graphics)    |
| PAGE2 on   | `$C055`       | R/W    | Select text page 2              |
| PAGE2 off  | `$C054`       | R/W    | Select text page 1              |
| MIXED on   | `$C053`       | R/W    | Enable mixed mode (4 text rows) |
| MIXED off  | `$C052`       | R/W    | Disable mixed mode              |
| 80COL on   | `$C00D`       | W      | Enable 80-column mode           |
| 80COL off  | `$C00C`       | W      | Disable 80-column mode          |
| 80STORE on | `$C001`       | W      | Enable 80STORE memory mapping   |
| 80STORE off| `$C000`       | W      | Disable 80STORE memory mapping  |
| ALTCHAR on | `$C00F`       | W      | Select alternate character set  |
| ALTCHAR off| `$C00E`       | W      | Select primary character set    |

For 40-column text mode specifically, TEXT must be on and 80COL must be off.
When TEXT is on, the display renders from the text page regardless of the HIRES
setting.

## 6.2. Status Reads

Programs can query the current display state by reading the following
addresses. Bit 7 of the returned byte is set (value `$80`) when the mode is
active, and clear (value `$00`) when inactive:

| Read Address | Returns Bit 7 High When |
|-------------|-------------------------|
| `$C01A`     | TEXT is on               |
| `$C01B`     | MIXED is on              |
| `$C01C`     | PAGE2 is on              |
| `$C01E`     | ALTCHAR is on            |
| `$C01F`     | 80COL is on              |
| `$C018`     | 80STORE is on            |

## 6.3. Vertical Blank

The VBL (vertical blank) status is readable at `$C019`. It returns bit 7 high
during the vertical blank period. The full scan cycle takes 17,030 CPU cycles:
12,480 cycles for the active display and 4,550 cycles for the vertical blank
retrace. Programs use VBL timing to synchronize screen updates and avoid
flicker.

## 6.4. Redraw Flag

Every soft switch write that changes display state sets a "redraw" flag. The
rendering pipeline checks this flag each frame and skips rendering entirely if
nothing has changed. Writes to display memory (`$0400`-`$0BFF`) also set this
flag.

# 7. Rendering Pipeline

## 7.1. Overview

The rendering pipeline converts the contents of text page memory into pixels
on screen. It runs once per frame when the redraw flag is set, and consists of
these stages:

1. Snapshot display memory
2. Select rendering mode
3. Iterate over the text page
4. Map each address to a screen position
5. Look up the glyph
6. Blit the glyph to the framebuffer
7. Present the framebuffer

## 7.2. Display Memory Snapshot

Before rendering, the current contents of display memory are copied into a
snapshot buffer. This prevents tearing that would occur if the CPU modifies
display memory while the renderer is reading it. The snapshot copies the
appropriate 1 KB region based on the current page and 80STORE settings:

- **80STORE off, PAGE2 off**: copy `$0400`-`$07FF` from main memory
- **80STORE on, PAGE2 on**: copy `$0400`-`$07FF` from auxiliary memory
- **80STORE off, PAGE2 on**: copy `$0800`-`$0BFF` from main memory

## 7.3. Mode Dispatch

The renderer checks the display state flags to determine which mode to use.
When the TEXT flag is set, 40-column text rendering is selected. The mode
priority is:

1. If TEXT is on: render text
2. If HIRES is on: render hi-res graphics
3. Otherwise: render lo-res graphics

Mixed mode (MIXED on with TEXT off) is a variation where the bottom 4 rows of
the screen display text while the upper portion displays graphics. This is
handled by the graphics renderers, not the text renderer.

## 7.4. Character Rendering Loop

The text renderer iterates over every address in the 1024-byte text page
(offsets `$000` through `$3FF`):

```
for offset = 0 to 1023:
    row = addressRows[offset]
    col = addressCols[offset]

    if row == -1 or col == -1:
        continue    // hole byte, skip

    char = snapshot.get($0400 + offset)
    glyph = font.glyph(char)

    x = col * 14   // glyphWidth for 40-col
    y = row * 16   // glyphHeight for 40-col

    screen.blit(x, y, glyph)
```

## 7.5. Monochrome Modes

The renderer supports optional monochrome color modes that recolor the display
to simulate period-accurate monitors:

- **Green screen**: foreground pixels are recolored to RGB `(152, 255, 152)`
- **Amber screen**: foreground pixels are recolored to RGB `(255, 191, 0)`

When a monochrome mode is active, each glyph is post-processed: any white
pixel (`255, 255, 255`) is replaced with the monochrome color. Black pixels
are left unchanged.

## 7.6. Framebuffer Output

The framebuffer is a 560 x 384 array of RGBA pixels. After all glyphs are
blitted, the framebuffer's pixel data is written to an Ebiten image for
display. An optional CRT shader may be applied during this final render step
to simulate the appearance of a CRT monitor (scanlines, curvature).

# 8. Flash

Characters in the `$40`-`$7F` range alternate between normal and inverse
display. On real hardware, the flash state toggles every ~16 VBL periods,
producing a rate of approximately 1.9 Hz. Flash only applies when the primary
character set is active; the alternate character set replaces these byte
values with MouseText glyphs that display without flashing.

## 8.1. Flash State

A boolean flash state tracks whether flash characters are currently shown in
their normal or inverse rendering. This state is derived from the CPU cycle
counter: divide the cycle counter by the scan cycle count (17,030) to obtain
an NTSC scan number (~60 per second), then divide by 16 to obtain a flash
phase. The flash state is the parity of the flash phase:

```
frameNumber = cycleCounter / ScanCycleCount
flashPhase  = frameNumber / 16
flashOn     = (flashPhase % 2) == 0
```

When `flashOn` is true, flash characters render as inverse (the current
behavior for all frames). When `flashOn` is false, flash characters render as
normal.

The cycle counter is the existing `CPU.CycleCounter()` value already used for
VBL timing. No additional counter or state key is needed.

## 8.2. Dual-Font Approach

The font system already builds pre-rendered glyphs for every character code at
font creation time. To support flash, the text renderer uses two fonts rather
than one: the existing primary font (where `$40`-`$7F` are inverse) and a
second "flash-alternate" font (where `$40`-`$7F` are normal). The
flash-alternate font is identical to the primary font in all other ranges.

The flash-alternate 40-column font is constructed exactly like `SystemFont40`,
except the `$40`-`$5F` range uses uppercase glyphs with no mask (normal
rendering) and the `$60`-`$7F` range uses special/punctuation glyphs with no
mask:

| Offset | Glyph Source | Mask (primary) | Mask (flash-alternate) |
|--------|-------------|----------------|------------------------|
| `$40`  | Uppercase   | Inverse        | None                   |
| `$60`  | Special     | Inverse        | None                   |

All other offsets are identical between the two fonts.

The same approach applies to the 80-column primary font.

The alternate character set fonts (`SystemFont40Alt`, `SystemFont80Alt`) are
not affected; they do not display flash characters.

## 8.3. Font Selection in the Render Path

The text renderer (`a2text.Render`) accepts two font arguments: the primary
font and the flash-alternate font. It also accepts the current flash state,
which is computed by `Computer.Render` as described in section 8.4. If
`flashOn` is true, it uses the primary font (flash characters appear inverse).
If `flashOn` is false, it uses the flash-alternate font (flash characters
appear normal).

This is a whole-font swap per frame, not a per-character branch. The cost is
one comparison before the loop rather than a conditional inside the loop for
every character.

## 8.4. Triggering Redraws

Flash state changes must cause the screen to redraw even when no memory has
been written. The `Computer.Render` method computes the current flash state
before checking the redraw flag. If the flash state differs from the
previously stored value, the redraw flag is set. This ensures the screen
updates at the flash toggle points without unnecessary redraws between them.

The previous flash state is stored as a boolean in the state map using a
`DisplayFlash` key.

# 9. Design Considerations

## 9.1. Lookup Tables vs. Computation

The address-to-position mapping uses precomputed lookup tables rather than
computing `row = ...` and `col = ...` from the address at render time. While
the formula in section 3.4 is straightforward, the lookup tables are a single
indexed read per address with no division or modulo operations. Given that
this code runs for every byte in the text page on every frame redraw, the
table approach is both simpler to verify and faster to execute.

## 9.2. Glyph Doubling at Font Creation Time

The 7x8 to 14x16 doubling is performed once when the font is created, not on
every render. Each character's doubled glyph is stored as a pre-built
framebuffer in the font's glyph map. At render time, the renderer simply looks
up the framebuffer by character code and blits it, with no per-pixel scaling
logic in the hot path.

## 9.3. Snapshot-Based Rendering

Copying display memory into a snapshot before rendering adds a small cost
(copying 1 KB) but eliminates an entire class of visual artifacts. Without the
snapshot, the CPU could write to display memory between two glyph blits in the
same frame, causing the top and bottom halves of the screen to show different
states. The snapshot guarantees a consistent view for the duration of the
render.

## 9.4. Flash via Font Swap

Flash could be implemented by checking each character's byte range inside the
rendering loop and selecting the appropriate glyph. However, this adds a
conditional to the inner loop for every character on every frame. The
dual-font approach moves the decision out of the loop entirely: one comparison
before the loop selects which font to use, and the loop body remains
unchanged. The trade-off is storing a second font (an additional 256 glyph
framebuffers), but since only the `$40`-`$7F` range differs, the memory cost
is modest and the rendering path needs no per-character flash check.

## 9.5. Shared Framebuffer with Other Modes

The 560 x 384 framebuffer is shared across all display modes (text, lo-res,
hi-res, double hi-res). The 40-column text renderer fills the entire
framebuffer since 40 columns x 14 pixels = 560 and 24 rows x 16 pixels = 384.
No padding or centering is needed.
