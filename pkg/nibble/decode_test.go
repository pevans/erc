package nibble

import (
	"testing"

	"github.com/pevans/erc/pkg/data"
	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	// A funny thing about the encode procedure is it's identical to the
	// decode procedure, so the decode test should suffice.

	s, err := Decode(data.NewSegment(100))
	assert.NotNil(t, s)
	assert.NoError(t, err)
}

func TestNibbleCopier(t *testing.T) {
	s := data.NewSegment(100)
	d, err := nibbleCopier(s)

	assert.NoError(t, err)
	assert.NotNil(t, d)

	assert.Equal(t, d.Size(), s.Size())
	assert.Equal(t, d.Mem, s.Mem)
}
