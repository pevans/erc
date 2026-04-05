package a2hires_test

import (
	"image/color"
	"testing"

	"github.com/pevans/erc/a2/a2hires"
	"github.com/pevans/erc/a2/a2mono"
	"github.com/pevans/erc/gfx"
	"github.com/stretchr/testify/assert"
)

// TestRenderMixedStopsAtRow160 verifies that mixed mode only renders the top
// 160 hires rows and leaves the bottom 32 rows untouched.
func TestRenderMixedStopsAtRow160(t *testing.T) {
	setupScreen()
	mem := &mockSegment{}

	// Fill a sentinel color into the bottom region of the screen so we can
	// verify it is not overwritten.
	sentinel := color.RGBA{0xAB, 0xCD, 0xEF, 0xFF}
	for y := uint(320); y < 384; y++ {
		for x := uint(0); x < 560; x++ {
			_ = gfx.Screen.SetCell(x, y, sentinel)
		}
	}

	// Fill all hires memory with on bits
	fillMemoryWithByte(mem, 0x7F)

	a2hires.Render(mem, 0, 0, a2mono.None, true)

	// The graphics area (rows 0-159, pixels 0-319) should be rendered
	nonBlackCount := 0
	for y := uint(0); y < 320; y += 20 {
		for x := uint(0); x < 560; x += 20 {
			pixel := gfx.Screen.GetPixel(x, y)
			if pixel != (color.RGBA{0, 0, 0, 0xff}) {
				nonBlackCount++
			}
		}
	}
	assert.Greater(t, nonBlackCount, 0, "graphics area should have rendered pixels")

	// The bottom region (pixels 320-383) should still have the sentinel
	for y := uint(320); y < 384; y += 8 {
		for x := uint(0); x < 560; x += 40 {
			pixel := gfx.Screen.GetPixel(x, y)
			assert.Equal(t, sentinel, pixel,
				"pixel at (%d,%d) should be untouched sentinel in mixed mode", x, y)
		}
	}
}

// TestRenderMixedMonochrome verifies that mixed mode works with monochrome
// rendering and still stops at row 160.
func TestRenderMixedMonochrome(t *testing.T) {
	setupScreen()
	mem := &mockSegment{}

	sentinel := color.RGBA{0xAB, 0xCD, 0xEF, 0xFF}
	for y := uint(320); y < 384; y++ {
		for x := uint(0); x < 560; x++ {
			_ = gfx.Screen.SetCell(x, y, sentinel)
		}
	}

	fillMemoryWithByte(mem, 0x7F)

	a2hires.Render(mem, 0, 0, a2mono.GreenScreen, true)

	// Graphics area should have green pixels
	green := a2mono.Green
	greenCount := 0
	for y := uint(0); y < 320; y += 20 {
		for x := uint(0); x < 560; x += 20 {
			if gfx.Screen.GetPixel(x, y) == green {
				greenCount++
			}
		}
	}
	assert.Greater(t, greenCount, 0, "graphics area should have green pixels")

	// Bottom region should be untouched
	for y := uint(320); y < 384; y += 8 {
		for x := uint(0); x < 560; x += 40 {
			pixel := gfx.Screen.GetPixel(x, y)
			assert.Equal(t, sentinel, pixel,
				"pixel at (%d,%d) should be untouched in mixed mode", x, y)
		}
	}
}

// TestRenderNonMixedRendersAllRows verifies that when mixed is false, all 192
// rows are rendered as before.
func TestRenderNonMixedRendersAllRows(t *testing.T) {
	setupScreen()
	mem := &mockSegment{}

	fillMemoryWithByte(mem, 0x7F)

	a2hires.Render(mem, 0, 0, a2mono.GreenScreen, false)

	green := a2mono.Green

	// Check that pixels in the bottom region (rows 160-191) are also rendered
	for y := uint(320); y < 384; y += 8 {
		for x := uint(0); x < 560; x += 40 {
			pixel := gfx.Screen.GetPixel(x, y)
			assert.Equal(t, green, pixel,
				"pixel at (%d,%d) should be rendered in non-mixed mode", x, y)
		}
	}
}

// TestRenderNonMixedRendersAllRowsColor verifies that when mixed is false,
// all 192 rows are rendered in color mode.
func TestRenderNonMixedRendersAllRowsColor(t *testing.T) {
	setupScreen()
	mem := &mockSegment{}

	fillMemoryWithByte(mem, 0x7F)

	a2hires.Render(mem, 0, 0, a2mono.None, false)

	// With all dots on in color mode, most pixels should be white
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
		"bottom rows should be rendered in non-mixed color mode")
}
