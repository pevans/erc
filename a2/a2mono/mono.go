package a2mono

import "image/color"

const (
	None        = iota
	GreenScreen // Show color in the style of a green monochrome monitor
	AmberScreen // Show color in the style of an amber monochrome monitor
)

var (
	Black = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff} // basic black
	White = color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff} // basic white
	Green = color.RGBA{R: 0x98, G: 0xff, B: 0x98, A: 0xff} // a minty green
	Amber = color.RGBA{R: 0xff, G: 0xbf, B: 0x00, A: 0xff} // a true amber
)
