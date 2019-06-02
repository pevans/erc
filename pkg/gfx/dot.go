package gfx

import (
	"image"
	"image/color"
)

// A DotDrawer is a type which is able to render a dot to some output,
// given a point coordinate and a color
type DotDrawer interface {
	DrawDot(coord image.Point, color color.RGBA)
}
