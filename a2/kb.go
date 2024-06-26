package a2

import (
	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/internal/metrics"
	"github.com/pevans/erc/memory"
)

const (
	kbDataAndStrobe int = 0xC000
	kbAnyKeyDown    int = 0xC010
)

func kbReadSwitches() []int {
	return []int{
		kbDataAndStrobe,
		kbAnyKeyDown,
	}
}

func kbWriteSwitches() []int {
	return []int{
		kbAnyKeyDown,
	}
}

func kbSwitchRead(addr int, stm *memory.StateMap) uint8 {
	switch addr {
	case kbDataAndStrobe:
		metrics.Increment("soft_read_kb_data_and_strobe", 1)
		return stm.Uint8(a2state.KBLastKey) | stm.Uint8(a2state.KBStrobe)
	case kbAnyKeyDown:
		metrics.Increment("soft_read_kb_any_key_down", 1)
		stm.SetUint8(a2state.KBStrobe, 0)
		return stm.Uint8(a2state.KBKeyDown)
	}

	// Nothing else can really come in here, but if something did...
	return 0
}

func kbSwitchWrite(addr int, val uint8, stm *memory.StateMap) {
	if addr == kbAnyKeyDown {
		metrics.Increment("soft_write_kb_any_key_down", 1)
		stm.SetUint8(a2state.KBStrobe, 0)
	}
}

func kbUseDefaults(c *Computer) {
	c.State.SetUint8(a2state.KBLastKey, 0)
	c.State.SetUint8(a2state.KBKeyDown, 0)
	c.State.SetUint8(a2state.KBStrobe, 0)
}

func (c *Computer) PressKey(key uint8) {
	// There can only be 7-bit ASCII in an Apple II, so we explicitly
	// take off the high bit.
	c.State.SetUint8(a2state.KBLastKey, key&0x7F)

	// We need to set the strobe bit, which (when returned) is always
	// with the high bit at 1.
	c.State.SetUint8(a2state.KBStrobe, 0x80)

	// This flag (again with the high bit set to 1) is set _while_ a key
	// is pressed.
	c.State.SetUint8(a2state.KBKeyDown, 0x80)
}

func (c *Computer) ClearKeys() {
	c.State.SetUint8(a2state.KBKeyDown, 0)
}
