package a2kb

import (
	"sync"

	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/internal/metrics"
	"github.com/pevans/erc/memory"
)

const (
	dataAndStrobe int = 0xC000
	anyKeyDown    int = 0xC010
)

func ReadSwitches() []int {
	switches := []int{
		dataAndStrobe,
		anyKeyDown,
	}

	// Addresses $C001-$C00F are write-only switches for other subsystems.
	// Reads from them return the keyboard data latch (open bus behavior).
	for a := 0xC001; a <= 0xC00F; a++ {
		switches = append(switches, a)
	}

	return switches
}

func WriteSwitches() []int {
	return []int{
		anyKeyDown,
	}
}

func kbMutex(stm *memory.StateMap) *sync.Mutex {
	return stm.Any(a2state.KBMutex).(*sync.Mutex)
}

func SwitchRead(addr int, stm *memory.StateMap) uint8 {
	mu := kbMutex(stm)
	mu.Lock()
	defer mu.Unlock()

	switch addr {
	case dataAndStrobe:
		metrics.Increment("soft_read_kb_data_and_strobe", 1)
		return stm.Uint8(a2state.KBLastKey) | stm.Uint8(a2state.KBStrobe)
	case anyKeyDown:
		metrics.Increment("soft_read_kb_any_key_down", 1)
		stm.SetUint8(a2state.KBStrobe, 0)
		return stm.Uint8(a2state.KBKeyDown)
	}

	// Open bus: addresses $C001-$C00F return the keyboard data latch
	if addr >= 0xC001 && addr <= 0xC00F {
		return stm.Uint8(a2state.KBLastKey) | stm.Uint8(a2state.KBStrobe)
	}

	return 0
}

func SwitchWrite(addr int, val uint8, stm *memory.StateMap) {
	if addr == anyKeyDown {
		mu := kbMutex(stm)
		mu.Lock()
		defer mu.Unlock()

		metrics.Increment("soft_write_kb_any_key_down", 1)
		stm.SetUint8(a2state.KBStrobe, 0)
	}
}

func UseDefaults(state *memory.StateMap) {
	state.SetUint8(a2state.KBLastKey, 0)
	state.SetUint8(a2state.KBKeyDown, 0)
	state.SetUint8(a2state.KBStrobe, 0)
}
