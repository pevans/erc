package a2

import (
	"github.com/pevans/erc/a2/a2state"
)

// Process executes a single execution of an opcode in the Apple II.
func (c *Computer) Process() (int, error) {
	err := c.CPU.Execute()
	if err != nil {
		return c.CPU.OpcodeCycles(c.CPU.Opcode()), err
	}

	// Check if this is was a knock-knock on one of our bank switches
	switch c.CPU.EffAddr {
	case 0xC081, 0xC083, 0xC085, 0xC087, 0xC089, 0xC08B, 0xC08D, 0xC08F:
		if c.State.Bool(a2state.InstructionReadOp) {
			c.State.SetInt(a2state.BankReadAttempts, c.State.Int(a2state.BankReadAttempts)+1)
			return c.CPU.OpcodeCycles(c.CPU.Opcode()), nil
		}
	}

	c.State.SetInt(a2state.BankReadAttempts, 0)

	return c.CPU.OpcodeCycles(c.CPU.Opcode()), nil
}
