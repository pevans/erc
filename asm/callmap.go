package asm

import (
	"fmt"
	"os"
	"sort"
)

type CallMap map[string]int

func (cm CallMap) Add(line string) {
	if _, ok := cm[line]; !ok {
		cm[line] = 1
		return
	}

	cm[line]++
}

func (cm CallMap) Lines() []string {
	lines := make([]string, len(cm))

	for line, count := range cm {
		lines = append(lines, fmt.Sprintf("%v (%v)\n", line, count))
	}

	sort.Strings(lines)

	return lines
}

func (cm CallMap) WriteToFile(file string) error {
	lines := cm.Lines()

	fp, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0644)
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
