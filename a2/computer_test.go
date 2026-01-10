package a2

// mockAudioStream is a mock implementation of AudioStream for testing.
type mockAudioStream struct {
	lastVolume float32
}

func (m *mockAudioStream) SetVolume(v float32) {
	m.lastVolume = v
}

func (s *a2Suite) TestNewComputer() {
	c := NewComputer(123)

	s.NotNil(c)
	s.NotNil(c.Aux)
	s.NotNil(c.Main)
	s.NotNil(c.ROM)
	s.NotNil(c.smap)
	s.NotNil(c.Drive1)
	s.NotNil(c.Drive2)
	s.Equal(c.SelectedDrive, c.Drive1)
	s.NotNil(c.CPU)
	s.Equal(c.CPU.RMem, c)
	s.Equal(c.CPU.WMem, c)
}

func (s *a2Suite) TestDimensions() {
	w, h := s.comp.Dimensions()
	s.NotZero(w)
	s.NotZero(h)
}

func (s *a2Suite) TestVolumeUp() {
	mock := &mockAudioStream{}
	comp := NewComputer(1)
	comp.SetAudioStream(mock)

	cases := []struct {
		name           string
		initialVolume  int
		initialMuted   bool
		amount         int
		expectedVolume int
		expectedMuted  bool
		expectedFloat  float32
	}{
		{
			name:           "increase from 50% to 60%",
			initialVolume:  50,
			initialMuted:   false,
			amount:         10,
			expectedVolume: 60,
			expectedMuted:  false,
			expectedFloat:  0.6,
		},
		{
			name:           "cap at 100%",
			initialVolume:  95,
			initialMuted:   false,
			amount:         10,
			expectedVolume: 100,
			expectedMuted:  false,
			expectedFloat:  1.0,
		},
		{
			name:           "unmute when increasing",
			initialVolume:  30, // preserved level
			initialMuted:   true,
			amount:         10,
			expectedVolume: 10, // starts from effective volume (0) when muted
			expectedMuted:  false,
			expectedFloat:  0.1,
		},
		{
			name:           "increase from 0%",
			initialVolume:  0,
			initialMuted:   false,
			amount:         10,
			expectedVolume: 10,
			expectedMuted:  false,
			expectedFloat:  0.1,
		},
		{
			name:           "increase from 0% muted (preserves 10)",
			initialVolume:  10, // preserved from before muting
			initialMuted:   true,
			amount:         10,
			expectedVolume: 10, // should go to 10%, not 20%
			expectedMuted:  false,
			expectedFloat:  0.1,
		},
	}

	for _, c := range cases {
		s.Run(c.name, func() {
			comp.volumeLevel = c.initialVolume
			comp.volumeMuted = c.initialMuted

			comp.VolumeUp(c.amount)

			s.Equal(c.expectedVolume, comp.volumeLevel)
			s.Equal(c.expectedMuted, comp.volumeMuted)
			s.InDelta(c.expectedFloat, mock.lastVolume, 0.001)
		})
	}
}

func (s *a2Suite) TestVolumeDown() {
	mock := &mockAudioStream{}
	comp := NewComputer(1)
	comp.SetAudioStream(mock)

	cases := []struct {
		name              string
		initialVolume     int
		initialMuted      bool
		amount            int
		expectedVolume    int
		expectedMuted     bool
		expectedFloat     float32
		volumeLevelKept   bool
		expectedKeptLevel int
	}{
		{
			name:            "decrease from 50% to 40%",
			initialVolume:   50,
			initialMuted:    false,
			amount:          10,
			expectedVolume:  40,
			expectedMuted:   false,
			expectedFloat:   0.4,
			volumeLevelKept: false,
		},
		{
			name:              "floor at 0% and set muted",
			initialVolume:     10,
			initialMuted:      false,
			amount:            10,
			expectedVolume:    10, // volumeLevel preserved at last non-zero
			expectedMuted:     true,
			expectedFloat:     0.0,
			volumeLevelKept:   true,
			expectedKeptLevel: 10,
		},
		{
			name:              "floor at 0% when going below",
			initialVolume:     5,
			initialMuted:      false,
			amount:            10,
			expectedVolume:    5, // volumeLevel preserved at last non-zero
			expectedMuted:     true,
			expectedFloat:     0.0,
			volumeLevelKept:   true,
			expectedKeptLevel: 5,
		},
		{
			name:              "stay muted when decreasing from muted state",
			initialVolume:     50, // preserved level
			initialMuted:      true,
			amount:            10,
			expectedVolume:    50, // stays at preserved level since we hit 0
			expectedMuted:     true,
			expectedFloat:     0.0,
			volumeLevelKept:   true,
			expectedKeptLevel: 50,
		},
	}

	for _, c := range cases {
		s.Run(c.name, func() {
			comp.volumeLevel = c.initialVolume
			comp.volumeMuted = c.initialMuted

			comp.VolumeDown(c.amount)

			if c.volumeLevelKept {
				s.Equal(c.expectedKeptLevel, comp.volumeLevel, "volumeLevel should be preserved")
			} else {
				s.Equal(c.expectedVolume, comp.volumeLevel)
			}
			s.Equal(c.expectedMuted, comp.volumeMuted)
			s.InDelta(c.expectedFloat, mock.lastVolume, 0.001)
		})
	}
}

