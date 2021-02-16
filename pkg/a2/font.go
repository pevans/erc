package a2

import "github.com/pevans/erc/pkg/gfx"

const (
	sysFontWidth  uint = 7
	sysFontHeight uint = 8
)

// SystemFont returns a font object that contains all the glyphs of the Apple II
// system font
func SystemFont() *gfx.Font {
	f := gfx.NewFont(
		sysFontWidth,
		sysFontHeight,
	)

	font00(f)
	font20(f)
	font40(f)
	font60(f)

	return f
}
