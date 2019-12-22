package gfx

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten"
)

type Screen struct {
	Image *ebiten.Image
}

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

func (s *Screen) DrawDot(coord image.Point, color color.RGBA) {
	s.Image.Set(coord.X, coord.Y, color)
}
