package mos_test

import (
	"github.com/stretchr/testify/assert"
)

func (s *mosSuite) TestPushStack() {
	cases := []struct {
		want uint8
	}{
		{0xFF},
		{0x00},
		{0x12},
		{0x34},
	}

	s.cpu.S = 0xFF

	for _, cc := range cases {
		s.cpu.PushStack(cc.want)

		s.cpu.S++
		assert.Equal(s.T(), cc.want, s.cpu.Get(uint16(0x100)+uint16(s.cpu.S)))
	}
}

func (s *mosSuite) TestPopStack() {
	cases := []struct {
		want uint8
	}{
		{0x00},
		{0x12},
		{0xFF},
		{0x34},
	}

	s.cpu.S = 0xFF

	for _, cc := range cases {
		s.cpu.Set(uint16(0x100)+uint16(s.cpu.S), cc.want)
		s.cpu.S--

		assert.Equal(s.T(), cc.want, s.cpu.PopStack())
	}
}
