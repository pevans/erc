package main

import (
	"fmt"
	"os"

	"github.com/pevans/erc/pkg/mach"
	"github.com/pevans/erc/pkg/mach/a2"
	"github.com/pkg/errors"
)

// ConfigFile is the default (relative) location of our configuration file.
const ConfigFile = `.erc/config.toml`

var emulator *mach.Emulator

func main() {
	var (
		homeDir    = os.Getenv("HOME")
		configFile = fmt.Sprintf("%s/%s", homeDir, ConfigFile)
	)

	// Let's see if we can figure out our config situation
	conf, err := NewConfig(configFile)
	if err != nil {
		fmt.Println(errors.Wrapf(err, "unable to read config file %s", configFile))
		os.Exit(1)
	}

	// And, if we need to be logging, where that goes
	if err := setLogging(conf.Log.File, conf.Log.Level); err != nil {
		fmt.Printf("unable to set logging: %v", err)
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		fmt.Println("you must pass the name of a file to load")
		os.Exit(1)
	}

	// Right now, there's only one machine to emulate.
	emulator = a2.NewEmulator()

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
