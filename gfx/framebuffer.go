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

	Image *ebiten.Image

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
	fb.Image = ebiten.NewImage(int(fb.Width), int(fb.Height))

	return fb
}

// Invert returns a framebuffer that is the opposite (or inverted) version of
// the receiver. This is useful for cases where you might want an inverse video
// effect.
func (fb *FrameBuffer) Invert() *FrameBuffer {
	inv := NewFrameBuffer(fb.Width, fb.Height)

	for i, px := range fb.pixels {
		if (i+1)%4 == 0 {
			continue
		}

		inv.pixels[i] = px ^ 0xff
	}

	return inv
}

// cell returns the index of a cell within the Cells slice. In essence,
// given X rows and Y columns, you can think of the slice of cells as Y
// cells in a single row, followed another row, and another row...
func (fb *FrameBuffer) cell(x, y uint) uint {
	return 4 * ((y * fb.Width) + x)
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
	fb.Image.WritePixels(fb.pixels)

	// TODO: maybe we could apply filters/shaders/etc. to modify how the
	// screen is rendered
	img.DrawImage(fb.Image, nil)

	return nil
}

// Blit will essentially copy the entire source framebuffer into the receiver,
// starting from a specific point.
func (fb *FrameBuffer) Blit(x, y uint, src *FrameBuffer) error {
	for sy := uint(0); sy < src.Height; sy++ {
		if err := fb.blitFromY(x, y+sy, sy, src); err != nil {
			return err
		}
	}

	return nil
}

// blitFromY is a helper method for blit; basically it encapsulates the logic of
// blitting a single row.
func (fb *FrameBuffer) blitFromY(x, y, sy uint, src *FrameBuffer) error {
	// Where we're writing to
	di := fb.cell(x, y)

	// Where we're writing from
	si := src.cell(0, sy)

	writeLength := src.Width * 4

	if fb.pixelsLength-di < writeLength {
		return fmt.Errorf(
			"destination out of bounds (pl[%d]-di[%d] < wl[%d]",
			fb.pixelsLength, di, writeLength,
		)
	}

	if src.pixelsLength-si < writeLength {
		return fmt.Errorf(
			"source out of bounds (pl[%d]-si[%d] < wl[%d]",
			src.pixelsLength, si, writeLength,
		)
	}

	// Remember that there are 4 pixels for every "cell" we need to copy!
	for slen := src.Width * 4; slen > 0; slen-- {
		fb.pixels[di] = src.pixels[si]
		di++
		si++
	}

	return nil
}
