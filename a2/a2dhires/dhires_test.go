package a2dhires_test

import (
	"image/color"
	"testing"

	"github.com/pevans/erc/a2/a2dhires"
	"github.com/pevans/erc/a2/a2mono"
	"github.com/pevans/erc/gfx"
	"github.com/stretchr/testify/assert"
)

// mockMemory implements the memoryGetter interface for testing.
type mockMemory struct {
	main [0x4000]uint8
	aux  [0x4000]uint8
}

func (m *mockMemory) GetMain(addr int) uint8 {
	if addr < 0x2000 || addr >= 0x4000 {
		return 0
	}

	return m.main[addr]
}

func (m *mockMemory) GetAux(addr int) uint8 {
	if addr < 0x2000 || addr >= 0x4000 {
		return 0
	}

	return m.aux[addr]
}

// setupScreen creates a fresh framebuffer for testing.
func setupScreen() {
	gfx.Screen = gfx.NewFrameBuffer(560, 384)
}

// setPattern sets up memory to produce a specific 4-bit color pattern. The
// pattern repeats across the entire screen.
//
// Pattern bits map to phases:
// - bit 3 -> phase 0
// - bit 2 -> phase 1
// - bit 1 -> phase 2
// - bit 0 -> phase 3
func setPattern(m *mockMemory, pattern uint8) {
	// Create a repeating bit pattern that produces our desired 4-bit pattern
	// at every dot position.
	for row := range 192 {
		baseAddr := rowAddresses[row]

		for memOffset := range 40 {
			addr := baseAddr + memOffset

			// For simplicity, we'll create a pattern where the same 4-bit
			// pattern repeats. Each byte holds 7 bits, and we want to align
			// with the 4-dot color cycle.
			//
			// The easiest approach is to set bits based on which phase each
			// bit position represents in the final dot array.

			// Aux byte goes to screen columns memOffset*2, contributing 7
			// dots starting at position (memOffset*2)*7
			auxStartDot := memOffset * 2 * 7
			auxByte := uint8(0)

			for bit := range 7 {
				dotPos := auxStartDot + bit
				phase := dotPos % 4
				bitInPattern := 3 - phase
				if pattern&(1<<bitInPattern) != 0 {
					auxByte |= (1 << bit)
				}
			}

			m.aux[addr] = auxByte

			// Main byte goes to screen columns memOffset*2+1, contributing 7
			// dots starting at position (memOffset*2+1)*7
			mainStartDot := (memOffset*2 + 1) * 7
			mainByte := uint8(0)

			for bit := range 7 {
				dotPos := mainStartDot + bit
				phase := dotPos % 4
				bitInPattern := 3 - phase
				if pattern&(1<<bitInPattern) != 0 {
					mainByte |= (1 << bit)
				}
			}

			m.main[addr] = mainByte
		}
	}
}

