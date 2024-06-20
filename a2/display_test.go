package a2

import "github.com/pevans/erc/statemap"

func (s *a2Suite) TestDisplaySwitcherUseDefaults() {
	displayUseDefaults(s.comp)

	s.Equal(true, s.comp.State.Bool(statemap.DisplayText))
	s.Equal(false, s.comp.State.Bool(statemap.DisplayAltChar))
	s.Equal(false, s.comp.State.Bool(statemap.DisplayCol80))
	s.Equal(false, s.comp.State.Bool(statemap.DisplayDoubleHigh))
	s.Equal(false, s.comp.State.Bool(statemap.DisplayHires))
	s.Equal(false, s.comp.State.Bool(statemap.DisplayIou))
	s.Equal(false, s.comp.State.Bool(statemap.DisplayMixed))
	s.Equal(false, s.comp.State.Bool(statemap.DisplayPage2))
	s.Equal(false, s.comp.State.Bool(statemap.DisplayStore80))
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

		test(s.comp, statemap.DisplayAltChar, rdAltChar)
		test(s.comp, statemap.DisplayCol80, rd80Col)
		test(s.comp, statemap.DisplayDoubleHigh, rdDHires)
		test(s.comp, statemap.DisplayHires, rdHires)
		test(s.comp, statemap.DisplayIou, rdIOUDis)
		test(s.comp, statemap.DisplayMixed, rdMixed)
		test(s.comp, statemap.DisplayPage2, rdPage2)
		test(s.comp, statemap.DisplayStore80, rd80Store)
		test(s.comp, statemap.DisplayText, rdText)
	})

	s.Run("reads turn stuff on", func() {
		onfn := func(c *Computer, key int, a int) {
			c.State.SetBool(key, false)
			displaySwitchRead(a, s.comp.State)
			s.True(c.State.Bool(key))
		}

		onfn(s.comp, statemap.DisplayPage2, onPage2)
		onfn(s.comp, statemap.DisplayText, onText)
		onfn(s.comp, statemap.DisplayMixed, onMixed)
		onfn(s.comp, statemap.DisplayHires, onHires)

		// doubleHigh will only be set true if iou is true
		s.comp.State.SetBool(statemap.DisplayIou, true)
		onfn(s.comp, statemap.DisplayDoubleHigh, onDHires)

		// But it would be nice to demonstrate the inverse, that we won't set it
		// true
		s.comp.State.SetBool(statemap.DisplayIou, false)
		s.comp.State.SetBool(statemap.DisplayDoubleHigh, false)
		displaySwitchRead(onDHires, s.comp.State)
		s.False(s.comp.State.Bool(statemap.DisplayDoubleHigh))
	})

	s.Run("reads turn stuff off", func() {
		offfn := func(c *Computer, key int, a int) {
			c.State.SetBool(key, true)
			displaySwitchRead(a, s.comp.State)
			s.False(c.State.Bool(key))
		}

		offfn(s.comp, statemap.DisplayPage2, offPage2)
		offfn(s.comp, statemap.DisplayText, offText)
		offfn(s.comp, statemap.DisplayMixed, offMixed)
		offfn(s.comp, statemap.DisplayHires, offHires)

		// Same as for the on-switches, this will only turn off if iou is true
		s.comp.State.SetBool(statemap.DisplayIou, true)
		offfn(s.comp, statemap.DisplayDoubleHigh, offDHires)

		s.comp.State.SetBool(statemap.DisplayIou, false)
		s.comp.State.SetBool(statemap.DisplayDoubleHigh, true)
		displaySwitchRead(offDHires, s.comp.State)
		s.True(s.comp.State.Bool(statemap.DisplayDoubleHigh))
	})
}

