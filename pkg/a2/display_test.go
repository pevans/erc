package a2

func (s *a2Suite) TestDisplaySwitcherUseDefaults() {
	var ds displaySwitcher

	ds.UseDefaults(s.comp)

	s.Equal(true, s.comp.state.Bool(displayText))
	s.Equal(false, s.comp.state.Bool(displayAltChar))
	s.Equal(false, s.comp.state.Bool(displayCol80))
	s.Equal(false, s.comp.state.Bool(displayDoubleHigh))
	s.Equal(false, s.comp.state.Bool(displayHires))
	s.Equal(false, s.comp.state.Bool(displayIou))
	s.Equal(false, s.comp.state.Bool(displayMixed))
	s.Equal(false, s.comp.state.Bool(displayPage2))
	s.Equal(false, s.comp.state.Bool(displayStore80))
}

func (s *a2Suite) TestDisplaySwitcherSwitchRead() {
	var (
		ds displaySwitcher
		hi uint8 = 0x80
		lo uint8 = 0x00
	)

	s.Run("high on bit 7", func() {
		test := func(c *Computer, key int, a int) {
			c.state.SetBool(key, true)
			s.Equal(hi, ds.SwitchRead(s.comp, a))
			c.state.SetBool(key, false)
			s.Equal(lo, ds.SwitchRead(s.comp, a))
		}

		test(s.comp, displayAltChar, rdAltChar)
		test(s.comp, displayCol80, rd80Col)
		test(s.comp, displayDoubleHigh, rdDHires)
		test(s.comp, displayHires, rdHires)
		test(s.comp, displayIou, rdIOUDis)
		test(s.comp, displayMixed, rdMixed)
		test(s.comp, displayPage2, rdPage2)
		test(s.comp, displayStore80, rd80Store)
		test(s.comp, displayText, rdText)
	})

	s.Run("reads turn stuff on", func() {
		onfn := func(c *Computer, key int, a int) {
			c.state.SetBool(key, false)
			ds.SwitchRead(s.comp, a)
			s.True(c.state.Bool(key))
		}

		onfn(s.comp, displayPage2, onPage2)
		onfn(s.comp, displayText, onText)
		onfn(s.comp, displayMixed, onMixed)
		onfn(s.comp, displayHires, onHires)

		// doubleHigh will only be set true if iou is true
		s.comp.state.SetBool(displayIou, true)
		onfn(s.comp, displayDoubleHigh, onDHires)

		// But it would be nice to demonstrate the inverse, that we won't set it
		// true
		s.comp.state.SetBool(displayIou, false)
		s.comp.state.SetBool(displayDoubleHigh, false)
		ds.SwitchRead(s.comp, onDHires)
		s.False(s.comp.state.Bool(displayDoubleHigh))
	})

	s.Run("reads turn stuff off", func() {
		offfn := func(c *Computer, key int, a int) {
			c.state.SetBool(key, true)
			ds.SwitchRead(s.comp, a)
			s.False(c.state.Bool(key))
		}

		offfn(s.comp, displayPage2, offPage2)
		offfn(s.comp, displayText, offText)
		offfn(s.comp, displayMixed, offMixed)
		offfn(s.comp, displayHires, offHires)

		// Same as for the on-switches, this will only turn off if iou is true
		s.comp.state.SetBool(displayIou, true)
		offfn(s.comp, displayDoubleHigh, offDHires)

		s.comp.state.SetBool(displayIou, false)
		s.comp.state.SetBool(displayDoubleHigh, true)
		ds.SwitchRead(s.comp, offDHires)
		s.True(s.comp.state.Bool(displayDoubleHigh))
	})
}

