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
	assert.Equal(t, val, s.Mem[addr.Addr()].Val)
}

func TestGet(t *testing.T) {
	s := NewSegment(100)
	addr := DByte(1)
	val := Byte(123)

	s.Mem[addr.Addr()].Val = val
	assert.Equal(t, val, s.Get(addr))
}

func TestSetReadFunc(t *testing.T) {
	s := NewSegment(100)
	addr := DByte(1)
	val := Byte(123)

	err := s.SetReadFunc(addr, func(seg *Segment, adr Addressor) Byte {
		return val
	})

	assert.Equal(t, nil, err)
	assert.Equal(t, false, s.Mem[addr.Addr()].ReadFn == nil)

	assert.Equal(t, val, s.Get(addr))
}

func TestSetWriteFunc(t *testing.T) {
	s := NewSegment(100)
	addr := DByte(1)
	val := Byte(123)

	err := s.SetWriteFunc(addr, func(seg *Segment, adr Addressor, v Byte) {
		seg.Mem[adr.Addr()].Val = v + Byte(1)
	})

	assert.Equal(t, nil, err)
	assert.Equal(t, false, s.Mem[addr.Addr()].WriteFn == nil)

	s.Set(addr, val)
	assert.Equal(t, Byte(1)+val, s.Mem[addr.Addr()].Val)
}
