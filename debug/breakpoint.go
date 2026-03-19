package debug

import (
	"fmt"
	"strconv"
	"strings"
)

// A set of addresses at which we should automatically enter the debugger.
var breakpoints = make(map[int]bool)

func AddBreakpoint(addr int) {
	breakpoints[addr] = true
}

func HasBreakpoint(addr int) bool {
	return breakpoints[addr]
}

// ParseBreakpoints parses a comma-separated list of hex addresses and adds
// each as a breakpoint. It returns an error if any address is invalid.
func ParseBreakpoints(flagVal string) error {
	for addrStr := range strings.SplitSeq(flagVal, ",") {
		addrStr = strings.TrimSpace(addrStr)
		if addrStr == "" {
			continue
		}
		addr, err := strconv.ParseUint(addrStr, 16, 16)
		if err != nil {
			return fmt.Errorf("invalid breakpoint address %q: %v", addrStr, err)
		}
		AddBreakpoint(int(addr))
	}
	return nil
}
