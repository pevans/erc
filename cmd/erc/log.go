package main

import (
	"os"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// logLevel will return the logrus logging level that corresponds with a
// given error level. If that level isn't a "real" one, it will return a
// default of the ErrorLevel value, as opposed to passing the error back
// up.
func logLevel(lvl string) log.Level {
	level, err := log.ParseLevel(lvl)
	if err != nil {
		return log.ErrorLevel
	}

	return level
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

// setLogging attempts to set the file name and level for logging within
// logrus, and returns an error if it was unable to do so.
func setLogging(fileName, levelName string) error {
	file, err := openLogFile(fileName)
	if err != nil {
		return errors.Wrapf(err, "could not open file %s for logging", fileName)
	}

	if file != nil {
		log.SetOutput(file)
	}

	log.SetLevel(logLevel(levelName))

	return nil
}
