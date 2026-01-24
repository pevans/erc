package a2video

import (
	"image/color"

	"github.com/pevans/erc/a2/a2mono"
	"github.com/pevans/erc/gfx"
	"github.com/pevans/erc/memory"
)

// recolorGlyph returns a framebuffer that modifies everything to match the
// provided monochrome color.
func recolorGlyph(glyph *gfx.FrameBuffer, monochromeColor color.RGBA) *gfx.FrameBuffer {
	recolored := gfx.NewFrameBuffer(glyph.Width, glyph.Height)

	for y := range uint(glyph.Height) {
		for x := range uint(glyph.Width) {
			pixel := glyph.GetPixel(x, y)

			if pixel.R == 255 && pixel.G == 255 && pixel.B == 255 {
				_ = recolored.SetCell(x, y, monochromeColor)
			} else {
				_ = recolored.SetCell(x, y, pixel)
			}
		}
	}

	return recolored
}

// RenderText will draw text in the framebuffer starting from a specific
// memory range, and ending at a specific memory range.
func RenderText(seg memory.Getter, font *gfx.Font, start, end int, monochromeMode int) {
	var monochromeColor color.RGBA
	useMonochrome := false

	switch monochromeMode {
	case a2mono.GreenScreen:
		monochromeColor = a2mono.Green
		useMonochrome = true
	case a2mono.AmberScreen:
		monochromeColor = a2mono.Amber
		useMonochrome = true
	}

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

		if useMonochrome {
			glyph = recolorGlyph(glyph, monochromeColor)
		}

		_ = gfx.Screen.Blit(x, y, glyph)
	}
}
