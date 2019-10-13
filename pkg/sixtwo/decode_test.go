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
