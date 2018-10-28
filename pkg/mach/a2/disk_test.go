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

func (s *a2Suite) TestPosition() {
	// When there is no valid segment in s.drive.Data, the position
	// should be zero.
	assert.Equal(s.T(), 0, s.drive.Position())

	s.drive.Data = mach.NewSegment(PhysTrackLen * 2)
	cases := []struct {
		tpos int
		spos int
		want int
	}{
		{0, 0, 0},
		{1, 0, 0},
		{1, 250, 250},
		{2, 0, PhysTrackLen},
		{2, 250, PhysTrackLen + 250},
		{34, 250, (PhysTrackLen * 17) + 250},
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
		{false, 1, PhysTrackLen, 0},
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

func (s *a2Suite) TestStep() {
	cases := []struct {
		tpos   int
		offset int
		want   int
	}{
		{0, 0, 0},
		{0, 1, 1},
		{2, -1, 1},
		{5, DDMaxSteps, DDMaxSteps},
		{5, -10, 0},
	}

	for _, c := range cases {
		s.drive.TrackPos = c.tpos
		s.drive.Step(c.offset)

		assert.Equal(s.T(), c.want, s.drive.TrackPos)
	}
}
