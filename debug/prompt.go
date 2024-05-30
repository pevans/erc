package debug

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/pevans/erc/a2"
	"github.com/pkg/errors"
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
	// data
	case "get":
		get(comp, tokens)
	case "reg":
		reg(comp, tokens)
	case "set":
		set(comp, tokens)
	case "statemap":
		stateMap(comp)
	case "status":
		status(comp)

		// execution
	case "step":
		step(comp, tokens)
	case "until":
		until(comp, tokens)

		// the rest
	case "help":
		help()
	case "quit":
		comp.Shutdown()
		os.Exit(0)
	case "resume":
		comp.Debugger = false
		say("resuming emulation")

	default:
		say(fmt.Sprintf("unknown command: \"%v\"", tokens[0]))
		help()
	}
}

func help() {
	say("list of commands")
	say("  [data]")
	say("    get <addr> ......... print the value at address <addr>")
	say("    reg <r> <val> ...... write <val> to register <r>")
	say("    set <addr> <val> ... write <val> at address <addr>")
	say("    status ............. show registers and next execution")
	say("  [execution]")
	say("    step <times> ....... execute <times> instructions")
	say("  [the rest]")
	say("    help ............... print this message")
	say("    quit ............... quit the emulator")
	say("    resume ............. resume emulation")
}

func hex(token string, bits int) (int, error) {
	ui64, err := strconv.ParseUint(token, 16, bits)
	if err != nil {
		return 0, errors.Wrapf(err, "invalid hex: \"%v\"", token)
	}

	return int(ui64), nil
}

func integer(token string) (int, error) {
	ui64, err := strconv.ParseUint(token, 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "invalid integer: \"%v\"", token)
	}

	return int(ui64), nil
}
