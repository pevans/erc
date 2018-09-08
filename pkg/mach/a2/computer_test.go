package a2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewComputer(t *testing.T) {
	comp := NewComputer()

	assert.NotEqual(t, nil, comp.Main)
	assert.NotEqual(t, nil, comp.Aux)
	assert.NotEqual(t, nil, comp.ROM)
	assert.NotEqual(t, nil, comp.CPU)
	assert.NotEqual(t, nil, comp.CPU.RSeg)
	assert.NotEqual(t, nil, comp.CPU.WSeg)
}
