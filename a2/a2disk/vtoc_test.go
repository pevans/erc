package a2disk

import (
	"testing"

	"github.com/pevans/erc/a2/a2enc"
	"github.com/pevans/erc/memory"
	"github.com/stretchr/testify/assert"
)

func TestVTOC_Parse(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*memory.Segment)
		expected VTOC
	}{
		{
			name: "typical DOS 3.3 VTOC",
			setup: func(seg *memory.Segment) {
				offset := a2enc.LogTrackLen * 17

				seg.Set(offset+0x1, 0x11)
				seg.Set(offset+0x2, 0x0F)
				seg.Set(offset+0x3, 0x03)
				seg.Set(offset+0x6, 0xFE)
				seg.Set(offset+0x27, 0x7A)
				seg.Set(offset+0x30, 0x12)
				seg.Set(offset+0x31, 0x01)
				seg.Set(offset+0x34, 0x23)
				seg.Set(offset+0x35, 0x10)
				seg.Set(offset+0x36, 0x00)
				seg.Set(offset+0x37, 0x01)

				seg.Set(offset+0x38, 0xFF)
				seg.Set(offset+0x39, 0xFF)
			},
			expected: VTOC{
				FirstCatalogSectorTrackNumber:  0x11,
				FirstCatalogSectorSectorNumber: 0x0F,
				ReleaseNumberOfDOS:             0x03,
				DisketteVolume:                 0xFE,
				MaxTrackSectorPairs:            0x7A,
				LastTrackAllocated:             0x12,
				DirectionOfAllocation:          1,
				TracksPerDiskette:              0x23,
				SectorsPerTrack:                0x10,
				BytesPerSector:                 256,
				FreeSectors: map[int]string{
					0: "FEDCBA98 76543210",
				},
			},
		},
		{
			name: "negative allocation direction",
			setup: func(seg *memory.Segment) {
				offset := a2enc.LogTrackLen * 17

				seg.Set(offset+0x1, 0x11)
				seg.Set(offset+0x2, 0x0F)
				seg.Set(offset+0x3, 0x03)
				seg.Set(offset+0x6, 0x01)
				seg.Set(offset+0x27, 0x7A)
				seg.Set(offset+0x30, 0x12)
				seg.Set(offset+0x31, 0xFF)
				seg.Set(offset+0x34, 0x23)
				seg.Set(offset+0x35, 0x10)
				seg.Set(offset+0x36, 0x00)
				seg.Set(offset+0x37, 0x01)

				seg.Set(offset+0x38, 0x33)
				seg.Set(offset+0x39, 0x22)
			},
			expected: VTOC{
				FirstCatalogSectorTrackNumber:  0x11,
				FirstCatalogSectorSectorNumber: 0x0F,
				ReleaseNumberOfDOS:             0x03,
				DisketteVolume:                 0x01,
				MaxTrackSectorPairs:            0x7A,
				LastTrackAllocated:             0x12,
				DirectionOfAllocation:          -1,
				TracksPerDiskette:              0x23,
				SectorsPerTrack:                0x10,
				BytesPerSector:                 256,
				FreeSectors: map[int]string{
					0: "..DC..98 ..5...1.",
				},
			},
		},
		{
			name: "different bytes per sector",
			setup: func(seg *memory.Segment) {
				offset := a2enc.LogTrackLen * 17

				seg.Set(offset+0x1, 0x11)
				seg.Set(offset+0x2, 0x0F)
				seg.Set(offset+0x3, 0x03)
				seg.Set(offset+0x6, 0x01)
				seg.Set(offset+0x27, 0x7A)
				seg.Set(offset+0x30, 0x12)
				seg.Set(offset+0x31, 0x01)
				seg.Set(offset+0x34, 0x23)
				seg.Set(offset+0x35, 0x10)
				seg.Set(offset+0x36, 0x00)
				seg.Set(offset+0x37, 0x02)

				seg.Set(offset+0x38, 0x00)
				seg.Set(offset+0x39, 0x55)
			},
			expected: VTOC{
				FirstCatalogSectorTrackNumber:  0x11,
				FirstCatalogSectorSectorNumber: 0x0F,
				ReleaseNumberOfDOS:             0x03,
				DisketteVolume:                 0x01,
				MaxTrackSectorPairs:            0x7A,
				LastTrackAllocated:             0x12,
				DirectionOfAllocation:          1,
				TracksPerDiskette:              0x23,
				SectorsPerTrack:                0x10,
				BytesPerSector:                 512,
				FreeSectors: map[int]string{
					0: "........ .6.4.2.0",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a segment large enough to hold a disk image
			seg := memory.NewSegment(a2enc.LogTrackLen * 35)
			tt.setup(seg)

			var vtoc VTOC

			err := vtoc.Parse(seg)
			assert.NoError(t, err)

			assert.Equal(t, tt.expected.FirstCatalogSectorTrackNumber, vtoc.FirstCatalogSectorTrackNumber)
			assert.Equal(t, tt.expected.FirstCatalogSectorSectorNumber, vtoc.FirstCatalogSectorSectorNumber)
			assert.Equal(t, tt.expected.ReleaseNumberOfDOS, vtoc.ReleaseNumberOfDOS)
			assert.Equal(t, tt.expected.DisketteVolume, vtoc.DisketteVolume)
			assert.Equal(t, tt.expected.MaxTrackSectorPairs, vtoc.MaxTrackSectorPairs)
			assert.Equal(t, tt.expected.LastTrackAllocated, vtoc.LastTrackAllocated)
			assert.Equal(t, tt.expected.DirectionOfAllocation, vtoc.DirectionOfAllocation)
			assert.Equal(t, tt.expected.TracksPerDiskette, vtoc.TracksPerDiskette)
			assert.Equal(t, tt.expected.SectorsPerTrack, vtoc.SectorsPerTrack)
			assert.Equal(t, tt.expected.BytesPerSector, vtoc.BytesPerSector)
			assert.Equal(t, tt.expected.FreeSectors[0], vtoc.FreeSectors[0])
		})
	}
}
