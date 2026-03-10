package record

import "strings"

// A Recorder collects state entries emitted by observers over a sequence of
// execution steps.
type Recorder struct {
	entries   []Entry
	step      int
	observers []Observer
}

// Add registers one or more observers with the recorder.
func (r *Recorder) Add(observers ...Observer) {
	r.observers = append(r.observers, observers...)
}

// Step calls fn (typically one emulator instruction execution), bracketing it
// with before/after snapshots from all registered observers. Any changes
// detected are recorded as entries.
func (r *Recorder) Step(fn func()) {
	for _, obs := range r.observers {
		obs.Before()
	}

	fn()

	r.step++

	for _, obs := range r.observers {
		r.entries = append(r.entries, obs.Observe(r.step)...)
	}
}

// Entries returns the full ordered list of recorded state entries.
func (r *Recorder) Entries() []Entry {
	return r.entries
}

// CurrentStep returns the number of steps executed so far.
func (r *Recorder) CurrentStep() int {
	return r.step
}

// String renders all recorded entries as newline-separated text.
func (r *Recorder) String() string {
	lines := make([]string, len(r.entries))
	for i, e := range r.entries {
		lines[i] = e.String()
	}

	return strings.Join(lines, "\n")
}
