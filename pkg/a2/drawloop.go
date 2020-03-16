package a2

func (c *Computer) DrawLoop() error {
	return c.drawer.Draw(c, nil)
}
