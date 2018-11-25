package a2

import "github.com/pevans/erc/pkg/mach"

const (
	// PCExpROM allows access to expansion ROM. When this is on, memory
	// in the $C800..$CFFF range is mapped to expansion ROM.
	PCExpROM = 0x20

	// PCSlotCxROM tells us to map $C100..$C7FF to the peripheral ROM
	// area of system ROM.
	PCSlotCxROM = 0x40

	// PCSlotC3ROM maps just the $C300 page of memory to peripheral
	// ROM.
	PCSlotC3ROM = 0x80
)

// Here we compute the position of our rom address based on certain
// memory modes. Generally we're looking at addresses in the CX00 range
// and trying to pull that data from where we actually store it in our
// memory segments.
func pcROMAddr(addr mach.DByte, mode int) mach.DByte {
	romAddr := addr - 0xC000

	if (mode&PCExpROM) > 0 && romAddr >= 0x0800 && romAddr < 0x1000 {
		romAddr += 0x4000
	}

	if (mode&PCSlotCxROM) > 0 && romAddr >= 0x0100 && romAddr < 0x0800 {
		romAddr += 0x4000
	}

	if (mode&PCSlotC3ROM) > 0 && romAddr >= 0x0300 && romAddr < 0x0400 {
		romAddr += 0x4000
	}

	return romAddr
}

func pcRead(c *Computer, addr mach.Addressor) mach.Byte {
	dbyte := mach.DByte(addr.Addr())
	return c.ROM.Get(pcROMAddr(dbyte, c.MemMode))
}

func pcWrite(c *Computer, addr mach.Addressor, val mach.Byte) {
	// Do nothing
}

func newPCSwitchCheck() *SwitchCheck {
	return &SwitchCheck{mode: pcMode, setMode: pcSetMode}
}

func pcMode(c *Computer) int {
	return c.PCMode
}

func pcSetMode(c *Computer, mode int) {
	c.PCMode = mode
}
