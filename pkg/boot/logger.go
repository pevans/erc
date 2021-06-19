package boot

import (
	"log"
	"os"

	"github.com/pkg/errors"
)

// A LogLevel is some level at which, if a message is at or exceeds its
// threshold, we will log the message being passed to us.
type LogLevel int

// A Logger is a device that can log messages with respect to a given log level.
type Logger struct {
	logger *log.Logger
	Level  LogLevel
}

const (
	// LogNothing essentially is the level at which nothing will be logged out.
	LogNothing LogLevel = iota

	// LogError is the main error type you'd use when something is effectively
	// "wrong" with execution.
	LogError

	// LogDebug is the error type you'd use for informational/development work.
	LogDebug
)

// LogLevel will return the LogDebug level if the level is specifically
// "debug"; otherwise it will always return LogError.
func (c *Config) LogLevel() LogLevel {
	if c.Log.Level == "debug" {
		return LogDebug
	}

	return LogError
}

// NewLogger will create a new logger from a configuration.
func (c *Config) NewLogger() (*Logger, error) {
	var (
		writer = os.Stdout
		l      = new(Logger)
	)

	l.Level = c.LogLevel()

	file, err := os.OpenFile(c.Log.File, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil && c.Log.File != "" {
		return nil, errors.Wrapf(err, "could not open log file %s", c.Log.File)
	}

	if file != nil {
		writer = file
	}

	l.logger = log.New(writer, "", log.LstdFlags)

	return l, nil
}

// UseOutput will tell the main log system to use our writer. This
// function is intended to be scaffolding until we can convert to a
// system where your logger is always an instantiation and not a global
// variable.
func (l *Logger) UseOutput() {
	log.SetOutput(l.logger.Writer())
}

// CanLog simply returns true if the log level configured by the logger
// would allow a given log level to be logged.
func (l *Logger) CanLog(lvl LogLevel) bool {
	return lvl <= l.Level
}

// Errorf will log an error message (if allowed).
func (l *Logger) Errorf(fmt string, vals ...interface{}) {
	if l.CanLog(LogError) {
		l.logger.Printf("error: "+fmt, vals...)
	}
}

// Error will log an error message (if allowed).
func (l *Logger) Error(vals ...interface{}) {
	if l.CanLog(LogError) {
		l.logger.Println(vals...)
	}
}

// Fatal will do what Error would do, with the added step of it'll exit the
// program afterward.
func (l *Logger) Fatal(vals ...interface{}) {
	l.Error(vals...)
	os.Exit(1)
}

// Debugf will log a debug message (if allowed).
func (l *Logger) Debugf(fmt string, vals ...interface{}) {
	if l.CanLog(LogDebug) {
		l.logger.Printf(fmt, vals...)
	}
}

// Debug will log a debug message (if allowed).
func (l *Logger) Debug(vals ...interface{}) {
	if l.CanLog(LogDebug) {
        l.logger.Println(vals...)
	}
}
