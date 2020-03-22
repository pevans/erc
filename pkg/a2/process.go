package a2

import "time"

// Process executes a single execution of an opcode in the Apple II.
func (c *Computer) Process() error {
	return c.CPU.Execute()
}

func (c *Computer) ProcessLoop() {
	for {
		if err := c.Process(); err != nil {
			c.log.Error(err)
			return
		}

		time.Sleep(100 * time.Nanosecond)
	}
}
