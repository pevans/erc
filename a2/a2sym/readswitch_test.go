package a2sym

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadSwitch(t *testing.T) {
	t.Run("unknown switch", func(t *testing.T) {
		s := ReadSwitch(0)
		assert.Equal(t, ModeNone, s.Mode)
	})

	t.Run("known switches", func(t *testing.T) {
		s := ReadSwitch(0xC011)
		assert.Equal(t, ModeR7, s.Mode)
		assert.Equal(t, "RDBNK2", s.Name)
		assert.NotEmpty(t, s.Description)

		s = ReadSwitch(0xC000)
		assert.Equal(t, ModeR, s.Mode)
		assert.Empty(t, s.Name)
		assert.NotEmpty(t, s.Description)
	})
}
