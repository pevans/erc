package a2

import (
	"testing"

	"github.com/pevans/erc/pkg/data"
	"github.com/stretchr/testify/assert"
)

func TestKBDefaults(t *testing.T) {
	var ks kbSwitcher

	ks.UseDefaults()
	assert.Zero(t, ks.lastKey)
	assert.Zero(t, ks.strobe)
	assert.Zero(t, ks.keyDown)
}

func (s *a2Suite) TestClearKeys() {
	s.comp.kb.keyDown = 128
	s.comp.ClearKeys()
	s.Zero(s.comp.kb.keyDown)
}

func (s *a2Suite) TestPressKey() {
	s.Run("clears the high bit and saves the low bits", func() {
		s.comp.PressKey(0xff)
		s.Equal(data.Byte(0x7f), s.comp.kb.lastKey)
	})

	s.Run("sets the strobe", func() {
		s.comp.PressKey(0)
		s.Equal(data.Byte(0x80), s.comp.kb.strobe)
	})

	s.Run("sets key down", func() {
		s.comp.PressKey(0)
		s.Equal(data.Byte(0x80), s.comp.kb.keyDown)
	})
}

func (s *a2Suite) TestKBSwitchRead() {
	var (
		in  data.Byte = 0x55
		hi  data.Byte = 0x80
		out data.Byte = in | hi
	)

	s.Run("data and strobe", func() {
		s.comp.PressKey(in)
		s.Equal(out, s.comp.kb.SwitchRead(s.comp, kbDataAndStrobe))
	})

	s.Run("any key down", func() {
		s.comp.kb.strobe = hi
		s.Equal(hi, s.comp.kb.SwitchRead(s.comp, kbAnyKeyDown))
		s.Zero(s.comp.kb.strobe)
	})
}

func (s *a2Suite) TestKBSwitchWrite() {
	var hi data.Byte = 0x80

	s.Run("any key down", func() {
		s.comp.kb.strobe = hi
		s.comp.kb.SwitchWrite(s.comp, kbAnyKeyDown, 0)
		s.Zero(s.comp.kb.strobe)
	})
}
