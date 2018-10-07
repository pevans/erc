package a2

// Here we set up all the soft switches that we'll use in the computer,
// which is a lot!
func (c *Computer) defineSoftSwitches() {
	c.MapRange(0x0, 0x200, zeroPageRead, zeroPageWrite)

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
}
