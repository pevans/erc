package render

import (
	"fmt"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/pevans/erc/a2"
	"github.com/pevans/erc/a2/a2audio"
	"github.com/pevans/erc/gfx"
	"github.com/pevans/erc/input"
)

// A game is just a small struct which ebiten will use to run the draw loop for us.
type game struct {
	comp           *a2.Computer
	keys           []ebiten.Key
	inputEvent     input.Event
	lastInputEvent input.Event
	audioPlayer    *audio.Player
}

// DrawLoop executes the logic to render our graphics according to some cadence
// (which is generally x frames per second).
func DrawLoop(comp *a2.Computer, shaderName string) error {
	w, h := comp.Dimensions()

	ebiten.SetWindowSize(int(w*2), int(h*2))
	ebiten.SetWindowTitle("erc")

	// Set up the shader if requested
	if err := gfx.Screen.SetShader(shaderName); err != nil {
		return err
	}

	// Set up audio
	audioCtx := audio.NewContext(a2audio.SampleRate)
	audioStream := a2audio.NewStream(comp.Speaker, comp)
	audioPlayer, err := audioCtx.NewPlayer(audioStream)
	if err != nil {
		slog.Error(fmt.Sprintf("could not create audio player: %v", err))
	}

	g := &game{
		comp:        comp,
		keys:        []ebiten.Key{},
		audioPlayer: audioPlayer,
	}

	// Start audio playback with reduced buffer for lower latency
	if g.audioPlayer != nil {
		// Smaller buffer = lower latency but higher risk of glitches
		// Default is typically ~46ms, we use ~23ms (1024 samples at 44100 Hz)
		g.audioPlayer.SetBufferSize(a2audio.BufferSamples * 4) // 4 bytes per stereo sample
		g.audioPlayer.Play()
	}

	return ebiten.RunGame(g)
}

// Layout returns the logical dimensions that ebiten should use.
func (g *game) Layout(outWidth, outHeight int) (scrWidth, scrHeight int) {
	w, h := g.comp.Dimensions()
	return int(w), int(h)
}

// Draw executes the render logic for the framebuffer.
func (g *game) Draw(screen *ebiten.Image) {
	if err := gfx.Screen.Render(screen); err != nil {
		slog.Error(fmt.Sprintf("could not render framebuffer: %v", err))
		return
	}

	gfx.StatusOverlay.Draw(screen)
	gfx.PrefixOverlay.Draw(screen)
	gfx.TextNotification.Draw(screen)
}

// Update handles logic once for every frame that is rendered, but this
// method is not _the method_ that renders the screen.
func (g *game) Update() error {
	if ebiten.IsWindowBeingClosed() {
		err := g.comp.Shutdown()
		if err != nil {
			return fmt.Errorf("could not properly shut down: %w", err)
		}

		return nil
	}

	g.keys = g.keys[:0]
	g.keys = inpututil.AppendPressedKeys(g.keys)

	if len(g.keys) > 0 {
		g.pushInputEvent()
	} else {
		g.inputEvent = input.Event{}
		g.lastInputEvent = input.Event{}
		g.comp.ClearKeys()
	}

	g.comp.Render()

	gfx.StatusOverlay.Update()
	gfx.PrefixOverlay.Update()
	gfx.TextNotification.Update()

	return nil
}

func (g *game) pushInputEvent() {
	for _, k := range g.keys {
		// If we see a modifier among the keys, we set the input event's
		// modifier. It's possible we've seen multiple modifiers -- if
		// so, the previous modifier is clobbered.
		mod := modifier(k)
		if mod != input.ModNone {
			g.inputEvent.Modifier = mod
			continue
		}

		g.inputEvent.Key = gfx.KeyToRune(k, g.inputEvent.Modifier)
	}

	// Don't allow repeat keystrokes
	if g.inputEvent == g.lastInputEvent &&
		g.inputEvent.Key != input.KeyNone {
		return
	}

	// If we got through the key slice with some valid key, we'll push
	// the input event with whatever modifier is left. If we only saw a
	// modifier, we'll hold it in inputEvent until the next time Update
	// is called.
	if g.inputEvent.Key != input.KeyNone {
		input.PushEvent(g.inputEvent)

		g.lastInputEvent = g.inputEvent
	}
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
