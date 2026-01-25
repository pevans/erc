package a2kb

import (
	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/internal/metrics"
	"github.com/pevans/erc/memory"
)

const (
	dataAndStrobe int = 0xC000
	anyKeyDown    int = 0xC010
)

func ReadSwitches() []int {
	return []int{
		dataAndStrobe,
		anyKeyDown,
	}
}

func WriteSwitches() []int {
	return []int{
		anyKeyDown,
	}
}

func SwitchRead(addr int, stm *memory.StateMap) uint8 {
	switch addr {
	case dataAndStrobe:
		metrics.Increment("soft_read_kb_data_and_strobe", 1)
		return stm.Uint8(a2state.KBLastKey) | stm.Uint8(a2state.KBStrobe)
	case anyKeyDown:
		metrics.Increment("soft_read_kb_any_key_down", 1)
		stm.SetUint8(a2state.KBStrobe, 0)
		return stm.Uint8(a2state.KBKeyDown)
	}

	// Nothing else can really come in here, but if something did...
	return 0
}

func SwitchWrite(addr int, val uint8, stm *memory.StateMap) {
	if addr == anyKeyDown {
		metrics.Increment("soft_write_kb_any_key_down", 1)
		stm.SetUint8(a2state.KBStrobe, 0)
	}
}

func UseDefaults(state *memory.StateMap) {
	state.SetUint8(a2state.KBLastKey, 0)
	state.SetUint8(a2state.KBKeyDown, 0)
	state.SetUint8(a2state.KBStrobe, 0)
}
