package a2

// Here we set up all the soft switches that we'll use in the computer,
// which is a lot!
func (c *Computer) defineSoftSwitches() {
	c.MapRange(0x0, 0x200, zeroPageRead, zeroPageWrite)

	c.MapRange(0xD000, 0x10000, bankRead, bankWrite)

	msc := newMemorySwitchCheck()
	c.RMap[0xC013] = msc.IsSetter(MemReadAux)
	c.RMap[0xC014] = msc.IsSetter(MemWriteAux)
	c.RMap[0xC018] = msc.IsSetter(Mem80Store)
	c.RMap[0xC01C] = msc.IsSetter(MemPage2)
	c.RMap[0xC01D] = msc.IsSetter(MemHires)
	c.RMap[0xC054] = msc.UnSetterR(MemPage2)
	c.RMap[0xC055] = msc.ReSetterR(MemPage2)
	c.RMap[0xC056] = msc.UnSetterR(MemHires)
	c.RMap[0xC057] = msc.ReSetterR(MemHires)
	c.WMap[0xC000] = msc.UnSetterW(Mem80Store)
	c.WMap[0xC001] = msc.ReSetterW(Mem80Store)
	c.WMap[0xC002] = msc.UnSetterW(MemReadAux)
	c.WMap[0xC003] = msc.ReSetterW(MemReadAux)
	c.WMap[0xC004] = msc.UnSetterW(MemWriteAux)
	c.WMap[0xC005] = msc.ReSetterW(MemWriteAux)
	c.WMap[0xC054] = msc.UnSetterW(MemPage2)
	c.WMap[0xC055] = msc.ReSetterW(MemPage2)
	c.WMap[0xC056] = msc.UnSetterW(MemHires)
	c.WMap[0xC057] = msc.ReSetterW(MemHires)

	bsc := newBankSwitchCheck()
	c.RMap[0xC080] = bsc.SetterR(BankRAM | BankRAM2)
	c.RMap[0xC081] = bsc.SetterR(BankWrite | BankRAM2)
	c.RMap[0xC082] = bsc.SetterR(BankRAM2)
	c.RMap[0xC083] = bsc.SetterR(BankRAM | BankWrite | BankRAM2)
	c.RMap[0xC088] = bsc.SetterR(BankRAM)
	c.RMap[0xC089] = bsc.SetterR(BankWrite)
	c.RMap[0xC08A] = bsc.SetterR(BankDefault)
	c.RMap[0xC08B] = bsc.SetterR(BankRAM | BankWrite)
	c.RMap[0xC011] = bsc.IsSetter(BankRAM2)
	c.RMap[0xC012] = bsc.IsSetter(BankRAM)
	c.RMap[0xC016] = bsc.IsSetter(BankAuxiliary)
	c.WMap[0xC008] = bsc.UnSetterW(BankAuxiliary)
	c.WMap[0xC009] = bsc.ReSetterW(BankAuxiliary)
}
