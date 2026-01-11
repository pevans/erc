package a2

import (
	"fmt"
	"image/color"
	"os"
	"strings"

	"github.com/pevans/erc/gfx"
)

// A ScreenFrame is a container of a single visual frame (as would have been
// rendered on screen) as text.
type ScreenFrame struct {
	Timestamp float64
	Rows      [192]string
}

// A ScreenLog is a collection of ScreenFrames, and represents everything
// we've visually captured of the software.
type ScreenLog struct {
	frames []ScreenFrame
}

// NewScreenLog returns a newly allocated ScreenLog that is ready to add new
// frames.
func NewScreenLog() *ScreenLog {
	return &ScreenLog{
		frames: make([]ScreenFrame, 0),
	}
}

// CaptureFrame records a ScreenFrame from the given FrameBuffer. That frame
// is appended to the ScreenLog receiver.
func (s *ScreenLog) CaptureFrame(fb *gfx.FrameBuffer, timestamp float64) {
	frame := ScreenFrame{Timestamp: timestamp}

	for y := range 192 {
		var row strings.Builder
		for x := range 280 {
			// Sample the doubled display at (x*2, y*2)
			c := fb.GetPixel(uint(x*2), uint(y*2))
			row.WriteRune(colorToChar(c))
		}
		frame.Rows[y] = row.String()
	}

	s.frames = append(s.frames, frame)
}

// colorToChar converts a given screen color to a character reprsentation in a
// ScreenFrame
func colorToChar(c color.RGBA) rune {
	switch {
	case c.R == 0xff && c.G == 0xff && c.B == 0xff:
		return 'W'
	case c.R == 0x2f && c.G == 0x95 && c.B == 0xe5:
		return 'B'
	case c.R == 0xd0 && c.G == 0x6a && c.B == 0x1a:
		return 'O'
	case (c.R == 0x2f && c.G == 0xbc && c.B == 0x1a) ||
		(c.R == 0x3f && c.G == 0x4c && c.B == 0x12) ||
		(c.R == 0xbd && c.G == 0xea && c.B == 0x86):
		return 'G'
	case (c.R == 0xd0 && c.G == 0x43 && c.B == 0xe5) ||
		(c.R == 0x3e && c.G == 0x31 && c.B == 0x79) ||
		(c.R == 0xbb && c.G == 0xaf && c.B == 0xf6):
		return 'P'
	default:
		return ' '
	}
}

// WriteToFile will write the contents of a ScreenLog to a file with the
// provided filename. If that can't be done, an error is returned.
func (s *ScreenLog) WriteToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close() //nolint:errcheck

	for _, frame := range s.frames {
		_, err := fmt.Fprintf(file, "FRAME %.6f\n", frame.Timestamp)
		if err != nil {
			return err
		}

		for _, row := range frame.Rows {
			_, err := fmt.Fprintln(file, row)
			if err != nil {
				return err
			}
		}

		_, err = fmt.Fprintln(file)
		if err != nil {
			return err
		}
	}

	return nil
}
