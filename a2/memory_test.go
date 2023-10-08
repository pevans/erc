package a2

import "github.com/pevans/erc/memory"

func (s *a2Suite) TestMemSwitcherUseDefaults() {
	memUseDefaults(s.comp)

	s.Equal(memMain, s.comp.state.Int(memRead))
	s.Equal(memMain, s.comp.state.Int(memWrite))
}

func (s *a2Suite) TestMemSwitcherSwitchRead() {
	var (
		c013 int   = 0xC013
		c014 int   = 0xC014
		hi   uint8 = 0x80
		lo   uint8 = 0x00
	)

	s.Run("read profile", func() {
		s.comp.state.SetInt(memRead, memAux)
		s.Equal(hi, memSwitchRead(c013, s.comp.state))

		s.comp.state.SetInt(memRead, memMain)
		s.Equal(lo, memSwitchRead(c013, s.comp.state))
	})

	s.Run("write profile", func() {
		s.comp.state.SetInt(memWrite, memAux)
		s.Equal(hi, memSwitchRead(c014, s.comp.state))

		s.comp.state.SetInt(memWrite, memMain)
		s.Equal(lo, memSwitchRead(c014, s.comp.state))
	})
}

func (s *a2Suite) TestMemSwitcherSwitchWrite() {
	var (
		c002 int = 0xC002
		c003 int = 0xC003
		c004 int = 0xC004
		c005 int = 0xC005
	)

	s.Run("set aux works", func() {
		s.comp.state.SetInt(memRead, memMain)
		memSwitchWrite(c003, 0, s.comp.state)
		s.Equal(memAux, s.comp.state.Int(memRead))

		s.comp.state.SetInt(memWrite, memMain)
		memSwitchWrite(c005, 0, s.comp.state)
		s.Equal(memAux, s.comp.state.Int(memWrite))
	})

	s.Run("set main works", func() {
		s.comp.state.SetInt(memRead, memAux)
		memSwitchWrite(c002, 0, s.comp.state)
		s.Equal(memMain, s.comp.state.Int(memRead))

		s.comp.state.SetInt(memWrite, memAux)
		memSwitchWrite(c004, 0, s.comp.state)
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
	s.comp.state.SetInt(memWrite, memMain)
	WriteSegment(s.comp.state).DirectSet(idx, val)
	s.Equal(val, s.comp.Main.Mem[idx])

	// test a set from wmap
	var target uint8
	s.comp.smap.SetWrite(uidx, func(addr int, val uint8, stm *memory.StateMap) {
		target = val
	})
	s.comp.Set(idx, val)
	s.Equal(target, val)

	// test a get from aux
	s.comp.state.SetInt(memWrite, memAux)
	WriteSegment(s.comp.state).DirectSet(idx, val)
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
