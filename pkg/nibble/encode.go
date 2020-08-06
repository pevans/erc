package nibble

import "github.com/pevans/erc/pkg/data"

// Encode copies and nibblizes a source segment. Because we presume here
// that the source segment is already nibblized, this is essentially a
// straight copy.
func Encode(src *data.Segment) (*data.Segment, error) {
	// Encode and Decode essentially do the same thing; so while we
	// expose an API for encoding, we don't really need to do anything
	// different than what we do in the Decode function.
	return Decode(src)
}
