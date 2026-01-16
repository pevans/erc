package a2video

import "image/color"

// Define the likely shifts that may happen on dots that exist on byte
// boundaries.
var boundaryShiftColorMap = map[color.RGBA]color.RGBA{
	HiresPurple: HiresDarkPurple,
	HiresGreen:  HiresDarkGreen,
	HiresBlue:   HiresLightPurple,
	HiresOrange: HiresLightGreen,
}

// It was possible for high resolution displays to have darker or lighter
// colored dots, but this seems to have been an artifact of how NTSC worked,
// or at least CRT screens of the era. This function tries to replicate these
// color shifts.
func shiftBoundaryDots(left, right HiresDot) (HiresDot, HiresDot) {
	// There's no need to shift colors if the palettes don't change
	if left.Palette == right.Palette {
		return left, right
	}

	switch right.Color {
	case HiresPurple, HiresGreen:
		// In this case, both left and right dots seem to shift color, and
		// only with regard to the right-hand dot.
		left.Color = boundaryShiftColorMap[right.Color]
		right.Color = boundaryShiftColorMap[right.Color]

	case HiresBlue, HiresOrange:
		// There's some curious logic here, and I have low confidence that
		// this is right in all cases. Only a single dot shifts, and only when
		// there isn't a black dot in the left-hand side.
		if left.Color != HiresBlack {
			right.Color = boundaryShiftColorMap[right.Color]
		}
	}

	return left, right
}
