package sixtwo

const (
	// VolumeMarker is a hard-coded number that encoding software will
	// expect to see in the track padding.
	VolumeMarker = 0xFE

	// NumTracks is the number of tracks that can be contained on a
	// disk.
	NumTracks = 35

	// NumSectors is the number of sectors that each track contains.
	NumSectors = 16

	// LogSectorLen is the length of a logical sector, which is 256
	// bytes.
	LogSectorLen = 0x100

	// LogTrackLen is the length of a logical track, consisting of 16
	// logical sectors, which are each 256 bytes long. It thus holds 4
	// kilobytes of data.
	LogTrackLen = 0x1000

	// PhysSectorLen is the length of a physical sector
	PhysSectorLen = 0x1A0

	// PhysSectorHeader is the length of a sector header
	PhysSectorHeader = 0x13

	// PhysTrackLen is the length of a physical track, consisting of 16
	// physical sectors.
	PhysTrackLen = 0x1A00

	// PhysTrackHeader is the length of a track header.
	PhysTrackHeader = 0x30

	// MaxSteps is the maximum number of steps we can move the drive
	// head before running out of tracks on the disk. (Note that steps
	// are half of the length of a track; 35 tracks, 70 steps.)
	MaxSteps = 70

	// MaxSectorPos is the highest sector position that we can allow
	// within a given track. (0xFFF = 4k - 1.)
	MaxSectorPos = 0xFFF

	// DosSize is the number of bytes in 140 kilobytes.
	DosSize = 143360

	// NibSize is the capacity of the segment we will create for
	// nibblized data, whether from 140k logical data or just any-old
	// NIB file.
	NibSize = 234640
)

const (
	// DOS33 is the image type for DOS 3.3, which is the
	// generally-used image type for Apple II DOS images. There are
	// other DOS versions, which are formatted differently, but we don't
	// handle them here.
	DOS33 = iota

	// ProDOS indicates that the image type is ProDOS.
	ProDOS
)
