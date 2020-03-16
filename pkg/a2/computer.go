package a2

import (
	"io"

	"github.com/pevans/erc/pkg/asmrec"
	"github.com/pevans/erc/pkg/asmrec/a2rec"
	"github.com/pevans/erc/pkg/boot"
	"github.com/pevans/erc/pkg/data"
	"github.com/pevans/erc/pkg/emu"
	"github.com/pevans/erc/pkg/proc/mos65c02"
)

// ReadMapFn is a function which can execute a soft switch procedure on
// read.
type ReadMapFn func(*Computer, data.Addressor) data.Byte

// WriteMapFn is a function which can execute a soft switch procedure on
// write.
type WriteMapFn func(*Computer, data.Addressor, data.Byte)

// A Computer is our abstraction of an Apple //e ("enhanced") computer.
type Computer struct {
	drawer emu.Drawer

	// The CPU of the Apple //e was an MOS 65C02 processor.
	CPU *mos65c02.CPU

	// reDraw is set to true when a screen redraw is necessary, and set
	// to false the redraw is done.
	reDraw bool

	// There are three primary segments of memory in an Apple //e; main
	// memory, read-only memory, and auxiliary memory. Each are
	// accessible through a mechanism called bank-switching.
	Main *data.Segment
	ROM  *data.Segment
	Aux  *data.Segment

	Drive1        *Drive
	Drive2        *Drive
	SelectedDrive *Drive

	// RMap and WMap are the read and write address maps. These contain
	// functions which emulate the "soft switches" that Apple IIs used
	// to implement special functionality.
	RMap map[int]ReadMapFn
	WMap map[int]WriteMapFn

	// MemMode is a collection of bit flags which tell us what state of
	// memory we have.
	MemMode int

	// BankMode is the set of bit flags which are the memory banks that
	// we are accessing right now.
	BankMode int

	// PCMode is the peripheral card mode we have for memory, which
	// governs the range of $C100 - $CFFF.
	PCMode int

	// DisplayMode is the state that our display output is currently in.
	// (For example, text mode, hires, lores, etc.)
	DisplayMode int

	log *boot.Logger

	recWriter io.Writer
	rec       asmrec.Recorder
}

const (
	// AuxMemorySize is the length of memory for auxiliary memory in the
	// Apple II, which was implemented through a peripheral called a
	// "language card" installed in the back. Bank-switches let you swap
	// in and out auxiliary memory for main memory. Note that auxiliary
	// memory is only 64k bytes large.
	AuxMemorySize = 0x10000

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

// NewComputer returns an Apple //e computer value, which essentially
// encompasses all of the things that an Apple II would need to run.
func NewComputer() *Computer {
	comp := &Computer{}

	comp.Aux = data.NewSegment(AuxMemorySize)
	comp.Main = data.NewSegment(MainMemorySize)
	comp.ROM = data.NewSegment(RomMemorySize)

	comp.Drive1 = NewDrive()
	comp.Drive2 = NewDrive()
	comp.SelectedDrive = comp.Drive1

	comp.CPU = new(mos65c02.CPU)
	comp.CPU.WMem = comp
	comp.CPU.RMem = comp

	comp.RMap = make(map[int]ReadMapFn)
	comp.WMap = make(map[int]WriteMapFn)

	return comp
}

func (c *Computer) SetLogger(l *boot.Logger) {
	c.log = l
}

func (c *Computer) SetRecorderWriter(w io.Writer) {
	c.recWriter = w
	c.rec = new(a2rec.Recorder)
}

// Dimensions returns the screen dimensions of an Apple II.
func (c *Computer) Dimensions() (width, height int) {
	return 280, 192
}
