package data

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSegment(t *testing.T) {
	assert.NotEqual(t, nil, NewSegment(1))
}

func TestSet(t *testing.T) {
	s := NewSegment(100)
	addr := 1
	val := Byte(123)

	s.Set(addr, val)
	assert.Equal(t, val, s.Mem[addr])

	assert.Panics(t, func() {
		s.Set(-1, val)
	})

	assert.Panics(t, func() {
		s.Set(cap(s.Mem)+1, val)
	})
}

func TestGet(t *testing.T) {
	s := NewSegment(100)
	addr := 1
	val := Byte(123)

	s.Mem[addr] = val
	assert.Equal(t, val, s.Get(addr))

	assert.Panics(t, func() {
		_ = s.Get(-1)
	})

	assert.Panics(t, func() {
		_ = s.Get(cap(s.Mem) + 1)
	})
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

func TestWriteFile(t *testing.T) {
	const (
		size  = 1
		value = 0x33
		file  = "/tmp/segment.writefile"
	)

	s := NewSegment(size)

	s.Set(0, value)

	// We should not be able to write a bad file
	assert.Error(t, s.WriteFile(""))

	// But should be able to write a good file
	assert.NoError(t, s.WriteFile(file))

	// We should be able to see the file, if we looked...
	ns := NewSegment(size)
	assert.NoError(t, ns.ReadFile(file))
	assert.Equal(t, Byte(value), ns.Get(0))
	os.Remove(file)
}

func TestReadFile(t *testing.T) {
	s := NewSegment(256)

	// See if we return an error for a bad file of some kind
	assert.Error(t, s.ReadFile(""))

	// See that we don't return an error for a real file
	assert.NoError(t, s.ReadFile("../../data/logical.sector"))

	// Make sure we have some real data
	assert.NotEqual(t, Byte(0), s.Get(0))
}
