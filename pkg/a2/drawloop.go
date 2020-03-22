package a2

func (c *Computer) DrawLoop() {
	if err := c.drawer.Draw(c, nil); err != nil {
		c.log.Error(err)
	}
}
