package obj

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlice(t *testing.T) {
	slen := len(storeData)

	errCases := []struct {
		fail       bool
		start, end int
	}{
		{true, -1, -1},
		{true, -1, slen - 1},
		{true, 0, slen},
		{false, 0, slen - 1},
	}

	for _, c := range errCases {
		_, err := Slice(c.start, c.end)

		if c.fail {
			assert.NotEqual(t, nil, err)
		} else {
			assert.Equal(t, nil, err)
		}
	}
}
