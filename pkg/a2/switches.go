package a2

import "github.com/pevans/erc/pkg/data"

// A Switcher is a type which provides a way to handle soft switch reads and
// writes in a relatively generic way.
type Switcher interface {
	SwitchRead(c *Computer, addr data.DByte) data.Byte
	SwitchWrite(c *Computer, addr data.DByte, val data.Byte)
}

var bankReadSwitches = []int{
	0xC011,
	0xC012,
	0xC016,
	0xC080,
	0xC081,
	0xC082,
	0xC083,
	0xC088,
	0xC089,
	0xC08A,
	0xC08B,
}

var bankWriteSwitches = []int{
	0xC008,
	0xC009,
}

// MapSoftSwitches will add several mappings for the soft switches that our
// computer uses.
func (c *Computer) MapSoftSwitches() {
	c.MapRange(0x0, 0x200, BankZPRead, BankZPWrite)
	c.MapRange(0x0400, 0x0800, displayRead, displayWrite)
	c.MapRange(0x2000, 0x4000, displayRead, displayWrite)
	// Note that there are other peripheral slots beginning with $C090, all the
	// way until $C100. We just don't emulate them right now.
	c.MapRange(0xC0E0, 0xC100, diskRead, diskWrite)
	c.MapRange(0xC100, 0xD000, pcRead, pcWrite)
	c.MapRange(0xD000, 0x10000, BankDFRead, BankDFWrite)

	for _, addr := range []int{0xC013, 0xC014} {
		c.RMap[addr] = memSwitchRead
	}

	for _, addr := range []int{0xC002, 0xC003, 0xC004, 0xC005} {
		c.WMap[addr] = memSwitchWrite
	}

	for _, addr := range bankReadSwitches {
		c.RMap[addr] = bankSwitchRead
	}

	for _, addr := range bankWriteSwitches {
		c.WMap[addr] = bankSwitchWrite
	}

	psc := newPCSwitchCheck()
	c.RMap[0xC015] = psc.IsOpSetter(PCSlotCxROM)
	c.RMap[0xC017] = psc.IsSetter(PCSlotC3ROM)
	c.WMap[0xC006] = psc.ReSetterW(PCSlotCxROM)
	c.WMap[0xC007] = psc.UnSetterW(PCSlotCxROM)
	c.WMap[0xC00A] = psc.UnSetterW(PCSlotC3ROM)
	c.WMap[0xC00B] = psc.ReSetterW(PCSlotC3ROM)

	dsc := newDisplaySwitchCheck()
	c.RMap[0xC018] = dsc.IsSetter(Display80Store)
	c.RMap[0xC01A] = dsc.IsSetter(DisplayText)
	c.RMap[0xC01B] = dsc.IsSetter(DisplayMixed)
	c.RMap[0xC01C] = dsc.IsSetter(DisplayPage2)
	c.RMap[0xC01D] = dsc.IsSetter(DisplayHires)
	c.RMap[0xC01E] = dsc.IsSetter(DisplayAltCharset)
	c.RMap[0xC01F] = dsc.IsSetter(Display80Col)
	c.RMap[0xC050] = dsc.UnSetterR(DisplayText)
	c.RMap[0xC051] = dsc.ReSetterR(DisplayText)
	c.RMap[0xC052] = dsc.UnSetterR(DisplayMixed)
	c.RMap[0xC053] = dsc.ReSetterR(DisplayMixed)
	c.RMap[0xC054] = dsc.UnSetterR(DisplayPage2)
	c.RMap[0xC055] = dsc.ReSetterR(DisplayPage2)
	c.RMap[0xC056] = dsc.UnSetterR(DisplayHires)
	c.RMap[0xC057] = dsc.ReSetterR(DisplayHires)
	c.RMap[0xC05E] = dsc.ReSetterR(DisplayDHires)
	c.RMap[0xC05F] = dsc.UnSetterR(DisplayDHires)
	c.RMap[0xC07E] = dsc.IsSetter(DisplayIOU)
	c.RMap[0xC07F] = dsc.IsSetter(DisplayDHires)
	c.WMap[0xC000] = dsc.UnSetterW(Display80Store)
	c.WMap[0xC001] = dsc.ReSetterW(Display80Store)
	c.WMap[0xC00C] = dsc.UnSetterW(Display80Col)
	c.WMap[0xC00D] = dsc.ReSetterW(Display80Col)
	c.WMap[0xC00E] = dsc.UnSetterW(DisplayAltCharset)
	c.WMap[0xC00F] = dsc.ReSetterW(DisplayAltCharset)
	c.WMap[0xC050] = dsc.UnSetterW(DisplayText)
	c.WMap[0xC051] = dsc.ReSetterW(DisplayText)
	c.WMap[0xC052] = dsc.UnSetterW(DisplayMixed)
	c.WMap[0xC053] = dsc.ReSetterW(DisplayMixed)
	c.WMap[0xC054] = dsc.UnSetterW(DisplayPage2)
	c.WMap[0xC055] = dsc.ReSetterW(DisplayPage2)
	c.WMap[0xC056] = dsc.UnSetterW(DisplayHires)
	c.WMap[0xC057] = dsc.ReSetterW(DisplayHires)
	c.WMap[0xC05E] = dsc.ReSetterW(DisplayDHires)
	c.WMap[0xC05F] = dsc.UnSetterW(DisplayDHires)
	c.WMap[0xC07E] = dsc.ReSetterW(DisplayIOU)
	c.WMap[0xC07F] = dsc.UnSetterW(DisplayIOU)
}
