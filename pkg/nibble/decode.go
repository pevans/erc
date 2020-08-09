package nibble

import "github.com/pevans/erc/pkg/data"

// Decode will return the nibblized version of the source data segment
// in a new copy. (Note that since the source segment is _already_
// nibblized, this is basically a straight copy.)
func Decode(src *data.Segment) (*data.Segment, error) {
	return nibbleCopier(src)
}

func nibbleCopier(src *data.Segment) (*data.Segment, error) {
	dst := data.NewSegment(src.Size())

	_, err := dst.CopySlice(0, src.Mem)
	if err != nil {
		return nil, err
	}

	return dst, nil
}
