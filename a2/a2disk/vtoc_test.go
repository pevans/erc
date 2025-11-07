package a2disk_test

import (
	"testing"

	"github.com/pevans/erc/a2/a2disk"
	"github.com/pevans/erc/a2/a2enc"
	"github.com/pevans/erc/memory"
	"github.com/stretchr/testify/assert"
)

func TestVTOC_Parse(t *testing.T) {
	cases := []struct {
		name     string
		setup    func(*memory.Segment)
		expected a2disk.VTOC
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
			expected: a2disk.VTOC{
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
			expected: a2disk.VTOC{
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
			expected: a2disk.VTOC{
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

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Create a segment large enough to hold a disk image
			seg := memory.NewSegment(a2enc.LogTrackLen * 35)
			c.setup(seg)

			var vtoc a2disk.VTOC

			err := vtoc.Parse(seg)
			assert.NoError(t, err)

			assert.Equal(t, c.expected.FirstCatalogSectorTrackNumber, vtoc.FirstCatalogSectorTrackNumber)
			assert.Equal(t, c.expected.FirstCatalogSectorSectorNumber, vtoc.FirstCatalogSectorSectorNumber)
			assert.Equal(t, c.expected.ReleaseNumberOfDOS, vtoc.ReleaseNumberOfDOS)
			assert.Equal(t, c.expected.DisketteVolume, vtoc.DisketteVolume)
			assert.Equal(t, c.expected.MaxTrackSectorPairs, vtoc.MaxTrackSectorPairs)
			assert.Equal(t, c.expected.LastTrackAllocated, vtoc.LastTrackAllocated)
			assert.Equal(t, c.expected.DirectionOfAllocation, vtoc.DirectionOfAllocation)
			assert.Equal(t, c.expected.TracksPerDiskette, vtoc.TracksPerDiskette)
			assert.Equal(t, c.expected.SectorsPerTrack, vtoc.SectorsPerTrack)
			assert.Equal(t, c.expected.BytesPerSector, vtoc.BytesPerSector)
			assert.Equal(t, c.expected.FreeSectors[0], vtoc.FreeSectors[0])
		})
	}
}

func TestVTOC_Valid(t *testing.T) {
	cases := []struct {
		name   string
		vtoc   a2disk.VTOC
		testfn assert.BoolAssertionFunc
	}{
		{
			name: "valid VTOC with typical DOS 3.3 values",
			vtoc: a2disk.VTOC{
				DisketteVolume:      0xFE,
				MaxTrackSectorPairs: 122,
			},
			testfn: assert.True,
		},
		{
			name: "invalid VTOC with wrong diskette volume",
			vtoc: a2disk.VTOC{
				DisketteVolume:      0x01,
				MaxTrackSectorPairs: 122,
			},
			testfn: assert.False,
		},
		{
			name: "invalid VTOC with wrong max track sector pairs",
			vtoc: a2disk.VTOC{
				DisketteVolume:      0xFE,
				MaxTrackSectorPairs: 100,
			},
			testfn: assert.False,
		},
		{
			name: "invalid VTOC with both wrong",
			vtoc: a2disk.VTOC{
				DisketteVolume:      0x00,
				MaxTrackSectorPairs: 0,
			},
			testfn: assert.False,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.testfn(t, c.vtoc.Valid())
		})
	}
}
