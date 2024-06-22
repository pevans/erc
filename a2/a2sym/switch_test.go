package a2sym

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSwitchModeString(t *testing.T) {
	cases := []struct {
		mode   SwitchMode
		expect string
	}{
		{mode: ModeNone, expect: "unknown switch mode"},
		{mode: ModeR, expect: "read"},
		{mode: ModeR7, expect: "read, result in high bit"},
		{mode: ModeRR, expect: "read, twice consecutively"},
		{mode: ModeRW, expect: "read or write"},
		{mode: ModeW, expect: "write"},
	}

	for _, c := range cases {
		assert.Equal(t, c.expect, c.mode.String())
	}
}

func TestSwitchString(t *testing.T) {
	t.Run("zero value", func(t *testing.T) {
		var s Switch
		assert.Empty(t, s.String())
	})

	t.Run("no name", func(t *testing.T) {
		s := Switch{
			Mode:        ModeR,
			Description: "abc",
		}

		assert.Contains(t, s.String(), s.Mode.String())
		assert.Contains(t, s.String(), s.Description)
	})

	t.Run("no description", func(t *testing.T) {
		s := Switch{
			Mode: ModeR,
			Name: "THING",
		}

		assert.Contains(t, s.String(), s.Mode.String())
		assert.Contains(t, s.String(), s.Name)
	})

	t.Run("everything", func(t *testing.T) {
		s := Switch{
			Mode:        ModeR,
			Name:        "THING",
			Description: "abc",
		}

		assert.Contains(t, s.String(), s.Mode.String())
		assert.Contains(t, s.String(), s.Name)
		assert.Contains(t, s.String(), s.Description)
	})
}
