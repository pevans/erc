package gfx

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// BorderOverlay draws a border around the screen edge.
type BorderOverlay struct {
	alpha        float64
	active       bool
	fadeRate     float64
	screenWidth  int
	screenHeight int
	pixel        *ebiten.Image
}

// PrefixOverlay is a global border overlay. We use it to flash an indicator
// when the user begins a keyboard shortcut.
var PrefixOverlay *BorderOverlay

func init() {
	PrefixOverlay = NewBorderOverlay()
}

// NewBorderOverlay creates a new border overlay.
func NewBorderOverlay() *BorderOverlay {
	pixel := ebiten.NewImage(1, 1)
	pixel.Fill(color.White)
	return &BorderOverlay{
		fadeRate: 1.0 / (borderFadeDuration * overlayTPS),
		pixel:    pixel,
	}
}

// Show activates the border overlay.
func (b *BorderOverlay) Show(screenWidth, screenHeight int) {
	b.screenWidth = screenWidth
	b.screenHeight = screenHeight
	b.alpha = 1.0
	b.active = true
}

// Hide immediately removes the border overlay.
func (b *BorderOverlay) Hide() {
	b.active = false
}

// Update handles the fade animation for the border.
func (b *BorderOverlay) Update() {
	if !b.active {
		return
	}

	b.alpha -= b.fadeRate
	if b.alpha <= 0 {
		b.alpha = 0
		b.active = false
	}
}

// Draw renders the border overlay.
func (b *BorderOverlay) Draw(screen *ebiten.Image) {
	if !b.active {
		return
	}

	a := float32(b.alpha)
	opts := &ebiten.DrawImageOptions{}
	opts.ColorScale.Scale(a, a, a, a)

	// Top edge
	opts.GeoM.Reset()
	opts.GeoM.Scale(float64(b.screenWidth), borderWidth)
	screen.DrawImage(b.pixel, opts)

	// Bottom edge
	opts.GeoM.Reset()
	opts.GeoM.Scale(float64(b.screenWidth), borderWidth)
	opts.GeoM.Translate(0, float64(b.screenHeight-borderWidth))
	screen.DrawImage(b.pixel, opts)

	// Left edge
	opts.GeoM.Reset()
	opts.GeoM.Scale(borderWidth, float64(b.screenHeight))
	screen.DrawImage(b.pixel, opts)

	// Right edge
	opts.GeoM.Reset()
	opts.GeoM.Scale(borderWidth, float64(b.screenHeight))
	opts.GeoM.Translate(float64(b.screenWidth-borderWidth), 0)
	screen.DrawImage(b.pixel, opts)
}

// IsActive returns whether the border overlay is visible.
func (b *BorderOverlay) IsActive() bool {
	return b.active
}
