package mach

type Drawer interface {
	Draw() error
	Dimensions() (width, height int)
}
