package a2

import (
	"testing"

	"github.com/pevans/erc/pkg/mach"
	"github.com/stretchr/testify/assert"
)

func TestNewDecoder(t *testing.T) {
	seg := mach.NewSegment(1)
	typ := 3

	dec := NewDecoder(typ, seg)
	assert.NotEqual(t, nil, dec)
	assert.Equal(t, typ, dec.imageType)
	assert.Equal(t, seg, dec.src)
}

func (s *encSuite) TestDecodeNIB() {
	err := loadFile(s.dec.src, s.baseDir+"/physical.disk")
	assert.Equal(s.T(), nil, err)

	dst, err := s.dec.DecodeNIB()
	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), true, fileMatches(dst, s.baseDir+"/physical.disk"))
}

func (s *encSuite) TestDecodeSector() {
	err := loadFile(s.dec.src, s.baseDir+"/physical.sector")
	assert.Equal(s.T(), nil, err)

	assert.Equal(s.T(), LogSectorLen, s.dec.DecodeSector(0, 0, 0, 0))
	assert.Equal(s.T(), true, fileMatches(s.dec.dst, s.baseDir+"/logical.sector"))
}
