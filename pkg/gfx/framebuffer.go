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

	Rows uint
	Cols uint
}

// NewFrameBuffer returns a new frame buffer that contains a set of
// logical rows and columns. These rows and columns should match
// whatever system you are emulating, as opposed to what might
// necessarily be shown on screen.
func NewFrameBuffer(rows, cols uint) *FrameBuffer {
	fb := new(FrameBuffer)

	fb.Rows = rows
	fb.Cols = cols
	fb.pixelsLength = rows * cols * 4
	fb.pixels = make([]byte, fb.pixelsLength)

	return fb
}

// cell returns the index of a cell within the Cells slice. In essence,
// given X rows and Y columns, you can think of the slice of cells as Y
// cells in a single row, followed another row, and another row...
func (fb *FrameBuffer) cell(row, col uint) uint {
	return (row * fb.Cols * 4) + (col * 4)
}

// getCell returns a cell's color, if one exists, or an error if not. This
// essentially translates the underlying cell structure into something similar
// to what gets passed in with SetCell.
func (fb *FrameBuffer) getCell(row, col uint) (color.RGBA, error) {
	i := fb.cell(row, col)

	if i > fb.pixelsLength {
		return color.RGBA{}, fmt.Errorf("out of bounds: (row %d, col %d)", row, col)
	}

	return color.RGBA{
		R: fb.pixels[i+0],
		G: fb.pixels[i+1],
		B: fb.pixels[i+2],
		A: fb.pixels[i+3],
	}, nil
}

// SetCell will assign the color of a single cell
func (fb *FrameBuffer) SetCell(row, col uint, clr color.RGBA) error {
	cellIndex := fb.cell(row, col)

	if cellIndex > fb.pixelsLength {
		return fmt.Errorf("out of bounds: (row %d, col %d)", row, col)
	}

	fb.pixels[cellIndex+0] = byte(clr.R)
	fb.pixels[cellIndex+1] = byte(clr.G)
	fb.pixels[cellIndex+2] = byte(clr.B)
	fb.pixels[cellIndex+3] = byte(clr.A)

	return nil
}

// ClearCells will set a color on every cell of the frame buffer
func (fb *FrameBuffer) ClearCells(clr color.RGBA) {
	for i := uint(0); i < fb.pixelsLength; i += 4 {
		fb.pixels[i+0] = byte(clr.R)
		fb.pixels[i+1] = byte(clr.G)
		fb.pixels[i+2] = byte(clr.B)
		fb.pixels[i+3] = byte(clr.A)
	}
}

// Render will accept an ebiten image and ~do something with it~ to render the
// contents of our frame buffer.
func (fb *FrameBuffer) Render(img *ebiten.Image) error {
	img.ReplacePixels(fb.pixels)
	return nil
}
