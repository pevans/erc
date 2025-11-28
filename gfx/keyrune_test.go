package gfx

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pevans/erc/input"
	"github.com/stretchr/testify/assert"
)

func TestKeyToRune(t *testing.T) {
	cases := []struct {
		explanation string
		key         ebiten.Key
		modifier    int
		want        rune
	}{
		{"lowercase letter", ebiten.KeyA, 0, 'A'},
		{"shift input", ebiten.KeyDigit1, input.ModShift, '!'},
		{"uppercase letter (same as lowercase)", ebiten.KeyA, input.ModShift, 'A'},
		{"unmapped key", ebiten.KeyF1, 0, rune(0)},
	}

	for _, c := range cases {
		t.Run(c.explanation, func(t *testing.T) {
			got := KeyToRune(c.key, c.modifier)
			assert.Equal(t, c.want, got)
		})
	}
}
