package obj

import (
	"fmt"

	"github.com/pevans/erc/pkg/data"
)

// Slice returns a slice of data from the object store for a given set
// of coordinates. Nil is returned if there is an error.
func Slice(start, end int) ([]data.Byte, error) {
	if start < 0 || start >= len(storeData) || end >= len(storeData) {
		return nil, fmt.Errorf("coordinates out of bounds: %v:%v", start, end)
	}

	return storeData[start:end], nil
}
