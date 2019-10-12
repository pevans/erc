package main

import (
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func processorLoop() {
	for {
		if err := emulator.Processor.Process(); err != nil {
			log.Error(errors.Wrap(err, "main loop received error from processor"))
			return
		}

		time.Sleep(100 * time.Nanosecond)
	}
}
