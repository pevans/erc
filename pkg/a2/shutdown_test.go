package a2

func (s *a2Suite) TestComputerShutdown() {
	s.NoError(s.comp.Shutdown())
}
