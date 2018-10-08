package a2

// Here we set up all the soft switches that we'll use in the computer,
// which is a lot!
func (c *Computer) defineSoftSwitches() {
	c.MapRange(0x0, 0x200, zeroPageRead, zeroPageWrite)

	c.MapRange(0xD000, 0x10000, bankRead, bankWrite)

	c.RMap[0xC013] = c.memorySwitchIsSetR(MemReadAux)
	c.RMap[0xC014] = c.memorySwitchIsSetR(MemWriteAux)
	c.RMap[0xC018] = c.memorySwitchIsSetR(Mem80Store)
	c.RMap[0xC01C] = c.memorySwitchIsSetR(MemPage2)
	c.RMap[0xC01D] = c.memorySwitchIsSetR(MemHires)
	c.RMap[0xC054] = c.memorySwitchUnsetR(MemPage2)
	c.RMap[0xC055] = c.memorySwitchSetR(MemPage2)
	c.RMap[0xC056] = c.memorySwitchUnsetR(MemHires)
	c.RMap[0xC057] = c.memorySwitchSetR(MemHires)

	c.WMap[0xC000] = c.memorySwitchUnsetW(Mem80Store)
	c.WMap[0xC001] = c.memorySwitchSetW(Mem80Store)
	c.WMap[0xC002] = c.memorySwitchUnsetW(MemReadAux)
	c.WMap[0xC003] = c.memorySwitchSetW(MemReadAux)
	c.WMap[0xC004] = c.memorySwitchUnsetW(MemWriteAux)
	c.WMap[0xC005] = c.memorySwitchSetW(MemWriteAux)
	c.WMap[0xC054] = c.memorySwitchUnsetW(MemPage2)
	c.WMap[0xC055] = c.memorySwitchSetW(MemPage2)
	c.WMap[0xC056] = c.memorySwitchUnsetW(MemHires)
	c.WMap[0xC057] = c.memorySwitchSetW(MemHires)

	c.RMap[0xC080] = c.bankSwitchSetR(BankRAM | BankRAM2)
	c.RMap[0xC081] = c.bankSwitchSetR(BankWrite | BankRAM2)
	c.RMap[0xC082] = c.bankSwitchSetR(BankRAM2)
	c.RMap[0xC083] = c.bankSwitchSetR(BankRAM | BankWrite | BankRAM2)
	c.RMap[0xC088] = c.bankSwitchSetR(BankRAM)
	c.RMap[0xC089] = c.bankSwitchSetR(BankWrite)
	c.RMap[0xC08A] = c.bankSwitchSetR(BankDefault)
	c.RMap[0xC08B] = c.bankSwitchSetR(BankRAM | BankWrite)
	c.RMap[0xC011] = c.bankSwitchIsSetR(BankRAM2)
	c.RMap[0xC012] = c.bankSwitchIsSetR(BankRAM)
	c.RMap[0xC016] = c.bankSwitchIsSetR(BankAuxiliary)

	c.WMap[0xC008] = c.bankSwitchUnsetW(BankAuxiliary)
	c.WMap[0xC009] = c.bankSwitchSetW(BankAuxiliary)
}
