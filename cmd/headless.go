package cmd

import (
	"encoding/binary"
	"fmt"
	"image/color"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/pevans/erc/a2"
	"github.com/pevans/erc/a2/a2audio"
	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/record"
	"github.com/spf13/cobra"
)

var (
	headlessStepsFlag        int
	headlessWatchMemFlag     string
	headlessWatchRegFlag     string
	headlessWatchCompFlag    string
	headlessRecordAudioFlag  bool
	headlessCaptureVideoFlag string
	headlessOutputFlag       string
	headlessStartAtFlag      string
)

var headlessCmd = &cobra.Command{
	Use:   "headless [image...]",
	Short: "Run the emulator without a graphical window",
	Long:  "Execute a fixed number of instruction steps and record state, audio, and video without requiring a display",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runHeadless(args)
	},
}

func init() {
	rootCmd.AddCommand(headlessCmd)

	headlessCmd.Flags().IntVar(&headlessStepsFlag, "steps", 0, "Number of instruction steps to execute (required)")
	headlessCmd.MarkFlagRequired("steps") //nolint:errcheck

	headlessCmd.Flags().StringVar(
		&headlessWatchMemFlag,
		"watch-mem",
		"",
		"Comma-separated memory address ranges (e.g. 0400-07FF,2000-3FFF or 013F)",
	)
	headlessCmd.Flags().StringVar(
		&headlessWatchRegFlag,
		"watch-reg",
		"",
		"Comma-separated registers to observe (e.g. A,X,P,PC)",
	)
	headlessCmd.Flags().StringVar(
		&headlessWatchCompFlag,
		"watch-comp",
		"",
		"Comma-separated computer state names to observe (e.g. DisplayHires,BankWriteRAM)",
	)
	headlessCmd.Flags().BoolVar(
		&headlessRecordAudioFlag,
		"record-audio",
		false,
		"Attach an audio recorder",
	)
	headlessCmd.Flags().StringVar(
		&headlessCaptureVideoFlag,
		"capture-video",
		"",
		"Comma-separated step numbers at which to capture video frames",
	)
	headlessCmd.Flags().StringVar(
		&headlessOutputFlag,
		"output",
		".",
		"Directory for output files",
	)
	headlessCmd.Flags().StringVar(
		&headlessStartAtFlag,
		"start-at",
		"",
		"Hex address at which to begin counting steps (e.g. 0801); warm-up runs without recording until PC reaches this address",
	)
}

