package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alecthomas/kong"
	"github.com/pevans/erc/a2"
	"github.com/pevans/erc/gfx"
	"github.com/pevans/erc/input"
	"github.com/pevans/erc/shortcut"

	"github.com/pkg/profile"
)

type cli struct {
	Profile bool   `help:"Write out a profile trace"`
	Image   string `arg`
	Speed   int    `default:"1" help:"Starting speed of the emulator (more is faster)"`
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
		if shortcut.Check(ev, comp) {
			return
		}

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

// Return hertz based on some given abstract speed. A larger speed
// should return a larger hertz.
func clockspeed(speed int) int64 {
	// You should not consider this number to correlate with how fast an
	// Apple II might have run.
	hertz := int64(2_000_000)

	// Let's not allow the caller to get too crazy
	if speed > 5 {
		speed = 5
	}

	for i := 1; i < speed; i++ {
		hertz *= 2
	}

	return hertz
}
