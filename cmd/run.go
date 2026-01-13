package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/peterh/liner"
	"github.com/pevans/erc/a2"
	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/a2/a2video"
	"github.com/pevans/erc/debug"
	"github.com/pevans/erc/input"
	"github.com/pevans/erc/render"
	"github.com/pevans/erc/shortcut"
	"github.com/pkg/profile"
	"github.com/spf13/cobra"
)

var (
	debugImageFlag   bool
	debugBreakFlag   string
	profileFlag      bool
	speedFlag        int
	writeProtectFlag bool
	shaderFlag       string
	monochromeFlag   string
)

var runCmd = &cobra.Command{
	Use:   "run [image...]",
	Short: "Emulate a disk image",
	Long:  "Emulate an Apple //e computer and boot with the specified disk image file(s)",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runEmulator(args)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().BoolVar(&debugImageFlag, "debug-image", false, "Write out debugging files to debug image loading")
	runCmd.Flags().StringVar(&debugBreakFlag, "debug-break", "", "Set breakpoints for a comma-separated list of addresses (eg 3FC8,9D94)")
	runCmd.Flags().BoolVar(&profileFlag, "profile", false, "Write out a profile trace")
	runCmd.Flags().IntVar(&speedFlag, "speed", 1, "Starting speed of the emulator (more is faster)")
	runCmd.Flags().BoolVar(&writeProtectFlag, "write-protect", false, "Whether to write-protect the image")
	runCmd.Flags().StringVar(&shaderFlag, "shader", "softcrt", "Shader to apply (none, softcrt, curvedcrt, hardcrt)")
	runCmd.Flags().StringVar(&monochromeFlag, "monochrome", "", "Render in monochrome (green or amber)")
}

func runEmulator(images []string) {
	if profileFlag {
		defer profile.Start().Stop()
	}

	// Parse monochrome flag
	monochromeMode := a2video.MonochromeNone
	switch monochromeFlag {
	case "green":
		monochromeMode = a2video.MonochromeGreen
	case "amber":
		monochromeMode = a2video.MonochromeAmber
	case "":
		monochromeMode = a2video.MonochromeNone
	default:
		fail("monochrome flag must be either 'green' or 'amber'")
	}

	// Parse and add breakpoints if provided
	if debugBreakFlag != "" {
		for addrStr := range strings.SplitSeq(debugBreakFlag, ",") {
			addrStr = strings.TrimSpace(addrStr)
			addr, err := strconv.ParseInt(addrStr, 16, 17)
			if err != nil {
				fail(fmt.Sprintf("invalid breakpoint address '%s': %v", addrStr, err))
			}
			debug.AddBreakpoint(int(addr))
		}
	}

	// Build the computer and screen objects
	comp := a2.NewComputer(speedFlag)

	// Set up a signal handler for graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals

		fmt.Printf("Received signal %v\n", sig)
		err := comp.Shutdown()
		if err != nil {
			fail(fmt.Sprintf("shutdown was not clean: %v", err))
		}

		os.Exit(1)
	}()

	// Load the image files
	comp.State.SetBool(a2state.DebugImage, debugImageFlag)
	comp.State.SetInt(a2state.DisplayMonochrome, monochromeMode)

	for _, filename := range images {
		if err := comp.Disks.Append(filename); err != nil {
			fail(fmt.Sprintf("could not open file %s: %v", filename, err))
		}
	}

	if err := comp.LoadFirst(); err != nil {
		fail(fmt.Sprintf("could not load file %s: %v", images[0], err))
	}

	if writeProtectFlag {
		comp.Drive1.SetWriteProtect(true)
	}

	if err := comp.Boot(); err != nil {
		fail(fmt.Sprintf("could not boot emulator: %v", err))
	}

	// Set up keyboard input handler
	go input.Listen(func(ev input.Event) {
		found, err := shortcut.Check(ev, comp)
		if err != nil {
			fail(fmt.Sprintf("shortcut failed: %v", err))
		}

		// We found a shortcut, and that did something, so don't register this
		// as a keypress
		if found {
			// What if we tried to quit...?
			if comp.WillShutDown {
				os.Exit(0)
			}

			return
		}

		key := ev.Key
		if ev.Modifier == input.ModControl {
			// We want the Apple software to see this as a literal
			// control character
			key = ev.Key & 0x1F
		}

		comp.PressKey(uint8(key))
	})

	// Set up the process loop with debugger
	line := liner.NewLiner()
	defer line.Close() //nolint:errcheck

	emulator := comp.ClockEmulator

	emulator.SetDebuggerEntry(func() {
		debug.Prompt(comp, line)
	})

	emulator.SetBreakpointCheck(func() {
		if debug.HasBreakpoint(int(comp.CPU.PC)) {
			comp.State.SetBool(a2state.Debugger, true)
		}
	})

	go emulator.ProcessLoop(comp)

	// Run the draw loop in the main thread
	if err := render.DrawLoop(comp, shaderFlag); err != nil {
		fail(fmt.Sprintf("failed to execute draw loop: %v", err))
	}

	// Shutdown if we exit the draw loop
	if err := comp.Shutdown(); err != nil {
		fail(fmt.Sprintf("could not properly shut down emulator: %v", err))
	}
}

func fail(reason string) {
	fmt.Fprintln(os.Stderr, reason)
	os.Exit(1)
}

// clockspeed returns hertz based on the given abstract speed.
// Relatively larger speeds imply a larger hertz; i.e. clockspeed(2) > clockspeed(1).
func clockspeed(speed int) int64 {
	// Use the basic clockspeed of an Apple IIe as a starting point
	hertz := int64(1_023_000)

	// Don't allow the caller to get too crazy
	if speed > 5 {
		speed = 5
	}

	for i := 1; i < speed; i++ {
		hertz *= 2
	}

	return hertz
}
