package a2sym

var subroutineMap = map[int]string{
	0xC305: "BASICIN",
	0xC307: "BASICOUT",
	0xC312: "AUXMOVE",
	0xC314: "XFER",
	0xF800: "PLOT",
	0xF819: "HLINE",
	0xF828: "VLINE",
	0xF832: "CLRSCR",
	0xF836: "CLRTOP",
	0xF85F: "NEXTCOL",
	0xF864: "SETCOL",
	0xF871: "SCRN",
	0xF941: "PRNTAX",
	0xF948: "PRBLNK",
	0xF94A: "PRBL2",
	0xFB1E: "PREAD",
	0xFBDD: "BELL1",
	0xFC42: "CLREOP",
	0xFC58: "HOME",
	0xFC9C: "CLREOL",
	0xFC9E: "CLEOLZ",
	0xFCA8: "WAIT",
	0xFD0C: "RDKEY",
	0xFD1B: "KEYIN",
	0xFD35: "RDCHAR",
	0xFD67: "GETLNZ",
	0xFD6A: "GETLN",
	0xFD6F: "GETLN1",
	0xFD8B: "CROUT1",
	0xFD8E: "CROUT",
	0xFDDA: "PRBYTE",
	0xFDE3: "PRHEX",
	0xFDED: "COUT",
	0xFDF0: "COUT1",
	0xFE2C: "MOVE",
	0xFE36: "VERIFY",
	0xFE80: "SETINV",
	0xFE84: "SETNORM",
	0xFECD: "WRITE",
	0xFEFD: "READ",
	0xFF2D: "PRERR",
	0xFF3A: "BELL",
	0xFF3F: "IOREST",
	0xFF4A: "IOSAVE",
}

// Subroutine returns the name of a subroutine that Apple documented in
// their technical reference, if one exists, for any given address. If
// one does not exist, it will return an empty string.
func Subroutine(addr int) string {
	name, ok := subroutineMap[addr]
	if !ok {
		return ""
	}

	return name
}
