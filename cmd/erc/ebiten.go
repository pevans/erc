package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/pevans/erc/pkg/a2"
	"github.com/pevans/erc/pkg/clog"
	"github.com/pevans/erc/pkg/gfx"
	"github.com/pevans/erc/pkg/input"
	"github.com/pkg/errors"
)

// A game is just a small struct which ebiten will use to run the draw loop for
// us.
type game struct {
	comp *a2.Computer
}

// drawLoop executes the logic to render our graphics according to some cadence
// (which is generally x frames per second).
func drawLoop(comp *a2.Computer) error {
	w, h := comp.Dimensions()

	ebiten.SetWindowSize(w*3, h*3)
	ebiten.SetWindowTitle("erc")

	g := &game{
		comp: comp,
	}

	return ebiten.RunGame(g)
}

// Layout returns the logical dimensions that ebiten should use.
func (g *game) Layout(outWidth, outHeight int) (scrWidth, scrHeight int) {
	return g.comp.Dimensions()
}

// Draw executes the render logic for the framebuffer.
func (g *game) Draw(screen *ebiten.Image) {
	if err := gfx.Screen.Render(screen); err != nil {
		clog.Error(errors.Wrap(err, "could not render framebuffer"))
	}
}

// Update is kind of a noop for us. Nominally you could use it to execute game
// logic, but it will run as often as the frames on screen will update--this
// ends up being too infrequently for us to make use of it.
func (g *game) Update() error {
	for _, k := range inpututil.PressedKeys() {
		input.PushEvent(input.Event{
			Key: gfx.KeyToRune(k),
		})
	}

	g.comp.Render()
	return nil
}
