package a2

func (s *a2Suite) TestMemSwitcherUseDefaults() {
	mem := memSwitcher{}
	mem.UseDefaults(s.comp)

	s.Equal(memMain, s.comp.state.Int(memRead))
	s.Equal(memMain, s.comp.state.Int(memWrite))
}

func (s *a2Suite) TestMemSwitcherSwitchRead() {
	var (
		c013 int   = 0xC013
		c014 int   = 0xC014
		hi   uint8 = 0x80
		lo   uint8 = 0x00
		ms   memSwitcher
	)

	s.Run("read profile", func() {
		s.comp.state.SetInt(memRead, memAux)
		s.Equal(hi, ms.SwitchRead(s.comp, c013))

		s.comp.state.SetInt(memRead, memMain)
		s.Equal(lo, ms.SwitchRead(s.comp, c013))
	})

	s.Run("write profile", func() {
		s.comp.state.SetInt(memWrite, memAux)
		s.Equal(hi, ms.SwitchRead(s.comp, c014))

		s.comp.state.SetInt(memWrite, memMain)
		s.Equal(lo, ms.SwitchRead(s.comp, c014))
	})
}

func (s *a2Suite) TestMemSwitcherSwitchWrite() {
	var (
		c002 int = 0xC002
		c003 int = 0xC003
		c004 int = 0xC004
		c005 int = 0xC005
		ms   memSwitcher
	)

	s.Run("set aux works", func() {
		s.comp.state.SetInt(memRead, memMain)
		ms.SwitchWrite(s.comp, c003, 0)
		s.Equal(memAux, s.comp.state.Int(memRead))

		s.comp.state.SetInt(memWrite, memMain)
		ms.SwitchWrite(s.comp, c005, 0)
		s.Equal(memAux, s.comp.state.Int(memWrite))
	})

	s.Run("set main works", func() {
		s.comp.state.SetInt(memRead, memAux)
		ms.SwitchWrite(s.comp, c002, 0)
		s.Equal(memMain, s.comp.state.Int(memRead))

		s.comp.state.SetInt(memWrite, memAux)
		ms.SwitchWrite(s.comp, c004, 0)
		s.Equal(memMain, s.comp.state.Int(memWrite))
	})
}

func (s *a2Suite) TestComputerGet() {
	idx := 0x1
	val := uint8(0x12)

	s.comp.Main.DirectSet(idx, val)
	s.comp.state.SetInt(memRead, memMain)
	s.comp.state.SetSegment(memReadSegment, s.comp.Main)
	s.Equal(val, s.comp.Get(idx))

	s.comp.Aux.DirectSet(idx, val)
	s.comp.state.SetInt(memRead, memAux)
	s.comp.state.SetSegment(memReadSegment, s.comp.Aux)
	s.Equal(val, s.comp.Get(idx))
}

func (s *a2Suite) TestComputerSet() {
	idx := 0x1
	uidx := int(idx)
	val := uint8(0x12)

	// test a normal set
	delete(s.comp.WMap, uidx)
	s.comp.state.SetInt(memWrite, memMain)
	s.comp.Set(idx, val)
	s.Equal(val, s.comp.Main.Mem[idx])

	// test a set from wmap
	var target uint8
	s.comp.WMap[uidx] = func(c *Computer, addr int, val uint8) {
		target = val
	}
	s.comp.Set(idx, val)
	s.Equal(target, val)

	// test a get from aux
	delete(s.comp.WMap, uidx)
	s.comp.state.SetInt(memWrite, memAux)
	s.comp.Set(idx, val)
	s.Equal(val, s.comp.Aux.Mem[idx])
}

func (s *a2Suite) TestReadSegment() {
	s.comp.state.SetInt(memRead, memMain)
	s.comp.state.SetSegment(memReadSegment, s.comp.Main)
	s.Equal(s.comp.Main, ReadSegment(s.comp.state))

	s.comp.state.SetInt(memRead, memAux)
	s.comp.state.SetSegment(memReadSegment, s.comp.Aux)
	s.Equal(s.comp.Aux, ReadSegment(s.comp.state))
}

func (s *a2Suite) TestWriteSegment() {
	s.comp.state.SetInt(memWrite, memMain)
	s.comp.state.SetSegment(memWriteSegment, s.comp.Main)
	s.Equal(s.comp.Main, WriteSegment(s.comp.state))

	s.comp.state.SetInt(memWrite, memAux)
	s.comp.state.SetSegment(memWriteSegment, s.comp.Aux)
	s.Equal(s.comp.Aux, WriteSegment(s.comp.state))
}
