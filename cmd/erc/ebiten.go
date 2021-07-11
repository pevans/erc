package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/pevans/erc/pkg/a2"
	"github.com/pevans/erc/pkg/boot"
	"github.com/pevans/erc/pkg/data"
	"github.com/pevans/erc/pkg/gfx"
	"github.com/pkg/errors"
)

// A game is just a small struct which ebiten will use to run the draw loop for
// us.
type game struct {
	comp *a2.Computer
	log  *boot.Logger
}

// drawLoop executes the logic to render our graphics according to some cadence
// (which is generally x frames per second).
func drawLoop(comp *a2.Computer, log *boot.Logger) error {
	w, h := comp.Dimensions()

	ebiten.SetWindowSize(w*3, h*3)
	ebiten.SetWindowTitle("erc")

	g := &game{
		comp: comp,
		log:  log,
	}

	return ebiten.RunGame(g)
}

// Layout returns the logical dimensions that ebiten should use.
func (g *game) Layout(outWidth, outHeight int) (scrWidth, scrHeight int) {
	return g.comp.Dimensions()
}

// Draw executes the render logic for the framebuffer.
func (g *game) Draw(screen *ebiten.Image) {
	if err := g.comp.FrameBuffer.Render(screen); err != nil {
		log.Fatal(errors.Wrap(err, "could not render framebuffer"))
	}
}

// Update is kind of a noop for us. Nominally you could use it to execute game
// logic, but it will run as often as the frames on screen will update--this
// ends up being too infrequently for us to make use of it.
func (g *game) Update() error {
	// TODO: call inpututil.PressedKeys()? Get a slice of keys that are
	// pressed, and send the last one (I suppose) to s.comp.PressKey().
	for _, k := range inpututil.PressedKeys() {
		g.comp.PressKey(data.Byte(gfx.KeyToRune(k)))
	}

	g.comp.Render()
	return nil
}
