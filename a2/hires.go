package a2

import (
	"fmt"
	"image/color"

	"github.com/pevans/erc/gfx"
)

const (
	palettePurpleGreen = iota
	paletteBlueOrange
)

type hiresDot struct {
	on      bool
	palette int
	clr     color.RGBA
}

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
	hiresPurple,
	hiresGreen,
	hiresWhite,
}

var hiresPalette1 = []color.RGBA{
	hiresBlack,
	hiresBlue,
	hiresOrange,
	hiresWhite,
}

func (c *Computer) hiresRender(start, end int) {
	dots := make([]hiresDot, 280)

	for y := uint(0); y < 192; y++ {
		c.HiresDots(y, dots)
		for x, dot := range dots {
			gfx.Screen.SetCell(uint(x), y, dot.clr)
		}
	}
}

// HiresDots fills a slice of dots with color information that indicates
// how to render a hires graphics row. If the length of dots is not
// sufficient to contain all the dots in such a row, this function will
// return an error.
func (c *Computer) HiresDots(row uint, dots []hiresDot) error {
	if len(dots) != 280 {
		return fmt.Errorf("dots slice must contain 280 items")
	}

	addr := hiresAddrs[row]

	for byteOffset := 0; byteOffset < 40; byteOffset++ {
		byt := c.Get(int(addr) + byteOffset)
		pal := palettePurpleGreen

		// The high bit tells us which palette to use; if it's 1, we
		// switch to the blue/orange palette.
		if byt&0x80 > 0 {
			pal = paletteBlueOrange
		}

		dotOffset := int(byteOffset) * 7

		// Loop through the bits and set dots to on or off for each bit.
		for byteColumn := 0; byteColumn < 7; byteColumn++ {
			dots[dotOffset+byteColumn].on = byte(byt&1) > 0
			dots[dotOffset+byteColumn].palette = pal

			byt >>= 1
		}
	}

	for i, _ := range dots {
		var (
			white  = 3
			black  = 0
			color1 = 1
			color2 = 2
		)

		thisOn := dots[i].on
		prevOn := (i-1 >= 0) && dots[i-1].on
		colors := hiresPalette0

		if dots[i].palette == paletteBlueOrange {
			colors = hiresPalette1

			if i == 0 {
				prevOn = true
			}
		}

		switch {
		case thisOn && prevOn:
			dots[i].clr = colors[white]

		case thisOn && !prevOn:
			thisColor := color1

			if i%2 > 0 {
				thisColor = color2
			}

			dots[i].clr = colors[thisColor]

		case prevOn && !thisOn:
			thisColor := color1

			if i%2 == 0 {
				thisColor = color2
			}

			dots[i].clr = colors[thisColor]

		default:
			dots[i].clr = colors[black]
		}
	}

	return nil
}