func (s *a2Suite) TestDisplaySwitcherSwitchWrite() {
	var ds displaySwitcher

	s.Run("writes turn stuff on", func() {
		on := func(c *Computer, key int, a int) {
			c.state.SetBool(key, false)
			ds.SwitchWrite(s.comp, a, 0x0)
			s.True(c.state.Bool(key))
		}

		on(s.comp, displayPage2, onPage2)
		on(s.comp, displayText, onText)
		on(s.comp, displayMixed, onMixed)
		on(s.comp, displayHires, onHires)
		on(s.comp, displayAltChar, onAltChar)
		on(s.comp, displayCol80, on80Col)
		on(s.comp, displayStore80, on80Store)
		on(s.comp, displayIou, onIOUDis)

		// doubleHigh will only be set true if iou is true
		s.comp.state.SetBool(displayIou, true)
		on(s.comp, displayDoubleHigh, onDHires)

		// But it would be nice to demonstrate the inverse, that we won't set it
		// true
		s.comp.state.SetBool(displayIou, false)
		s.comp.state.SetBool(displayDoubleHigh, false)
		ds.SwitchWrite(s.comp, onDHires, 0x0)
		s.False(s.comp.state.Bool(displayDoubleHigh))
	})

	s.Run("writes turn stuff off", func() {
		off := func(c *Computer, key int, a int) {
			c.state.SetBool(key, true)
			ds.SwitchWrite(s.comp, a, 0x0)
			s.False(c.state.Bool(key))
		}

		off(s.comp, displayPage2, offPage2)
		off(s.comp, displayText, offText)
		off(s.comp, displayMixed, offMixed)
		off(s.comp, displayHires, offHires)
		off(s.comp, displayAltChar, offAltChar)
		off(s.comp, displayCol80, off80Col)
		off(s.comp, displayStore80, off80Store)
		off(s.comp, displayIou, offIOUDis)

		// doubleHigh will only be set true if iou is true
		s.comp.state.SetBool(displayIou, true)
		off(s.comp, displayDoubleHigh, offDHires)

		// But it would be nice to demonstrate the inverse, that we won't set it
		// true
		s.comp.state.SetBool(displayIou, false)
		s.comp.state.SetBool(displayDoubleHigh, true)
		ds.SwitchWrite(s.comp, offDHires, 0x0)
		s.True(s.comp.state.Bool(displayDoubleHigh))
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
		s.comp.state.SetBool(displayStore80, false)
		WriteSegment(s.comp.state).Set(p1addr, val)
		WriteSegment(s.comp.state).Set(p2addr, val)
		WriteSegment(s.comp.state).Set(other, val)
		s.Equal(val, DisplaySegment(up1addr, s.comp.state).Get(p1addr))
		s.Equal(val, DisplaySegment(up2addr, s.comp.state).Get(p2addr))
		s.Equal(val, DisplaySegment(uother, s.comp.state).Get(other))
	})

	s.Run("80store uses aux", func() {
		s.comp.state.SetBool(displayStore80, true)
		WriteSegment(s.comp.state).Set(p1addr, val)
		WriteSegment(s.comp.state).Set(p2addr, val)
		WriteSegment(s.comp.state).Set(other, val)

		// References outside of the display pages should be unaffected
		s.Equal(val, DisplaySegment(uother, s.comp.state).Get(other))

		// We should be able to show that we use a different memory segment if
		// highRes is on
		s.comp.state.SetBool(displayHires, false)
		s.Equal(val, DisplaySegment(up1addr, s.comp.state).Get(p1addr))
		s.comp.state.SetBool(displayHires, true)
		s.NotEqual(val, DisplaySegment(up1addr, s.comp.state).Get(p1addr))

		// We need both double high resolution _and_ page2 in order to get a
		// different segment in the page 2 address space.
		s.comp.state.SetBool(displayDoubleHigh, false)
		s.comp.state.SetBool(displayPage2, false)
		s.Equal(val, DisplaySegment(up2addr, s.comp.state).Get(p2addr))
		s.comp.state.SetBool(displayDoubleHigh, true)
		s.Equal(val, DisplaySegment(up2addr, s.comp.state).Get(p2addr))
		s.comp.state.SetBool(displayPage2, true)
		s.NotEqual(val, DisplaySegment(up2addr, s.comp.state).Get(p2addr))
	})
}

func (s *a2Suite) TestDisplayRead() {
	var (
		addr  = 0x1111
		uaddr = int(addr)
		val   = uint8(0x22)
	)

	DisplaySegment(uaddr, s.comp.state).Set(addr, val)
	s.Equal(val, DisplayRead(uaddr, s.comp.state))
}

func (s *a2Suite) TestDisplayWrite() {
	var (
		addr  = 0x1112
		uaddr = int(addr)
		val   = uint8(0x23)
	)

	s.comp.state.SetBool(displayRedraw, false)
	DisplayWrite(uaddr, val, s.comp.state)
	s.Equal(val, DisplaySegment(uaddr, s.comp.state).Get(addr))
	s.True(s.comp.state.Bool(displayRedraw))
}
