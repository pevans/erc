package a2

import (
	"fmt"
	"os"
)

type DiskSet struct {
	images  []string
	current int
}

func NewDiskSet() *DiskSet {
	set := new(DiskSet)
	set.images = make([]string, 0)
	set.current = 0

	return set
}

func (set *DiskSet) Append(file string) error {
	_, err := os.Stat(file)
	if err != nil {
		return fmt.Errorf("could not append file %v to diskset: %w", file, err)
	}

	set.images = append(set.images, file)

	return nil
}

func (set *DiskSet) Disk(index int) (*os.File, string, error) {
	if index < 0 || index >= len(set.images) {
		return nil, "", fmt.Errorf("no disk at index %v", index)
	}

	file := set.images[index]

	reader, err := os.OpenFile(file, os.O_RDWR, 0o644)
	if err != nil {
		return nil, "", fmt.Errorf("could not open file %v: %w", file, err)
	}

	return reader, set.images[index], nil
}

// Reset the diskset position to the first file and return that
func (set *DiskSet) First() (*os.File, string, error) {
	set.current = 0
	return set.Current()
}

func (set *DiskSet) Current() (*os.File, string, error) {
	return set.Disk(set.current)
}

func (set *DiskSet) Next() (*os.File, string, error) {
	set.current++
	if set.current >= len(set.images) {
		set.current = 0
	}

	return set.Current()
}
