package gfx

import (
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/stretchr/testify/assert"
)

var (
	red   = color.RGBA{R: 255}
	green = color.RGBA{G: 255}
	blue  = color.RGBA{B: 255}
)

func TestNewFrameBuffer(t *testing.T) {
	var (
		rows uint = 123
		cols uint = 234
	)

	fb := NewFrameBuffer(rows, cols)
	assert.NotNil(t, fb)
	assert.Equal(t, uint(len(fb.pixels)), fb.pixelsLength)
	assert.Equal(t, rows, fb.Rows)
	assert.Equal(t, cols, fb.Cols)
}

func TestSetCell(t *testing.T) {
	var (
		rows uint = 123
		cols uint = 222
	)

	fb := NewFrameBuffer(rows, cols)

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
	assert.NoError(t, fb.SetCell(rows-1, cols-1, green))
	c, _ = fb.getCell(rows-1, cols-1)
	assert.Equal(t, green, c)

	// This should be out of bounds
	assert.Error(t, fb.SetCell(rows, cols, blue))
}

func TestClearCells(t *testing.T) {
	var (
		rows uint = 111
		cols uint = 222
	)

	fb := NewFrameBuffer(rows, cols)

	fb.ClearCells(blue)
	beg, _ := fb.getCell(0, 0)
	end, _ := fb.getCell(rows-1, cols-1)

	assert.Equal(t, blue, beg)
	assert.Equal(t, blue, end)
}

func TestRender(t *testing.T) {
	var (
		rows uint = 111
		cols uint = 222
	)

	fb := NewFrameBuffer(rows, cols)
	img := ebiten.NewImage(int(rows), int(cols))

	// This probably doesn't matter, but 🤷‍♂️
	fb.ClearCells(red)

	assert.NoError(t, fb.Render(img))
}
