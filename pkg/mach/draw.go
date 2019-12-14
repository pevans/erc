package mach

import (
	"github.com/pevans/erc/pkg/gfx"
)

// A Drawer is a type which can draw dots on a screen surface.
type Drawer interface {
	// Draw will actually render something on a given screen.
	Draw(screen gfx.DotDrawer) error

	// We allow Emulators to define their own dimensions. These aren't
	// necessarily what a viewer would see -- you might double or more
	// the actual render space from what this method might return.
	Dimensions() (width, height int)
}
