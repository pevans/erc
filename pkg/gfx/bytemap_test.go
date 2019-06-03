package gfx

import (
	"image"
	"image/color"
	"testing"

	"github.com/pevans/erc/pkg/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockDotDrawer struct {
	mock.Mock
}

func (m *mockDotDrawer) DrawDot(coord image.Point, color color.RGBA) {
	m.Called(coord, color)
}

func TestBytemapValid(t *testing.T) {
	type test struct {
		width   int
		height  int
		dotSize int
		boolFn  assert.BoolAssertionFunc
	}

	cases := map[string]test{
		"square":     {3, 3, 9, assert.True},
		"square bad": {3, 3, 8, assert.False},
		"rect":       {8, 7, 56, assert.True},
	}

	for desc, c := range cases {
		t.Run(desc, func(t *testing.T) {
			bmap := Bytemap{
				Width:  c.width,
				Height: c.height,
				Dots:   make([]data.Byte, c.dotSize),
			}

			c.boolFn(t, bmap.Valid())
		})
	}
}

func TestBytemapDraw(t *testing.T) {
	type test struct {
		bmap  Bytemap
		coord image.Point
		color color.RGBA
	}

	m := new(mockDotDrawer)

	m.On("DrawDot", mock.Anything, mock.Anything).
		Return()

	cases := map[string]test{
		"blank": {
			bmap: Bytemap{
				Width:  2,
				Height: 2,
				Dots: []data.Byte{
					0, 1,
					1, 0,
				},
			},
			coord: image.Point{X: 1, Y: 1},
			color: color.RGBA{R: 1, G: 0, B: 0, A: 0},
		},

		"invalid": {
			bmap: Bytemap{
				Width:  2,
				Height: 2,
				Dots:   []data.Byte{0},
			},

			coord: image.Point{X: 0, Y: 0},
			color: color.RGBA{R: 2, G: 1, B: 0, A: 0},
		},
	}

	totalDots := 0

	for desc, c := range cases {
		if c.bmap.Valid() {
			totalDots += len(c.bmap.Dots)
		}

		t.Run(desc, func(t *testing.T) {
			c.bmap.Draw(m, c.coord, c.color)
		})
	}

	m.AssertNumberOfCalls(t, "DrawDot", totalDots)
}
