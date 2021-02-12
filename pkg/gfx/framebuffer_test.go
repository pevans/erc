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
