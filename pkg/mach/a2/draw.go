package a2

import (
	"image"
	"image/color"

	"github.com/pevans/erc/pkg/gfx"
)

func (c *Computer) Draw(screen gfx.DotDrawer) error {
	for i := 0; i < 10; i++ {
		screen.DrawDot(
			image.Point{X: 50 + i, Y: 50 + i},
			color.RGBA{R: 0, G: 0xFF, B: 0, A: 0},
		)
	}

	return nil
}

func (c *Computer) DrawHires(screen gfx.DotDrawer) {
}

func (c *Computer) Dimensions() (width, height int) {
	return 280, 192
}
