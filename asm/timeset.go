package asm

import (
	"fmt"
	"os"
	"slices"
	"sync"
	"time"

	"golang.org/x/exp/maps"
)

// A TimeSetEntry records information about the execution of some instruction
// that is recorded in a TimeSet.
type TimeSetEntry struct {
	// executed is the number of times an instruction was executed
	executed int

	// cyclesConsumed is the total number cycles that were consumed by this
	// instruction
	cyclesConsumed int
}

// A TimeSet is a map of the instructions that were executed during emulation.
// Each instruction is paired with some metadata about the execution that can
// be used to determine the rough time spent on execution for each
// instruction.
type TimeSet struct {
	// calls is a map of instructions (the key) paired with TimeSetEntry
	// values.
	calls map[string]TimeSetEntry

	// mu is a mutex. We're likely to have a lot of writes and reads on the
	// map that happen concurrently during operation.
	mu sync.Mutex

	// timePerCycle is our belief of the time that each cycle takes to be
	// executed. This is something that is defined for us when the TimeSet is
	// allocated.
	timePerCycle time.Duration
}

// NewTimeSet returns a TimeSet oriented by the provided timePerCycle -- that
// is, we will assume each cycle consumes timePerCycle, and will use that to
// estimate the total cost of execution for each instruction.
func NewTimeset(timePerCycle time.Duration) *TimeSet {
	set := new(TimeSet)
	set.timePerCycle = timePerCycle

	set.calls = make(map[string]TimeSetEntry)

	return set
}

// formatted returns a formatted string representation of a TimeSetEntry, with
// respect to a given TimeSet (i.e. which has a particular alignment of
// timePerCycle).
func (entry *TimeSetEntry) formatted(set *TimeSet) string {
	return fmt.Sprintf(
		"run %v cyc %v spent %v",
		entry.executed,
		entry.cyclesConsumed,
		time.Duration(entry.cyclesConsumed)*set.timePerCycle,
	)
}

// Record will add a call to the TimeSet and set the number of cycles consumed
// by the provided arguments.
func (set *TimeSet) Record(call string, cycles int) {
	set.mu.Lock()
	defer set.mu.Unlock()

	// If no entry existed previously, this will give us a blank record.
	entry := set.calls[call]

	entry.executed++
	entry.cyclesConsumed += cycles

	set.calls[call] = entry
}

// WriteToFile will write the contents of the timeset to some provided
// filename. If that cannot be done, an error is returned.
func (set *TimeSet) WriteToFile(file string) error {
	set.mu.Lock()
	defer set.mu.Unlock()

	calls := maps.Keys(set.calls)
	slices.Sort(calls)

	fp, err := os.Create(file)
	if err != nil {
		return err
	}

	defer fp.Close() //nolint:errcheck

	for _, call := range calls {
		entry := set.calls[call]
		output := fmt.Sprintf("%-30v | %v\n", call, entry.formatted(set))

		if _, err := fp.WriteString(output); err != nil {
			return err
		}
	}

	return nil
}
