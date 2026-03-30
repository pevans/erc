package a2text

import (
	"image/color"
	"testing"

	"github.com/pevans/erc/a2/a2mono"
	"github.com/pevans/erc/gfx"
	"github.com/stretchr/testify/assert"
)

// transparent is the color written by an unlit glyph pixel. Font.Glyph
// returns a blank FrameBuffer for undefined characters, whose pixels are all
// zero (fully transparent), not opaque black.
var transparent = color.RGBA{}

// mockMem80 implements memoryGetter for 80-column text tests.
type mockMem80 struct {
	main [0x400]uint8
	aux  [0x400]uint8
}

func (m *mockMem80) GetMain(addr int) uint8 {
	if addr >= 0x400 && addr < 0x800 {
		return m.main[addr-0x400]
	}

	return 0
}

func (m *mockMem80) GetAux(addr int) uint8 {
	if addr >= 0x400 && addr < 0x800 {
		return m.aux[addr-0x400]
	}

	return 0
}

// makeTestFont80 builds a minimal 7x16 font for testing. Glyph 0x01 is all
// white; all other characters use the default transparent glyph.
func makeTestFont80() *gfx.Font {
	f := gfx.NewFont(7, 16)
	pts := make([]byte, 7*16)
	for i := range pts {
		pts[i] = 1
	}
	f.DefineGlyph(0x01, pts)

	return f
}

// TestRender80AuxGoesLeft verifies that the aux character is blitted to the
// left (even) screen column of the pair.
func TestRender80AuxGoesLeft(t *testing.T) {
	gfx.Screen = gfx.NewFrameBuffer(560, 384)

	font := makeTestFont80()
	mem := &mockMem80{}

	// $0400 -> row=0, col=0 -> leftX=0 (aux), rightX=7 (main)
	mem.aux[0] = 0x01  // all-white glyph at left column
	mem.main[0] = 0x00 // default (transparent) at right column

	Render80(mem, font, font, true, 0x400, 0x800, a2mono.None)

	assert.Equal(t, a2mono.White, gfx.Screen.GetPixel(0, 0), "aux char should render at left column")
	assert.Equal(t, transparent, gfx.Screen.GetPixel(7, 0), "main char should render at right column")
}

// TestRender80MainGoesRight verifies that the main character is blitted to
// the right (odd) screen column of the pair.
func TestRender80MainGoesRight(t *testing.T) {
	gfx.Screen = gfx.NewFrameBuffer(560, 384)

	font := makeTestFont80()
	mem := &mockMem80{}

	// $0400 -> row=0, col=0 -> leftX=0 (aux), rightX=7 (main)
	mem.aux[0] = 0x00  // default (transparent) at left column
	mem.main[0] = 0x01 // all-white glyph at right column

	Render80(mem, font, font, true, 0x400, 0x800, a2mono.None)

	assert.Equal(t, transparent, gfx.Screen.GetPixel(0, 0), "aux char should render at left column")
	assert.Equal(t, a2mono.White, gfx.Screen.GetPixel(7, 0), "main char should render at right column")
}

// TestRender80FlashAltFontUsedWhenFlashOff verifies that flashAltFont is
// selected when flashOn is false, and the primary font when flashOn is true.
func TestRender80FlashAltFontUsedWhenFlashOff(t *testing.T) {
	// primary: glyph 0x01 is white; flashAlt: glyph 0x01 undefined
	// (transparent).
	primary := makeTestFont80()
	flashAlt := gfx.NewFont(7, 16)

	mem := &mockMem80{}
	mem.aux[0] = 0x01

	t.Run("flashOff uses flashAltFont", func(t *testing.T) {
		gfx.Screen = gfx.NewFrameBuffer(560, 384)
		Render80(mem, primary, flashAlt, false, 0x400, 0x800, a2mono.None)
		assert.Equal(t, transparent, gfx.Screen.GetPixel(0, 0))
	})

	t.Run("flashOn uses primary font", func(t *testing.T) {
		gfx.Screen = gfx.NewFrameBuffer(560, 384)
		Render80(mem, primary, flashAlt, true, 0x400, 0x800, a2mono.None)
		assert.Equal(t, a2mono.White, gfx.Screen.GetPixel(0, 0))
	})
}

// TestRender80GreenScreen verifies that green monochrome recolors white
// pixels to green and leaves transparent pixels unchanged.
func TestRender80GreenScreen(t *testing.T) {
	gfx.Screen = gfx.NewFrameBuffer(560, 384)

	font := makeTestFont80()
	mem := &mockMem80{}
	mem.aux[0] = 0x01  // all-white glyph -> should become green
	mem.main[0] = 0x00 // default (transparent) -> stays transparent

	Render80(mem, font, font, true, 0x400, 0x800, a2mono.GreenScreen)

	assert.Equal(t, a2mono.Green, gfx.Screen.GetPixel(0, 0), "white pixels should become green")
	assert.Equal(t, transparent, gfx.Screen.GetPixel(7, 0), "transparent pixels should stay transparent")
}

// TestRender80AmberScreen verifies that amber monochrome recolors white
// pixels to amber.
func TestRender80AmberScreen(t *testing.T) {
	gfx.Screen = gfx.NewFrameBuffer(560, 384)

	font := makeTestFont80()
	mem := &mockMem80{}
	mem.aux[0] = 0x01

	Render80(mem, font, font, true, 0x400, 0x800, a2mono.AmberScreen)

	assert.Equal(t, a2mono.Amber, gfx.Screen.GetPixel(0, 0), "white pixels should become amber")
}

// TestRender80SkipsHoleBytes verifies that addresses with row=-1 or col=-1 in
// the lookup tables are silently skipped without writing to the screen. Hole
// bytes start at $0478 (offset 0x78).
func TestRender80SkipsHoleBytes(t *testing.T) {
	gfx.Screen = gfx.NewFrameBuffer(560, 384)

	font := makeTestFont80()
	mem := &mockMem80{}

	// $0478 (offset 0x78=120) is a hole: addressRows[120]=-1. Set it to white
	// in both banks; it should not appear on screen.
	mem.aux[0x78] = 0x01
	mem.main[0x78] = 0x01

	// Completes without panicking.
	Render80(mem, font, font, true, 0x400, 0x800, a2mono.None)

	// The last valid column in the $0470 range is col=39 at addr $0477
	// (offset 0x77, row=16). Because both banks are zero there, those pixels
	// are transparent. Checking these neighbors confirms the hole was
	// skipped.
	lastValidLeftX := uint(39 * 2 * 7)  // = 546
	lastValidRightX := uint(39*2+1) * 7 // = 553
	assert.Equal(t, transparent, gfx.Screen.GetPixel(lastValidLeftX, 16*16))
	assert.Equal(t, transparent, gfx.Screen.GetPixel(lastValidRightX, 16*16))
}

// TestRender80SecondColumn verifies column placement for the second byte
// position ($0401, col=1): aux at leftX=14, main at rightX=21.
func TestRender80SecondColumn(t *testing.T) {
	gfx.Screen = gfx.NewFrameBuffer(560, 384)

	font := makeTestFont80()
	mem := &mockMem80{}

	// $0401 -> row=0, col=1 -> leftX=14, rightX=21
	mem.aux[1] = 0x01
	mem.main[1] = 0x00

	Render80(mem, font, font, true, 0x400, 0x800, a2mono.None)

	assert.Equal(t, a2mono.White, gfx.Screen.GetPixel(14, 0), "aux at col 1 should be at x=14")
	assert.Equal(t, transparent, gfx.Screen.GetPixel(21, 0), "main at col 1 should be at x=21")
}
