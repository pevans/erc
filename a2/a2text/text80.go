package a2text

import (
	"github.com/pevans/erc/gfx"
)

// memoryGetter provides access to both main and auxiliary memory banks for
// 80-column text rendering.
type memoryGetter interface {
	GetMain(addr int) uint8
	GetAux(addr int) uint8
}

// Render80 draws 80-column text from both memory banks. For each 40-byte
// position in the text page, the auxiliary byte is placed in the left (even)
// screen column and the main byte is placed in the right (odd) screen column.
// flashAltFont is used when flashOn is false.
func Render80(
	seg memoryGetter,
	font *gfx.Font,
	flashAltFont *gfx.Font,
	flashOn bool,
	start, end int,
	monochromeMode int,
) {
	monochromeColor, useMonochrome := monochromeSetup(monochromeMode)

	activeFont := font
	if !flashOn {
		activeFont = flashAltFont
	}

	for addr := start; addr < end; addr++ {
		offset := addr - start
		row := addressRows[offset]
		col := addressCols[offset]

		if row < 0 || col < 0 {
			continue
		}

		auxChar := seg.GetAux(addr)
		mainChar := seg.GetMain(addr)

		auxGlyph := activeFont.Glyph(int(auxChar))
		mainGlyph := activeFont.Glyph(int(mainChar))

		if useMonochrome {
			auxGlyph = recolorGlyph(auxGlyph, monochromeColor)
			mainGlyph = recolorGlyph(mainGlyph, monochromeColor)
		}

		leftX := uint(col*2) * activeFont.GlyphWidth
		rightX := uint(col*2+1) * activeFont.GlyphWidth
		y := uint(row) * activeFont.GlyphHeight

		_ = gfx.Screen.Blit(leftX, y, auxGlyph)
		_ = gfx.Screen.Blit(rightX, y, mainGlyph)
	}
}
