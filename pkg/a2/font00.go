package a2

import "github.com/pevans/erc/pkg/gfx"

func font00(f *gfx.Font) {
	for ch := 0x0; ch < 0x20; ch++ {
		f.DefineGlyphAsBuffer(ch, f.Glyph(ch+0x40).Invert())
	}
}
