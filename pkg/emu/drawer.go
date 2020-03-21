package emu

import "github.com/pevans/erc/pkg/data"

// A DrawContext is a very generic map type which allows the caller to
// pass some data into the Draw method that might be meaningful to it.
type DrawContext map[string]interface{}

// Drawer interfaces are ones which can produce some kind of graphical
// display (as a side-effect, not as a return value).
type Drawer interface {
	Draw(data.Getter, DrawContext) error
	Init() error
}
