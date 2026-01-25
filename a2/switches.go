package a2

import (
	"github.com/pevans/erc/a2/a2bank"
	"github.com/pevans/erc/a2/a2display"
	"github.com/pevans/erc/a2/a2kb"
	"github.com/pevans/erc/a2/a2peripheral"
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
	c.MapRange(0x0, 0x200, a2bank.ZPRead, a2bank.ZPWrite)
	c.MapRange(0x0400, 0x0800, a2display.Read, a2display.Write)
	c.MapRange(0x2000, 0x4000, a2display.Read, a2display.Write)
	// Note that there are other peripheral slots beginning with $C090, all
	// the way until $C100. We just don't emulate them right now.
	c.MapRange(0xC0E0, 0xC100, diskRead, diskWrite)
	c.MapRange(0xC100, 0xD000, a2peripheral.Read, a2peripheral.Write)
	c.MapRange(0xD000, 0x10000, a2bank.DFRead, a2bank.DFWrite)

	for _, a := range a2kb.ReadSwitches() {
		c.smap.SetRead(a, a2kb.SwitchRead)
	}

	for _, a := range a2kb.WriteSwitches() {
		c.smap.SetWrite(a, a2kb.SwitchWrite)
	}

	for _, a := range memReadSwitches() {
		c.smap.SetRead(a, memSwitchRead)
	}

	for _, a := range memWriteSwitches() {
		c.smap.SetWrite(a, memSwitchWrite)
	}

	for _, a := range a2bank.ReadSwitches() {
		c.smap.SetRead(a, a2bank.SwitchRead)
	}

	for _, a := range a2bank.WriteSwitches() {
		c.smap.SetWrite(a, a2bank.SwitchWrite)
	}

	for _, a := range a2peripheral.ReadSwitches() {
		c.smap.SetRead(a, a2peripheral.SwitchRead)
	}

	for _, a := range a2peripheral.WriteSwitches() {
		c.smap.SetWrite(a, a2peripheral.SwitchWrite)
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
