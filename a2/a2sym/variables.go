package a2sym

var variablesMap = map[int]string{
	// zero page variables
	0x0000: "LOC0",
	0x0001: "LOC1",
	0x0020: "WNDLFT",
	0x0021: "WNDWDTH",
	0x0022: "WNDTOP",
	0x0023: "WNDBTM",
	0x0024: "CH",
	0x0025: "CV",
	0x0026: "GBASL",
	0x0027: "GBASH",
	0x0028: "BASL",
	0x0029: "BASH",
	0x002A: "BAS2L",
	0x002B: "BAS2H",
	0x002C: "H2/LMNEM",
	0x002D: "V2/MMNEM",
	0x002E: "MASK/CHKSUM/FORMAT",
	0x002F: "LASTIN/LENGTH/SIGN",
	0x0030: "COLOR",
	0x0031: "MODE",
	0x0032: "INVFLG",
	0x0033: "PROMPT",
	0x0034: "YSAV",
	0x0035: "YSAV1",
	0x0036: "CSWL",
	0x0037: "CSWH",
	0x0038: "KSWL",
	0x0039: "KSWH",
	0x003A: "PCL",
	0x003B: "PCH",
	0x003C: "A1L",
	0x003D: "A1H",
	0x003E: "A2L",
	0x003F: "A2H",
	0x0040: "A3L",
	0x0041: "A3H",
	0x0042: "A4L",
	0x0043: "A4H",
	0x0044: "A5L",
	0x0045: "A5H/ACC",
	0x0046: "XREG",
	0x0047: "YREG",
	0x0048: "STATUS",
	0x0049: "SPNT",
	0x004E: "RNDL",
	0x004F: "RNDH",
	0x0095: "PICK",

	// out of the zero page
	0x0200: "IN",
	0x03F0: "BRKV",
	0x03F2: "SOFTEV",
	0x03F4: "PWREDUP",
	0x03F5: "AMPERV",
	0x03F8: "USRADR",
	0x03FB: "NMI",
	0x03FE: "IRQLOC",

	// page 1 memory
	0x0400: "LINE1",

	0x07FB: "MSLOT",

	0xE000: "BASIC",
	0xE003: "BASIC2",
}

func Variable(addr int) string {
	name, ok := variablesMap[addr]
	if !ok {
		return ""
	}

	return name
}
