package a2

import (
	"image/color"

	"github.com/pevans/erc/gfx"
)

const (
	loresBlockWidth  = 7
	loresBlockHeight = 4
)

var (
	loresColors = []*gfx.FrameBuffer{
		/* loresBlockBlack      */ loresNewBlock(color.RGBA{0x00, 0x00, 0x00, 0x00}),
		/* loresBlockMagenta    */ loresNewBlock(color.RGBA{0x90, 0x17, 0x40, 0x00}),
		/* loresBlockDarkBlue   */ loresNewBlock(color.RGBA{0x40, 0x2c, 0xa5, 0x00}),
		/* loresBlockPurple     */ loresNewBlock(color.RGBA{0xd0, 0x43, 0xe5, 0x00}),
		/* loresBlockDarkGreen  */ loresNewBlock(color.RGBA{0x00, 0x69, 0x40, 0x00}),
		/* loresBlockGray1      */ loresNewBlock(color.RGBA{0x80, 0x80, 0x80, 0x00}),
		/* loresBlockMediumBlue */ loresNewBlock(color.RGBA{0x2f, 0x95, 0xe5, 0x00}),
		/* loresBlockLightBlue  */ loresNewBlock(color.RGBA{0xbf, 0xab, 0xff, 0x00}),
		/* loresBlockBrown      */ loresNewBlock(color.RGBA{0x40, 0x54, 0x00, 0x00}),
		/* loresBlockOrange     */ loresNewBlock(color.RGBA{0xd0, 0x6a, 0x1a, 0x00}),
		/* loresBlockGray2      */ loresNewBlock(color.RGBA{0x80, 0x80, 0x80, 0x00}),
		/* loresBlockPink       */ loresNewBlock(color.RGBA{0xff, 0x96, 0xbf, 0x00}),
		/* loresBlockLightGreen */ loresNewBlock(color.RGBA{0x2f, 0xbc, 0x1a, 0x00}),
		/* loresBlockYellow     */ loresNewBlock(color.RGBA{0xbf, 0xd3, 0x5a, 0x00}),
		/* loresBlockAquamarine */ loresNewBlock(color.RGBA{0x6f, 0xe8, 0xbf, 0x00}),
		/* loresBlockWhite      */ loresNewBlock(color.RGBA{0xff, 0xff, 0xff, 0x00}),
	}
)

// Return a solid rectangle composed of a given color
func loresNewBlock(clr color.RGBA) *gfx.FrameBuffer {
	fbuf := gfx.NewFrameBuffer(loresBlockWidth, loresBlockHeight)

	for y := uint(0); y < loresBlockHeight; y++ {
		for x := uint(0); x < loresBlockWidth; x++ {
			fbuf.SetCell(x, y, clr)
		}
	}

	return fbuf
}

// Return a color rectangle that matches the color suggested by the
// given pattern of bits
func loresBlock(bitPattern uint8) *gfx.FrameBuffer {
	// Use a bitmask to prevent us from index something outside the
	// bounds of loresColors
	return loresColors[bitPattern&0xf]
}

func (c *Computer) loresRender(start, end int) {
	for addr := start; addr < end; addr++ {
		row := textAddressRows[addr-start]
		col := textAddressCols[addr-start]

		if row < 0 || col < 0 {
			continue
		}

		byt := c.Get(int(addr))

		x := uint(col) * loresBlockWidth
		y := uint(row) * loresBlockHeight

		block := loresBlock(byt >> 4)
		_ = gfx.Screen.Blit(x, y, block)

		y += 4
		block = loresBlock(byt & 0xf)
		_ = gfx.Screen.Blit(x, y, block)
	}
}
