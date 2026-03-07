package record_test

import (
	"image/color"
	"testing"

	"github.com/pevans/erc/record"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseVideoAssertion_screen(t *testing.T) {
	lines := []string{
		"step 512: video screen 4x3",
		"colors: . = 000000, # = FFFFFF",
		"..##",
		".##.",
		"##..",
	}

	a, err := record.ParseVideoAssertion(lines)
	require.NoError(t, err)

	assert.Equal(t, 512, a.Step)
	assert.Equal(t, record.VideoScreen, a.Kind)
	assert.Equal(t, 4, a.GridW)
	assert.Equal(t, 3, a.GridH)
	assert.Len(t, a.Expected, 3)
	assert.Equal(t, "..##", a.Expected[0])

	assert.Equal(t, color.RGBA{R: 0, G: 0, B: 0, A: 0xff}, a.Legend.Forward['.'])
	assert.Equal(t, color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}, a.Legend.Forward['#'])
}

func TestParseVideoAssertion_region(t *testing.T) {
	lines := []string{
		"step 100: video region 10,5 3x2 280x192",
		"colors: . = 000000, # = FFFFFF",
		"#.#",
		".#.",
	}

	a, err := record.ParseVideoAssertion(lines)
	require.NoError(t, err)

	assert.Equal(t, 100, a.Step)
	assert.Equal(t, record.VideoRegion, a.Kind)
	assert.Equal(t, 10, a.RegionX)
	assert.Equal(t, 5, a.RegionY)
	assert.Equal(t, 3, a.RegionW)
	assert.Equal(t, 2, a.RegionH)
	assert.Equal(t, 280, a.GridW)
	assert.Equal(t, 192, a.GridH)
	assert.Len(t, a.Expected, 2)
}

func TestParseVideoAssertion_row(t *testing.T) {
	lines := []string{
		"step 200: video row 96 4x3",
		"colors: . = 000000, # = FFFFFF",
		".##.",
	}

	a, err := record.ParseVideoAssertion(lines)
	require.NoError(t, err)

	assert.Equal(t, 200, a.Step)
	assert.Equal(t, record.VideoRow, a.Kind)
	assert.Equal(t, 96, a.RowIndex)
	assert.Equal(t, 4, a.GridW)
	assert.Equal(t, 3, a.GridH)
	assert.Len(t, a.Expected, 1)
	assert.Equal(t, ".##.", a.Expected[0])
}

func TestParseColorLegend(t *testing.T) {
	lines := []string{
		"step 1: video screen 2x1",
		"colors: . = 000000, # = FFFFFF, P = D043E5",
		".#",
	}

	a, err := record.ParseVideoAssertion(lines)
	require.NoError(t, err)

	assert.Equal(t, color.RGBA{R: 0xD0, G: 0x43, B: 0xE5, A: 0xff}, a.Legend.Forward['P'])

	// Reverse lookup
	ch, ok := a.Legend.Reverse[color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}]
	assert.True(t, ok)
	assert.Equal(t, byte('#'), ch)
}

func TestParseVideoAssertion_errors(t *testing.T) {
	cases := []struct {
		name  string
		lines []string
	}{
		{
			name:  "too few lines",
			lines: []string{"step 1: video screen 4x3"},
		},
		{
			name: "bad header",
			lines: []string{
				"not a step",
				"colors: . = 000000",
				"....",
			},
		},
		{
			name: "bad colors",
			lines: []string{
				"step 1: video screen 4x1",
				"not colors",
				"....",
			},
		},
		{
			name: "wrong row count",
			lines: []string{
				"step 1: video screen 4x2",
				"colors: . = 000000",
				"....",
			},
		},
		{
			name: "wrong row width",
			lines: []string{
				"step 1: video screen 4x1",
				"colors: . = 000000",
				"...",
			},
		},
		{
			name: "bad hex color",
			lines: []string{
				"step 1: video screen 2x1",
				"colors: . = ZZZZZZ",
				"..",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := record.ParseVideoAssertion(tc.lines)
			assert.Error(t, err)
		})
	}
}
