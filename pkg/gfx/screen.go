package gfx

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
)

type Point struct {
	X, Y int
}

type Renderer interface {
	DrawDot(coord Point, rgba color.RGBA)
}

type Screen struct {
	Image *ebiten.Image
}

var Scr Screen

func (s *Screen) DrawDot(coord Point, rgba color.RGBA) {
	s.Image.Set(coord.X, coord.Y, rgba)
}

func SetImage(img *ebiten.Image) {
	Scr.Image = img
}
