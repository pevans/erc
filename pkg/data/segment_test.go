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
	type test struct {
		want []Byte
		in   []byte
	}

	cases := map[string]test{
		"empty":    {want: []Byte{}, in: []byte{}},
		"nonempty": {want: []Byte{1, 2, 3}, in: []byte{1, 2, 3}},
	}

	for desc, c := range cases {
		t.Run(desc, func(t *testing.T) {
			assert.Equal(t, c.want, ByteSlice(c.in))
		})
	}
}

func TestSize(t *testing.T) {
	s := NewSegment(100)
	assert.Equal(t, 100, s.Size())
}

func TestCopySlice(t *testing.T) {
	type test struct {
		written   int
		size      int
		start     int
		byteSlice []Byte
		errfn     assert.ErrorAssertionFunc
	}

	cases := map[string]test{
		"negative start": {
			written:   0,
			size:      1,
			start:     -1,
			byteSlice: []Byte{},
			errfn:     assert.Error,
		},
		"no size": {
			written:   0,
			size:      0,
			start:     1,
			byteSlice: ByteSlice([]byte{1, 2}),
			errfn:     assert.Error,
		},
		"normal": {
			written:   3,
			size:      5,
			start:     0,
			byteSlice: ByteSlice([]byte{1, 2, 3}),
			errfn:     assert.NoError,
		},
	}

	for desc, c := range cases {
		t.Run(desc, func(t *testing.T) {
			s := NewSegment(c.size)

			writ, err := s.CopySlice(c.start, c.byteSlice)
			assert.Equal(t, c.written, writ)
			c.errfn(t, err)
		})
	}
}
