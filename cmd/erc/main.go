package main

import (
	"log"

	"github.com/pevans/erc/pkg/mach/a2"
)

func main() {
	emu := a2.NewEmulator()
	err := emu.Booter.Boot()

	if err != nil {
		log.Fatal(err)
	}
}
