package debug

import (
	"fmt"
	"os"

	"github.com/pevans/erc/a2"
)

func disk(comp *a2.Computer, tokens []string) {
	if len(tokens) < 2 {
		say("disk requires an argument that is a file to load")
		return
	}

	image := tokens[1]

	data, err := os.OpenFile(image, os.O_RDWR, 0644)
	if err != nil {
		say(fmt.Sprintf("couldn't open file %v: %v", image, err))
		return
	}

	if err := comp.Load(data, image); err != nil {
		say(fmt.Sprintf("couldn't load file: %v", err))
		return
	}

	say(fmt.Sprintf("loaded %v into drive", image))
}

// Toggle write protection on drive 1. This is _probably_ not the right
// interface to toggle write protection, but it's easy to implement for
// testing.
func writeProtect(comp *a2.Computer, tokens []string) {
	comp.Drive1.WriteProtect = !comp.Drive1.WriteProtect

	status := "OFF"
	if comp.Drive1.WriteProtect {
		status = "ON"
	}

	say(fmt.Sprintf("write protect on drive 1 is %v", status))
}
