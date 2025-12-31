package a2drive

import (
	"testing"

	"github.com/pevans/erc/a2/a2enc"
	"github.com/stretchr/testify/assert"
)

func TestDrivePosition(t *testing.T) {
	d := NewDrive()

	assert.Equal(t, 0, d.dataPosition())

	// In a zero track position, the drive position should be equal
	// exactly to the sector position.
	d.sectorPos = 123
	assert.Equal(t, d.sectorPos, d.dataPosition())

	// test track position
	d.trackPos = 6
	assert.Equal(t, (a2enc.PhysTrackLen*d.trackPos/2)+d.sectorPos, d.dataPosition())
}

func TestDriveShift(t *testing.T) {
	d := NewDrive()

	d.sectorPos = 0

	// Positive shift
	d.Shift(10)
	assert.Equal(t, 10, d.sectorPos)

	// Negative shift
	d.Shift(-3)
	assert.Equal(t, 7, d.sectorPos)

	// We should not be able to shift below the zero boundary for a sector
	d.Shift(-10)
	assert.Equal(t, a2enc.PhysTrackLen-3, d.sectorPos)

	// We can shift up but not including the length of a track
	d.Shift(a2enc.PhysTrackLen - 1)
	assert.Equal(t, a2enc.PhysTrackLen-4, d.sectorPos)
	d.Shift(4)
	assert.Equal(t, 0, d.sectorPos)
}

func TestDriveStep(t *testing.T) {
	d := NewDrive()

	// Positive step, plus note that we always reset the sector position
	d.sectorPos = 123
	d.Step(2)
	assert.Equal(t, 2, d.trackPos)
	assert.Equal(t, 123, d.sectorPos)

	// Negative step
	d.Step(-1)
	assert.Equal(t, 1, d.trackPos)

	// No matter our starting point, if a step would go beyond MaxSteps,
	// we should be left _at_ the MaxSteps position
	d.Step(a2enc.MaxSteps + 1)
	assert.Equal(t, a2enc.MaxSteps-1, d.trackPos)

	// Any negative step that goes below zero should keep us at zero
	d.Step(-a2enc.MaxSteps * 2)
	assert.Equal(t, 0, d.trackPos)
}
