package a2disk

import (
	"fmt"

	"github.com/pevans/erc/a2/a2enc"
	"github.com/pevans/erc/memory"
)

// This is an admittedly experimental structure. There are several ways to
// consider a disk image. One is to strictly carve it up into tracks and
// sectors. Another is to read the VTOC, then the catalog, and discover
// discrete files. Although we can begin with the former for now, it's my goal
// to support the latter over time.
type Image struct {
	Tracks []*memory.Segment
}

func (img *Image) Parse(seg *memory.Segment) error {
	maxTracks := a2enc.MaxSteps / 2

	for track := 0; track < maxTracks; track++ {
		tseg := memory.NewSegment(a2enc.LogTrackLen)

		count, err := tseg.ExtractFrom(seg, track*a2enc.LogTrackLen, (track+1)*a2enc.LogTrackLen)
		if err != nil {
			return fmt.Errorf("failed to extract data from disk image: %w", err)
		}

		if count != a2enc.LogTrackLen {
			return fmt.Errorf("did not extract the number of expected bytes")
		}
	}

	return nil
}
