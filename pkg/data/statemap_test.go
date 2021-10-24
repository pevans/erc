package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStatemap(t *testing.T) {
	sm := NewStateMap()
	assert.NotNil(t, sm)
	assert.NotNil(t, sm.m)
}

func TestStateMapInt(t *testing.T) {
	var (
		sm   = NewStateMap()
		k    = 1
		v    = 123
		zero = 0
	)

	assert.Equal(t, zero, sm.Int(k))

	sm.SetInt(k, v)
	assert.Equal(t, v, sm.Int(k))
}

func TestStateMapUint8(t *testing.T) {
	var (
		sm         = NewStateMap()
		k          = 1
		v    uint8 = 111
		zero uint8 = 0
	)

	assert.Equal(t, zero, sm.Uint8(k))
	sm.SetUint8(k, v)
	assert.Equal(t, v, sm.Uint8(k))
}

func TestStateMapUint16(t *testing.T) {
	var (
		sm          = NewStateMap()
		k           = 1
		v    uint16 = 111
		zero uint16 = 0
	)

	assert.Equal(t, zero, sm.Uint16(k))
	sm.SetUint16(k, v)
	assert.Equal(t, v, sm.Uint16(k))
}

func TestStateMapBool(t *testing.T) {
	var (
		sm       = NewStateMap()
		k        = 1
		tru bool = true
		fal bool = false
	)

	assert.Equal(t, fal, sm.Bool(k))
	sm.SetBool(k, tru)
	assert.Equal(t, tru, sm.Bool(k))
}
