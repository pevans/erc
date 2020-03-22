package main

import (
	"fmt"
	"os"

	"github.com/pevans/erc/pkg/a2"
	"github.com/pevans/erc/pkg/boot"

	"github.com/pkg/errors"
)

// ConfigFile is the default (relative) location of our configuration file.
const ConfigFile = `.erc/config.toml`

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
		log.Error("you must pass the name of a file to load")
		os.Exit(1)
	}

	inputFile := os.Args[1]

	if conf.InstructionLog.File != "" {
		instLogFile, err = boot.OpenFile(conf.InstructionLog.File)
		if err != nil {
			log.Errorf("unable to open file for instruction logging: %v", err)
		}
	}

	comp := a2.NewComputer()
	comp.SetLogger(log)
	comp.SetRecorderWriter(instLogFile)

	data, err := boot.OpenFile(inputFile)
	if err != nil {
		log.Error(errors.Wrapf(err, "could not open file %s", inputFile))
		os.Exit(1)
	}

	if err := comp.Load(data, inputFile); err != nil {
		log.Error(errors.Wrapf(err, "could not load file %s", inputFile))
		os.Exit(1)
	}

	// Attempt a cold boot
	if err := comp.Boot(); err != nil {
		log.Error(errors.Wrapf(err, "could not boot emulator"))
		os.Exit(1)
	}

	// In another goroutine, execute the process loop
	go comp.ProcessLoop()

	// And in the main thread, execute the draw loop
	comp.DrawLoop()

	// Shutdown
	if err := comp.Shutdown(); err != nil {
		log.Error(errors.Wrapf(err, "could not properly shut down emulator"))
		os.Exit(1)
	}
}
