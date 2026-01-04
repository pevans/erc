package gfx

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
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
}

// DiskNotification is the global text overlay for disk swap notifications.
var DiskNotification *TextOverlay

func init() {
	DiskNotification = NewTextOverlay()
}

// NewTextOverlay creates a new text overlay.
func NewTextOverlay() *TextOverlay {
	return &TextOverlay{
		fadeRate: 1.0 / (textOverlayFadeDuration * textOverlayTPS),
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

	face := basicfont.Face7x13

	// Measure text width for centering
	bounds, _ := font.BoundString(face, t.text)
	textWidth := (bounds.Max.X - bounds.Min.X).Ceil()

	x := (t.screenWidth - textWidth) / 2
	y := t.screenHeight - textOverlayMarginBottom

	// Create a temporary image to render text with alpha
	textImg := ebiten.NewImage(textWidth+4, 20)
	text.Draw(textImg, t.text, face, 2, 14, color.White)

	// Draw the text image with alpha
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(x), float64(y))

	a := float32(t.alpha)
	opts.ColorScale.Scale(a, a, a, a)

	screen.DrawImage(textImg, opts)
}
