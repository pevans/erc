package a2drive

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDriveRead(t *testing.T) {
	d := NewDrive()

	dat, err := os.Open("../../data/logical.disk")
	require.NoError(t, err)
	defer dat.Close()

	require.NoError(t, d.Load(dat, "something.dsk"))

	// Note that software expects the high bit to be set on all data coming
	// from the drive (so any test data needs Latch >= 0x80)
	d.latch = 0x81
	d.latchWasRead = false

	// With newLatchData is true, we should get the same value back unmodified
	assert.Equal(t, uint8(0x81), d.ReadLatch())

	// Once you've read the latch, we unset newLatchData, and expect the
	// return value to be the same _except_ that the high bit is unset
	assert.Equal(t, uint8(0x81&0x7F), d.ReadLatch())
}

func TestDriveWrite(t *testing.T) {
	d := NewDrive()

	dat, err := os.Open("../../data/logical.disk")
	require.NoError(t, err)
	defer dat.Close()

	require.NoError(t, d.Load(dat, "something.dsk"))

	d.SetWriteMode()
	d.StartMotor()

	// If Latch < 0x80, Write should not write data, but position still shifts
	d.latch = 0x11
	d.sectorPos = 0
	d.WriteLatch()
	assert.NotEqual(t, d.latch, d.data.Get(d.dataPosition()))

	// Write should do something here
	d.latch = 0x81
	d.WriteLatch()
	assert.Equal(t, d.latch, d.data.Get(d.dataPosition()))
}
