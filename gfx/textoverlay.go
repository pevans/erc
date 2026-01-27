package gfx

import (
	"bytes"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

const (
	textOverlayFadeDuration = 3.0
	textOverlayTPS          = 60.0
	textOverlayMarginBottom = 40
)

// TextOverlay displays text on screen with fade animation.
type TextOverlay struct {
	text         string
	alpha        float64
	active       bool
	fadeRate     float64
	screenWidth  int
	screenHeight int
	faceSource   *text.GoTextFaceSource
}

// TextNotification is the global text overlay for disk swap notifications.
var TextNotification *TextOverlay

func init() {
	TextNotification = NewTextOverlay()
}

// NewTextOverlay creates a new text overlay.
func NewTextOverlay() *TextOverlay {
	faceSource, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		panic("failed to create font source: " + err.Error())
	}

	return &TextOverlay{
		fadeRate:   1.0 / (textOverlayFadeDuration * textOverlayTPS),
		faceSource: faceSource,
	}
}

// Show displays text on the screen with fade animation.
func (t *TextOverlay) Show(message string, width, height int) {
	t.text = message
	t.alpha = 1.0
	t.active = true
	t.screenWidth = width
	t.screenHeight = height
}

// Update handles the fade animation, called each frame.
func (t *TextOverlay) Update() {
	if !t.active {
		return
	}

	t.alpha -= t.fadeRate
	if t.alpha <= 0 {
		t.alpha = 0
		t.active = false
		t.text = ""
	}
}

// Draw renders the text overlay to the screen.
func (t *TextOverlay) Draw(screen *ebiten.Image) {
	if !t.active || t.text == "" {
		return
	}

	face := &text.GoTextFace{
		Source: t.faceSource,
		Size:   16,
	}

	// Measure text width for centering
	width, _ := text.Measure(t.text, face, 0)

	x := float64(t.screenWidth-int(width)) / 2
	y := float64(t.screenHeight - textOverlayMarginBottom)

	// Draw text directly to screen with alpha applied via ColorScale
	opts := &text.DrawOptions{}
	opts.GeoM.Translate(x, y)

	// Apply fade alpha
	a := float32(t.alpha)
	opts.ColorScale.Scale(a, a, a, a)

	text.Draw(screen, t.text, face, opts)
}
