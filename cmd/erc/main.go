package main

import (
	"log"
	"os"

	"github.com/pevans/erc/pkg/mach/a2"
	"github.com/pkg/errors"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("You must pass the name of a file to load")
	}

	// Right now, there's only one machine to emulate.
	emu := a2.NewEmulator()

	// At this stage, we need to decide what we should be loading.
	if err := emu.Loader.Load(os.Args[1]); err != nil {
		log.Fatal(errors.Wrapf(err, "could not load file %s", os.Args[1]))
	}

	// Attempt a cold boot
	if err := emu.Booter.Boot(); err != nil {
		log.Fatal(errors.Wrapf(err, "could not boot emulator"))
	}

	// This sets up our processor loop
	for {
		if err := emu.Processor.Process(); err != nil {
			log.Println(errors.Wrapf(err, "received error from processor"))
			break
		}
	}

	// Shutdown
	if err := emu.Ender.End(); err != nil {
		log.Fatal(errors.Wrapf(err, "could not properly shut down emulator"))
	}
}
