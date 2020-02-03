package font

import (
	"bytes"
	"image"
	"image/png"

	"github.com/hajimehoshi/ebiten"
	"github.com/pkg/errors"
)

// A Bitmap is a representation of a font. In practice, fonts are PNG
// images which are a set of rendered characters, or glyphs. This also
// means that fonts are essentially bitmapped, and they can only scale
// in accordance to the aspect ratio of the display.
type Bitmap struct {
	// The source image of the font bitmap.
	img *ebiten.Image

	// This is the size of an individual glyph, which in a bitmap must
	// be uniform across all glyphs.
	size image.Point

	// Not every bitmap can uniformly draw every character you imagine;
	// the mask ensures that there is some limit to what we would
	// allow.
	mask int

	submap map[rune]*Glyph
}

// NewBitmap returns a new font type, from which you can render text on
// the screen. The possible fonts for this function are defined as
// consts.
func NewBitmap(fn Name) (*Bitmap, error) {
	// Each font is hard-coded into the binary as a slice of bytes.
	info, err := fontInfo(fn)
	if err != nil {
		return nil, err
	}

	// Which is then decoded.
	pngImage, err := png.Decode(bytes.NewReader(info.byts))
	if err != nil {
		return nil, errors.Wrapf(err, "font %+v", fn)
	}

	ebiImage, err := ebiten.NewImageFromImage(pngImage, ebiten.FilterLinear)
	if err != nil {
		return nil, errors.Wrapf(err, "font %+v", fn)
	}

	return &Bitmap{
		img:    ebiImage,
		size:   info.size,
		mask:   info.mask,
		submap: make(map[rune]*Glyph),
	}, nil
}
