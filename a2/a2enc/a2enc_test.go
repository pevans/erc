package a2enc_test

import (
	"testing"

	"github.com/pevans/erc/a2/a2enc"
	"github.com/pevans/erc/memory"
	"github.com/stretchr/testify/assert"
)

func TestSize(t *testing.T) {
	cases := []struct {
		name      string
		imageType int
		wantSize  int
		errfn     assert.ErrorAssertionFunc
	}{
		{
			name:      "DOS33 returns DosSize",
			imageType: a2enc.DOS33,
			wantSize:  a2enc.DosSize,
			errfn:     assert.NoError,
		},
		{
			name:      "ProDOS returns DosSize",
			imageType: a2enc.ProDOS,
			wantSize:  a2enc.DosSize,
			errfn:     assert.NoError,
		},
		{
			name:      "Nibble returns NibSize",
			imageType: a2enc.Nibble,
			wantSize:  a2enc.NibSize,
			errfn:     assert.NoError,
		},
		{
			name:      "unknown image type returns error",
			imageType: 99,
			wantSize:  -1,
			errfn:     assert.Error,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			size, err := a2enc.Size(c.imageType)

			assert.Equal(t, c.wantSize, size)
			c.errfn(t, err)
		})
	}
}

func TestDecode(t *testing.T) {
	logSeg := memory.NewSegment(a2enc.DosSize)
	for i := range logSeg.Size() {
		logSeg.Set(i, uint8(i%256))
	}

	encodedDOS33, err := a2enc.Encode(a2enc.DOS33, logSeg)
	assert.NoError(t, err)

	encodedProDOS, err := a2enc.Encode(a2enc.ProDOS, logSeg)
	assert.NoError(t, err)

	nibSeg := memory.NewSegment(a2enc.NibSize)

	cases := []struct {
		name      string
		imageType int
		seg       *memory.Segment
		equalfn   assert.ComparisonAssertionFunc
		nilfn     assert.ValueAssertionFunc
		errfn     assert.ErrorAssertionFunc
	}{
		{
			name:      "DOS33 decodes successfully",
			imageType: a2enc.DOS33,
			seg:       encodedDOS33,
			equalfn:   assert.NotEqual,
			nilfn:     assert.NotNil,
			errfn:     assert.NoError,
		},
		{
			name:      "ProDOS decodes successfully",
			imageType: a2enc.ProDOS,
			seg:       encodedProDOS,
			equalfn:   assert.NotEqual,
			nilfn:     assert.NotNil,
			errfn:     assert.NoError,
		},
		{
			name:      "Nibble returns same segment",
			imageType: a2enc.Nibble,
			seg:       nibSeg,
			equalfn:   assert.Equal,
			nilfn:     assert.NotNil,
			errfn:     assert.NoError,
		},
		{
			name:      "unknown image type returns error",
			imageType: 99,
			seg:       logSeg,
			equalfn:   assert.NotEqual,
			nilfn:     assert.Nil,
			errfn:     assert.Error,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result, err := a2enc.Decode(c.imageType, c.seg)

			c.errfn(t, err)
			c.nilfn(t, result)

			if result != nil {
				c.equalfn(t, c.seg, result)
			}
		})
	}
}
