package a2sym

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteSwitch(t *testing.T) {
	t.Run("unknown switch", func(t *testing.T) {
		s := WriteSwitch(0)
		assert.Equal(t, ModeNone, s.Mode)
	})

	t.Run("known switches", func(t *testing.T) {
		s := WriteSwitch(0xC010)
		assert.Equal(t, ModeW, s.Mode)
		assert.Equal(t, "clear strobe", s.Description)

		s = WriteSwitch(0xC050)
		assert.Equal(t, ModeRW, s.Mode)
		assert.Equal(t, "TEXT", s.Name)
		assert.NotEmpty(t, s.Description)
	})
}
