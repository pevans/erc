package record

import (
	"image/color"

	"github.com/pevans/erc/gfx"
)

// VideoRecorder implements Observer and captures framebuffer snapshots at
// declared execution steps.
type VideoRecorder struct {
	screen       *gfx.FrameBuffer
	captureSteps map[int]bool
	frames       map[int]*frameSnapshot
}

// A frameSnapshot stores a copy of the framebuffer pixel data along with its
// dimensions.
type frameSnapshot struct {
	pixels []byte
	width  uint
	height uint
}

// NewVideoRecorder returns a VideoRecorder that captures from the given
// framebuffer.
func NewVideoRecorder(screen *gfx.FrameBuffer) *VideoRecorder {
	return &VideoRecorder{
		screen:       screen,
		captureSteps: make(map[int]bool),
		frames:       make(map[int]*frameSnapshot),
	}
}

// CaptureAt declares which steps to capture.
func (v *VideoRecorder) CaptureAt(steps ...int) {
	for _, s := range steps {
		v.captureSteps[s] = true
	}
}

// Before is a no-op for the video recorder.
func (v *VideoRecorder) Before() {}

// Observe captures the framebuffer if the current step is in the capture set.
func (v *VideoRecorder) Observe(step int) []Entry {
	if !v.captureSteps[step] {
		return nil
	}

	v.frames[step] = &frameSnapshot{
		pixels: v.screen.Pixels(),
		width:  v.screen.Width,
		height: v.screen.Height,
	}

	return nil
}

// NeedsCapture reports whether the given step is in the capture set.
func (v *VideoRecorder) NeedsCapture(step int) bool {
	return v.captureSteps[step]
}

// Frame returns the captured snapshot for a given step, or nil if no capture
// exists.
func (v *VideoRecorder) Frame(step int) *frameSnapshot {
	return v.frames[step]
}

// Width returns the width of the frame snapshot.
func (f *frameSnapshot) Width() uint {
	return f.width
}

// Height returns the height of the frame snapshot.
func (f *frameSnapshot) Height() uint {
	return f.height
}

// GetPixel returns the color at (x, y) in a frame snapshot.
func (f *frameSnapshot) GetPixel(x, y uint) color.RGBA {
	if x >= f.width || y >= f.height {
		return color.RGBA{}
	}

	i := 4 * ((y * f.width) + x)

	return color.RGBA{
		R: f.pixels[i+0],
		G: f.pixels[i+1],
		B: f.pixels[i+2],
		A: f.pixels[i+3],
	}
}
