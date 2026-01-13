package a2

import "github.com/pevans/erc/a2/a2state"

func (s *a2Suite) TestDisplaySwitcherUseDefaults() {
	displayUseDefaults(s.comp)

	s.Equal(true, s.comp.State.Bool(a2state.DisplayText))
	s.Equal(false, s.comp.State.Bool(a2state.DisplayAltChar))
	s.Equal(false, s.comp.State.Bool(a2state.DisplayCol80))
	s.Equal(false, s.comp.State.Bool(a2state.DisplayDoubleHigh))
	s.Equal(false, s.comp.State.Bool(a2state.DisplayHires))
	s.Equal(false, s.comp.State.Bool(a2state.DisplayIou))
	s.Equal(false, s.comp.State.Bool(a2state.DisplayMixed))
	s.Equal(false, s.comp.State.Bool(a2state.DisplayPage2))
	s.Equal(false, s.comp.State.Bool(a2state.DisplayStore80))
}

func (s *a2Suite) TestDisplaySwitcherSwitchRead() {
	var (
		hi uint8 = 0x80
		lo uint8 = 0x00
	)

	s.Run("high on bit 7", func() {
		test := func(c *Computer, key int, a int) {
			c.State.SetBool(key, true)
			s.Equal(hi, displaySwitchRead(a, s.comp.State))
			c.State.SetBool(key, false)
			s.Equal(lo, displaySwitchRead(a, s.comp.State))
		}

		test(s.comp, a2state.DisplayAltChar, rdAltChar)
		test(s.comp, a2state.DisplayCol80, rd80Col)
		test(s.comp, a2state.DisplayDoubleHigh, rdDHires)
		test(s.comp, a2state.DisplayHires, rdHires)
		test(s.comp, a2state.DisplayIou, rdIOUDis)
		test(s.comp, a2state.DisplayMixed, rdMixed)
		test(s.comp, a2state.DisplayPage2, rdPage2)
		test(s.comp, a2state.DisplayStore80, rd80Store)
		test(s.comp, a2state.DisplayText, rdText)
	})

	s.Run("reads turn stuff on", func() {
		onfn := func(c *Computer, key int, a int) {
			c.State.SetBool(key, false)
			displaySwitchRead(a, s.comp.State)
			s.True(c.State.Bool(key))
		}

		onfn(s.comp, a2state.DisplayPage2, onPage2)
		onfn(s.comp, a2state.DisplayText, onText)
		onfn(s.comp, a2state.DisplayMixed, onMixed)
		onfn(s.comp, a2state.DisplayHires, onHires)

		// doubleHigh will only be set true if iou is true
		s.comp.State.SetBool(a2state.DisplayIou, true)
		onfn(s.comp, a2state.DisplayDoubleHigh, onDHires)

		// But it would be nice to demonstrate the inverse, that we won't set it
		// true
		s.comp.State.SetBool(a2state.DisplayIou, false)
		s.comp.State.SetBool(a2state.DisplayDoubleHigh, false)
		displaySwitchRead(onDHires, s.comp.State)
		s.False(s.comp.State.Bool(a2state.DisplayDoubleHigh))
	})

	s.Run("reads turn stuff off", func() {
		offfn := func(c *Computer, key int, a int) {
			c.State.SetBool(key, true)
			displaySwitchRead(a, s.comp.State)
			s.False(c.State.Bool(key))
		}

		offfn(s.comp, a2state.DisplayPage2, offPage2)
		offfn(s.comp, a2state.DisplayText, offText)
		offfn(s.comp, a2state.DisplayMixed, offMixed)
		offfn(s.comp, a2state.DisplayHires, offHires)

		// Same as for the on-switches, this will only turn off if iou is true
		s.comp.State.SetBool(a2state.DisplayIou, true)
		offfn(s.comp, a2state.DisplayDoubleHigh, offDHires)

		s.comp.State.SetBool(a2state.DisplayIou, false)
		s.comp.State.SetBool(a2state.DisplayDoubleHigh, true)
		displaySwitchRead(offDHires, s.comp.State)
		s.True(s.comp.State.Bool(a2state.DisplayDoubleHigh))
	})
}

func (s *a2Suite) TestDisplaySwitcherSwitchWrite() {
	s.Run("writes turn stuff on", func() {
		on := func(c *Computer, key int, a int) {
			c.State.SetBool(key, false)
			displaySwitchWrite(a, 0x0, s.comp.State)
			s.True(c.State.Bool(key))
		}

		on(s.comp, a2state.DisplayPage2, onPage2)
		on(s.comp, a2state.DisplayText, onText)
		on(s.comp, a2state.DisplayMixed, onMixed)
		on(s.comp, a2state.DisplayHires, onHires)
		on(s.comp, a2state.DisplayAltChar, onAltChar)
		on(s.comp, a2state.DisplayCol80, on80Col)
		on(s.comp, a2state.DisplayStore80, on80Store)
		on(s.comp, a2state.DisplayIou, onIOUDis)

		// doubleHigh will only be set true if iou is true
		s.comp.State.SetBool(a2state.DisplayIou, true)
		on(s.comp, a2state.DisplayDoubleHigh, onDHires)

		// But it would be nice to demonstrate the inverse, that we won't set it
		// true
		s.comp.State.SetBool(a2state.DisplayIou, false)
		s.comp.State.SetBool(a2state.DisplayDoubleHigh, false)
		displaySwitchWrite(onDHires, 0x0, s.comp.State)
		s.False(s.comp.State.Bool(a2state.DisplayDoubleHigh))
	})

	s.Run("writes turn stuff off", func() {
		off := func(c *Computer, key int, a int) {
			c.State.SetBool(key, true)
			displaySwitchWrite(a, 0x0, s.comp.State)
			s.False(c.State.Bool(key))
		}

		off(s.comp, a2state.DisplayPage2, offPage2)
		off(s.comp, a2state.DisplayText, offText)
		off(s.comp, a2state.DisplayMixed, offMixed)
		off(s.comp, a2state.DisplayHires, offHires)
		off(s.comp, a2state.DisplayAltChar, offAltChar)
		off(s.comp, a2state.DisplayCol80, off80Col)
		off(s.comp, a2state.DisplayStore80, off80Store)
		off(s.comp, a2state.DisplayIou, offIOUDis)

		// doubleHigh will only be set true if iou is true
		s.comp.State.SetBool(a2state.DisplayIou, true)
		off(s.comp, a2state.DisplayDoubleHigh, offDHires)

		// But it would be nice to demonstrate the inverse, that we won't set it
		// true
		s.comp.State.SetBool(a2state.DisplayIou, false)
		s.comp.State.SetBool(a2state.DisplayDoubleHigh, true)
		displaySwitchWrite(offDHires, 0x0, s.comp.State)
		s.True(s.comp.State.Bool(a2state.DisplayDoubleHigh))
	})
}

