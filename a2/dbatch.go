package a2

import (
	"time"

	"github.com/pevans/erc/elog"
)

// StartDbatch starts a debug batch recording session.
func (c *Computer) StartDbatch() {
	c.dbatchMode = true
	c.dbatchTime = time.Now()
	c.instDiffMap = elog.NewInstructionMap()
}

// StopDbatch stops the debug batch recording and writes the instruction diff to file.
func (c *Computer) StopDbatch() error {
	c.dbatchMode = false
	c.dbatchEnded = time.Now()

	if c.instDiffMap != nil && c.instDiffMapFileName != "" {
		err := c.instDiffMap.WriteToFile(c.instDiffMapFileName)
		c.instDiffMap = nil
		return err
	}

	return nil
}
