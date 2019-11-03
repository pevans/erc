package font

import (
	"fmt"
	"image"
)

// Glyph returns an image which is a subset of the total font graphic
// that contains only the text that is indicated by the given rune.
func (b *Bitmap) Glyph(ch rune) (image.Image, error) {
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

	return img, nil
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