// rowAddresses copied from coords.go for test convenience.
var rowAddresses = []int{
	0x2000, 0x2400, 0x2800, 0x2C00, 0x3000, 0x3400, 0x3800, 0x3C00, // 0-7
	0x2080, 0x2480, 0x2880, 0x2C80, 0x3080, 0x3480, 0x3880, 0x3C80, // 8-15
	0x2100, 0x2500, 0x2900, 0x2D00, 0x3100, 0x3500, 0x3900, 0x3D00, // 16-23
	0x2180, 0x2580, 0x2980, 0x2D80, 0x3180, 0x3580, 0x3980, 0x3D80, // 24-31
	0x2200, 0x2600, 0x2A00, 0x2E00, 0x3200, 0x3600, 0x3A00, 0x3E00, // 32-39
	0x2280, 0x2680, 0x2A80, 0x2E80, 0x3280, 0x3680, 0x3A80, 0x3E80, // 40-47
	0x2300, 0x2700, 0x2B00, 0x2F00, 0x3300, 0x3700, 0x3B00, 0x3F00, // 48-55
	0x2380, 0x2780, 0x2B80, 0x2F80, 0x3380, 0x3780, 0x3B80, 0x3F80, // 56-63
	0x2028, 0x2428, 0x2828, 0x2C28, 0x3028, 0x3428, 0x3828, 0x3C28, // 64-71
	0x20A8, 0x24A8, 0x28A8, 0x2CA8, 0x30A8, 0x34A8, 0x38A8, 0x3CA8, // 72-79
	0x2128, 0x2528, 0x2928, 0x2D28, 0x3128, 0x3528, 0x3928, 0x3D28, // 80-87
	0x21A8, 0x25A8, 0x29A8, 0x2DA8, 0x31A8, 0x35A8, 0x39A8, 0x3DA8, // 88-95
	0x2228, 0x2628, 0x2A28, 0x2E28, 0x3228, 0x3628, 0x3A28, 0x3E28, // 96-103
	0x22A8, 0x26A8, 0x2AA8, 0x2EA8, 0x32A8, 0x36A8, 0x3AA8, 0x3EA8, // 104-111
	0x2328, 0x2728, 0x2B28, 0x2F28, 0x3328, 0x3728, 0x3B28, 0x3F28, // 112-119
	0x23A8, 0x27A8, 0x2BA8, 0x2FA8, 0x33A8, 0x37A8, 0x3BA8, 0x3FA8, // 120-127
	0x2050, 0x2450, 0x2850, 0x2C50, 0x3050, 0x3450, 0x3850, 0x3C50, // 128-135
	0x20D0, 0x24D0, 0x28D0, 0x2CD0, 0x30D0, 0x34D0, 0x38D0, 0x3CD0, // 136-143
	0x2150, 0x2550, 0x2950, 0x2D50, 0x3150, 0x3550, 0x3950, 0x3D50, // 144-151
	0x21D0, 0x25D0, 0x29D0, 0x2DD0, 0x31D0, 0x35D0, 0x39D0, 0x3DD0, // 152-159
	0x2250, 0x2650, 0x2A50, 0x2E50, 0x3250, 0x3650, 0x3A50, 0x3E50, // 160-167
	0x22D0, 0x26D0, 0x2AD0, 0x2ED0, 0x32D0, 0x36D0, 0x3AD0, 0x3ED0, // 168-175
	0x2350, 0x2750, 0x2B50, 0x2F50, 0x3350, 0x3750, 0x3B50, 0x3F50, // 176-183
	0x23D0, 0x27D0, 0x2BD0, 0x2FD0, 0x33D0, 0x37D0, 0x3BD0, 0x3FD0, // 184-191
}

// TestRenderCompositeColors tests that each of the 16 composite colors
// renders correctly.
func TestRenderCompositeColors(t *testing.T) {
	colors := []struct {
		pattern uint8
		name    string
		color   color.RGBA
	}{
		{0x0, "Black", color.RGBA{0x00, 0x00, 0x00, 0xff}},
		{0x1, "Magenta", color.RGBA{0x90, 0x17, 0x40, 0xff}},
		{0x2, "Brown", color.RGBA{0x40, 0x54, 0x00, 0xff}},
		{0x3, "Orange", color.RGBA{0xd0, 0x6a, 0x1a, 0xff}},
		{0x4, "Dark Green", color.RGBA{0x00, 0x69, 0x40, 0xff}},
		{0x5, "Gray 1", color.RGBA{0x80, 0x80, 0x80, 0xff}},
		{0x6, "Green", color.RGBA{0x2f, 0xbc, 0x1a, 0xff}},
		{0x7, "Yellow", color.RGBA{0xbf, 0xd3, 0x5a, 0xff}},
		{0x8, "Dark Blue", color.RGBA{0x40, 0x2c, 0xa5, 0xff}},
		{0x9, "Purple", color.RGBA{0xd0, 0x43, 0xe5, 0xff}},
		{0xA, "Gray 2", color.RGBA{0x80, 0x80, 0x80, 0xff}},
		{0xB, "Pink", color.RGBA{0xff, 0x96, 0xbf, 0xff}},
		{0xC, "Medium Blue", color.RGBA{0x2f, 0x95, 0xe5, 0xff}},
		{0xD, "Light Blue", color.RGBA{0xbf, 0xab, 0xff, 0xff}},
		{0xE, "Aqua", color.RGBA{0x6f, 0xe8, 0xbf, 0xff}},
		{0xF, "White", color.RGBA{0xff, 0xff, 0xff, 0xff}},
	}

	for _, tc := range colors {
		t.Run(tc.name, func(t *testing.T) {
			setupScreen()

			mem := &mockMemory{}
			setPattern(mem, tc.pattern)

			a2dhires.Render(mem, a2mono.None)

			// Sample several pixels across the screen to verify the color We
			// check multiple rows and columns to ensure consistency
			for y := uint(0); y < 384; y += 50 {
				for x := uint(0); x < 560; x += 50 {
					pixel := gfx.Screen.GetPixel(x, y)
					assert.Equal(t, tc.color, pixel,
						"pixel at (%d,%d) should be %s", x, y, tc.name)
				}
			}
		})
	}
}

