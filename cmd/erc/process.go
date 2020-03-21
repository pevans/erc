package main

import (
	"time"

	"github.com/pevans/erc/pkg/boot"
	"github.com/pevans/erc/pkg/emu"

	"github.com/pkg/errors"
)

func processorLoop(comp emu.Computer, log *boot.Logger) {
	for {
		if err := comp.Process(); err != nil {
			log.Error(errors.Wrap(err, "main loop received error from processor"))
			return
		}

		time.Sleep(100 * time.Nanosecond)
	}
}
