package a2

import (
	"os"
)

func (s *a2Suite) TestComputerLoad() {
	dat, _ := os.Open("../data/logical.disk")
	s.NoError(s.comp.Load(dat, "something.dsk"))

	s.Error(s.comp.Load(nil, "bad"))
}
