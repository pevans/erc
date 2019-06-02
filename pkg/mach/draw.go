package mach

import (
	"github.com/pevans/erc/pkg/gfx"
)

type Drawer interface {
	Draw(screen gfx.DotDrawer) error
	Dimensions() (width, height int)
}
