package a2

import (
	"testing"

	"github.com/pevans/erc/pkg/mach"
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

func (s *a2Suite) TestPosition() {
	// When there is no valid segment in s.drive.Data, the position
	// should be zero.
	assert.Equal(s.T(), 0, s.drive.Position())

	s.drive.Data = mach.NewSegment(EncTrackLen * 2)
	cases := []struct {
		tpos int
		spos int
		want int
	}{
		{0, 0, 0},
		{1, 0, 0},
		{1, 250, 250},
		{2, 0, EncTrackLen},
		{2, 250, EncTrackLen + 250},
		{34, 250, (EncTrackLen * 17) + 250},
	}

	for _, c := range cases {
		s.drive.TrackPos = c.tpos
		s.drive.SectorPos = c.spos

		assert.Equal(s.T(), c.want, s.drive.Position())
	}
}

func (s *a2Suite) TestShift() {
	cases := []struct {
		locked bool
		spos   int
		offset int
		want   int
	}{
		{false, 0, 0, 0},
		{false, 1, 0, 1},
		{false, 1, 1, 2},
		{false, 1, EncTrackLen, 0},
		{false, 2, -1, 1},
		{false, 2, -3, 0},
		{true, 1, 1, 1},
	}

	for _, c := range cases {
		s.drive.Locked = c.locked
		s.drive.SectorPos = c.spos

		s.drive.Shift(c.offset)
		assert.Equal(s.T(), c.want, s.drive.SectorPos)
	}
}
