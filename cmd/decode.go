package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/pevans/erc/a2/a2drive"
	"github.com/pevans/erc/a2/a2enc"
	"github.com/pevans/erc/memory"
	"github.com/spf13/cobra"
)

var decodeOutputFlag string

var decodeCmd = &cobra.Command{
	Use:   "decode [encoded-file]",
	Short: "Decode a physical disk image back to logical format",
	Long:  "Decode a physically encoded disk image (6-and-2 nibblized) back to a logical disk image (DOS 3.3 or ProDOS).",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if decodeOutputFlag == "" {
			fail("output file must be specified with -o flag")
		}
		decodeImage(args[0], decodeOutputFlag)
	},
}

func init() {
	rootCmd.AddCommand(decodeCmd)

	decodeCmd.Flags().StringVarP(&decodeOutputFlag, "output", "o", "", "Output file path (.dsk, .do, or .po) (required)")
	decodeCmd.MarkFlagRequired("output") //nolint:errcheck
}

func decodeImage(inputPath, outputPath string) {
	imageType, err := a2drive.ImageType(outputPath)
	if err != nil {
		fail(fmt.Sprintf("could not determine image type from output extension: %v", err))
	}

	if imageType == a2enc.Nibble {
		fail("cannot decode to nibble format (.nib)")
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

	if len(bytes) != a2enc.EncodedSize {
		fail(fmt.Sprintf(
			"input file has unexpected size: %d (given) != %d (expected)",
			len(bytes), a2enc.EncodedSize,
		))
	}

	physicalSeg := memory.NewSegment(len(bytes))
	_, err = physicalSeg.CopySlice(0, []uint8(bytes))
	if err != nil {
		fail(fmt.Sprintf("could not copy bytes to segment: %v", err))
	}

	logicalSeg, err := a2enc.Decode(imageType, physicalSeg)
	if err != nil {
		fail(fmt.Sprintf("could not decode image: %v", err))
	}

	if err := logicalSeg.WriteFile(outputPath); err != nil {
		fail(fmt.Sprintf("could not write output file: %v", err))
	}

	fmt.Printf(
		"successfully decoded %s to %s\n", inputPath, outputPath,
	)
}
