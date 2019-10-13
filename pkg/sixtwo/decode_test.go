package sixtwo

import (
	"github.com/pevans/erc/pkg/data"
	"github.com/stretchr/testify/assert"
)

func (s *sixtwoSuite) TestDecodeWriteSector() {
	dec := decoder{
		ps:        s.physSector,
		ls:        data.NewSegment(LogSectorLen),
		imageType: s.imageType,
		decMap:    newDecodeMap(),
	}

	dec.writeSector(0, 0)
	assert.Equal(s.T(), s.logSector, dec.ls)
}

func (s *sixtwoSuite) TestDecodeWriteTrack() {
	dec := decoder{
		ps:        s.physTrack,
		ls:        data.NewSegment(LogTrackLen),
		imageType: s.imageType,
		decMap:    newDecodeMap(),
	}

	dec.writeTrack(0)
	assert.Equal(s.T(), s.logTrack, dec.ls)
}

func (s *sixtwoSuite) TestDecode() {
	ls, err := Decode(s.imageType, s.physDisk)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), s.logDisk, ls)
}
