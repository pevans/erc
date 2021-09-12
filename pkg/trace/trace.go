package trace

import (
	"fmt"
	"io"
)

type Trace struct {
	// Where this trace is located when it is run (e.g. an address in
	// memory).
	Location string

	// Instruction is the string form of the instruction being executed
	Instruction string

	// Operand is the string form of the operand to the instruction.
	// This may include several operands, formatted in whatever way
	// makes sense given the context.
	Operand string

	// State is a collection of state variables that describe the
	// context of the machine as this trace was recorded
	State string

	// Counter is a counter for how many times we've recorded a trace
	Counter int
}

func (t *Trace) Write(w io.Writer) {
	fmt.Fprintf(w, "%-8s%-6s%-12s ; %s +%d\n",
		t.Location, t.Instruction,
		t.Operand, t.State, t.Counter,
	)
}
