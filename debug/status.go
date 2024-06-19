package debug

import (
	"fmt"

	"github.com/pevans/erc/a2"
)

func status(comp *a2.Computer) {
	var (
		regfmt  = "registers .......... %s"
		nextfmt = "last instruction ... %s"
	)

	say(fmt.Sprintf(regfmt, comp.CPU.Status()))
	say(fmt.Sprintf(nextfmt, comp.CPU.LastInstruction()))
}
