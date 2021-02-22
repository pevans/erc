package mos65c02

// Lda implements the LDA (load A) instruction, which saves the EffVal
// into the A register.
func Lda(c *CPU) {
	c.ApplyNZ(c.EffVal)
	c.A = c.EffVal
}

// Ldx implements the LDX (load X) instruction, which saves EffVal into
// X.
func Ldx(c *CPU) {
	c.ApplyNZ(c.EffVal)
	c.X = c.EffVal
}

// Ldy implements the LDY (load Y) instruction, which saves EffVal into
// Y.
func Ldy(c *CPU) {
	c.ApplyNZ(c.EffVal)
	c.Y = c.EffVal
}

// Pha implements the PHA (push A) instruction, which pushes the current
// value of the A register onto the stack.
func Pha(c *CPU) {
	c.PushStack(c.A)
}

// Php implements the PHP (push P) instruction, which pushes the current
// value of the P register onto the stack. It has no connection to
// PHP-the-language.
func Php(c *CPU) {
	c.PushStack(c.P)
}

// Phx implements the PHX (push X) instruction, which pushes the current
// value of the X register onto the stack.
func Phx(c *CPU) {
	c.PushStack(c.X)
}

// Phy implements the PHY (push Y) instruction, which pushes the current
// value of the Y register onto the stack.
func Phy(c *CPU) {
	c.PushStack(c.Y)
}

// Pla implements the PLA (pull A) instruction, which pops the top of
// the stack and saves that value in the A register.
func Pla(c *CPU) {
	c.A = c.PopStack()
	c.ApplyNZ(c.A)
}

// Plp implements the PLP (pull P) instruction, which pops the top of
// the stack and saves that value in the P register.
func Plp(c *CPU) {
	c.P = c.PopStack()
}

// Plx implements the PLX (pull X) instruction, which pops the top of
// the stack and saves that value in the X register.
func Plx(c *CPU) {
	c.X = c.PopStack()
	c.ApplyNZ(c.X)
}

// Ply implements the PLY (pull Y) instruction, which pops the top of
// the stack and saves that value in the Y register.
func Ply(c *CPU) {
	c.Y = c.PopStack()
	c.ApplyNZ(c.Y)
}

// Sta implements the STA (store A) instruction, which saves the current
// value of the A register at the effective address.
func Sta(c *CPU) {
	c.Set(c.EffAddr, c.A)
}

// Stx implements the STX (store X) instruction, which saves the current
// value of the X register at the effective address.
func Stx(c *CPU) {
	c.Set(c.EffAddr, c.X)
}

// Sty implements the STY (store Y) instruction, which saves the current
// value of the Y register at the effective address.
func Sty(c *CPU) {
	c.Set(c.EffAddr, c.Y)
}

// Stz implements the STZ (store zero) instruction, which sets the byte
// at the effective address to zero.
func Stz(c *CPU) {
	c.Set(c.EffAddr, 0)
}

// Tax implements the TAX (transfer A to X) instruction, which sets the
// X register equal to the value of the A register.
func Tax(c *CPU) {
	c.ApplyNZ(c.A)
	c.X = c.A
}

// Tay implements the TAY (transfer A to Y) instruction, which sets the
// Y register equal to the value of the A register.
func Tay(c *CPU) {
	c.ApplyNZ(c.A)
	c.Y = c.A
}

// Tsx implements the TSX (transfer S to X) instruction, which sets the
// X register equal to the value of the S register.
func Tsx(c *CPU) {
	c.ApplyNZ(c.S)
	c.X = c.S
}

// Txa implements the TXA (transfer X to A) instruction, which sets the
// A register equal to the value of the X register.
func Txa(c *CPU) {
	c.ApplyNZ(c.X)
	c.A = c.X
}

// Txs implements the TXS (transfer X to S) instruction, which sets the
// S register equal to the value of the X register.
func Txs(c *CPU) {
	c.ApplyNZ(c.X)
	c.S = c.X
}

// Tya implements the TYA (transfer Y to A) instruction, which sets the
// A register equal to the value of the Y register.
func Tya(c *CPU) {
	c.ApplyNZ(c.Y)
	c.A = c.Y
}
