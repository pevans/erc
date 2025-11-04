package a2

import (
	"time"

	"github.com/pevans/erc/a2/a2font"
	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/asm"
	"github.com/pevans/erc/clock"
	"github.com/pevans/erc/gfx"
	"github.com/pevans/erc/memory"
	"github.com/pevans/erc/mos"
)

// ReadMapFn is a function which can execute a soft switch procedure on
// read.
type ReadMapFn func(*Computer, int) uint8

// WriteMapFn is a function which can execute a soft switch procedure on
// write.
type WriteMapFn func(*Computer, int, uint8)

// A Computer is our abstraction of an Apple //e ("enhanced") computer.
type Computer struct {
	// The CPU of the Apple //e was an MOS 65C02 processor.
	CPU *mos.CPU

	ClockEmulator *clock.Emulator

	// When did the computer boot? This also includes when the computer is
	// soft-booted (i.e. reset with a blank context after having been running
	// for a while).
	BootTime time.Time

	// There are three primary segments of memory in an Apple //e; main
	// memory, read-only memory, and auxiliary memory. Each are
	// accessible through a mechanism called bank-switching.
	Main *memory.Segment
	ROM  *memory.Segment
	Aux  *memory.Segment

	Screen *gfx.FrameBuffer

	Drive1        *Drive
	Drive2        *Drive
	SelectedDrive *Drive
	diskLog       *asm.DiskLog

	smap  *memory.SoftMap
	State *memory.StateMap

	InstructionLog *asm.CallMap

	// Where to write the instruction log
	InstructionLogFileName string

	TimeSet         *asm.TimeSet
	TimeSetFileName string

	MetricsFileName string

	// MemMode is a collection of bit flags which tell us what state of
	// memory we have.
	MemMode int

	// DisplayMode is the state that our display output is currently in.
	// (For example, text mode, hires, lores, etc.)
	DisplayMode int

	// We use two different fonts for 40-column and 80-column text
	Font40 *gfx.Font
	Font80 *gfx.Font
}

const (
	// AuxMemorySize is the length of memory for auxiliary memory in the
	// Apple II, which was implemented through a peripheral called a
	// "language card" installed in the back. Bank-switches let you swap
	// in and out auxiliary memory for main memory. Note that auxiliary
	// memory is only 64k bytes large.
	AuxMemorySize = 0x11000

	// MainMemorySize is the length of memory for so-called "main
	// memory" in an Apple II. It consists of 68k of RAM; although only
	// 64k is addressible at a time, the last 4k can be accessed via
	// bank-switches.
	MainMemorySize = 0x11000

	// RomMemorySize is the length of system read-only memory.
	RomMemorySize = 0x5000

	// SysRomOffset is the spot in memory where system ROM can be found.
	SysRomOffset = 0xC000
)

const (
	// This width and height is actually double the size of a typical Apple II
	// display. We're doing this so that it's easier for us to manage
	// 80 column text and double-high resolution.
	screenWidth  uint = 560
	screenHeight uint = 384
)

// NewComputer returns an Apple //e computer value, which essentially
// encompasses all of the things that an Apple II would need to run.
func NewComputer(hertz int64) *Computer {
	comp := &Computer{}

	comp.Aux = memory.NewSegment(AuxMemorySize)
	comp.Main = memory.NewSegment(MainMemorySize)
	comp.ROM = memory.NewSegment(RomMemorySize)
	comp.smap = memory.NewSoftMap(0x20000)
	comp.State = memory.NewStateMap()
	comp.smap.UseState(comp.State)

	comp.Aux.UseSoftMap(comp.smap)
	comp.Main.UseSoftMap(comp.smap)
	comp.ROM.UseSoftMap(comp.smap)

	comp.Drive1 = NewDrive()
	comp.Drive2 = NewDrive()
	comp.SelectedDrive = comp.Drive1

	comp.CPU = new(mos.CPU)
	comp.CPU.RMem = comp
	comp.CPU.WMem = comp
	comp.CPU.State = comp.State

	// Note that hertz is treated as a unit of cycles per second, but
	// the number may not feel precisely accurate to how an Apple II
	// might have run. I've found that if I use 1.023 MHz, the Apple IIe
	// speed, things feel much slower than I'd expect. In practice,
	// something approximately double that number feels more right.
	comp.ClockEmulator = clock.NewEmulator(hertz)

	comp.Font40 = a2font.SystemFont40()
	comp.Font80 = a2font.SystemFont80()

	comp.Screen = gfx.NewFrameBuffer(screenWidth, screenHeight)
	gfx.Screen = comp.Screen

	return comp
}

// Dimensions returns the screen dimensions of an Apple II.
func (c *Computer) Dimensions() (width, height uint) {
	return screenWidth, screenHeight
}

func (c *Computer) NeedsRender() bool {
	return c.State.Bool(a2state.DisplayRedraw)
}

func (c *Computer) StateMap() *memory.StateMap {
	return c.State
}
