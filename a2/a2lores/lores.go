package a2lores

import (
	"image/color"

	"github.com/pevans/erc/a2/a2mono"
	"github.com/pevans/erc/a2/a2video"
	"github.com/pevans/erc/gfx"
	"github.com/pevans/erc/memory"
)

const (
	blockWidth  = 14
	blockHeight = 8
)

// Color palette for color mode rendering. Although there are 16 color blocks,
// technically Gray1 and Gray2 are the same color mask.
var colors = []*gfx.FrameBuffer{
	/* black      */ newBlock(color.RGBA{0x00, 0x00, 0x00, 0xff}),
	/* magenta    */ newBlock(color.RGBA{0x90, 0x17, 0x40, 0xff}),
	/* darkBlue   */ newBlock(color.RGBA{0x40, 0x2c, 0xa5, 0xff}),
	/* purple     */ newBlock(color.RGBA{0xd0, 0x43, 0xe5, 0xff}),
	/* darkGreen  */ newBlock(color.RGBA{0x00, 0x69, 0x40, 0xff}),
	/* gray1      */ newBlock(color.RGBA{0x80, 0x80, 0x80, 0xff}),
	/* mediumBlue */ newBlock(color.RGBA{0x2f, 0x95, 0xe5, 0xff}),
	/* lightBlue  */ newBlock(color.RGBA{0xbf, 0xab, 0xff, 0xff}),
	/* brown      */ newBlock(color.RGBA{0x40, 0x54, 0x00, 0xff}),
	/* orange     */ newBlock(color.RGBA{0xd0, 0x6a, 0x1a, 0xff}),
	/* gray2      */ newBlock(color.RGBA{0x80, 0x80, 0x80, 0xff}),
	/* pink       */ newBlock(color.RGBA{0xff, 0x96, 0xbf, 0xff}),
	/* lightGreen */ newBlock(color.RGBA{0x2f, 0xbc, 0x1a, 0xff}),
	/* yellow     */ newBlock(color.RGBA{0xbf, 0xd3, 0x5a, 0xff}),
	/* aquamarine */ newBlock(color.RGBA{0x6f, 0xe8, 0xbf, 0xff}),
	/* white      */ newBlock(color.RGBA{0xff, 0xff, 0xff, 0xff}),
}

const (
	shadeLight  = iota // Light colors (75% intensity)
	shadeMedium        // Medium colors (50% intensity)
	shadeDark          // Dark colors (25% intensity)
)

var colorShades = []int{
	shadeMedium, // Black (special case, always black)
	shadeMedium, // Magenta
	shadeDark,   // Dark Blue
	shadeMedium, // Purple
	shadeDark,   // Dark Green
	shadeMedium, // Gray1
	shadeMedium, // Medium Blue
	shadeLight,  // Light Blue
	shadeMedium, // Brown
	shadeMedium, // Orange
	shadeMedium, // Gray2
	shadeMedium, // Pink
	shadeLight,  // Light Green
	shadeMedium, // Yellow
	shadeMedium, // Aquamarine
	shadeMedium, // White (special case, always full monochrome color)
}

var (
	monochromeGreenBlocks [16]*gfx.FrameBuffer
	monochromeAmberBlocks [16]*gfx.FrameBuffer
)

func init() {
	for i := range 16 {
		monochromeGreenBlocks[i] = newMonochromeBlock(uint8(i), a2mono.Green)
		monochromeAmberBlocks[i] = newMonochromeBlock(uint8(i), a2mono.Amber)
	}
}

// shadedMonochromeColor returns the version of a given color but with some
// percent-shading applied to it.
func shadedMonochromeColor(baseColor color.RGBA, intensity float64) color.RGBA {
	return color.RGBA{
		R: uint8(float64(baseColor.R) * intensity),
		G: uint8(float64(baseColor.G) * intensity),
		B: uint8(float64(baseColor.B) * intensity),
		A: 0xff,
	}
}

// newMonochromeBlock returns a framebuffer that is filled with the monochrome
// color that applies to the normal lo-res color that would be displayed. That
// is, each color that isn't white or black would be shown as a darker hue of
// the monochrome color -- the question, "how dark?", is answered based on the
// whether that color is indexed as light, dark, or no qualifier.
func newMonochromeBlock(bitPattern uint8, monochromeColor color.RGBA) *gfx.FrameBuffer {
	var clr color.RGBA

	index := bitPattern & 0xf

	switch index {
	case 0:
		clr = color.RGBA{0x00, 0x00, 0x00, 0xff}
	case 15:
		clr = monochromeColor
	default:
		shade := colorShades[index]
		switch shade {
		case shadeLight:
			clr = shadedMonochromeColor(monochromeColor, 0.75)
		case shadeMedium:
			clr = shadedMonochromeColor(monochromeColor, 0.50)
		case shadeDark:
			clr = shadedMonochromeColor(monochromeColor, 0.25)
		}
	}

	return newBlock(clr)
}

// newBlock returns a new framebuffer filled with the provided lo-res color.
func newBlock(clr color.RGBA) *gfx.FrameBuffer {
	fbuf := gfx.NewFrameBuffer(blockWidth, blockHeight)

	for y := range uint(blockHeight) {
		for x := range uint(blockWidth) {
			err := fbuf.SetCell(x, y, clr)
			if err != nil {
				panic(err)
			}
		}
	}

	return fbuf
}

// Block returns a color rectangle that matches the color suggested by the
// given pattern of bits
func Block(bitPattern uint8, monochromeMode int) *gfx.FrameBuffer {
	index := bitPattern & 0xf

	switch monochromeMode {
	case a2mono.GreenScreen:
		return monochromeGreenBlocks[index]
	case a2mono.AmberScreen:
		return monochromeAmberBlocks[index]
	default:
		return colors[index]
	}
}

// Render takes the data in the lo-res display buffer (essentially the text
// page) and writes that to the Screen framebuffer.
func Render(seg memory.Getter, start, end int, monochromeMode int) {
	for addr := start; addr < end; addr++ {
		row := a2video.LoresAddressRows[addr-start]
		col := a2video.LoresAddressCols[addr-start]

		if row < 0 || col < 0 {
			continue
		}

		x := uint(col) * blockWidth
		y := uint(row) * blockHeight

		byt := seg.Get(int(addr))

		// The Apple IIe technical reference (p. 22) states that we should
		// show the low-order nibble in the top row, and high-order nibble in
		// the bottom row.
		_ = gfx.Screen.Blit(x, y, Block(byt&0xf, monochromeMode))
		_ = gfx.Screen.Blit(x, y+8, Block(byt>>4, monochromeMode))
	}
}
