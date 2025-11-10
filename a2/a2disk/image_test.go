package a2disk_test

import (
	"testing"

	"github.com/pevans/erc/a2/a2disk"
	"github.com/pevans/erc/a2/a2enc"
	"github.com/pevans/erc/memory"
	"github.com/stretchr/testify/assert"
)

func TestNewImage(t *testing.T) {
	assert.NotNil(t, a2disk.NewImage())
}

func TestImage_Parse(t *testing.T) {
	cases := []struct {
		name  string
		size  int
		setup func(*memory.Segment)
		errfn assert.ErrorAssertionFunc
	}{
		{
			name: "valid full DOS 3.3 disk image",
			size: a2enc.DosSize,
			setup: func(seg *memory.Segment) {
				for i := range seg.Size() {
					seg.Set(i, uint8(i%256))
				}
			},
			errfn: assert.NoError,
		},
		{
			name: "segment too small",
			size: a2enc.LogTrackLen - 1,
			setup: func(seg *memory.Segment) {
				for i := range seg.Size() {
					seg.Set(i, 0xFF)
				}
			},
			errfn: assert.Error,
		},
		{
			name: "empty segment",
			size: 0,
			setup: func(seg *memory.Segment) {
			},
			errfn: assert.Error,
		},
		{
			name: "single track",
			size: a2enc.LogTrackLen,
			setup: func(seg *memory.Segment) {
				for i := range seg.Size() {
					seg.Set(i, 0xAA)
				}
			},
			errfn: assert.Error,
		},
		{
			name: "partial disk image",
			size: a2enc.LogTrackLen * 20,
			setup: func(seg *memory.Segment) {
				for i := range seg.Size() {
					seg.Set(i, uint8(i&0xFF))
				}
			},
			errfn: assert.Error,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			seg := memory.NewSegment(c.size)
			c.setup(seg)

			img := a2disk.NewImage()
			err := img.Parse(seg)
			c.errfn(t, err)
		})
	}
}
