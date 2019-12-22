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
	updatedScreen := emulator.Drawer.Draw()

	return screen.DrawImage(updatedScreen, nil)
}
