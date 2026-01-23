package a2video

import (
	"image/color"

	"github.com/pevans/erc/gfx"
)

// DHiresGetter is an interface for reading from both main and auxiliary
// memory. GetMain returns main memory at 0x2000-0x3FFF regardless of page
// settings.
type DHiresGetter interface {
	GetMain(addr int) uint8
	GetAux(addr int) uint8
}

// dhiresColors maps 4-bit patterns to colors we'll use in double hi-res.
var dhiresColors = []color.RGBA{
	/* 0000 Black       */ {0x00, 0x00, 0x00, 0xff},
	/* 0001 Magenta     */ {0x90, 0x17, 0x40, 0xff},
	/* 0010 Brown       */ {0x40, 0x54, 0x00, 0xff},
	/* 0011 Orange      */ {0xd0, 0x6a, 0x1a, 0xff},
	/* 0100 Dark Green  */ {0x00, 0x69, 0x40, 0xff},
	/* 0101 Gray 1      */ {0x80, 0x80, 0x80, 0xff},
	/* 0110 Green       */ {0x2f, 0xbc, 0x1a, 0xff},
	/* 0111 Yellow      */ {0xbf, 0xd3, 0x5a, 0xff},
	/* 1000 Dark Blue   */ {0x40, 0x2c, 0xa5, 0xff},
	/* 1001 Purple      */ {0xd0, 0x43, 0xe5, 0xff},
	/* 1010 Gray 2      */ {0x80, 0x80, 0x80, 0xff},
	/* 1011 Pink        */ {0xff, 0x96, 0xbf, 0xff},
	/* 1100 Medium Blue */ {0x2f, 0x95, 0xe5, 0xff},
	/* 1101 Light Blue  */ {0xbf, 0xab, 0xff, 0xff},
	/* 1110 Aqua        */ {0x6f, 0xe8, 0xbf, 0xff},
	/* 1111 White       */ {0xff, 0xff, 0xff, 0xff},
}

// RenderDHires draws double hi-res graphics from both main and auxiliary
// memory. Double hi-res provides 560x192 monochrome or 140x192 with 16
// colors.
func RenderDHires(seg DHiresGetter, monochromeMode int) {
	for y := range uint(192) {
		renderDHiresRow(seg, y, monochromeMode)
	}
}

// renderDHiresRow renders a single row of double hi-res graphics.
func renderDHiresRow(seg DHiresGetter, row uint, monochromeMode int) {
	addr := int(hiresAddrs[row])

	if monochromeMode != MonochromeNone {
		renderDHiresRowMono(seg, row, addr, monochromeMode)
		return
	}

	renderDHiresRowColor(seg, row, addr)
}

// renderDHiresRowMono renders a row in monochrome mode (560 dots). This ends
// up working a lot like monochrome for hi-res. I never used double hi-res
// software on a monochrome monitor, so this is more of a guess as to how it
// should look.
func renderDHiresRowMono(seg DHiresGetter, row uint, addr int, monochromeMode int) {
	monochromeColor := HiresMonochromeGreen
	if monochromeMode == MonochromeAmber {
		monochromeColor = HiresMonochromeAmber
	}

	ypos := row * 2

	// Double hi-res has 80 screen byte columns: aux[0], main[0], aux[1],
	// main[1], ... Each memory offset (0-39) contributes 2 screen columns
	// (aux and main).
	for memOffset := range 40 {
		byteAddr := addr + memOffset
		auxByte := seg.GetAux(byteAddr)
		mainByte := seg.GetMain(byteAddr)

		// Aux byte goes to even screen column, main byte to odd screen column
		auxScreenCol := memOffset * 2
		mainScreenCol := memOffset*2 + 1

		// Render aux byte (7 dots)
		baseX := uint(auxScreenCol * 7)

		for bit := range 7 {
			xpos := baseX + uint(bit)

			var clr color.RGBA
			if auxByte&(1<<bit) != 0 {
				clr = monochromeColor
			} else {
				clr = HiresBlack
			}

			_ = gfx.Screen.SetCell(xpos, ypos, clr)
			_ = gfx.Screen.SetCell(xpos, ypos+1, clr)
		}

		// Render main byte (7 dots)
		baseX = uint(mainScreenCol * 7)

		for bit := range 7 {
			xpos := baseX + uint(bit)

			var clr color.RGBA
			if mainByte&(1<<bit) != 0 {
				clr = monochromeColor
			} else {
				clr = HiresBlack
			}

			_ = gfx.Screen.SetCell(xpos, ypos, clr)
			_ = gfx.Screen.SetCell(xpos, ypos+1, clr)
		}
	}
}

// renderDHiresRowColor renders a row in color mode.
func renderDHiresRowColor(seg DHiresGetter, row uint, addr int) {
	ypos := row * 2

	// Build array of all 560 dots in this row
	var dots [560]bool

	for memOffset := range 40 {
		byteAddr := addr + memOffset
		auxByte := seg.GetAux(byteAddr)
		mainByte := seg.GetMain(byteAddr)

		// Aux byte goes to even screen column, main to odd screen columns:
		// aux[0], main[0], aux[1], main[1], ...
		auxScreenCol := memOffset * 2
		mainScreenCol := memOffset*2 + 1

		// Aux byte: 7 dots starting at auxScreenCol * 7
		auxBaseX := auxScreenCol * 7
		for bit := range 7 {
			dots[auxBaseX+bit] = (auxByte & (1 << bit)) != 0
		}

		// Main byte: 7 dots starting at mainScreenCol * 7
		mainBaseX := mainScreenCol * 7
		for bit := range 7 {
			dots[mainBaseX+bit] = (mainByte & (1 << bit)) != 0
		}
	}

	renderDHiresRowWithSlidingWindow(dots, ypos)
}

// renderDHiresRowWithSlidingWindow renders using a per-dot sliding window.
// Each of the 560 dots gets its own color based on a 4-dot window starting at
// that position. The result should be something that resembles NTSC color
// compositing.
func renderDHiresRowWithSlidingWindow(dots [560]bool, ypos uint) {
	for dot := range 560 {
		pattern := 0

		// Look at 4 dots starting at this position. Map each dot to its bit
		// position based on its phase in the color cycle
		for i := range 4 {
			idx := dot + i
			if idx < 560 && dots[idx] {
				// Phase 0 is bit 3, phase 1 is bit 2, etc.
				phase := idx % 4
				pattern |= (1 << (3 - phase))
			}
		}

		clr := dhiresColors[pattern]

		_ = gfx.Screen.SetCell(uint(dot), ypos, clr)
		_ = gfx.Screen.SetCell(uint(dot), ypos+1, clr)
	}
}
