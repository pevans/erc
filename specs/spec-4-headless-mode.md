---
Specification: 4
Category: Tests
Drafted At: 2026-03-08
Authors:
  - Peter Evans
---

# 1. Overview

This spec describes a headless execution mode for erc. Headless mode runs the
emulator without a graphical window, executing a fixed number of instruction
steps and recording state, audio, and video as configured. This enables
automated testing using the record package (specs 1-3) without requiring a
display server or user interaction.

# 2. Problem

Two categories of ebiten dependency block headless execution:

## 2.1. FrameBuffer Constructor

`gfx.NewFrameBuffer()` calls `ebiten.NewImage()` at construction time.
`a2.NewComputer()` creates a FrameBuffer, so any code that instantiates a
Computer -- including test code -- triggers an ebiten dependency. However, the
pixel storage (`pixels` slice, `SetCell`, `GetPixel`, `Blit`, `ClearCells`,
`Invert`, `Pixels`) has no ebiten dependency. Only `Image`, `Render()`, and
`SetShader()` use ebiten types.

## 2.2. Package-Level init() Functions

Two `init()` functions in the `gfx` package create ebiten resources at import
time:

- `gfx/border.go:init()` creates a `BorderOverlay` whose constructor calls
  `ebiten.NewImage(1, 1)`.
- `gfx/textoverlay.go:init()` creates a `TextOverlay` whose constructor calls
  `text.NewGoTextFaceSource()` from ebiten's text package.

A third init in `gfx/status.go` creates an `Overlay` via `NewOverlay()`, which
does not call any ebiten functions and is already safe.

## 2.3. What Is Already Headless-Compatible

Everything else works without ebiten:

- The execution loop (`clock.Emulator.ProcessLoop`)
- Audio (`a2audio.Stream.Read()`)
- The record package (recorder, observers, audio/video recorders)
- Core emulation (`a2.Process`, `a2.Boot`, `a2.Load`)
- Display renderers (`a2hires`, `a2text`, etc.) write to the `gfx.Screen`
  pixel buffer, which is just a byte slice

# 3. Lazy Ebiten Resource Creation

## 3.1. FrameBuffer

Remove the `ebiten.NewImage()` call from `NewFrameBuffer()`. The `Image` field
should be created lazily on the first call to `Render()`. This is safe because:

- All pixel-manipulation methods (`SetCell`, `GetPixel`, `Blit`, `ClearCells`,
  `Invert`, `Pixels`) only touch the `pixels` byte slice.
- `Render()` is only called from ebiten's draw loop.
- `SetShader()` compiles a shader but does not use `Image`.

## 3.2. BorderOverlay

Remove the `ebiten.NewImage(1, 1)` call from `NewBorderOverlay()`. The `pixel`
field should be created lazily on the first call to `Draw()`. The `pixel` is
only used in `Draw()` to render scaled rectangles, so lazy creation there is
natural.

## 3.3. TextOverlay

Remove the `text.NewGoTextFaceSource()` call from `NewTextOverlay()`. The
`faceSource` field should be created lazily on the first call to `Draw()`. The
font source is only used in `Draw()` to construct a `GoTextFace` for text
rendering.

Note: `a2.Computer.ShowText()` calls `gfx.TextNotification.Show()`
unconditionally. In headless mode this is harmless -- `Show()` only sets fields
on the struct. No ebiten call happens until `Draw()`, which is never called in
headless mode.

# 4. Headless Execution Loop

## 4.1. Run Function

The record package should provide a `Run()` function that drives the emulator
step-by-step at full speed (no wall-clock timing). It accepts configuration
specifying the computer, recorder, optional audio/video recorders, step count,
and an optional render function.

## 4.2. Per-Step Behavior

For each step:

1. The recorder's `Step()` method executes one instruction and records state
   changes (per spec 1).
2. If a video recorder is attached and the current step is in its capture set,
   the render function is called to update the framebuffer before the video
   recorder's `Observe()` captures it.

## 4.3. Render Function

The render function is provided by the caller -- typically
`comp.Render()`. This updates the framebuffer pixel data from display memory
without requiring ebiten. Video capture only works at declared steps, so the
render function is called sparingly.

## 4.4. Audio

The caller creates an `a2audio.Stream` without an ebiten audio player and
passes it to the `AudioRecorder` as an `io.Reader`. The stream generates PCM
samples from speaker toggle events synchronized to CPU cycles, which works
entirely without ebiten.

## 4.5. Audio and Full-Speed Mode

The audio stream discards speaker events and outputs silence whenever the
emulator is in full-speed mode (e.g. during disk I/O). This behavior is the
same in headless mode as in normal graphical execution -- full-speed periods
produce no meaningful audio. Audio tests should target step ranges after boot
and disk loading have completed, when the emulator is running at normal clock
speed.

# 5. Supporting Methods

## 5.1. Recorder.CurrentStep()

