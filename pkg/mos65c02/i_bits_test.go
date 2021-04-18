package mos65c02

import (
	"github.com/pevans/erc/pkg/data"
)

// And implements the AND instruction, which performs a bitwise-and on A
// and the effective value and saves the result there.
func (s *mosSuite) TestAnd() {
	var (
		d127 data.Byte = 127
		d63  data.Byte = 63
		d255 data.Byte = 255
		d0   data.Byte = 0
	)

	and := func(a, val data.Byte) {
		s.cpu.A = a
		s.cpu.EffVal = val
		s.cpu.P = 0
		And(s.cpu)
	}

	s.Run("two equal values result in said values", func() {
		and(d127, d127)
		s.Equal(d127, s.cpu.A)
	})

	s.Run("different values with binary overlaps result with the smaller subset", func() {
		and(d127, d63)
		s.Equal(d63, s.cpu.A)
	})

	s.Run("result >= 0x80 sets the negative flag", func() {
		and(d255, d255)
		s.Equal(NEGATIVE, s.cpu.P&NEGATIVE)
	})

	s.Run("result = 0 sets the zero flag", func() {
		and(d255, d0)
		s.Equal(ZERO, s.cpu.P&ZERO)
	})
}

func (s *mosSuite) TestAsl() {
	var (
		d0   data.Byte  = 0
		d64  data.Byte  = 64
		d128 data.Byte  = 128
		addr data.DByte = 1
	)

	aslacc := func(a, val data.Byte) {
		s.cpu.AddrMode = amAcc
		s.cpu.A = a
		s.cpu.EffVal = val
		s.cpu.P = 0
		Asl(s.cpu)
	}

	s.Run("result is shifted left by one bit position", func() {
		aslacc(d0, d64)
		s.Equal(d128, s.cpu.A)
	})

	s.Run("result >= 0x80 sets the negative flag", func() {
		s.Equal(NEGATIVE, s.cpu.P&NEGATIVE)
	})

	s.Run("result of zero sets zero flag", func() {
		aslacc(d0, d128)
		s.Equal(d0, s.cpu.A)
		s.Equal(ZERO, s.cpu.P&ZERO)
	})

	s.Run("when bit 7 is high, the carry flag should be set", func() {
		s.Equal(CARRY, s.cpu.P&CARRY)
	})

	s.cpu.AddrMode = amAbs
	s.Run("absolute address mode sets the value in the right place", func() {
		s.cpu.EffVal = d64
		s.cpu.EffAddr = addr
		s.cpu.Set(addr, d0)
		Asl(s.cpu)

		s.Equal(d128, s.cpu.Get(addr))
	})
}

func (s *mosSuite) TestBit() {
	var (
		d63  data.Byte = 63
		d127 data.Byte = 127
		d128 data.Byte = 128
	)

	bit := func(a, val data.Byte) {
		s.cpu.A = a
		s.cpu.EffVal = val
		s.cpu.P = 0
		Bit(s.cpu)
	}

	s.Run("sets zero flag when there are no overlapping bits", func() {
		bit(d127, d128)
		s.Equal(ZERO, s.cpu.P&ZERO)
	})

	s.Run("sets negative flag when bit 7 is high", func() {
		s.Equal(NEGATIVE, s.cpu.P&NEGATIVE)
	})

	s.Run("sets overflow flag when bit 6 is high", func() {
		bit(0, d127)
		s.Equal(OVERFLOW, s.cpu.P&OVERFLOW)
	})

	s.Run("does not set overflow if bit 6 is low", func() {
		bit(0, d63)
		s.NotEqual(OVERFLOW, s.cpu.P&OVERFLOW)
	})

	s.Run("does not set a zero flag when there are overlapping bits", func() {
		bit(d127, d127)
		s.NotEqual(ZERO, s.cpu.P&ZERO)
	})

	s.Run("does not set the negative flag if bit 7 is low", func() {
		s.NotEqual(NEGATIVE, s.cpu.P&NEGATIVE)
	})
}

func (s *mosSuite) TestBim() {
	var (
		d1 data.Byte = 1
		d2 data.Byte = 2
	)

	bim := func(a, val data.Byte) {
		s.cpu.A = a
		s.cpu.EffVal = val
		s.cpu.P = 0
		Bim(s.cpu)
	}

	s.Run("sets zero flag when a&val has zero result", func() {
		bim(d1, d2)
		s.Equal(ZERO, s.cpu.P&ZERO)
	})

	s.Run("does not set zero flag when a&val has nonzero result", func() {
		bim(d1, d1)
		s.NotEqual(ZERO, s.cpu.P&ZERO)
	})
}

