package a2display

import (
	"testing"

	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/gfx"
	"github.com/pevans/erc/memory"
	"github.com/stretchr/testify/suite"
)

type displaySuite struct {
	suite.Suite

	state *memory.StateMap
	main  *memory.Segment
	aux   *memory.Segment
}

func (s *displaySuite) SetupTest() {
	// Create memory segments
	s.main = memory.NewSegment(0x10000)
	s.aux = memory.NewSegment(0x10000)

	// Create state map
	s.state = memory.NewStateMap()

	// Set up segment mappings
	s.state.SetSegment(a2state.MemMainSegment, s.main)
	s.state.SetSegment(a2state.MemAuxSegment, s.aux)
	s.state.SetSegment(a2state.MemReadSegment, s.main)
	s.state.SetSegment(a2state.MemWriteSegment, s.main)
	s.state.SetSegment(a2state.DisplayAuxSegment, s.aux)

	// Create and store mock Computer for ComputerState interface
	mockComp := &mockComputer{screen: gfx.NewFrameBuffer(560, 192)}
	s.state.SetAny(a2state.Computer, mockComp)
}

// mockComputer implements ComputerState interface
type mockComputer struct {
	screen *gfx.FrameBuffer
	vblank bool
}

func (m *mockComputer) IsVerticalBlank() bool {
	return m.vblank
}

func (m *mockComputer) GetScreen() *gfx.FrameBuffer {
	return m.screen
}

func TestDisplaySuite(t *testing.T) {
	suite.Run(t, new(displaySuite))
}

func (s *displaySuite) TestDisplaySwitcherSwitchRead() {
	var (
		hi uint8 = 0x80
		lo uint8 = 0x00
	)

	s.Run("high on bit 7", func() {
		test := func(key int, a int) {
			s.state.SetBool(key, true)
			s.Equal(hi, SwitchRead(a, s.state))
			s.state.SetBool(key, false)
			s.Equal(lo, SwitchRead(a, s.state))
		}

		test(a2state.DisplayAltChar, RdAltChar)
		test(a2state.DisplayCol80, Rd80Col)
		test(a2state.DisplayDoubleHigh, RdDHires)
		test(a2state.DisplayHires, RdHires)
		test(a2state.DisplayIou, RdIOUDis)
		test(a2state.DisplayMixed, RdMixed)
		test(a2state.DisplayPage2, RdPage2)
		test(a2state.DisplayStore80, Rd80Store)
		test(a2state.DisplayText, RdText)
	})

	s.Run("reads turn stuff on", func() {
		onfn := func(key int, a int) {
			s.state.SetBool(key, false)
			SwitchRead(a, s.state)
			s.True(s.state.Bool(key))
		}

		onfn(a2state.DisplayPage2, OnPage2)
		onfn(a2state.DisplayText, OnText)
		onfn(a2state.DisplayMixed, OnMixed)
		onfn(a2state.DisplayHires, OnHires)

		// doubleHigh is set regardless of IOUDIS state
		onfn(a2state.DisplayDoubleHigh, OnDHires)
	})

	s.Run("reads turn stuff off", func() {
		offfn := func(key int, a int) {
			s.state.SetBool(key, true)
			SwitchRead(a, s.state)
			s.False(s.state.Bool(key))
		}

		offfn(a2state.DisplayPage2, OffPage2)
		offfn(a2state.DisplayText, OffText)
		offfn(a2state.DisplayMixed, OffMixed)
		offfn(a2state.DisplayHires, OffHires)

		// doubleHigh is cleared regardless of IOUDIS state
		offfn(a2state.DisplayDoubleHigh, OffDHires)
	})
}

