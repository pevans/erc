package a2text

import (
	"testing"

	"github.com/pevans/erc/a2/a2mono"
	"github.com/pevans/erc/gfx"
	"github.com/stretchr/testify/assert"
)

// mockMem implements memory.Getter for text render tests.
type mockMem struct {
	data [0x800]uint8
}

func (m *mockMem) Get(addr int) uint8 {
	if addr >= 0x400 && addr < 0x800 {
		return m.data[addr]
	}
	return 0
}

func (m *mockMem) Get16(addr int) uint16 {
	lo := uint16(m.Get(addr))
	hi := uint16(m.Get(addr + 1))
	return lo | (hi << 8)
}

// makeTestFont builds a minimal 14x16 font for testing. Glyph 0x01 is all
// white; other characters use the default transparent glyph.
func makeTestFont() *gfx.Font {
	f := gfx.NewFont(14, 16)
	pts := make([]byte, 14*16)
	for i := range pts {
		pts[i] = 1
	}
	f.DefineGlyph(0x01, pts)
	return f
}

// TestRenderStartRowZeroRendersAll verifies that startRow=0 renders all 24
// text rows.
func TestRenderStartRowZeroRendersAll(t *testing.T) {
	gfx.Screen = gfx.NewFrameBuffer(560, 384)
	font := makeTestFont()
	mem := &mockMem{}

	// Put the all-white glyph at row 0 col 0 ($0400)
	mem.data[0x400] = 0x01

	Render(mem, font, font, true, 0x400, 0x800, a2mono.None, 0)

	// Row 0, col 0 should have a white pixel
	assert.Equal(t, a2mono.White, gfx.Screen.GetPixel(0, 0),
		"row 0 should be rendered with startRow=0")
}

// TestRenderStartRowSkipsEarlierRows verifies that startRow=20 skips text
// rows 0-19 and only renders rows 20-23.
func TestRenderStartRowSkipsEarlierRows(t *testing.T) {
	gfx.Screen = gfx.NewFrameBuffer(560, 384)
	font := makeTestFont()
	mem := &mockMem{}

	// Put the all-white glyph at row 0 col 0 ($0400) and row 20 col 0
	// ($0650). Row 20 base address is $0650 in the text page.
	mem.data[0x400] = 0x01
	mem.data[0x650] = 0x01

	Render(mem, font, font, true, 0x400, 0x800, a2mono.None, 20)

	// Row 0 (pixel y=0) should NOT be rendered
	assert.Equal(t, transparent, gfx.Screen.GetPixel(0, 0),
		"row 0 should be skipped with startRow=20")

	// Row 20 (pixel y=320) should be rendered
	assert.Equal(t, a2mono.White, gfx.Screen.GetPixel(0, 320),
		"row 20 should be rendered with startRow=20")
}

// TestRenderStartRowMonochrome verifies that startRow works correctly with
// monochrome mode.
func TestRenderStartRowMonochrome(t *testing.T) {
	gfx.Screen = gfx.NewFrameBuffer(560, 384)
	font := makeTestFont()
	mem := &mockMem{}

	mem.data[0x400] = 0x01
	mem.data[0x650] = 0x01

	Render(mem, font, font, true, 0x400, 0x800, a2mono.GreenScreen, 20)

	// Row 0 should be skipped
	assert.Equal(t, transparent, gfx.Screen.GetPixel(0, 0),
		"row 0 should be skipped in monochrome mixed mode")

	// Row 20 should be rendered in green
	assert.Equal(t, a2mono.Green, gfx.Screen.GetPixel(0, 320),
		"row 20 should be rendered in green with startRow=20")
}
