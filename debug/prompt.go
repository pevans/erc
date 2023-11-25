package debug

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pevans/erc/a2"
)

func Prompt(comp *a2.Computer) {
	fmt.Printf("debug> ")

	cmd, err := read()
	if err != nil {
		say("couldn't read input")
		return
	}

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

func say(message string) {
	fmt.Printf(" --> %v\n", message)
}

func execute(comp *a2.Computer, cmd string) {
	tokens := strings.Fields(cmd)

	if len(tokens) == 0 {
		say("no command given")
		return
	}

	switch tokens[0] {
	case "help":
		help()

	case "resume":
		comp.Debugger = false
		say("resuming emulation")
	default:
		say(fmt.Sprintf(`unknown command: "%v"`, tokens[0]))
		help()
	}
}

func help() {
	say("list of commands")
	say("  help ..... print this message")
	say("  resume ... resume emulation")
}
