package a2

import (
	"github.com/pevans/erc/statemap"
)

// Process executes a single execution of an opcode in the Apple II.
func (c *Computer) Process() (int, error) {
	err := c.CPU.Execute()
	if err != nil {
		return c.CPU.Cycles(), err
	}

	// Check if this is was a knock-knock on one of our bank switches
	switch c.CPU.EffAddr {
	case 0xC081, 0xC083, 0xC085, 0xC087, 0xC089, 0xC08B, 0xC08D, 0xC08F:
		if c.State.Bool(statemap.InstructionReadOp) {
			c.State.SetInt(statemap.BankReadAttempts, c.State.Int(statemap.BankReadAttempts)+1)
			return c.CPU.Cycles(), nil
		}
	}

	c.State.SetInt(statemap.BankReadAttempts, 0)

	return c.CPU.Cycles(), nil
}
