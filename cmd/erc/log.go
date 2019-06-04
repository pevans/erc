package main

import (
	"os"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// setLogging attempts to set the file name and level for logging within
// logrus, and returns an error if it was unable to do so.
func setLogging(fileName, levelName string) error {
	// No logging is necessary!
	if fileName == "" {
		return nil
	}

	// If a file was given, but no level, then assume they want error
	// logging
	if levelName == "" {
		levelName = "error"
	}

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		return errors.Wrapf(err, "could not open file %s for logging", fileName)
	}

	level, err := log.ParseLevel(levelName)
	if err != nil {
		return errors.Wrapf(err, "could not recognize level %s for logging", levelName)
	}

	log.SetOutput(file)
	log.SetLevel(level)

	return nil
}
