package a2hires_test

import (
	"testing"

	"github.com/pevans/erc/a2"
	"github.com/pevans/erc/a2/a2hires"
	"github.com/pevans/erc/a2/a2mono"
	"github.com/stretchr/testify/assert"
)

func TestPrepareRow(t *testing.T) {
	c := a2.NewComputer(123)
	assert.NoError(t, c.Boot())

	dots := make([]a2hires.Dot, 280)
	emptyDots := []a2hires.Dot{}

	t.Run("an insufficient length row of dots will error", func(t *testing.T) {
		assert.Error(t, a2hires.PrepareRow(c, 0, emptyDots, a2mono.None))
	})

	t.Run("a sufficient length row of dots will work", func(t *testing.T) {
		assert.NoError(t, a2hires.PrepareRow(c, 0, dots, a2mono.None))
	})
}
