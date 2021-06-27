package a2

import (
	"image/color"

	"github.com/pevans/erc/pkg/data"
)

var (
	hiresBlack  = color.RGBA{R: 0x00, G: 0x00, B: 0x00}
	hiresWhite  = color.RGBA{R: 0xff, G: 0xff, B: 0xff}
	hiresGreen  = color.RGBA{R: 0x2f, G: 0xbc, B: 0x1a}
	hiresPurple = color.RGBA{R: 0xd0, G: 0x43, B: 0xe5}
	hiresBlue   = color.RGBA{R: 0x2f, G: 0x95, B: 0xe5}
	hiresOrange = color.RGBA{R: 0xd0, G: 0x6a, B: 0x1a}
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
	for y := uint(0); y < 192; y++ {
		addr := hiresAddrs[y]

		for i := 0; i < 40; i++ {
			dots := HiresDots(c.Get(addr + data.DByte(i)))

			for x, clr := range dots {
				c.FrameBuffer.SetCell(uint((i*7)+x), y, clr)
			}
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
