package a2

func (s *a2Suite) TestComputerShutdown() {
	// Don't try to save anything on shutdown
	s.comp.Drive(1).RemoveDisk()

	s.NoError(s.comp.Shutdown())
}
