package a2

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pevans/erc/a2/a2enc"
	"github.com/pevans/erc/memory"
	"github.com/stretchr/testify/assert"
)

func TestNewDrive(t *testing.T) {
	d := NewDrive()

	assert.NotNil(t, d)
	assert.Equal(t, ReadMode, d.Mode)
	assert.Equal(t, a2enc.DOS33, d.ImageType)
}

func TestMotorOn(t *testing.T) {
	const cycles uint64 = 123

	d := NewDrive()

	assert.NotNil(t, d)
	assert.False(t, d.MotorOn())

	d.StartMotor(cycles)
	assert.True(t, d.MotorOn())
	assert.Equal(t, cycles, d.cyclesSinceLastSpin)

	d.StopMotor()
	assert.False(t, d.MotorOn())
}

func (s *a2Suite) TestDrivePosition() {
	d := NewDrive()

	s.Equal(0, d.Position())

	// In a zero track position, the drive position should be equal
	// exactly to the sector position.
	d.SectorPos = 123
	s.Equal(d.SectorPos, d.Position())

	// test track position
	d.TrackPos = 6
	s.Equal((a2enc.PhysTrackLen*d.TrackPos/2)+d.SectorPos, d.Position())
}

func (s *a2Suite) TestDriveShift() {
	d := NewDrive()

	d.SectorPos = 0

	// Positive shift
	d.Shift(10)
	s.Equal(10, d.SectorPos)

	// Negative shift
	d.Shift(-3)
	s.Equal(7, d.SectorPos)

	// We should not be able to shift below the zero boundary for a sector
	d.Shift(-10)
	s.Equal(a2enc.PhysTrackLen-3, d.SectorPos)

	// We can shift up but not including the length of a track
	d.Shift(a2enc.PhysTrackLen - 1)
	s.Equal(a2enc.PhysTrackLen-4, d.SectorPos)
	d.Shift(4)
	s.Equal(0, d.SectorPos)

	// And if the drive is locked, we shouldn't be able to shift at all
	d.Locked = true
	d.Shift(3)
	s.Equal(0, d.SectorPos)
}

func (s *a2Suite) TestDriveStep() {
	d := NewDrive()

	// Positive step, plus note that we always reset the sector position
	d.SectorPos = 123
	d.Step(2)
	s.Equal(2, d.TrackPos)
	s.Equal(123, d.SectorPos)

	// Negative step
	d.Step(-1)
	s.Equal(1, d.TrackPos)

	// No matter our starting point, if a step would go beyond MaxSteps,
	// we should be left _at_ the MaxSteps position
	d.Step(a2enc.MaxSteps + 1)
	s.Equal(a2enc.MaxSteps-1, d.TrackPos)

	// Any negative step that goes below zero should keep us at zero
	d.Step(-a2enc.MaxSteps * 2)
	s.Equal(0, d.TrackPos)
}

/*
func (s *a2Suite) TestDriveStepPhase() {
	d := NewDrive()

	// Any step based on a negative phase should do nothing
	d.Phase = 1
	d.TrackPos = 0
	d.StepPhase(2)
	s.Equal(0, d.TrackPos)
	s.Equal(1, d.Phase)

	// Positive nonoverflow step phase
	d.Phase = 1
	d.StepPhase(3)
	s.Equal(1, d.TrackPos)
	s.Equal(2, d.Phase)

	// Negative nonoverflow step phase
	d.StepPhase(1)
	s.Equal(0, d.TrackPos)
	s.Equal(1, d.Phase)

	// Negative overflow step phase
	d.TrackPos = 5
	d.StepPhase(7)
	s.Equal(4, d.TrackPos)
	s.Equal(4, d.Phase)

	// Positive overflow step phase
	d.StepPhase(1)
	s.Equal(5, d.TrackPos)
	s.Equal(1, d.Phase)
}
*/

func (s *a2Suite) TestImageType() {
	type test struct {
		fname string
		want  int
		efn   assert.ErrorAssertionFunc
	}

	cases := map[string]test{
		"do file": {
			fname: "something.do",
			want:  a2enc.DOS33,
			efn:   assert.NoError,
		},
		"dsk file": {
			fname: "something.dsk",
			want:  a2enc.DOS33,
			efn:   assert.NoError,
		},
		"nib file": {
			fname: "something.nib",
			want:  a2enc.Nibble,
			efn:   assert.NoError,
		},
		"po file": {
			fname: "something.po",
			want:  a2enc.ProDOS,
			efn:   assert.NoError,
		},
		"bad file": {
			fname: "bad",
			want:  -1,
			efn:   assert.Error,
		},
	}

	for desc, c := range cases {
		s.T().Run(desc, func(t *testing.T) {
			typ, err := ImageType(c.fname)
			assert.Equal(t, c.want, typ)
			c.efn(t, err)
		})
	}
}

func (s *a2Suite) TestDriveLoad() {
	d := NewDrive()

	data, _ := os.Open("../data/logical.disk")
	s.NoError(d.Load(data, "something.dsk"))

	s.Equal(a2enc.DOS33, d.ImageType)
	s.NotNil(d.Image)
	s.NotNil(d.Data)
}

