package main

import (
	"fmt"
	"os"

	"github.com/pevans/erc/pkg/boot"
	"github.com/pevans/erc/pkg/mach"
	"github.com/pevans/erc/pkg/mach/a2"
	"github.com/pkg/errors"
)

// ConfigFile is the default (relative) location of our configuration file.
const ConfigFile = `.erc/config.toml`

var emulator *mach.Emulator

func main() {
	var (
		homeDir     = os.Getenv("HOME")
		configFile  = fmt.Sprintf("%s/%s", homeDir, ConfigFile)
		instLogFile *os.File
	)

	// Let's see if we can figure out our config situation
	conf, err := boot.NewConfig(configFile)
	if err != nil {
		fmt.Println(errors.Wrapf(err, "unable to read config file %s", configFile))
		os.Exit(1)
	}

	log, err := conf.NewLogger()
	if err != nil {
		fmt.Println(errors.Wrap(err, "unable to create logger"))
		os.Exit(1)
	}

	log.UseOutput()

	if len(os.Args) < 2 {
		fmt.Println("you must pass the name of a file to load")
		os.Exit(1)
	}

	if conf.InstructionLog.File != "" {
		instLogFile, err = openLogFile(conf.InstructionLog.File)
		if err != nil {
			fmt.Printf("unable to open file for instruction logging: %v", err)
		}
	}

	// Right now, there's only one machine to emulate.
	emulator = a2.NewEmulator(instLogFile)

	// At this stage, we need to decide what we should be loading.
	if err := emulator.Loader.Load(os.Args[1]); err != nil {
		fmt.Println(errors.Wrapf(err, "could not load file %s", os.Args[1]))
		os.Exit(1)
	}

	// Attempt a cold boot
	if err := emulator.Booter.Boot(); err != nil {
		fmt.Println(errors.Wrapf(err, "could not boot emulator"))
		os.Exit(1)
	}

	go processorLoop()

	if err := gameLoop(); err != nil {
		fmt.Println(errors.Wrap(err, "run loop failed"))
	}

	// Shutdown
	if err := emulator.Ender.End(); err != nil {
		fmt.Println(errors.Wrapf(err, "could not properly shut down emulator"))
		os.Exit(1)
	}
}
