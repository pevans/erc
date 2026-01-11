package a2enc_test

import (
	"github.com/pevans/erc/a2/a2enc"
)

func (s *sixtwoSuite) TestEncode62() {
	// Test that encoding the logical disk succeeds
	ps, err := a2enc.Encode62(s.imageType, s.logDisk)

	s.NoError(err)
	s.NotNil(ps)
	s.Equal(a2enc.EncodedSize, ps.Size())
}
