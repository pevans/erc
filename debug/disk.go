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
	return
}
