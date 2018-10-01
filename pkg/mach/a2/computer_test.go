package a2

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type a2Suite struct {
	suite.Suite

	comp *Computer
}

func (s *a2Suite) SetupSuite() {
	s.comp = NewComputer()
}

func (s *a2Suite) SetupTest() {
	_ = s.comp.Boot()
}

func TestNewComputer(t *testing.T) {
	comp := NewComputer()

	assert.NotEqual(t, nil, comp.Main)
	assert.NotEqual(t, nil, comp.Aux)
	assert.NotEqual(t, nil, comp.ROM)
	assert.NotEqual(t, nil, comp.CPU)
	assert.NotEqual(t, nil, comp.CPU.RMem)
	assert.NotEqual(t, nil, comp.CPU.WMem)
	assert.NotEqual(t, nil, comp.RMap)
	assert.NotEqual(t, nil, comp.WMap)
}
