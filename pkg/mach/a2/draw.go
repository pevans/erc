package a2

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten"
)

// Draw will figure out what to render on-screen at any given moment.
func (c *Computer) Draw() *ebiten.Image {
	if !c.reDraw {
		return c.Screen.Image
	}

	for i := 0; i < 10; i++ {
		c.Screen.DrawDot(
			image.Point{X: 50 + i, Y: 50 + i},
			color.RGBA{R: 0xFF, G: 0xFF, B: 0x00, A: 0xFF},
		)
	}

	c.reDraw = false
	return c.Screen.Image
}

// Dimensions returns the screen dimensions of an Apple II.
func (c *Computer) Dimensions() (width, height int) {
	return 280, 192
}
