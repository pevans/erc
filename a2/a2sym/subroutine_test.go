package a2sym

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubroutine(t *testing.T) {
	t.Run("unknown routine", func(t *testing.T) {
		assert.Empty(t, Subroutine(0))
	})

	t.Run("known routines", func(t *testing.T) {
		assert.Equal(t, "SETNORM", Subroutine(0xFE84))
		assert.Equal(t, "BELL", Subroutine(0xFF3A))
	})
}
