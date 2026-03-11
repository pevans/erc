package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/pevans/erc/assembler"
)

func main() {
	var outputPath string

	flag.StringVar(&outputPath, "o", "", "Output .dsk file path (omit or - for stdout)")
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: erc-assembler [-o output.dsk] input.s")
		os.Exit(1)
	}

	inputPath := args[0]

	var (
		src      []byte
		err      error
		filename string
	)

	if inputPath == "-" {
		src, err = io.ReadAll(os.Stdin)
		filename = "<stdin>"
	} else {
		src, err = os.ReadFile(inputPath)
		filename = inputPath
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "could not read input: %v\n", err)
		os.Exit(1)
	}

	image, err := assembler.Assemble(src, filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	if outputPath == "" || outputPath == "-" {
		_, err = os.Stdout.Write(image)
	} else {
		err = os.WriteFile(outputPath, image, 0o644)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "could not write output: %v\n", err)
		os.Exit(1)
	}
}