func (s *displaySuite) TestDisplaySwitcherSwitchWrite() {
	s.Run("writes turn stuff on", func() {
		on := func(key int, a int) {
			s.state.SetBool(key, false)
			SwitchWrite(a, 0x0, s.state)
			s.True(s.state.Bool(key))
		}

		on(a2state.DisplayPage2, OnPage2)
		on(a2state.DisplayText, OnText)
		on(a2state.DisplayMixed, OnMixed)
		on(a2state.DisplayHires, OnHires)
		on(a2state.DisplayAltChar, OnAltChar)
		on(a2state.DisplayCol80, On80Col)
		on(a2state.DisplayStore80, On80Store)
		on(a2state.DisplayIou, OnIOUDis)

		// doubleHigh is set regardless of IOUDIS state
		on(a2state.DisplayDoubleHigh, OnDHires)
	})

	s.Run("writes turn stuff off", func() {
		off := func(key int, a int) {
			s.state.SetBool(key, true)
			SwitchWrite(a, 0x0, s.state)
			s.False(s.state.Bool(key))
		}

		off(a2state.DisplayPage2, OffPage2)
		off(a2state.DisplayText, OffText)
		off(a2state.DisplayMixed, OffMixed)
		off(a2state.DisplayHires, OffHires)
		off(a2state.DisplayAltChar, OffAltChar)
		off(a2state.DisplayCol80, Off80Col)
		off(a2state.DisplayStore80, Off80Store)
		off(a2state.DisplayIou, OffIOUDis)

		// doubleHigh is cleared regardless of IOUDIS state
		off(a2state.DisplayDoubleHigh, OffDHires)
	})
}

func (s *displaySuite) TestDisplaySegment() {
	var (
		p1addr  = 0x401
		up1addr = int(p1addr)
		p2addr  = 0x2001
		up2addr = int(p2addr)
		other   = 0x301
		uother  = int(other)
		val     = uint8(0x12)
	)

	s.Run("read from main memory", func() {
		s.state.SetBool(a2state.DisplayStore80, false)
		s.state.Segment(a2state.MemWriteSegment).Set(p1addr, val)
		s.state.Segment(a2state.MemWriteSegment).Set(p2addr, val)
		s.state.Segment(a2state.MemWriteSegment).Set(other, val)
		s.Equal(val, Segment(up1addr, s.state, a2state.MemReadSegment).Get(p1addr))
		s.Equal(val, Segment(up2addr, s.state, a2state.MemReadSegment).Get(p2addr))
		s.Equal(val, Segment(uother, s.state, a2state.MemReadSegment).Get(other))
	})

	s.Run("80store uses aux", func() {
		s.state.SetBool(a2state.DisplayStore80, true)
		s.state.Segment(a2state.MemWriteSegment).Set(p1addr, val)
		s.state.Segment(a2state.MemWriteSegment).Set(p2addr, val)
		s.state.Segment(a2state.MemWriteSegment).Set(other, val)

		// References outside of the display pages should be unaffected
		s.Equal(val, Segment(uother, s.state, a2state.MemReadSegment).Get(other))

		// We should be able to show that we use a different memory segment if
		// highRes is on
		s.state.SetBool(a2state.DisplayPage2, false)
		s.Equal(val, Segment(up1addr, s.state, a2state.MemReadSegment).Get(p1addr))
		s.state.SetBool(a2state.DisplayPage2, true)
		s.NotEqual(val, Segment(up1addr, s.state, a2state.MemReadSegment).Get(p1addr))

		// We need both double high resolution _and_ page2 in order to get a
		// different segment in the page 2 address space.
		s.state.SetBool(a2state.DisplayHires, false)
		s.state.SetBool(a2state.DisplayPage2, false)
		s.Equal(val, Segment(up2addr, s.state, a2state.MemReadSegment).Get(p2addr))
		s.state.SetBool(a2state.DisplayHires, true)
		s.Equal(val, Segment(up2addr, s.state, a2state.MemReadSegment).Get(p2addr))
		s.state.SetBool(a2state.DisplayPage2, true)
		s.NotEqual(val, Segment(up2addr, s.state, a2state.MemReadSegment).Get(p2addr))
	})
}

func (s *displaySuite) TestDisplayRead() {
	var (
		addr  = 0x1111
		uaddr = int(addr)
		val   = uint8(0x22)
	)

	Segment(uaddr, s.state, a2state.MemWriteSegment).Set(addr, val)
	s.Equal(val, Read(uaddr, s.state))
}

func (s *displaySuite) TestDisplayWrite() {
	var (
		addr  = 0x1112
		uaddr = int(addr)
		val   = uint8(0x23)
	)

	s.state.SetBool(a2state.DisplayRedraw, false)
	Write(uaddr, val, s.state)
	s.Equal(val, Segment(uaddr, s.state, a2state.MemReadSegment).Get(addr))
	s.True(s.state.Bool(a2state.DisplayRedraw))
}