func (s *mosSuite) TestEor() {
	var (
		d0   data.Byte = 0
		d3   data.Byte = 3
		d4   data.Byte = 4
		d7   data.Byte = 7
		d128 data.Byte = 128
	)

	eor := func(a, val data.Byte) {
		s.cpu.A = a
		s.cpu.EffVal = val
		s.cpu.P = 0
		Eor(s.cpu)
	}

	s.Run("exclusive-ors two values", func() {
		eor(d4, d7)
		s.Equal(d3, s.cpu.A)
	})

	s.Run("does not set zero flag when there is a nonzero result", func() {
		s.NotEqual(ZERO, s.cpu.P&ZERO)
	})

	s.Run("sets zero flag when result is zero", func() {
		eor(d4, d4)
		s.Equal(d0, s.cpu.A)
		s.Equal(ZERO, s.cpu.P&ZERO)
	})

	s.Run("does not set negative flag when there is a non-negative result", func() {
		s.NotEqual(NEGATIVE, s.cpu.P&NEGATIVE)
	})

	s.Run("sets negative flag when there is a negative result", func() {
		eor(d3, d128)
		s.Equal(NEGATIVE, s.cpu.P&NEGATIVE)
	})
}

func (s *mosSuite) TestLsr() {
	var (
		d0 data.Byte = 0
		d1 data.Byte = 1
		d2 data.Byte = 2
	)

	lsracc := func(val data.Byte) {
		s.cpu.AddrMode = amAcc
		s.cpu.EffVal = val
		s.cpu.P = 0
		Lsr(s.cpu)
	}

	s.Run("shifts value right by one bit position", func() {
		lsracc(d2)
		s.Equal(d1, s.cpu.A)
	})

	s.Run("does not set negative flag", func() {
		s.NotEqual(NEGATIVE, s.cpu.P&NEGATIVE)
	})

	s.Run("does not set zero flag when result is not zero", func() {
		s.NotEqual(ZERO, s.cpu.P&ZERO)
	})

	s.Run("does not set carry flag when bit 0 is low", func() {
		s.NotEqual(CARRY, s.cpu.P&CARRY)
	})

	s.Run("sets zero flag when result is zero", func() {
		lsracc(d1)
		s.Equal(d0, s.cpu.A)
		s.Equal(ZERO, s.cpu.P&ZERO)
	})

	s.Run("sets carry flag when result has bit 0 hi", func() {
		s.Equal(CARRY, s.cpu.P&CARRY)
	})

	s.cpu.AddrMode = amAbs
	s.Run("sets value in memory when not in accumulator addr mode", func() {
		addr := data.DByte(1)
		s.cpu.EffVal = d2
		s.cpu.EffAddr = addr
		s.cpu.P = 0
		Lsr(s.cpu)

		s.Equal(d1, s.cpu.Get(addr))
	})
}

func (s *mosSuite) TestOra() {
	var (
		d0   data.Byte = 0
		d1   data.Byte = 1
		d2   data.Byte = 2
		d3   data.Byte = 3
		d128 data.Byte = 128
	)

	ora := func(a, val data.Byte) {
		s.cpu.A = a
		s.cpu.EffVal = val
		s.cpu.P = 0
		Ora(s.cpu)
	}

	s.Run("does a bitwise or of a and effval", func() {
		ora(d1, d2)
		s.Equal(d3, s.cpu.A)
	})

	s.Run("does not set zero flag when result is nonzero", func() {
		s.NotEqual(ZERO, s.cpu.P&ZERO)
	})

	s.Run("does not set negative flag when result is non-negative", func() {
		s.NotEqual(NEGATIVE, s.cpu.P&NEGATIVE)
	})

	s.Run("sets negative flag when result is negative", func() {
		ora(d1, d128)
		s.Equal(NEGATIVE, s.cpu.P&NEGATIVE)
	})

	s.Run("sets zero flag when result is zero", func() {
		ora(d0, d0)
		s.Equal(ZERO, s.cpu.P&ZERO)
	})
}

