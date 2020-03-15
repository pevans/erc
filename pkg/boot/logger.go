package boot

import (
	"log"
	"os"

	"github.com/pkg/errors"
)

type Logger struct {
	log   *log.Logger
	Level int
}

const (
	LogError = iota
	LogDebug
)

func (c *Config) LogLevel() int {
	if c.Log.Level == "debug" {
		return LogDebug
	}

	return LogError
}

func (c *Config) NewLogger() (*Logger, error) {
	l := new(Logger)
	l.Level = c.LogLevel()

	file, err := openLogFile(c.Log.File)
	if err != nil {
		return nil, errors.Wrapf(err, "could not open log file %s", c.Log.File)
	}

	l.log.SetOutput(file)

	return l, nil
}

func (l *Logger) UseOutput() {
	log.SetOutput(l.log.Writer())
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
