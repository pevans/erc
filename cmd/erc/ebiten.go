package main

import (
	"github.com/hajimehoshi/ebiten"
)

func gameLoop() error {
	var (
		width, height = emulator.Drawer.Dimensions()
	)

	// We want things to execute even when the window doesn't have
	// active focus.
	ebiten.SetRunnableInBackground(true)

	return ebiten.Run(ebitenLoop, width, height, 3, "erc")
}

// ebitenLoop is the "run loop" of our graphics logic, which we use both
// to implement processor speed and frame rate.
func ebitenLoop(screen *ebiten.Image) error {
	// ebiten gives us a way to see if this result would even be
	// considered, and if not, we can return immediately.
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	// The Draw method will give us what the current screen should look
	// like.
	updatedScreen := emulator.Drawer.Draw()

	// Which is then flashed to ebiten.
	return screen.DrawImage(updatedScreen, nil)
}
