package a2video

import (
	"github.com/pevans/erc/gfx"
	"github.com/pevans/erc/memory"
)

// RenderText will draw text in the framebuffer starting from a specific
// memory range, and ending at a specific memory range.
func RenderText(seg memory.Getter, font *gfx.Font, start, end int) {
	for addr := start; addr < end; addr++ {
		// Try to figure out where the text should be displayed
		row := textAddressRows[addr-start]
		col := textAddressCols[addr-start]

		// This address does not contain displayable data, and should be
		// skipped.
		if row < 0 || col < 0 {
			continue
		}

		// Convert the row and column into the framebuffer grid
		x := uint(col) * font.GlyphWidth
		y := uint(row) * font.GlyphHeight

		// Figure out what glyph to render
		char := seg.Get(int(addr))
		glyph := font.Glyph(int(char))

		_ = gfx.Screen.Blit(x, y, glyph)
	}
}
