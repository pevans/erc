package a2sym

// readSwitchMap contains a set of memory addresses that correspond to
// soft switches defined in the Apple //e technical manual
var readSwitchMap = map[int]Switch{
	// Keyboard switches
	0xC000: {Mode: ModeR, Description: "keyboard data and strobe"},
	0xC010: {Mode: ModeR, Description: "any-key-down flag and clear-strobe switch"},

	// Display switches
	0xC018: {Mode: ModeR7, Name: "RD80STORE", Description: "1 = on"},
	0xC01A: {Mode: ModeR7, Name: "RDTEXT", Description: "1 = on"},
	0xC01B: {Mode: ModeR7, Name: "RDMIXED", Description: "1 = on"},
	0xC01C: {Mode: ModeR7, Name: "RDPAGE2", Description: "1 = on"},
	0xC01D: {Mode: ModeR7, Name: "RDHIRES", Description: "1 = on"},
	0xC01E: {Mode: ModeR7, Name: "RDALTCHAR", Description: "1 = on"},
	0xC01F: {Mode: ModeR7, Name: "RD80COL", Description: "1 = on"},
	0xC07E: {Mode: ModeR7, Name: "RDIOUDIS", Description: "1 = off"}, // "off" is not a typo
	0xC07F: {Mode: ModeR7, Name: "RDDHIRES", Description: "1 = on"},
	// These are duplicated in writeSwitchMap given we support both
	// reads and writes
	0xC050: {Mode: ModeRW, Name: "TEXT", Description: "off: display graphics or, if MIXED on, mixed"},
	0xC051: {Mode: ModeRW, Name: "TEXT", Description: "on: display text"},
	0xC052: {Mode: ModeRW, Name: "MIXED", Description: "off: display only text or only graphics"},
	0xC053: {Mode: ModeRW, Name: "MIXED", Description: "on: if TEXT off, display text and graphics"},
	0xC054: {Mode: ModeRW, Name: "PAGE2", Description: "off: select display page 1"},
	0xC055: {Mode: ModeRW, Name: "PAGE2", Description: "off: select display page 2 or, if 80STORE on, display page 1 in auxiliary memory"},
	0xC056: {Mode: ModeRW, Name: "HIRES", Description: "off: if TEXT off, display low-resolution graphics"},
	0xC057: {Mode: ModeRW, Name: "HIRES", Description: "on: if TEXT off, display high-resolution or, if DHIRES on, double-high-resolution graphics"},
	0xC05E: {Mode: ModeRW, Name: "DHIRES", Description: "on: if IOUDIS on, turn on double-high res."},
	0xC05F: {Mode: ModeRW, Name: "DHIRES", Description: "off: if IOUDIS on, turn off double-high res."},

	// Bank Select switches
	0xC011: {Mode: ModeR7, Name: "RDBNK2", Description: "read whether $D000 bank 2 (1) or bank 1 (0)"},
	0xC012: {Mode: ModeR7, Name: "RDLCRAM", Description: "read ram (1) or rom (0)"},
	0xC016: {Mode: ModeR7, Name: "RDALTZP", Description: "read whether auxiliary (1) or main (0) bank for page 0 and 1"},
	0xC080: {Mode: ModeR, Description: "read ram, write none, bank 2"},
	0xC081: {Mode: ModeRR, Description: "read rom, write ram, bank 2"},
	0xC082: {Mode: ModeR, Description: "read rom, write none, bank 2"},
	0xC083: {Mode: ModeRR, Description: "read ram, write ram, bank 2"},
	0xC084: {Mode: ModeR, Description: "read ram, write none, bank 2"},
	0xC085: {Mode: ModeRR, Description: "read rom, write ram, bank 2"},
	0xC086: {Mode: ModeR, Description: "read rom, write none, bank 2"},
	0xC087: {Mode: ModeRR, Description: "read ram, write ram, bank 2"},
	0xC088: {Mode: ModeR, Description: "read ram, write none, bank 1"},
	0xC089: {Mode: ModeRR, Description: "read rom, write ram, bank 1"},
	0xC08A: {Mode: ModeR, Description: "read rom, write none, bank 1"},
	0xC08B: {Mode: ModeRR, Description: "read ram, write ram, bank 1"},
	0xC08C: {Mode: ModeR, Description: "read ram, write none, bank 1"},
	0xC08D: {Mode: ModeRR, Description: "read rom, write ram, bank 1"},
	0xC08E: {Mode: ModeR, Description: "read rom, write none, bank 1"},
	0xC08F: {Mode: ModeRR, Description: "read ram, write ram, bank 1"},

	// Auxiliary memory
	0xC013: {Mode: ModeR, Name: "RAMRD"},
	0xC014: {Mode: ModeR, Name: "RAMWRT"},

	// I/O (PC) memory
	0xC015: {Mode: ModeR, Name: "SLOTCXROM"},
	0xC017: {Mode: ModeR, Name: "SLOTC3ROM"},
}

// ReadSwitch returns a Switch for a given address. If there is no known
// switch, then an empty Switch is returned (Mode == ModeNone).
func ReadSwitch(addr int) Switch {
	readSwitch, ok := readSwitchMap[addr]
	if !ok {
		return Switch{}
	}

	return readSwitch
}
