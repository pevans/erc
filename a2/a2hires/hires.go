package a2hires

import (
	"fmt"
	"image/color"

	"github.com/pevans/erc/a2/a2mono"
	"github.com/pevans/erc/a2/a2video"
	"github.com/pevans/erc/gfx"
	"github.com/pevans/erc/memory"
)

const (
	palettePurpleGreen = iota
	paletteBlueOrange
)

// A Dot represents any dot on screen in a HIRES mode. In itself, a dot
// doesn't have enough information to determine its color; you must also look
// at a neighboring dot to figure that out.
type Dot struct {
	on       bool
	boundary bool
	palette  int
	color    color.RGBA
}

var (
	black       = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff}
	white       = color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	green       = color.RGBA{R: 0x2f, G: 0xbc, B: 0x1a, A: 0xff}
	purple      = color.RGBA{R: 0xd0, G: 0x43, B: 0xe5, A: 0xff}
	blue        = color.RGBA{R: 0x2f, G: 0x95, B: 0xe5, A: 0xff}
	orange      = color.RGBA{R: 0xd0, G: 0x6a, B: 0x1a, A: 0xff}
	darkGreen   = color.RGBA{R: 0x3f, G: 0x4c, B: 0x12, A: 0xff}
	darkPurple  = color.RGBA{R: 0x3e, G: 0x31, B: 0x79, A: 0xff}
	lightGreen  = color.RGBA{R: 0xbd, G: 0xea, B: 0x86, A: 0xff}
	lightPurple = color.RGBA{R: 0xbb, G: 0xaf, B: 0xf6, A: 0xff}
)

var purpleGreen = []color.RGBA{
	purple,
	green,
}

var blueOrange = []color.RGBA{
	blue,
	orange,
}

// Render draws dots in memory onto the screen
func Render(seg memory.Getter, start, end int, monochromeMode int) {
	dots := make([]Dot, 280)

	for y := range uint(192) {
		err := PrepareRow(seg, y, dots, monochromeMode)
		if err != nil {
			// This should really never happen...
			panic(err)
		}

		for x, dot := range dots {
			xpos := x * 2
			ypos := y * 2
			_ = gfx.Screen.SetCell(uint(xpos), ypos, dot.color)
			_ = gfx.Screen.SetCell(uint(xpos), ypos+1, dot.color)
			_ = gfx.Screen.SetCell(uint(xpos+1), ypos, dot.color)
			_ = gfx.Screen.SetCell(uint(xpos+1), ypos+1, dot.color)
		}
	}
}

// PrepareRow fills a slice of dots with color information that indicates how
// to render a hires graphics row. If the length of dots is not sufficient to
// contain all the dots in such a row, this function will return an error.
func PrepareRow(seg memory.Getter, row uint, dots []Dot, monochromeMode int) error {
	if err := fillDots(seg, row, dots); err != nil {
		return err
	}

	if monochromeMode != a2mono.None {
		monochromeColor := a2mono.Green
		if monochromeMode == a2mono.AmberScreen {
			monochromeColor = a2mono.Amber
		}

		for i := range dots {
			if dots[i].on {
				dots[i].color = monochromeColor
			} else {
				dots[i].color = a2mono.Black
			}
		}

		return nil
	}

	// This is technically a double scan of the row. We could probably make
	// this faster, but on modern hardware, this hasn't been a problem worth
	// solving.
	for i := range dots {
		thisOn := dots[i].on
		prevOn := (i-1 >= 0) && dots[i-1].on
		colors := purpleGreen
		colorIndex := i % 2

		if dots[i].palette == paletteBlueOrange {
			colors = blueOrange
		}

		switch {
		case thisOn && prevOn:
			dots[i].color = white

		case thisOn && !prevOn:
			dots[i].color = colors[colorIndex]

		case !thisOn && prevOn:
			// The XOR just flips the position of the color our index would
			// use; so if it would have been purple, now it's green, or
			// whatever.
			dots[i].color = colors[colorIndex^1]

		default:
			dots[i].color = black
		}

		if i > 0 && dots[i-1].boundary {
			dots[i-1], dots[i] = shiftBoundaryDots(dots[i-1], dots[i])
		}
	}

	return nil
}

// fillDots fills in the boolean on/off state of each dot based on a 40 byte
// region that is implied by the given row.
func fillDots(seg memory.Getter, row uint, dots []Dot) error {
	if len(dots) != 280 {
		return fmt.Errorf("dots slice must contain 280 items")
	}

	addr := a2video.HiresAddrs[row]

	for byteOffset := range 40 {
		byt := seg.Get(int(addr) + byteOffset)
		pal := palettePurpleGreen

		// The high bit tells us which palette to use; if it's 1, we switch to
		// the blue/orange palette.
		if byt&0x80 > 0 {
			pal = paletteBlueOrange
		}

		dotOffset := int(byteOffset) * 7

		// Loop through the bits and set dots to on or off for each bit.
		for byteColumn := range 7 {
			dots[dotOffset+byteColumn].on = byte(byt&1) > 0
			dots[dotOffset+byteColumn].palette = pal

			byt >>= 1
		}

		dots[dotOffset+6].boundary = true
	}

	return nil
}
