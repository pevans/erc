package gfx

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	red   = color.RGBA{R: 255}
	green = color.RGBA{G: 255}
	blue  = color.RGBA{B: 255}
)

func TestNewFrameBuffer(t *testing.T) {
	var (
		width  uint = 123
		height uint = 234
	)

	fb := NewFrameBuffer(width, height)
	assert.NotNil(t, fb)
	assert.Equal(t, uint(len(fb.pixels)), fb.pixelsLength)
	assert.Equal(t, width, fb.Width)
	assert.Equal(t, height, fb.Height)
}

func TestFrameBufferInvert(t *testing.T) {
	t.Run("invert returns a frame buffer with inverted colors", func(t *testing.T) {
		fb := NewFrameBuffer(1, 1)
		red := color.RGBA{R: 255}
		blueGreen := color.RGBA{G: 255, B: 255}

		fb.ClearCells(red)
		inv := fb.Invert()

		cbg, _ := inv.getCell(0, 0)
		assert.Equal(t, blueGreen, cbg)
	})
}

func TestSetCell(t *testing.T) {
	var (
		width  uint = 123
		height uint = 222
	)

	fb := NewFrameBuffer(width, height)

	assert.NoError(t, fb.SetCell(0, 0, green))
	c, _ := fb.getCell(0, 0)
	assert.Equal(t, green, c)

	assert.NoError(t, fb.SetCell(1, 0, blue))
	c, _ = fb.getCell(1, 0)
	assert.Equal(t, blue, c)

	assert.NoError(t, fb.SetCell(2, 0, red))
	c, _ = fb.getCell(2, 0)
	assert.Equal(t, red, c)

	// This would be the maximum position
	assert.NoError(t, fb.SetCell(width-1, height-1, green))
	c, _ = fb.getCell(width-1, height-1)
	assert.Equal(t, green, c)

	// This should be out of bounds
	assert.Error(t, fb.SetCell(width, height, blue))
}

func TestClearCells(t *testing.T) {
	var (
		width  uint = 111
		height uint = 222
	)

	fb := NewFrameBuffer(width, height)

	fb.ClearCells(blue)
	beg, _ := fb.getCell(0, 0)
	end, _ := fb.getCell(width-1, height-1)

	assert.Equal(t, blue, beg)
	assert.Equal(t, blue, end)
}

func TestBlit(t *testing.T) {
	var (
		srcWidth  uint = 111
		srcHeight uint = 222
		red            = color.RGBA{R: 240}
		black          = color.RGBA{}
		src            = NewFrameBuffer(srcWidth, srcHeight)
	)

	src.ClearCells(red)

	t.Run("copy src into an equal size fb", func(t *testing.T) {
		fb := NewFrameBuffer(srcWidth, srcHeight)
		assert.NoError(t, fb.Blit(0, 0, src))
	})

	t.Run("copy into portion of dest", func(t *testing.T) {
		fb := NewFrameBuffer(srcWidth, srcHeight+1)

		// Skip the first row in the copy
		assert.NoError(t, fb.Blit(0, 1, src))

		// Make sure that first row is still the default color, but the second
		// should be red
		cb, _ := fb.getCell(0, 0)
		cr, _ := fb.getCell(1, 1)
		assert.Equal(t, black, cb)
		assert.Equal(t, red, cr)
	})

	t.Run("can't copy beyond boundaries", func(t *testing.T) {
		var (
			w uint = 1
			h uint = 1
		)

		fb := NewFrameBuffer(w, h)
		assert.Error(t, fb.Blit(0, 0, src))
	})
}
