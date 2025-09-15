package a2video_test

import (
	"testing"

	"github.com/pevans/erc/a2"
	"github.com/pevans/erc/a2/a2video"
	"github.com/stretchr/testify/assert"
)

func TestPrepareHiresRow(t *testing.T) {
	c := a2.NewComputer(123)
	assert.NoError(t, c.Boot())

	dots := make([]a2video.HiresDot, 280)
	emptyDots := []a2video.HiresDot{}

	t.Run("an insufficient length row of dots will error", func(t *testing.T) {
		assert.Error(t, a2video.PrepareHiresRow(c, 0, emptyDots))
	})

	t.Run("a sufficient length row of dots will work", func(t *testing.T) {
		assert.NoError(t, a2video.PrepareHiresRow(c, 0, dots))
	})
}
