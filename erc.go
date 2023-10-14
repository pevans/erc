package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alecthomas/kong"
	"github.com/pevans/erc/a2"
	"github.com/pevans/erc/clog"
	"github.com/pevans/erc/gfx"
	"github.com/pevans/erc/input"

	"github.com/pkg/profile"
)

type cli struct {
	Profile bool   `help:"Write out a profile trace"`
	Image   string `arg`
}

// ConfigFile is the default (relative) location of our configuration file.
const ConfigFile = `.erc/config.toml`

func main() {
	var (
		cli cli
	)

	_ = kong.Parse(&cli)

	if cli.Profile {
		defer profile.Start().Stop()
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	clog.Init(os.Stdout)

	if len(os.Args) < 2 {
		fail("you must pass the name of a file to load")
	}

	comp := a2.NewComputer()
	comp.SetFont(a2.SystemFont())
	gfx.Screen = a2.NewScreen()

	go func() {
		sig := <-signals

		fmt.Printf("Received signal %v", sig)
		comp.Shutdown()
		os.Exit(1)
	}()

	inputFile := cli.Image
	data, err := os.OpenFile(inputFile, os.O_RDWR, 0644)
	if err != nil {
		fail(fmt.Sprintf("could not open file %s: %v", inputFile, err))
	}

	if err := comp.Load(data, inputFile); err != nil {
		fail(fmt.Sprintf("could not load file %s: %v", inputFile, err))
	}

	// Attempt a cold boot
	if err := comp.Boot(); err != nil {
		fail(fmt.Sprintf("could not boot emulator: %v", err))
	}

	// Set up a listener event that funnels through our keyboard handler
	go input.Listen(func(ev input.Event) {
		comp.PressKey(uint8(ev.Key))
	})

	delay := 10 * time.Nanosecond
	go processLoop(comp, delay)

	if err := drawLoop(comp); err != nil {
		fail(fmt.Sprintf("failed to execute draw loop: %v", err))
	}

	// Shutdown
	if err := comp.Shutdown(); err != nil {
		fail(fmt.Sprintf("could not properly shut down emulator: %v", err))
	}
}

func fail(reason string) {
	fmt.Println(reason)
	os.Exit(1)
}
