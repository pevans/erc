package a2

import (
	"fmt"
	"sync"
	"time"

	"github.com/pevans/erc/a2/a2drive"
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

	ShutdownMutex sync.Mutex
	WillShutDown  bool

	ClockEmulator *clock.Emulator

	// speed is some abstract number to indicate how fast we're going. Refer
	// to ClockSpeed to see how this number is used to set our actual hertz.
	speed int

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

	Drive1        *a2drive.Drive
	Drive2        *a2drive.Drive
	SelectedDrive *a2drive.Drive

	// diskLogFileName is the name we'll use to write the diskLog
	diskLogFileName string

	// diskLog is the log of all disk operations recorded by the file in
	// drive1
	diskLog *asm.DiskLog

	// When the computer is booted up, this will be a set of disks that we
	// might use to run software. There are often cases where you need to swap
	// disks, but we constrain that to a small set of disks that is knowable
	// at boot-time.
	Disks *DiskSet

	smap  *memory.SoftMap
	State *memory.StateMap

	InstructionLog *asm.InstructionMap

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

// ClockSpeed returns hertz based on the given abstract speed. Relatively
// larger speeds imply a larger hertz; i.e. ClockSpeed(2) > ClockSpeed(1).
func ClockSpeed(speed int) int64 {
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

	comp.Drive1 = a2drive.NewDrive()
	comp.Drive2 = a2drive.NewDrive()
	comp.SelectedDrive = comp.Drive1

	comp.Disks = NewDiskSet()

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
	comp.speed = 1

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

// NeedsRender returns true if we think the screen needs to be redrawn.
func (c *Computer) NeedsRender() bool {
	return c.State.Bool(a2state.DisplayRedraw)
}

// SetSpeed changes the rate of emulation to a precise speed number (with
// hertz determined by ClockSpeed()). This can be dangerous; if you set a
// speed of 0, the emulated computer will halt.
func (c *Computer) SetSpeed(n int) {
	c.ShowText(fmt.Sprintf("speed: %v", n))
	c.speed = n
	c.ClockEmulator.ChangeHertz(ClockSpeed(n))
}

// ShowText flashes a text message on the screen using the gfx package's
// TextNotification
func (c *Computer) ShowText(message string) {
	w, h := c.Dimensions()
	gfx.TextNotification.Show(message, int(w), int(h))
}

// SpeedUp changes the rate of emulation by increasing the relative speed to
// one number greater than it is right now. If you're at speed 1, it becomes
// speed 2; etc. This method will not set the rate of speed higher than 5.
func (c *Computer) SpeedUp() {
	// We can't go higher than ~5mhz right now. (We could, of course, but...
	// we don't want to.)
	if c.speed > 4 {
		c.ShowText(fmt.Sprintf("speed: %v", c.speed))
		return
	}

	c.SetSpeed(c.speed + 1)
}

// SpeedDown changes the rate of emulation by decreasing the relative speed to
// one number fewer than it is right now. If you're at speed 3, it becomes
// speed 2. This method will not set a speed below 1.
func (c *Computer) SpeedDown() {
	// We can't go lower than 1mhz. (Well, we could, but we don't do this
	// right now.)
	if c.speed < 2 {
		c.ShowText(fmt.Sprintf("speed: %v", c.speed))
		return
	}

	c.SetSpeed(c.speed - 1)
}

// StateMap returns the computer's available state map.
func (c *Computer) StateMap() *memory.StateMap {
	return c.State
}
