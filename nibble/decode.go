package nibble

import "github.com/pevans/erc/memory"

// Decode will return the nibblized version of the source data segment
// in a new copy. (Note that since the source segment is _already_
// nibblized, this is basically a straight copy.)
func Decode(src *memory.Segment) (*memory.Segment, error) {
	return nibbleCopier(src)
}

func nibbleCopier(src *memory.Segment) (*memory.Segment, error) {
	dst := memory.NewSegment(src.Size())

	_, err := dst.CopySlice(0, src.Mem)
	if err != nil {
		return nil, err
	}

	return dst, nil
}
