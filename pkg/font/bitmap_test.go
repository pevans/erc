package font

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBitmap(t *testing.T) {
	type test struct {
		fn    Name
		errFn assert.ErrorAssertionFunc
	}

	cases := map[string]test{
		"a known font": test{
			fn:    A2System,
			errFn: assert.NoError,
		},

		"an unknown font": test{
			fn:    maxFontName,
			errFn: assert.Error,
		},
	}

	for desc, c := range cases {
		t.Run(desc, func(tt *testing.T) {
			_, err := NewBitmap(c.fn)
			c.errFn(tt, err)
		})
	}
}
