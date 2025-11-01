package debug

// A set of addresses at which we should automatically enter the debugger.
var breakpoints = make(map[int]bool)

func AddBreakpoint(addr int) {
	breakpoints[addr] = true
}

func HasBreakpoint(addr int) bool {
	return breakpoints[addr]
}
