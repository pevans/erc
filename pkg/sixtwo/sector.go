package sixtwo

// This is the sector table for DOS 3.3.
var dosSectorTable = []int{
	0x0, 0x7, 0xe, 0x6, 0xd, 0x5, 0xc, 0x4,
	0xb, 0x3, 0xa, 0x2, 0x9, 0x1, 0x8, 0xf,
}

// This is the sector table for ProDOS.
var proSectorTable = []int{
	0x0, 0x8, 0x1, 0x9, 0x2, 0xa, 0x3, 0xb,
	0x4, 0xc, 0x5, 0xd, 0x6, 0xe, 0x7, 0xf,
}

// logicalSector returns the logical sector number, given the current
// image type and a physical sector number (sect).
func logicalSector(imageType, sect int) int {
	if sect < 0 || sect > 15 {
		return 0
	}

	switch imageType {
	case DOS33:
		return dosSectorTable[sect]

	case ProDOS:
		return proSectorTable[sect]
	}

	// Note: logical nibble sectors are the same as the "physical"
	// sectors.
	return sect
}
