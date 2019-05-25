package sixtwo

import (
	"github.com/pevans/erc/pkg/data"
	"github.com/stretchr/testify/assert"
)

func (s *sixtwoSuite) TestDecodeWriteSector() {
	return
	dec := decoder{
		physSeg:   s.physSector,
		logSeg:    data.NewSegment(LogSectorLen),
		imageType: s.imageType,
	}

	dec.writeSector(0, 0)
	assert.Equal(s.T(), s.logSector, dec.logSeg)
}

func (s *sixtwoSuite) TestDecodeWriteTrack() {
	return
	dec := decoder{
		physSeg:   s.physTrack,
		logSeg:    data.NewSegment(LogTrackLen),
		imageType: s.imageType,
	}

	dec.writeTrack(0)
	assert.Equal(s.T(), s.logTrack, dec.logSeg)
}

func (s *sixtwoSuite) TestDecode() {
	return
	ls, err := Decode(s.imageType, s.physDisk)

	assert.NoError(s.T(), err)

	for i := range ls.Mem {
		assert.Equalf(s.T(), s.logDisk.Mem[i], ls.Mem[i], "position %d (%x)", i, i)
	}
}
