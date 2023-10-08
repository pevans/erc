package mos65c02

import (
	"github.com/stretchr/testify/assert"
)

func (s *mosSuite) TestStackAddr() {
	cases := []struct {
		s    uint8
		want uint16
	}{
		{0, 0x100},
		{0xFF, 0x1FF},
		{0x82, 0x182},
	}

	for _, cc := range cases {
		s.cpu.S = cc.s
		assert.Equal(s.T(), cc.want, s.cpu.stackAddr())
	}
}

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
		assert.Equal(s.T(), cc.want, s.cpu.Get(s.cpu.stackAddr()))
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
		s.cpu.Set(s.cpu.stackAddr(), cc.want)
		s.cpu.S--

		assert.Equal(s.T(), cc.want, s.cpu.PopStack())
	}
}
