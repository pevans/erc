package gfx

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten"
)

// A Screen is a simple abstraction that hides the details of how we
// gather and send back the pixels to be rendered by our graphics
// library.
type Screen struct {
	Image *ebiten.Image
}

// NewScreen returns a new screen object that is ready to be rendered
// with ebiten.
func NewScreen(width, height int) *Screen {
	var (
		err error
		s   = new(Screen)
	)

	s.Image, err = ebiten.NewImage(width, height, ebiten.FilterLinear)
	if err != nil {
		panic(err)
	}

	return s
}

// DrawDot will set the value of a pixel at a given coordinate to the
// given color.
func (s *Screen) DrawDot(coord image.Point, color color.RGBA) {
	s.Image.Set(coord.X, coord.Y, color)
}
