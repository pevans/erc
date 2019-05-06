package sixtwo

import (
	"github.com/stretchr/testify/assert"
)

func (s *sixtwoSuite) TestDecode() {
	dst, err := Decode(DOS33, s.physDisk)
	assert.Nil(s.T(), err)

	assert.True(s.T(), fileMatches(dst, s.baseDir+"/logical.disk"))
}
