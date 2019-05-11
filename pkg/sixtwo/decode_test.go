package sixtwo

import (
	"github.com/pevans/erc/pkg/data"
	"github.com/stretchr/testify/assert"
)

func (s *sixtwoSuite) TestDecode() {
	_, err := Decode(DOS33, s.physDisk)
	assert.Nil(s.T(), err)

	//assert.True(s.T(), fileMatches(dst, s.baseDir+"/logical.disk"))
}

func (s *sixtwoSuite) TestWriteSector() {
	var (
		sect = 0
		seg  = data.NewSegment(LogSectorLen)
		dec  = decoder{
			physSeg:   s.physDisk,
			logSeg:    data.NewSegment(DosSize),
			imageType: DOS33,
		}

		logStart = logicalSector(dec.imageType, sect)
		logEnd   = logStart + LogSectorLen
	)

	dec.writeSector(0, 0)
	dec.logSeg.WriteFile("/tmp/writeSector.seg")

	_, err := seg.CopySlice(0, dec.logSeg.Mem[logStart:logEnd])
	assert.NoError(s.T(), err)

	assert.True(s.T(), fileMatches(seg, s.logSectorPath))
}
