package a2

import (
	"fmt"
	"os"
)

// A DiskSet is a container of disk image filenames, with some tracking for
// the current disk in the set. Many old software packages had more than one
// disk, so the idea is a single DiskSet can contain every disk you would need
// to operate the software.
type DiskSet struct {
	images  []string
	current int
}

// NewDiskSet returns a newly allocated empty diskset.
func NewDiskSet() *DiskSet {
	set := new(DiskSet)
	set.images = make([]string, 0)
	set.current = 0

	return set
}

// Append adds a disk to the diskset. Given some file, we will test that it's
// there, and then append the filename to the diskset. If no image file exists
// at the given filename, we return an error.
func (set *DiskSet) Append(file string) error {
	_, err := os.Stat(file)
	if err != nil {
		return fmt.Errorf("could not append file %v to diskset: %w", file, err)
	}

	set.images = append(set.images, file)

	return nil
}

// Disk returns the disk image at a given index. If the index is not valid, an
// error is returned.
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

// Current returns the current disk in the diskset (according to its index).
func (set *DiskSet) Current() (*os.File, string, error) {
	return set.Disk(set.current)
}

// Name is the name of the entire diskset, which we take to be the name of the
// first filename loaded in the set.
func (set *DiskSet) Name() string {
	if len(set.images) == 0 {
		return ""
	}

	return set.images[0]
}

// Next returns the next disk in the diskset (the index one after the current
// index). If we're at the end of the diskset, this will wrap around to the
// first disk in the set.
func (set *DiskSet) Next() (*os.File, string, error) {
	set.current++
	if set.current >= len(set.images) {
		set.current = 0
	}

	return set.Current()
}

// Previous returns the previous disk in the diskset (the index one earlier
// from the current index). If we're at the beginning of the diskset, this
// will wrap around to the last disk in the set.
func (set *DiskSet) Previous() (*os.File, string, error) {
	set.current--
	if set.current < 0 {
		set.current = len(set.images) - 1
	}

	return set.Current()
}

// CurrentIndex returns the index of the current disk within the diskset.
func (set *DiskSet) CurrentIndex() int {
	return set.current
}
