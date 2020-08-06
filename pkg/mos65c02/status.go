package mos65c02

import (
	"github.com/pevans/erc/pkg/data"
)

// This block defines the flags that we recognize within the status
// register.
const (
	CARRY     = data.Byte(1)
	ZERO      = data.Byte(2)
	INTERRUPT = data.Byte(4)
	DECIMAL   = data.Byte(8)
	BREAK     = data.Byte(16)
	UNUSED    = data.Byte(32)
	OVERFLOW  = data.Byte(64)
	NEGATIVE  = data.Byte(128)
)

// ApplyStatus will make a status update for the given flag based upon
// cond being true or not.
func (c *CPU) ApplyStatus(cond bool, flag data.Byte) {
	c.P &= ^flag
	if cond {
		c.P |= flag
	}
}

// ApplyN will apply the normal negative status check (which is whether
// the eighth bit is high or not).
func (c *CPU) ApplyN(val data.Byte) {
	c.ApplyStatus(val&0x80 > 0, NEGATIVE)
}

// ApplyZ will apply the normal zero status check, which is literally if
// val is zero or not.
func (c *CPU) ApplyZ(val data.Byte) {
	c.ApplyStatus(val == 0, ZERO)
}

// ApplyNZ will apply both the normal negative and zero checks.
func (c *CPU) ApplyNZ(val data.Byte) {
	c.ApplyN(val)
	c.ApplyZ(val)
}

// Compare will compute the difference between the given base and the
// current EffVal value of c. ApplyNZ is called on the result. CARRY is
// set if the result is greater than zero.
func Compare(c *CPU, base data.Byte) {
	res := base - c.EffVal

	c.ApplyNZ(res)
	c.ApplyStatus(res > 0, CARRY)
}
