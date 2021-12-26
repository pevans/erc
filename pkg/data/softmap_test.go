package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	softTestVal  uint8 = 123
	softTestAddr int   = 111
	softTestSize int   = 256
)

func TestSoftMapRead(t *testing.T) {
	sm := NewSoftMap(softTestSize)
	val := softTestVal
	fn := func(x int, _ *StateMap) uint8 {
		return val
	}

	_, ok := sm.Read(softTestAddr)
	assert.False(t, ok)

	sm.SetRead(softTestAddr, fn)
	assert.NotNil(t, sm.reads[softTestAddr])

	rval, ok := sm.Read(softTestAddr)
	assert.True(t, ok)
	assert.Equal(t, softTestVal, rval)
}

func TestSoftMapWrite(t *testing.T) {
	var val uint8

	sm := NewSoftMap(softTestSize)
	fn := func(x int, y uint8, _ *StateMap) {
		val = y
	}

	ok := sm.Write(softTestAddr, softTestVal)
	assert.False(t, ok)
	assert.Zero(t, val)

	sm.SetWrite(softTestAddr, fn)
	assert.NotNil(t, sm.writes[softTestAddr])

	ok = sm.Write(softTestAddr, softTestVal)
	assert.True(t, ok)
	assert.Equal(t, softTestVal, val)
}
