package font

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlyph(t *testing.T) {
	type test struct {
		ch    rune
		bm    *Bitmap
		errFn assert.ErrorAssertionFunc
	}

	bmp, _ := NewBitmap(A2System)
	bad, _ := NewBitmap(A2System)

	// This basically recycles the memory of the "bad" bitmap now. We
	// don't care if it errors out.
	_ = bad.img.Dispose()

	cases := map[string]test{
		"a character we have": test{
			ch:    '@',
			bm:    bmp,
			errFn: assert.NoError,
		},

		"a character we don't": test{
			ch:    '🍔',
			bm:    bmp,
			errFn: assert.Error,
		},

		// Really the only way to get ebiten's SubImage method to return
		// something resembling an error is to destroy the image
		// beforehand.
		"a disposed bitmap": test{
			ch:    'h',
			bm:    bad,
			errFn: assert.Error,
		},
	}

	for desc, c := range cases {
		t.Run(desc, func(tt *testing.T) {
			_, err := c.bm.Glyph(c.ch)
			c.errFn(tt, err)
		})
	}
}