func (s *a2Suite) TestDisplaySwitcherSwitchWrite() {
	s.Run("writes turn stuff on", func() {
		on := func(c *Computer, key int, a int) {
			c.State.SetBool(key, false)
			displaySwitchWrite(a, 0x0, s.comp.State)
			s.True(c.State.Bool(key))
		}

		on(s.comp, statemap.DisplayPage2, onPage2)
		on(s.comp, statemap.DisplayText, onText)
		on(s.comp, statemap.DisplayMixed, onMixed)
		on(s.comp, statemap.DisplayHires, onHires)
		on(s.comp, statemap.DisplayAltChar, onAltChar)
		on(s.comp, statemap.DisplayCol80, on80Col)
		on(s.comp, statemap.DisplayStore80, on80Store)
		on(s.comp, statemap.DisplayIou, onIOUDis)

		// doubleHigh will only be set true if iou is true
		s.comp.State.SetBool(statemap.DisplayIou, true)
		on(s.comp, statemap.DisplayDoubleHigh, onDHires)

		// But it would be nice to demonstrate the inverse, that we won't set it
		// true
		s.comp.State.SetBool(statemap.DisplayIou, false)
		s.comp.State.SetBool(statemap.DisplayDoubleHigh, false)
		displaySwitchWrite(onDHires, 0x0, s.comp.State)
		s.False(s.comp.State.Bool(statemap.DisplayDoubleHigh))
	})

	s.Run("writes turn stuff off", func() {
		off := func(c *Computer, key int, a int) {
			c.State.SetBool(key, true)
			displaySwitchWrite(a, 0x0, s.comp.State)
			s.False(c.State.Bool(key))
		}

		off(s.comp, statemap.DisplayPage2, offPage2)
		off(s.comp, statemap.DisplayText, offText)
		off(s.comp, statemap.DisplayMixed, offMixed)
		off(s.comp, statemap.DisplayHires, offHires)
		off(s.comp, statemap.DisplayAltChar, offAltChar)
		off(s.comp, statemap.DisplayCol80, off80Col)
		off(s.comp, statemap.DisplayStore80, off80Store)
		off(s.comp, statemap.DisplayIou, offIOUDis)

		// doubleHigh will only be set true if iou is true
		s.comp.State.SetBool(statemap.DisplayIou, true)
		off(s.comp, statemap.DisplayDoubleHigh, offDHires)

		// But it would be nice to demonstrate the inverse, that we won't set it
		// true
		s.comp.State.SetBool(statemap.DisplayIou, false)
		s.comp.State.SetBool(statemap.DisplayDoubleHigh, true)
		displaySwitchWrite(offDHires, 0x0, s.comp.State)
		s.True(s.comp.State.Bool(statemap.DisplayDoubleHigh))
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
		s.comp.State.SetBool(statemap.DisplayStore80, false)
		WriteSegment(s.comp.State).Set(p1addr, val)
		WriteSegment(s.comp.State).Set(p2addr, val)
		WriteSegment(s.comp.State).Set(other, val)
		s.Equal(val, DisplaySegment(up1addr, s.comp.State, ReadSegment).Get(p1addr))
		s.Equal(val, DisplaySegment(up2addr, s.comp.State, ReadSegment).Get(p2addr))
		s.Equal(val, DisplaySegment(uother, s.comp.State, ReadSegment).Get(other))
	})

	s.Run("80store uses aux", func() {
		s.comp.State.SetBool(statemap.DisplayStore80, true)
		WriteSegment(s.comp.State).Set(p1addr, val)
		WriteSegment(s.comp.State).Set(p2addr, val)
		WriteSegment(s.comp.State).Set(other, val)

		// References outside of the display pages should be unaffected
		s.Equal(val, DisplaySegment(uother, s.comp.State, ReadSegment).Get(other))

		// We should be able to show that we use a different memory segment if
		// highRes is on
		s.comp.State.SetBool(statemap.DisplayPage2, false)
		s.Equal(val, DisplaySegment(up1addr, s.comp.State, ReadSegment).Get(p1addr))
		s.comp.State.SetBool(statemap.DisplayPage2, true)
		s.NotEqual(val, DisplaySegment(up1addr, s.comp.State, ReadSegment).Get(p1addr))

		// We need both double high resolution _and_ page2 in order to get a
		// different segment in the page 2 address space.
		s.comp.State.SetBool(statemap.DisplayHires, false)
		s.comp.State.SetBool(statemap.DisplayPage2, false)
		s.Equal(val, DisplaySegment(up2addr, s.comp.State, ReadSegment).Get(p2addr))
		s.comp.State.SetBool(statemap.DisplayHires, true)
		s.Equal(val, DisplaySegment(up2addr, s.comp.State, ReadSegment).Get(p2addr))
		s.comp.State.SetBool(statemap.DisplayPage2, true)
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

	s.comp.State.SetBool(statemap.DisplayRedraw, false)
	DisplayWrite(uaddr, val, s.comp.State)
	s.Equal(val, DisplaySegment(uaddr, s.comp.State, ReadSegment).Get(addr))
	s.True(s.comp.State.Bool(statemap.DisplayRedraw))
}
