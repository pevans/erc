package a2video

import "image/color"

const (
	MonochromeNone = iota
	MonochromeGreen
	MonochromeAmber
)

var (
	HiresBlack           = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff}
	HiresWhite           = color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	HiresMonochromeGreen = color.RGBA{R: 0x98, G: 0xff, B: 0x98, A: 0xff}
	HiresMonochromeAmber = color.RGBA{R: 0xff, G: 0xbf, B: 0x00, A: 0xff}
)
