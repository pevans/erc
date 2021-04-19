package mos65c02

import (
	"github.com/pevans/erc/pkg/data"
)

type with struct {
	a    data.Byte
	p    data.Byte
	s    data.Byte
	x    data.Byte
	y    data.Byte
	pc   data.DByte
	val  data.Byte
	addr data.DByte
	mode int
}

func (s *mosSuite) op(fn func(*CPU), c with) {
	s.cpu.A = c.a
	s.cpu.P = c.p
	s.cpu.S = c.s
	s.cpu.X = c.x
	s.cpu.Y = c.y
	s.cpu.PC = c.pc
	s.cpu.EffVal = c.val
	s.cpu.EffAddr = c.addr
	s.cpu.AddrMode = c.mode
	fn(s.cpu)
}

// And implements the AND instruction, which performs a bitwise-and on A
// and the effective value and saves the result there.
func (s *mosSuite) TestAnd() {
	var (
		d127 data.Byte = 127
		d63  data.Byte = 63
		d255 data.Byte = 255
		d0   data.Byte = 0
	)

	s.Run("two equal values result in said values", func() {
		s.op(And, with{a: d127, val: d127})
		s.Equal(d127, s.cpu.A)
	})

	s.Run("different values with binary overlaps result with the smaller subset", func() {
		s.op(And, with{a: d127, val: d63})
		s.Equal(d63, s.cpu.A)
	})

	s.Run("result >= 0x80 sets the negative flag", func() {
		s.op(And, with{a: d255, val: d255})
		s.Equal(NEGATIVE, s.cpu.P&NEGATIVE)
	})

	s.Run("result = 0 sets the zero flag", func() {
		s.op(And, with{a: d255, val: d0})
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

	s.Run("result is shifted left by one bit position", func() {
		s.op(Asl, with{a: d0, val: d64, mode: amAcc})
		s.Equal(d128, s.cpu.A)
	})

	s.Run("result >= 0x80 sets the negative flag", func() {
		s.Equal(NEGATIVE, s.cpu.P&NEGATIVE)
	})

	s.Run("result of zero sets zero flag", func() {
		s.op(Asl, with{a: d0, val: d128})
		s.Equal(d0, s.cpu.A)
		s.Equal(ZERO, s.cpu.P&ZERO)
	})

	s.Run("when bit 7 is high, the carry flag should be set", func() {
		s.Equal(CARRY, s.cpu.P&CARRY)
	})

	s.Run("absolute address mode sets the value in the right place", func() {
		s.op(Asl, with{a: d0, val: d64, addr: addr, mode: amAbs})
		s.Equal(d128, s.cpu.Get(addr))
	})
}

func (s *mosSuite) TestBit() {
	var (
		d63  data.Byte = 63
		d127 data.Byte = 127
		d128 data.Byte = 128
	)

	s.Run("sets zero flag when there are no overlapping bits", func() {
		s.op(Bit, with{a: d127, val: d128})
		s.Equal(ZERO, s.cpu.P&ZERO)
	})

	s.Run("sets negative flag when bit 7 is high", func() {
		s.Equal(NEGATIVE, s.cpu.P&NEGATIVE)
	})

	s.Run("sets overflow flag when bit 6 is high", func() {
		s.op(Bit, with{val: d127})
		s.Equal(OVERFLOW, s.cpu.P&OVERFLOW)
	})

	s.Run("does not set overflow if bit 6 is low", func() {
		s.op(Bit, with{val: d63})
		s.NotEqual(OVERFLOW, s.cpu.P&OVERFLOW)
	})

	s.Run("does not set a zero flag when there are overlapping bits", func() {
		s.op(Bit, with{a: d127, val: d127})
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

	s.Run("sets zero flag when a&val has zero result", func() {
		s.op(Bim, with{a: d1, val: d2})
		s.Equal(ZERO, s.cpu.P&ZERO)
	})

	s.Run("does not set zero flag when a&val has nonzero result", func() {
		s.op(Bim, with{a: d1, val: d1})
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

	s.Run("exclusive-ors two values", func() {
		s.op(Eor, with{a: d4, val: d7})
		s.Equal(d3, s.cpu.A)
	})

	s.Run("does not set zero flag when there is a nonzero result", func() {
		s.NotEqual(ZERO, s.cpu.P&ZERO)
	})

	s.Run("sets zero flag when result is zero", func() {
		s.op(Eor, with{a: d4, val: d4})
		s.Equal(d0, s.cpu.A)
		s.Equal(ZERO, s.cpu.P&ZERO)
	})

	s.Run("does not set negative flag when there is a non-negative result", func() {
		s.NotEqual(NEGATIVE, s.cpu.P&NEGATIVE)
	})

	s.Run("sets negative flag when there is a negative result", func() {
		s.op(Eor, with{a: d3, val: d128})
		s.Equal(NEGATIVE, s.cpu.P&NEGATIVE)
	})
}

func (s *mosSuite) TestLsr() {
	var (
		d0 data.Byte = 0
		d1 data.Byte = 1
		d2 data.Byte = 2
	)

	s.Run("shifts value right by one bit position", func() {
		s.op(Lsr, with{a: d0, val: d2, mode: amAcc})
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
		s.op(Lsr, with{a: d0, val: d1})
		s.Equal(d0, s.cpu.A)
		s.Equal(ZERO, s.cpu.P&ZERO)
	})

	s.Run("sets carry flag when result has bit 0 hi", func() {
		s.Equal(CARRY, s.cpu.P&CARRY)
	})

	s.Run("sets value in memory when not in accumulator addr mode", func() {
		addr := data.DByte(1)
		s.op(Lsr, with{a: d0, val: d2, addr: addr, mode: amAbs})
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

	s.Run("does a bitwise or of a and effval", func() {
		s.op(Ora, with{a: d1, val: d2})
		s.Equal(d3, s.cpu.A)
	})

	s.Run("does not set zero flag when result is nonzero", func() {
		s.NotEqual(ZERO, s.cpu.P&ZERO)
	})

	s.Run("does not set negative flag when result is non-negative", func() {
		s.NotEqual(NEGATIVE, s.cpu.P&NEGATIVE)
	})

	s.Run("sets negative flag when result is negative", func() {
		s.op(Ora, with{a: d1, val: d128})
		s.Equal(NEGATIVE, s.cpu.P&NEGATIVE)
	})

	s.Run("sets zero flag when result is zero", func() {
		s.op(Ora, with{a: d0, val: d0})
		s.Equal(ZERO, s.cpu.P&ZERO)
	})
}

func (s *mosSuite) TestRol() {
	var (
		d1   data.Byte = 1
		d2   data.Byte = 2
		d3   data.Byte = 3
		d64  data.Byte = 64
		d128 data.Byte = 128
	)

	s.Run("rotate performs 9-bit rotation", func() {
		s.op(Rol, with{val: d1, mode: amAcc})
		s.Equal(d2, s.cpu.A)
		s.NotEqual(CARRY, s.cpu.P&CARRY)

		s.op(Rol, with{val: d1, p: CARRY, mode: amAcc})
		s.Equal(d3, s.cpu.A)
		s.NotEqual(CARRY, s.cpu.P&CARRY)

		s.op(Rol, with{val: d128, p: CARRY, mode: amAcc})
		s.Equal(d1, s.cpu.A)
		s.Equal(CARRY, s.cpu.P&CARRY)
	})

	s.Run("sets negative flag when result is negative", func() {
		s.op(Rol, with{val: d64, p: CARRY, mode: amAcc})
		s.Equal(NEGATIVE, s.cpu.P&NEGATIVE)
	})

	s.Run("sets zero flag when result is zero", func() {
		s.op(Rol, with{val: d128, mode: amAcc})
		s.Equal(ZERO, s.cpu.P&ZERO)
	})

	s.Run("sets result in memory when not in accumulator mode", func() {
		addr := data.DByte(1)
		s.op(Rol, with{val: d1, addr: addr, mode: amAbs})
		s.Equal(d2, s.cpu.Get(addr))
	})
}

func (s *mosSuite) TestRor() {
	var (
		d1   data.Byte = 1
		d2   data.Byte = 2
		d64  data.Byte = 64
		d128 data.Byte = 128
		d192 data.Byte = 192
	)

	s.Run("rotate performs 9-bit rotation", func() {
		s.op(Ror, with{val: d2, mode: amAcc})
		s.Equal(d1, s.cpu.A)
		s.NotEqual(CARRY, s.cpu.P&CARRY)

		s.op(Ror, with{val: d1, p: CARRY, mode: amAcc})
		s.Equal(d128, s.cpu.A)
		s.Equal(CARRY, s.cpu.P&CARRY)

		s.op(Ror, with{val: d128, p: CARRY, mode: amAcc})
		s.Equal(d192, s.cpu.A)
		s.NotEqual(CARRY, s.cpu.P&CARRY)
	})

	s.Run("sets negative flag when result is negative", func() {
		s.op(Ror, with{val: d64, p: CARRY})
		s.Equal(NEGATIVE, s.cpu.P&NEGATIVE)
	})

	s.Run("sets zero flag when result is zero", func() {
		s.op(Ror, with{val: d1, mode: amAcc})
		s.Equal(ZERO, s.cpu.P&ZERO)
	})

	s.Run("sets result in memory when not in accumulator mode", func() {
		addr := data.DByte(1)
		s.op(Ror, with{val: d2, addr: addr, mode: amAbs})
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

	s.Run("a&val result is stored as zero flag", func() {
		s.op(Trb, with{a: d4, val: d2, addr: addr})
		s.Equal(ZERO, s.cpu.P&ZERO)

		s.op(Trb, with{a: d3, val: d2, addr: addr})
		s.NotEqual(ZERO, s.cpu.P&ZERO)
	})

	s.Run("the result of (a^$FF)&val is saved in addr", func() {
		s.op(Trb, with{a: d4, val: d2, addr: addr})
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

	s.Run("a&val result is stored as zero flag", func() {
		s.op(Tsb, with{a: d4, val: d2, addr: addr})
		s.Equal(ZERO, s.cpu.P&ZERO)

		s.op(Tsb, with{a: d3, val: d2, addr: addr})
		s.NotEqual(ZERO, s.cpu.P&ZERO)
	})

	s.Run("the result of a|val is saved in addr", func() {
		s.op(Tsb, with{a: d4, val: d2, addr: addr})
		s.Equal(d6, s.cpu.Get(addr))
	})
}
