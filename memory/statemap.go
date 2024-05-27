package memory

import "sync"

type StateMap struct {
	m sync.Map
}

func NewStateMap() *StateMap {
	sm := new(StateMap)
	return sm
}

func (sm *StateMap) Int(k int) int {
	v, ok := sm.m.Load(k)
	if !ok {
		return 0
	}

	return v.(int)
}

func (sm *StateMap) SetInt(k, v int) {
	sm.m.Store(k, v)
}

func (sm *StateMap) Uint8(k int) uint8 {
	v, ok := sm.m.Load(k)
	if !ok {
		return 0
	}

	return v.(uint8)
}

func (sm *StateMap) SetUint8(k int, v uint8) {
	sm.m.Store(k, v)
}

func (sm *StateMap) Uint16(k int) uint16 {
	v, ok := sm.m.Load(k)
	if !ok {
		return 0
	}

	return v.(uint16)
}

func (sm *StateMap) SetUint16(k int, v uint16) {
	sm.m.Store(k, v)
}

func (sm *StateMap) Bool(k int) bool {
	v, ok := sm.m.Load(k)
	if !ok {
		return false
	}

	return v.(bool)
}

func (sm *StateMap) SetBool(k int, v bool) {
	sm.m.Store(k, v)
}

func (sm *StateMap) Segment(k int) *Segment {
	v, ok := sm.m.Load(k)
	if !ok {
		return nil
	}

	return v.(*Segment)
}

func (sm *StateMap) SetSegment(k int, v *Segment) {
	sm.m.Store(k, v)
}

func (sm *StateMap) Any(k int) interface{} {
	v, ok := sm.m.Load(k)
	if !ok {
		return nil
	}

	return v
}

func (sm *StateMap) SetAny(k int, v interface{}) {
	sm.m.Store(k, v)
}

func (sm *StateMap) Map(keyConvert func(any) string) map[string]any {
	plainMap := make(map[string]any)

	sm.m.Range(func(key, val any) bool {
		plainMap[keyConvert(key)] = val
		return true
	})

	return plainMap
}