func (s *a2Suite) TestDisplaySegment() {
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
		s.comp.State.SetBool(a2state.DisplayStore80, false)
		WriteSegment(s.comp.State).Set(p1addr, val)
		WriteSegment(s.comp.State).Set(p2addr, val)
		WriteSegment(s.comp.State).Set(other, val)
		s.Equal(val, DisplaySegment(up1addr, s.comp.State, ReadSegment).Get(p1addr))
		s.Equal(val, DisplaySegment(up2addr, s.comp.State, ReadSegment).Get(p2addr))
		s.Equal(val, DisplaySegment(uother, s.comp.State, ReadSegment).Get(other))
	})

	s.Run("80store uses aux", func() {
		s.comp.State.SetBool(a2state.DisplayStore80, true)
		WriteSegment(s.comp.State).Set(p1addr, val)
		WriteSegment(s.comp.State).Set(p2addr, val)
		WriteSegment(s.comp.State).Set(other, val)

		// References outside of the display pages should be unaffected
		s.Equal(val, DisplaySegment(uother, s.comp.State, ReadSegment).Get(other))

		// We should be able to show that we use a different memory segment if
		// highRes is on
		s.comp.State.SetBool(a2state.DisplayPage2, false)
		s.Equal(val, DisplaySegment(up1addr, s.comp.State, ReadSegment).Get(p1addr))
		s.comp.State.SetBool(a2state.DisplayPage2, true)
		s.NotEqual(val, DisplaySegment(up1addr, s.comp.State, ReadSegment).Get(p1addr))

		// We need both double high resolution _and_ page2 in order to get a
		// different segment in the page 2 address space.
		s.comp.State.SetBool(a2state.DisplayHires, false)
		s.comp.State.SetBool(a2state.DisplayPage2, false)
		s.Equal(val, DisplaySegment(up2addr, s.comp.State, ReadSegment).Get(p2addr))
		s.comp.State.SetBool(a2state.DisplayHires, true)
		s.Equal(val, DisplaySegment(up2addr, s.comp.State, ReadSegment).Get(p2addr))
		s.comp.State.SetBool(a2state.DisplayPage2, true)
		s.NotEqual(val, DisplaySegment(up2addr, s.comp.State, ReadSegment).Get(p2addr))
	})
}

func (s *a2Suite) TestDisplayRead() {
	var (
		addr  = 0x1111
		uaddr = int(addr)
		val   = uint8(0x22)
	)

	DisplaySegment(uaddr, s.comp.State, WriteSegment).Set(addr, val)
	s.Equal(val, DisplayRead(uaddr, s.comp.State))
}

func (s *a2Suite) TestDisplayWrite() {
	var (
		addr  = 0x1112
		uaddr = int(addr)
		val   = uint8(0x23)
	)

	s.comp.State.SetBool(a2state.DisplayRedraw, false)
	DisplayWrite(uaddr, val, s.comp.State)
	s.Equal(val, DisplaySegment(uaddr, s.comp.State, ReadSegment).Get(addr))
	s.True(s.comp.State.Bool(a2state.DisplayRedraw))
}

func (s *a2Suite) TestIsVerticalBlank() {
	s.Run("returns consistent result based on cycle position", func() {
		_ = s.comp.Boot()

		cycles := s.comp.CPU.CycleCounter() % scanCycleCount
		result := s.comp.IsVerticalBlank()

		if cycles < 12480 {
			s.False(result, "cycles %d < 12480, should not be in vblank", cycles)
		} else {
			s.True(result, "cycles %d >= 12480, should be in vblank", cycles)
		}
	})

	s.Run("changes state after executing instructions", func() {
		_ = s.comp.Boot()

		seenFalse := false
		seenTrue := false

		for range 9000 {
			s.comp.Main.Set(int(s.comp.CPU.PC), 0xEA)
			_ = s.comp.CPU.Execute()

			cycles := s.comp.CPU.CycleCounter() % scanCycleCount
			result := s.comp.IsVerticalBlank()

			expectedInVBlank := cycles >= 12480
			s.Equal(expectedInVBlank, result,
				"at %d cycles (mod %d), expected vblank=%v but got %v",
				s.comp.CPU.CycleCounter(), cycles, expectedInVBlank, result)

			if result {
				seenTrue = true
			} else {
				seenFalse = true
			}

			if seenTrue && seenFalse {
				return
			}
		}

		s.True(seenTrue && seenFalse, "should have seen both true and false states")
	})
}
