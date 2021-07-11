package a2

func diskReadWrite(c *Computer, addr uint16, val *uint8) {
	var (
		nib = addr & 0xF
	)

	switch nib {
	case 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7:
		// Set the drive phase, thus adjusting the track position
		c.SelectedDrive.SwitchPhase(nib)
		c.log.Debugf("track position is now %v", c.SelectedDrive.TrackPos)

	case 0x8:
		// Turn both drives on
		c.log.Debug("both drives online")
		c.Drive1.Online = true
		c.Drive2.Online = true

	case 0x9:
		// Turn only the selected drive on
		c.log.Debug("selected drive online")
		c.SelectedDrive.Online = true

	case 0xA:
		// Set the selected drive to drive 1
		c.log.Debug("switch selected drive to drive1")
		c.SelectedDrive = c.Drive1

	case 0xB:
		// Set the selected drive to drive 2
		c.log.Debug("switch selected drive to drive2")
		c.SelectedDrive = c.Drive2

	case 0xC:
		if c.SelectedDrive.Mode == ReadMode || c.SelectedDrive.WriteProtect {
			*val = c.SelectedDrive.Read()
		} else if c.SelectedDrive.Mode == WriteMode {
			// Write the value currently in the latch
			c.SelectedDrive.Write()
		}

	case 0xD:
		// Set the latch value (for writes) to val
		if c.SelectedDrive.Mode == WriteMode {
			c.SelectedDrive.Latch = *val
		}

	case 0xE:
		// Set the selected drive mode to read
		c.log.Debug("selected drive is now in read mode")
		c.SelectedDrive.Mode = ReadMode

	case 0xF:
		// Set the selected drive mode to write
		c.log.Debug("selected drive is now in write mode")
		c.SelectedDrive.Mode = WriteMode
	}

	if nib%2 == 0 {
		*val = c.SelectedDrive.Latch
	}
}

func diskRead(c *Computer, addr uint16) uint8 {
	// With reads, we pass a byte value for the ReadWrite function to
	// modify.
	val := uint8(0)

	diskReadWrite(c, addr, &val)

	return val
}

func diskWrite(c *Computer, addr uint16, val uint8) {
	// Compared to Read, we pass the val exactly as it comes in.
	diskReadWrite(c, addr, &val)
}
