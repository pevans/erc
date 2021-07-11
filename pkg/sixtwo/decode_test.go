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

func (s *sixtwoSuite) TestNewDecodeMap() {
	assert.NotNil(s.T(), newDecodeMap())
}

func (s *sixtwoSuite) TestLogByte() {
	// Test a decoder with a nominal decode map
	dec := decoder{
		decMap: newDecodeMap(),
	}
	s.Equal(uint8(0x3F), dec.logByte(uint8(0xFF)))

	// Test a decoder with no decode map
	dec = decoder{}
	s.Panics(func() {
		dec.logByte(uint8(0xFF))
	})
}

func (s *sixtwoSuite) TestDecodeWriteByte() {
	byt := uint8(123)
	dec := decoder{
		ls:        data.NewSegment(LogSectorLen),
		imageType: s.imageType,
	}

	dec.writeByte(byt)
	assert.Equal(s.T(), byt, dec.ls.Get(0))
}
