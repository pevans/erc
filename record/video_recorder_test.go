package record_test

import (
	"image/color"
	"testing"

	"github.com/pevans/erc/gfx"
	"github.com/pevans/erc/record"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVideoRecorder_captureAtDeclaredSteps(t *testing.T) {
	fb := gfx.NewFrameBuffer(4, 4)
	vr := record.NewVideoRecorder(fb)
	vr.CaptureAt(3, 5)

	for step := 1; step <= 6; step++ {
		vr.Before()
		fb.ClearCells(color.RGBA{R: uint8(step * 40), A: 0xff})
		vr.Observe(step)
	}

	assert.Nil(t, vr.Frame(1))
	assert.Nil(t, vr.Frame(2))
	require.NotNil(t, vr.Frame(3))
	assert.Nil(t, vr.Frame(4))
	require.NotNil(t, vr.Frame(5))
	assert.Nil(t, vr.Frame(6))
}

func TestVideoRecorder_snapshotsAreIndependent(t *testing.T) {
	fb := gfx.NewFrameBuffer(2, 2)
	vr := record.NewVideoRecorder(fb)
	vr.CaptureAt(1, 2)

	// Step 1: set to red
	fb.ClearCells(color.RGBA{R: 0xff, A: 0xff})
	vr.Observe(1)

	// Step 2: set to blue
	fb.ClearCells(color.RGBA{B: 0xff, A: 0xff})
	vr.Observe(2)

	// Frame 1 should still be red
	f1 := vr.Frame(1)
	require.NotNil(t, f1)
	c1 := f1.GetPixel(0, 0)
	assert.Equal(t, uint8(0xff), c1.R)
	assert.Equal(t, uint8(0), c1.B)

	// Frame 2 should be blue
	f2 := vr.Frame(2)
	require.NotNil(t, f2)
	c2 := f2.GetPixel(0, 0)
	assert.Equal(t, uint8(0), c2.R)
	assert.Equal(t, uint8(0xff), c2.B)
}

func TestVideoRecorder_withRecorder(t *testing.T) {
	fb := gfx.NewFrameBuffer(2, 2)
	vr := record.NewVideoRecorder(fb)
	vr.CaptureAt(2)

	var r record.Recorder
	r.Add(vr)

	r.Step(func() {
		fb.ClearCells(color.RGBA{R: 0x10, A: 0xff})
	})
	r.Step(func() {
		fb.ClearCells(color.RGBA{R: 0x20, A: 0xff})
	})

	assert.Nil(t, vr.Frame(1))
	require.NotNil(t, vr.Frame(2))

	c := vr.Frame(2).GetPixel(0, 0)
	assert.Equal(t, uint8(0x20), c.R)
}
