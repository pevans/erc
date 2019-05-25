package main

import (
	"fmt"
	"os"

	"github.com/hajimehoshi/ebiten"
	"github.com/pevans/erc/pkg/mach/a2"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	ConfigFile = `.erc/config.toml`
)

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

	var (
		width, height = emu.Drawer.Dimensions()

		loop = func(screen *ebiten.Image) error {
			if err := emu.Processor.Process(); err != nil {
				log.Error(errors.Wrapf(err, "received error from processor"))
			}

			return nil
		}
	)

	if err := ebiten.Run(loop, width, height, 3, "erc"); err != nil {
		fmt.Println(errors.Wrap(err, "run loop failed"))
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
