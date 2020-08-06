package a2

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
		s.Equal(c.want, memoryMode(s.comp))
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
		s.Equal(c.want, s.comp.MemMode)
	}
}
