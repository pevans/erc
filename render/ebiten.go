package render

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/pevans/erc/a2"
	"github.com/pevans/erc/a2/a2audio"
	"github.com/pevans/erc/gfx"
	"github.com/pevans/erc/input"
)

// A game is just a small struct which ebiten will use to run the draw loop
// for us.
type game struct {
	comp           *a2.Computer
	keys           []ebiten.Key
	inputEvent     input.Event
	lastInputEvent input.Event
	audioPlayer    *audio.Player
	keyPressTime   time.Time
	lastRepeatTime time.Time
}

// DrawLoop executes the logic to render our graphics according to some
// cadence (which is generally x frames per second).
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
	audioStream := a2audio.NewStream(comp.Speaker(), comp)
	audioPlayer, err := audioCtx.NewPlayerF32(audioStream)
	if err != nil {
		slog.Error(fmt.Sprintf("could not create audio player: %v", err))
	}

	// Give the computer access to the audio stream for volume control
	comp.SetAudioStream(audioStream)

	// Set audio logger if available (for debug-image mode)
	if comp.AudioLog != nil {
		audioStream.SetAudioLogger(comp.AudioLog)
	}

	g := &game{
		comp:        comp,
		keys:        []ebiten.Key{},
		audioPlayer: audioPlayer,
	}

	// Start audio playback
	if g.audioPlayer != nil {
		g.audioPlayer.SetBufferSize(0) // 0 = use Ebiten's default buffer size
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

// Update handles logic once for every frame that is rendered, but this method
// is not _the method_ that renders the screen.
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
		g.keyPressTime = time.Time{}
		g.lastRepeatTime = time.Time{}
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
		// modifier. It's possible we've seen multiple modifiers -- if so, the
		// previous modifier is clobbered.
		mod := modifier(k)
		if mod != input.ModNone {
			g.inputEvent.Modifier = mod
			continue
		}

		g.inputEvent.Key = gfx.KeyToRune(k, g.inputEvent.Modifier)
	}

	if g.inputEvent.Key == input.KeyNone {
		return
	}

	now := time.Now()

	// Check if this is a different event (key or modifier changed)
	if g.inputEvent != g.lastInputEvent {
		// New key press - send immediately
		input.PushEvent(g.inputEvent)
		g.lastInputEvent = g.inputEvent
		g.keyPressTime = now
		g.lastRepeatTime = time.Time{}
		return
	}

	// This is a repeat key. See if we should send the repeats through.
	timeSincePress := now.Sub(g.keyPressTime)

	// Wait at least this long...
	if timeSincePress < 500*time.Millisecond {
		return
	}

	// If they're still going, figure out how long we've been holding down.
	var timeSinceLastRepeat time.Duration
	if g.lastRepeatTime.IsZero() {
		timeSinceLastRepeat = timeSincePress - 500*time.Millisecond
	} else {
		timeSinceLastRepeat = now.Sub(g.lastRepeatTime)
	}

	// We don't want to send _too_ many repeat key presses, so space it out
	if timeSinceLastRepeat >= 100*time.Millisecond {
		input.PushEvent(g.inputEvent)
		g.lastRepeatTime = now
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
