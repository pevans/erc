package a2

// Process executes a single execution of an opcode in the Apple II.
func (c *Computer) Process() error {
	return c.CPU.Execute()
}
