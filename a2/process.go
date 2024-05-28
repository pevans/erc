package a2

// Process executes a single execution of an opcode in the Apple II.
func (c *Computer) Process() error {
	err := c.CPU.Execute()
	if err != nil {
		return err
	}

	// Check if this is was a knock-knock on one of our bank switches
	switch c.CPU.EffAddr {
	case 0xC081, 0xC083, 0xC085, 0xC087, 0xC089, 0xC08B, 0xC08D, 0xC08F:
		c.state.SetInt(bankReadAttempts, c.state.Int(bankReadAttempts)+1)
	default:
		c.state.SetInt(bankReadAttempts, 0)
	}

	return nil
}
