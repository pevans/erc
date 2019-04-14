package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSegment(t *testing.T) {
	assert.NotEqual(t, nil, NewSegment(1))
}

func TestSet(t *testing.T) {
	s := NewSegment(100)
	addr := DByte(1)
	val := Byte(123)

	s.Set(addr, val)
	assert.Equal(t, val, s.Mem[addr.Addr()])
}

func TestGet(t *testing.T) {
	s := NewSegment(100)
	addr := DByte(1)
	val := Byte(123)

	s.Mem[addr.Addr()] = val
	assert.Equal(t, val, s.Get(addr))
}

func TestByteSlice(t *testing.T) {
	cases := []struct {
		want []Byte
		in   []byte
	}{
		{[]Byte{}, []byte{}},
		{[]Byte{1, 2, 3}, []byte{1, 2, 3}},
	}

	for _, c := range cases {
		assert.Equal(t, c.want, ByteSlice(c.in))
	}
}

func TestSize(t *testing.T) {
	s := NewSegment(100)
	assert.Equal(t, 100, s.Size())
}

func TestCopySlice(t *testing.T) {
	cases := []struct {
		wantWritten int
		wantError   bool
		segmentSize int
		start       int
		byteSlice   []Byte
	}{
		{0, true, 1, -1, []Byte{}},
		{0, true, 0, 1, ByteSlice([]byte{1, 2})},
		{3, false, 5, 0, ByteSlice([]byte{1, 2, 3})},
	}

	for _, c := range cases {
		s := NewSegment(c.segmentSize)

		writ, err := s.CopySlice(c.start, c.byteSlice)
		assert.Equal(t, c.wantWritten, writ)

		if c.wantError {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
	}
}
