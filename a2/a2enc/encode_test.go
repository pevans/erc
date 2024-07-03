package a2enc_test

import (
	"github.com/pevans/erc/a2/a2enc"
	"github.com/stretchr/testify/assert"
)

func (s *sixtwoSuite) TestEncode() {
	ps, err := a2enc.Encode(s.imageType, s.logDisk)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), s.physDisk, ps)
}
