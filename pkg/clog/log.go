package clog

import (
	"fmt"
	"io"
	"os"
	"time"
)

const (
	LevelDebug = iota
	LevelInfo
	LevelError
)

var (
	Level    = LevelError
	stdlog   = NewChannel(os.Stdout)
	Shutdown = make(chan bool)
)

// Init starts up the package logger.
func Init(w io.Writer) {
	stdlog.Writer = w
	go stdlog.WriteLoop(Shutdown)
}

// Debug sends a message to the stdlog channel if the level is >= debug.
func Debugf(format string, vals ...interface{}) {
	if Level >= LevelDebug {
		mesg := fmt.Sprintf(format, vals...)
		stdlog.Printf("[%v] <debug> %v", time.Now(), mesg)
	}
}

func Debug(v interface{}) {
	Debugf("%v", v)
}

// Info sends a message to the stdlog channel if the level is >= info.
func Infof(format string, vals ...interface{}) {
	if Level >= LevelInfo {
		mesg := fmt.Sprintf(format, vals...)
		stdlog.Printf("[%v] <info> %v", time.Now(), mesg)
	}
}

func Info(v interface{}) {
	Infof("%v", v)
}

// Error will _always_ send a message to the stdlog channel. It will
// also panic, effectively shutting down process without some recovery
// attempt.
func Errorf(format string, vals ...interface{}) {
	mesg := fmt.Sprintf(format, vals...)
	stdlog.Printf("[%v] <error> %v", time.Now(), mesg)

	// Sending a true to the Shutdown channel may not be strictly
	// necessary, but we want to show the intention that we exit the
	// write loop.
	Shutdown <- true

	panic("error, exiting")
}

func Error(v interface{}) {
	Errorf("%v", v)
}
