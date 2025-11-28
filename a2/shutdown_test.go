package a2

func (s *a2Suite) TestComputerShutdown() {
	// Don't try to save anything on shutdown
	s.comp.Drive1.ImageName = ""

	s.NoError(s.comp.Shutdown())
}
