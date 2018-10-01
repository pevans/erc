package a2

import (
	"github.com/stretchr/testify/assert"
)

func (s *a2Suite) TestDefineSoftSwitches() {
	var ok bool

	for addr := 0x0; addr < 0x200; addr++ {
		_, ok = s.comp.RMap[addr]
		assert.Equal(s.T(), true, ok)

		_, ok = s.comp.WMap[addr]
		assert.Equal(s.T(), true, ok)
	}
}
