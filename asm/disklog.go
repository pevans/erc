package asm

import (
	"fmt"
	"os"
	"time"
)

const (
	// DiskRead is a mode for a DiskOp that indicates that an operation is
	// some kind of read
	DiskRead = iota

	// DiskWrite is a mode for a DiskOp that indicates some write took place
	DiskWrite
)

// A DiskOp records some disk operation for later analysis (e.g. for
// debugging).
type DiskOp struct {
	// Elapsed is the time between this DiskOp and some user-defined starting
	// point. One example is the time when the emulated computer was booted.
	Elapsed time.Duration

	// Mode is the mode of operation that was performed. By convention, this
	// would be either DiskRead or DiskWrite.
	Mode int

	// Track is the track number where the operation occurred.
	Track int

	// Sector is the sector number where the operation occurred.
	Sector int

	// SectorPosition is the raw offset from the beginning of the track where
	// the operation occurred.
	SectorPosition int

	// Byte is the byte involved in the operation (for example -- what was
	// read, or what was written).
	Byte uint8

	// Instruction is the string representation of the instruction and operand
	// that caused the disk operation
	Instruction string
}

// A DiskLog is a collection of DiskOps. Put together, it would describe a
// series of actions on a given disk.
type DiskLog struct {
	// ops is a slice of the DiskOps that record the actions taken on some
	// disk.
	ops []DiskOp
}

// NewDiskLog returns a newly allocated DiskLog that can be used for recording
// purposes.
func NewDiskLog() *DiskLog {
	return new(DiskLog)
}

// Add appends a given disk operation to the DiskLog.
func (l *DiskLog) Add(op *DiskOp) {
	l.ops = append(l.ops, *op)
}

// WriteToFile will write the contents of the disklog to a file. If this
// cannot be done, an error is returned.
func (l *DiskLog) WriteToFile(filename string) error {
	fp, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer fp.Close() //nolint:errcheck

	for _, op := range l.ops {
		mode := "RD"
		if op.Mode == DiskWrite {
			mode = "WR"
		}

		logLine := fmt.Sprintf(
			"[%-10v] %s T:%02X S:%01X P:%04X B:$%02X | %v\n",
			op.Elapsed.Round(time.Millisecond), // time since boot
			mode,
			op.Track,
			op.Sector,
			op.SectorPosition,
			op.Byte,
			op.Instruction,
		)

		if _, err := fp.WriteString(logLine); err != nil {
			return err
		}
	}

	return nil
}
