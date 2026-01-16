// a2font provides fonts that can be used for the Apple II. These are based on
// an examination of the font used by the computer, so there may either be
// errors or minor discrepancies.
package a2font

import "github.com/pevans/erc/gfx"

const (
	// These are the dimensions of any glyph we define in this package,
	// regardless of how it might be rendered on screen. You can think of the
	// units for these numbers as "dots" as they'd be rendered on a screen.
	glyphWidth  = 7
	glyphHeight = 8

	// The dimensions of a font rendered for 40-column text, which are
	// effectively double the size of the original glyphs.
	sysFont40Width  = 14
	sysFont40Height = 16

	// The dimensions of a font rendered for 80-column text, which are double
	// the height but the same width as the original glyphs.
	sysFont80Width  = 7
	sysFont80Height = 16
)

type (
	maskFunc  func([]byte) []byte
	glyphFunc func(*gfx.Font, int, maskFunc, []byte)
)

// Apply a mask to the font so that the dots are inverted from their
// definition (e.g. instead of an white "A" rendered a black field, a black
// "A" rendered on a white field).
func invert(b []byte) []byte {
	for i := range b {
		b[i] ^= 1
	}

	return b
}
