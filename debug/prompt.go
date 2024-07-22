package debug

import (
	"bufio"
	"os"

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

func read() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)

	// Since we're just scanning one line, we don't really care what the
	// result of Scan is
	_ = scanner.Scan()

	line := scanner.Text()
	err := scanner.Err()

	return line, err
}
