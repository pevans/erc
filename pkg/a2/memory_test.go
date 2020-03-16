package a2

import "github.com/stretchr/testify/assert"

func (s *a2Suite) TestMemoryMode() {
	cases := []struct {
		memMode int
		want    int
	}{
		{4, 4},
		{0, 0},
	}

	for _, c := range cases {
		s.comp.MemMode = c.memMode
		assert.Equal(s.T(), c.want, memoryMode(s.comp))
	}
}

func (s *a2Suite) TestMemorySetMode() {
	cases := []struct {
		memMode int
		newMode int
		want    int
	}{
		{4, 3, 3},
		{0, 2, 2},
		{1, 0, 0},
	}

	for _, c := range cases {
		s.comp.MemMode = c.memMode
		memorySetMode(s.comp, c.newMode)
		assert.Equal(s.T(), c.want, s.comp.MemMode)
	}
}