// TestRenderMonochromeGreen tests green monochrome rendering.
func TestRenderMonochromeGreen(t *testing.T) {
	setupScreen()
	mem := &mockMemory{}

	// Create a pattern with some bits on and some off. We'll use a
	// checkerboard-like pattern for variety.
	for row := range 192 {
		baseAddr := rowAddresses[row]
		for memOffset := range 40 {
			addr := baseAddr + memOffset
			// Alternate between 0xAA (10101010) and 0x55 (01010101)
			if memOffset%2 == 0 {
				mem.aux[addr] = 0xAA
				mem.main[addr] = 0x55
			} else {
				mem.aux[addr] = 0x55
				mem.main[addr] = 0xAA
			}
		}
	}

	a2dhires.Render(mem, a2mono.GreenScreen)

	green := a2mono.Green
	black := a2mono.Black

	// Check that we see both green and black pixels (not all one color)
	greenCount := 0
	blackCount := 0

	for y := uint(0); y < 384; y += 10 {
		for x := uint(0); x < 560; x += 10 {
			pixel := gfx.Screen.GetPixel(x, y)
			switch pixel {
			case green:
				greenCount++
			case black:
				blackCount++
			default:
				t.Errorf("unexpected color at (%d,%d): got %+v, want green or black", x, y, pixel)
			}
		}
	}

	assert.Greater(t, greenCount, 0, "should have some green pixels")
	assert.Greater(t, blackCount, 0, "should have some black pixels")
}

// TestRenderMonochromeAmber tests amber monochrome rendering.
func TestRenderMonochromeAmber(t *testing.T) {
	setupScreen()
	mem := &mockMemory{}

	// Create a pattern with some bits on and some off.
	for row := range 192 {
		baseAddr := rowAddresses[row]
		for memOffset := range 40 {
			addr := baseAddr + memOffset
			// Alternate between 0xAA (10101010) and 0x55 (01010101)
			if memOffset%2 == 0 {
				mem.aux[addr] = 0xAA
				mem.main[addr] = 0x55
			} else {
				mem.aux[addr] = 0x55
				mem.main[addr] = 0xAA
			}
		}
	}

	a2dhires.Render(mem, a2mono.AmberScreen)

	amber := a2mono.Amber
	black := a2mono.Black

	// Check that we see both amber and black pixels (not all one color)
	amberCount := 0
	blackCount := 0

	for y := uint(0); y < 384; y += 10 {
		for x := uint(0); x < 560; x += 10 {
			pixel := gfx.Screen.GetPixel(x, y)
			switch pixel {
			case amber:
				amberCount++
			case black:
				blackCount++
			default:
				t.Errorf("unexpected color at (%d,%d): got %+v, want amber or black", x, y, pixel)
			}
		}
	}

	assert.Greater(t, amberCount, 0, "should have some amber pixels")
	assert.Greater(t, blackCount, 0, "should have some black pixels")
}

// TestRenderMonochromeAllOn tests monochrome with all bits on.
func TestRenderMonochromeAllOn(t *testing.T) {
	setupScreen()
	mem := &mockMemory{}

	// Set all bits to 1
	for row := range 192 {
		baseAddr := rowAddresses[row]
		for memOffset := range 40 {
			addr := baseAddr + memOffset
			mem.aux[addr] = 0x7F  // All 7 bits on
			mem.main[addr] = 0x7F // All 7 bits on
		}
	}

	a2dhires.Render(mem, a2mono.GreenScreen)

	green := a2mono.Green

	// All pixels should be green
	for y := uint(0); y < 384; y += 20 {
		for x := uint(0); x < 560; x += 20 {
			pixel := gfx.Screen.GetPixel(x, y)
			assert.Equal(t, green, pixel, "all pixels should be green when all bits are on")
		}
	}
}

// TestRenderMonochromeAllOff tests monochrome with all bits off.
func TestRenderMonochromeAllOff(t *testing.T) {
	setupScreen()
	mem := &mockMemory{}

	// All memory is zero by default

	a2dhires.Render(mem, a2mono.AmberScreen)

	black := a2mono.Black

	// All pixels should be black
	for y := uint(0); y < 384; y += 20 {
		for x := uint(0); x < 560; x += 20 {
			pixel := gfx.Screen.GetPixel(x, y)
			assert.Equal(t, black, pixel, "all pixels should be black when all bits are off")
		}
	}
}
