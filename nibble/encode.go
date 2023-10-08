package nibble

import "github.com/pevans/erc/memory"

// Encode copies and nibblizes a source segment. Because we presume here
// that the source segment is already nibblized, this is essentially a
// straight copy.
func Encode(src *memory.Segment) (*memory.Segment, error) {
	return nibbleCopier(src)
}
