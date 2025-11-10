package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/pevans/erc/a2"
	"github.com/pevans/erc/a2/a2enc"
	"github.com/pevans/erc/memory"
	"github.com/spf13/cobra"
)

var outputFlag string

var encodeCmd = &cobra.Command{
	Use:   "encode [image]",
	Short: "Encode a logical disk image to physical nibblized format",
	Long:  "Encode a logically formatted disk image (DOS 3.3 or ProDOS) into a physical image file that is 6-and-2 encoded.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if outputFlag == "" {
			fail("output file must be specified with -o flag")
		}
		encodeImage(args[0], outputFlag)
	},
}

func init() {
	rootCmd.AddCommand(encodeCmd)

	encodeCmd.Flags().StringVarP(&outputFlag, "output", "o", "", "Output file path (required)")
	encodeCmd.MarkFlagRequired("output") //nolint:errcheck
}

func encodeImage(inputPath, outputPath string) {
	imageType, err := a2.ImageType(inputPath)
	if err != nil {
		fail(fmt.Sprintf("could not determine image type: %v", err))
	}

	if imageType == a2enc.Nibble {
		fail("input file is already in nibble format (.nib)")
	}

	inputFile, err := os.Open(inputPath)
	if err != nil {
		fail(fmt.Sprintf("could not open input file %s: %v", inputPath, err))
	}
	defer inputFile.Close() //nolint:errcheck

	bytes, err := io.ReadAll(inputFile)
	if err != nil {
		fail(fmt.Sprintf("could not read input file: %v", err))
	}

	if len(bytes) != a2enc.DosSize {
		fail(fmt.Sprintf(
			"input file has unexpected size: %d (given) != %d (expected)",
			len(bytes), a2enc.DosSize,
		))
	}

	// We need to turn the image file bytes into a segment for the encoder to
	// function.
	logicalSeg := memory.NewSegment(len(bytes))
	_, err = logicalSeg.CopySlice(0, []uint8(bytes))
	if err != nil {
		fail(fmt.Sprintf("could not copy bytes to segment: %v", err))
	}

	physicalSeg, err := a2enc.Encode(imageType, logicalSeg)
	if err != nil {
		fail(fmt.Sprintf("could not encode image: %v", err))
	}

	if err := physicalSeg.WriteFile(outputPath); err != nil {
		fail(fmt.Sprintf("could not write output file: %v", err))
	}

	fmt.Printf(
		"successfully encoded %s to %s\n", inputPath, outputPath,
	)
}
