package a2speaker

import (
	"testing"

	"github.com/pevans/erc/a2/a2audio"
	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/memory"
	"github.com/stretchr/testify/suite"
)

type speakerSuite struct {
	suite.Suite

	state *memory.StateMap
}

type mockSpeaker struct {
	pushCount int
}

func (m *mockSpeaker) Push(cycle uint64, state bool) { m.pushCount++ }
func (m *mockSpeaker) Pop() *a2audio.ToggleEvent     { return nil }
func (m *mockSpeaker) Peek() *a2audio.ToggleEvent    { return nil }
func (m *mockSpeaker) Len() int                      { return 0 }
func (m *mockSpeaker) IsActive() bool                { return false }

type mockComputer struct {
	spk     *mockSpeaker
	counter uint64
}

func (m *mockComputer) CycleCounter() uint64 { return m.counter }
func (m *mockComputer) Speaker() Speaker     { return m.spk }

func (s *speakerSuite) SetupTest() {
	s.state = memory.NewStateMap()
	s.state.SetAny(a2state.Computer, &mockComputer{spk: &mockSpeaker{}})
	UseDefaults(s.state)
}

func TestSpeakerSuite(t *testing.T) {
	suite.Run(t, new(speakerSuite))
}

func (s *speakerSuite) TestSwitchRead() {
	s.state.SetBool(a2state.SpeakerState, false)
	SwitchRead(speakerToggle, s.state)
	s.True(s.state.Bool(a2state.SpeakerState))

	SwitchRead(speakerToggle, s.state)
	s.False(s.state.Bool(a2state.SpeakerState))
}

func (s *speakerSuite) TestSwitchWrite() {
	s.Run("write toggles speaker", func() {
		s.state.SetBool(a2state.SpeakerState, false)
		SwitchWrite(speakerToggle, 0, s.state)
		s.True(s.state.Bool(a2state.SpeakerState))

		SwitchWrite(speakerToggle, 0, s.state)
		s.False(s.state.Bool(a2state.SpeakerState))
	})
}

func (s *speakerSuite) TestUseDefaults() {
	UseDefaults(s.state)
	s.False(s.state.Bool(a2state.SpeakerState))
}
