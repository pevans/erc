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
		{"unshifted letter", ebiten.KeyA, 0, 'a'},
		{"shift input", ebiten.KeyDigit1, input.ModShift, '!'},
		{"shifted letter produces uppercase", ebiten.KeyA, input.ModShift, 'A'},
		{"unmapped key", ebiten.KeyF1, 0, rune(0)},
	}

	for _, c := range cases {
		t.Run(c.explanation, func(t *testing.T) {
			got := KeyToRune(c.key, c.modifier)
			assert.Equal(t, c.want, got)
		})
	}
}
