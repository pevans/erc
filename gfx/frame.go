package gfx

import "image/color"

// A FrameManager is something which can manage a graphics frame, which is a
// two-dimensional set of rows and columns (composed of "cells") that contain
// color information.
type FrameManager interface {
	SetCell(row, col int, clr color.RGBA) error
	ClearCells(clr color.RGBA)
}

// A FrameFilter is some device which, given a specific cell, can return a
// filtered color. For example, think of a filter which could render a frame
// in monochrome.
type FrameFilter interface {
	Filter(row, col int, clr color.RGBA) color.RGBA
}
