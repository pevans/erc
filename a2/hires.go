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
	on       bool
	boundary bool
	palette  int
	clr      color.RGBA
}

var (
	hiresBlack       = color.RGBA{R: 0x00, G: 0x00, B: 0x00}
	hiresWhite       = color.RGBA{R: 0xff, G: 0xff, B: 0xff}
	hiresGreen       = color.RGBA{R: 0x2f, G: 0xbc, B: 0x1a}
	hiresPurple      = color.RGBA{R: 0xd0, G: 0x43, B: 0xe5}
	hiresBlue        = color.RGBA{R: 0x2f, G: 0x95, B: 0xe5}
	hiresOrange      = color.RGBA{R: 0xd0, G: 0x6a, B: 0x1a}
	hiresDarkGreen   = color.RGBA{R: 0x3f, G: 0x4c, B: 0x12}
	hiresDarkPurple  = color.RGBA{R: 0x3e, G: 0x31, B: 0x79}
	hiresLightGreen  = color.RGBA{R: 0xbd, G: 0xea, B: 0x86}
	hiresLightPurple = color.RGBA{R: 0xbb, G: 0xaf, B: 0xf6}
)

var hiresPalette0 = []color.RGBA{
	hiresPurple,
	hiresGreen,
}

var hiresPalette1 = []color.RGBA{
	hiresBlue,
	hiresOrange,
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
	if err := hiresFillDots(c, row, dots); err != nil {
		return err
	}

	// This is technically a double scan of the row. We could probably
	// make this faster, but on modern hardware, this hasn't been a
	// problem worth solving.
	for i, _ := range dots {
		thisOn := dots[i].on
		prevOn := (i-1 >= 0) && dots[i-1].on
		colors := hiresPalette0
		colorIndex := i % 2

		if dots[i].palette == paletteBlueOrange {
			colors = hiresPalette1
		}

		switch {
		case thisOn && prevOn:
			dots[i].clr = hiresWhite

		case thisOn && !prevOn:
			dots[i].clr = colors[colorIndex]

		case !thisOn && prevOn:
			// The XOR just flips the position of the color our index
			// would use; so if it would have been purple, now it's
			// green, or whatever.
			dots[i].clr = colors[colorIndex^1]

		default:
			dots[i].clr = hiresBlack
		}

		if i > 0 && dots[i-1].boundary {
			dots[i-1], dots[i] = shiftBoundaryDots(dots[i-1], dots[i])
		}
	}

	return nil
}

// Fill in the boolean on/off state of each dot based on a 40 byte
// region that is implied by the given row.
func hiresFillDots(comp *Computer, row uint, dots []hiresDot) error {
	if len(dots) != 280 {
		return fmt.Errorf("dots slice must contain 280 items")
	}

	addr := hiresAddrs[row]

	for byteOffset := 0; byteOffset < 40; byteOffset++ {
		byt := comp.Get(int(addr) + byteOffset)
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

		dots[dotOffset+6].boundary = true
	}

	return nil
}
