package a2

import (
	"image"
	"image/color"

	"github.com/pevans/erc/pkg/data"
	"github.com/pevans/erc/pkg/gfx"
)

var (
	HiresGreen      = color.RGBA{R: 0x2F, G: 0xBC, B: 0x1A}
	HiresPurple     = color.RGBA{R: 0xD0, G: 0x43, B: 0xE5}
	HiresOrange     = color.RGBA{R: 0xD0, G: 0x6A, B: 0x1A}
	HiresBlue       = color.RGBA{R: 0x2F, G: 0x95, B: 0xE5}
	HiresBlack      = color.RGBA{R: 0x00, G: 0x00, B: 0x00}
	HiresWhite      = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF}
	HiresColorTable = []color.RGBA{
		HiresGreen,
		HiresPurple,
		HiresOrange,
		HiresBlue,
	}

	// This table maps a row number to a base address in the hires graphics
	// buffer. From there, (base + i) maps to column i in that row.
	HiresAddresses = []data.DByte{
		//   0       1       2       3       4       5       6       7
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
)

func (c *Computer) HiresRowDots(row int) []data.Byte {
	addr := HiresAddresses[row%192]
	dots := make([]data.Byte, 280)

	i := 0
	for bytePos := data.DByte(0); bytePos < 40; bytePos++ {
		byt := c.Get(addr + bytePos)

		for pos := uint(0); pos < 7; pos++ {
			dots[i] = 0
			if byt&0x80 > 0 {
				dots[i] = 2
			}

			dots[i] |= (byt >> pos) & 1

			i++
		}
	}

	return dots
}

func hiresColor(curr, next data.Byte, pos int) color.RGBA {
	if curr > 0 && next > 0 {
		return HiresWhite
	}

	if curr == 0 && next == 0 {
		return HiresBlack
	}

	fn := func(dot data.Byte, pos int) color.RGBA {
		var cindex int

		if pos%2 == 0 {
			cindex++
		}

		if dot&2 > 0 {
			cindex += 2
		}

		return HiresColorTable[cindex]
	}

	if next > 0 {
		return fn(next, pos+1)
	}

	return fn(curr, pos)
}

func (c *Computer) DrawHiresRow(screen gfx.DotDrawer, row int) {
	var (
		curr, next data.Byte
		useColor   color.RGBA
	)

	dots := c.HiresRowDots(row)

	for i := 0; i < 279; i++ {
		curr = dots[i] & 1
		next = dots[i+1] & 1
		useColor = hiresColor(curr, next, i)

		screen.DrawDot(image.Point{X: i, Y: row}, useColor)
	}
}
