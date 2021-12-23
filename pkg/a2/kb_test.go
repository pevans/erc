package a2

func (s *a2Suite) TestKBDefaults() {
	var zero uint8 = 0

	kbUseDefaults(s.comp)
	s.Equal(zero, s.comp.state.Uint8(kbLastKey))
	s.Equal(zero, s.comp.state.Uint8(kbKeyDown))
	s.Equal(zero, s.comp.state.Uint8(kbStrobe))
}

func (s *a2Suite) TestClearKeys() {
	s.comp.state.SetUint8(kbKeyDown, 128)
	s.comp.ClearKeys()
	s.Zero(s.comp.state.Uint8(kbKeyDown))
}

func (s *a2Suite) TestPressKey() {
	s.Run("clears the high bit and saves the low bits", func() {
		s.comp.PressKey(0xff)
		s.Equal(uint8(0x7f), s.comp.state.Uint8(kbLastKey))
	})

	s.Run("sets the strobe", func() {
		s.comp.PressKey(0)
		s.Equal(uint8(0x80), s.comp.state.Uint8(kbStrobe))
	})

	s.Run("sets key down", func() {
		s.comp.PressKey(0)
		s.Equal(uint8(0x80), s.comp.state.Uint8(kbKeyDown))
	})
}

func (s *a2Suite) TestKBSwitchRead() {
	var (
		in  uint8 = 0x55
		hi  uint8 = 0x80
		out uint8 = in | hi
	)

	s.Run("data and strobe", func() {
		s.comp.PressKey(in)
		s.Equal(out, kbSwitchRead(kbDataAndStrobe, s.comp.state))
	})

	s.Run("any key down", func() {
		s.comp.state.SetUint8(kbStrobe, hi)
		s.Equal(hi, kbSwitchRead(kbAnyKeyDown, s.comp.state))
		s.Zero(s.comp.state.Uint8(kbStrobe))
	})
}

func (s *a2Suite) TestKBSwitchWrite() {
	var hi uint8 = 0x80

	s.Run("any key down", func() {
		s.comp.state.SetUint8(kbStrobe, hi)
		kbSwitchWrite(kbAnyKeyDown, 0, s.comp.state)
		s.Zero(s.comp.state.Uint8(kbStrobe))
	})
}
