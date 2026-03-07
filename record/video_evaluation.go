package record

import (
	"image/color"
)

// EvaluateVideoAssertions evaluates all video assertions against the captured
// frames in the recorder.
func EvaluateVideoAssertions(
	assertions []VideoAssertion,
	recorder *VideoRecorder,
) []VideoAssertionResult {
	results := make([]VideoAssertionResult, len(assertions))

	for i, a := range assertions {
		results[i] = evaluateVideoAssertion(a, recorder)
	}

	return results
}

func evaluateVideoAssertion(a VideoAssertion, recorder *VideoRecorder) VideoAssertionResult {
	frame := recorder.Frame(a.Step)
	if frame == nil {
		return VideoAssertionResult{
			Assertion: a,
			Passed:    false,
			Failure:   &VideoMismatch{MissingFrame: true},
		}
	}

	switch a.Kind {
	case VideoScreen:
		return compareGrid(a, frame, 0, 0, a.GridW, a.GridH)
	case VideoRegion:
		return compareGrid(a, frame, a.RegionX, a.RegionY, a.RegionW, a.RegionH)
	case VideoRow:
		return compareGrid(a, frame, 0, a.RowIndex, a.GridW, 1)
	default:
		return VideoAssertionResult{
			Assertion: a,
			Passed:    false,
			Failure:   &VideoMismatch{},
		}
	}
}

func samplePixelChar(
	frame *frameSnapshot,
	legend ColorLegend,
	fx, fy uint,
) (byte, color.RGBA, bool) {
	pixel := frame.GetPixel(fx, fy)
	rgb := color.RGBA{R: pixel.R, G: pixel.G, B: pixel.B, A: 0xff}

	ch, ok := legend.Reverse[rgb]
	return ch, rgb, ok
}

func compareGrid(
	a VideoAssertion,
	frame *frameSnapshot,
	startX, startY, w, h int,
) VideoAssertionResult {
	result := VideoAssertionResult{Assertion: a, Passed: true}

	fw := int(frame.width)
	fh := int(frame.height)

	for gy := range h {
		expectedRow := a.Expected[gy]
		fy := uint((startY + gy) * fh / a.GridH)

		for gx := range w {
			fx := uint((startX + gx) * fw / a.GridW)

			ch, rgb, ok := samplePixelChar(frame, a.Legend, fx, fy)
			if !ok {
				result.Passed = false
				result.Failure = &VideoMismatch{
					X:             gx,
					Y:             gy,
					Expected:      expectedRow[gx],
					ActualRGB:     rgb,
					ExpectedRow:   expectedRow,
					UnmappedColor: true,
				}
				return result
			}

			if ch != expectedRow[gx] {
				actualRow := buildActualRow(frame, a, fw, fh, startX, startY, gy, w)
				copy(actualRow[:gx+1], []byte(string([]byte{ch})))
				// We already know chars 0..gx-1 matched
				for px := range gx {
					pfx := uint((startX + px) * fw / a.GridW)
					pch, _, _ := samplePixelChar(frame, a.Legend, pfx, fy)
					actualRow[px] = pch
				}
				actualRow[gx] = ch

				result.Passed = false
				result.Failure = &VideoMismatch{
					X:           gx,
					Y:           gy,
					Expected:    expectedRow[gx],
					Actual:      ch,
					ActualRGB:   rgb,
					ExpectedRow: expectedRow,
					ActualRow:   string(actualRow),
				}
				return result
			}
		}
	}

	return result
}

func buildActualRow(
	frame *frameSnapshot,
	a VideoAssertion,
	fw, fh, startX, startY, gy, w int,
) []byte {
	row := make([]byte, w)
	fy := uint((startY + gy) * fh / a.GridH)

	for gx := range w {
		fx := uint((startX + gx) * fw / a.GridW)
		ch, _, ok := samplePixelChar(frame, a.Legend, fx, fy)
		if ok {
			row[gx] = ch
		} else {
			row[gx] = '?'
		}
	}

	return row
}
