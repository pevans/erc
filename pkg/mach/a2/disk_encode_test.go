package a2

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/pevans/erc/pkg/mach"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type encSuite struct {
	suite.Suite

	enc     *Encoder
	dos     *mach.Segment
	baseDir string
}

func (s *encSuite) SetupSuite() {
	s.enc = NewEncoder(0, nil)
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	s.baseDir = dir + "/../../../data"
}

func (s *encSuite) SetupTest() {
	s.enc.imageType = DDDOS33
	s.enc.src = mach.NewSegment(DD140K)
	s.enc.dst = mach.NewSegment(DD140KNib)
}

func TestEncoderSuite(t *testing.T) {
	suite.Run(t, new(encSuite))
}

func TestNewEncoder(t *testing.T) {
	seg := mach.NewSegment(1)
	typ := 3

	enc := NewEncoder(typ, seg)
	assert.NotEqual(t, nil, enc)
	assert.Equal(t, seg, enc.src)
	assert.Equal(t, typ, enc.imageType)
}

func (s *encSuite) TestLogicalSector() {
	cases := []struct {
		imgType int
		psect   int
		want    int
	}{
		{0, 0, 0},
		{DDDOS33, -1, 0},
		{DDDOS33, 16, 0},
		{DDDOS33, 0x0, 0x0},
		{DDDOS33, 0x1, 0x7},
		{DDDOS33, 0xE, 0x8},
		{DDDOS33, 0xF, 0xF},
		{DDProDOS, 0x0, 0x0},
		{DDProDOS, 0x1, 0x8},
		{DDProDOS, 0xE, 0x7},
		{DDProDOS, 0xF, 0xF},
		{DDNibble, 1, 1},
	}

	for _, c := range cases {
		s.enc.imageType = c.imgType
		assert.Equal(s.T(), c.want, s.enc.LogicalSector(c.psect))
	}
}

func (s *encSuite) TestEncodeNIB() {
	_, _ = s.enc.src.CopySlice(0, []mach.Byte{0x1, 0x2, 0x3})

	dst, err := s.enc.EncodeNIB()
	assert.Equal(s.T(), nil, err)

	for i := 0; i < dst.Size(); i++ {
		assert.Equal(s.T(), s.enc.src.Mem[i], dst.Mem[i])
	}
}

func (s *encSuite) TestEncodeDOS() {
	bytes, err := ioutil.ReadFile(s.baseDir + "/logical.disk")
	assert.Equal(s.T(), nil, err)

	s.enc.src = mach.NewSegment(DD140K)
	for i, b := range bytes {
		s.enc.src.Set(mach.DByte(i), mach.Byte(b))
	}

	dst, err := s.enc.EncodeDOS()
	assert.Equal(s.T(), nil, err)

	bytes, err = ioutil.ReadFile(s.baseDir + "/physical.disk")
	assert.Equal(s.T(), nil, err)

	for i, b := range bytes {
		assert.Equal(s.T(), mach.Byte(b), dst.Mem[i])
	}
}

func (s *encSuite) TestWrite() {
	bytes := []mach.Byte{0x1, 0x2, 0x3}
	_, _ = s.enc.src.CopySlice(0, bytes)

	assert.Equal(s.T(), 3, s.enc.Write(0, bytes))

	for i := 0; i < len(bytes); i++ {
		assert.Equal(s.T(), s.enc.src.Mem[i], s.enc.dst.Mem[i])
	}
}

func (s *encSuite) TestEncode4n4() {
	cases := []struct {
		in    mach.Byte
		want1 mach.Byte
		want2 mach.Byte
	}{
		{0xFE, 0xFF, 0xFE},
		{0x37, 0xBB, 0xBF},
	}

	for _, c := range cases {
		assert.Equal(s.T(), 2, s.enc.Encode4n4(0, c.in))
		assert.Equal(s.T(), c.want1, s.enc.dst.Mem[0])
		assert.Equal(s.T(), c.want2, s.enc.dst.Mem[1])
	}
}

func (s *encSuite) TestEncodeSector() {
	bytes, err := ioutil.ReadFile(s.baseDir + "/logical.sector")
	assert.Equal(s.T(), nil, err)

	mbytes := make([]mach.Byte, len(bytes))
	for i, b := range bytes {
		mbytes[i] = mach.Byte(b)
	}

	s.enc.src.CopySlice(0, mbytes)
	assert.Equal(s.T(), PhysSectorLen, s.enc.EncodeSector(0, 0, 0, 0))

	bytes, err = ioutil.ReadFile(s.baseDir + "/physical.sector")
	assert.Equal(s.T(), nil, err)

	for i, b := range bytes {
		assert.Equal(s.T(), mach.Byte(b), s.enc.dst.Mem[i])
	}
}

func (s *encSuite) TestEncodeTrack() {
	bytes, err := ioutil.ReadFile(s.baseDir + "/logical.track")
	assert.Equal(s.T(), nil, err)

	mbytes := make([]mach.Byte, len(bytes))
	for i, b := range bytes {
		mbytes[i] = mach.Byte(b)
	}

	s.enc.src.CopySlice(0, mbytes)
	assert.Equal(s.T(), PhysTrackLen, s.enc.EncodeTrack(0, 0))

	bytes, err = ioutil.ReadFile(s.baseDir + "/physical.track")
	assert.Equal(s.T(), nil, err)

	for i, b := range bytes {
		assert.Equal(s.T(), mach.Byte(b), s.enc.dst.Mem[i])
	}
}
