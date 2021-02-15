package gfx

import (
	"fmt"
	"image/color"
)

// A Font is a bitmapped font which can represent glyphs, or renderings of a
// single character, as FrameBuffers that can then be blitted onto a larger
// FrameBuffer.
type Font struct {
	// GlyphWidth is the width of a certain glyph within the font
	GlyphWidth uint

	// GlyphHeight is the height of a certain glyph within the font
	GlyphHeight uint

	defaultGlyph *FrameBuffer
	glyphMap     map[rune]*FrameBuffer
}

// NewFont returns a new font in which each glyph will have dimensions based on
// the given width and height.
func NewFont(width, height uint) *Font {
	f := new(Font)

	f.GlyphWidth = width
	f.GlyphHeight = height
	f.defaultGlyph = NewFrameBuffer(f.GlyphWidth, f.GlyphHeight)
	f.glyphMap = make(map[rune]*FrameBuffer)

	return f
}

// Glyph returns a glyph for a given rune. If no such rune exists in our font,
// we will return the font's _default glyph_, rather than an error.
func (f *Font) Glyph(ch rune) *FrameBuffer {
	fb, ok := f.glyphMap[ch]
	if !ok {
		return f.defaultGlyph
	}

	return fb
}

// DefineGlyph will define a new glyph in the font, or replace an existing
// glyph, for a given rune. Points deserves special attention: it's expected to
// be a sequence of zeroes and ones, where zero indicates a point in the bitmap
// font that should not be drawn, and one indicates a point that should be
// drawn. The length of points should be equal to the product of width x height
// and, if it isn't, DefineGlyph will panic.
func (f *Font) DefineGlyph(ch rune, points []byte) {
	if len(points) != int(f.GlyphWidth)*int(f.GlyphHeight) {
		panic(fmt.Sprintf(
			"invalid points length for font (pl[%d] != w[%d] x h[%d]",
			len(points), f.GlyphWidth, f.GlyphHeight,
		))
	}

	fb := NewFrameBuffer(f.GlyphWidth, f.GlyphHeight)

	for i, pt := range points {
		ui := uint(i)

		// It's ok to ignore the error return here, since the only error
		// condition of SetCell can occur if you attempt an out-of-bounds set.
		// We confirmed that we can't by the check above on the length of
		// points.
		_ = fb.SetCell(ui%f.GlyphWidth, ui/f.GlyphWidth, gcolor(pt))
	}

	f.glyphMap[ch] = fb
}

// gcolor will return a color that, for our purposes, will suffice to indicate
// that one cell should be drawn or another should not.
func gcolor(b byte) color.RGBA {
	if b == 0 {
		return color.RGBA{}
	}

	return color.RGBA{R: 255, G: 255, B: 255}
}
