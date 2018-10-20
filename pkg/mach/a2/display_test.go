package a2

import "github.com/stretchr/testify/assert"

func (s *a2Suite) TestDisplayMode() {
	s.comp.DisplayMode = 123
	assert.Equal(s.T(), 123, displayMode(s.comp))

	s.comp.DisplayMode = 124
	assert.Equal(s.T(), 124, displayMode(s.comp))
}

func (s *a2Suite) TestDisplaySetMode() {
	displaySetMode(s.comp, 111)
	assert.Equal(s.T(), 111, s.comp.DisplayMode)

	displaySetMode(s.comp, 222)
	assert.Equal(s.T(), 222, s.comp.DisplayMode)
}
