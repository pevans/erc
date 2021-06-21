package a2

import (
	"fmt"
	"image/color"
	"testing"

	"github.com/pevans/erc/pkg/data"
	"github.com/stretchr/testify/assert"
)

func TestHiresDots(t *testing.T) {
	t.Run("all black row", func(_ *testing.T) {
		assert.Equal(t, []color.RGBA{
			hiresBlack,
			hiresBlack,
			hiresBlack,
			hiresBlack,
			hiresBlack,
			hiresBlack,
			hiresBlack,
		}, HiresDots(0x00))
	})

	t.Run("all white row", func(_ *testing.T) {
		assert.Equal(t, []color.RGBA{
			hiresWhite,
			hiresWhite,
			hiresWhite,
			hiresWhite,
			hiresWhite,
			hiresWhite,
			hiresWhite,
		}, HiresDots(0xFF))
	})

	// There are some other permutations I can test here, but want to
	// get the hires visuals done so I can eyeball them
}

func TestHiresPoint(t *testing.T) {
	fn := func(a data.DByte, x, y int) {
		t.Run(
			fmt.Sprintf("hires_point_addr=%04X_x=%d_y=%d", a, x, y),
			func(t *testing.T) {
				rx, ry := HiresPoint(a)
				assert.Equal(t, x, rx, "x")
				assert.Equal(t, y, ry, "y")
			},
		)
	}

	fn(0x2000, 0, 0)
	fn(0x2001, 1, 0)
	fn(0x2028, 0, 8)
	fn(0x2308, 8, 6)
}