func runHeadless(images []string) {
	comp := a2.NewComputer(1)

	for _, filename := range images {
		if err := comp.Disks.Append(filename); err != nil {
			fail(fmt.Sprintf("could not open file %s: %v", filename, err))
		}
	}

	if err := comp.LoadFirst(); err != nil {
		fail(fmt.Sprintf("could not load file: %v", err))
	}

	if err := comp.Boot(); err != nil {
		fail(fmt.Sprintf("could not boot emulator: %v", err))
	}

	if headlessStartAtFlag != "" {
		startAddr, err := strconv.ParseUint(headlessStartAtFlag, 16, 16)
		if err != nil {
			fail(fmt.Sprintf("invalid --start-at address %q: %v", headlessStartAtFlag, err))
		}
		const maxWarmup = 10_000_000
		reached := false
		for range maxWarmup {
			if comp.CPU.PC == uint16(startAddr) {
				reached = true
				break
			}
			comp.Process() //nolint:errcheck
		}
		if !reached {
			fail(fmt.Sprintf("--start-at address %04X not reached within %d steps", uint16(startAddr), maxWarmup))
		}
	}

	rec := &record.Recorder{}

	if headlessWatchMemFlag != "" {
		for rangeStr := range strings.SplitSeq(headlessWatchMemFlag, ",") {
			rangeStr = strings.TrimSpace(rangeStr)
			if rangeStr == "" {
				continue
			}
			observers, err := parseHeadlessMemRange(comp, rangeStr)
			if err != nil {
				fail(fmt.Sprintf("invalid --watch-mem range %q: %v", rangeStr, err))
			}
			rec.Add(observers...)
		}
	}

	if headlessWatchRegFlag != "" {
		for regStr := range strings.SplitSeq(headlessWatchRegFlag, ",") {
			regStr = strings.TrimSpace(regStr)
			if regStr == "" {
				continue
			}
			obs, err := headlessRegObserver(comp, regStr)
			if err != nil {
				fail(fmt.Sprintf("invalid --watch-reg register %q: %v", regStr, err))
			}
			rec.Add(obs)
		}
	}

	if headlessWatchCompFlag != "" {
		for stateName := range strings.SplitSeq(headlessWatchCompFlag, ",") {
			stateName = strings.TrimSpace(stateName)
			if stateName == "" {
				continue
			}
			obs, err := headlessCompStateObserver(comp, stateName)
			if err != nil {
				fail(fmt.Sprintf("invalid --watch-comp state %q: %v", stateName, err))
			}
			rec.Add(obs)
		}
	}

	var audioRec *record.AudioRecorder
	if headlessRecordAudioFlag {
		stream := a2audio.NewStream(comp.Speaker(), comp)
		audioRec = record.NewAudioRecorder(
			"audio",
			stream,
			comp.CycleCounter,
			comp.CPUClockRate(),
			a2audio.SampleRate,
		)
		rec.Add(audioRec)
	}

	var videoRec *record.VideoRecorder
	var captureSteps []int
	if headlessCaptureVideoFlag != "" {
		videoRec = record.NewVideoRecorder(comp.Screen)
		for stepStr := range strings.SplitSeq(headlessCaptureVideoFlag, ",") {
			stepStr = strings.TrimSpace(stepStr)
			if stepStr == "" {
				continue
			}
			n, err := strconv.Atoi(stepStr)
			if err != nil {
				fail(fmt.Sprintf("invalid --capture-video step %q: %v", stepStr, err))
			}
			captureSteps = append(captureSteps, n)
		}
		videoRec.CaptureAt(captureSteps...)
	}

	record.Run(
		func() { comp.Process() }, //nolint:errcheck
		rec,
		videoRec,
		headlessStepsFlag,
		comp.Render,
	)

	if err := os.MkdirAll(headlessOutputFlag, 0o755); err != nil {
		fail(fmt.Sprintf("could not create output directory: %v", err))
	}

	if len(rec.Entries()) > 0 {
		path := filepath.Join(headlessOutputFlag, "state.log")
		if err := os.WriteFile(path, []byte(rec.String()), 0o644); err != nil {
			fail(fmt.Sprintf("could not write state.log: %v", err))
		}
	}

	if audioRec != nil {
		samples := audioRec.Samples()
		if len(samples) > 0 {
			path := filepath.Join(headlessOutputFlag, "audio.pcm")
			buf := make([]byte, len(samples)*4)
			for i, s := range samples {
				binary.LittleEndian.PutUint32(buf[i*4:], math.Float32bits(s))
			}
			if err := os.WriteFile(path, buf, 0o644); err != nil {
				fail(fmt.Sprintf("could not write audio.pcm: %v", err))
			}
		}
	}

	if videoRec != nil && len(captureSteps) > 0 {
		content := buildVideoFrameFile(videoRec, captureSteps)
		if content != "" {
			path := filepath.Join(headlessOutputFlag, "video.frame")
			if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
				fail(fmt.Sprintf("could not write video.frame: %v", err))
			}
		}
	}
}

func parseHeadlessMemRange(comp *a2.Computer, rangeStr string) ([]record.Observer, error) {
	parts := strings.SplitN(rangeStr, "-", 2)
	if len(parts) == 1 {
		addr, err := strconv.ParseInt(parts[0], 16, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid address: %w", err)
		}
		return []record.Observer{record.MemObserver(comp.Main, int(addr))}, nil
	}

	start, err := strconv.ParseInt(parts[0], 16, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid start address: %w", err)
	}

	end, err := strconv.ParseInt(parts[1], 16, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid end address: %w", err)
	}

	obs := make([]record.Observer, 0, int(end-start+1))
	for addr := int(start); addr <= int(end); addr++ {
		obs = append(obs, record.MemObserver(comp.Main, addr))
	}

	return obs, nil
}

func headlessRegObserver(comp *a2.Computer, reg string) (record.Observer, error) {
	switch strings.ToUpper(reg) {
	case "A":
		return record.NewObserver(record.TagReg, "A", func() any { return comp.CPU.A }), nil
	case "X":
		return record.NewObserver(record.TagReg, "X", func() any { return comp.CPU.X }), nil
	case "Y":
		return record.NewObserver(record.TagReg, "Y", func() any { return comp.CPU.Y }), nil
	case "P":
		return record.NewObserver(record.TagReg, "P", func() any { return comp.CPU.P }), nil
	case "S":
		return record.NewObserver(record.TagReg, "S", func() any { return comp.CPU.S }), nil
	case "PC":
		return record.NewObserver(record.TagReg, "PC", func() any { return comp.CPU.PC }), nil
	default:
		return nil, fmt.Errorf("unknown register %q", reg)
	}
}

