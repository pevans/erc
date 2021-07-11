package a2

func (s *a2Suite) TestMemSwitcherUseDefaults() {
	mem := memSwitcher{}
	mem.UseDefaults()

	s.Equal(memMain, mem.read)
	s.Equal(memMain, mem.write)
}

func (s *a2Suite) TestMemSwitcherSwitchRead() {
	var (
		c013 uint16 = 0xC013
		c014 uint16 = 0xC014
		hi   uint8  = 0x80
		lo   uint8  = 0x00
		ms   memSwitcher
	)

	s.Run("read profile", func() {
		ms.read = memAux
		s.Equal(hi, ms.SwitchRead(s.comp, c013))

		ms.read = memMain
		s.Equal(lo, ms.SwitchRead(s.comp, c013))
	})

	s.Run("write profile", func() {
		ms.write = memAux
		s.Equal(hi, ms.SwitchRead(s.comp, c014))

		ms.write = memMain
		s.Equal(lo, ms.SwitchRead(s.comp, c014))
	})
}

func (s *a2Suite) TestMemSwitcherSwitchWrite() {
	var (
		c002 uint16 = 0xC002
		c003 uint16 = 0xC003
		c004 uint16 = 0xC004
		c005 uint16 = 0xC005
		ms   memSwitcher
	)

	s.Run("set aux works", func() {
		ms.read = memMain
		ms.SwitchWrite(s.comp, c003, 0)
		s.Equal(memAux, ms.read)

		ms.write = memMain
		ms.SwitchWrite(s.comp, c005, 0)
		s.Equal(memAux, ms.write)
	})

	s.Run("set main works", func() {
		ms.read = memAux
		ms.SwitchWrite(s.comp, c002, 0)
		s.Equal(memMain, ms.read)

		ms.write = memAux
		ms.SwitchWrite(s.comp, c004, 0)
		s.Equal(memMain, ms.write)
	})
}

func (s *a2Suite) TestComputerGet() {
	idx := 0x1
	uidx := uint16(idx)
	val := uint8(0x12)

	// test a normal get
	delete(s.comp.RMap, uidx)
	s.comp.Main.Mem[idx] = val
	s.comp.mem.read = memMain
	s.Equal(val, s.comp.Get(idx))

	// test a get from rmap
	s.comp.Main.Mem[idx] = uint8(0)
	s.comp.RMap[uidx] = func(c *Computer, addr uint16) uint8 {
		return val
	}
	s.Equal(val, s.comp.Get(idx))

	// test a get from aux
	delete(s.comp.RMap, uidx)
	s.comp.Aux.Mem[idx] = val
	s.comp.mem.read = memAux
	s.Equal(val, s.comp.Get(idx))
}

func (s *a2Suite) TestComputerSet() {
	idx := 0x1
	uidx := uint16(idx)
	val := uint8(0x12)

	// test a normal set
	delete(s.comp.WMap, uidx)
	s.comp.mem.write = memMain
	s.comp.Set(idx, val)
	s.Equal(val, s.comp.Main.Mem[idx])

	// test a set from wmap
	var target uint8
	s.comp.WMap[uidx] = func(c *Computer, addr uint16, val uint8) {
		target = val
	}
	s.comp.Set(idx, val)
	s.Equal(target, val)

	// test a get from aux
	delete(s.comp.WMap, uidx)
	s.comp.mem.write = memAux
	s.comp.Set(idx, val)
	s.Equal(val, s.comp.Aux.Mem[idx])
}

func (s *a2Suite) TestReadSegment() {
	s.comp.mem.read = memMain
	s.Equal(s.comp.Main, s.comp.ReadSegment())

	s.comp.mem.read = memAux
	s.Equal(s.comp.Aux, s.comp.ReadSegment())
}

func (s *a2Suite) TestWriteSegment() {
	s.comp.mem.write = memMain
	s.Equal(s.comp.Main, s.comp.WriteSegment())

	s.comp.mem.write = memAux
	s.Equal(s.comp.Aux, s.comp.WriteSegment())
}
