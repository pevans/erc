package boot

import (
	"log"
	"os"

	"github.com/pkg/errors"
)

type LogLevel int

type Logger struct {
	log   *log.Logger
	Level LogLevel
}

const (
	LogError LogLevel = iota
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

	file, err := openLogFile(c.Log.File)
	if err != nil {
		return nil, errors.Wrapf(err, "could not open log file %s", c.Log.File)
	}

	if file != nil {
		writer = file
	}

	l.log = log.New(writer, "", log.LstdFlags)

	return l, nil
}

// UseOutput will tell the main log system to use our writer. This
// function is intended to be scaffolding until we can convert to a
// system where your logger is always an instantiation and not a global
// variable.
func (l *Logger) UseOutput() {
	log.SetOutput(l.log.Writer())
}

// CanLog simply returns true if the log level configured by the logger
// would allow a given log level to be logged.
func (l *Logger) CanLog(lvl LogLevel) bool {
	return lvl <= l.Level
}

// Error will log an error message (if allowed).
func (l *Logger) Error(fmt string, vals ...interface{}) {
	if l.CanLog(LogError) {
		l.log.Printf("error: "+fmt, vals...)
	}
}

// Debug will log a debug message (if allowed).
func (l *Logger) Debug(fmt string, vals ...interface{}) {
	if l.CanLog(LogDebug) {
		l.log.Printf("debug: "+fmt, vals...)
	}
}

// openLogFile will attempt to open the given filename for writing log
// data out. If the filename is empty, it will not return an error, but
// will return a nil file value.
func openLogFile(fileName string) (*os.File, error) {
	if fileName == "" {
		return nil, nil
	}

	return os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
}
