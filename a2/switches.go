package a2

import (
	"github.com/pevans/erc/a2/a2display"
)

// A Switcher is a type which provides a way to handle soft switch reads and
// writes in a relatively generic way.
type Switcher interface {
	SwitchRead(c *Computer, addr int) uint8
	SwitchWrite(c *Computer, addr int, val uint8)
}

// MapSoftSwitches will add several mappings for the soft switches that our
// computer uses.
func (c *Computer) MapSoftSwitches() {
	c.MapRange(0x0, 0x200, BankZPRead, BankZPWrite)
	c.MapRange(0x0400, 0x0800, a2display.Read, a2display.Write)
	c.MapRange(0x2000, 0x4000, a2display.Read, a2display.Write)
	// Note that there are other peripheral slots beginning with $C090, all
	// the way until $C100. We just don't emulate them right now.
	c.MapRange(0xC0E0, 0xC100, diskRead, diskWrite)
	c.MapRange(0xC100, 0xD000, PCRead, PCWrite)
	c.MapRange(0xD000, 0x10000, BankDFRead, BankDFWrite)

	for _, a := range kbReadSwitches() {
		c.smap.SetRead(a, kbSwitchRead)
	}

	for _, a := range kbWriteSwitches() {
		c.smap.SetWrite(a, kbSwitchWrite)
	}

	for _, a := range memReadSwitches() {
		c.smap.SetRead(a, memSwitchRead)
	}

	for _, a := range memWriteSwitches() {
		c.smap.SetWrite(a, memSwitchWrite)
	}

	for _, a := range bankReadSwitches() {
		c.smap.SetRead(a, bankSwitchRead)
	}

	for _, a := range bankWriteSwitches() {
		c.smap.SetWrite(a, bankSwitchWrite)
	}

	for _, a := range pcReadSwitches() {
		c.smap.SetRead(a, pcSwitchRead)
	}

	for _, a := range pcWriteSwitches() {
		c.smap.SetWrite(a, pcSwitchWrite)
	}

	for _, a := range a2display.ReadSwitches() {
		c.smap.SetRead(a, a2display.SwitchRead)
	}

	for _, a := range a2display.WriteSwitches() {
		c.smap.SetWrite(a, a2display.SwitchWrite)
	}

	for _, a := range speakerReadSwitches() {
		c.smap.SetRead(a, speakerSwitchRead)
	}
}
