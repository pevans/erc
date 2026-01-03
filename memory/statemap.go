package memory

import "sync"

// A StateMap is a map of some state (i.e. from the a2state package) to some
// other value. We don't use generics because we allow a mixed set of types in
// the map. StateMaps are concurrent-safe.
type StateMap struct {
	m sync.Map
}

// NewStateMap returns a new StateMap ready to use.
func NewStateMap() *StateMap {
	sm := new(StateMap)
	return sm
}

// Int returns a integer key value if one is available. If not, we'll return a
// zero value. If the value is not an int, we'll panic.
func (sm *StateMap) Int(k int) int {
	v, ok := sm.m.Load(k)
	if !ok {
		return 0
	}

	return v.(int)
}

// SetInt assigns some integer value v to some key k.
func (sm *StateMap) SetInt(k, v int) {
	sm.m.Store(k, v)
}

// Int64 returns a int64 value for some key k. If the key doesn't exist, a
// zero is returned. If it's not an int64 value, we'll panic.
func (sm *StateMap) Int64(k int) int64 {
	v, ok := sm.m.Load(k)
	if !ok {
		return 0
	}

	return v.(int64)
}

// SetInt64 assigns some 64-bit integer value v to some key k.
func (sm *StateMap) SetInt64(k int, v int64) {
	sm.m.Store(k, v)
}

// Uint8 returns a uint8 value for some key k. If the key doesn't exist, a
// zero is returned. If it's not a uint8 value, we'll panic.
func (sm *StateMap) Uint8(k int) uint8 {
	v, ok := sm.m.Load(k)
	if !ok {
		return 0
	}

	return v.(uint8)
}

// SetUint8 assigns some 8-bit unsigned integer value v to some key k.
func (sm *StateMap) SetUint8(k int, v uint8) {
	sm.m.Store(k, v)
}

// Uint16 returns a uint16 value for some key k. If the key doesn't exist, a
// zero is returned. If it's not a uint16 value, we'll panic.
func (sm *StateMap) Uint16(k int) uint16 {
	v, ok := sm.m.Load(k)
	if !ok {
		return 0
	}

	return v.(uint16)
}

// SetUint16 assigns some 16-bit unsigned integer value v to some key k.
func (sm *StateMap) SetUint16(k int, v uint16) {
	sm.m.Store(k, v)
}

// Bool returns a boolean value for some key k. If the key doesn't exist,
// we'll return false. If it's not the correct value, we'll panic.
func (sm *StateMap) Bool(k int) bool {
	v, ok := sm.m.Load(k)
	if !ok {
		return false
	}

	return v.(bool)
}

// SetBool assigns some boolean value v to some key k.
func (sm *StateMap) SetBool(k int, v bool) {
	sm.m.Store(k, v)
}

// Segment returns a Segment-typed value for some key k. If no key exists,
// we'll return nil. If it's not the right type, we'll panic.
func (sm *StateMap) Segment(k int) *Segment {
	v, ok := sm.m.Load(k)
	if !ok {
		return nil
	}

	return v.(*Segment)
}

// SetSegment assigns some Segment v to some key k.
func (sm *StateMap) SetSegment(k int, v *Segment) {
	sm.m.Store(k, v)
}

// Any will return an any-typed value for some key k. It will be up to the
// caller to understand how to assert the type of the value. You normally
// don't want to use this method, but because StateMap is part of memory and
// that is imported by so many packages, referring to types from other
// packages could easily create a circular dependency.
func (sm *StateMap) Any(k int) any {
	v, ok := sm.m.Load(k)
	if !ok {
		return nil
	}

	return v
}

// SetAny will take an "any" interface value v and assign that to some key k.
func (sm *StateMap) SetAny(k int, v any) {
	sm.m.Store(k, v)
}

// Map produces a form of the StateMap such that the keys are converted to
// readable strings (i.e. not raw integers). Mostly useful for debugging.
func (sm *StateMap) Map(keyConvert func(any) string) map[string]any {
	plainMap := make(map[string]any)

	sm.m.Range(func(key, val any) bool {
		plainMap[keyConvert(key)] = val
		return true
	})

	return plainMap
}
