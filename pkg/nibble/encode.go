package nibble

import "github.com/pevans/erc/pkg/data"

// Encode copies and nibblizes a source segment. Because we presume here
// that the source segment is already nibblized, this is essentially a
// straight copy.
func Encode(src *data.Segment) (*data.Segment, error) {
	dst := data.NewSegment(src.Size())
	_, err := dst.CopySlice(0, src.Mem)

	if err != nil {
		return nil, err
	}

	return dst, nil
}
