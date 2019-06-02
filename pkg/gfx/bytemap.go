package gfx

import (
	"image"
	"image/color"

	"github.com/pevans/erc/pkg/data"
)

type Bytemap struct {
	Dots          []data.Byte
	Width, Height int
}

// Valid returns true if the number of dots a Bytemap contains is equal
// to the product of its width and height.
func (b *Bytemap) Valid() bool {
	return len(b.Dots) == b.Width*b.Height
}

// Draw will render all of the dots in the bytemap from a given point,
// using a given color.
func (b *Bytemap) Draw(screen DotDrawer, from image.Point, color color.RGBA) {
	if !b.Valid() {
		return
	}

	var offset image.Point

	for i := 0; offset.Y < b.Height; offset.Y++ {
		for ; offset.X < b.Width; offset.X++ {
			var useColor = Black

			// We only render the given color if the dot at our position
			// is nonzero
			if b.Dots[i] != 0 {
				useColor = color
			}

			screen.DrawDot(from.Add(offset), useColor)

			i++
		}
	}
}
