package a2hires

import "image/color"

// boundaryShiftColorMap are the likely shifts that may happen on dots that
// exist on byte boundaries.
var boundaryShiftColorMap = map[color.RGBA]color.RGBA{
	purple: darkPurple,
	green:  darkGreen,
	blue:   lightPurple,
	orange: lightGreen,
}

// shiftBoundaryDots updates the edges of colors in a very naive attempt to
// emulate NTSC color.
func shiftBoundaryDots(left, right Dot) (Dot, Dot) {
	// There's no need to shift colors if the palettes don't change
	if left.palette == right.palette {
		return left, right
	}

	switch right.color {
	case purple, green:
		// In this case, both left and right dots seem to shift color, and
		// only with regard to the right-hand dot.
		left.color = boundaryShiftColorMap[right.color]
		right.color = boundaryShiftColorMap[right.color]

	case blue, orange:
		// There's some curious logic here, and I have low confidence that
		// this is right in all cases. Only a single dot shifts, and only when
		// there isn't a black dot in the left-hand side.
		if left.color != black {
			right.color = boundaryShiftColorMap[right.color]
		}
	}

	return left, right
}
