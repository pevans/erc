package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/pevans/erc/a2/a2disk"
	"github.com/pevans/erc/a2/a2enc"
	"github.com/pevans/erc/memory"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info [image]",
	Short: "Display VTOC information from a disk image",
	Long:  "Parse and display the Volume Table of Contents (VTOC) from a disk image",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		displayInfo(args[0])
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

func displayInfo(imagePath string) {
	file, err := os.Open(imagePath)
	if err != nil {
		fail(fmt.Sprintf("could not open file %s: %v", imagePath, err))
	}
	defer file.Close() //nolint:errcheck

	fileInfo, err := file.Stat()
	if err != nil {
		fail(fmt.Sprintf("could not stat file: %v", err))
	}

	if fileInfo.Size() != a2enc.DosSize {
		fail("file does not appear to be a disk image")
	}

	seg := memory.NewSegment(a2enc.DosSize)
	if err := seg.ReadFile(imagePath); err != nil {
		fail(fmt.Sprintf("could not read image file: %v", err))
	}

	var vtoc a2disk.VTOC
	if err := vtoc.Parse(seg); err != nil {
		fail(fmt.Sprintf("could not parse VTOC: %v", err))
	}

	if vtoc.DisketteVolume != 254 && vtoc.MaxTrackSectorPairs != 122 {
		fail(`file does not contain a valid volume table of contents.
this does not mean the disk is bad! it's normal for disks 
to omit this data, e.g. to increase free space on the disk.`,
		)
	}

	// Display the results
	fmt.Printf("VTOC Information for:      %s\n", imagePath)
	fmt.Printf("---------------------------%s\n", strings.Repeat("-", len(imagePath)))
	fmt.Printf(
		"First Catalog Track:       %d (0x%02X)\n",
		vtoc.FirstCatalogSectorTrackNumber,
		vtoc.FirstCatalogSectorTrackNumber,
	)
	fmt.Printf(
		"First Catalog Sector:      %d (0x%02X)\n",
		vtoc.FirstCatalogSectorSectorNumber,
		vtoc.FirstCatalogSectorSectorNumber,
	)
	fmt.Printf(
		"DOS Release Number:        %d\n",
		vtoc.ReleaseNumberOfDOS,
	)
	fmt.Printf(
		"Diskette Volume:           %d (0x%02X)\n",
		vtoc.DisketteVolume,
		vtoc.DisketteVolume,
	)
	fmt.Printf(
		"Max Track/Sector Pairs:    %d\n",
		vtoc.MaxTrackSectorPairs,
	)
	fmt.Printf(
		"Last Track Allocated:      %d\n",
		vtoc.LastTrackAllocated,
	)
	fmt.Printf(
		"Direction of Allocation:   %+d\n",
		vtoc.DirectionOfAllocation,
	)
	fmt.Printf(
		"Tracks per Diskette:       %d\n",
		vtoc.TracksPerDiskette,
	)
	fmt.Printf(
		"Sectors per Track:         %d\n",
		vtoc.SectorsPerTrack,
	)
	fmt.Printf(
		"Bytes per Sector:          %d\n",
		vtoc.BytesPerSector,
	)

	fmt.Println("\nTrack   Free  Sectors")
	fmt.Println("----- -----------------")

	for track := 0; track < int(vtoc.TracksPerDiskette); track++ {
		if freeSectors, ok := vtoc.FreeSectors[track*4]; ok {
			fmt.Printf("%2d    %s\n", track, freeSectors)
		}
	}
}