The `Recorder` should expose a `CurrentStep()` method returning the current
step count. The `step` field already exists but is unexported.

## 5.2. VideoRecorder.NeedsCapture()

The `VideoRecorder` should expose a `NeedsCapture(step)` method that returns
whether the given step is in its capture set. This is a simple lookup in the
existing `captureSteps` map.

# 6. CLI Subcommand

## 6.1. Command

A new `erc headless` subcommand accepts one or more disk image arguments and
the following flags:

- `--steps N` (required) -- number of instruction steps to execute
- `--watch-mem RANGES` -- comma-separated memory address ranges to observe
  (e.g. `0400-07FF,2000-3FFF` or `013F`)
- `--watch-reg LIST` -- comma-separated registers to observe (e.g. `A,X,P,PC`)
- `--watch-comp LIST` -- comma-separated computer state names to observe (e.g.
  `DisplayHires,BankWriteRAM`)
- `--record-audio` -- attach an audio recorder
- `--capture-video STEPS` -- comma-separated step numbers at which to capture
  video frames
- `--output DIR` -- directory for output files (defaults to current directory)
- `--start-at ADDR` -- hex address at which to begin counting steps (see 6.3)
- `--monochrome MODE` -- render in monochrome (`green` or `amber`); applies to
  video captures

## 6.2. Execution Flow

1. Create an `a2.Computer` (headless-safe after lazy init changes).
2. Load and boot the disk image(s).
3. If `--monochrome` is given, set `DisplayMonochrome` on the computer state.
4. If `--start-at` is given, run the warm-up loop (see 6.3).
5. Create a recorder and register observers based on flags.
6. If `--record-audio`, create an `a2audio.Stream` and `AudioRecorder`.
7. If `--capture-video`, create a `VideoRecorder` with declared capture steps.
8. Call the headless `Run()` function.
9. Write results to the output directory.

## 6.3. Warm-Up Loop (--start-at)

When `--start-at ADDR` is given, the emulator executes instructions without
recording, stopping when `comp.CPU.PC` first equals `ADDR`. The `--steps N`
count then begins from that point, so step 1 is the instruction at `ADDR`.

This is useful when user code begins at a known address (e.g. `$0801` for
Apple II programs assembled with the default origin) and the test should not
count ROM boot cycles. Without `--start-at`, `--steps` counts from the very
first instruction after `comp.Boot()`.

The warm-up loop is capped at 10,000,000 iterations. If `ADDR` is not reached
within that limit, the command fails with an error. This prevents an infinite
loop when the address is never executed (e.g. wrong disk image, no code at
that address).

The warm-up loop calls `comp.Process()` directly, bypassing the recorder.
Observers are not registered until after the warm-up completes, so no state
entries are produced for ROM boot instructions.

## 6.4. Output Files

- `state.log` -- state entries as text (spec 1 format), written only if
  entries exist
- `audio.pcm` -- raw mono float32 little-endian PCM at 44100 Hz, written only
  if `--record-audio` produced samples
- `video.frame` -- all captured frames written to a single file, each as a
  text grid using the format defined in spec 3 (section 7.1), with a color
  legend derived from the distinct colors in the captured frame. Written only
  if `--capture-video` produced captures.

# 7. Files

| File | Change |
|------|--------|
| `gfx/framebuffer.go` | Lazy `ebiten.Image` creation in `Render()` |
| `gfx/border.go` | Lazy `ebiten.Image` pixel in `Draw()` |
| `gfx/textoverlay.go` | Lazy font source in `Draw()` |
| `record/run.go` | New -- headless execution loop |
| `record/recorder.go` | Add `CurrentStep()` method |
| `record/video_recorder.go` | Add `NeedsCapture()` method |
| `cmd/headless.go` | New -- `erc headless` cobra command |

# 8. What Does Not Change

- `record/*.go` (assertion, segment, observer code) -- already
  headless-compatible
- `clock/emulator.go` -- not used in headless mode (step-by-step, no timing)
- `a2/process.go`, `a2/boot.go`, `a2/load.go` -- core emulation, no graphics
  dependencies
- `emu/computer.go` -- interface is unchanged
- Display renderers (`a2hires`, `a2text`, etc.) -- write to `gfx.Screen`
  pixel buffer
- `a2/a2audio/audio.go` -- `Stream.Read()` works without ebiten audio context

# 9. Verification

1. `just lint` passes.
2. Existing tests pass (`go test ./...`) -- lazy ebiten changes are
   transparent to code that does call `Render()`/`Draw()`.
3. A headless integration test can boot the emulator, run N steps, and check
   that state entries are produced.
4. Video capture: run steps, render at a capture step, confirm
   `VideoRecorder.Frame()` returns pixel data.
5. Audio capture: create a `Stream`, attach an `AudioRecorder`, run steps,
   check that `Samples()` is non-empty.
