package font

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

var (
	apple40Face font.Face
	apple80Face font.Face
)

func newFontFace(bytes []byte) font.Face {
	// If this doesn't parse, then I (specifically) did something wrong
	fon, err := truetype.Parse(bytes)
	if err != nil {
		panic(err)
	}

	return truetype.NewFace(fon, &truetype.Options{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}

func Apple40() font.Face {
	if apple40Face == nil {
		apple40Face = newFontFace(apple40Col)
	}

	return apple40Face
}

func Apple80() font.Face {
	if apple80Face == nil {
		apple80Face = newFontFace(apple80Col)
	}

	return apple80Face
}
