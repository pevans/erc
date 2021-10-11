package data

type StateMap struct {
	m map[int]interface{}
}

func NewStateMap() *StateMap {
	sm := new(StateMap)
	sm.m = make(map[int]interface{})
	return sm
}

func (sm *StateMap) Int(k int) int {
	v, ok := sm.m[k]
	if !ok {
		return 0
	}

	return v.(int)
}

func (sm *StateMap) SetInt(k, v int) {
	sm.m[k] = v
}

func (sm *StateMap) Uint8(k int) uint8 {
	v, ok := sm.m[k]
	if !ok {
		return 0
	}

	return v.(uint8)
}

func (sm *StateMap) SetUint8(k int, v uint8) {
	sm.m[k] = v
}

func (sm *StateMap) Uint16(k int) uint16 {
	v, ok := sm.m[k]
	if !ok {
		return 0
	}

	return v.(uint16)
}

func (sm *StateMap) SetUint16(k int, v uint16) {
	sm.m[k] = v
}
