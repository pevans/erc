package main

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type screen struct {
	*ebiten.Image
}

// DrawDot satisfies the DotDrawer interface, and sets a single pixel in
// an ebiten image (screen).
func (s screen) DrawDot(coord image.Point, color color.RGBA) {
	s.Set(coord.X, coord.Y, color)
}

func gameLoop() error {
	var (
		width, height = emulator.Drawer.Dimensions()
	)

	return ebiten.Run(ebitenLoop, width, height, 3, "erc")
}

// ebitenLoop is the "run loop" of our graphics logic, which we use both
// to implement processor speed and frame rate.
func ebitenLoop(image *ebiten.Image) error {
	scr := screen{image}

	if err := emulator.Processor.Process(); err != nil {
		log.Error(errors.Wrap(err, "main loop received error from processor"))
	}

	if err := emulator.Drawer.Draw(scr); err != nil {
		log.Error(errors.Wrap(err, "main loop received error from drawer"))
	}

	return nil
}
