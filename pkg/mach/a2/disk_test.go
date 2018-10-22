package a2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDiskDrive(t *testing.T) {
	drive := NewDiskDrive()

	assert.NotEqual(t, nil, drive)
	assert.Equal(t, DDRead, drive.Mode)
	assert.Equal(t, DDDOS33, drive.ImageType)
}

func (s *a2Suite) TestLogicalSector() {
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
		s.drive.ImageType = c.imgType
		assert.Equal(s.T(), c.want, s.drive.LogicalSector(c.psect))
	}
}
