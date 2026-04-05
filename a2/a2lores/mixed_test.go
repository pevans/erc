package a2lores_test

import (
	"image/color"
	"testing"

	"github.com/pevans/erc/a2/a2lores"
	"github.com/pevans/erc/a2/a2mono"
	"github.com/pevans/erc/gfx"
	"github.com/stretchr/testify/assert"
)

// mockSegment implements memory.Getter for testing.
type mockSegment struct {
	data [0x800]uint8
}

func (m *mockSegment) Get(addr int) uint8 {
	if addr >= 0x400 && addr < 0x800 {
		return m.data[addr]
	}
	return 0
}

func (m *mockSegment) Get16(addr int) uint16 {
	lo := uint16(m.Get(addr))
	hi := uint16(m.Get(addr + 1))
	return lo | (hi << 8)
}

func setupScreen() {
	gfx.Screen = gfx.NewFrameBuffer(560, 384)
}

// TestRenderMixedSkipsBottomRows verifies that mixed mode skips lores rows
// 40-47, leaving the bottom portion of the screen untouched.
func TestRenderMixedSkipsBottomRows(t *testing.T) {
	setupScreen()
	mem := &mockSegment{}

	// Fill all text/lores memory with a non-zero pattern (color 5 top, color
	// 10 bottom -- both are gray)
	for addr := 0x400; addr < 0x800; addr++ {
		mem.data[addr] = 0xA5
	}

	// Fill the bottom region with a sentinel color
	sentinel := color.RGBA{0xAB, 0xCD, 0xEF, 0xFF}
	for y := uint(320); y < 384; y++ {
		for x := uint(0); x < 560; x++ {
			_ = gfx.Screen.SetCell(x, y, sentinel)
		}
	}

	a2lores.Render(mem, 0x400, 0x800, a2mono.None, true)

	// The graphics area (lores rows 0-39, pixel rows 0-319) should be
	// rendered with the gray color
	gray := color.RGBA{0x80, 0x80, 0x80, 0xff}
	grayCount := 0
	for y := uint(0); y < 320; y += 20 {
		for x := uint(0); x < 560; x += 20 {
			if gfx.Screen.GetPixel(x, y) == gray {
				grayCount++
			}
		}
	}
	assert.Greater(t, grayCount, 0, "graphics area should have gray pixels")

	// The bottom region (pixel rows 320-383) should still have the sentinel
	// color
	for y := uint(320); y < 384; y += 8 {
		for x := uint(0); x < 560; x += 40 {
			pixel := gfx.Screen.GetPixel(x, y)
			assert.Equal(t, sentinel, pixel,
				"pixel at (%d,%d) should be untouched sentinel in mixed mode", x, y)
		}
	}
}

// TestRenderNonMixedRendersAllRows verifies that when mixed is false, all 48
// lores rows are rendered.
func TestRenderNonMixedRendersAllRows(t *testing.T) {
	setupScreen()
	mem := &mockSegment{}

	for addr := 0x400; addr < 0x800; addr++ {
		mem.data[addr] = 0xFF
	}

	a2lores.Render(mem, 0x400, 0x800, a2mono.None, false)

	// White (color 15) should appear in the bottom region too
	white := color.RGBA{0xff, 0xff, 0xff, 0xff}
	whiteCount := 0
	for y := uint(320); y < 384; y += 8 {
		for x := uint(0); x < 560; x += 40 {
			if gfx.Screen.GetPixel(x, y) == white {
				whiteCount++
			}
		}
	}
	assert.Greater(t, whiteCount, 0,
		"bottom rows should be rendered in non-mixed mode")
}

// TestRenderMixedMonochrome verifies mixed mode works with monochrome.
func TestRenderMixedMonochrome(t *testing.T) {
	setupScreen()
	mem := &mockSegment{}

	for addr := 0x400; addr < 0x800; addr++ {
		mem.data[addr] = 0xFF
	}

	sentinel := color.RGBA{0xAB, 0xCD, 0xEF, 0xFF}
	for y := uint(320); y < 384; y++ {
		for x := uint(0); x < 560; x++ {
			_ = gfx.Screen.SetCell(x, y, sentinel)
		}
	}

	a2lores.Render(mem, 0x400, 0x800, a2mono.GreenScreen, true)

	// Bottom region should be untouched
	for y := uint(320); y < 384; y += 8 {
		for x := uint(0); x < 560; x += 40 {
			pixel := gfx.Screen.GetPixel(x, y)
			assert.Equal(t, sentinel, pixel,
				"pixel at (%d,%d) should be untouched in mixed monochrome mode", x, y)
		}
	}
}
