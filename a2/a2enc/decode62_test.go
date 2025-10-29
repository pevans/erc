package a2enc_test

import (
	"github.com/pevans/erc/a2/a2enc"
)

func (s *sixtwoSuite) TestDecode62() {
	return
	ls, err := a2enc.Decode62(s.imageType, s.physDisk)

	s.NoError(err)
	s.NotNil(ls)
	s.Equal(a2enc.DosSize, ls.Size())
}

func (s *sixtwoSuite) TestDecode62RoundTrip() {
	return
	encoded, err := a2enc.Encode62(s.imageType, s.logDisk)
	s.NoError(err)

	decoded, err := a2enc.Decode62(s.imageType, encoded)
	s.NoError(err)
	s.Equal(s.logDisk, decoded)
}
