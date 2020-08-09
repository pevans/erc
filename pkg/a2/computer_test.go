package a2

import (
	"os"

	"github.com/pevans/erc/pkg/boot"
)

func (s *a2Suite) TestNewComputer() {
	c := NewComputer()

	s.NotNil(c)
	s.NotNil(c.Aux)
	s.NotNil(c.Main)
	s.NotNil(c.ROM)
	s.NotNil(c.Drive1)
	s.NotNil(c.Drive2)
	s.Equal(c.SelectedDrive, c.Drive1)
	s.NotNil(c.CPU)
	s.Equal(c.CPU.WMem, c)
	s.Equal(c.CPU.RMem, c)
	s.NotNil(c.RMap)
	s.NotNil(c.WMap)
}

func (s *a2Suite) TestSetLogger() {
	l, _ := boot.DefaultConfig().NewLogger()
	s.comp.SetLogger(l)
	s.NotNil(s.comp.log)
}

func (s *a2Suite) TestRecorderWriter() {
	s.comp.SetRecorderWriter(os.Stdout)

	s.Equal(os.Stdout, s.comp.recWriter)
	s.NotNil(s.comp.rec)
}

func (s *a2Suite) TestDimensions() {
	w, h := s.comp.Dimensions()
	s.Equal(280, w)
	s.Equal(192, h)
}
