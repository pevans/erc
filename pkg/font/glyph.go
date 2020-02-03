package font

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten"
)

type Glyph struct {
	Image image.Image
}

// Glyph returns an image which is a subset of the total font graphic
// that contains only the text that is indicated by the given rune.
func (b *Bitmap) NewGlyph(ch rune) (*Glyph, error) {
	// If we already have this glyph in our subimage map, just return
	// that.
	g, ok := b.submap[ch]
	if ok {
		return g, nil
	}

	// mask is essentially the integer bitmask of what's allowed; but if
	// simply compare ch to the mask, we can quickly say that the given
	// rune is non-renderable.
	if int(ch) > b.mask {
		return nil, fmt.Errorf("non-renderable character: %v", ch)
	}

	// offset is where we would find the glyph in the bitmap.
	offset := b.offset(ch)

	// rectangle represents the exact dimensions of the glyph itself;
	// Min would be the top-left coordinate of the glyph, and Max would
	// be the bottom-right.
	rect := image.Rectangle{
		Min: offset,
		Max: image.Point{
			X: offset.X + b.size.X,
			Y: offset.Y + b.size.Y,
		},
	}

	img := b.img.SubImage(rect)

	// Because ebiten does not return an error here, we need to double
	// check if the returned img is nil.
	if img == nil {
		return nil, fmt.Errorf("unable to acquire SubImage for Glyph: %+v", rect)
	}

	g = &Glyph{Image: img}
	b.submap[ch] = g

	return g, nil
}

func row(ch rune) int {
	return (int(ch) & 0xfffffff0) >> 4
}

func col(ch rune) int {
	return int(ch) & 0x0f
}

func (b *Bitmap) offset(ch rune) image.Point {
	return image.Point{
		X: col(ch) * b.size.X,
		Y: row(ch) * b.size.Y,
	}
}

func (g *Glyph) Draw(coord image.Point) {
	var (
		op ebiten.DrawImageOptions
		fx = float64(coord.X)
		fy = float64(coord.Y)
	)

	op.GeoM.Translate(fx, fy)
}
