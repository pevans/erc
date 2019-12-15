package a2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEmulator(t *testing.T) {
	emu := NewEmulator(nil)

	assert.NotEqual(t, nil, emu)
	assert.NotEqual(t, nil, emu.Booter)
}
