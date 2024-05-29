package a2

import (
	"github.com/pevans/erc/statemap"
)

// Process executes a single execution of an opcode in the Apple II.
func (c *Computer) Process() error {
	err := c.CPU.Execute()
	if err != nil {
		return err
	}

	/*
		fmt.Printf("ugh: %s\n", c.CPU.NextInstruction())
		if c.CPU.NextInstruction() == "INC $C083,X" {
			c.Debugger = true
		}
	*/

	// Check if this is was a knock-knock on one of our bank switches
	switch c.CPU.EffAddr {
	case 0xC081, 0xC083, 0xC085, 0xC087, 0xC089, 0xC08B, 0xC08D, 0xC08F:
		if c.state.Bool(statemap.InstructionReadOp) {
			c.state.SetInt(statemap.BankReadAttempts, c.state.Int(statemap.BankReadAttempts)+1)
			return nil
		}
	}

	c.state.SetInt(statemap.BankReadAttempts, 0)

	return nil
}