func (s *a2Suite) TestVolumeToggle() {
	mock := &mockAudioStream{}
	comp := NewComputer(1)
	comp.SetAudioStream(mock)

	cases := []struct {
		name           string
		initialVolume  int
		initialMuted   bool
		expectedVolume int
		expectedMuted  bool
		expectedFloat  float32
	}{
		{
			name:           "mute from unmuted",
			initialVolume:  60,
			initialMuted:   false,
			expectedVolume: 60, // volume level preserved
			expectedMuted:  true,
			expectedFloat:  0.0,
		},
		{
			name:           "unmute from muted",
			initialVolume:  60,
			initialMuted:   true,
			expectedVolume: 60,
			expectedMuted:  false,
			expectedFloat:  0.6,
		},
		{
			name:           "toggle with 100% volume",
			initialVolume:  100,
			initialMuted:   false,
			expectedVolume: 100,
			expectedMuted:  true,
			expectedFloat:  0.0,
		},
	}

	for _, c := range cases {
		s.Run(c.name, func() {
			comp.volumeLevel = c.initialVolume
			comp.volumeMuted = c.initialMuted

			comp.VolumeToggle()

			s.Equal(c.expectedVolume, comp.volumeLevel)
			s.Equal(c.expectedMuted, comp.volumeMuted)
			s.InDelta(c.expectedFloat, mock.lastVolume, 0.001)
		})
	}
}

func (s *a2Suite) TestGetVolume() {
	comp := NewComputer(1)
	comp.SetAudioStream(&mockAudioStream{})

	cases := []struct {
		name          string
		volumeLevel   int
		volumeMuted   bool
		expectedValue int
	}{
		{
			name:          "unmuted returns volume level",
			volumeLevel:   75,
			volumeMuted:   false,
			expectedValue: 75,
		},
		{
			name:          "muted returns 0",
			volumeLevel:   75,
			volumeMuted:   true,
			expectedValue: 0,
		},
		{
			name:          "unmuted at 0%",
			volumeLevel:   0,
			volumeMuted:   false,
			expectedValue: 0,
		},
	}

	for _, c := range cases {
		s.Run(c.name, func() {
			comp.volumeLevel = c.volumeLevel
			comp.volumeMuted = c.volumeMuted

			result := comp.GetVolume()

			s.Equal(c.expectedValue, result)
		})
	}
}

func (s *a2Suite) TestIsMuted() {
	comp := NewComputer(1)
	comp.SetAudioStream(&mockAudioStream{})

	cases := []struct {
		name          string
		volumeMuted   bool
		expectedValue bool
	}{
		{
			name:          "not muted",
			volumeMuted:   false,
			expectedValue: false,
		},
		{
			name:          "muted",
			volumeMuted:   true,
			expectedValue: true,
		},
	}

	for _, c := range cases {
		s.Run(c.name, func() {
			comp.volumeMuted = c.volumeMuted

			result := comp.IsMuted()

			s.Equal(c.expectedValue, result)
		})
	}
}

func (s *a2Suite) TestSetAudioStream() {
	mock := &mockAudioStream{}
	comp := NewComputer(1)

	comp.SetAudioStream(mock)

	s.Equal(50, comp.volumeLevel, "should initialize to 50%")
	s.Equal(mock, comp.audioStream)
}
