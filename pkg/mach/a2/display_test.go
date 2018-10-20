package a2

import (
	"github.com/pevans/erc/pkg/mach"
	"github.com/stretchr/testify/assert"
)

func (s *a2Suite) TestDisplayMode() {
	s.comp.DisplayMode = 123
	assert.Equal(s.T(), 123, displayMode(s.comp))

	s.comp.DisplayMode = 124
	assert.Equal(s.T(), 124, displayMode(s.comp))
}

func (s *a2Suite) TestDisplaySetMode() {
	displaySetMode(s.comp, 111)
	assert.Equal(s.T(), 111, s.comp.DisplayMode)

	displaySetMode(s.comp, 222)
	assert.Equal(s.T(), 222, s.comp.DisplayMode)
}

func (s *a2Suite) TestDisplayAuxSegment() {
	cases := []struct {
		memMode int
		addr    mach.DByte
		want    *mach.Segment
	}{
		{0, 0x444, nil},
		{Mem80Store, 0x444, nil},
		{Mem80Store | MemHires, 0x444, s.comp.Aux},
		{0, 0x2444, nil},
		{Mem80Store, 0x2444, nil},
		{Mem80Store | MemHires, 0x2444, nil},
		{Mem80Store | MemHires | MemPage2, 0x2444, s.comp.Aux},
	}

	for _, c := range cases {
		s.comp.MemMode = c.memMode
		assert.Equal(s.T(), c.want, displayAuxSegment(s.comp, c.addr))
	}
}

func (s *a2Suite) TestDisplayRead() {
	cases := []struct {
		memMode int
		addr    mach.DByte
		seg     *mach.Segment
		want    mach.Byte
	}{
		{0, 0x444, s.comp.Main, 111},
		{Mem80Store | MemHires, 0x445, s.comp.Aux, 111},
		{0, 0x2444, s.comp.Main, 111},
		{Mem80Store | MemHires | MemPage2, 0x2445, s.comp.Aux, 111},
	}

	for _, c := range cases {
		s.comp.MemMode = c.memMode
		c.seg.Set(c.addr, c.want)

		assert.Equal(s.T(), c.want, displayRead(s.comp, c.addr))
	}
}

func (s *a2Suite) TestDisplayWrite() {
	cases := []struct {
		memMode int
		addr    mach.DByte
		seg     *mach.Segment
		want    mach.Byte
	}{
		{0, 0x444, s.comp.Main, 111},
		{Mem80Store | MemHires, 0x445, s.comp.Aux, 111},
		{0, 0x2444, s.comp.Main, 111},
		{Mem80Store | MemHires | MemPage2, 0x2445, s.comp.Aux, 111},
	}

	for _, c := range cases {
		s.comp.MemMode = c.memMode
		displayWrite(s.comp, c.addr, c.want)

		assert.Equal(s.T(), c.want, c.seg.Get(c.addr))
	}
}
