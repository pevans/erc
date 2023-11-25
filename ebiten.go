package main

import (
	"fmt"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/pevans/erc/a2"
	"github.com/pevans/erc/gfx"
	"github.com/pevans/erc/input"
)

// A game is just a small struct which ebiten will use to run the draw loop for
// us.
type game struct {
	comp       *a2.Computer
	keys       []ebiten.Key
	inputEvent input.Event
}

// drawLoop executes the logic to render our graphics according to some cadence
// (which is generally x frames per second).
func drawLoop(comp *a2.Computer) error {
	w, h := comp.Dimensions()

	ebiten.SetWindowSize(w*3, h*3)
	ebiten.SetWindowTitle("erc")

	g := &game{
		comp: comp,
		keys: []ebiten.Key{},
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
		slog.Error(fmt.Sprintf("could not render framebuffer: %v", err))
		return
	}
}

// Update is kind of a noop for us. Nominally you could use it to execute game
// logic, but it will run as often as the frames on screen will update--this
// ends up being too infrequently for us to make use of it.
func (g *game) Update() error {
	g.keys = inpututil.AppendPressedKeys(g.keys)

	for _, k := range g.keys {
		// If we see a modifier among the keys, we set the input event's
		// modifier. It's possible we've seen multiple modifiers -- if
		// so, the previous modifier is clobbered.
		if mod := modifier(k); mod != input.ModNone {
			g.inputEvent.Modifier = mod
			continue
		}

		g.inputEvent.Key = gfx.KeyToRune(k)
	}

	// If we got through the key slice with some valid key, we'll push
	// the input event with whatever modifier is left. If we only saw a
	// modifier, we'll hold it in inputEvent until the next time Update
	// is called.
	if g.inputEvent.Key != input.KeyNone {
		input.PushEvent(g.inputEvent)
		g.inputEvent = input.Event{}
	}

	// Wipe the slice without freeing its capacity
	g.keys = g.keys[:0]

	g.comp.Render()
	return nil
}

func modifier(key ebiten.Key) int {
	switch key {
	case ebiten.KeyControl, ebiten.KeyControlLeft, ebiten.KeyControlRight:
		return input.ModControl
	case ebiten.KeyShift, ebiten.KeyShiftLeft, ebiten.KeyShiftRight:
		return input.ModShift
	case ebiten.KeyAlt, ebiten.KeyAltLeft, ebiten.KeyAltRight:
		return input.ModOption
	case ebiten.KeyMeta, ebiten.KeyMetaLeft, ebiten.KeyMetaRight:
		return input.ModCommand
	}

	return input.ModNone
}
