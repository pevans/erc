package cmd

import (
	"fmt"
	"os"

	"github.com/pevans/erc/a2/a2disk"
	"github.com/pevans/erc/a2/a2enc"
	"github.com/pevans/erc/memory"
	"github.com/spf13/cobra"
)

var disCmd = &cobra.Command{
	Use:   "dis [image] [output]",
	Short: "Disassemble a disk image",
	Long:  "Disassemble a disk image and write the assembly to a file",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		disassemble(args[0], args[1])
	},
}

func init() {
	rootCmd.AddCommand(disCmd)
}

func disassemble(imagePath, outputPath string) {
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

	img := a2disk.NewImage()
	if err := img.Parse(seg); err != nil {
		fail("could not parse disk image")
	}

	if err := img.Disassemble(); err != nil {
		fail(fmt.Sprintf("could not disassemble disk image: %v", err))
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		fail(fmt.Sprintf("could not create output file %s: %v", outputPath, err))
	}
	defer outFile.Close() //nolint:errcheck

	for _, line := range img.Code {
		if _, err := fmt.Fprintln(outFile, line.String()); err != nil {
			fail(fmt.Sprintf("could not write to output file: %v", err))
		}
	}

	fmt.Printf("Disassembled %d instructions to %s\n", len(img.Code), outputPath)
}
