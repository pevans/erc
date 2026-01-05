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

// Return a solid rectangle composed of a given color
func newLoresBlock(clr color.RGBA) *gfx.FrameBuffer {
	fbuf := gfx.NewFrameBuffer(loresBlockWidth, loresBlockHeight)

	for y := uint(0); y < loresBlockHeight; y++ {
		for x := uint(0); x < loresBlockWidth; x++ {
			err := fbuf.SetCell(x, y, clr)
			if err != nil {
				// This should really never happen...
				panic(err)
			}
		}
	}

	return fbuf
}

// LoresBlock will return a color rectangle that matches the color
// suggested by the given pattern of bits
func LoresBlock(bitPattern uint8) *gfx.FrameBuffer {
	// Use a bitmask to prevent us from index something outside the
	// bounds of loresColors
	return loresColors[bitPattern&0xf]
}

func RenderLores(seg memory.Getter, start, end int) {
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
		_ = gfx.Screen.Blit(x, y, LoresBlock(byt&0xf))
		_ = gfx.Screen.Blit(x, y+8, LoresBlock(byt>>4))
	}
}
