// a2font provides fonts that can be used for the Apple II. These are
// based on an examination of the font used by the computer, so there
// may either be errors or minor discrepancies.
package a2font

import "github.com/pevans/erc/gfx"

const (
	sysFontWidth  uint = 7
	sysFontHeight uint = 8
)

type maskFunc func([]byte) []byte

// SystemFont returns a font object that contains all the glyphs of the Apple II
// system font
func SystemFont() *gfx.Font {
	f := gfx.NewFont(
		sysFontWidth,
		sysFontHeight,
	)

	fontUpperCase(f, 0x00, invert)
	fontSpecial(f, 0x20, invert)

	// TODO: these should be "flashing" characters
	fontUpperCase(f, 0x40, invert)
	fontSpecial(f, 0x60, invert)

	fontUpperCase(f, 0x80, nil)
	fontSpecial(f, 0xa0, nil)
	fontUpperCase(f, 0xc0, nil)
	fontLowerCase(f, 0xe0, nil)

	return f
}

func invert(b []byte) []byte {
	for i := range b {
		b[i] ^= 1
	}

	return b
}
