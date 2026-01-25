package a2kb

import (
	"testing"

	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/memory"
	"github.com/stretchr/testify/suite"
)

type kbSuite struct {
	suite.Suite

	state *memory.StateMap
}

func (s *kbSuite) SetupTest() {
	s.state = memory.NewStateMap()
}

func TestKBSuite(t *testing.T) {
	suite.Run(t, new(kbSuite))
}

func (s *kbSuite) TestUseDefaults() {
	var zero uint8 = 0

	UseDefaults(s.state)
	s.Equal(zero, s.state.Uint8(a2state.KBLastKey))
	s.Equal(zero, s.state.Uint8(a2state.KBKeyDown))
	s.Equal(zero, s.state.Uint8(a2state.KBStrobe))
}

func (s *kbSuite) TestSwitchRead() {
	var (
		in  uint8 = 0x55
		hi  uint8 = 0x80
		out       = in | hi
	)

	s.Run("data and strobe", func() {
		s.state.SetUint8(a2state.KBLastKey, in&0x7F)
		s.state.SetUint8(a2state.KBStrobe, 0x80)
		s.Equal(out, SwitchRead(dataAndStrobe, s.state))
	})

	s.Run("any key down", func() {
		s.state.SetUint8(a2state.KBStrobe, hi)
		s.state.SetUint8(a2state.KBKeyDown, hi)
		s.Equal(hi, SwitchRead(anyKeyDown, s.state))
		s.Zero(s.state.Uint8(a2state.KBStrobe))
	})
}

func (s *kbSuite) TestSwitchWrite() {
	var hi uint8 = 0x80

	s.Run("any key down", func() {
		s.state.SetUint8(a2state.KBStrobe, hi)
		SwitchWrite(anyKeyDown, 0, s.state)
		s.Zero(s.state.Uint8(a2state.KBStrobe))
	})
}
