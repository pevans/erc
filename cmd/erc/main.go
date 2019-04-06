package main

import (
	"fmt"
	"os"

	"github.com/pevans/erc/pkg/mach/a2"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func main() {
	var (
		logFile  = os.Getenv("LOG_FILE")
		logLevel = os.Getenv("LOG_LEVEL")
	)

	if err := setLogging(logFile, logLevel); err != nil {
		fmt.Printf("unable to set logging: %v", err)
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		fmt.Printf("you must pass the name of a file to load")
		os.Exit(1)
	}

	// Right now, there's only one machine to emulate.
	emu := a2.NewEmulator()

	// At this stage, we need to decide what we should be loading.
	if err := emu.Loader.Load(os.Args[1]); err != nil {
		fmt.Println(errors.Wrapf(err, "could not load file %s", os.Args[1]))
		os.Exit(1)
	}

	// Attempt a cold boot
	if err := emu.Booter.Boot(); err != nil {
		fmt.Println(errors.Wrapf(err, "could not boot emulator"))
		os.Exit(1)
	}

	// This sets up our processor loop
	for {
		if err := emu.Processor.Process(); err != nil {
			log.Error(errors.Wrapf(err, "received error from processor"))
			break
		}
	}

	// Shutdown
	if err := emu.Ender.End(); err != nil {
		fmt.Println(errors.Wrapf(err, "could not properly shut down emulator"))
		os.Exit(1)
	}
}

// setLogging attempts to set the file name and level for logging within
// logrus, and returns an error if it was unable to do so.
func setLogging(fileName, levelName string) error {
	// No logging is necessary!
	if fileName == "" {
		return nil
	}

	// If a file was given, but no level, then assume they want error
	// logging
	if levelName == "" {
		levelName = "error"
	}

	file, err := os.OpenFile(fileName, os.O_WRONLY, 0755)
	if err != nil {
		return errors.Wrapf(err, "could not open file %s for logging", fileName)
	}

	level, err := log.ParseLevel(levelName)
	if err != nil {
		return errors.Wrapf(err, "could not recognize level %s for logging", levelName)
	}

	log.SetOutput(file)
	log.SetLevel(level)

	return nil
}
