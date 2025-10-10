package asm

import (
	"os"
	"sort"
	"sync"
)

type CallMap struct {
	m  map[string]int
	mu sync.Mutex
}

func NewCallMap() *CallMap {
	return &CallMap{
		m:  make(map[string]int),
		mu: sync.Mutex{},
	}
}

func (cm *CallMap) Add(line string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if _, ok := cm.m[line]; !ok {
		cm.m[line] = 1
		return
	}

	cm.m[line]++
}

func (cm *CallMap) Lines() []string {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	lines := make([]string, len(cm.m))

	for line, _ := range cm.m {
		lines = append(lines, line+"\n")
	}

	sort.Strings(lines)

	return lines
}

func (cm *CallMap) WriteToFile(file string) error {
	lines := cm.Lines()

	fp, err := os.Create(file)
	if err != nil {
		return err
	}

	defer fp.Close() //nolint:errcheck

	for _, line := range lines {
		if _, err := fp.WriteString(line); err != nil {
			return err
		}
	}

	return nil
}
