package a2

import (
	"github.com/pevans/erc/pkg/data"
	"github.com/pevans/erc/pkg/mach/a2/disk"
)

func diskReadWrite(c *Computer, addr data.DByte, val *data.Byte) {
	var (
		dbyte = data.DByte(addr.Addr())
		nib   = dbyte & 0xF
		drive = c.SelectedDrive
	)

	switch nib {
	case 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7:
		// Set the Drive phase
		drive.Phase = disk.Phase(nib)

	case 0x8:
		// Turn both drives on
		c.Drive1.Online = true
		c.Drive2.Online = true

	case 0x9:
		// Turn only the selected drive on
		drive.Online = true

	case 0xA:
		// Set the selected drive to drive 1
		c.SelectedDrive = c.Drive1

	case 0xB:
		// Set the selected drive to drive 2
		c.SelectedDrive = c.Drive2

	case 0xC:
		// read or write
	case 0xD:
		// set latch

	case 0xE:
		// Set the selected drive mode to read
		drive.Mode = disk.ReadMode

	case 0xF:
		// Set the selected drive mode to write
		drive.Mode = disk.WriteMode
	}
}

func diskRead(c *Computer, addr data.Addressor) data.Byte {
	// Since we won't be in a WriteMode situation, we simply pass a
	// dummy value to readWrite.
	val := data.Byte(0)

	diskReadWrite(c, data.DByte(addr.Addr()), &val)

	// A random number could/should be returned here, but for now, we
	// hard-code something
	return 0xFF
}

func diskWrite(c *Computer, addr data.Addressor, val data.Byte) {
	// Compared to Read, we pass the val exactly as it comes in.
	diskReadWrite(c, data.DByte(addr.Addr()), &val)
}
