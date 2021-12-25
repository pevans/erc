package a2

func (s *a2Suite) TestNewComputer() {
	c := NewComputer()

	s.NotNil(c)
	s.NotNil(c.Aux)
	s.NotNil(c.Main)
	s.NotNil(c.ROM)
	s.NotNil(c.smap)
	s.NotNil(c.Drive1)
	s.NotNil(c.Drive2)
	s.Equal(c.SelectedDrive, c.Drive1)
	s.NotNil(c.CPU)
	s.Equal(c.CPU.Memory, c.Main)
}

func (s *a2Suite) TestDimensions() {
	w, h := s.comp.Dimensions()
	s.Equal(280, w)
	s.Equal(192, h)
}
