package debug

import (
	"fmt"

	"github.com/pevans/erc/a2"
)

func status(comp *a2.Computer) {
	var (
		regfmt  = "registers .......... A:$%02x X:%02x Y:%02x S:$%02x P:$%02x PC:$%04x"
		nextfmt = "next instruction ..."
	)

	say(fmt.Sprintf(regfmt, comp.CPU.A, comp.CPU.X, comp.CPU.Y, comp.CPU.S, comp.CPU.P, comp.CPU.PC))
	say(fmt.Sprintf(nextfmt))
}
