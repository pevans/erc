package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pevans/erc/pkg/a2"
	"github.com/pevans/erc/pkg/boot"
	"github.com/pkg/errors"
)

type game struct {
	comp *a2.Computer
	log  *boot.Logger
}

func (g *game) Layout(outWidth, outHeight int) (scrWidth, scrHeight int) {
	return g.comp.Dimensions()
}

func (g *game) Draw(screen *ebiten.Image) {
	if err := g.comp.FrameBuffer.Render(screen); err != nil {
		log.Fatal(errors.Wrap(err, "could not render framebuffer"))
	}
}

func (g *game) Update() error {
	return nil
}
