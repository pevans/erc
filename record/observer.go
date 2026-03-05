package record

import (
	"fmt"

	"github.com/pevans/erc/memory"
)

// An Observer watches a single piece of state and emits entries when that
// state changes.
type Observer interface {
	// Before captures a snapshot of the observed value prior to a step.
	Before()

	// Observe compares the current value against the snapshot taken by
	// Before. If the value changed, it returns one Entry; otherwise nil.
	Observe(step int) []Entry
}

// A FuncObserver is a generic observer backed by a getter function. The
// getter returns the current value of whatever is being watched.
type FuncObserver struct {
	tag    string
	name   string
	get    func() any
	before any
}

// NewObserver returns a FuncObserver that reports changes under the given tag
// and name. The get function is called before and after each step to detect
// changes.
func NewObserver(tag, name string, get func() any) *FuncObserver {
	return &FuncObserver{tag: tag, name: name, get: get}
}

// MemObserver returns a FuncObserver that watches a single byte in memory at
// the given address.
func MemObserver(getter memory.Getter, addr int) *FuncObserver {
	return NewObserver(TagMem, fmt.Sprintf("$%04X", addr), func() any {
		return getter.Get(addr)
	})
}

// Before captures the current value.
func (o *FuncObserver) Before() {
	o.before = o.get()
}

// Observe compares the current value to the snapshot. If they differ it
// returns a slice containing a single Entry.
func (o *FuncObserver) Observe(step int) []Entry {
	after := o.get()
	if after == o.before {
		return nil
	}

	return []Entry{{
		Step: step,
		Tag:  o.tag,
		Name: o.name,
		Old:  o.before,
		New:  after,
	}}
}
