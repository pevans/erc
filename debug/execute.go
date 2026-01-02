package debug

import (
	"fmt"
	"os"
	"strings"

	"github.com/pevans/erc/a2"
	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/gfx"
	"github.com/pevans/erc/obj"
)

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
	case "state":
		stateMap(comp)
	case "status":
		status(comp)

		// execution
	case "step":
		step(comp, tokens)
	case "until":
		until(comp, tokens)

		// simulation
	case "keypress":
		keypress(comp, tokens)

		// the rest
	case "disk":
		disk(comp, tokens)
	case "writeprotect":
		writeProtect(comp, tokens)
	case "help":
		help()
	case "quit":
		err := comp.Shutdown()
		if err != nil {
			say(fmt.Sprintf("shutdown was not clean: %v", err))
		}
		os.Exit(0)
	case "resume":
		comp.State.SetBool(a2state.Debugger, false)
		gfx.ShowStatus(obj.ResumePNG())
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
	say("    state .............. print the apple II state")
	say("    status ............. show registers and next execution")
	say("  [execution]")
	say("    step <times> ....... execute <times> instructions")
	say("    until <instruction>  execute until <instruction>")
	say("  [simulation]")
	say("    keypress <val> ..... simulate a keypress with hex ascii code <val>")
	say("  [the rest]")
	say("    disk <file> ........ load <file> into drive")
	say("    writeprotect ....... toggle write protect on drive 1")
	say("    help ............... print this message")
	say("    quit ............... quit the emulator")
	say("    resume ............. resume emulation")
}
