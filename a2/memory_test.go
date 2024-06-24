package a2

import (
	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/memory"
)

func (s *a2Suite) TestMemSwitcherUseDefaults() {
	memUseDefaults(s.comp)

	s.False(s.comp.State.Bool(a2state.MemReadAux))
	s.False(s.comp.State.Bool(a2state.MemWriteAux))
}

func (s *a2Suite) TestMemSwitcherSwitchRead() {
	var (
		c013 int   = 0xC013
		c014 int   = 0xC014
		hi   uint8 = 0x80
		lo   uint8 = 0x00
	)

	s.Run("read profile", func() {
		s.comp.State.SetBool(a2state.MemReadAux, true)
		s.Equal(hi, memSwitchRead(c013, s.comp.State))

		s.comp.State.SetBool(a2state.MemReadAux, false)
		s.Equal(lo, memSwitchRead(c013, s.comp.State))
	})

	s.Run("write profile", func() {
		s.comp.State.SetBool(a2state.MemWriteAux, true)
		s.Equal(hi, memSwitchRead(c014, s.comp.State))

		s.comp.State.SetBool(a2state.MemWriteAux, false)
		s.Equal(lo, memSwitchRead(c014, s.comp.State))
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
		s.comp.State.SetBool(a2state.MemReadAux, false)
		memSwitchWrite(c003, 0, s.comp.State)
		s.True(s.comp.State.Bool(a2state.MemReadAux))

		s.comp.State.SetBool(a2state.MemWriteAux, false)
		memSwitchWrite(c005, 0, s.comp.State)
		s.True(s.comp.State.Bool(a2state.MemWriteAux))
	})

	s.Run("set main works", func() {
		s.comp.State.SetBool(a2state.MemReadAux, true)
		memSwitchWrite(c002, 0, s.comp.State)
		s.False(s.comp.State.Bool(a2state.MemReadAux))

		s.comp.State.SetBool(a2state.MemWriteAux, true)
		memSwitchWrite(c004, 0, s.comp.State)
		s.False(s.comp.State.Bool(a2state.MemWriteAux))
	})
}

func (s *a2Suite) TestComputerGet() {
	idx := 0x1
	val := uint8(0x12)

	s.comp.Main.DirectSet(idx, val)
	s.comp.State.SetBool(a2state.MemReadAux, false)
	s.comp.State.SetSegment(a2state.MemReadSegment, s.comp.Main)
	s.Equal(val, s.comp.Get(idx))

	s.comp.Aux.DirectSet(idx, val)
	s.comp.State.SetBool(a2state.MemReadAux, true)
	s.comp.State.SetSegment(a2state.MemReadSegment, s.comp.Aux)
	s.Equal(val, s.comp.Get(idx))
}

func (s *a2Suite) TestComputerSet() {
	idx := 0x1
	uidx := int(idx)
	val := uint8(0x12)

	// test a normal set
	s.comp.State.SetBool(a2state.MemWriteAux, false)
	WriteSegment(s.comp.State).DirectSet(idx, val)
	s.Equal(val, s.comp.Main.Mem[idx])

	// test a set from wmap
	var target uint8
	s.comp.smap.SetWrite(uidx, func(addr int, val uint8, stm *memory.StateMap) {
		target = val
	})
	s.comp.Set(idx, val)
	s.Equal(target, val)

	// test a get from aux
	s.comp.State.SetBool(a2state.MemWriteAux, true)
	WriteSegment(s.comp.State).DirectSet(idx, val)
	s.Equal(val, s.comp.Aux.Mem[idx])
}

func (s *a2Suite) TestReadSegment() {
	s.comp.State.SetBool(a2state.MemReadAux, false)
	s.comp.State.SetSegment(a2state.MemReadSegment, s.comp.Main)
	s.Equal(s.comp.Main, ReadSegment(s.comp.State))

	s.comp.State.SetBool(a2state.MemReadAux, true)
	s.comp.State.SetSegment(a2state.MemReadSegment, s.comp.Aux)
	s.Equal(s.comp.Aux, ReadSegment(s.comp.State))
}

func (s *a2Suite) TestWriteSegment() {
	s.comp.State.SetBool(a2state.MemWriteAux, false)
	s.comp.State.SetSegment(a2state.MemWriteSegment, s.comp.Main)
	s.Equal(s.comp.Main, WriteSegment(s.comp.State))

	s.comp.State.SetBool(a2state.MemWriteAux, true)
	s.comp.State.SetSegment(a2state.MemWriteSegment, s.comp.Aux)
	s.Equal(s.comp.Aux, WriteSegment(s.comp.State))
}
