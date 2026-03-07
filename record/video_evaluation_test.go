package record_test

import (
	"image/color"
	"testing"

	"github.com/pevans/erc/gfx"
	"github.com/pevans/erc/record"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeVideoRecorderWithFrame(t *testing.T, w, h uint, fill func(fb *gfx.FrameBuffer)) *record.VideoRecorder {
	t.Helper()

	fb := gfx.NewFrameBuffer(w, h)
	fill(fb)

	vr := record.NewVideoRecorder(fb)
	vr.CaptureAt(1)
	vr.Observe(1)

	return vr
}

func TestEvaluateVideoAssertions_screenPass(t *testing.T) {
	black := color.RGBA{R: 0, G: 0, B: 0, A: 0xff}
	white := color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}

	vr := makeVideoRecorderWithFrame(t, 4, 2, func(fb *gfx.FrameBuffer) {
		// Row 0: black black white white Row 1: white white black black
		fb.ClearCells(black)
		_ = fb.SetCell(2, 0, white)
		_ = fb.SetCell(3, 0, white)
		_ = fb.SetCell(0, 1, white)
		_ = fb.SetCell(1, 1, white)
	})

	lines := []string{
		"step 1: video screen 4x2",
		"colors: . = 000000, # = FFFFFF",
		"..##",
		"##..",
	}

	a, err := record.ParseVideoAssertion(lines)
	require.NoError(t, err)

	results := record.EvaluateVideoAssertions([]record.VideoAssertion{a}, vr)
	require.Len(t, results, 1)
	assert.True(t, results[0].Passed, "expected pass, got failure: %+v", results[0].Failure)
}

func TestEvaluateVideoAssertions_screenMismatch(t *testing.T) {
	black := color.RGBA{R: 0, G: 0, B: 0, A: 0xff}

	vr := makeVideoRecorderWithFrame(t, 4, 2, func(fb *gfx.FrameBuffer) {
		fb.ClearCells(black)
	})

	lines := []string{
		"step 1: video screen 4x2",
		"colors: . = 000000, # = FFFFFF",
		"..#.",
		"....",
	}

	a, err := record.ParseVideoAssertion(lines)
	require.NoError(t, err)

	results := record.EvaluateVideoAssertions([]record.VideoAssertion{a}, vr)
	require.Len(t, results, 1)
	assert.False(t, results[0].Passed)

	m := results[0].Failure
	require.NotNil(t, m)
	assert.Equal(t, 2, m.X)
	assert.Equal(t, 0, m.Y)
	assert.Equal(t, byte('#'), m.Expected)
	assert.Equal(t, byte('.'), m.Actual)
}

func TestEvaluateVideoAssertions_unmappedColor(t *testing.T) {
	red := color.RGBA{R: 0xff, A: 0xff}

	vr := makeVideoRecorderWithFrame(t, 2, 1, func(fb *gfx.FrameBuffer) {
		fb.ClearCells(red)
	})

	lines := []string{
		"step 1: video screen 2x1",
		"colors: . = 000000, # = FFFFFF",
		".#",
	}

	a, err := record.ParseVideoAssertion(lines)
	require.NoError(t, err)

	results := record.EvaluateVideoAssertions([]record.VideoAssertion{a}, vr)
	require.Len(t, results, 1)
	assert.False(t, results[0].Passed)
	assert.True(t, results[0].Failure.UnmappedColor)
}

func TestEvaluateVideoAssertions_regionPass(t *testing.T) {
	black := color.RGBA{R: 0, G: 0, B: 0, A: 0xff}
	white := color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}

	// 8x4 framebuffer, grid is 4x2, so each grid cell = 2x2 framebuffer
	// pixels
	vr := makeVideoRecorderWithFrame(t, 8, 4, func(fb *gfx.FrameBuffer) {
		fb.ClearCells(black)
		// Set grid cell (2, 1) to white -- that's fb pixels (4,2) and (5,2)
		// etc.
		for dy := range 2 {
			for dx := range 2 {
				_ = fb.SetCell(uint(4+dx), uint(2+dy), white)
			}
		}
	})

	lines := []string{
		"step 1: video region 2,1 2x1 4x2",
		"colors: . = 000000, # = FFFFFF",
		"#.",
	}

	a, err := record.ParseVideoAssertion(lines)
	require.NoError(t, err)

	results := record.EvaluateVideoAssertions([]record.VideoAssertion{a}, vr)
	require.Len(t, results, 1)
	assert.True(t, results[0].Passed, "expected pass, got failure: %+v", results[0].Failure)
}

func TestEvaluateVideoAssertions_rowPass(t *testing.T) {
	black := color.RGBA{R: 0, G: 0, B: 0, A: 0xff}
	white := color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}

	// 4x4 framebuffer, grid 4x4 (1:1 mapping)
	vr := makeVideoRecorderWithFrame(t, 4, 4, func(fb *gfx.FrameBuffer) {
		fb.ClearCells(black)
		_ = fb.SetCell(1, 2, white)
		_ = fb.SetCell(2, 2, white)
	})

	lines := []string{
		"step 1: video row 2 4x4",
		"colors: . = 000000, # = FFFFFF",
		".##.",
	}

	a, err := record.ParseVideoAssertion(lines)
	require.NoError(t, err)

	results := record.EvaluateVideoAssertions([]record.VideoAssertion{a}, vr)
	require.Len(t, results, 1)
	assert.True(t, results[0].Passed, "expected pass, got failure: %+v", results[0].Failure)
}

func TestEvaluateVideoAssertions_missingFrame(t *testing.T) {
	fb := gfx.NewFrameBuffer(4, 4)
	vr := record.NewVideoRecorder(fb)

	lines := []string{
		"step 1: video screen 4x4",
		"colors: . = 000000",
		"....",
		"....",
		"....",
		"....",
	}

	a, err := record.ParseVideoAssertion(lines)
	require.NoError(t, err)

	results := record.EvaluateVideoAssertions([]record.VideoAssertion{a}, vr)
	require.Len(t, results, 1)
	assert.False(t, results[0].Passed)
	assert.True(t, results[0].Failure.MissingFrame)
}

func TestEvaluateVideoAssertions_sampling(t *testing.T) {
	black := color.RGBA{R: 0, G: 0, B: 0, A: 0xff}
	white := color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}

	// 8x4 framebuffer sampled at 4x2 grid -- each grid cell maps to every
	// other framebuffer pixel
	vr := makeVideoRecorderWithFrame(t, 8, 4, func(fb *gfx.FrameBuffer) {
		fb.ClearCells(black)
		// Grid cell (1,0) samples fb pixel (2,0)
		_ = fb.SetCell(2, 0, white)
	})

	lines := []string{
		"step 1: video screen 4x2",
		"colors: . = 000000, # = FFFFFF",
		".#..",
		"....",
	}

	a, err := record.ParseVideoAssertion(lines)
	require.NoError(t, err)

	results := record.EvaluateVideoAssertions([]record.VideoAssertion{a}, vr)
	require.Len(t, results, 1)
	assert.True(t, results[0].Passed, "expected pass, got failure: %+v", results[0].Failure)
}
