package a2

import (
	"image/color"
	"testing"

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
