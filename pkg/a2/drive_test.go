package a2

import (
	"os"
	"testing"

	"github.com/pevans/erc/pkg/data"
	"github.com/pevans/erc/pkg/sixtwo"
	"github.com/stretchr/testify/assert"
)

func TestNewDrive(t *testing.T) {
	d := NewDrive()

	assert.NotNil(t, d)
	assert.Equal(t, ReadMode, d.Mode)
	assert.Equal(t, sixtwo.DOS33, d.ImageType)
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
	s.Equal((sixtwo.PhysTrackLen*d.TrackPos/2)+d.SectorPos, d.Position())
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
	s.Equal(0, d.SectorPos)

	// We can shift up but not including the length of a track
	d.Shift(sixtwo.PhysTrackLen - 1)
	s.Equal(sixtwo.PhysTrackLen-1, d.SectorPos)
	d.Shift(1)
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
	s.Equal(0, d.SectorPos)

	// Negative step
	d.Step(-1)
	s.Equal(1, d.TrackPos)

	// No matter our starting point, if a step would go beyond MaxSteps,
	// we should be left _at_ the MaxSteps position
	d.Step(sixtwo.MaxSteps + 1)
	s.Equal(sixtwo.MaxSteps, d.TrackPos)

	// Any negative step that goes below zero should keep us at zero
	d.Step(-sixtwo.MaxSteps * 2)
	s.Equal(0, d.TrackPos)
}

func (s *a2Suite) TestDrivePhase() {
	for i := 0x0; i < 0x10; i++ {
		p := Phase(data.DByte(i))
		switch i {
		case 0x1:
			s.Equal(1, p)
		case 0x3:
			s.Equal(2, p)
		case 0x5:
			s.Equal(3, p)
		case 0x7:
			s.Equal(4, p)
		default:
			s.Equal(-1, p)
		}
	}
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
			want:  sixtwo.DOS33,
			efn:   assert.NoError,
		},
		"dsk file": {
			fname: "something.dsk",
			want:  sixtwo.DOS33,
			efn:   assert.NoError,
		},
		"nib file": {
			fname: "something.nib",
			want:  Nibble,
			efn:   assert.NoError,
		},
		"po file": {
			fname: "something.po",
			want:  sixtwo.ProDOS,
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

	data, _ := os.Open("../../data/logical.disk")
	s.NoError(d.Load(data, "something.dsk"))

	s.Equal(sixtwo.DOS33, d.ImageType)
	s.NotNil(d.Image)
	s.NotNil(d.Data)
}

func (s *a2Suite) TestDriveRead() {
	d := NewDrive()

	dat, _ := os.Open("../../data/logical.disk")
	s.NoError(d.Load(dat, "something.dsk"))

	d.Data.Set(d.Position(), 0x11)
	spos := d.SectorPos

	b := d.Read()

	s.Equal(data.Byte(0x11), b)
	s.Equal(data.Byte(0x11), d.Latch)
	s.Equal(spos+1, d.SectorPos)
}

func (s *a2Suite) TestDriveWrite() {
	d := NewDrive()

	dat, _ := os.Open("../../data/logical.disk")
	s.NoError(d.Load(dat, "something.dsk"))

	// If Latch < 0x80, Write should do nothing
	d.Latch = 0x11
	d.SectorPos = 0
	d.Write()
	s.Equal(0, d.SectorPos)
	s.NotEqual(d.Latch, d.Data.Get(d.Position()))

	// Write should do something here
	d.Latch = 0x81
	d.Write()
	s.Equal(1, d.SectorPos)
	s.Equal(d.Latch, d.Data.Get(d.Position()-1))
}
