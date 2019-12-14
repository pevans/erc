package disk

import (
	"testing"

	"github.com/pevans/erc/pkg/data"
	"github.com/pevans/erc/pkg/sixtwo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type diskSuite struct {
	suite.Suite
	drive *Drive
}

func (s *diskSuite) SetupTest() {
	s.drive = NewDrive()
}

func TestNewDrive(t *testing.T) {
	drive := NewDrive()

	assert.NotEqual(t, nil, drive)
	assert.Equal(t, ReadMode, drive.Mode)
	assert.Equal(t, sixtwo.DOS33, drive.ImageType)
}

func TestDiskSuite(t *testing.T) {
	suite.Run(t, new(diskSuite))
}

func (s *diskSuite) TestPosition() {
	// When there is no valid segment in s.drive.Data, the position
	// should be zero.
	assert.Equal(s.T(), data.Int(0), s.drive.Position())

	s.drive.Data = data.NewSegment(sixtwo.PhysTrackLen * 2)
	cases := []struct {
		tpos int
		spos int
		want data.Int
	}{
		{0, 0, 0},
		{1, 0, 0},
		{1, 250, 250},
		{2, 0, sixtwo.PhysTrackLen},
		{2, 250, sixtwo.PhysTrackLen + 250},
		{34, 250, (sixtwo.PhysTrackLen * 17) + 250},
	}

	for _, c := range cases {
		s.drive.TrackPos = c.tpos
		s.drive.SectorPos = c.spos

		assert.Equal(s.T(), c.want, s.drive.Position())
	}
}

func (s *diskSuite) TestShift() {
	cases := []struct {
		locked bool
		spos   int
		offset int
		want   int
	}{
		{false, 0, 0, 0},
		{false, 1, 0, 1},
		{false, 1, 1, 2},
		{false, 1, sixtwo.PhysTrackLen, 0},
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

func (s *diskSuite) TestStep() {
	cases := []struct {
		tpos   int
		offset int
		want   int
	}{
		{0, 0, 0},
		{0, 1, 1},
		{2, -1, 1},
		{5, sixtwo.MaxSteps, sixtwo.MaxSteps},
		{5, -10, 0},
	}

	for _, c := range cases {
		s.drive.TrackPos = c.tpos
		s.drive.Step(c.offset)

		assert.Equal(s.T(), c.want, s.drive.TrackPos)
	}
}
