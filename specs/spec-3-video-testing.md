---
Specification: 3
Drafted At: 2026-03-06
Authors:
  - Peter Evans
---

# 1. Overview

This spec describes a system for recording and testing video output during
emulator execution. Where spec 1 (state testing) captures discrete transitions
-- a byte changing, a register updating -- and spec 2 (audio testing) captures
a continuous signal summarized by fingerprints, video occupies a middle
ground. The screen is a large piece of state (hundreds of thousands of
pixels); what matters for testing is the complete picture at a given moment.

This system captures the rendered framebuffer at specific execution steps and
compares it against expected content using exact frame matching. A video
assertion says: "at step N, the screen looks precisely like this."

# 2. Concepts

## 2.1. Video Recorder

A video recorder captures snapshots of the rendered framebuffer at specified
execution steps. It sits downstream of the display rendering pipeline -- after
all mode-specific rendering has produced the final image in `gfx.Screen`. The
recorder captures what the user would see, not the underlying memory layout.

The video recorder operates alongside the state recorder from spec 1 and the
audio recorder from spec 2. All three share the same step counter as their
timeline.

## 2.2. Frame Snapshot

A frame snapshot is a capture of the screen content at a single execution
step. It contains the rendered pixel grid stored as RGBA color values -- the
same data held by `gfx.Screen`. A snapshot is taken after the instruction at
the given step has executed and after the display has been re-rendered.

## 2.3. Color Legend

A color legend maps single-character codes to specific RGB color values. Each
video assertion defines its own legend, declaring the colors that appear in
the expected grid. This keeps the spec independent of any particular machine's
palette -- the test author picks characters that are meaningful for the colors
in their test case.

A legend entry maps one character to one RGB hex value:

```
colors: . = 000000, # = FFFFFF, P = D043E5, g = 2FBC1A
```

The character `.` conventionally represents black, and `#` conventionally
represents white, but these are conventions, not requirements. The test author
is free to use any printable non-whitespace character.

During comparison, each pixel in the captured frame is matched against the
legend by exact RGB value (alpha is ignored). If a captured pixel's color does
not appear in the legend, the comparison fails with an "unmapped color" error,
indicating that the legend is incomplete for the actual frame content.

## 2.4. Sampling Resolution

The framebuffer (`gfx.Screen`) has fixed pixel dimensions that may be larger
than the emulated machine's logical output. For example, a framebuffer might
render at 560x384 to accommodate double-width modes, even when the current
display content only uses 280x192 distinct values.

A video assertion specifies the grid dimensions at which the comparison
operates. The recorder samples the framebuffer to produce a grid of the
requested size. For a grid of width W and height H over a framebuffer of
width FW and height FH, each grid cell (gx, gy) is sampled from the
framebuffer pixel at:

```
fx = gx * FW / W
fy = gy * FH / H
```

This integer-scaled sampling works cleanly when the grid dimensions evenly
divide the framebuffer dimensions (e.g. 280x192 over 560x384 samples every
other pixel). Assertions should use grid dimensions that divide evenly to
avoid sampling artifacts.

## 2.5. Video Assertion

A video assertion declares the expected content of the screen (or a region of
it) at a specific execution step. Unlike audio assertions which use
tolerances, video assertions are exact -- every cell in the expected grid must
match the captured frame.

# 3. Recording

## 3.1. Capture Point

The video recorder captures from `gfx.Screen` -- the same framebuffer that
the rendering engine displays. This is after all rendering has been applied.
The recorder sees the final rendered image.

The recorder does not capture overlays (volume indicators, status text, pause
notifications, etc.) since those are drawn on top of `gfx.Screen` during the
display engine's draw phase and are not part of the framebuffer data. It also
does not capture shader effects (CRT simulation, scanlines) since those are
cosmetic post-processing. Both are user-facing presentation concerns that do
not reflect the emulated machine's output.

## 3.2. Capture Timing

The recorder maintains a set of step numbers at which to capture. These are
declared by the test before execution begins. At each step in the set:

