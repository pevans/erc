package a2

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type a2Suite struct {
	suite.Suite

	comp *Computer
}

func (s *a2Suite) SetupSuite() {
	s.comp = NewComputer()
}

func (s *a2Suite) SetupTest() {
	_ = s.comp.Boot("")
}

func TestA2Suite(t *testing.T) {
	suite.Run(t, new(a2Suite))
}
