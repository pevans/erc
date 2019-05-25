package nibble

import "github.com/pevans/erc/pkg/data"

func Encode(src *data.Segment) (*data.Segment, error) {
	dst := data.NewSegment(src.Size())
	_, err := dst.CopySlice(0, src.Mem)

	if err != nil {
		return nil, err
	}

	return dst, nil
}
