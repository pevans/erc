package a2enc_test

import (
	"github.com/pevans/erc/a2/a2enc"
	"github.com/stretchr/testify/assert"
)

func (s *sixtwoSuite) TestDecode() {
	ls, err := a2enc.Decode(s.imageType, s.physDisk)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), s.logDisk, ls)
}