// headlessStateNameToKey maps a2state names to their integer keys for use
// with --watch-comp.
var headlessStateNameToKey = map[string]int{
	"BankDFBlockBank2":    a2state.BankDFBlockBank2,
	"BankROMSegment":      a2state.BankROMSegment,
	"BankReadRAM":         a2state.BankReadRAM,
	"BankSysBlockAux":     a2state.BankSysBlockAux,
	"BankSysBlockSegment": a2state.BankSysBlockSegment,
	"BankWriteRAM":        a2state.BankWriteRAM,
	"DisplayAltChar":      a2state.DisplayAltChar,
	"DisplayAuxSegment":   a2state.DisplayAuxSegment,
	"DisplayCol80":        a2state.DisplayCol80,
	"DisplayDoubleHigh":   a2state.DisplayDoubleHigh,
	"DisplayHires":        a2state.DisplayHires,
	"DisplayIou":          a2state.DisplayIou,
	"DisplayMixed":        a2state.DisplayMixed,
	"DisplayMonochrome":   a2state.DisplayMonochrome,
	"DisplayPage2":        a2state.DisplayPage2,
	"DisplayRedraw":       a2state.DisplayRedraw,
	"DisplayStore80":      a2state.DisplayStore80,
	"DisplayText":         a2state.DisplayText,
	"KBKeyDown":           a2state.KBKeyDown,
	"KBLastKey":           a2state.KBLastKey,
	"KBStrobe":            a2state.KBStrobe,
	"MemAuxSegment":       a2state.MemAuxSegment,
	"MemMainSegment":      a2state.MemMainSegment,
	"MemReadAux":          a2state.MemReadAux,
	"MemReadSegment":      a2state.MemReadSegment,
	"MemWriteAux":         a2state.MemWriteAux,
	"MemWriteSegment":     a2state.MemWriteSegment,
	"Paused":              a2state.Paused,
	"SpeakerState":        a2state.SpeakerState,
}

func headlessCompStateObserver(comp *a2.Computer, name string) (record.Observer, error) {
	key, ok := headlessStateNameToKey[name]
	if !ok {
		return nil, fmt.Errorf("unknown state %q", name)
	}
	return record.NewObserver(record.TagComp, name, func() any {
		return comp.State.Any(key)
	}), nil
}

func buildVideoFrameFile(vr *record.VideoRecorder, steps []int) string {
	var sb strings.Builder
	first := true

	for _, step := range steps {
		frame := vr.Frame(step)
		if frame == nil {
			continue
		}

		if !first {
			sb.WriteByte('\n')
		}

		first = false
		writeFrameText(&sb, step, frame, vr)
	}

	return sb.String()
}

// framePixelGetter abstracts pixel access on a frame returned by
// VideoRecorder.Frame.
type framePixelGetter interface {
	Width() uint
	Height() uint
	GetPixel(x, y uint) color.RGBA
}

func writeFrameText(sb *strings.Builder, step int, frame framePixelGetter, _ *record.VideoRecorder) {
	w := int(frame.Width())
	h := int(frame.Height())

	// Enumerate distinct colors.
	colorToChar := make(map[color.RGBA]byte)
	var nextChar byte = '!'

	for y := range h {
		for x := range w {
			clr := frame.GetPixel(uint(x), uint(y))
			clr.A = 0xff
			if _, ok := colorToChar[clr]; !ok {
				colorToChar[clr] = nextChar
				if nextChar < '~' {
					nextChar++
				}
			}
		}
	}

	// Build sorted legend for deterministic output.
	type colorEntry struct {
		ch  byte
		clr color.RGBA
	}

	entries := make([]colorEntry, 0, len(colorToChar))
	for clr, ch := range colorToChar {
		entries = append(entries, colorEntry{ch, clr})
	}

	sort.Slice(entries, func(i, j int) bool { return entries[i].ch < entries[j].ch })

	legendParts := make([]string, len(entries))
	for i, e := range entries {
		legendParts[i] = fmt.Sprintf("%c=%02X%02X%02X", e.ch, e.clr.R, e.clr.G, e.clr.B)
	}

	fmt.Fprintf(sb, "step %d: video screen %dx%d\n", step, w, h)
	fmt.Fprintf(sb, "colors: %s\n", strings.Join(legendParts, ", "))

	row := make([]byte, w)
	for y := range h {
		for x := range w {
			clr := frame.GetPixel(uint(x), uint(y))
			clr.A = 0xff
			row[x] = colorToChar[clr]
		}
		sb.Write(row)
		sb.WriteByte('\n')
	}
}