func (s *a2Suite) TestDriveRead() {
	d := NewDrive()

	dat, _ := os.Open("../data/logical.disk")
	s.NoError(d.Load(dat, "something.dsk"))

	// Note that software expects the high bit to be set on all data coming
	// from the drive (so any test data needs Latch >= 0x80)
	d.Latch = 0x81
	d.newLatchData = true

	// With newLatchData is true, we should get the same value back unmodified
	s.Equal(uint8(0x81), d.Read())

	// Once you've read the latch, we unset newLatchData, and expect the
	// return value to be the same _except_ that the high bit is unset
	s.Equal(uint8(0x81&0x7F), d.Read())
}

func (s *a2Suite) TestDriveWrite() {
	d := NewDrive()

	dat, _ := os.Open("../data/logical.disk")
	s.NoError(d.Load(dat, "something.dsk"))

	d.Mode = WriteMode
	d.StartMotor(0)

	// If Latch < 0x80, Write should not write data, but position still shifts
	d.Latch = 0x11
	d.SectorPos = 0
	d.Write()
	s.NotEqual(d.Latch, d.Data.Get(d.Position()))

	// Write should do something here
	d.Latch = 0x81
	d.Write()
	s.Equal(d.Latch, d.Data.Get(d.Position()))
}

func (s *a2Suite) TestDriveSave() {
	logSeg := memory.NewSegment(a2enc.DosSize)
	for i := range logSeg.Size() {
		logSeg.Set(i, uint8(i%256))
	}

	encodedDOS33, err := a2enc.Encode(a2enc.DOS33, logSeg)
	s.NoError(err)

	encodedProDOS, err := a2enc.Encode(a2enc.ProDOS, logSeg)
	s.NoError(err)

	// Create a modified logical segment to test that changes to encoded data
	// are properly saved
	modifiedLogSeg := memory.NewSegment(a2enc.DosSize)
	for i := range modifiedLogSeg.Size() {
		modifiedLogSeg.Set(i, uint8(i%256))
	}

	modifiedLogSeg.Set(0, 0xAA)
	modifiedLogSeg.Set(100, 0xBB)
	modifiedLogSeg.Set(1000, 0xCC)

	encodedModified, err := a2enc.Encode(a2enc.DOS33, modifiedLogSeg)
	s.NoError(err)

	cases := []struct {
		name          string
		imageName     string
		imageType     int
		data          *memory.Segment
		expectedLog   *memory.Segment
		checkFileStat bool
		errfn         assert.ErrorAssertionFunc
	}{
		{
			name:          "empty ImageName does nothing",
			imageName:     "",
			imageType:     a2enc.DOS33,
			data:          encodedDOS33,
			expectedLog:   logSeg,
			checkFileStat: false,
			errfn:         assert.NoError,
		},
		{
			name:          "empty Data does nothing",
			imageName:     "something!",
			imageType:     a2enc.DOS33,
			data:          nil,
			expectedLog:   nil,
			checkFileStat: false,
			errfn:         assert.NoError,
		},
		{
			name:          "DOS33 saves successfully",
			imageName:     filepath.Join(s.T().TempDir(), "test_dos33.dsk"),
			imageType:     a2enc.DOS33,
			data:          encodedDOS33,
			expectedLog:   logSeg,
			checkFileStat: true,
			errfn:         assert.NoError,
		},
		{
			name:          "ProDOS saves successfully",
			imageName:     filepath.Join(s.T().TempDir(), "test_prodos.po"),
			imageType:     a2enc.ProDOS,
			data:          encodedProDOS,
			expectedLog:   logSeg,
			checkFileStat: true,
			errfn:         assert.NoError,
		},
		{
			name:          "modified data is saved correctly",
			imageName:     filepath.Join(s.T().TempDir(), "test_modified.dsk"),
			imageType:     a2enc.DOS33,
			data:          encodedModified,
			expectedLog:   modifiedLogSeg,
			checkFileStat: true,
			errfn:         assert.NoError,
		},
		{
			name:          "invalid image type returns error",
			imageName:     filepath.Join(s.T().TempDir(), "test_invalid.dsk"),
			imageType:     99,
			data:          encodedDOS33,
			expectedLog:   nil,
			checkFileStat: false,
			errfn:         assert.Error,
		},
	}

	for _, c := range cases {
		s.Run(c.name, func() {
			d := NewDrive()

			d.ImageName = c.imageName
			d.ImageType = c.imageType
			d.Image = logSeg
			d.Data = c.data

			err := d.Save()
			c.errfn(s.T(), err)

			if c.checkFileStat {
				_, statErr := os.Stat(c.imageName)
				s.NoError(statErr)

				// Read the saved file back and verify contents
				savedBytes, err := os.ReadFile(c.imageName)
				s.NoError(err)
				s.Equal(c.expectedLog.Size(), len(savedBytes))

				// Verify every byte matches the expected logical segment
				for i := range c.expectedLog.Size() {
					s.Equal(
						c.expectedLog.Get(i), savedBytes[i],
						"byte mismatch at offset %v", i,
					)
				}
			}
		})
	}
}
