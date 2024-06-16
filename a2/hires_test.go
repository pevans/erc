package a2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: we should do some better tests to show the colors are what we
// want them to be.
func TestHiresDots(t *testing.T) {
	c := NewComputer(123)
	c.Boot()

	dots := make([]hiresDot, 280)
	emptyDots := []hiresDot{}

	t.Run("an insufficient length row of dots will error", func(t *testing.T) {
		assert.Error(t, c.HiresDots(0, emptyDots))
	})

	t.Run("a sufficient length row of dots will work", func(t *testing.T) {
		assert.NoError(t, c.HiresDots(0, dots))
	})
}
