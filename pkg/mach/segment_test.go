package mach

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
