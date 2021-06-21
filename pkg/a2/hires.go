package a2

import (
	"image/color"

	"github.com/pevans/erc/pkg/data"
)

var (
	hiresBlack  = color.RGBA{R: 0x00, G: 0x00, B: 0x00}
	hiresWhite  = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF}
	hiresGreen  = color.RGBA{R: 0x00, G: 0xFF, B: 0x00}
	hiresPurple = color.RGBA{R: 0x00, G: 0xFF, B: 0xFF}
	hiresBlue   = color.RGBA{R: 0x00, G: 0x00, B: 0xFF}
	hiresOrange = color.RGBA{R: 0xFF, G: 0xFF, B: 0x00}
)

var hiresPalette0 = []color.RGBA{
	hiresBlack,
	hiresGreen,
	hiresPurple,
	hiresWhite,
}

var hiresPalette1 = []color.RGBA{
	hiresBlack,
	hiresBlue,
	hiresOrange,
	hiresWhite,
}

func (c *Computer) hiresRender(start, end data.DByte) {
	for addr := start; addr < end; addr++ {
		// Each byte consists of a set of dots to render
		byt := c.Get(addr)

		x, y := HiresPoint(addr)

		// Turns out this address does not map to a real point on the
		// screen. Note that there _should_ never be a time when x < y
		// and y >= 0, or vice versa, but...
		if x < 0 || y < 0 {
			continue
		}

		// Dots are always horizontally contiguous, so whatever range we
		// get, we want to render them left to right. (Hence the use of
		// x+i.)
		for i, clr := range HiresDots(byt) {
			c.FrameBuffer.SetCell(uint(x+i), uint(y), clr)
		}
	}
}

// HiresDots returns a set of colors (representing dots) based on a
// given byte. These colors only make sense in a single hires context.
func HiresDots(b data.Byte) []color.RGBA {
	dots := make([]color.RGBA, 7)
	pal := hiresPalette0

	if b&0x80 > 0 {
		pal = hiresPalette1
	}

	for i := 0; i < 7; i++ {
		pair := (b >> i) & 0x3
		dots[i] = pal[pair]

		if dots[i] == hiresWhite && i < 6 {
			dots[i+1] = dots[i]
			i++
		}
	}

	return dots
}

// HiresPoint returns an x,y coordinate (column, row) for a given high
// resolution address.
func HiresPoint(a data.DByte) (int, int) {
	var (
		off = a - 0x2000
		x   = hiresCols[off]
		y   = hiresRows[off]
	)

	return x, y
}
