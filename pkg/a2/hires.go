package a2

import (
	"fmt"
	"image/color"
)

type hiresDot struct {
	bits    byte
	palette []color.RGBA
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

func (c *Computer) hiresRender(start, end uint16) {
	dots := make([]hiresDot, 280)

	for y := uint(0); y < 192; y++ {
		c.HiresDots(y, dots)
		for x, dot := range dots {
			c.FrameBuffer.SetCell(uint(x), y, dot.clr)
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

	for i := 0; i < 40; i++ {
		byt := c.Get(int(addr) + i)
		pal := hiresPalette0

		if byt&0x80 > 0 {
			pal = hiresPalette1
		}

		for d := 0; d < 7; d++ {
			dots[(int(i)*7)+d].bits = byte(byt & 3)
			dots[(int(i)*7)+d].palette = pal
			byt >>= 1
		}
	}

	evenOffset := 0
	for i, dot := range dots {
		palIndex := int(dot.bits)
		if dot.bits > 0 && dot.bits < 3 {
			palIndex += evenOffset
		}

		dots[i].clr = dot.palette[palIndex]
		evenOffset ^= 1
	}

	return nil
}
