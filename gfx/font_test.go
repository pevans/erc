package gfx

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFont(t *testing.T) {
	assert.NotNil(t, NewFont(1, 1))
}

func TestFontGlyph(t *testing.T) {
	var (
		f     = NewFont(1, 1)
		fb    = NewFrameBuffer(1, 1)
		black = color.RGBA{}
		red   = color.RGBA{R: 255}
	)

	fb.ClearCells(red)
	f.glyphMap[2] = fb

	t.Run("a glyph for a nonexistent rune should be the default glyph", func(t *testing.T) {
		c, _ := f.Glyph(1).getCell(0, 0)
		assert.Equal(t, black, c)
	})

	t.Run("a glyph for an existent rune should match what we defined", func(t *testing.T) {
		c, _ := f.Glyph(2).getCell(0, 0)
		assert.Equal(t, red, c)
	})
}

func TestFontDefineGlyph(t *testing.T) {
	var (
		f     = NewFont(2, 2)
		white = color.RGBA{R: 255, G: 255, B: 255, A: 0xff}
		black = color.RGBA{A: 0xff}
	)

	t.Run("a new glyph should work", func(t *testing.T) {
		f.DefineGlyph(1, []byte{1, 0, 0, 1})

		// Testing the first point in the slice is white
		c, _ := f.Glyph(1).getCell(0, 0)
		assert.Equal(t, white, c)

		// Testing the third point in the slice is black
		c, _ = f.Glyph(1).getCell(0, 1)
		assert.Equal(t, black, c)
	})

	t.Run("a glyph with bad points size should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			f.DefineGlyph(2, []byte{1})
		})
	})
}

func TestFontDefineGlyphAsBuffer(t *testing.T) {
	var (
		f   = NewFont(2, 2)
		red = color.RGBA{R: 255}
	)

	t.Run("a framebuffer with matching dimensions should work", func(t *testing.T) {
		fb := NewFrameBuffer(2, 2)
		fb.ClearCells(red)

		f.DefineGlyphAsBuffer(1, fb)

		c, _ := f.Glyph(1).getCell(0, 0)
		assert.Equal(t, red, c)
	})

	t.Run("a framebuffer with mismatched dimensions should panic", func(t *testing.T) {
		fb := NewFrameBuffer(3, 3)

		assert.Panics(t, func() {
			f.DefineGlyphAsBuffer(2, fb)
		})
	})
}

func TestFontWrite(t *testing.T) {
	var (
		f     = NewFont(2, 2)
		red   = color.RGBA{R: 255}
		green = color.RGBA{G: 255}
		black = color.RGBA{}
	)

	redGlyph := NewFrameBuffer(2, 2)
	redGlyph.ClearCells(red)
	f.DefineGlyphAsBuffer('A', redGlyph)

	greenGlyph := NewFrameBuffer(2, 2)
	greenGlyph.ClearCells(green)
	f.DefineGlyphAsBuffer('B', greenGlyph)

	t.Run("write a string to framebuffer", func(t *testing.T) {
		fb := NewFrameBuffer(10, 2)

		assert.NoError(t, f.Write("AB", 0, 0, fb))

		c, _ := fb.getCell(0, 0)
		assert.Equal(t, red, c)

		c, _ = fb.getCell(2, 0)
		assert.Equal(t, green, c)
	})

	t.Run("cursor advances correctly", func(t *testing.T) {
		fb := NewFrameBuffer(10, 2)

		assert.NoError(t, f.Write("AB", 1, 0, fb))

		c, _ := fb.getCell(0, 0)
		assert.Equal(t, black, c)

		c, _ = fb.getCell(1, 0)
		assert.Equal(t, red, c)

		c, _ = fb.getCell(3, 0)
		assert.Equal(t, green, c)
	})

	t.Run("writing out of bounds returns error", func(t *testing.T) {
		fb := NewFrameBuffer(2, 2)

		assert.Error(t, f.Write("AB", 0, 0, fb))
	})
}
