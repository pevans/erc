package a2drive

import "github.com/pevans/erc/a2/a2enc"

// phaseTransition will make use of a given drive phase to step the drive
// forward or backward by some number of tracks. This is always going to be -1
// (step backward); 1 (step forward); or 0 (no change).
func (d *Drive) phaseTransition(phase int) {
	if phase < 0 || phase > 4 {
		return
	}

	offset := phaseTable[(d.phase*5)+phase]

	// Because of the above check, we can assert that the formula we use for
	// the phase transition ((curPhase * 5) + phase) will match something, so
	// we step immediately.
	d.Step(offset)

	// We also have to update our current phase
	d.phase = phase
}

// trackLen returns the track length in bytes based on the loaded image type.
func (d *Drive) trackLen() int {
	if d.imageType == a2enc.Nibble {
		return a2enc.NibTrackLen
	}

	return a2enc.PhysTrackLen
}

// dataPosition returns the segment position that the drive is currently at,
// based upon track and sector position.
func (d *Drive) dataPosition() int {
	return (d.Track() * d.trackLen()) + d.sectorPos
}

// Sector returns the current sector that the drive head is positioned over.
// Reads and writes will occur in the returned sector.
func (d *Drive) Sector() int {
	return d.sectorPos / 0x1A0
}

// SectorPosition returns the raw offset from the beginning of the track at
// which the drive head is now positioned. This is always a number greater
// than or equal to zero, but less than the length of a physical track.
func (d *Drive) SectorPosition() int {
	return d.sectorPos
}

// Shift updates the sector position of the drive forward or backward by the
// given offset in bytes. Since tracks are circular and the disk is spinning,
// offsets that carry us beyond the bounds of the track instead bring us to
// the other end of the track.
func (d *Drive) Shift(offset int) {
	d.sectorPos += offset

	trackLen := d.trackLen()

	// In practice, these for loops are mutually exclusive; only one of them
	// would ever be entered.
	for d.sectorPos >= trackLen {
		d.sectorPos -= trackLen
	}

	for d.sectorPos < 0 {
		d.sectorPos += trackLen
	}

	d.diskShifted = true
}

// Step moves the track position forward or backward, depending on the sign of
// the offset. This simulates the stepper motor that moves the drive head
// further into the center of the disk platter (offset > 0) or further out
// (offset < 0).
func (d *Drive) Step(offset int) {
	d.trackPos += offset

	switch {
	case d.trackPos >= a2enc.MaxSteps:
		d.trackPos = a2enc.MaxSteps - 1
	case d.trackPos < 0:
		d.trackPos = 0
	}
}

// SwitchPhase will figure out what phase we should be moving to based on a
// given address.
func (d *Drive) SwitchPhase(addr int) {
	phase := -1

	switch addr & 0xf {
	case 0x1:
		phase = 1
	case 0x3:
		phase = 2
	case 0x5:
		phase = 3
	case 0x7:
		phase = 4
	}

	d.phaseTransition(phase)
}

// Track returns the current track in which the drive head is positioned.
// Reads and writes will occur within the provided track.
func (d *Drive) Track() int {
	return d.trackPos / 2
}
