package a2

import (
	"github.com/pevans/erc/pkg/data"
)

func (s *a2Suite) TestDisplaySwitcherUseDefaults() {
	var ds displaySwitcher

	ds.UseDefaults()

	s.Equal(true, ds.text)
	s.Equal(false, ds.altChar)
	s.Equal(false, ds.col80)
	s.Equal(false, ds.doubleHigh)
	s.Equal(false, ds.highRes)
	s.Equal(false, ds.iou)
	s.Equal(false, ds.mixed)
	s.Equal(false, ds.page2)
	s.Equal(false, ds.store80)
}

func (s *a2Suite) TestDisplaySwitcherSwitchRead() {
	var (
		ds displaySwitcher
		hi data.Byte = 0x80
		lo data.Byte = 0x00
	)

	s.Run("high on bit 7", func() {
		test := func(b *bool, a data.Addressor) {
			*b = true
			s.Equal(hi, ds.SwitchRead(s.comp, a))
			*b = false
			s.Equal(lo, ds.SwitchRead(s.comp, a))
		}

		test(&ds.altChar, rdAltChar)
		test(&ds.col80, rd80Col)
		test(&ds.doubleHigh, rdDHires)
		test(&ds.highRes, rdHires)
		test(&ds.iou, rdIOUDis)
		test(&ds.mixed, rdMixed)
		test(&ds.page2, rdPage2)
		test(&ds.store80, rd80Store)
		test(&ds.text, rdText)
	})

	s.Run("reads turn stuff on", func() {
		onfn := func(b *bool, a data.Addressor) {
			*b = false
			ds.SwitchRead(s.comp, a)
			s.True(*b)
		}

		onfn(&ds.page2, onPage2)
		onfn(&ds.text, onText)
		onfn(&ds.mixed, onMixed)
		onfn(&ds.highRes, onHires)

		// doubleHigh will only be set true if iou is true
		ds.iou = true
		onfn(&ds.doubleHigh, onDHires)

		// But it would be nice to demonstrate the inverse, that we won't set it
		// true
		ds.iou = false
		ds.doubleHigh = false
		ds.SwitchRead(s.comp, onDHires)
		s.False(ds.doubleHigh)
	})

	s.Run("reads turn stuff off", func() {
		offfn := func(b *bool, a data.Addressor) {
			*b = true
			ds.SwitchRead(s.comp, a)
			s.False(*b)
		}

		offfn(&ds.page2, offPage2)
		offfn(&ds.text, offText)
		offfn(&ds.mixed, offMixed)
		offfn(&ds.highRes, offHires)

		// Same as for the on-switches, this will only turn off if iou is true
		ds.iou = true
		offfn(&ds.doubleHigh, offDHires)

		ds.iou = false
		ds.doubleHigh = true
		ds.SwitchRead(s.comp, offDHires)
		s.True(ds.doubleHigh)
	})
}

func (s *a2Suite) TestDisplaySwitcherSwitchWrite() {
	var ds displaySwitcher

	s.Run("writes turn stuff on", func() {
		on := func(b *bool, a data.Addressor) {
			*b = false
			ds.SwitchWrite(s.comp, a, 0x0)
			s.True(*b)
		}

		on(&ds.page2, onPage2)
		on(&ds.text, onText)
		on(&ds.mixed, onMixed)
		on(&ds.highRes, onHires)
		on(&ds.altChar, onAltChar)
		on(&ds.col80, on80Col)
		on(&ds.store80, on80Store)
		on(&ds.iou, onIOUDis)

		// doubleHigh will only be set true if iou is true
		ds.iou = true
		on(&ds.doubleHigh, onDHires)

		// But it would be nice to demonstrate the inverse, that we won't set it
		// true
		ds.iou = false
		ds.doubleHigh = false
		ds.SwitchWrite(s.comp, onDHires, 0x0)
		s.False(ds.doubleHigh)
	})

	s.Run("writes turn stuff off", func() {
		off := func(b *bool, a data.Addressor) {
			*b = true
			ds.SwitchWrite(s.comp, a, 0x0)
			s.False(*b)
		}

		off(&ds.page2, offPage2)
		off(&ds.text, offText)
		off(&ds.mixed, offMixed)
		off(&ds.highRes, offHires)
		off(&ds.altChar, offAltChar)
		off(&ds.col80, off80Col)
		off(&ds.store80, off80Store)
		off(&ds.iou, offIOUDis)

		// doubleHigh will only be set true if iou is true
		ds.iou = true
		off(&ds.doubleHigh, offDHires)

		// But it would be nice to demonstrate the inverse, that we won't set it
		// true
		ds.iou = false
		ds.doubleHigh = true
		ds.SwitchWrite(s.comp, offDHires, 0x0)
		s.True(ds.doubleHigh)
	})
}

func (s *a2Suite) TestDisplaySegment() {
	var (
		p1addr = data.DByte(0x401)
		p2addr = data.DByte(0x2001)
		other  = data.DByte(0x301)
		val    = data.Byte(0x12)
	)

	s.Run("read from main memory", func() {
		s.comp.disp.store80 = false
		s.comp.WriteSegment().Set(p1addr, val)
		s.comp.WriteSegment().Set(p2addr, val)
		s.comp.WriteSegment().Set(other, val)
		s.Equal(val, s.comp.DisplaySegment(p1addr).Get(p1addr))
		s.Equal(val, s.comp.DisplaySegment(p2addr).Get(p2addr))
		s.Equal(val, s.comp.DisplaySegment(other).Get(other))
	})

	s.Run("80store uses aux", func() {
		s.comp.disp.store80 = true
		s.comp.WriteSegment().Set(p1addr, val)
		s.comp.WriteSegment().Set(p2addr, val)
		s.comp.WriteSegment().Set(other, val)

		// References outside of the display pages should be unaffected
		s.Equal(val, s.comp.DisplaySegment(other).Get(other))

		// We should be able to show that we use a different memory segment if
		// highRes is on
		s.comp.disp.highRes = false
		s.Equal(val, s.comp.DisplaySegment(p1addr).Get(p1addr))
		s.comp.disp.highRes = true
		s.NotEqual(val, s.comp.DisplaySegment(p1addr).Get(p1addr))

		// We need both double high resolution _and_ page2 in order to get a
		// different segment in the page 2 address space.
		s.comp.disp.doubleHigh = false
		s.comp.disp.page2 = false
		s.Equal(val, s.comp.DisplaySegment(p2addr).Get(p2addr))
		s.comp.disp.doubleHigh = true
		s.Equal(val, s.comp.DisplaySegment(p2addr).Get(p2addr))
		s.comp.disp.page2 = true
		s.NotEqual(val, s.comp.DisplaySegment(p2addr).Get(p2addr))
	})
}

func (s *a2Suite) TestDisplayRead() {
	var (
		addr = data.DByte(0x1111)
		val  = data.Byte(0x22)
	)

	s.comp.DisplaySegment(addr).Set(addr, val)
	s.Equal(val, DisplayRead(s.comp, addr))
}

func (s *a2Suite) TestDisplayWrite() {
	var (
		addr = data.DByte(0x1112)
		val  = data.Byte(0x23)
	)

	s.comp.reDraw = false
	DisplayWrite(s.comp, addr, val)
	s.Equal(val, s.comp.DisplaySegment(addr).Get(addr))
	s.True(s.comp.reDraw)
}
