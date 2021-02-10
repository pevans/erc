package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pevans/erc/pkg/a2"
)

type game struct {
	comp *a2.Computer
}

func (g *game) Layout(outWidth, outHeight int) (scrWidth, scrHeight int) {
	return g.comp.Dimensions()
}

func (g *game) Draw(screen *ebiten.Image) {
	g.comp.FrameBuffer.Render(screen)
}

func (g *game) Update() error {
	err := g.comp.Process()
	return err
}
