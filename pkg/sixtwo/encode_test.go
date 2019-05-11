package sixtwo

import (
	"testing"

	"github.com/pevans/erc/pkg/data"
	"github.com/stretchr/testify/assert"
)

func (s *sixtwoSuite) TestEncodeWrite() {
	type test struct {
		bytes     []data.Byte
		startPoff int
		wantPoff  int
	}

	cases := map[string]test{
		"no bytes": {
			bytes:     []data.Byte{},
			startPoff: 0,
			wantPoff:  0,
		},

		"some bytes": {
			bytes:     []data.Byte{0x23, 0x34, 0x45},
			startPoff: 5,
			wantPoff:  8,
		},
	}

	for desc, c := range cases {
		enc := newEncoder(0, c.startPoff+len(c.bytes))
		enc.poff = c.startPoff

		s.T().Run(desc, func(t *testing.T) {
			enc.write(c.bytes)
			assert.Equal(t, c.wantPoff, enc.poff)
		})
	}
}

func (s *sixtwoSuite) TestWriteByte() {
	bytes := []data.Byte{0, 1, 2}

	enc := newEncoder(0, len(bytes))
	for i, b := range bytes {
		enc.writeByte(b)
		assert.Equal(s.T(), i+1, enc.poff)
	}
}

func (s *sixtwoSuite) TestWrite4n4() {
	cases := []struct {
		byt  data.Byte
		want []data.Byte
	}{
		{0x32, []data.Byte{0xBB, 0xBA}},
		{0xFE, []data.Byte{0xFF, 0xFE}},
		{0x45, []data.Byte{0xAA, 0xEF}},
	}

	for _, c := range cases {
		enc := newEncoder(0, 2)
		enc.write4n4(c.byt)

		assert.Equal(s.T(), c.want, enc.ps.Mem)
	}
}
