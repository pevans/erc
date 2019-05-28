package a2

import (
	"image/color"

	"github.com/pevans/erc/pkg/gfx"
)

func (c *Computer) Draw() error {
	for i := 0; i < 10; i++ {
		gfx.Scr.DrawDot(
			gfx.Point{X: 50 + i, Y: 50 + i},
			color.RGBA{R: 0, G: 0xFF, B: 0, A: 0},
		)
	}

	return nil
}

func (c *Computer) Dimensions() (width, height int) {
	return 280, 192
}
