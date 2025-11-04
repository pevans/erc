package a2font

import "github.com/pevans/erc/gfx"

// In 40-column text, the fonts we define actually need to be twice the size
// look correct on the screen.
func font40Glyph(original []byte) []byte {
	// This should _probably_ never happen...
	if len(original) != glyphWidth*glyphHeight {
		panic("7x8 glyph byte slice is required")
	}

	doubled := make([]byte, sysFont40Width*sysFont40Height)
	dpos := 0

	for row := 0; row < 8; row++ {
		rowStart := row * 7

		// We're going to copy the column row twice. We don't really need to
		// use i at all; we only need to run the column copy loop twice.
		for i := 0; i < 2; i++ {
			for col := 0; col < 7; col++ {
				dot := original[rowStart+col]

				doubled[dpos] = dot
				doubled[dpos+1] = dot

				dpos += 2
			}
		}
	}

	return doubled
}

// Given some font, define a glyph at the given offset, using a provided
// masking function and the underlying slice of bytes that represent the
// glyph.
func define40Glyph(font *gfx.Font, offset int, mask maskFunc, b []byte) {
	glyph := font40Glyph(b)

	if mask != nil {
		glyph = mask(glyph)
	}

	font.DefineGlyph(offset, glyph)
}

// SystemFont40 returns a font object that contains all the glyphs of the Apple II
// system font that is suitable for 40-column text
func SystemFont40() *gfx.Font {
	f := gfx.NewFont(
		sysFont40Width,
		sysFont40Height,
	)

	fontUpperCase(f, 0x00, invert, define40Glyph)
	fontSpecial(f, 0x20, invert, define40Glyph)

	// TODO: these should be "flashing" characters
	fontUpperCase(f, 0x40, invert, define40Glyph)
	fontSpecial(f, 0x60, invert, define40Glyph)

	fontUpperCase(f, 0x80, nil, define40Glyph)
	fontSpecial(f, 0xa0, nil, define40Glyph)
	fontUpperCase(f, 0xc0, nil, define40Glyph)
	fontLowerCase(f, 0xe0, nil, define40Glyph)

	return f
}
