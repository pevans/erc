package a2

import (
	"os"
)

func (s *a2Suite) TestComputerLoad() {
	dat, _ := os.Open("../data/logical.disk")
	s.NoError(s.comp.Load(dat, "something.dsk"))

	// If we run Load again without calling RemoveDisk, it'll try and save
	// "something.dsk"
	s.comp.SelectedDrive.RemoveDisk()

	s.Error(s.comp.Load(nil, "bad"))
}
