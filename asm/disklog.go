package asm

import (
	"fmt"
	"os"
	"time"
)

type DiskRead struct {
	// The time between reads
	Elapsed time.Duration

	// Where the read occurred on the disk
	HalfTrack int
	Sector    int

	// What byte was read
	Byte uint8

	// The full instruction and operand that caused the disk read
	Instruction string
}

type DiskLog struct {
	Reads []DiskRead
	Name  string
}

func NewDiskLog(name string) *DiskLog {
	log := new(DiskLog)
	log.Name = name

	return log
}

func (l *DiskLog) Add(read *DiskRead) {
	l.Reads = append(l.Reads, *read)
}

func (l *DiskLog) WriteToFile() error {
	file := fmt.Sprintf("%v.disklog", l.Name)

	fp, err := os.Create(file)
	if err != nil {
		return err
	}

	defer fp.Close() //nolint:errcheck

	for _, read := range l.Reads {
		logLine := fmt.Sprintf(
			"[%-10v] T:%02X S:%01X P:%04X B:$%02X | %v\n",
			read.Elapsed.Round(time.Millisecond), // time since boot
			read.HalfTrack>>1,                    // track
			read.Sector/0x1A0, read.Sector,       // sect and pos
			read.Byte, read.Instruction,
		)

		if _, err := fp.WriteString(logLine); err != nil {
			return err
		}
	}

	return nil
}
