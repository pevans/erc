package a2sym

var writeSwitchMap = map[int]Switch{
	// Keyboard switches -- see the note in the readSwitchMap on how
	// this switch was named.
	0xC010: {Mode: ModeW, Description: "clear strobe"},

	// Display switches
	0xC000: {Mode: ModeW, Name: "80STORE", Description: "off: cause PAGE2 on to select auxiliary ram"},
	0xC001: {Mode: ModeW, Name: "80STORE", Description: "on: allow PAGE2 to switch main RAM areas"},
	0xC00C: {Mode: ModeW, Name: "80COL", Description: "off: display 40 columns"},
	0xC00D: {Mode: ModeW, Name: "80COL", Description: "on: display 80 columns"},
	0xC00E: {Mode: ModeW, Name: "ALTCHAR", Description: "off: display text using primary character set"},
	0xC00F: {Mode: ModeW, Name: "ALTCHAR", Description: "on: display text using alternate character set"},
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
	0xC07E: {Mode: ModeW, Name: "IOUDIS", Description: "on: disable IOU access for addresses $C058 to $C05F; enable access to DHIRES switch"},
	0xC07F: {Mode: ModeW, Name: "IOUDIS", Description: "off: enable IOU access for addresses $C058 to $C05F; disable access to DHIRES switch"},

	// Bank memory
	0xC008: {Mode: ModeW, Name: "ALTZP", Description: "off: use main bank, page 0 and 1"},
	0xC009: {Mode: ModeW, Name: "ALTZP", Description: "on: use auxiliary bank, page 0 and 1"},

	// Auxiliary memory
	0xC002: {Mode: ModeW, Name: "RAMRD", Description: "read main memory"},
	0xC003: {Mode: ModeW, Name: "RAMRD", Description: "read auxiliary memory"},
	0xC004: {Mode: ModeW, Name: "RAMWRT", Description: "write main memory"},
	0xC005: {Mode: ModeW, Name: "RAMWRT", Description: "write auxiliary memory"},

	// I/O (PC) memory -- the ordering of these switch addresses is not
	// a typo
	0xC00A: {Mode: ModeW, Name: "SLOTC3ROM", Description: "internal rom at $C300"},
	0xC00B: {Mode: ModeW, Name: "SLOTC3ROM", Description: "slot rom at $C300"},
	0xC006: {Mode: ModeW, Name: "SLOTCXROM", Description: "slot rom at $Cx00"},
	0xC007: {Mode: ModeW, Name: "SLOTCXROM", Description: "internal rom at $Cx00"},
}

func WriteSoftSwitch(addr int) Switch {
	writeSwitch, ok := writeSwitchMap[addr]
	if !ok {
		return Switch{}
	}

	return writeSwitch
}
