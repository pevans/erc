package asm

import (
	"fmt"
	"os"
	"slices"
	"sync"
	"time"

	"golang.org/x/exp/maps"
)

type TimeSetEntry struct {
	// number of times an instruction was executed
	Executed int

	// How many total cycles were consumed by this instruction
	CyclesConsumed int
}

type TimeSet struct {
	calls map[string]TimeSetEntry
	mu    sync.Mutex

	timePerCycle time.Duration
}

func NewTimeset(timePerCycle time.Duration) *TimeSet {
	set := new(TimeSet)
	set.timePerCycle = timePerCycle

	set.calls = make(map[string]TimeSetEntry)

	return set
}

func (set *TimeSet) Record(call string, cycles int) {
	set.mu.Lock()
	defer set.mu.Unlock()

	// If no entry existed previously, this will give us a blank record.
	entry := set.calls[call]

	entry.Executed++
	entry.CyclesConsumed += cycles

	set.calls[call] = entry
}

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
		output := fmt.Sprintf("%-30v | %v\n", call, entry.Formatted(set))

		if _, err := fp.WriteString(output); err != nil {
			return err
		}
	}

	return nil
}

func (entry *TimeSetEntry) Formatted(set *TimeSet) string {
	return fmt.Sprintf(
		"run %v cyc %v spent %v",
		entry.Executed,
		entry.CyclesConsumed,
		time.Duration(entry.CyclesConsumed)*set.timePerCycle,
	)
}
