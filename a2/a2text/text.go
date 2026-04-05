package a2text

import (
	"image/color"

	"github.com/pevans/erc/a2/a2mono"
	"github.com/pevans/erc/gfx"
	"github.com/pevans/erc/memory"
)

// monochromeSetup returns the monochrome color and whether monochrome mode is
// active for the given mode constant.
func monochromeSetup(monochromeMode int) (color.RGBA, bool) {
	switch monochromeMode {
	case a2mono.GreenScreen:
		return a2mono.Green, true
	case a2mono.AmberScreen:
		return a2mono.Amber, true
	}

	return color.RGBA{}, false
}

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

// Render will draw text in the framebuffer starting from a specific memory
// range, and ending at a specific memory range. flashAltFont is used when
// flashOn is false (flash characters appear normal rather than inverse).
func Render(
	seg memory.Getter,
	font *gfx.Font,
	flashAltFont *gfx.Font,
	flashOn bool,
	start, end int,
	monochromeMode int,
	startRow int,
) {
	monochromeColor, useMonochrome := monochromeSetup(monochromeMode)

	activeFont := font
	if !flashOn {
		activeFont = flashAltFont
	}

	for addr := start; addr < end; addr++ {
		// Try to figure out where the text should be displayed
		row := addressRows[addr-start]
		col := addressCols[addr-start]

		// This address does not contain displayable data, and should be
		// skipped.
		if row < 0 || col < 0 {
			continue
		}

		if row < startRow {
			continue
		}

		// Convert the row and column into the framebuffer grid
		x := uint(col) * activeFont.GlyphWidth
		y := uint(row) * activeFont.GlyphHeight

		// Figure out what glyph to render
		char := seg.Get(int(addr))
		glyph := activeFont.Glyph(int(char))

		if useMonochrome {
			glyph = recolorGlyph(glyph, monochromeColor)
		}

		_ = gfx.Screen.Blit(x, y, glyph)
	}
}
