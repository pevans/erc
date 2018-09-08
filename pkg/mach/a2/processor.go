package a2

func (c *Computer) Process() error {
	return c.CPU.Execute()
}