func (s *mosSuite) TestRol() {
	var (
		d1           data.Byte = 1
		d2           data.Byte = 2
		d3           data.Byte = 3
		d64          data.Byte = 64
		d128         data.Byte = 128
		withCarry    data.Byte = CARRY
		withoutCarry data.Byte = 0
	)

	rolacc := func(val, p data.Byte) {
		s.cpu.A = 0
		s.cpu.P = p
		s.cpu.EffVal = val
		s.cpu.AddrMode = amAcc
		Rol(s.cpu)
	}

	s.Run("rotate performs 9-bit rotation", func() {
		rolacc(d1, withoutCarry)
		s.Equal(d2, s.cpu.A)
		s.NotEqual(CARRY, s.cpu.P&CARRY)

		rolacc(d1, withCarry)
		s.Equal(d3, s.cpu.A)
		s.NotEqual(CARRY, s.cpu.P&CARRY)

		rolacc(d128, withCarry)
		s.Equal(d1, s.cpu.A)
		s.Equal(CARRY, s.cpu.P&CARRY)
	})

	s.Run("sets negative flag when result is negative", func() {
		rolacc(d64, withCarry)
		s.Equal(NEGATIVE, s.cpu.P&NEGATIVE)
	})

	s.Run("sets zero flag when result is zero", func() {
		rolacc(d128, withoutCarry)
		s.Equal(ZERO, s.cpu.P&ZERO)
	})

	s.Run("sets result in memory when not in accumulator mode", func() {
		addr := data.DByte(1)
		s.cpu.AddrMode = amAbs
		s.cpu.EffAddr = addr
		s.cpu.EffVal = d1
		s.cpu.P = 0
		Rol(s.cpu)
		s.Equal(d2, s.cpu.Get(addr))
	})
}

func (s *mosSuite) TestRor() {
	var (
		d1 data.Byte = 1
		d2 data.Byte = 2
		//d3           data.Byte = 3
		d64          data.Byte = 64
		d128         data.Byte = 128
		d192         data.Byte = 192
		withCarry    data.Byte = CARRY
		withoutCarry data.Byte = 0
	)

	roracc := func(val, p data.Byte) {
		s.cpu.A = 0
		s.cpu.P = p
		s.cpu.EffVal = val
		s.cpu.AddrMode = amAcc
		Ror(s.cpu)
	}

	s.Run("rotate performs 9-bit rotation", func() {
		roracc(d2, withoutCarry)
		s.Equal(d1, s.cpu.A)
		s.NotEqual(CARRY, s.cpu.P&CARRY)

		roracc(d1, withCarry)
		s.Equal(d128, s.cpu.A)
		s.Equal(CARRY, s.cpu.P&CARRY)

		roracc(d128, withCarry)
		s.Equal(d192, s.cpu.A)
		s.NotEqual(CARRY, s.cpu.P&CARRY)
	})

	s.Run("sets negative flag when result is negative", func() {
		roracc(d64, withCarry)
		s.Equal(NEGATIVE, s.cpu.P&NEGATIVE)
	})

	s.Run("sets zero flag when result is zero", func() {
		roracc(d1, withoutCarry)
		s.Equal(ZERO, s.cpu.P&ZERO)
	})

	s.Run("sets result in memory when not in accumulator mode", func() {
		addr := data.DByte(1)
		s.cpu.AddrMode = amAbs
		s.cpu.EffAddr = addr
		s.cpu.EffVal = d2
		s.cpu.P = 0
		Ror(s.cpu)
		s.Equal(d1, s.cpu.Get(addr))
	})
}

func (s *mosSuite) TestTrb() {
	var (
		d4   data.Byte  = 4
		d3   data.Byte  = 3
		d2   data.Byte  = 2
		addr data.DByte = 1
	)

	trb := func(a, val data.Byte, addr data.DByte) {
		s.cpu.A = a
		s.cpu.EffVal = val
		s.cpu.EffAddr = addr
		s.cpu.P = 0
		Trb(s.cpu)
	}

	s.Run("a&val result is stored as zero flag", func() {
		trb(d4, d2, addr)
		s.Equal(ZERO, s.cpu.P&ZERO)

		trb(d3, d2, addr)
		s.NotEqual(ZERO, s.cpu.P&ZERO)
	})

	s.Run("the result of (a^$FF)&val is saved in addr", func() {
		trb(d4, d2, addr)
		s.Equal(d2, s.cpu.Get(addr))
	})
}

func (s *mosSuite) TestTsb() {
	var (
		d6   data.Byte  = 6
		d4   data.Byte  = 4
		d3   data.Byte  = 3
		d2   data.Byte  = 2
		addr data.DByte = 1
	)

	tsb := func(a, val data.Byte, addr data.DByte) {
		s.cpu.A = a
		s.cpu.EffVal = val
		s.cpu.EffAddr = addr
		s.cpu.P = 0
		Tsb(s.cpu)
	}

	s.Run("a&val result is stored as zero flag", func() {
		tsb(d4, d2, addr)
		s.Equal(ZERO, s.cpu.P&ZERO)

		tsb(d3, d2, addr)
		s.NotEqual(ZERO, s.cpu.P&ZERO)
	})

	s.Run("the result of a|val is saved in addr", func() {
		tsb(d4, d2, addr)
		s.Equal(d6, s.cpu.Get(addr))
	})
}
