package a2font

import "github.com/pevans/erc/gfx"

func font80Glyph(original []byte) []byte {
	if len(original) != glyphWidth*glyphHeight {
		panic("7x8 glyph byte slice is required")
	}

	doubled := make([]byte, sysFont80Width*sysFont80Height)
	dpos := 0

	for row := range 8 {
		rowStart := row * 7

		// We're going to copy each row twice to double the height
		for range 2 {
			for col := range 7 {
				dot := original[rowStart+col]

				doubled[dpos] = dot
				dpos++
			}
		}
	}

	return doubled
}

func define80Glyph(font *gfx.Font, offset int, mask maskFunc, b []byte) {
	glyph := font80Glyph(b)

	if mask != nil {
		glyph = mask(glyph)
	}

	font.DefineGlyph(offset, glyph)
}

// SystemFont80 returns a font object that contains all the glyphs of the
// Apple II system font that is suitable for 80-column text
func SystemFont80() *gfx.Font {
	f := gfx.NewFont(
		sysFont80Width,
		sysFont80Height,
	)

	fontUpperCase(f, 0x00, invert, define80Glyph)
	fontSpecial(f, 0x20, invert, define80Glyph)

	// TODO: these should be "flashing" characters
	fontUpperCase(f, 0x40, invert, define80Glyph)
	fontSpecial(f, 0x60, invert, define80Glyph)

	fontUpperCase(f, 0x80, nil, define80Glyph)
	fontSpecial(f, 0xa0, nil, define80Glyph)
	fontUpperCase(f, 0xc0, nil, define80Glyph)
	fontLowerCase(f, 0xe0, nil, define80Glyph)

	return f
}
