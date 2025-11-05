package debug

import (
	"fmt"
	"strings"

	"github.com/pevans/erc/a2"
)

func status(comp *a2.Computer) {
	var (
		regfmt  = "registers .......... %s"
		lastfmt = "last instruction ... %s"
		nextfmt = "next instruction ... %s"
	)

	say(fmt.Sprintf(regfmt, strings.TrimSpace(comp.CPU.Status())))
	say(fmt.Sprintf(lastfmt, strings.TrimSpace(comp.CPU.LastInstruction())))
	say(fmt.Sprintf(nextfmt, strings.TrimSpace(comp.CPU.NextInstruction())))
}
