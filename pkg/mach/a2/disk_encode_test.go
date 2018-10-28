package a2

import (
	"testing"

	"github.com/pevans/erc/pkg/mach"
	"github.com/stretchr/testify/assert"
)

func TestNewEncoder(t *testing.T) {
	seg := mach.NewSegment(1)
	typ := 3

	enc := NewEncoder(typ, seg)
	assert.NotEqual(t, nil, enc)
	assert.Equal(t, seg, enc.src)
	assert.Equal(t, typ, enc.imageType)
}

func (s *a2Suite) TestLogicalSector() {
	cases := []struct {
		imgType int
		psect   int
		want    int
	}{
		{0, 0, 0},
		{DDDOS33, -1, 0},
		{DDDOS33, 16, 0},
		{DDDOS33, 0x0, 0x0},
		{DDDOS33, 0x1, 0x7},
		{DDDOS33, 0xE, 0x8},
		{DDDOS33, 0xF, 0xF},
		{DDProDOS, 0x0, 0x0},
		{DDProDOS, 0x1, 0x8},
		{DDProDOS, 0xE, 0x7},
		{DDProDOS, 0xF, 0xF},
		{DDNibble, 1, 1},
	}

	seg := mach.NewSegment(100)
	for _, c := range cases {
		enc := NewEncoder(c.imgType, seg)
		assert.Equal(s.T(), c.want, enc.LogicalSector(c.psect))
	}
}

func TestEncodeNIB(t *testing.T) {
	seg := mach.NewSegment(3)
	_, _ = seg.CopySlice(0, []mach.Byte{0x1, 0x2, 0x3})

	enc := NewEncoder(DDNibble, seg)
	dst, err := enc.EncodeNIB()
	assert.Equal(t, nil, err)

	for i := 0; i < dst.Size(); i++ {
		assert.Equal(t, seg.Mem[i], dst.Mem[i])
	}
}

func TestWrite(t *testing.T) {
	seg := mach.NewSegment(3)
	enc := NewEncoder(0, seg)
	enc.dst = mach.NewSegment(3)

	bytes := []mach.Byte{0x1, 0x2, 0x3}
	_, _ = seg.CopySlice(0, bytes)

	assert.Equal(t, 3, enc.Write(0, bytes))

	for i := 0; i < enc.dst.Size(); i++ {
		assert.Equal(t, seg.Mem[i], enc.dst.Mem[i])
	}
}
