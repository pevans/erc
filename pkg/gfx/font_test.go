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
		white = color.RGBA{R: 255, G: 255, B: 255}
		black = color.RGBA{}
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
