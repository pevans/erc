package main

import (
	"fmt"
	"os"
	"time"

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
		log.Fatal("you must pass the name of a file to load")
	}

	comp := a2.NewComputer()
	comp.SetFont(a2.SystemFont())
	comp.SetLogger(log)

	if conf.InstructionLog.File != "" {
		instLogFile, err = os.OpenFile(
			conf.InstructionLog.File,
			os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
			0755,
		)
		if err != nil {
			log.Errorf("unable to open file for instruction logging: %v", err)
		}

		comp.SetRecorderWriter(instLogFile)
	}

	inputFile := os.Args[1]
	data, err := os.OpenFile(inputFile, os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "could not open file %s", inputFile))
	}

	if err := comp.Load(data, inputFile); err != nil {
		log.Fatal(errors.Wrapf(err, "could not load file %s", inputFile))
	}

	// Attempt a cold boot
	if err := comp.Boot(); err != nil {
		log.Fatal(errors.Wrapf(err, "could not boot emulator"))
	}

	delay := 10 * time.Nanosecond
	go processLoop(comp, log, delay)

	if err := drawLoop(comp, log); err != nil {
		log.Fatal(errors.Wrap(err, "failed to execute draw loop"))
	}

	// Shutdown
	if err := comp.Shutdown(); err != nil {
		log.Fatal(errors.Wrapf(err, "could not properly shut down emulator"))
	}
}
