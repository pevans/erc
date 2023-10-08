package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alecthomas/kong"
	"github.com/pevans/erc/a2"
	"github.com/pevans/erc/asmrec"
	"github.com/pevans/erc/clog"
	"github.com/pevans/erc/disasm"
	"github.com/pevans/erc/gfx"
	"github.com/pevans/erc/input"

	"github.com/pkg/errors"
	"github.com/pkg/profile"
)

type cli struct {
	ExecTrace   string `help:"Write an execution trace to a file"`
	Disassembly string `help:"Write disassembled instructions to a file"`
	Profile     bool   `help:"Write out a profile trace"`
	Image       string `arg`
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
		clog.Error("you must pass the name of a file to load")
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

	if cli.ExecTrace != "" {
		instLogFile, err := os.OpenFile(
			cli.ExecTrace, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755,
		)
		if err != nil {
			clog.Errorf("unable to open file for instruction logging: %v", err)
		}

		asmrec.Init(instLogFile)
	}

	if cli.Disassembly != "" {
		disFile, err := os.OpenFile(
			cli.Disassembly, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755,
		)
		if err != nil {
			clog.Errorf("unable to open file for disassembly: %v", err)
		}

		disasm.Init(disFile)
	}

	inputFile := cli.Image
	data, err := os.OpenFile(inputFile, os.O_RDWR, 0644)
	if err != nil {
		clog.Error(errors.Wrapf(err, "could not open file %s", inputFile))
	}

	if err := comp.Load(data, inputFile); err != nil {
		clog.Error(errors.Wrapf(err, "could not load file %s", inputFile))
	}

	// Attempt a cold boot
	if err := comp.Boot(cli.Disassembly); err != nil {
		clog.Error(errors.Wrapf(err, "could not boot emulator"))
	}

	// Set up a listener event that funnels through our keyboard handler
	go input.Listen(func(ev input.Event) {
		comp.PressKey(uint8(ev.Key))
	})

	delay := 10 * time.Nanosecond
	go processLoop(comp, delay)

	if err := drawLoop(comp); err != nil {
		clog.Error(errors.Wrap(err, "failed to execute draw loop"))
	}

	// Shutdown
	if err := comp.Shutdown(); err != nil {
		clog.Error(errors.Wrapf(err, "could not properly shut down emulator"))
	}
}
