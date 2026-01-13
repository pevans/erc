package a2video

import (
	"image/color"

	"github.com/pevans/erc/gfx"
	"github.com/pevans/erc/memory"
)

// Low resolution graphics are rendered in "blocks", 7 dots wide by 4
// dots high. Each byte in Display Page 1 holds two blocks, which are
// rendered top-to-bottom. Unlike text mode, low resolution mode's
// column spans in memory encode two rows at a time.

const (
	loresBlockWidth  = 14
	loresBlockHeight = 8
)

// Although there are 16 color blocks, technically Gray1 and Gray2
// are the same color mask
var loresColors = []*gfx.FrameBuffer{
	/* loresBlockBlack      */ newLoresBlock(color.RGBA{0x00, 0x00, 0x00, 0xff}),
	/* loresBlockMagenta    */ newLoresBlock(color.RGBA{0x90, 0x17, 0x40, 0xff}),
	/* loresBlockDarkBlue   */ newLoresBlock(color.RGBA{0x40, 0x2c, 0xa5, 0xff}),
	/* loresBlockPurple     */ newLoresBlock(color.RGBA{0xd0, 0x43, 0xe5, 0xff}),
	/* loresBlockDarkGreen  */ newLoresBlock(color.RGBA{0x00, 0x69, 0x40, 0xff}),
	/* loresBlockGray1      */ newLoresBlock(color.RGBA{0x80, 0x80, 0x80, 0xff}),
	/* loresBlockMediumBlue */ newLoresBlock(color.RGBA{0x2f, 0x95, 0xe5, 0xff}),
	/* loresBlockLightBlue  */ newLoresBlock(color.RGBA{0xbf, 0xab, 0xff, 0xff}),
	/* loresBlockBrown      */ newLoresBlock(color.RGBA{0x40, 0x54, 0x00, 0xff}),
	/* loresBlockOrange     */ newLoresBlock(color.RGBA{0xd0, 0x6a, 0x1a, 0xff}),
	/* loresBlockGray2      */ newLoresBlock(color.RGBA{0x80, 0x80, 0x80, 0xff}),
	/* loresBlockPink       */ newLoresBlock(color.RGBA{0xff, 0x96, 0xbf, 0xff}),
	/* loresBlockLightGreen */ newLoresBlock(color.RGBA{0x2f, 0xbc, 0x1a, 0xff}),
	/* loresBlockYellow     */ newLoresBlock(color.RGBA{0xbf, 0xd3, 0x5a, 0xff}),
	/* loresBlockAquamarine */ newLoresBlock(color.RGBA{0x6f, 0xe8, 0xbf, 0xff}),
	/* loresBlockWhite      */ newLoresBlock(color.RGBA{0xff, 0xff, 0xff, 0xff}),
}

const (
	loresShadeLight = iota // Light colors (75% intensity)
	loreShadeMedium        // Medium colors (50% intensity)
	loreShadeDark          // Dark colors (25% intensity)
)

var loresColorShades = []int{
	loreShadeMedium, // Black (special case, always black)
	loreShadeMedium, // Magenta
	loreShadeDark,   // Dark Blue
	loreShadeMedium, // Purple
	loreShadeDark,   // Dark Green
	loreShadeMedium, // Gray1
	loreShadeMedium, // Medium Blue
	loresShadeLight, // Light Blue
	loreShadeMedium, // Brown
	loreShadeMedium, // Orange
	loreShadeMedium, // Gray2
	loreShadeMedium, // Pink
	loresShadeLight, // Light Green
	loreShadeMedium, // Yellow
	loreShadeMedium, // Aquamarine
	loreShadeMedium, // White (special case, always full monochrome color)
}

var (
	loresMonochromeGreenBlocks [16]*gfx.FrameBuffer
	loresMonochromeAmberBlocks [16]*gfx.FrameBuffer
)

func init() {
	for i := range 16 {
		loresMonochromeGreenBlocks[i] = newMonochromeLoresBlock(uint8(i), HiresMonochromeGreen)
		loresMonochromeAmberBlocks[i] = newMonochromeLoresBlock(uint8(i), HiresMonochromeAmber)
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

// newMonochromeLoresBlock returns a framebuffer that is filled with the
// monochrome color that applies to the normal lo-res color that would be
// displayed. That is, each color that isn't white or black would be shown as
// a darker hue of the monochrome color -- the question, "how dark?", is
// answered based on the whether that color is indexed as light, dark, or no
// qualifier.
func newMonochromeLoresBlock(bitPattern uint8, monochromeColor color.RGBA) *gfx.FrameBuffer {
	var clr color.RGBA

	index := bitPattern & 0xf

	switch index {
	case 0:
		clr = color.RGBA{0x00, 0x00, 0x00, 0xff}
	case 15:
		clr = monochromeColor
	default:
		shade := loresColorShades[index]
		switch shade {
		case loresShadeLight:
			clr = shadedMonochromeColor(monochromeColor, 0.75)
		case loreShadeMedium:
			clr = shadedMonochromeColor(monochromeColor, 0.50)
		case loreShadeDark:
			clr = shadedMonochromeColor(monochromeColor, 0.25)
		}
	}

	return newLoresBlock(clr)
}

// newLoresBlock returns a new framebuffer filled with the provided lo-res
// color.
func newLoresBlock(clr color.RGBA) *gfx.FrameBuffer {
	fbuf := gfx.NewFrameBuffer(loresBlockWidth, loresBlockHeight)

	for y := range uint(loresBlockHeight) {
		for x := range uint(loresBlockWidth) {
			err := fbuf.SetCell(x, y, clr)
			if err != nil {
				panic(err)
			}
		}
	}

	return fbuf
}

// LoresBlock will return a color rectangle that matches the color
// suggested by the given pattern of bits
func LoresBlock(bitPattern uint8, monochromeMode int) *gfx.FrameBuffer {
	index := bitPattern & 0xf

	switch monochromeMode {
	case MonochromeGreen:
		return loresMonochromeGreenBlocks[index]
	case MonochromeAmber:
		return loresMonochromeAmberBlocks[index]
	default:
		return loresColors[index]
	}
}

// RenderLores takes the data in the lo-res display buffer (essentially the
// text page) and writes that to the Screen framebuffer.
func RenderLores(seg memory.Getter, start, end int, monochromeMode int) {
	for addr := start; addr < end; addr++ {
		row := loresAddressRows[addr-start]
		col := loresAddressCols[addr-start]

		if row < 0 || col < 0 {
			continue
		}

		x := uint(col) * loresBlockWidth
		y := uint(row) * loresBlockHeight

		byt := seg.Get(int(addr))

		// The Apple IIe technical reference (p. 22) states that we
		// should show the low-order nibble in the top row, and
		// high-order nibble in the bottom row.
		_ = gfx.Screen.Blit(x, y, LoresBlock(byt&0xf, monochromeMode))
		_ = gfx.Screen.Blit(x, y+8, LoresBlock(byt>>4, monochromeMode))
	}
}
