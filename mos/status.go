package mos

// This block defines the flags that we recognize within the status
// register.
const (
	CARRY     = uint8(1)
	ZERO      = uint8(2)
	INTERRUPT = uint8(4)
	DECIMAL   = uint8(8)
	BREAK     = uint8(16)
	UNUSED    = uint8(32)
	OVERFLOW  = uint8(64)
	NEGATIVE  = uint8(128)
)

// ApplyStatus will make a status update for the given flag based upon
// cond being true or not.
func (c *CPU) ApplyStatus(cond bool, flag uint8) {
	c.P &= ^flag
	if cond {
		c.P |= flag
	}
}

// ApplyN will apply the normal negative status check (which is whether
// the eighth bit is high or not).
func (c *CPU) ApplyN(val uint8) {
	c.ApplyStatus(val&0x80 > 0, NEGATIVE)
}

// ApplyZ will apply the normal zero status check, which is literally if
// val is zero or not.
func (c *CPU) ApplyZ(val uint8) {
	c.ApplyStatus(val == 0, ZERO)
}

// ApplyNZ will apply both the normal negative and zero checks.
func (c *CPU) ApplyNZ(val uint8) {
	c.ApplyN(val)
	c.ApplyZ(val)
}

// Compare will compute the difference between the given base and the
// current EffVal value of c. ApplyNZ is called on the result. CARRY is
// set if the result is greater than zero.
func Compare(c *CPU, base uint8) {
	res := base - c.EffVal

	c.ApplyNZ(res)
	c.ApplyStatus(base >= c.EffVal, CARRY)
}
