package sixtwo

/*
func TestEncoderSuite(t *testing.T) {
	suite.Run(t, new(encSuite))
}

func TestNewEncoder(t *testing.T) {
	seg := data.NewSegment(1)
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
		{DOS33, -1, 0},
		{DOS33, 16, 0},
		{DOS33, 0x0, 0x0},
		{DOS33, 0x1, 0x7},
		{DOS33, 0xE, 0x8},
		{DOS33, 0xF, 0xF},
		{ProDOS, 0x0, 0x0},
		{ProDOS, 0x1, 0x8},
		{ProDOS, 0xE, 0x7},
		{ProDOS, 0xF, 0xF},
		{Nibble, 1, 1},
	}

	for _, c := range cases {
		s.enc.imageType = c.imgType
		assert.Equal(s.T(), c.want, LogicalSector(s.enc.imageType, c.psect))
	}
}

func (s *encSuite) TestEncodeNIB() {
	_, _ = s.enc.src.CopySlice(0, []data.Byte{0x1, 0x2, 0x3})

	dst, err := s.enc.EncodeNIB()
	assert.Equal(s.T(), nil, err)

	for i := 0; i < dst.Size(); i++ {
		assert.Equal(s.T(), s.enc.src.Mem[i], dst.Mem[i])
	}
}

func (s *encSuite) TestEncodeDOS() {
	err := loadFile(s.enc.src, s.baseDir+"/logical.disk")
	assert.Equal(s.T(), nil, err)

	dst, err := s.enc.EncodeDOS()
	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), true, fileMatches(dst, s.baseDir+"/physical.disk"))
}

func (s *encSuite) TestWrite() {
	bytes := []data.Byte{0x1, 0x2, 0x3}
	_, _ = s.enc.src.CopySlice(0, bytes)

	assert.Equal(s.T(), 3, s.enc.Write(0, bytes))

	for i := 0; i < len(bytes); i++ {
		assert.Equal(s.T(), s.enc.src.Mem[i], s.enc.dst.Mem[i])
	}
}

func (s *encSuite) TestEncode4n4() {
	cases := []struct {
		in    data.Byte
		want1 data.Byte
		want2 data.Byte
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
	err := loadFile(s.enc.src, s.baseDir+"/logical.sector")
	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), PhysSectorLen, s.enc.EncodeSector(0, 0, 0, 0))
	assert.Equal(s.T(), true, fileMatches(s.enc.dst, s.baseDir+"/physical.sector"))
}

func (s *encSuite) TestEncodeTrack() {
	err := loadFile(s.enc.src, s.baseDir+"/logical.track")
	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), PhysTrackLen, s.enc.EncodeTrack(0, 0))
	assert.Equal(s.T(), true, fileMatches(s.enc.dst, s.baseDir+"/physical.track"))
}
*/
