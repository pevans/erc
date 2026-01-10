package gfx

import (
	_ "embed"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed shaders/softcrt.kage
var softcrtShaderSource []byte

//go:embed shaders/hardcrt.kage
var hardcrtShaderSource []byte

// A FrameBuffer is a set of cells which contain color information.
type FrameBuffer struct {
	// pixels is a slice of contiguous bytes that represent the on-screen
	// pixels that we'll render.
	pixels []byte

	// pixelsLength is the effective size necessary to hold the width/height
	// of the FrameBuffer.
	pixelsLength uint

	// Image is the result of the FrameBuffer's pixels -- it those pixels
	// turned into a raw graphic that our engine, Ebiten, can render
	Image *ebiten.Image

	shader     *ebiten.Shader // the shader we'll use to alter our graphics before render
	shaderName string         // the name of the shader

	Width  uint // the effective width of the FrameBuffer
	Height uint // the effective height of the FrameBuffer
}

// NewFrameBuffer returns a new frame buffer that contains a set of logical
// rows and columns. These rows and columns should match whatever system you
// are emulating, as opposed to what might necessarily be shown on screen.
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
// the receiver. This is useful for cases where you might want an inverse
// video effect.
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

// cell returns the index of a cell within the Cells slice. In essence, given
// X rows and Y columns, you can think of the slice of cells as Y cells in a
// single row, followed another row, and another row...
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

// GetPixel returns the color at a given pixel coordinate. If the coordinate
// is out of bounds, it returns black.
func (fb *FrameBuffer) GetPixel(x, y uint) color.RGBA {
	c, err := fb.getCell(x, y)
	if err != nil {
		return color.RGBA{R: 0, G: 0, B: 0, A: 0xff}
	}
	return c
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

// SetShader loads and sets a shader for the framebuffer. Pass "none" or empty
// string to disable shaders, "softcrt" for soft CRT shader, "hardcrt" for
// hard scanlines, or "curvedcrt" for curved CRT with scanlines.
func (fb *FrameBuffer) SetShader(shaderName string) error {
	if shaderName == "none" || shaderName == "" {
		fb.shader = nil
		fb.shaderName = "none"
		return nil
	}

	var shaderSource []byte
	switch shaderName {
	case "softcrt", "curvedcrt":
		shaderSource = softcrtShaderSource
	case "hardcrt":
		shaderSource = hardcrtShaderSource
	default:
		return fmt.Errorf("unknown shader: %s", shaderName)
	}

	shader, err := ebiten.NewShader(shaderSource)
	if err != nil {
		return fmt.Errorf("failed to compile %s shader: %w", shaderName, err)
	}

	fb.shader = shader
	fb.shaderName = shaderName
	return nil
}

// Render will accept an ebiten image and ~do something with it~ to render the
// contents of our frame buffer.
func (fb *FrameBuffer) Render(img *ebiten.Image) error {
	fb.Image.WritePixels(fb.pixels)

	if fb.shader != nil {
		opts := &ebiten.DrawRectShaderOptions{}
		opts.Images[0] = fb.Image

		// Set uniforms based on shader type
		switch fb.shaderName {
		case "softcrt":
			opts.Uniforms = map[string]any{
				"Curvature":      float32(0.0),  // No barrel distortion (flat screen)
				"ScanlineWeight": float32(0.20), // Subtle scanlines
			}
		case "curvedcrt":
			opts.Uniforms = map[string]any{
				"Curvature":      float32(0.15), // 15% barrel distortion (curved screen)
				"ScanlineWeight": float32(0.20), // Subtle scanlines
			}
		case "hardcrt":
			opts.Uniforms = map[string]any{
				"ScanlineWeight": float32(0.6), // Show a darker effect on the screen than you'd see for softcrt
			}
		}

		w, h := fb.Image.Bounds().Dx(), fb.Image.Bounds().Dy()
		img.DrawRectShader(w, h, fb.shader, opts)
	} else {
		img.DrawImage(fb.Image, nil)
	}

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
