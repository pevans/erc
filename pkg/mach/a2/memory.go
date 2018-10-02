package a2

import "github.com/pevans/erc/pkg/mach"

const (
	// BankDefault is the default bank-switching scheme: reads in
	// bs-memory go to ROM; writes to RAM are disallowed; bank 1 memory
	// is used.
	BankDefault = 0x00

	// BankRAM indicates that reads are from RAM rather than ROM.
	BankRAM = 0x01

	// BankWrite tells us that we can write to RAM in bs-memory.
	BankWrite = 0x02

	// BankRAM2 tells us to read from bank 2 memory for $D000..$DFFF.
	BankRAM2 = 0x04

	// BankAuxiliary indicates that we should reads and writes in the
	// zero page AND stack page will be done in auxiliary memory rather
	// than main memory. This flag ALSO indicates that reads and/or
	// writes to bs-memory are done in auxiliary memory.
	BankAuxiliary = 0x08

	// MemDefault tells us to read and write only to main memory.
	MemDefault = 0x00

	// MemReadAux will tell us to read the core first 48k memory from
	// auxiliary memory.
	MemReadAux = 0x01

	// MemWriteAux is the switch that tells us write to auxiliary memory
	// in the core 48k memory range.
	MemWriteAux = 0x02

	// Mem80Store is an "enabling" switch for MemPage2 and MemHires
	// below.  If this bit is not on, then those two other bits don't do
	// anything, and all aux memory access is governed by MemWriteAux
	// and MemReadAux above.
	Mem80Store = 0x04

	// MemPage2 allows access to auxiliary memory for the display page,
	// which is $0400..$07FF. This switch only works if Mem80Store is
	// also enabled.
	MemPage2 = 0x08

	// MemHires allows auxiliary memory access for $2000..$3FFF, as long
	// as MemPage2 and Mem80Store are also enabled.
	MemHires = 0x10

	// MemExpROM allows access to expansion ROM. When this is on, memory
	// in the $C800..$CFFF range is mapped to expansion ROM.
	MemExpROM = 0x20

	// MemSlotCxROM tells us to map $C100..$C7FF to the peripheral ROM
	// area of system ROM.
	MemSlotCxROM = 0x40

	// MemSlotC3ROM maps just the $C300 page of memory to peripheral
	// ROM.
	MemSlotC3ROM = 0x80
)

// Get will return the byte at addr, or will execute a read switch if
// one is present at the given address.
func (c *Computer) Get(addr mach.Addressor) mach.Byte {
	if fn, ok := c.RMap[addr.Addr()]; ok {
		return fn(c, addr)
	}

	return c.Main.Get(addr)
}

// Set will set the byte at addr to val, or will execute a write switch
// if one is present at the given address.
func (c *Computer) Set(addr mach.Addressor, val mach.Byte) {
	if fn, ok := c.WMap[addr.Addr()]; ok {
		fn(c, addr, val)
		return
	}

	c.Main.Set(addr, val)
}

// Here we set up all the soft switches that we'll use in the computer,
// which is a lot!
func (c *Computer) defineSoftSwitches() {
	for addr := 0x0; addr < 0x200; addr++ {
		c.RMap[addr] = zeroPageRead
		c.WMap[addr] = zeroPageWrite
	}
}

func zeroPageRead(c *Computer, addr mach.Addressor) mach.Byte {
	seg := c.Main
	if c.BankMode&BankAuxiliary > 0 {
		seg = c.Aux
	}

	return seg.Get(addr)
}

func zeroPageWrite(c *Computer, addr mach.Addressor, val mach.Byte) {
	seg := c.Main
	if c.BankMode&BankAuxiliary > 0 {
		seg = c.Aux
	}

	seg.Set(addr, val)
}