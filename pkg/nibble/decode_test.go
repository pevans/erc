package nibble

import (
	"testing"

	"github.com/pevans/erc/pkg/data"
	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	s := data.NewSegment(100)
	d, err := Decode(s)

	assert.NoError(t, err)
	assert.NotNil(t, d)

	assert.Equal(t, d.Size(), s.Size())
	assert.Equal(t, d.Mem, s.Mem)
}
