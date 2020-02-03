package a2

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/pevans/erc/pkg/data"
)

const (
	Text40ColStart data.DByte = 0x400
	Text40ColEnd   data.DByte = 0x800
)

// DrawText40 will, given a target image, render what needs to be shown
// for text on the screen. This method makes the assumption that the
// text should be rendered in a 40-column display.
func (c *Computer) DrawText40(img *ebiten.Image) {
	for addr := Text40ColStart; addr < Text40ColEnd; addr++ {
	}
}
