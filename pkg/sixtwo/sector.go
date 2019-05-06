package sixtwo

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
