package mos

// Clc implements the CLC (clear carry) instruction, which unsets the
// carry flag.
func Clc(c *CPU) {
	c.P &^= CARRY
}

// Cld implements the CLD (clear decimal) instruction, which unsets the
// decimal flag.
func Cld(c *CPU) {
	c.P &^= DECIMAL
}

// Cli implements the CLI (clear interrupt) instruction, which unsets the
// interrupt flag.
func Cli(c *CPU) {
	c.P &^= INTERRUPT
}

// Clv implements the CLV (clear overflow) instruction, which unsets the
// overflow flag.
func Clv(c *CPU) {
	c.P &^= OVERFLOW
}

// Sec implements the SEC (set carry) instruction, which sets the carry
// flag.
func Sec(c *CPU) {
	c.P |= CARRY
}

// Sed implements the SED (set decimal) instruction, which sets the
// decimal flag.
func Sed(c *CPU) {
	c.P |= DECIMAL
}

// Sei implements the SEI (set interrupt) instruction, which sets the
// interrupt flag.
func Sei(c *CPU) {
	c.P |= INTERRUPT
}