1. The emulator executes the instruction.
2. The step counter increments.
3. The display re-renders (the emulator's `Render()` method is called).
4. The recorder reads `gfx.Screen` and stores a copy of the pixel data.

Captures happen only at the declared steps, not at every step. This avoids the
cost of snapshotting the framebuffer on every instruction.

# 4. Assertion Format

## 4.1. Full Screen Assertion

A full screen assertion specifies a step, a grid size, a color legend, and the
expected pixel grid. The header line declares the step and grid dimensions.
The color legend follows. Then the grid data appears as one line per row, with
each character representing one pixel using the legend's codes.

```
step 512: video screen 280x192
colors: . = 000000, # = FFFFFF
............................................................
............................................................
............................................................
..............................##............................
.............................####...........................
............................##..##..........................
...........................##....##.........................
..........................##########........................
.........................##........##.......................
........................##..........##......................
............................................................
............................................................
```

Each row must contain exactly as many characters as the grid width. The number
of rows must equal the grid height.

## 4.2. Region Assertion

A region assertion captures a rectangular portion of the screen. This is the
common case -- most tests do not need to verify the entire screen.

```
step 512: video region 10,5 20x10 280x192
colors: . = 000000, # = FFFFFF
....................
..........##........
.........####.......
........##..##......
.......##....##.....
......##########....
.....##........##...
....##..........##..
....................
....................
```

The format is `video region <x>,<y> <w>x<h> <grid_w>x<grid_h>` where
`<x>,<y>` is the top-left corner in grid coordinates, `<w>x<h>` is the region
size, and `<grid_w>x<grid_h>` is the sampling resolution for the full screen.
The region coordinates and dimensions are in the grid's coordinate space, not
the framebuffer's. The grid data must contain exactly `<h>` rows of exactly
`<w>` characters each.

## 4.3. Row Assertion

For cases where only a few rows matter, a row assertion avoids reproducing the
entire grid:

```
step 512: video row 96 280x192
colors: . = 000000, # = FFFFFF
............................##..##..........................
```

This asserts that row 96 (at the specified sampling resolution) matches the
given content. Multiple row assertions can appear for the same step.

# 5. Comparison

## 5.1. Cell-by-Cell Matching

The recorded frame is sampled at the assertion's grid resolution, and each
sampled pixel's RGB value is looked up in the color legend to produce a
character. The resulting character grid is compared cell by cell against the
expected grid, left to right, top to bottom. A test passes if and only if
every cell matches. A single mismatch causes a failure.

If a sampled pixel's RGB value does not match any entry in the legend, the
comparison fails immediately with an "unmapped color" error. This catches
cases where the legend is missing a color that actually appears in the
captured frame.

## 5.2. Failure Reporting

When a comparison fails, the report includes:

- The step number
- The coordinates of the first mismatched cell
- The expected and actual color codes at that cell
- The actual RGB value of the mismatched pixel
- A rendering of the mismatched row showing the difference

For example:

```
step 512: video mismatch at (14, 7)
  expected: ......##########........
  actual:   ......###.######........
                     ^ (14, 7): expected '#' (FFFFFF), got '.' (000000)
```

## 5.3. Diff Output

For larger mismatches, the recorder can produce a full diff of the expected
and actual grids. Rows that match are omitted; rows that differ are printed
with markers highlighting the differing cells.

# 6. Testing Model

## 6.1. Video Test Structure

A video test consists of:

1. A program or disk image to load into the emulator
2. A set of state observers to register (from spec 1, if needed)
3. An audio recorder attached to the audio stream (from spec 2, if needed)
4. A video recorder with a set of capture steps declared
5. A number of steps to execute
6. A list of video assertions, each specifying a step and expected content

The video recorder begins in an idle state and only captures at the declared
steps. There is no continuous recording -- unlike audio, video frames are
captured on demand.

## 6.2. Combining with State and Audio Tests

Video assertions, state entry assertions (from spec 1), and audio assertions
(from spec 2) can coexist in the same test. All three systems share the same
step counter and the same execution run. This allows a test to verify, for
example, that a program writes a value to video memory, produces a tone, *and*
renders the expected image on screen -- all at the same step.

## 6.3. Test-Only Concern

As with the state and audio recording systems, the video recorder is only
instantiated during testing. Production execution paths do not create a video
recorder and incur no overhead. The recorder reads from the existing
`gfx.Screen` framebuffer and does not require changes to the rendering code.

# 7. Raw Frame Capture

## 7.1. Debugging Failed Tests

When a video assertion fails, the full captured frame can be written to a file
for inspection. The recorder supports two output formats:

- **Text grid** -- a character grid using the assertion's color legend, written
  to a `.frame` file for easy diffing
- **PNG image** -- the raw framebuffer pixels written as a PNG file for visual
  inspection in an image viewer

Both outputs are produced for the failing step, allowing the developer to see
exactly what the emulator rendered.

## 7.2. Frame Capture Without Assertions

The recorder can also capture frames at specified steps without any
assertions. This is useful for generating initial expected data -- run the
emulator, capture the frames, inspect them visually, and then paste the grid
data into the test as the expected content. The recorder can also emit a color
legend derived from the distinct colors present in the captured frame.

# 8. Design Considerations

## 8.1. Why Not Per-Pixel Transitions

Spec 1's observer model watches individual state items and records
transitions. Applying this to video -- watching every pixel in the framebuffer
as an individual observer -- would produce an overwhelming volume of entries.
A single frame of graphics involves thousands of pixel changes. The
transitions themselves are not informative; what matters is the resulting
image. Frame snapshots capture that result directly.

## 8.2. Why Not Fingerprinting

Spec 2 uses fingerprints for audio because minor sample-level variations do
not affect the audible result. Video does not have this property -- a single
wrong pixel is a rendering bug. Given the same program state, the rendered
output is deterministic. Exact matching is both feasible and appropriate.

## 8.3. Machine Independence

The spec does not define a fixed color palette or a fixed set of display
modes. The color legend is declared per assertion, so any set of colors can be
represented. The sampling resolution is an explicit parameter, so any
framebuffer geometry can be tested. A test for one machine (e.g. a 280x192
display with 16 colors) uses the same assertion format as a test for another
(e.g. a 320x200 display with 256 colors) -- only the legend and grid
dimensions differ.

## 8.4. Region Assertions as the Common Case

Full screen assertions are supported but uncommon in practice. Most tests
focus on a specific area of the screen -- a line of text, a sprite, a status
indicator. Region assertions keep the test data small and focused, making
failures easier to read and expected data easier to author.

## 8.5. Step Counting

As with specs 1 and 2, the step counter counts instructions executed, not
clock cycles. The video recorder uses the same counter maintained by the
state recorder. Captures happen after the instruction at the given step has
completed and the display has been re-rendered, ensuring the frame reflects
the state produced by that instruction.
