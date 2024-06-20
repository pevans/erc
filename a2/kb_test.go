package a2

import "github.com/pevans/erc/statemap"

func (s *a2Suite) TestKBDefaults() {
	var zero uint8 = 0

	kbUseDefaults(s.comp)
	s.Equal(zero, s.comp.State.Uint8(statemap.KBLastKey))
	s.Equal(zero, s.comp.State.Uint8(statemap.KBKeyDown))
	s.Equal(zero, s.comp.State.Uint8(statemap.KBStrobe))
}

func (s *a2Suite) TestClearKeys() {
	s.comp.State.SetUint8(statemap.KBKeyDown, 128)
	s.comp.ClearKeys()
	s.Zero(s.comp.State.Uint8(statemap.KBKeyDown))
}

func (s *a2Suite) TestPressKey() {
	s.Run("clears the high bit and saves the low bits", func() {
		s.comp.PressKey(0xff)
		s.Equal(uint8(0x7f), s.comp.State.Uint8(statemap.KBLastKey))
	})

	s.Run("sets the strobe", func() {
		s.comp.PressKey(0)
		s.Equal(uint8(0x80), s.comp.State.Uint8(statemap.KBStrobe))
	})

	s.Run("sets key down", func() {
		s.comp.PressKey(0)
		s.Equal(uint8(0x80), s.comp.State.Uint8(statemap.KBKeyDown))
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
		s.Equal(out, kbSwitchRead(kbDataAndStrobe, s.comp.State))
	})

	s.Run("any key down", func() {
		s.comp.State.SetUint8(statemap.KBStrobe, hi)
		s.Equal(hi, kbSwitchRead(kbAnyKeyDown, s.comp.State))
		s.Zero(s.comp.State.Uint8(statemap.KBStrobe))
	})
}

func (s *a2Suite) TestKBSwitchWrite() {
	var hi uint8 = 0x80

	s.Run("any key down", func() {
		s.comp.State.SetUint8(statemap.KBStrobe, hi)
		kbSwitchWrite(kbAnyKeyDown, 0, s.comp.State)
		s.Zero(s.comp.State.Uint8(statemap.KBStrobe))
	})
}
