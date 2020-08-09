package nibble

import (
	"testing"

	"github.com/pevans/erc/pkg/data"
	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	// A funny thing about the encode procedure is it's identical to the
	// decode procedure, so the decode test should suffice.

	s, err := Encode(data.NewSegment(100))
	assert.NotNil(t, s)
	assert.NoError(t, err)
}
