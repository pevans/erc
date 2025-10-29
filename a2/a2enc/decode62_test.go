package a2enc_test

import (
	"fmt"

	"github.com/pevans/erc/a2/a2enc"
)

func (s *sixtwoSuite) TestDecode62RoundTrip() {
	encoded, err := a2enc.Encode62(s.imageType, s.logDisk)
	s.NoError(err)
	s.NotNil(encoded)

	decoded, err := a2enc.Decode62(s.imageType, encoded)
	s.NoError(err)
	s.NotNil(decoded)

	s.Equal(s.logDisk.Size(), decoded.Size())

	for i := 0; i < s.logDisk.Size(); i++ {
		s.Equal(
			s.logDisk.Get(i), decoded.Get(i),
			fmt.Sprintf(
				"byte mismatch at offset %v", i,
			),
		)
	}
}
