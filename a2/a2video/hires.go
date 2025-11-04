package a2video

import (
	"fmt"
	"image/color"

	"github.com/pevans/erc/gfx"
	"github.com/pevans/erc/memory"
)

const (
	PalettePurpleGreen = iota
	PaletteBlueOrange
)

// A HiresDot represents any dot on screen in a HIRES mode. In itself, a
// dot doesn't have enough information to determine its color; you must
// also look at a neighboring dot to figure that out.
type HiresDot struct {
	On       bool
	Boundary bool
	Palette  int
	Color    color.RGBA
}

var (
	HiresBlack       = color.RGBA{R: 0x00, G: 0x00, B: 0x00}
	HiresWhite       = color.RGBA{R: 0xff, G: 0xff, B: 0xff}
	HiresGreen       = color.RGBA{R: 0x2f, G: 0xbc, B: 0x1a}
	HiresPurple      = color.RGBA{R: 0xd0, G: 0x43, B: 0xe5}
	HiresBlue        = color.RGBA{R: 0x2f, G: 0x95, B: 0xe5}
	HiresOrange      = color.RGBA{R: 0xd0, G: 0x6a, B: 0x1a}
	HiresDarkGreen   = color.RGBA{R: 0x3f, G: 0x4c, B: 0x12}
	HiresDarkPurple  = color.RGBA{R: 0x3e, G: 0x31, B: 0x79}
	HiresLightGreen  = color.RGBA{R: 0xbd, G: 0xea, B: 0x86}
	HiresLightPurple = color.RGBA{R: 0xbb, G: 0xaf, B: 0xf6}
)

var purpleGreen = []color.RGBA{
	HiresPurple,
	HiresGreen,
}

var blueOrange = []color.RGBA{
	HiresBlue,
	HiresOrange,
}

// RenderHires draws dots in memory onto the screen
func RenderHires(seg memory.Getter, start, end int) {
	dots := make([]HiresDot, 280)

	for y := uint(0); y < 192; y++ {
		err := PrepareHiresRow(seg, y, dots)
		if err != nil {
			// This should really never happen...
			panic(err)
		}

		for x, dot := range dots {
			xpos := x * 2
			ypos := y * 2
			_ = gfx.Screen.SetCell(uint(xpos), ypos, dot.Color)
			_ = gfx.Screen.SetCell(uint(xpos), ypos+1, dot.Color)
			_ = gfx.Screen.SetCell(uint(xpos+1), ypos, dot.Color)
			_ = gfx.Screen.SetCell(uint(xpos+1), ypos+1, dot.Color)
		}
	}
}

// PrepareHiresRow fills a slice of dots with color information that
// indicates how to render a hires graphics row. If the length of dots
// is not sufficient to contain all the dots in such a row, this
// function will return an error.
func PrepareHiresRow(seg memory.Getter, row uint, dots []HiresDot) error {
	if err := fillHiresDots(seg, row, dots); err != nil {
		return err
	}

	// This is technically a double scan of the row. We could probably
	// make this faster, but on modern hardware, this hasn't been a
	// problem worth solving.
	for i := range dots {
		thisOn := dots[i].On
		prevOn := (i-1 >= 0) && dots[i-1].On
		colors := purpleGreen
		colorIndex := i % 2

		if dots[i].Palette == PaletteBlueOrange {
			colors = blueOrange
		}

		switch {
		case thisOn && prevOn:
			dots[i].Color = HiresWhite

		case thisOn && !prevOn:
			dots[i].Color = colors[colorIndex]

		case !thisOn && prevOn:
			// The XOR just flips the position of the color our index
			// would use; so if it would have been purple, now it's
			// green, or whatever.
			dots[i].Color = colors[colorIndex^1]

		default:
			dots[i].Color = HiresBlack
		}

		if i > 0 && dots[i-1].Boundary {
			dots[i-1], dots[i] = shiftBoundaryDots(dots[i-1], dots[i])
		}
	}

	return nil
}

// Fill in the boolean on/off state of each dot based on a 40 byte
// region that is implied by the given row.
func fillHiresDots(seg memory.Getter, row uint, dots []HiresDot) error {
	if len(dots) != 280 {
		return fmt.Errorf("dots slice must contain 280 items")
	}

	addr := hiresAddrs[row]

	for byteOffset := 0; byteOffset < 40; byteOffset++ {
		byt := seg.Get(int(addr) + byteOffset)
		pal := PalettePurpleGreen

		// The high bit tells us which palette to use; if it's 1, we
		// switch to the blue/orange palette.
		if byt&0x80 > 0 {
			pal = PaletteBlueOrange
		}

		dotOffset := int(byteOffset) * 7

		// Loop through the bits and set dots to on or off for each bit.
		for byteColumn := 0; byteColumn < 7; byteColumn++ {
			dots[dotOffset+byteColumn].On = byte(byt&1) > 0
			dots[dotOffset+byteColumn].Palette = pal

			byt >>= 1
		}

		dots[dotOffset+6].Boundary = true
	}

	return nil
}
