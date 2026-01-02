package gfx

import (
	"bytes"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	overlayMargin       = 10
	overlaySize         = 32
	overlayFadeDuration = 3.0
	overlayTPS          = 60.0

	borderWidth        = 2
	borderFadeDuration = 5.0
)

// Overlay manages a PNG image that can be displayed on screen.
type Overlay struct {
	image    *ebiten.Image
	alpha    float64
	active   bool
	fade     bool
	fadeRate float64
	scale    float64
}

// NewOverlay creates a new overlay manager.
func NewOverlay() *Overlay {
	return &Overlay{
		fadeRate: 1.0 / (overlayFadeDuration * overlayTPS),
	}
}

// Show displays a PNG image (which it assumes are the provided bytes). The
// image will be positioned in the top-left corner and scaled to fit within
// the overlay size.
func (o *Overlay) Show(data []byte) error {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return err
	}

	o.image = ebiten.NewImageFromImage(img)
	o.alpha = 1.0
	o.active = true
	o.fade = true

	o.calculateScale()

	return nil
}

// ShowPersistent displays a PNG image without fading. The overlay will remain
// visible until Hide() is called.
func (o *Overlay) ShowPersistent(data []byte) error {
	if err := o.Show(data); err != nil {
		return err
	}
	o.fade = false
	return nil
}

// Hide immediately removes the overlay from the screen.
func (o *Overlay) Hide() {
	o.active = false
	o.image = nil
}

// calculateScale determines the scale factor to fit within the overlay size
// while maintaining aspect ratio.
func (o *Overlay) calculateScale() {
	if o.image == nil {
		return
	}

	imgBounds := o.image.Bounds()
	imgWidth := float64(imgBounds.Dx())
	imgHeight := float64(imgBounds.Dy())

	scaleX := overlaySize / imgWidth
	scaleY := overlaySize / imgHeight

	o.scale = scaleX
	if scaleY < scaleX {
		o.scale = scaleY
	}
}

// Update should be called each frame to handle the fade animation.
func (o *Overlay) Update() {
	if !o.active || !o.fade {
		return
	}

	o.alpha -= o.fadeRate
	if o.alpha <= 0 {
		o.alpha = 0
		o.active = false
		o.image = nil
	}
}

// Draw renders the overlay to the screen with the current alpha.
func (o *Overlay) Draw(screen *ebiten.Image) {
	if !o.active || o.image == nil {
		return
	}

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(o.scale, o.scale)
	opts.GeoM.Translate(overlayMargin, overlayMargin)

	a := float32(o.alpha)
	opts.ColorScale.Scale(a, a, a, a)

	screen.DrawImage(o.image, opts)
}

// IsActive returns whether the overlay is currently visible.
func (o *Overlay) IsActive() bool {
	return o.active
}
