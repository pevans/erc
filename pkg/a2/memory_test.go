package a2

func (s *a2Suite) TestMemSwitcherUseDefaults() {
	mem := memSwitcher{}
	mem.UseDefaults()

	s.Equal(memMain, mem.read)
	s.Equal(memMain, mem.write)
}

/*
func (s *a2Suite) TestMemoryMode() {
	cases := []struct {
		memMode int
		want    int
	}{
		{4, 4},
		{0, 0},
	}

	for _, c := range cases {
		s.comp.MemMode = c.memMode
		s.Equal(c.want, memoryMode(s.comp))
	}
}

func (s *a2Suite) TestMemorySetMode() {
	cases := []struct {
		memMode int
		newMode int
		want    int
	}{
		{4, 3, 3},
		{0, 2, 2},
		{1, 0, 0},
	}

	for _, c := range cases {
		s.comp.MemMode = c.memMode
		memorySetMode(s.comp, c.newMode)
		s.Equal(c.want, s.comp.MemMode)
	}
}

func (s *a2Suite) TestComputerGet() {
	idx := data.DByte(0x1)
	val := data.Byte(0x12)

	// test a normal get
	delete(s.comp.RMap, int(idx))
	s.comp.Main.Mem[idx] = val
	s.Equal(val, s.comp.Get(idx))

	// test a get from rmap
	s.comp.Main.Mem[idx] = data.Byte(0)
	s.comp.RMap[int(idx)] = func(c *Computer, addr data.Addressor) data.Byte {
		return val
	}
	s.Equal(val, s.comp.Get(idx))

	// test a get from aux
	delete(s.comp.RMap, int(idx))
	s.comp.Aux.Mem[idx] = val
	s.comp.MemMode = MemReadAux
	s.Equal(val, s.comp.Get(idx))
}

func (s *a2Suite) TestComputerSet() {
	idx := data.DByte(0x1)
	val := data.Byte(0x12)

	// test a normal set
	delete(s.comp.WMap, int(idx))
	s.comp.MemMode = MemDefault
	s.comp.Set(idx, val)
	s.Equal(val, s.comp.Main.Mem[idx])

	// test a set from wmap
	var target data.Byte
	s.comp.WMap[int(idx)] = func(c *Computer, addr data.Addressor, val data.Byte) {
		target = val
	}
	s.comp.Set(idx, val)
	s.Equal(target, val)

	// test a get from aux
	delete(s.comp.WMap, int(idx))
	s.comp.MemMode = MemWriteAux
	s.comp.Set(idx, val)
	s.Equal(val, s.comp.Aux.Mem[idx])
}

func (s *a2Suite) TestReadSegment() {
	s.comp.MemMode = MemDefault
	s.Equal(s.comp.Main, s.comp.ReadSegment())

	s.comp.MemMode = MemReadAux
	s.Equal(s.comp.Aux, s.comp.ReadSegment())
}

func (s *a2Suite) TestWriteSegment() {
	s.comp.MemMode = MemDefault
	s.Equal(s.comp.Main, s.comp.WriteSegment())

	s.comp.MemMode = MemWriteAux
	s.Equal(s.comp.Aux, s.comp.WriteSegment())
}
*/
