package sixtwo

import (
	"testing"

	"github.com/pevans/erc/memory"
	"github.com/stretchr/testify/assert"
)

func (s *sixtwoSuite) TestEncodeWrite() {
	type test struct {
		bytes     []uint8
		startPoff int
		wantPoff  int
	}

	cases := map[string]test{
		"no bytes": {
			bytes:     []uint8{},
			startPoff: 0,
			wantPoff:  0,
		},

		"some bytes": {
			bytes:     []uint8{0x23, 0x34, 0x45},
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
	bytes := []uint8{0, 1, 2}

	enc := newEncoder(0, len(bytes))
	for i, b := range bytes {
		enc.writeByte(b)
		assert.Equal(s.T(), i+1, enc.poff)
	}
}

func (s *sixtwoSuite) TestWrite4n4() {
	cases := []struct {
		byt  uint8
		want []uint8
	}{
		{0x32, []uint8{0xBB, 0xBA}},
		{0xFE, []uint8{0xFF, 0xFE}},
		{0x45, []uint8{0xAA, 0xEF}},
	}

	for _, c := range cases {
		enc := newEncoder(0, 2)
		enc.write4n4(c.byt)

		assert.Equal(s.T(), c.want, enc.ps.Mem)
	}
}

func (s *sixtwoSuite) TestEncodeWriteSector() {
	enc := encoder{
		ls:        s.logSector,
		ps:        memory.NewSegment(PhysSectorLen),
		imageType: s.imageType,
	}

	enc.writeSector(0, 0)
	assert.Equal(s.T(), s.physSector, enc.ps)
}

func (s *sixtwoSuite) TestEncodeWriteTrack() {
	enc := encoder{
		ls:        s.logTrack,
		ps:        memory.NewSegment(PhysTrackLen),
		imageType: s.imageType,
	}

	enc.writeTrack(0)
	assert.Equal(s.T(), s.physTrack, enc.ps)
}

func (s *sixtwoSuite) TestEncode() {
	ps, err := Encode(s.imageType, s.logDisk)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), s.physDisk, ps)
}
