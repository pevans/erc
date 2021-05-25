package a2

// Shutdown will execute whatever is necessary to basically cease operation of
// the computer.
func (c *Computer) Shutdown() error {
	c.CPU.SMap.WriteLog()
	return nil
}
