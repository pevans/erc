package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kong"
	"github.com/pevans/erc/a2"
	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/gfx"
	"github.com/pevans/erc/input"
	"github.com/pevans/erc/shortcut"

	"github.com/pkg/profile"
)

type cli struct {
	Profile    bool   `help:"Write out a profile trace"`
	DebugImage bool   `help:"Write out debugging files to debug image loading"`
	Image      string `arg:""`
	Speed      int    `default:"1" help:"Starting speed of the emulator (more is faster)"`
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

	if len(os.Args) < 2 {
		fail("you must pass the name of a file to load")
	}

	comp := a2.NewComputer(clockspeed(cli.Speed))
	gfx.Screen = a2.NewScreen()

	go func() {
		sig := <-signals

		fmt.Printf("Received signal %v", sig)
		err := comp.Shutdown()
		if err != nil {
			fail(fmt.Sprintf("shutdown was not clean: %v", err))
		}

		os.Exit(1)
	}()

	comp.State.SetBool(a2state.DebugImage, cli.DebugImage)

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
		if shortcut.Check(ev, comp) {
			return
		}

		comp.PressKey(uint8(ev.Key))
	})

	go processLoop(comp)

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

// Return hertz based on some given abstract speed. Relatively larger
// speeds imply a larger hertz; i.e. clockspeed(2) > clockspeed(1).
func clockspeed(speed int) int64 {
	// Use the basic clockspeed of an Apple IIe as a starting point
	hertz := int64(1_023_000)

	// Let's not allow the caller to get too crazy
	if speed > 5 {
		speed = 5
	}

	for i := 1; i < speed; i++ {
		hertz *= 2
	}

	return hertz
}
