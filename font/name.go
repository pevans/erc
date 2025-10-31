package font

import (
	"fmt"
	"image"
)

// A Name is simply a const that represents the font. Exactly where
// and how fonts are stored in the system is an abstraction that is not
// exported from the font package.
type Name int

type info struct {
	byts []byte
	size image.Point
	mask int
}

const (
	// A2System is the "system font" for the Apple II. Although it was
	// posible to have a narrow-width font, it turns out that the source
	// font has exactly the same dimensions in both cases; the only
	// difference is that twice as many pixels are rendered within the
	// same physical space.
	A2System Name = iota

	// A2Inverted is essentially the same as the system font, except
	// that all of the pixels are inverted white-to-black. The rendered
	// effect is that each glyph would be shown as a silhouette on a
	// white background.
	A2Inverted

	maxFontName
)

// Each font const is keyed to a slice of bytes that represents the
// font.
var fontNames = map[Name]info{
	A2System: {
		byts: apple2SystemFont,
		size: image.Point{X: 8, Y: 9},
		mask: 0x7F,
	},

	A2Inverted: {
		byts: apple2InverseFont,
		size: image.Point{X: 8, Y: 9},
		mask: 0x7F,
	},
}

func fontInfo(fn Name) (info, error) {
	inf, ok := fontNames[fn]

	// If this happened, we probably did something wrong by defining a
	// font const but not adding something to fontNames.
	if !ok {
		return info{}, fmt.Errorf("unknown font: %+v", fn)
	}

	return inf, nil
}
