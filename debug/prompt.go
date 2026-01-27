package debug

import (
	"github.com/peterh/liner"
	"github.com/pevans/erc/a2"
)

func Prompt(comp *a2.Computer, line *liner.State) {
	cmd, err := line.Prompt("debug> ")
	if err != nil {
		say("couldn't read input")
		return
	}

	line.AppendHistory(cmd)
	execute(comp, cmd)
}
