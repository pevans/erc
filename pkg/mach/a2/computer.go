// Package a2 defines all of the code necessary to describe an Apple
// II computer. Specifically, it builds what would qualify as an Apple
// //e computer, which is the Apple IIe "enhanced" model. (Contrast with
// the original Apple ][, Apple ][+ ("plus"), Apple ][e, and Apple //c.)
//
// Ideally, you don't need to directly call any of this code in order to
// emulator the machine aside from the NewEmulator() function. An
// Emulator is a very abstract type which abstracts all of the "things"
// a computer (that we can emulate) can do.
package a2

import (
	"github.com/pevans/erc/pkg/mach"
	"github.com/pevans/erc/pkg/proc/mos65c02"
)

// A Computer is our abstraction of an Apple //e ("enhanced") computer.
type Computer struct {
	// The CPU of the Apple //e was an MOS 65C02 processor.
	CPU *mos65c02.CPU

	// There are three primary segments of memory in an Apple //e; main
	// memory, read-only memory, and auxiliary memory. Each are
	// accessible through a mechanism called bank-switching.
	Main *mach.Segment
	ROM  *mach.Segment
	Aux  *mach.Segment

	// MemMode is a collection of bit flags which tell us what state of
	// memory we have.
	MemMode int
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
)

// NewComputer returns an Apple //e computer value, which essentially
// encompasses all of the things that an Apple II would need to run.
func NewComputer() *Computer {
	comp := &Computer{}

	comp.Aux = mach.NewSegment(AuxMemorySize)
	comp.Main = mach.NewSegment(MainMemorySize)
	comp.ROM = mach.NewSegment(RomMemorySize)

	comp.CPU = new(mos65c02.CPU)
	comp.CPU.WSeg = comp.Main
	comp.CPU.RSeg = comp.Main

	return comp
}
