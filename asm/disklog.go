package asm

import (
	"fmt"
	"os"
	"time"
)

const (
	DiskRead = iota
	DiskWrite
)

type DiskOp struct {
	// The time between reads
	Elapsed time.Duration

	// What kind of operation is this (DiskRead, DiskWrite)
	Mode int

	// Where the read occurred on the disk
	Track          int
	Sector         int
	SectorPosition int

	// What byte was read
	Byte uint8

	// The full instruction and operand that caused the disk read
	Instruction string
}

type DiskLog struct {
	Ops  []DiskOp
	Name string
}

func NewDiskLog(name string) *DiskLog {
	log := new(DiskLog)
	log.Name = name

	return log
}

func (l *DiskLog) Add(read *DiskOp) {
	l.Ops = append(l.Ops, *read)
}

func (l *DiskLog) WriteToFile() error {
	file := fmt.Sprintf("%v.disklog", l.Name)

	fp, err := os.Create(file)
	if err != nil {
		return err
	}

	defer fp.Close() //nolint:errcheck

	for _, op := range l.Ops {
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
