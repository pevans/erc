package gfx

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// A FrameBuffer is a set of cells which contain color information.
type FrameBuffer struct {
	pixels       []byte
	pixelsLength uint

	Width  uint
	Height uint
}

// NewFrameBuffer returns a new frame buffer that contains a set of
// logical rows and columns. These rows and columns should match
// whatever system you are emulating, as opposed to what might
// necessarily be shown on screen.
func NewFrameBuffer(width, height uint) *FrameBuffer {
	fb := new(FrameBuffer)

	fb.Width = width
	fb.Height = height
	fb.pixelsLength = width * height * 4
	fb.pixels = make([]byte, fb.pixelsLength)

	return fb
}

// cell returns the index of a cell within the Cells slice. In essence,
// given X rows and Y columns, you can think of the slice of cells as Y
// cells in a single row, followed another row, and another row...
func (fb *FrameBuffer) cell(x, y uint) uint {
	return (y * fb.Height * 4) + (x * 4)
}

// getCell returns a cell's color, if one exists, or an error if not. This
// essentially translates the underlying cell structure into something similar
// to what gets passed in with SetCell.
func (fb *FrameBuffer) getCell(x, y uint) (color.RGBA, error) {
	i := fb.cell(x, y)

	if i > fb.pixelsLength {
		return color.RGBA{}, fmt.Errorf("out of bounds: (x %d, y %d)", x, y)
	}

	return color.RGBA{
		R: fb.pixels[i+0],
		G: fb.pixels[i+1],
		B: fb.pixels[i+2],
		A: fb.pixels[i+3],
	}, nil
}

// SetCell will assign the color of a single cell
func (fb *FrameBuffer) SetCell(x, y uint, clr color.RGBA) error {
	i := fb.cell(x, y)

	if i > fb.pixelsLength {
		return fmt.Errorf("out of bounds: (x %d, y %d)", x, y)
	}

	fb.pixels[i+0] = clr.R
	fb.pixels[i+1] = clr.G
	fb.pixels[i+2] = clr.B
	fb.pixels[i+3] = clr.A

	return nil
}

// ClearCells will set a color on every cell of the frame buffer
func (fb *FrameBuffer) ClearCells(clr color.RGBA) {
	for i := uint(0); i < fb.pixelsLength; i += 4 {
		fb.pixels[i+0] = clr.R
		fb.pixels[i+1] = clr.G
		fb.pixels[i+2] = clr.B
		fb.pixels[i+3] = clr.A
	}
}

// Render will accept an ebiten image and ~do something with it~ to render the
// contents of our frame buffer.
func (fb *FrameBuffer) Render(img *ebiten.Image) error {
	img.ReplacePixels(fb.pixels)
	return nil
}
